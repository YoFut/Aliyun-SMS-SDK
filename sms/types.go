// 本文件定义短信 SDK 对外暴露的入参、回执与实体结构。
package sms

// SendSMSInput 单发短信入参。
type SendSMSInput struct {
	// PhoneNumbers 接收号码，支持逗号分隔（最多 1000 个），可带 +86 前缀。
	PhoneNumbers string
	// SignName 短信签名名称，须已审核通过。
	SignName string
	// TemplateCode 短信模板编号，须已审核通过。
	TemplateCode string
	// TemplateParam 模板变量 JSON 串，模板无变量时留空。
	TemplateParam string
	// OutID 外部业务标识，用于对账，可留空。
	OutID string
}

// SendSMSReceipt 单发短信回执。
type SendSMSReceipt struct {
	// BizID 发送回执编号，用于后续查询发送状态。
	BizID string
	// RequestID 阿里云请求编号，用于问题排查。
	RequestID string
}

// BatchSMSItem 群发短信中的一条发送项。
type BatchSMSItem struct {
	// PhoneNumber 接收号码。
	PhoneNumber string
	// SignName 签名名称，须与模板匹配。
	SignName string
	// TemplateParam 该号码的模板变量 JSON 串，无变量时留空。
	TemplateParam string
}

// SendBatchSMSInput 群发短信入参，各发送项按下标与号码一一对应。
type SendBatchSMSInput struct {
	// TemplateCode 短信模板编号。
	TemplateCode string
	// Items 发送项列表。
	Items []*BatchSMSItem
	// OutID 外部业务标识，可留空。
	OutID string
}

// SendBatchSMSReceipt 群发短信回执。
type SendBatchSMSReceipt struct {
	// BizID 发送回执编号。
	BizID string
	// RequestID 阿里云请求编号。
	RequestID string
}

// QuerySendDetailsInput 查询发送明细入参。
type QuerySendDetailsInput struct {
	// PhoneNumber 要查询的号码，单次只能查一个。
	PhoneNumber string
	// SendDate 发送日期，格式 yyyyMMdd，可查最近 30 天。
	SendDate string
	// CurrentPage 页码，从 1 开始，不填时默认 1。
	CurrentPage int64
	// PageSize 每页条数，1~50，不填时默认 20。
	PageSize int64
	// BizID 发送回执编号，可选，用于按回执筛选。
	BizID string
}

// SendDetail 单条发送明细。
type SendDetail struct {
	// Content 短信内容。
	Content string
	// ErrCode 运营商返回的错误码。
	ErrCode string
	// OutID 外部业务标识。
	OutID string
	// PhoneNumber 接收号码。
	PhoneNumber string
	// ReceiveDate 接收时间。
	ReceiveDate string
	// SendDate 发送时间。
	SendDate string
	// SendStatus 发送状态：1 成功，2 失败，3 处理中。
	SendStatus int64
	// TemplateCode 模板编号。
	TemplateCode string
}

// SendDetailList 发送明细集合。
type SendDetailList struct {
	// Items 明细列表。
	Items []*SendDetail
	// TotalCount 总条数。
	TotalCount int64
}

// SignFile 签名资质文件。
type SignFile struct {
	// FileContents base64 编码的文件内容。
	FileContents string
	// FileSuffix 文件后缀，如 jpg、png。
	FileSuffix string
}

// AddSignInput 新增签名入参。
type AddSignInput struct {
	// SignName 签名名称。
	SignName string
	// SignSource 签名来源：0 企事业单位，1 网站，2 APP，3 小程序，4 商标，5 政府机关事业单位。
	SignSource int32
	// SignType 签名类型：0 验证码，1 通知，2 推广。
	SignType int32
	// Remark 申请说明。
	Remark string
	// SignFiles 资质文件列表。
	SignFiles []*SignFile
}

// ModifySignInput 修改签名入参，字段含义同 AddSignInput。
type ModifySignInput struct {
	// SignName 签名名称。
	SignName string
	// SignSource 签名来源。
	SignSource int32
	// SignType 签名类型。
	SignType int32
	// Remark 申请说明。
	Remark string
	// SignFiles 资质文件列表。
	SignFiles []*SignFile
}

// DeleteSignInput 删除签名入参。
type DeleteSignInput struct {
	// SignName 签名名称。
	SignName string
}

