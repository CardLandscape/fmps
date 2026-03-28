package models

import "time"

type PenaltyPoint struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CaseID       uint       `json:"case_id"`
	MemberID     uint       `json:"member_id"`
	RuleText     string     `json:"rule_text"`
	ScoreDelta   int        `json:"score_delta"` // 负数
	Reason       string     `json:"reason"`
	CreatedAt    time.Time  `json:"created_at"`
	Revoked      bool       `json:"revoked"`
	RevokedAt    *time.Time `json:"revoked_at"`
	RevokeReason string     `json:"revoke_reason"`
}
