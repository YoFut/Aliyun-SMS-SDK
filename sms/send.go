// 本文件实现短信发送与发送明细查询能力。
package sms

import (
	"context"
	"encoding/json"
	"strconv"

	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	service "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// SendSMS 向指定号码发送一条短信。
//
// 参数：
//   - ctx：请求上下文，用于超时与取消
//   - input：发送入参，PhoneNumbers、SignName、TemplateCode 为必填
//
// 返回：
//   - *SendSMSReceipt：发送回执，BizID 可用于查询发送状态
//   - error：传输或业务失败时返回，业务失败按阿里云错误码映射为哨兵错误
func (c *Client) SendSMS(ctx context.Context, input *SendSMSInput) (*SendSMSReceipt, error) {
	req := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(input.PhoneNumbers),
		SignName:      tea.String(input.SignName),
		TemplateCode:  tea.String(input.TemplateCode),
		TemplateParam: tea.String(input.TemplateParam),
	}
	if input.OutID != "" {
		req.OutId = tea.String(input.OutID)
	}

	return c.sendSMS(ctx, req)
}

// sendSMS 发送单条短信并返回回执，供 SendSMS 与 SendVerificationCode 复用。
func (c *Client) sendSMS(ctx context.Context, req *dysmsapi.SendSmsRequest) (*SendSMSReceipt, error) {
	resp, err := c.dysmsapi.SendSmsWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if err := checkBody(resp.Body); err != nil {
		return nil, err
	}
	return &SendSMSReceipt{
		BizID:     tea.StringValue(resp.Body.BizId),
		RequestID: tea.StringValue(resp.Body.RequestId),
	}, nil
}

// SendBatchSMS 批量发送短信，各发送项按下标与号码一一对应。
//
// 参数：
//   - ctx：请求上下文
//   - input：群发入参，TemplateCode 与至少一个 Items 为必填
//
// 返回：
//   - *SendBatchSMSReceipt：发送回执
//   - error：传输或业务失败时返回
func (c *Client) SendBatchSMS(ctx context.Context, input *SendBatchSMSInput) (*SendBatchSMSReceipt, error) {
	numbers := make([]string, 0, len(input.Items))
	signs := make([]string, 0, len(input.Items))
	params := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		numbers = append(numbers, item.PhoneNumber)
		signs = append(signs, item.SignName)
		params = append(params, item.TemplateParam)
	}

	req := &dysmsapi.SendBatchSmsRequest{
		PhoneNumberJson:   tea.String(marshalStrings(numbers)),
		SignNameJson:      tea.String(marshalStrings(signs)),
		TemplateParamJson: tea.String(marshalStrings(params)),
		TemplateCode:      tea.String(input.TemplateCode),
	}
	if input.OutID != "" {
		req.OutId = tea.String(input.OutID)
	}

	resp, err := c.dysmsapi.SendBatchSmsWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if err := checkBody(resp.Body); err != nil {
		return nil, err
	}
	return &SendBatchSMSReceipt{
		BizID:     tea.StringValue(resp.Body.BizId),
		RequestID: tea.StringValue(resp.Body.RequestId),
	}, nil
}

// QuerySendDetails 查询指定号码的发送明细与送达状态。
//
// 参数：
//   - ctx：请求上下文
//   - input：查询入参，PhoneNumber、SendDate 为必填；CurrentPage、PageSize 不填时默认第 1 页每页 20 条
//
// 返回：
//   - *SendDetailList：发送明细集合
//   - error：传输或业务失败时返回
func (c *Client) QuerySendDetails(ctx context.Context, input *QuerySendDetailsInput) (*SendDetailList, error) {
	currentPage := input.CurrentPage
	if currentPage == 0 {
		currentPage = 1
	}
	pageSize := input.PageSize
	if pageSize == 0 {
		pageSize = 20
	}

	req := &dysmsapi.QuerySendDetailsRequest{
		PhoneNumber: tea.String(input.PhoneNumber),
		SendDate:    tea.String(input.SendDate),
		CurrentPage: tea.Int64(currentPage),
		PageSize:    tea.Int64(pageSize),
	}
	if input.BizID != "" {
		req.BizId = tea.String(input.BizID)
	}

	resp, err := c.dysmsapi.QuerySendDetailsWithContext(ctx, req, &service.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if err := checkBody(resp.Body); err != nil {
		return nil, err
	}

	list := &SendDetailList{}
	if dto := resp.Body.SmsSendDetailDTOs; dto != nil {
		list.Items = make([]*SendDetail, 0, len(dto.SmsSendDetailDTO))
		for _, d := range dto.SmsSendDetailDTO {
			list.Items = append(list.Items, &SendDetail{
				Content:      tea.StringValue(d.Content),
				ErrCode:      tea.StringValue(d.ErrCode),
				OutID:        tea.StringValue(d.OutId),
				PhoneNumber:  tea.StringValue(d.PhoneNum),
				ReceiveDate:  tea.StringValue(d.ReceiveDate),
				SendDate:     tea.StringValue(d.SendDate),
				SendStatus:   tea.Int64Value(d.SendStatus),
				TemplateCode: tea.StringValue(d.TemplateCode),
			})
		}
	}
	list.TotalCount, _ = strconv.ParseInt(tea.StringValue(resp.Body.TotalCount), 10, 64)
	return list, nil
}

// marshalStrings 将字符串切片序列化为 JSON 数组串，供阿里云 JSON 数组参数使用。
func marshalStrings(values []string) string {
	b, _ := json.Marshal(values)
	return string(b)
}
