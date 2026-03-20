package services

import (
	"database/sql"
	"fmps/internal/models"
	"time"
)

type AppealService struct {
	db *sql.DB
}

func NewAppealService(db *sql.DB) *AppealService {
	return &AppealService{db: db}
}

func (s *AppealService) GetAppealsByCase(caseId int64) ([]*models.Appeal, error) {
	rows, err := s.db.Query(`SELECT a.id, a.case_id, a.appellant_id, a.reason, a.status, a.reviewer_id, COALESCE(a.review_comment,''), a.created_at, a.reviewed_at, COALESCE(m.name,'') FROM appeals a LEFT JOIN members m ON a.appellant_id=m.id WHERE a.case_id=? ORDER BY a.created_at DESC`, caseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var appeals []*models.Appeal
	for rows.Next() {
		a := &models.Appeal{}
		rows.Scan(&a.ID, &a.CaseID, &a.AppellantID, &a.Reason, &a.Status, &a.ReviewerID, &a.ReviewComment, &a.CreatedAt, &a.ReviewedAt, &a.AppellantName)
		appeals = append(appeals, a)
	}
	if appeals == nil {
		appeals = []*models.Appeal{}
	}
	return appeals, nil
}

func (s *AppealService) CreateAppeal(appeal *models.Appeal) (*models.Appeal, error) {
	result, err := s.db.Exec(`INSERT INTO appeals (case_id, appellant_id, reason) VALUES (?,?,?)`,
		appeal.CaseID, appeal.AppellantID, appeal.Reason)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	s.db.Exec(`UPDATE cases SET status='appealed', updated_at=CURRENT_TIMESTAMP WHERE id=?`, appeal.CaseID)
	a := &models.Appeal{}
	s.db.QueryRow(`SELECT a.id, a.case_id, a.appellant_id, a.reason, a.status, a.reviewer_id, COALESCE(a.review_comment,''), a.created_at, a.reviewed_at, COALESCE(m.name,'') FROM appeals a LEFT JOIN members m ON a.appellant_id=m.id WHERE a.id=?`, id).
		Scan(&a.ID, &a.CaseID, &a.AppellantID, &a.Reason, &a.Status, &a.ReviewerID, &a.ReviewComment, &a.CreatedAt, &a.ReviewedAt, &a.AppellantName)
	return a, nil
}

func (s *AppealService) ReviewAppeal(appealId int64, reviewerId int64, approved bool, comment string) error {
	status := "rejected"
	if approved {
		status = "approved"
	}
	reviewedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.Exec(`UPDATE appeals SET status=?, reviewer_id=?, review_comment=?, reviewed_at=? WHERE id=?`,
		status, reviewerId, comment, reviewedAt, appealId)
	if err != nil {
		return err
	}
	if approved {
		var caseId int64
		s.db.QueryRow(`SELECT case_id FROM appeals WHERE id=?`, appealId).Scan(&caseId)
		s.db.Exec(`UPDATE cases SET status='cancelled', updated_at=CURRENT_TIMESTAMP WHERE id=?`, caseId)
	}
	return nil
}
