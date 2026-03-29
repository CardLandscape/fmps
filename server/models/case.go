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
	PunishmentProcess string     `json:"punishment_process"` // 原始文本，管道符分隔
	Status            string     `json:"status"`             // pending / active / completed
	StartTime         *time.Time `json:"start_time"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
