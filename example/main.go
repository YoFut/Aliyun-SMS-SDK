// Package main 提供阿里云短信 SDK 的可运行示例。
//
// 用法：
//
//	go run ./example
//
// 运行前需配置阿里云访问凭据，参见 README 的凭据说明。
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/YoFut/Aliyun-SMS/sms"
)

func main() {
	client, err := sms.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, "初始化客户端失败:", err)
		os.Exit(1)
	}

	receipt, err := client.SendSMS(context.Background(), &sms.SendSMSInput{
		PhoneNumbers:  "13900000000",
		SignName:      "阿里云",
		TemplateCode:  "SMS_153050000",
		TemplateParam: `{"name":"张三"}`,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "发送失败:", err)
		os.Exit(1)
	}
	fmt.Printf("发送成功，回执编号: %s，请求编号: %s\n", receipt.BizID, receipt.RequestID)
}
