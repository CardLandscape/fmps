package services

import (
	"database/sql"
	"fmps/internal/models"
)

type SettingService struct {
	db *sql.DB
}

func NewSettingService(db *sql.DB) *SettingService {
	return &SettingService{db: db}
}

func (s *SettingService) GetSetting(key string) (string, error) {
	var value string
	err := s.db.QueryRow(`SELECT COALESCE(value,'') FROM settings WHERE key=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *SettingService) SetSetting(key, value string) error {
	_, err := s.db.Exec(`INSERT INTO settings (key, value) VALUES (?,?) ON CONFLICT(key) DO UPDATE SET value=excluded.value, updated_at=CURRENT_TIMESTAMP`, key, value)
	return err
}

func (s *SettingService) GetAllSettings() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, COALESCE(value,'') FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		result[k] = v
	}
	return result, nil
}

func (s *SettingService) IsInitialized() bool {
	val, _ := s.GetSetting("initialized")
	return val == "true"
}

func (s *SettingService) MarkInitialized() error {
	return s.SetSetting("initialized", "true")
}

var _ = models.Setting{}
