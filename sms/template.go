// 本文件实现短信模板管理能力。
package sms

import (
	"context"

	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	service "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// AddTemplate 申请一个短信模板，提交后进入审核。
//
// 参数：
//   - ctx：请求上下文
//   - input：模板入参，TemplateName、TemplateType、TemplateContent 为必填
//
// 返回：error，传输或业务失败时返回，业务失败按阿里云错误码映射。
func (c *Client) AddTemplate(ctx context.Context, input *AddTemplateInput) error {
	req := &dysmsapi.AddSmsTemplateRequest{
		TemplateName:    tea.String(input.TemplateName),
		TemplateType:    tea.Int32(input.TemplateType),
		TemplateContent: tea.String(input.TemplateContent),
		Remark:          tea.String(input.Remark),
	}
	resp, err := c.dysmsapi.AddSmsTemplateWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return err
	}
	return checkBody(resp.Body)
}

// ModifyTemplate 修改一个已存在的短信模板，重新进入审核。
//
// 参数：
//   - ctx：请求上下文
//   - input：模板入参，TemplateCode 为必填
//
// 返回：error，传输或业务失败时返回。
func (c *Client) ModifyTemplate(ctx context.Context, input *ModifyTemplateInput) error {
	req := &dysmsapi.ModifySmsTemplateRequest{
		TemplateCode:    tea.String(input.TemplateCode),
		TemplateName:    tea.String(input.TemplateName),
		TemplateType:    tea.Int32(input.TemplateType),
		TemplateContent: tea.String(input.TemplateContent),
		Remark:          tea.String(input.Remark),
	}
	resp, err := c.dysmsapi.ModifySmsTemplateWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return err
	}
	return checkBody(resp.Body)
}

// DeleteTemplate 删除一个短信模板。
//
// 参数：
//   - ctx：请求上下文
//   - input：模板入参，TemplateCode 为必填
//
// 返回：error，传输或业务失败时返回。
func (c *Client) DeleteTemplate(ctx context.Context, input *DeleteTemplateInput) error {
	req := &dysmsapi.DeleteSmsTemplateRequest{
		TemplateCode: tea.String(input.TemplateCode),
	}
	resp, err := c.dysmsapi.DeleteSmsTemplateWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return err
	}
	return checkBody(resp.Body)
}

// QueryTemplate 查询单个短信模板的详情。
//
// 参数：
//   - ctx：请求上下文
//   - input：模板入参，TemplateCode 为必填
//
// 返回：
//   - *Template：模板详情
//   - error：传输或业务失败时返回
func (c *Client) QueryTemplate(ctx context.Context, input *QueryTemplateInput) (*Template, error) {
	req := &dysmsapi.QuerySmsTemplateRequest{
		TemplateCode: tea.String(input.TemplateCode),
	}
	resp, err := c.dysmsapi.QuerySmsTemplateWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if err := checkBody(resp.Body); err != nil {
		return nil, err
	}
	return &Template{
		TemplateCode:    tea.StringValue(resp.Body.TemplateCode),
		TemplateName:    tea.StringValue(resp.Body.TemplateName),
		TemplateType:    tea.Int32Value(resp.Body.TemplateType),
		TemplateContent: tea.StringValue(resp.Body.TemplateContent),
		TemplateStatus:  tea.Int32Value(resp.Body.TemplateStatus),
		CreateDate:      tea.StringValue(resp.Body.CreateDate),
		Reason:          tea.StringValue(resp.Body.Reason),
	}, nil
}

// ListTemplates 分页查询短信模板列表。
//
// 参数：
//   - ctx：请求上下文
//   - input：分页入参
//
// 返回：
//   - *TemplateList：模板列表集合
//   - error：传输或业务失败时返回
func (c *Client) ListTemplates(ctx context.Context, input *ListTemplatesInput) (*TemplateList, error) {
	req := &dysmsapi.QuerySmsTemplateListRequest{
		PageIndex: tea.Int32(input.PageIndex),
		PageSize:  tea.Int32(input.PageSize),
	}
	resp, err := c.dysmsapi.QuerySmsTemplateListWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if err := checkBody(resp.Body); err != nil {
		return nil, err
	}

	list := &TemplateList{TotalCount: tea.Int64Value(resp.Body.TotalCount)}
	list.Items = make([]*TemplateItem, 0, len(resp.Body.SmsTemplateList))
	for _, item := range resp.Body.SmsTemplateList {
		tpl := &TemplateItem{
			TemplateCode:    tea.StringValue(item.TemplateCode),
			TemplateName:    tea.StringValue(item.TemplateName),
			TemplateType:    tea.Int32Value(item.TemplateType),
			TemplateContent: tea.StringValue(item.TemplateContent),
			AuditStatus:     tea.StringValue(item.AuditStatus),
			CreateDate:      tea.StringValue(item.CreateDate),
		}
		if item.Reason != nil {
			tpl.Reason = tea.StringValue(item.Reason.RejectInfo)
		}
		list.Items = append(list.Items, tpl)
	}
	return list, nil
}
