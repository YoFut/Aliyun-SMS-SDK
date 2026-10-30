// 本文件实现短信签名管理能力。
package sms

import (
	"context"

	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	service "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// AddSign 申请一个短信签名，提交后进入审核。
//
// 参数：
//   - ctx：请求上下文
//   - input：签名入参，SignName、SignSource 为必填
//
// 返回：error，传输或业务失败时返回，业务失败按阿里云错误码映射。
func (c *Client) AddSign(ctx context.Context, input *AddSignInput) error {
	req := &dysmsapi.AddSmsSignRequest{
		SignName:     tea.String(input.SignName),
		SignSource:   tea.Int32(input.SignSource),
		SignType:     tea.Int32(input.SignType),
		Remark:       tea.String(input.Remark),
		SignFileList: make([]*dysmsapi.AddSmsSignRequestSignFileList, 0, len(input.SignFiles)),
	}
	for _, f := range input.SignFiles {
		req.SignFileList = append(req.SignFileList, &dysmsapi.AddSmsSignRequestSignFileList{
			FileContents: tea.String(f.FileContents),
			FileSuffix:   tea.String(f.FileSuffix),
		})
	}

	resp, err := c.dysmsapi.AddSmsSignWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return err
	}
	return checkBody(resp.Body)
}

// ModifySign 修改一个已存在的短信签名，重新进入审核。
//
// 参数：
//   - ctx：请求上下文
//   - input：签名入参，字段含义同 AddSignInput
//
// 返回：error，传输或业务失败时返回。
func (c *Client) ModifySign(ctx context.Context, input *ModifySignInput) error {
	req := &dysmsapi.ModifySmsSignRequest{
		SignName:     tea.String(input.SignName),
		SignSource:   tea.Int32(input.SignSource),
		SignType:     tea.Int32(input.SignType),
		Remark:       tea.String(input.Remark),
		SignFileList: make([]*dysmsapi.ModifySmsSignRequestSignFileList, 0, len(input.SignFiles)),
	}
	for _, f := range input.SignFiles {
		req.SignFileList = append(req.SignFileList, &dysmsapi.ModifySmsSignRequestSignFileList{
			FileContents: tea.String(f.FileContents),
			FileSuffix:   tea.String(f.FileSuffix),
		})
	}

	resp, err := c.dysmsapi.ModifySmsSignWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return err
	}
	return checkBody(resp.Body)
}

// DeleteSign 删除一个短信签名。
//
// 参数：
//   - ctx：请求上下文
//   - input：签名入参，SignName 为必填
//
// 返回：error，传输或业务失败时返回。
func (c *Client) DeleteSign(ctx context.Context, input *DeleteSignInput) error {
	req := &dysmsapi.DeleteSmsSignRequest{
		SignName: tea.String(input.SignName),
	}
	resp, err := c.dysmsapi.DeleteSmsSignWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return err
	}
	return checkBody(resp.Body)
}

// QuerySign 查询单个短信签名的详情。
//
// 参数：
//   - ctx：请求上下文
//   - input：签名入参，SignName 为必填
//
// 返回：
//   - *Sign：签名详情
//   - error：传输或业务失败时返回
func (c *Client) QuerySign(ctx context.Context, input *QuerySignInput) (*Sign, error) {
	req := &dysmsapi.QuerySmsSignRequest{
		SignName: tea.String(input.SignName),
	}
	resp, err := c.dysmsapi.QuerySmsSignWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if err := checkBody(resp.Body); err != nil {
		return nil, err
	}
	return &Sign{
		SignName:   tea.StringValue(resp.Body.SignName),
		CreateDate: tea.StringValue(resp.Body.CreateDate),
		SignStatus: tea.Int32Value(resp.Body.SignStatus),
		Reason:     tea.StringValue(resp.Body.Reason),
	}, nil
}

// ListSigns 分页查询短信签名列表。
//
// 参数：
//   - ctx：请求上下文
//   - input：分页入参
//
// 返回：
//   - *SignList：签名列表集合
//   - error：传输或业务失败时返回
func (c *Client) ListSigns(ctx context.Context, input *ListSignsInput) (*SignList, error) {
	req := &dysmsapi.QuerySmsSignListRequest{
		PageIndex: tea.Int32(input.PageIndex),
		PageSize:  tea.Int32(input.PageSize),
	}
	resp, err := c.dysmsapi.QuerySmsSignListWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if err := checkBody(resp.Body); err != nil {
		return nil, err
	}

	list := &SignList{TotalCount: tea.Int64Value(resp.Body.TotalCount)}
	list.Items = make([]*SignItem, 0, len(resp.Body.SmsSignList))
	for _, item := range resp.Body.SmsSignList {
		sign := &SignItem{
			SignName:    tea.StringValue(item.SignName),
			AuditStatus: tea.StringValue(item.AuditStatus),
			CreateDate:  tea.StringValue(item.CreateDate),
			OrderID:     tea.StringValue(item.OrderId),
		}
		if item.Reason != nil {
			sign.Reason = tea.StringValue(item.Reason.RejectInfo)
		}
		list.Items = append(list.Items, sign)
	}
	return list, nil
}
