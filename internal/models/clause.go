package models

type Clause struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    int    `json:"severity"`
	Category    string `json:"category"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}
