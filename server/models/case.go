package models

import "time"

type Case struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	MemberID          uint       `json:"member_id"`           // legacy field, kept for backward compat
	Member            Member     `gorm:"foreignKey:MemberID" json:"member"`
	ParentMemberID    uint       `json:"parent_member_id"`    // 家长成员ID
	ChildMemberID     uint       `json:"child_member_id"`     // 小孩成员ID
	ParentMember      *Member    `gorm:"foreignKey:ParentMemberID" json:"parent_member"`
	ChildMember       *Member    `gorm:"foreignKey:ChildMemberID" json:"child_member"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	PunishmentProcess string     `json:"punishment_process"` // 原始文本，管道符分隔（旧格式兼容）
	// 新惩罚流程字段
	PunishmentLevel string `json:"punishment_level"` // 惩罚级别：A/B/C/D
	PrepItems       string `json:"prep_items"`       // JSON 数组：准备物品列表
	ExecSteps       string `json:"exec_steps"`       // JSON 数组：执行步骤列表
	CurrentStep     int    `json:"current_step" gorm:"default:-1"` // 当前步骤索引（-1=准备阶段）
	FinalGrade      string `json:"final_grade"`      // 最终成绩：满分/优/良/达标/不达标/态度不端正
	Status          string `json:"status"`           // pending / active / completed
	StartTime       *time.Time `json:"start_time"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
