package services

import (
	"database/sql"
	"fmps/internal/models"
)

type TemplateService struct {
	db *sql.DB
}

func NewTemplateService(db *sql.DB) *TemplateService {
	return &TemplateService{db: db}
}

func (s *TemplateService) GetAllTemplates() ([]*models.PunishmentTemplate, error) {
	rows, err := s.db.Query(`SELECT id, name, clause_id, punishment_type, duration_minutes, COALESCE(description,''), is_active, created_at FROM punishment_templates ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var templates []*models.PunishmentTemplate
	for rows.Next() {
		t := &models.PunishmentTemplate{}
		if err := rows.Scan(&t.ID, &t.Name, &t.ClauseID, &t.PunishmentType, &t.DurationMinutes, &t.Description, &t.IsActive, &t.CreatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	if templates == nil {
		templates = []*models.PunishmentTemplate{}
	}
	return templates, nil
}

func (s *TemplateService) GetTemplate(id int64) (*models.PunishmentTemplate, error) {
	t := &models.PunishmentTemplate{}
	err := s.db.QueryRow(`SELECT id, name, clause_id, punishment_type, duration_minutes, COALESCE(description,''), is_active, created_at FROM punishment_templates WHERE id=?`, id).
		Scan(&t.ID, &t.Name, &t.ClauseID, &t.PunishmentType, &t.DurationMinutes, &t.Description, &t.IsActive, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TemplateService) CreateTemplate(template *models.PunishmentTemplate) (*models.PunishmentTemplate, error) {
	result, err := s.db.Exec(`INSERT INTO punishment_templates (name, clause_id, punishment_type, duration_minutes, description) VALUES (?,?,?,?,?)`,
		template.Name, template.ClauseID, template.PunishmentType, template.DurationMinutes, template.Description)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return s.GetTemplate(id)
}

func (s *TemplateService) UpdateTemplate(template *models.PunishmentTemplate) (*models.PunishmentTemplate, error) {
	_, err := s.db.Exec(`UPDATE punishment_templates SET name=?, clause_id=?, punishment_type=?, duration_minutes=?, description=?, is_active=? WHERE id=?`,
		template.Name, template.ClauseID, template.PunishmentType, template.DurationMinutes, template.Description, template.IsActive, template.ID)
	if err != nil {
		return nil, err
	}
	return s.GetTemplate(template.ID)
}

func (s *TemplateService) DeleteTemplate(id int64) error {
	_, err := s.db.Exec(`DELETE FROM punishment_templates WHERE id=?`, id)
	return err
}
