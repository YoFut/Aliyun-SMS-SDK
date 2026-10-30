# 命名规范 — Entity-First / Noun-First

> 触发时机：新增任何标识符（包 / 类型 / 函数 / 方法 / 表 / 字段 / 路由 / 参数）之前必读。

---

## 1. 两条核心原则

### 原则一：Entity-First（实体优先）

标识符以**业务实体**开头，不以**技术角色**开头。

```
✅ UserStorer        — 实体 User 的存储接口
❌ StorerUser / UserRepo / UserDAO

✅ ThreadController  ❌ ControllerThread
✅ CommentLike       ❌ LikeComment
```

### 原则二：Noun-First（名词优先）

- **名词性标识符**（类型、表、包、字段、参数）：**纯名词**，不得含动词
- **动作性标识符**（函数、方法）：**动词 + 名词**

```
✅ type User struct{}                      ❌ type CreateUser struct{}
✅ func (c *Core) CreateUser(...)          ❌ func (c *Core) UserCreate(...)
✅ func (c *Core) ListThreads(...)         ❌ func (c *Core) ThreadsList(...)
```

---

## 2. 分层命名映射表

| 层 | 规则 | 示例 |
|---|---|---|
| **领域包** | 单数名词、全小写、无下划线 | `identity` `upload` `forum` `comment` `moderation` |
| **实体类型** | 大驼峰名词，无 `Entity`/`Model` 后缀 | `User` `Enterprise` `Thread` `Comment` |
| **表模型**（`tables/`） | 同实体名（goddd 生成器约定） | `tables/forum/thread.go` → `type Thread struct` |
| **数据库表** | snake_case **复数** | `users` `enterprises` `threads` `comments` |
| **Storer 接口** | `<Entity>Storer` | `UserStorer` `ThreadStorer` |
| **Core 方法** | `动词 + <Entity>`（复数表示集合） | `CreateUser` `ListThreads` `ApproveComment` |
| **入参结构** | `<动作><Entity>Input` | `CreateUserInput` `ListThreadInput` |
| **出参结构** | `<Entity>` / `<Entity>List` | `User` / `UserList` |
| **端口接口** | `<能力>Port` | `ObjectStorer`（对象存储）· `UpstreamPort`（回源） |
| **适配器目录** | 供应商名小写 | `adapter/oss` `adapter/github` `adapter/casbin` |
| **错误变量** | `Err<条件>` | `ErrUserNotFound` `ErrUploadForbidden` |
| **API 路由** | `/api/v1/<复数名词>` | `POST /api/v1/users` · `GET /api/v1/threads/:id` |
| **权限点** | `<资源>:<动作>` | `resource:upload` `comment:approve` |
| **配置节** | TOML 大驼峰 | `[Upload]` `[Forum]` `[Moderation]` |

---

## 3. 禁用清单（泛词与角色后缀）

| 禁用 | 原因 | 替代 |
|---|---|---|
| `Manager` `Helper` `Util` `Utils` `Common` `Base` | 语义空洞，无所指 | 用实体名：`UploadPolicy` 而非 `UploadManager` |
| `Data` `Info` `Object` `Item` | 泛化 | 用具体实体名 |
| `Handler`（在 Core 层） | 与 API 层混淆 | Core 层用实体方法名 |
| `Service`（在 Core 层） | goddd 用 Core 而非 Service | 用 `<Entity>Core` 或直接领域包名 |
| `Impl` `Default` 后缀 | 接口与实现同名更清晰 | `casbin.Authorizer` 而非 `AuthorizerImpl` |
| 拼音 / 中英混拼 | 不可读 | 全英文（注释才是中文） |
| 单字母变量（除循环下标 `i`/`j`、错误 `err`、上下文 `ctx`） | 可读性差 | 完整单词 |
| 动词开头的类型名 | 违反 Noun-First | `SyncTask` 而非 `DoSync` |

---

## 4. YoFut 九个领域的实体清单

> 新增标识符前先在此对齐；确需新增实体时，先更新本表再写代码。

| 领域 | 核心实体 | 表名 |
|---|---|---|
| `identity` | `User` · `Enterprise` · `Role` · `Permission` | `users` `enterprises` `roles` `permissions` |
| `upload` | `Resource` · `Version` · `Asset` | `resources` `versions` `assets` |
| `mirror` | `MirrorFile` | `mirror_files` |
| `catalog` | `Software` · `SoftwareTag` | `softwares` `software_tags` |
| `syncjob` | `SyncTask` · `SyncLog` | `sync_tasks` `sync_logs` |
| `forum` | `Board` · `Thread` · `Post` | `boards` `threads` `posts` |
| `comment` | `Comment` · `CommentLike` | `comments` `comment_likes` |
| `moderation` | `ModerationTask` · `Report` | `moderation_tasks` `reports` |
| `audit` | `OperationLog` · `LoginLog` | `operation_logs` `login_logs` |

---

## 5. goddd 生成链路对命名的约束

`goddd gen -f tables/<domain>/<entity>.go` 会按文件与类型名推导一整条命名链，
因此**表模型命名是全局命名的源头**：

```
tables/forum/thread.go
  └─ type Thread struct{ ID; CreatedAt; UpdatedAt }
        ↓ goddd gen
  internal/core/forum/core.go            Core 结构体 + Storer 聚合接口
  internal/core/forum/thread.go          ThreadStorer 接口 + 业务方法
  internal/core/forum/thread.model.go    GORM 持久化模型
  internal/core/forum/thread.param.go    CreateThreadInput / ListThreadInput
  internal/core/forum/stores/forumdb/    数据库实现
  internal/core/forum/stores/forumcache/ 缓存实现
```

**结论**：`tables/` 里的类型名一旦定错，生成物全线污染。
**新增实体前必须先过本节与 §2 映射表。**

---

## 6. 命名自检

- [ ] 类型/表/包/字段是**纯名词**，函数/方法是**动词 + 名词**
- [ ] 无 `Manager`/`Helper`/`Util`/`Common`/`Base`/`Data`/`Info` 等泛词
- [ ] 实体名在前，技术角色在后（`UserStorer` 而非 `StorerUser`）
- [ ] 表名 snake_case 复数，与实体单数名对应
- [ ] 新实体已登记到 §4 清单
- [ ] 路由用复数名词资源（`/api/v1/threads`）
