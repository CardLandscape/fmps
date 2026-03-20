package models

type Member struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Role           string  `json:"role"`
	Avatar         string  `json:"avatar"`
	IsProtected    bool    `json:"is_protected"`
	ProtectedUntil *string `json:"protected_until"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
