# Git 工作流规范（WaveYo 风格）

> 触发时机：创建分支、编写 commit、创建 PR、执行 Code Review 时必读。
> 来源：YoMirrorSite 全流程会话（5 个 PR 完整生命周期）+ YoFut 项目适配。
> **本文件覆盖 goddd `git-workflow` skill 的对应部分**（差异见 §8）。

---

## 1. Commit 规范

### 1.1 格式

```
<type>(<scope>): <中文描述>
```

本项目全部 Git User 应当使用 WaveYo <waveyo@waveyo.cn>；
禁止添加任何额外“共同撰写人”

Gitea Token 可从全局系统环境变量 `GITEA_TOKEN` 获取直接加载

### 1.2 Type 枚举

| Type | 场景 | 示例 |
|---|---|---|
| `feat` | 新功能 | `feat(upload): 实现 STS 前缀锁定凭证签发` |
| `fix` | 缺陷 / 安全修复 | `fix(security): 移除对象存储调试日志中的凭证` |
| `refactor` | 重构（无功能变更） | `refactor(model): 将 FileInfo 下沉到 model 层` |
| `docs` | 文档变更 | `docs(readme): 补充 OSS 配置说明` |
| `chore` | 清理 / 构建 / 工具 / 格式 | `chore(cleanup): 清理临时文件并加固 .gitignore` |
| `test` | 测试相关 | `test(handler): 补充上传接口单元测试` |
| `perf` | 性能优化 | `perf(cache): 优化 LRU 淘汰策略` |

> 与 goddd 默认的差异：**新增 `docs`**（毕设文档量大，需独立类型）。

### 1.3 Scope 规范

Scope 按**领域 / 模块**命名，YoFut 的 scope 清单：

```
identity    — 用户 / 企业 / RBAC
upload      — 上传与五道闸门
mirror      — 文件列表 / 搜索 / 预签名下载
catalog     — 软件 / 版本 / 资产
syncjob     — GitHub 同步 / 调度 / 对账
forum       — 版块 / 主题 / 帖子
comment     — 评论
moderation  — 统一审核
audit       — 审计日志
storage     — 对象存储适配（OSS）
web         — 前端
config      — 配置管理
router      — 路由注册
ci          — CI/CD 配置
security    — 安全相关（多模块通用）
defects     — 缺陷修复（多模块通用）
cleanup     — 代码清理（多模块通用）
build       — 构建相关
delivery    — 部署交付
```

### 1.4 描述规则

- **使用中文**，简洁说明变更内容
- 多项同类修复用**顿号**分隔：
  `fix(defects): 修复静默错误吞掉、panic 改 error、goroutine recover`
- 一个 commit 对应一个逻辑单元，**不混合**不同性质的变更
- 不加句号，首字母不强求小写（中文描述直接书写）

### 1.5 破坏性变更

在 type 后加 `!`，footer 写迁移说明：

```
feat(auth)!: 变更令牌载荷结构

BREAKING CHANGE: JWT 载荷字段 uid 更名为 user_id，客户端需同步调整解析逻辑。
```

---

## 2. 原子提交拆分原则

当单个 PR 包含多类修复时，按严重度拆分为独立 commit：

```
commit 1: fix(security): ...          ← 最高优先级（安全漏洞）
commit 2: fix(defects): ...           ← 缺陷修复（运行时问题）
commit 3: chore(cleanup): ...         ← 风格 / 规范清理
```

**优先级排序**：`安全 > 缺陷 > 重构 > 风格 > 文档 > 清理`

**拆分粒度判断**：

| 情形 | 处理 |
|---|---|
| 修改同一函数的不同行 | 合并在一个 commit |
| 修改不同模块 / 不同性质的代码 | 拆分为多个 commit |
| 重构与行为修复混在一起 | **必须拆开**（重构不得改变语义） |

多个 commit 共享同一 PR，由人类通过 **squash merge** 合并为一条干净历史。

---

## 3. PR 标准

### 3.1 分支命名

```
feature/<描述>     — 功能开发 / 审查修复 / 清理
fix/<描述>         — 紧急修复
```

描述用 **kebab-case**。示例：

- `feature/upload-gate`
- `feature/identity-rbac`
- `feature/comment-system`
- `fix/go-compile-errors`

> 禁止个人化命名（`test123`、`tmp`、`johns-branch`）。

### 3.2 PR 标题

遵循 commit 规范格式，概括 PR 目的：

```
feat(upload): 实现上传权限五道闸门 — 路由隔离 + 中间件 + 租户断言 + 凭证前缀锁定
fix(backend): 代码审查修复 — 安全漏洞 + 运行时缺陷 + 配置加固
```

### 3.3 PR 正文模板

````markdown
## 变更概要
<一句话说明 PR 目的和变更范围>

## 修复清单

### 🔴 安全漏洞
| # | 文件:行 | 问题 | 修复 |
|---|---|---|---|
| S1 | `path/file.go:42` | <问题描述> | <修复说明> |

