package models

type Setting struct {
	ID        int64  `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at"`
}
