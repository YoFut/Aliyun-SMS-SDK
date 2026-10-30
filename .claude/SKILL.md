---
name: YoFut
description: >
  YoFut 项目通用开发规范，是 goddd skill 之上的**项目增量层**。
  当在 YoFut 仓库内编写代码、命名标识符、撰写注释、打日志、创建分支、
  编写 commit 或 PR、执行 Code Review、与 AI 协作交付时使用此技能。
  也应在以下隐含场景主动触发：新增文件/函数/类型、决定包名或表名、
  选择日志级别、判断是否该引入 Redis 等分布式组件、
  判断某个改动是否该走 PR、编写结构化审查报告、跨会话接续开发。
  即使用户没有提到"规范"，只要涉及本项目的编码风格、交付流程、
  协作方式，都应使用此技能。
  **本 skill 与 goddd 默认约定冲突时，以本 skill 为准**（冲突裁决见 §2）。
---

# YoFut 开发规范（项目专属）

> **定位**：goddd skill 管**架构**（六边形分层、Store/Adapter、Wire、缓存规范）；
> 本 skill 管**项目增量**（命名、注释、日志、原型策略、协作流程、Git 风格）。
> 两者叠加使用，无冲突处一律沿用 goddd。

---

## 1. 九条约束速查

| # | 约束 | 强度 | 专题文档 |
|---|---|---|---|
| C1 | **架构遵守 goddd skill** | 🔴 强制 | `.claude/skills/goddd/SKILL.md` + 12 份 references |
| C2 | **命名遵守 Entity-First / Noun-First** | 🔴 强制 | [`references/naming.md`](references/naming.md) |
| C3 | **注释：全中文 + 文件行头声明 + 函数详细注释** | 🔴 强制 | [`references/comment-style.md`](references/comment-style.md) |
| C4 | **日志：5 级体系（TRACE/DEBUG/INFO/WARN/ERROR）** | 🔴 强制 | [`references/logging.md`](references/logging.md) |
| C5 | **最小闭环先行**：真实 PostgreSQL + 简单 HTTP；连接读取 `configs/config.local.toml` 的 `[Data]`，SQLite 仅用于隔离单测 | 🔴 强制 | [`references/minimal-prototype.md`](references/minimal-prototype.md) |
| C6 | **AI 协作**：feature 分支 + PR + 结构化审查报告，人类独占合并权 | 🔴 强制 | [`references/ai-collaboration.md`](references/ai-collaboration.md) |
| C7 | **Git 风格**：`<type>(<scope>): <中文描述>` + 原子提交 + squash 线性历史 | 🔴 强制 | [`references/git-workflow.md`](references/git-workflow.md) |
| C8 | 敏感字段任何级别不落日志 | 🔴 强制 | `references/logging.md` §5 |
| C9 | 不碰 goddd `pkg/`，不手改 `wire_gen.go` | 🔴 强制 | goddd `SKILL.md` |
| C10 | **API 对外错误响应只暴露项目契约字段** | 🔴 强制 | 本节「API 错误边界」 |

---

## 2. 与 goddd 默认约定的冲突裁决

> goddd 的 `git-workflow` skill 面向多项目通用；
> 本项目按毕设实际情况作如下**覆盖**，未列出的部分一律沿用 goddd。

| 冲突点 | goddd 默认 | **YoFut 采用** | 理由 |
|---|---|---|---|
| 主线分支 | `main` + `dev` 双线 | **`dev` 集成、`main` 发布** | 已确认的自动版本与发布模型 |
| 个人分支命名 | `<initials>/<feature>` | **`feature/<描述>` / `fix/<描述>`** | 单人项目无需署名前缀 |
| 合并方式 | `dev → main` 用 `--no-ff` 保留合并节点 | **功能 PR squash 到 dev；dev → main 保留合并节点** | 保留发布范围内的语义提交与 PR 关联 |
| commit type 枚举 | 6 类（无 `docs`） | **7 类**（含 `docs`） | 毕设文档量大，需独立类型 |
| subject 语言 | 跟随交互语言 | **固定中文** | 本项目以中文交付 |
| 自动提交 | 完成后自动 commit、不询问 | **须走 PR；AI 不得直推 `main`** | 保留人工审查决策点 |
| 分支名风格 | kebab-case | **`feature/` 后用 kebab-case** | 一致 |

**仍沿用 goddd 的部分**：rebase 变基保持线性、rebase 后合并、Changelog 生成方式、
「一个 commit 只做一件事」、破坏性变更用 `!` 标注 + `BREAKING CHANGE` footer。

---

## 3. 硬性红线（违反即拒绝交付）

