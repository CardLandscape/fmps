package models

type Case struct {
	ID                  int64   `json:"id"`
	CaseNo              string  `json:"case_no"`
	MemberID            int64   `json:"member_id"`
	ClauseID            int64   `json:"clause_id"`
	TemplateID          *int64  `json:"template_id"`
	Status              string  `json:"status"`
	IncidentDescription string  `json:"incident_description"`
	IncidentTime        *string `json:"incident_time"`
	PunishmentDetail    string  `json:"punishment_detail"`
	StartedAt           *string `json:"started_at"`
	CompletedAt         *string `json:"completed_at"`
	CreatedBy           *int64  `json:"created_by"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	MemberName          string  `json:"member_name"`
	ClauseTitle         string  `json:"clause_title"`
	ClauseCode          string  `json:"clause_code"`
}

type CaseComment struct {
	ID         int64  `json:"id"`
	CaseID     int64  `json:"case_id"`
	MemberID   int64  `json:"member_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	MemberName string `json:"member_name"`
}

type CaseStats struct {
	TotalMembers  int64             `json:"total_members"`
	ActiveClauses int64             `json:"active_clauses"`
	MonthCases    int64             `json:"month_cases"`
	PendingCases  int64             `json:"pending_cases"`
	StatusCounts  map[string]int64  `json:"status_counts"`
	MemberCounts  []MemberCaseCount `json:"member_counts"`
}

type MemberCaseCount struct {
	MemberID   int64  `json:"member_id"`
	MemberName string `json:"member_name"`
	Count      int64  `json:"count"`
}
