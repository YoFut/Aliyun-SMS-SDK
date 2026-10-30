# 日志规范

> 触发时机：添加任何日志语句、调整日志级别、配置日志输出时必读。
> 基础：goddd `pkg/logger`（Zap 封装）。本规范在其之上补充 **TRACE 级别** 与 **字段标准**。

---

## 1. 五级日志体系

| 级别 | Zap 值 | 用途 | 环境 |
|---|---|---|---|
| **TRACE** | `-2`（自定义） | 函数入口/出口、变量 dump、缓存**逐 key** 命中 | 仅开发 |
| **DEBUG** | `-1` | 请求/响应 payload（脱敏）、SQL、Redis 命令、配置 dump | 仅开发 |
| **INFO** | `0` | 请求完成、服务启停、健康检查、缓存预热完成 | **生产最低** |
| **WARN** | `1` | 重试、降级触发、限流预警、**非安全类**认证失败 | 全环境 |
| **ERROR** | `2` | 未处理错误、依赖故障、数据不一致、panic 恢复 | 全环境 |

> 安全类认证失败（如签名伪造、越权尝试）属 **ERROR**，不属 WARN。

### 1.1 级别选择决策树

```
是否属于安全事件（伪造/越权/注入）？  → 是 → ERROR
是否为未处理错误或依赖故障？          → 是 → ERROR
是否为「可自愈」的异常？               → 是 → WARN  （重试 / 降级 / 限流预警）
是否为正常业务完成事件？               → 是 → INFO
是否需要排查细节（SQL/Redis/payload）？ → 是 → DEBUG
是否属于逐行跟踪（进出/命中）？        → 是 → TRACE
```

---

## 2. TRACE 级别的实现

goddd 的 `pkg/logger` 基于 Zap，Zap 支持自定义级别，**无需修改框架**，在项目侧扩展即可。

```go
// 本文件扩展 goddd pkg/logger，补充 TRACE 级别（-2）。
//
// 用法：
//
//	logger.Trace("函数入口", zap.String("func", "CreateVersion"))
package logging

import "go.uber.org/zap/zapcore"

// LevelTrace 定义 TRACE 级别（-2），低于 DEBUG，仅开发环境启用。
//
// 选择 -2 而非 -1 的原因：goddd 已占用 DEBUG(-1)，TRACE 需与其区分。
const LevelTrace = zapcore.Level(-2)

// Trace 输出 TRACE 级别日志，仅开发环境可见。
//
// 前提：Logger 已通过 EnableTrace() 开启 TRACE 级别。
// 副作用：无。
func Trace(msg string, fields ...zap.Field) {
	// 跳过 1 层调用栈，使日志位置指向业务代码而非本封装函数
	//
	//	业务代码 ──► Trace() ──► 底层 zap
	//	              ↑ 需跳过的这一层
	if Logger != nil {
		Logger.Log(LevelTrace, msg, fields...)
	}
}

// EnableTrace 将 Logger 的最低输出级别降至 TRACE。
//
// 调用时机：仅在 debug 模式下由 main.go 调用；生产环境绝不调用。
func EnableTrace() {
	if Logger != nil {
		Logger.Core().Enabled(LevelTrace)
	}
}
```

**开启方式**（`internal/web/api/api.go` 或 `main.go`）：

```go
// 仅开发环境开启 TRACE，生产环境保持 INFO 及以上
if conf.Runtime.Debug {
	logging.EnableTrace()
}
```

---

## 3. 结构化日志字段标准

### 3.1 HTTP 请求通用字段（必须）

| 字段 | 类型 | 说明 |
|---|---|---|
| `request_id` | string | 全链路追踪 ID，中间件生成并写入 context |
| `latency` | duration | 请求耗时（`time.Since(start)`） |
| `status` | int | HTTP 状态码 |
| `method` | string | HTTP 方法 |
| `path` | string | 请求路径（**不含 query，避免泄露凭证**） |

```go
logging.Info("请求完成",
	zap.String("request_id", rid),
	zap.Duration("latency", time.Since(start)),
	zap.Int("status", c.Writer.Status()),
	zap.String("method", c.Request.Method),
	zap.String("path", c.Path()),
)
```

