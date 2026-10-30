# 最小开发闭环

> 新需求先验证最小业务语义，直接复用已有真实 PostgreSQL；不是先搭建新基础设施。
> 2026-10-27 用户确认：开发验证与规范改用 PostgreSQL，保留隔离单测的 SQLite。

## 1. 开发基线

```text
需求 → ① 已有 PostgreSQL + 简单 HTTP / 进程内缓存
            ↓ 验证领域模型、事务与接口契约
       ② 按任务需要接入已有 Redis / OSS 等真实依赖
            ↓ 验证并发、一致性和实际错误
       ③ 有证据时才增加分布式能力
```

SQLite 仅用于无需外部服务的隔离单测、故障注入和兼容性回归，不是开发端到端验证的替代品。不能因 PostgreSQL 连接失败而回退 SQLite，再报告开发验证成功；也不要求把所有隔离单测迁移为外部数据库测试。

## 2. 配置来源

- 顺序与启动入口一致：默认值 → `configs/config.toml` → `configs/config.local.toml`，后者覆盖同名字段。
- 实际 PostgreSQL 连接及连接池读取本地文件的 `[Data.Database]`；`[Data.Redis]` 继续由既有 Redis 装配使用，数据库验证不因此强制新增服务。
- 本地文件必须存在且数据库配置须为 PostgreSQL；开发验证配置缺失、非法或连接失败时明确失败，不使用硬编码凭证、不另设一份开发 DSN。
- `config.local.toml` 不入库。不得输出文件全文、DSN、密码、Token、AK/SK或带配置内容的解析错误。

## 3. 最小验证入口

在 YoFut 根目录使用 PowerShell：

```powershell
$env:YoFut_DEVELOPMENT_POSTGRES = '1'
go test ./internal/web/api -tags=integration -run '^TestDevelopmentPostgres$' -count=1 -v
```

此入口读取本地 `[Data.Database]`，在真实 PostgreSQL 的随机独立 schema 执行相同 HTTP 断言：审核/申诉、管理/禁用、审计 CSV/JSONL。先核对 `current_schema()` 再创建业务表，结束后关闭验证连接、精确删除本轮 schema 并确认不存在；不迁移或清空业务 schema。不需要 OSS 凭证或额外 Redis 服务。

新领域仍按 goddd 生成模型及骨架，再在同类 PostgreSQL 隔离环境验证其最小业务闭环。该入口是现有模块的可执行范例，不代表未来领域无需自己的断言。

隔离回归命令保留：

```text
go build ./...
go vet ./...
go test ./... -short
```

短测试和 CI 可以保留 SQLite。它们与真实 PostgreSQL 开发验证互补，不能用“短测试通过”代替 PG 行锁、时钟、序列、类型和事务语义验证。

## 4. 不提前增加的能力

| 当前最小需求 | 首选实现 | 升级条件 |
|---|---|---|
| 持久化及事务 | 已配置的单机 PostgreSQL | 已有连接直接复用，不另建数据库服务 |
| 简单缓存 | 进程内缓存或已有领域缓存 | 有实际瓶颈才增加新缓存层 |
| 简单副作用 | 直接调用已有端口 | 有持久化/重试需求才引入队列 |
| 单实例互斥 | 数据库行锁或已有同步机制 | 多实例竞争有证据才增加分布式锁 |
| 部署 | `go run . -conf ./configs` | 任务明确要求才增加容器编排或多实例 |

不修改 goddd 框架、不手改 Wire 生成物，不以切换数据库为由提前建设分库、主从、Redis Cluster或消息系统。

## 5. 交付要求

- 写清最小业务语义、配置来源和真实 PG 命令。
- 记录实际通过/失败、驱动类别、隔离清理结果；失败不能通过删断言、加无理由 skip 或换 SQLite 掩盖。
- 数据库时钟、序列、隔离级别或 SQL 方言差异必须调查根因，行为修复与验证重构分别提交。
- 历史 SQLite 验收记录保留原事实，新增 PG 验证记录不得冒充过去已经执行。
- 说明未验证的 Redis/OSS、生产迁移和容量范围；不把单库 HTTP 验证写成完整生产验收。
