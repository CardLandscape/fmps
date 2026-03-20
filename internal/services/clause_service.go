package services

import (
	"database/sql"
	"fmps/internal/models"
)

type ClauseService struct {
	db *sql.DB
}

func NewClauseService(db *sql.DB) *ClauseService {
	return &ClauseService{db: db}
}

func (s *ClauseService) GetAllClauses() ([]*models.Clause, error) {
	rows, err := s.db.Query(`SELECT id, code, title, COALESCE(description,''), severity, COALESCE(category,''), is_active, created_at FROM clauses ORDER BY severity DESC, code ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clauses []*models.Clause
	for rows.Next() {
		c := &models.Clause{}
		if err := rows.Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.Severity, &c.Category, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		clauses = append(clauses, c)
	}
	if clauses == nil {
		clauses = []*models.Clause{}
	}
	return clauses, nil
}

func (s *ClauseService) GetClause(id int64) (*models.Clause, error) {
	c := &models.Clause{}
	err := s.db.QueryRow(`SELECT id, code, title, COALESCE(description,''), severity, COALESCE(category,''), is_active, created_at FROM clauses WHERE id=?`, id).
		Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.Severity, &c.Category, &c.IsActive, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ClauseService) CreateClause(clause *models.Clause) (*models.Clause, error) {
	result, err := s.db.Exec(`INSERT INTO clauses (code, title, description, severity, category) VALUES (?,?,?,?,?)`,
		clause.Code, clause.Title, clause.Description, clause.Severity, clause.Category)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return s.GetClause(id)
}

func (s *ClauseService) UpdateClause(clause *models.Clause) (*models.Clause, error) {
	_, err := s.db.Exec(`UPDATE clauses SET code=?, title=?, description=?, severity=?, category=?, is_active=? WHERE id=?`,
		clause.Code, clause.Title, clause.Description, clause.Severity, clause.Category, clause.IsActive, clause.ID)
	if err != nil {
		return nil, err
	}
	return s.GetClause(clause.ID)
}

func (s *ClauseService) DeleteClause(id int64) error {
	_, err := s.db.Exec(`DELETE FROM clauses WHERE id=?`, id)
	return err
}

func (s *ClauseService) GetClausesByCategory(category string) ([]*models.Clause, error) {
	rows, err := s.db.Query(`SELECT id, code, title, COALESCE(description,''), severity, COALESCE(category,''), is_active, created_at FROM clauses WHERE category=? ORDER BY severity DESC`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clauses []*models.Clause
	for rows.Next() {
		c := &models.Clause{}
		if err := rows.Scan(&c.ID, &c.Code, &c.Title, &c.Description, &c.Severity, &c.Category, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		clauses = append(clauses, c)
	}
	if clauses == nil {
		clauses = []*models.Clause{}
	}
	return clauses, nil
}
