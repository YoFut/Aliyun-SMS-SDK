# 注释规范

> 触发时机：编写或修改任何文件头、函数、方法、接口、结构体字段时必读。

---

## 1. 四条铁律

1. **全中文** —— 行内注释、函数文档、section header、`TODO` / `FIXME` 一律中文
   （标识符、协议名、命令、API 路径保持英文原文，不翻译）
2. **所有函数必须有注释**，按公开程度采用不同粒度（见 §3）
3. **不写冗余注释** —— 不解释「代码在做什么」，只解释「为什么这么做」
4. **复杂并发逻辑必须配 ASCII 图** 说明锁的获取顺序

---

## 2. 文件行头声明

每个 `.go` 文件顶部必须声明：**文件类型 + 用途 + 用法**。

```go
// Package upload 提供资源上传领域能力，并实现上传权限五道闸门。
//
// 文件类型：领域 Core 层 —— 上传业务方法与 Storer 接口定义
//
// 用法：
//
//	core := upload.NewCore(store, policy, credentialIssuer)
//	version, err := core.CreateVersion(ctx, actor, input)
//
// 依赖方向：本包只依赖 port.go 中定义的接口，不感知 OSS / GORM 等具体实现。
package upload
```

**非 package 文件**（`.model.go` / `.param.go` / 适配器实现）在文件头说明**本文件职责**：

```go
// 本文件定义 Upload 的 GORM 持久化模型。
// goddd gen 生成后手工补充字段说明；重新生成前请先备份本文件的自定义部分。
//
// 对应数据库表：resource_versions
package upload
```

---

## 3. 函数注释三档粒度

### 3.1 公开函数 / 方法 —— GoDoc 格式

**首行以函数名开头**，随后说明：用途 → 参数含义 → 返回值语义 → 错误条件。

```go
// CreateVersion 为指定资源创建新版本，并签发受前缀约束的上传凭证。
//
// 参数：
//   - actor：调用者身份，企业 ID 与内部员工标记均从此取值，不接受外部传入
//   - input：版本创建入参，其中 EnterpriseID 若与 actor 不一致将被拒绝
//
// 返回：
//   - *Version：创建成功的版本记录（状态为 draft，需审核后可见）
//   - error：ErrUploadForbidden 表示调用者无上传权限；
//     ErrTenantMismatch 表示入参企业 ID 与调用者不一致
//
// 副作用：写入 resource_versions 表、签发对象存储临时凭证、记录上传审计。
func (c *Core) CreateVersion(ctx context.Context, actor *auth.Actor, input *CreateVersionInput) (*Version, error) {
```

### 3.2 内部函数 —— 单行注释

```go
// canUpload 判断调用者是否满足上传策略，不查库、纯内存判定。
func (c *Core) canUpload(actor *auth.Actor) bool {
```

非显而易见的参数需补充说明：

```go
// sanitizeKey 清洗对象键，depth 表示需保留的路径层级数（企业私有前缀固定占 2 层）。
func sanitizeKey(key string, depth int) string {
```

### 3.3 接口方法 —— 说明契约

必须写明：**调用前提、副作用、错误条件**。

```go
// ObjectStorer 定义对象存储的被动适配器契约，由 adapter/oss 实现。
type ObjectStorer interface {
	// IssueUploadCredential 签发受前缀约束的上传凭证。
	//
	// 前提：调用者已完成权限校验，且 prefix 必须以 ent/<企业ID>/ 开头。
	// 副作用：无（仅签发凭证，不产生对象）。
	// 错误：ErrInvalidPrefix 表示前缀未通过校验。
	IssueUploadCredential(ctx context.Context, prefix string, ttl time.Duration) (*Credential, error)

	// DeleteObject 幂等删除对象。
	//
	// 前提：无。
	// 副作用：对象不可恢复；对象不存在时不报错。
	// 错误：仅返回底层存储故障，不返回「对象不存在」。
	DeleteObject(ctx context.Context, key string) error
}
```

---

## 4. 不写冗余注释（反例对照）

| ❌ 冗余 | ✅ 应改为 |
|---|---|
| `// 初始化变量` | 删除 |
| `// 返回结果` | 删除 |
| `// 遍历列表` | 删除（`for` 已说明） |
| `// 设置 name 字段` | 删除 |
| `// 如果 err 不为空则返回` | 删除 |
| `// 加锁` | `// 此处加锁防止并发签发同一版本凭证时重复上传` |
| `// 跳过` | `// 分片缓存缺少任一分片即视为不完整，整体重建而非返回半截数据` |

**判断标准**：注释删除后，读者是否还能理解**为什么**这么写？

- 能 → 是冗余注释，删除
- 不能 → 保留，并确保它解释的是「为什么」而非「是什么」

---

## 5. 复杂并发逻辑：ASCII 图说明锁顺序

涉及多把锁时，必须在函数上方画出获取顺序，避免死锁。

```go
// 本函数按固定顺序获取两把锁，顺序不可颠倒：
//
//	企业锁 (ent:<id>)  ──►  资源锁 (res:<id>)
//	     │                        │
//	     └── 必须先拿 ────────────┘
//
// 反向获取会与 ApproveVersion 形成死锁。
func (c *Core) PublishVersion(ctx context.Context, entID, resID int64) error {
```

---

## 6. section header

大型文件用中文分隔区块：

```go
// =============================================================================
// 闸门 2：中间件三重校验
// =============================================================================
```

---

## 7. TODO / FIXME 规范

必须中文，并标注**责任人与触发条件**：

```go
// TODO(邸维阳)：接入 ClamAV 病毒扫描，待 S7 阶段部署扫描服务后启用。
// FIXME(邸维阳)：当前敏感词词库为静态加载，热更新见 T4.4。
```

禁止无主 TODO（如 `// TODO: 优化`）。

---

## 8. 注释自检

- [ ] 全部注释为中文（标识符 / 协议 / 路径除外）
- [ ] 每个新文件有行头：文件类型 + 用途 + 用法
- [ ] 每个函数有注释，公开函数为 GoDoc 格式且首行以函数名开头
- [ ] 接口方法说明了「前提 / 副作用 / 错误」
- [ ] 无冗余注释（删除后不影响理解「为什么」）
- [ ] 并发逻辑配了 ASCII 锁顺序图
- [ ] `TODO` / `FIXME` 有责任人与触发条件
