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

	// 辅助证件
	AuxDocType   string `json:"aux_doc_type"`
	AuxDocNumber string `json:"aux_doc_number"`

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
