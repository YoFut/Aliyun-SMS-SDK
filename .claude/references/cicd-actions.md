# Gitea Actions CI/CD 基础设施规范（Runner + 物理机 Clash 代理）

> 触发时机：编写/修改 `.gitea/workflows/*.yml`、调整 runner 容器、CI 出现网络超时/TLS 握手失败时必读。
> 关联：执行手册 §3.1 T0.7（CI 接入）、§4.5（CI 门禁）；本机 Gitea 实例 `http://localhost:11451`。
> **关键事实**：`actions/*` 的拉取由 **act_runner 进程**执行（不在 job 容器内），因此代理必须配在 **runner 容器 env**，workflow `env:` 管不到 action 拉取——这是 TLS 超时的根因（§3）。

---

## 1. Runner 部署形态（3 个常驻 · linux x64）

单机 WSL2 Docker，3 个 act_runner 容器常驻、实例级注册，服务全实例所有仓库（含 `Graduation-Project/YoFut` 与个人测试仓）。

| 项 | 值 |
|---|---|
| 镜像 | `docker.gitea.com/act_runner:0.2.12`（amd64/linux，pin 版本勿用 `latest`） |
| 容器 | `gitea-runner` / `gitea-runner-2` / `gitea-runner-3`，`--restart unless-stopped` |
| 网络 | `pgsql_default`（与 `gitea` / `pgsql` 同网） |
| 实例地址 | `GITEA_INSTANCE_URL=http://gitea:3000`（**容器内网名 + 3000**，非 `localhost:11451`） |
| 注册级别 | 实例级（admin registration token），runner 各自独立 id |
| 名称 | `YoFut-runner-linux-x64[-2/-3]`，labels=`ubuntu-latest/22.04/20.04` |
| 挂载 | `/var/run/docker.sock`（生成 job 容器）+ `/home/wave/Downloads/temp/gitea/runner[-2/-3]:/data`（`.runner` 态 + `config.yaml` + `actcache/`） |
| 重启策略 | `unless-stopped`（`.runner` 态在卷里持久，重建容器不丢注册） |

**config.yaml 关键项**（每个 runner 各持一份，路径 `/data/config.yaml`）：

```yaml
container:
  network: "pgsql_default"   # job 容器加入同网，才能解析 gitea 主机名、checkout 代码
  privileged: false
cache:
  enabled: true
  dir: "/data/actcache"      # 缓存落卷，容器重建不丢
```

> 若 `container.network` 留空，job 容器走默认 bridge，无法解析 `gitea` → checkout 失败。

**重建一个 runner**（`.runner` 态持久，无需重新注册；仅首次或换实例时才需 registration token）：

```bash
docker rm -f gitea-runner-2
docker run -d --name gitea-runner-2 --restart unless-stopped --network pgsql_default \
  -e GITEA_INSTANCE_URL=http://gitea:3000 \
  -e GITEA_RUNNER_REGISTRATION_TOKEN=<实例级 token> \
  -e GITEA_RUNNER_NAME=YoFut-runner-linux-x64-2 \
  -e CONFIG_FILE=/data/config.yaml \
  -e HTTP_PROXY=http://172.26.160.1:7897 -e HTTPS_PROXY=http://172.26.160.1:7897 \
  -e http_proxy=http://172.26.160.1:7897 -e https_proxy=http://172.26.160.1:7897 \
  -e NO_PROXY=localhost,127.0.0.1,gitea,pgsql,YoFut-redis,172.23.0.0/16,172.26.0.0/16 \
  -e no_proxy=localhost,127.0.0.1,gitea,pgsql,YoFut-redis,172.23.0.0/16,172.26.0.0/16 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /home/wave/Downloads/temp/gitea/runner-2:/data \
  docker.gitea.com/act_runner:0.2.12
```

**实例级 registration token 取法**（需 Gitea admin token，见 memory `gitea-credential`）：

```bash
curl -X POST -H "Authorization: token <admin_token>" \
  http://localhost:11451/api/v1/admin/actions/runners/registration-token
```

**验证 runner 在线**：`docker logs gitea-runner-2 | grep "declare successfully"` 或 API `GET /api/v1/admin/actions/runners`（`status=online`）。

---

## 2. 物理机 Clash 代理（三层配置）

代理走物理机 Windows 上的 Clash（已开 **Allow LAN**）。job 容器在 `pgsql_default` 网络经 WSL2 NAT 网关转发到 Windows。

**代理地址**：`http://172.26.160.1:7897`

> ⚠️ `172.26.160.1` 是 **WSL2 NAT 网关**（= Windows 主机），WSL/Windows 重启后**可能变**。取当前值：
> ```bash
> wsl -e bash -c 'ip route show default | awk "{print \$3}"'
> ```
> 若变了，需更新 3 个 runner 容器 env 与 workflow `env:` 中的代理地址。

代理必须配在**三层**，缺一不可：

| 层 | 位置 | 作用对象 | 解决的问题 |
|---|---|---|---|
| ① runner 容器 env | `docker run -e HTTP_PROXY=...` | act_runner 进程的 git | **action 拉取**（`actions/checkout@v4` 等，从 github.com/gitea.com clone） |
| ② Gitea 服务端 | `app.ini [actions] DEFAULT_ACTIONS_URL = https://gitea.com` | action 解析源 | 把 `actions/*` 从 github.com 引到 gitea.com 镜像（国内可达） |
| ③ workflow `env:` | `.gitea/workflows/*.yml` job 级 env | job 容器内步骤（apt / go modules / curl） | job 内访问外网 |

