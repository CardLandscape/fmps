package models

import "time"

type Member struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	NameCn    string    `json:"name_cn"` // 中文姓名
	NameEn    string    `json:"name_en"` // 英文姓名
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`

	// 基本信息
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	BirthDate   string `json:"birth_date"`

	// 主要证件
	IdDocType        string `json:"id_doc_type"`
	IdDocNumber      string `json:"id_doc_number"`
	IdIssueDate      string `json:"id_issue_date"`
	IdExpiryDate     string `json:"id_expiry_date"`
	IdIssueAuthority string `json:"id_issue_authority"`

	// 辅助证件（旧字段，保留兼容）
	AuxDocType   string `json:"aux_doc_type"`
	AuxDocNumber string `json:"aux_doc_number"`

	// 辅助证件1
	Aux1DocType   string `json:"aux1_doc_type"`
	Aux1DocNumber string `json:"aux1_doc_number"`

	// 辅助证件2
	Aux2DocType   string `json:"aux2_doc_type"`
	Aux2DocNumber string `json:"aux2_doc_number"`

	// 主证件为04时的补充字段
	ProofDocType      string `json:"proof_doc_type"`       // 证明文件类型（94RV/94PV/94PC/94PE/94NP）
	ProofIssueCountry string `json:"proof_issue_country"`  // 签发国家（ISO alpha-3）

	// 学籍信息
	SchoolName        string `json:"school_name"`
	Grade             string `json:"grade"`
	ClassName         string `json:"class_name"`
	ClassTeacherName  string `json:"class_teacher_name"`
	ClassTeacherPhone string `json:"class_teacher_phone"`

	// 外出权限
	OutingPermission string `json:"outing_permission"` // 许可/不许可/受限
	OutingDates      string `json:"outing_dates"`      // JSON array
	OutingTimeRanges string `json:"outing_time_ranges"` // JSON array

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
