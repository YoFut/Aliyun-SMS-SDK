// 本文件定义错误变量与阿里云业务错误的映射。
package sms

import (
	"errors"
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
)

// 阿里云短信常见业务错误，调用方可用 errors.Is 判断。
var (
	// ErrPhoneNumbersInvalid 表示接收号码非法。
	ErrPhoneNumbersInvalid = errors.New("接收号码非法")
	// ErrSignNotFound 表示签名不存在或未审核通过。
	ErrSignNotFound = errors.New("签名不存在或未审核通过")
	// ErrTemplateNotFound 表示模板不存在或未审核通过。
	ErrTemplateNotFound = errors.New("模板不存在或未审核通过")
	// ErrDailyQuotaExceeded 表示当日发送量已达上限。
	ErrDailyQuotaExceeded = errors.New("当日发送量已达上限")
	// ErrSignatureNotMatch 表示签名与模板场景不匹配。
	ErrSignatureNotMatch = errors.New("签名与模板场景不匹配")
	// ErrService 表示服务响应缺失或未返回有效业务状态。
	ErrService = errors.New("短信服务响应异常")
)

// Error 表示一次阿里云短信调用失败，携带阿里云返回的错误码与信息。
type Error struct {
	// Code 阿里云错误码，如 isv.MOBILE_NUMBER_ILLEGAL。
	Code string `json:"-"`
	// Message 阿里云错误信息。
	Message string `json:"-"`
	// RequestID 用于后端关联供应商请求，不应直接返回浏览器。
	RequestID string `json:"-"`
	// Cause 保留错误分类或底层错误，供 errors.Is 和 errors.As 检查。
	Cause error `json:"-"`
}

// Unwrap 返回错误原因；调用方须在 HTTP 边界另外选择安全提示。
func (e *Error) Unwrap() error { return e.Cause }

// Error 返回错误描述。
func (e *Error) Error() string {
	if e.Code == "" && e.Message == "" && e.Cause != nil {
		return e.Cause.Error()
	}
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprintf("短信服务错误 %s: %s", e.Code, e.Message)
}

// responseBody 抽象阿里云响应体的公共错误字段，供 checkBody 统一校验。
type responseBody interface {
	GetCode() *string
	GetMessage() *string
}

// checkBody 校验阿里云响应体，业务失败时返回映射后的错误。
//
// 阿里云在业务失败时返回 HTTP 200，错误码位于响应体 Code 字段，
// 因此除传输层 err 之外，还需在此校验业务码。
func checkBody(body responseBody) error {
	if body == nil {
		return &Error{Message: "短信服务返回空响应", Cause: ErrService}
	}
	code := tea.StringValue(body.GetCode())
	if code == "OK" {
		return nil
	}
	err := apiError(code, tea.StringValue(body.GetMessage()))
	if receipt, ok := body.(interface{ GetRequestId() *string }); ok {
		err.RequestID = tea.StringValue(receipt.GetRequestId())
	}
	return err
}

// apiError 按阿里云错误码映射为哨兵错误，未识别的码包装为 Error。
func apiError(code, message string) *Error {
	var cause error
	switch code {
	case "":
		cause = ErrService
	case "isv.MOBILE_NUMBER_ILLEGAL":
		cause = ErrPhoneNumbersInvalid
	case "isv.SMS_SIGN_ILLEGAL", "isv.SMS_SIGNATURE_ILLEGAL":
		cause = ErrSignNotFound
	case "isv.SMS_TEMPLATE_ILLEGAL":
		cause = ErrTemplateNotFound
	case "isv.DAY_LIMIT_CONTROL", "isv.BUSINESS_LIMIT_CONTROL":
		cause = ErrDailyQuotaExceeded
	case "isv.SMS_SIGNATURE_SCENE_ILLEGAL":
		cause = ErrSignatureNotMatch
	}
	return &Error{Code: code, Message: message, Cause: cause}
}