### 🟠 缺陷修复
| # | 文件:行 | 问题 | 修复 |
|---|---|---|---|

### 🟡 风格/规范
| # | 文件 | 问题 | 修复 |
|---|---|---|---|

## 测试验证
- [ ] `go build ./...` 通过
- [ ] `go vet ./...` 无警告
- [ ] 单元测试通过（`go test ./... -short`）
- [ ] 集成测试通过（`go test ./... -tags=integration`）
- [ ] 前端构建通过（如涉及）

## 风险说明
- 破坏性变更：<无 / 列出具体变更>
- 依赖变更：<无 / 列出新增或升级的依赖>
- 演进触发条件（如引入新组件）：<压测数据 / 错误日志 / 实例数证据>

## 审核者
@WaveYo
````

### 3.4 提交数量建议

| 变更规模 | 建议 commit 数 | 示例 |
|---|---|---|
| 单文件小修 | 1 | 修正一个配置项 |
| 多文件同类别 | 1–3 | 模块级重构 |
| 跨模块多类别 | 3–5 | 安全 + 缺陷 + 风格综合修复 |
| 大型功能 | 5+ | 按功能拆分为多个 PR |

---

## 4. PR Review Checklist

人类审核 PR 时逐项确认：

```
[ ] CI 全部通过（build + vet + unit + integration + frontend）
[ ] PR 正文包含结构化审查报告（Security / Defect / Style 三级分类）
[ ] 每个问题标注了文件路径和行号
[ ] 修复方案与问题一一对应
[ ] 无破坏性变更（或已明确列出并说明影响范围）
[ ] 无依赖变更（或已明确列出变更原因）
[ ] Commit 信息符合规范（<type>(<scope>): <中文描述>）
[ ] 分支命名符合 feature/xxx 或 fix/xxx
[ ] 无敏感信息泄露（凭证 / 密钥 / 内网地址）
[ ] 命名符合 Entity-First / Noun-First
[ ] 注释全中文，新文件有行头，公开函数有 GoDoc
[ ] 日志级别正确，未输出 password/token/secret/key
[ ] 未提前引入 Redis / 分布式组件（最小原型先行）
```

---

## 5. 合并后清理步骤

PR 被 squash merge 后执行：

```bash
# 1. 切换到 dev 并拉取最新
git checkout dev && git pull gitea dev

# 2. 删除本地 feature 分支
git branch -D feature/<分支名>

# 3. 删除远端 feature 分支
git push gitea --delete feature/<分支名>

# 4. 清理本地跟踪引用
git remote prune gitea
```

> 本项目远端名为 `gitea`（`http://localhost:11451/Graduation-Project/YoFut.git`），非 `origin`。

---

## 6. 集成与发布原则

功能 PR 优先 squash 合入 `dev`；`dev → main` 通过人工审核的发布 PR 保留合并节点，便于版本流程读取完整语义变更：

```
  feature/* → dev → main → 自动 tag → 人工 Release
```

`main → dev` 同步不创建正式 tag。语义版本和发布操作见 `docs/version-release.md`；版本流程不依赖固定的 merge 标题格式。

---

## 7. 反模式（避免的行为）

- ❌ 直接在 `main` 上 commit 或 push
- ❌ 将 `dev → main` 的语义历史压成无语义标题且不设置 semver 标签
- ❌ 一个 PR 包含不相关的修复（混合安全漏洞修复与功能开发）
- ❌ PR 正文只有 "fixes #123" 没有结构化说明
- ❌ Commit 描述为 "fix bug" / "update code"（不具体）
- ❌ 分支名使用个人化命名
- ❌ 测试失败后修改断言以求「通过」

---

## 8. 与 goddd `git-workflow` 的差异裁决

| 冲突点 | goddd 默认 | **YoFut 采用**                      |
|---|---|-------------------------------------|
| 主线分支 | `main` + `dev` 双线 | **`dev` 集成、`main` 发布**           |
| 个人分支命名 | `<initials>/<feature>` | **`feature/<描述>` / `fix/<描述>`** |
| 合并方式 | `dev → main` 用 `--no-ff` | **功能 PR squash 到 dev；dev → main 保留合并节点** |
| commit type | 6 类（无 `docs`） | **7 类**（含 `docs`）               |
| subject 语言 | 跟随交互语言 | **固定中文**                        |
| 自动提交 | 完成后自动 commit、不询问 | **须走 PR；AI 不得直推 `main`**     |
| 分支名风格 | kebab-case | `feature/` 后用 kebab-case（一致）  |

**仍沿用 goddd 的部分**：

- rebase 变基保持线性；rebase 后合并
- Changelog 生成方式（`git log $(git describe --tags --abbrev=0)..HEAD --pretty=format:"%s" | sort`）
- 「一个 commit 只做一件事」
- 破坏性变更用 `!` 标注 + `BREAKING CHANGE` footer
- 版本号用 Git tag `v<major>.<minor>.<patch>`