1. **禁止直接 push `main`**——所有变更走 `feature/*` 或 `fix/*` 分支 + PR
2. **禁止英文注释**——包括行内注释、函数文档、section header、TODO/FIXME
3. **禁止敏感字段落日志**——`password` / `token` / `secret` / `key`，任何级别都不输出
4. **禁止早期过度工程化**——新需求先以已有 PostgreSQL + 简单 HTTP 跑通最小闭环；不为验证另建分布式组件，不能用 SQLite 单测替代真实数据库开发验收
5. **禁止修改 goddd `pkg/`**——确需扩展时新增文件，并在 `NOTICE.md` 登记
6. **禁止手改 `wire_gen.go`**——改 `wire.go` 后运行 `wire` 重新生成
7. **禁止手工创建领域目录**——必须走 `tables/<domain>/<entity>.go` → `goddd gen`
8. **禁止混合性质的 commit**——安全 / 缺陷 / 风格必须拆分为独立 commit

### 3.1 API 错误边界

- 对外错误响应固定只允许 `code`、`message`、`data`、`request_id` 四个字段；错误时 `data` 固定为 `null`。
- goddd 的 `reason`、`msg`、`details`、`trace_id` 属于服务端内部诊断字段，只能用于内部错误映射和受控日志，禁止返回给前端、第三方调用方或写入公开文档契约。
- `request_id` 只作为不含业务语义的关联标识；堆栈、SQL、表名、内部路由细节、资源存在性判断等信息不得进入公开错误消息或 `data`。
- 路由不存在、方法不支持及下游业务错误必须经过 YoFut API 错误边界统一适配，分别返回项目定义的 404、405 或其他业务状态，不得直接透传 goddd 原始响应、HTML 或未结构化字符串。
- `ErrNotFound` 及资源不可见统一按 404 处理；未知错误统一按 500 处理；参数绑定和业务参数错误按 400 处理。公开错误消息使用面向调用方的安全文案，内部原因只在服务端保留。
- OpenAPI 的 `ErrorResponse` 只能声明上述四个公开字段；新增或修改错误字段时，必须同步更新实现、测试和 API 文档。

---

## 4. 规范优先级

```
人类口头指令  >  本 skill（YoFut）  >  goddd skill  >  通用最佳实践
```

- 用户明确说"这次按 X 做" → 按 X，不必援引规范
- 本 skill 未覆盖 → 查 goddd skill 的 12 份 references
- 两者都未覆盖 → 用通用最佳实践，并在 PR 中说明判断依据

---

## 5. 交付前自检清单

每次提交前逐项确认：

- [ ] 命名符合 Entity-First（实体名词在前，无 `Manager`/`Helper`/`Util` 泛词）
- [ ] 新文件有中文行头（文件类型 + 用途 + 用法）
- [ ] 公开函数有 GoDoc 格式中文注释，接口方法说明了契约
- [ ] 日志级别选择正确，敏感字段未输出
- [ ] 对外错误响应只有 `code`、`message`、`data`、`request_id`，未泄露 goddd 内部诊断字段
- [ ] 新需求已使用本地 `[Data.Database]` 的真实 PostgreSQL 验证；SQLite 仅承担隔离单测，未提前引入新分布式组件
- [ ] 分支名 `feature/*` 或 `fix/*`，commit 消息符合 `<type>(<scope>): <中文>`
- [ ] 多类修复已按「安全 > 缺陷 > 重构 > 风格 > 文档 > 清理」拆分
- [ ] `go build ./...` 与 `go vet ./...` 通过

---

## 6. 专题文档索引

| 文档 | 何时必读 |
|---|---|
| [`references/naming.md`](references/naming.md) | 新增任何标识符（包/类型/函数/表/路由/参数）之前 |
| [`references/comment-style.md`](references/comment-style.md) | 编写或修改任何函数、文件头、接口 |
| [`references/logging.md`](references/logging.md) | 添加任何日志语句、调整日志级别 |
| [`references/minimal-prototype.md`](references/minimal-prototype.md) | 接到新需求、考虑引入新组件时 |
| [`references/ai-collaboration.md`](references/ai-collaboration.md) | AI 参与编码/审查/交付时 |
| [`references/git-workflow.md`](references/git-workflow.md) | 创建分支、写 commit、建 PR、Code Review |
| [`references/cicd-actions.md`](references/cicd-actions.md) | 编写 `.gitea/workflows/*.yml`、调整 runner、CI 网络超时/TLS 失败 |
| goddd `.cursor/skills/goddd/references/*` | 领域实现、缓存层、适配器、代码生成、重构迁移 |
