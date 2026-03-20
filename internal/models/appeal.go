package models

type Appeal struct {
	ID            int64   `json:"id"`
	CaseID        int64   `json:"case_id"`
	AppellantID   int64   `json:"appellant_id"`
	Reason        string  `json:"reason"`
	Status        string  `json:"status"`
	ReviewerID    *int64  `json:"reviewer_id"`
	ReviewComment string  `json:"review_comment"`
	CreatedAt     string  `json:"created_at"`
	ReviewedAt    *string `json:"reviewed_at"`
	AppellantName string  `json:"appellant_name"`
}
