# AI 协作规范

> 触发时机：WaveYo 以 AI 助手形态参与编码、审查或交付时必读。
> 核心前提：**人类独占合并决策权，AI 只负责执行、审查和报告。**

---

## 1. 交付通道（强制）

| 规则 | 说明 |
|---|---|
| **所有变更走分支 + PR** | `feature/<描述>` 或 `fix/<描述>`；**禁止直接 push `main`** |
| **AI 不执行合并** | AI 可创建 PR、可提交审查意见，但不点 merge |
| **人类审查后合并** | 按 `git-workflow.md` 的 Checklist 逐项确认 |

---

## 2. 结构化审查报告（每次 PR 必须附带）

按 **Security / Defect / Style** 三级分类，每条**必须含文件路径 + 行号**。

```markdown
## 审查报告

### 🔴 Security（安全）

| # | 文件:行 | 问题 | 严重度 | 修复方案 |
|---|---|---|---|---|
| S1 | `internal/core/upload/upload.go:87` | 企业 ID 取自请求体，可越权写入他人租户 | 高 | 改从 `actor.EnterpriseID` 取值 |
| S2 | `internal/web/api/api.go:112` | 调试日志打印完整请求头，含 Authorization | 中 | 改为只记 `has_authorization` 布尔值 |

### 🟠 Defect（缺陷）

| # | 文件:行 | 问题 | 修复方案 |
|---|---|---|---|
| D1 | `internal/core/forum/thread.go:45` | `ListThreads` 空结果返回 `null` 而非 `[]` | 用 `make([]*Thread, 0)` 初始化 |

### 🟡 Style（风格）

| # | 文件 | 问题 | 修复方案 |
|---|---|---|---|
| Y1 | `internal/core/comment/comment.go` | 命名 `CommentManager` 违反 Noun-First | 改为 `CommentPolicy` |
```

**无问题时**也要显式写「本次审查未发现问题」，不得省略该章节。

---

## 3. Atomic Commit 原则

一个 commit 对应一个**逻辑单元**；多类别修复按严重度拆分：

```
commit 1: fix(security): 修复上传越权与企业 ID 注入
commit 2: fix(defects): 修复空列表返回 null、context 泄漏
commit 3: chore(style): 统一命名与注释规范
```

**优先级排序**：

```
安全  >  缺陷  >  重构  >  风格  >  文档  >  清理
```

**拆分粒度判据**：

| 情形 | 处理 |
|---|---|
| 同一函数的不同行 | 合并在一个 commit |
| 不同模块 / 不同性质 | 拆分为多个 commit |
| 重构 + 行为修复混在一起 | **必须拆开**（重构不得改变语义） |

---

## 4. 跨会话延续提示词（Session Handoff）

每次会话结束或上下文即将耗尽时，**主动输出**如下结构化提示词：

````markdown
## 会话延续提示词

**HEAD commit**：`<commit-sha> <message>`
**分支**：`feature/<当前分支>`
**仓库**：`http://localhost:11451/Graduation-Project/YoFut.git`

### 已完成清单
- [x] <具体事项 1>（对应 WBS 任务号）
- [x] <具体事项 2>

### 待讨论方向
- [ ] <下一步要做什么>
- [ ] <未决问题>

### 本次会话确立的规范/决策
- <决策 1>：<理由>
- <决策 2>：<理由>

### 阻塞项
- <阻塞点>（如需人工确认的事项）

### 规范摘要（供新会话快速加载）
- 架构：goddd 六边形（`.cursor/skills/goddd/`）
- 项目规范：`.cursor/skills/YoFut/SKILL.md` + 6 份 references
- 命名：Entity-First / Noun-First
- 注释：全中文 + 文件行头 + GoDoc
- 日志：TRACE/DEBUG/INFO/WARN/ERROR 五级
- 交付：feature 分支 + PR + squash merge
````

---

## 5. 集成测试失败诊断顺序

按固定顺序排查，**不得跳步**：

```
① 客户端配置  ──►  ② 基础设施  ──►  ③ 替代方案
```

### ① 客户端配置（最常见）

- [ ] Endpoint / Region 是否正确（OSS 用 `s3.oss-<region>.aliyuncs.com`）
- [ ] `UsePathStyle` 是否为 OSS 所需的值
- [ ] AccessKey / SecretKey 是否有效且未过期
- [ ] Bucket 名称与权限策略是否匹配
- [ ] 网络代理环境变量（`HTTP_PROXY` / `HTTPS_PROXY`）是否干扰

### ② 基础设施

- [ ] PostgreSQL / Redis 容器是否健康（`docker compose ps`）
- [ ] 迁移是否已执行（表是否存在）
- [ ] 端口映射是否正确（WSL2 下 `wslrelay` 转发）
- [ ] CI 中的 service container 是否就绪（等待 healthcheck）

### ③ 替代方案

- [ ] 用本地 `[Data.Database]` 只读连接检查及 PostgreSQL 独立 schema 最小复现区分代码/环境问题；SQLite 仅作补充隔离单测，不能替代 PG 开发验收
- [ ] 换官方 CLI（`ossutil`）验证凭证可用性
- [ ] 降级为最小复现脚本，排除框架干扰

> **禁止行为**：测试失败后直接改断言或加 `skip`。
> 必须定位到根因，或在 PR 中明确说明「暂时跳过 + 原因 + 恢复条件」。

---

## 6. AI 行为红线

- ❌ 直接 push `main`
- ❌ 替人类做合并决策
- ❌ 未经确认执行破坏性操作（删除数据、强推、重置历史）
- ❌ 测试失败后修改断言以「通过」
- ❌ 声称「已完成」但未实际验证（必须先跑 `go build` / `go vet` / 测试）
- ❌ 把不确定的事说成确定的（不确定就标注「待验证」）
