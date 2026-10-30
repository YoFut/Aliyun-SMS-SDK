// 本文件定义客户端配置选项。
package sms

import (
	"time"

	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/aliyun/credentials-go/credentials"
)

// defaultEndpoint 阿里云短信公共接入点。
const defaultEndpoint = "dysmsapi.aliyuncs.com"

// defaultCodeTTL 验证码默认有效期。
const defaultCodeTTL = 5 * time.Minute

// options 客户端可选配置，由 Option 函数逐一填充。
type options struct {
	credential credentials.Credential // 访问凭据
	endpoint   string                 // 服务接入点
	regionID   string                 // 地域标识
	codeTTL    time.Duration          // 验证码有效期
	client     *dysmsapi.Client       // 已构造的底层客户端
}

// Option 客户端配置函数，遵循函数式选项模式。
type Option func(*options)

// defaultOptions 返回默认配置。
func defaultOptions() *options {
	return &options{
		endpoint: defaultEndpoint,
		codeTTL:  defaultCodeTTL,
	}
}

// WithCredential 显式注入访问凭据，不指定时走默认凭据链。
func WithCredential(credential credentials.Credential) Option {
	return func(o *options) { o.credential = credential }
}

// WithEndpoint 覆盖服务接入点，默认 dysmsapi.aliyuncs.com。
func WithEndpoint(endpoint string) Option {
	return func(o *options) { o.endpoint = endpoint }
}

// WithRegionID 覆盖地域标识，短信服务通常无需指定。
func WithRegionID(regionID string) Option {
	return func(o *options) { o.regionID = regionID }
}

// WithCodeTTL 覆盖验证码有效期，默认 5 分钟。
func WithCodeTTL(ttl time.Duration) Option {
	return func(o *options) { o.codeTTL = ttl }
}

// WithClient 注入已构造的底层客户端，主要用于测试与复用。
func WithClient(client *dysmsapi.Client) Option {
	return func(o *options) { o.client = client }
}
