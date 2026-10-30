// Package sms 提供阿里云短信（Aliyun SMS）服务的 SDK 封装，
// 覆盖短信发送、发送明细查询、签名管理与模板管理四类能力。
//
// 文件类型：SDK 入口 —— Client 结构与构造器
//
// 用法：
//
//	client, err := sms.NewClient()
//	if err != nil {
//		// 处理初始化错误
//	}
//	receipt, err := client.SendSMS(ctx, &sms.SendSMSInput{
//		PhoneNumbers:  "13900000000",
//		SignName:      "阿里云",
//		TemplateCode:  "SMS_153050000",
//		TemplateParam: `{"code":"123456"}`,
//	})
//
// 依赖方向：本包封装阿里云 dysmsapi 底层客户端，不依赖任何业务框架。
package sms

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

// Client 封装阿里云短信底层客户端，对外提供发送、查询、签名与模板管理能力。
type Client struct {
	// dysmsapi 阿里云短信底层客户端，仅包内使用。
	dysmsapi *dysmsapi.Client
}

// NewClient 构造阿里云短信客户端。
//
// 默认通过凭据链初始化访问凭据（环境变量 / ECS RAM Role 等），
// 可用 WithCredential、WithEndpoint、WithRegionID 覆盖，
// 或用 WithClient 注入已构造的底层客户端。
//
// 返回：
//   - *Client：就绪的客户端
//   - error：凭据初始化或客户端构造失败时返回
func NewClient(opts ...Option) (*Client, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.client != nil {
		return &Client{dysmsapi: o.client}, nil
	}

	cred := o.credential
	if cred == nil {
		var err error
		cred, err = credentials.NewCredential(nil)
		if err != nil {
			return nil, err
		}
	}

	config := &openapi.Config{
		Credential: cred,
		Endpoint:   tea.String(o.endpoint),
	}
	if o.regionID != "" {
		config.RegionId = tea.String(o.regionID)
	}
	client, err := dysmsapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{dysmsapi: client}, nil
}
