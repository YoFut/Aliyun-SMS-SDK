// 本文件实现短信验证码的生成、发送与核验能力。
package sms

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/json"
	"math/big"
	"sync"
	"time"

	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea/tea"
)

// SendVerificationCode 生成并发送一条短信验证码，验证码由本包生成并暂存，供 VerifyVerificationCode 核验。
//
// 参数：
//   - ctx：请求上下文
//   - input：发送入参，PhoneNumbers、SignName、TemplateCode 为必填；CodeLength 不填时默认 6
//
// 返回：
//   - *SendSMSReceipt：发送回执
//   - error：验证码生成、传输或业务失败时返回
func (c *Client) SendVerificationCode(ctx context.Context, input *SendVerificationCodeInput) (*SendSMSReceipt, error) {
	code, err := generateCode(input.CodeLength)
	if err != nil {
		return nil, err
	}
	templateParam, err := json.Marshal(map[string]string{"code": code})
	if err != nil {
		return nil, err
	}

	req := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(input.PhoneNumbers),
		SignName:      tea.String(input.SignName),
		TemplateCode:  tea.String(input.TemplateCode),
		TemplateParam: tea.String(string(templateParam)),
	}
	if input.OutID != "" {
		req.OutId = tea.String(input.OutID)
	}

	receipt, err := c.sendSMS(ctx, req)
	if err != nil {
		return nil, err
	}

	// 发送成功后才暂存，避免发送失败留下无法核验的记录。
	c.codes.set(input.PhoneNumbers, code)

	return receipt, nil
}

// VerifyVerificationCode 核验用户提交的验证码是否通过，命中后即失效（单次使用）。
//
// 号码未发送过、已过期或提交码不匹配时均返回 false。
//
// 参数：
//   - ctx：请求上下文（保留，当前为进程内存核验，未使用）
//   - input：核验入参，PhoneNumber 须与发送时一致
//
// 返回：bool，true 表示验证码匹配且未过期。
func (c *Client) VerifyVerificationCode(_ context.Context, input *VerifyVerificationCodeInput) bool {
	return c.codes.verify(input.PhoneNumber, input.Code)
}

// generateCode 生成指定长度的数字验证码，长度不足时默认 6 位。
func generateCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	b := make([]byte, length)
	for i := range b {
		n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		b[i] = byte('0' + n.Int64())
	}
	return string(b), nil
}

// codeEntry 一条验证码记录。
type codeEntry struct {
	code      string
	expiresAt time.Time
}

// codeStore 内存验证码存储，按号码记录最近一次发送的验证码。
//
// 存储为进程内存，多实例部署时各实例独立，需要跨实例共享时应改用外部存储。
type codeStore struct {
	mu    sync.Mutex
	ttl   time.Duration
	codes map[string]*codeEntry
}

// newCodeStore 构造验证码存储。
func newCodeStore(ttl time.Duration) *codeStore {
	return &codeStore{
		ttl:   ttl,
		codes: make(map[string]*codeEntry),
	}
}

// set 记录某号码的验证码，覆盖旧值。
func (s *codeStore) set(phone, code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[phone] = &codeEntry{code: code, expiresAt: time.Now().Add(s.ttl)}
}

// verify 校验号码与验证码是否匹配，命中后即删除（单次使用）。
func (s *codeStore) verify(phone, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.codes[phone]
	if !ok {
		return false
	}
	if time.Now().After(entry.expiresAt) {
		delete(s.codes, phone)
		return false
	}
	if entry.code != code {
		return false
	}
	delete(s.codes, phone)
	return true
}
