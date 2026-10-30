// 本文件定义凭据与服务接入配置，供调用方从自身配置管理器（如 config.toml）解析后传入。
package sms

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

// Config 阿里云访问凭据与服务接入配置。
//
// 可由调用方从自身配置管理器（如 config.toml）解析后，通过 WithConfig 传入 NewClient；
// 凭据相关字段全部留空时，走默认凭据链（环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID /
// ALIBABA_CLOUD_ACCESS_KEY_SECRET、ECS RAM Role 等）。
type Config struct {
	// AccessKeyID 访问密钥 ID，与 AccessKeySecret 成对使用（access_key / sts 类型）。
	AccessKeyID string `toml:"access_key_id"`
	// AccessKeySecret 访问密钥 Secret。
	AccessKeySecret string `toml:"access_key_secret"`
	// SecurityToken STS 临时凭证令牌，与 AccessKeyID/AccessKeySecret 一并填写时使用 sts 类型。
	SecurityToken string `toml:"security_token"`
	// RoleName ECS 实例 RAM 角色名，填写时使用 ecs_ram_role 类型，与访问密钥二选一。
	RoleName string `toml:"role_name"`
	// Endpoint 服务接入点，留空时用默认 dysmsapi.aliyuncs.com。
	Endpoint string `toml:"endpoint"`
	// RegionID 地域标识，留空时用 SDK 默认。
	RegionID string `toml:"region_id"`
}

// buildCredential 将配置解析为阿里云凭据，凭据字段全空时返回 nil 表示走默认链。
func (c Config) buildCredential() (credentials.Credential, error) {
	switch {
	case c.RoleName != "":
		return credentials.NewCredential(&credentials.Config{
			Type:     tea.String("ecs_ram_role"),
			RoleName: tea.String(c.RoleName),
		})
	case c.AccessKeyID != "" && c.AccessKeySecret != "":
		typ := "access_key"
		if c.SecurityToken != "" {
			typ = "sts"
		}
		return credentials.NewCredential(&credentials.Config{
			Type:            tea.String(typ),
			AccessKeyId:     tea.String(c.AccessKeyID),
			AccessKeySecret: tea.String(c.AccessKeySecret),
			SecurityToken:   tea.String(c.SecurityToken),
		})
	default:
		return nil, nil
	}
}
