package models

type PunishmentTemplate struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	ClauseID        *int64 `json:"clause_id"`
	PunishmentType  string `json:"punishment_type"`
	DurationMinutes int    `json:"duration_minutes"`
	Description     string `json:"description"`
	IsActive        bool   `json:"is_active"`
	CreatedAt       string `json:"created_at"`
}