// QuerySignInput 查询单个签名入参。
type QuerySignInput struct {
	// SignName 签名名称。
	SignName string
}

// Sign 短信签名信息。
type Sign struct {
	// SignName 签名名称。
	SignName string
	// CreateDate 创建时间。
	CreateDate string
	// SignStatus 审核状态。
	SignStatus int32
	// Reason 审核说明，驳回时为驳回原因。
	Reason string
}

// ListSignsInput 分页查询签名列表入参。
type ListSignsInput struct {
	// PageIndex 页码，从 1 开始。
	PageIndex int32
	// PageSize 每页条数。
	PageSize int32
}

// SignItem 签名列表项。
type SignItem struct {
	// SignName 签名名称。
	SignName string
	// AuditStatus 审核状态。
	AuditStatus string
	// CreateDate 创建时间。
	CreateDate string
	// OrderID 工单编号。
	OrderID string
	// Reason 驳回原因。
	Reason string
}

// SignList 签名列表集合。
type SignList struct {
	// Items 签名列表。
	Items []*SignItem
	// TotalCount 总数。
	TotalCount int64
}

// AddTemplateInput 新增模板入参。
type AddTemplateInput struct {
	// TemplateName 模板名称。
	TemplateName string
	// TemplateType 模板类型：0 验证码，1 通知，2 推广，3 国际。
	TemplateType int32
	// TemplateContent 模板内容，变量占位形如 ${code}。
	TemplateContent string
	// Remark 申请说明。
	Remark string
}

// ModifyTemplateInput 修改模板入参。
type ModifyTemplateInput struct {
	// TemplateCode 模板编号。
	TemplateCode string
	// TemplateName 模板名称。
	TemplateName string
	// TemplateType 模板类型。
	TemplateType int32
	// TemplateContent 模板内容。
	TemplateContent string
	// Remark 申请说明。
	Remark string
}

// DeleteTemplateInput 删除模板入参。
type DeleteTemplateInput struct {
	// TemplateCode 模板编号。
	TemplateCode string
}

// QueryTemplateInput 查询单个模板入参。
type QueryTemplateInput struct {
	// TemplateCode 模板编号。
	TemplateCode string
}

// Template 短信模板信息。
type Template struct {
	// TemplateCode 模板编号。
	TemplateCode string
	// TemplateName 模板名称。
	TemplateName string
	// TemplateType 模板类型。
	TemplateType int32
	// TemplateContent 模板内容。
	TemplateContent string
	// TemplateStatus 审核状态。
	TemplateStatus int32
	// CreateDate 创建时间。
	CreateDate string
	// Reason 审核说明。
	Reason string
}

// ListTemplatesInput 分页查询模板列表入参。
type ListTemplatesInput struct {
	// PageIndex 页码，从 1 开始。
	PageIndex int32
	// PageSize 每页条数。
	PageSize int32
}

// TemplateItem 模板列表项。
type TemplateItem struct {
	// TemplateCode 模板编号。
	TemplateCode string
	// TemplateName 模板名称。
	TemplateName string
	// TemplateType 模板类型。
	TemplateType int32
	// TemplateContent 模板内容。
	TemplateContent string
	// AuditStatus 审核状态。
	AuditStatus string
	// CreateDate 创建时间。
	CreateDate string
	// Reason 驳回原因。
	Reason string
}

// TemplateList 模板列表集合。
type TemplateList struct {
	// Items 模板列表。
	Items []*TemplateItem
	// TotalCount 总数。
	TotalCount int64
}

// SendVerificationCodeInput 发送验证码入参。
type SendVerificationCodeInput struct {
	// PhoneNumbers 接收号码，验证码场景通常为单个号码。
	PhoneNumbers string
	// SignName 短信签名名称。
	SignName string
	// TemplateCode 短信模板编号，模板须含 code 变量。
	TemplateCode string
	// CodeLength 验证码长度，不填时默认 6 位。
	CodeLength int
	// OutID 外部业务标识，可留空。
	OutID string
}

// VerifyVerificationCodeInput 核验验证码入参。
type VerifyVerificationCodeInput struct {
	// PhoneNumber 接收号码，须与发送时一致。
	PhoneNumber string
	// Code 用户提交的验证码。
	Code string
}
