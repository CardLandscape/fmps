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
	PunishmentProcess string     `json:"punishment_process"` // 原始文本，管道符分隔（兼容旧格式）
	// Punishment workflow fields
	PunishmentLevel   string     `json:"punishment_level"`    // A/B/C/D
	PrepItems         string     `json:"prep_items"`          // JSON array of preparation item strings
	ParsedSteps       string     `json:"parsed_steps"`        // JSON array of step strings
	// CurrentStepIndex tracks execution progress:
	//   -1 = preparation phase (items must be checked before starting)
	//   0+ = index of the currently active execution step
	CurrentStepIndex  int        `json:"current_step_index"`
	FinalGrade        string     `json:"final_grade"`         // 满分/优/良/达标/不达标/态度不端正
	TxtFilename       string     `json:"txt_filename"`        // original filename of imported TXT
	Status            string     `json:"status"`              // pending / active / completed
	StartTime         *time.Time `json:"start_time"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
