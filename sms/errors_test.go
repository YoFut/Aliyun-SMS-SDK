// 文件类型：错误链测试；用途：检查业务错误分类、供应商详情及本地传输错误。
// 用法：执行 go test ./sms；测试仅访问本地服务器，使用合成凭据。
package sms

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea/tea"
)

// TestBusinessErrorDetails 检查分类转换后仍可读取供应商错误码和消息。
func TestBusinessErrorDetails(t *testing.T) {
	for _, sample := range []struct {
		code string
		want error
	}{
		{"isv.MOBILE_NUMBER_ILLEGAL", ErrPhoneNumbersInvalid},
		{"isv.SMS_SIGN_ILLEGAL", ErrSignNotFound},
		{"isv.SMS_TEMPLATE_ILLEGAL", ErrTemplateNotFound},
		{"isv.BUSINESS_LIMIT_CONTROL", ErrDailyQuotaExceeded},
		{"isv.SMS_SIGNATURE_SCENE_ILLEGAL", ErrSignatureNotMatch},
		{"", ErrService},
	} {
		err := checkBody(&dysmsapi.SendSmsResponseBody{Code: tea.String(sample.code), Message: tea.String("供应商诊断"), RequestId: tea.String("request-id")})
		var detail *Error
		if !errors.Is(err, sample.want) || !errors.As(err, &detail) || detail.Code != sample.code || detail.Message != "供应商诊断" || detail.RequestID != "request-id" {
			t.Fatalf("错误详情或分类丢失: %v", err)
		}
		encoded, _ := json.Marshal(detail)
		if string(encoded) != "{}" {
			t.Fatal("诊断错误不得被意外序列化")
		}
	}
}

type localTransport struct {
	call func(*http.Request) (*http.Response, error)
}

// Call 注入本地 HTTP 行为，避免请求真实短信服务。
func (transport localTransport) Call(request *http.Request, _ *http.Transport) (*http.Response, error) {
	return transport.call(request)
}

// TestSendErrorChain 验证真实官方 SDK 调用路径保留传输及业务错误。
func TestSendErrorChain(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		io.WriteString(writer, `{"Code":"isv.SMS_TEMPLATE_ILLEGAL","Message":"模板未审核","RequestId":"local-request"}`)
	}))
	defer server.Close()
	underlying, err := dysmsapi.NewClient(&openapi.Config{Endpoint: tea.String(strings.TrimPrefix(server.URL, "https://")), AccessKeyId: tea.String("test-id"), AccessKeySecret: tea.String("test-secret")})
	if err != nil {
		t.Fatal(err)
	}
	underlying.HttpClient = localTransport{call: server.Client().Do}
	client, err := NewClient(WithClient(underlying))
	if err != nil {
		t.Fatal(err)
	}
	input := &SendSMSInput{PhoneNumbers: "13800000000", SignName: "测试", TemplateCode: "test", TemplateParam: `{"code":"000000"}`}
	_, err = client.SendSMS(context.Background(), input)
	var detail *Error
	if !errors.Is(err, ErrTemplateNotFound) || !errors.As(err, &detail) || detail.RequestID != "local-request" {
		t.Fatalf("实际 SDK 路径丢失业务详情: %v", err)
	}
	original := errors.New("本地模拟网络故障")
	underlying.HttpClient = localTransport{call: func(*http.Request) (*http.Response, error) { return nil, original }}
	_, err = client.SendSMS(context.Background(), input)
	if !errors.Is(err, original) {
		t.Fatalf("传输根因被替换: %v", err)
	}
}