### 3.2 业务字段扩展（按需，同为 snake_case）

| 场景 | 建议字段 |
|---|---|
| 领域操作 | `domain` · `entity` · `entity_id` · `actor_id` · `enterprise_id` |
| 上传 | `object_key` · `size_bytes` · `sha256` · `gate`（命中哪道闸门） |
| 审核 | `target_type` · `target_id` · `reviewer_id` · `risk_score` |
| 缓存 | `cache_layer`（l1/l2）· `cache_key` · `hit` |
| 同步 | `software_id` · `tag` · `new_versions` · `skipped` |

---

## 4. 环境级别

| 环境 | 最低输出级别 | 格式 | 说明 |
|---|---|---|---|
| **开发**（`Runtime.Debug=true`） | **TRACE** | Console 彩色 | 含请求体/响应体（已脱敏） |
| **Staging** | **DEBUG** | JSON | 记录 SQL 与 Redis 命令 |
| **生产** | **INFO** | JSON | 禁止 DEBUG 及以下 |

```toml
# configs/config.toml
[Log]
  Dir = './logs'
  Level = 'debug'          # 开发；生产改 'info'
  MaxAge = 7
  RotationTime = '8h0m0s'
  MaxSize = 50
  Compress = false
```

---

## 5. 敏感字段（任何级别都不输出）

### 5.1 黑名单

```
password · passwd · pwd
token · access_token · refresh_token · jwt
secret · client_secret · jwt_secret
key · access_key · secret_key · api_key
authorization · cookie · set-cookie
```

### 5.2 强制要求

1. **禁止** 直接 `zap.Any("body", req)` 打印整个请求体——必须逐字段挑选
2. **禁止** 打印 `Authorization` 头与 Cookie（含 `Set-Cookie`）
3. 需要记录敏感字段**存在性**时，只记布尔值：

```go
// ✅ 正确：只记是否存在，不记值
logging.Debug("鉴权头解析",
	zap.Bool("has_authorization", c.Get("Authorization") != ""),
)

// ❌ 错误：整包打印，token 会落盘
logging.Debug("请求详情", zap.Any("headers", c.GetReqHeaders()))
```

4. 对象存储凭证只记 **TTL 与前缀**，不记 `AccessKeyId` / `SecurityToken`：

```go
logging.Info("签发上传凭证",
	zap.String("prefix", cred.Prefix),
	zap.Duration("ttl", cred.TTL),
	// 不记录 cred.AccessKeyId / cred.SecurityToken
)
```

---

## 6. 各层级日志规范（贴合 goddd 六边形）

| 层 | 允许级别 | 示例 |
|---|---|---|
| **API 层**（`internal/web/api`） | INFO / WARN / ERROR | 请求完成、鉴权失败、限流触发 |
| **Core 层**（`internal/core/*`） | WARN / ERROR / DEBUG | 业务规则拒绝（WARN）、领域不变量违约（ERROR） |
| **Store 层**（`stores/*db` `*cache`） | DEBUG / WARN / ERROR | SQL 耗时（DEBUG）、缓存穿透（DEBUG）、DB 故障（ERROR） |
| **Adapter 层**（`adapter/*`） | DEBUG / WARN / ERROR | 外部调用耗时（DEBUG）、重试（WARN）、上游不可用（ERROR） |

> **注意**：Core 层**不应**出现 TRACE，避免领域逻辑被逐行跟踪污染；
> 需要逐行跟踪时在 Adapter 层或临时本地调试。

---

## 7. 日志自检

- [ ] 级别选择符合 §1.1 决策树
- [ ] HTTP 日志含 `request_id` / `latency` / `status` / `method` / `path`
- [ ] 未打印 `password` / `token` / `secret` / `key` 及其载体对象
- [ ] 未整包打印请求体或请求头
- [ ] 生产环境最低 INFO，DEBUG/TRACE 仅在开发与 Staging 启用
- [ ] 字段名为 snake_case，与 §3.2 的扩展字段对齐
- [ ] 错误日志带 `zap.Error(err)` 而不仅是字符串拼接
