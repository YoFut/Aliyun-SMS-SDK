# Aliyun SMS SDK

阿里云短信服务（Aliyun SMS）的 Go SDK，封装短信发送、发送明细查询、签名管理与模板管理四类能力。

## 安装

```go
import "github.com/YoFut/Aliyun-SMS/sms"
```

## 凭据

默认走阿里云凭据链，无需显式传入：依次读取环境变量
`ALIBABA_CLOUD_ACCESS_KEY_ID` / `ALIBABA_CLOUD_ACCESS_KEY_SECRET`、ECS RAM Role 等。

也可从调用方自身的配置管理器（如 `config.toml`）解析后，通过 `WithConfig` 传入 AK/SK、STS 或 ECS RAM Role 凭据（示例见 `configs/config.toml.example`）；凭据字段留空时回退到默认凭据链。需要更细粒度控制时可用 `WithCredential` 显式注入。

## 快速开始

```go
client, err := sms.NewClient()
if err != nil {
    log.Fatal(err)
}

receipt, err := client.SendSMS(ctx, &sms.SendSMSInput{
    PhoneNumbers:  "13900000000",
    SignName:      "阿里云",
    TemplateCode:  "SMS_153050000",
    TemplateParam: `{"code":"123456"}`,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(receipt.BizID) // 回执编号，可用于查询发送状态
```

## 能力清单

| 分类 | 方法 |
|---|---|
| 发送 | `SendSMS`、`SendBatchSMS` |
| 查询 | `QuerySendDetails` |
| 签名 | `AddSign`、`ModifySign`、`DeleteSign`、`QuerySign`、`ListSigns` |
| 模板 | `AddTemplate`、`ModifyTemplate`、`DeleteTemplate`、`QueryTemplate`、`ListTemplates` |
| 验证码 | `SendVerificationCode`、`VerifyVerificationCode` |

## 错误处理

业务失败会返回带阿里云错误码的错误，常见错误映射为哨兵变量，可用 `errors.Is` 判断：

```go
if errors.Is(err, sms.ErrDailyQuotaExceeded) {
    // 当日额度已满
}
```

## 验证码

发送验证码时，调用侧只需指定验证码长度，验证码由 SDK 生成、发送并暂存，随后用 `VerifyVerificationCode` 核验用户提交的码：

```go
_, err := client.SendVerificationCode(ctx, &sms.SendVerificationCodeInput{
    PhoneNumbers: "13900000000",
    SignName:     "阿里云",
    TemplateCode: "SMS_153050000",
    CodeLength:   6,
})
// 用户提交后核验
if client.VerifyVerificationCode(ctx, &sms.VerifyVerificationCodeInput{
    PhoneNumber: "13900000000",
    Code:        "123456",
}) {
    // 核验通过
}
```

验证码默认 5 分钟有效、命中后即失效，可用 `WithCodeTTL` 覆盖有效期；存储为进程内存。

## 运行示例

```bash
go run ./example
```