**NO_PROXY 必须包含内网**：`localhost,127.0.0.1,gitea,pgsql,YoFut-redis,172.23.0.0/16,172.26.0.0/16`，否则 gitea/pgsql/redis 流量被劫持到代理 → 内网调用失败。

> ② 与 ① 互为冗余：单配 ①（runner 走代理访问 github.com）已可工作；② 进一步把 action 源换到 gitea.com，更稳。两者都开最稳。

---

## 3. TLS handshake timeout 根因与规避

**现象**：workflow 在 `拉取代码`（`actions/checkout@v4`）步骤之前失败，job 日志末尾：

```
☁  git clone 'https://github.com/actions/checkout' # ref=v4
cloning https://github.com/actions/checkout to /root/.cache/act/...
Get "https://github.com/actions/checkout/info/refs?service=git-upload-pack": net/http: TLS handshake timeout
🏁  Job failed
```

**根因**：action 仓库的 clone 由 **act_runner 进程**（gitea-runner 容器内）执行，缓存到 `/root/.cache/act/`，**不在 job 容器内**。因此 workflow 的 `env:` 代理（只注入 job 容器）**管不到 action 拉取**。act_runner 直连 github.com → 国内 TLS 握手超时（时通时断，表现为偶发失败或 7 分钟慢成功）。

**规避**（按优先级）：

1. **配 runner 容器代理 env**（§2 ①）——act_runner 的 git 经 Clash 访问 github.com，clone 从超时/7min 降到 ~7s。
2. **设 `DEFAULT_ACTIONS_URL=https://gitea.com`**（§2 ②）——action 源切到 gitea.com（国内可达，http 303），绕开 github.com。
3. **job 容器代理**（§2 ③）——仅影响 job 内步骤，不影响 action 拉取，但 go modules / apt 等必需。

> 排障顺序：先看 `docker logs gitea-runner` 与 job 日志的 `git clone` 行；若超时，确认 runner 容器有 `HTTP_PROXY` env（`docker exec gitea-runner env | grep -i proxy`）。

---

## 4. workflow yml 代理模板

`.gitea/workflows/ci.yml` 在 job 级注入代理 env（对应 §2 ③ 层）：

```yaml
name: CI
on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    env:
      # 物理机 Clash 代理（job 容器经 pgsql_default 网关转发到 Windows）
      HTTP_PROXY: http://172.26.160.1:7897
      HTTPS_PROXY: http://172.26.160.1:7897
      http_proxy: http://172.26.160.1:7897
      https_proxy: http://172.26.160.1:7897
      # 内网不走代理：gitea/pgsql/redis 与 docker/WSL 子网
      NO_PROXY: localhost,127.0.0.1,gitea,pgsql,YoFut-redis,172.23.0.0/16,172.26.0.0/16
      no_proxy: localhost,127.0.0.1,gitea,pgsql,YoFut-redis,172.23.0.0/16,172.26.0.0/16
    steps:
      - uses: actions/checkout@v4
      - name: go build & test
        run: |
          go build ./...
          go vet ./...
          go test ./... -short
```

> `runs-on: ubuntu-latest` 对应 runner labels（`ubuntu-latest` → `docker.gitea.com/runner-images:ubuntu-latest`，job 容器自动加入 `pgsql_default`）。
> Go 模块代理可选加 `GOPROXY=https://goproxy.cn,direct`，与 Clash 代理叠加更稳。

---

## 5. 维护与排障速查

| 症状 | 排查 |
|---|---|
| run 一直 `queued`、无人拾取 | runner 是否 `online`（API 或 `docker logs`）；实例级 runner 才服务个人仓，org 级只服务该 org |
| job 日志 `git clone ... TLS handshake timeout` | runner 容器缺代理 env（§2 ①）；或设 `DEFAULT_ACTIONS_URL=gitea.com`（§2 ②） |
| job 内 `curl gitea` 超时 | `NO_PROXY` 漏了 `gitea`，内网流量被劫持到代理 |
| WSL 重启后代理失效 | 网关 IP 变了，`ip route show default` 取新值，更新 runner env + workflow env |
| `.runner` 态丢失导致重注册 | 重建容器前确认 `/data` 卷挂载正确；重注册会生成新 runner id，旧 id 需在 Gitea 后台删除 |
| 想节约资源改按需 runner | 见执行手册外「watchdog + ephemeral」方案；常驻 3 个的单机开销可忽略（act_runner 空闲 ~30MB） |

**同步本 skill 文档**（Windows 无符号链接，`.cursor` 与 `.claude` 是物理副本；**目标已存在时 `Copy-Item` 会嵌套创建 `YoFut/YoFut/`，必须先删目标**）：

```powershell
cd C:\Users\diwei\OneDrive\文档\NFZY\毕设\repos\YoFut
Remove-Item .claude/skills/YoFut -Recurse -Force
Copy-Item .cursor/skills/YoFut .claude/skills/YoFut -Recurse -Force
```
