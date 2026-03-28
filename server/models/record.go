package models

import "time"

type Record struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MemberID   uint      `json:"member_id"`
	Member     Member    `gorm:"foreignKey:MemberID" json:"member"`
	RuleID     uint      `json:"rule_id"`
	Rule       Rule      `gorm:"foreignKey:RuleID" json:"rule"`
	Points     int       `json:"points"`
	Note       string    `json:"note"`
	OccurredAt time.Time `json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
}
