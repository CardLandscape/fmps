package models

import "time"

type Case struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	MemberID          uint       `json:"member_id"`
	Member            Member     `gorm:"foreignKey:MemberID" json:"member"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	PunishmentProcess string     `json:"punishment_process"` // 原始文本，管道符分隔
	Status            string     `json:"status"`             // pending / active / completed
	StartTime         *time.Time `json:"start_time"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
