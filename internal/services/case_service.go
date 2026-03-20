package services

import (
	"database/sql"
	"errors"
	"fmt"
	"fmps/internal/models"
	"time"
)

type CaseService struct {
	db *sql.DB
}

func NewCaseService(db *sql.DB) *CaseService {
	return &CaseService{db: db}
}

func (s *CaseService) generateCaseNo() string {
	today := time.Now().Format("20060102")
	var count int64
	s.db.QueryRow(`SELECT COUNT(*) FROM cases WHERE case_no LIKE ?`, fmt.Sprintf("CASE-%s-%%", today)).Scan(&count)
	return fmt.Sprintf("CASE-%s-%03d", today, count+1)
}

func (s *CaseService) scanCase(row *sql.Row) (*models.Case, error) {
	c := &models.Case{}
	err := row.Scan(&c.ID, &c.CaseNo, &c.MemberID, &c.ClauseID, &c.TemplateID, &c.Status,
		&c.IncidentDescription, &c.IncidentTime, &c.PunishmentDetail,
		&c.StartedAt, &c.CompletedAt, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
		&c.MemberName, &c.ClauseTitle, &c.ClauseCode)
	return c, err
}

const caseQuery = `SELECT c.id, c.case_no, c.member_id, c.clause_id, c.template_id, c.status,
    COALESCE(c.incident_description,''), c.incident_time, COALESCE(c.punishment_detail,''),
    c.started_at, c.completed_at, c.created_by, c.created_at, c.updated_at,
    COALESCE(m.name,''), COALESCE(cl.title,''), COALESCE(cl.code,'')
    FROM cases c
    LEFT JOIN members m ON c.member_id = m.id
    LEFT JOIN clauses cl ON c.clause_id = cl.id`

func (s *CaseService) GetAllCases() ([]*models.Case, error) {
	rows, err := s.db.Query(caseQuery + ` ORDER BY c.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanCases(rows)
}

func (s *CaseService) scanCases(rows *sql.Rows) ([]*models.Case, error) {
	var cases []*models.Case
	for rows.Next() {
		c := &models.Case{}
		err := rows.Scan(&c.ID, &c.CaseNo, &c.MemberID, &c.ClauseID, &c.TemplateID, &c.Status,
			&c.IncidentDescription, &c.IncidentTime, &c.PunishmentDetail,
			&c.StartedAt, &c.CompletedAt, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
			&c.MemberName, &c.ClauseTitle, &c.ClauseCode)
		if err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	if cases == nil {
		cases = []*models.Case{}
	}
	return cases, nil
}

func (s *CaseService) GetCase(id int64) (*models.Case, error) {
	row := s.db.QueryRow(caseQuery+` WHERE c.id=?`, id)
	return s.scanCase(row)
}

func (s *CaseService) CreateCase(c *models.Case) (*models.Case, error) {
	var isProtected bool
	var protectedUntil *string
	s.db.QueryRow(`SELECT is_protected, protected_until FROM members WHERE id=?`, c.MemberID).
		Scan(&isProtected, &protectedUntil)
	if isProtected && protectedUntil != nil {
		until, _ := time.Parse("2006-01-02 15:04:05", *protectedUntil)
		if time.Now().Before(until) {
			return nil, errors.New("该成员处于保护期内，无法创建案件")
		}
	}

	caseNo := s.generateCaseNo()
	result, err := s.db.Exec(`INSERT INTO cases (case_no, member_id, clause_id, template_id, status, incident_description, incident_time, punishment_detail, created_by) VALUES (?,?,?,?,?,?,?,?,?)`,
		caseNo, c.MemberID, c.ClauseID, c.TemplateID, "pending", c.IncidentDescription, c.IncidentTime, c.PunishmentDetail, c.CreatedBy)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return s.GetCase(id)
}

func (s *CaseService) UpdateCaseStatus(id int64, status string) error {
	validStatuses := map[string]bool{"pending": true, "in_progress": true, "completed": true, "cancelled": true, "appealed": true}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}
	_, err := s.db.Exec(`UPDATE cases SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, status, id)
	return err
}

func (s *CaseService) GetCasesByMember(memberId int64) ([]*models.Case, error) {
	rows, err := s.db.Query(caseQuery+` WHERE c.member_id=? ORDER BY c.created_at DESC`, memberId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanCases(rows)
}

func (s *CaseService) GetRecentCases(limit int) ([]*models.Case, error) {
	rows, err := s.db.Query(caseQuery+` ORDER BY c.created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanCases(rows)
}

func (s *CaseService) GetCaseStats() (*models.CaseStats, error) {
	stats := &models.CaseStats{
		StatusCounts: make(map[string]int64),
	}

	s.db.QueryRow(`SELECT COUNT(*) FROM members`).Scan(&stats.TotalMembers)
	s.db.QueryRow(`SELECT COUNT(*) FROM clauses WHERE is_active=1`).Scan(&stats.ActiveClauses)
	s.db.QueryRow(`SELECT COUNT(*) FROM cases WHERE strftime('%Y-%m', created_at)=strftime('%Y-%m','now')`).Scan(&stats.MonthCases)
	s.db.QueryRow(`SELECT COUNT(*) FROM cases WHERE status='pending'`).Scan(&stats.PendingCases)

	rows, err := s.db.Query(`SELECT status, COUNT(*) FROM cases GROUP BY status`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var status string
			var count int64
			rows.Scan(&status, &count)
			stats.StatusCounts[status] = count
		}
	}

	memberRows, err := s.db.Query(`SELECT m.id, m.name, COUNT(c.id) as cnt FROM members m LEFT JOIN cases c ON m.id=c.member_id GROUP BY m.id, m.name ORDER BY cnt DESC`)
	if err == nil {
		defer memberRows.Close()
		for memberRows.Next() {
			mc := models.MemberCaseCount{}
			memberRows.Scan(&mc.MemberID, &mc.MemberName, &mc.Count)
			stats.MemberCounts = append(stats.MemberCounts, mc)
		}
	}

	return stats, nil
}

func (s *CaseService) StartPunishment(caseId int64) error {
	_, err := s.db.Exec(`UPDATE cases SET status='in_progress', started_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=?`, caseId)
	return err
}

func (s *CaseService) CompletePunishment(caseId int64) error {
	_, err := s.db.Exec(`UPDATE cases SET status='completed', completed_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=?`, caseId)
	return err
}

func (s *CaseService) GetCaseComments(caseId int64) ([]*models.CaseComment, error) {
	rows, err := s.db.Query(`SELECT cc.id, cc.case_id, cc.member_id, cc.content, cc.created_at, COALESCE(m.name,'') FROM case_comments cc LEFT JOIN members m ON cc.member_id=m.id WHERE cc.case_id=? ORDER BY cc.created_at ASC`, caseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.CaseComment
	for rows.Next() {
		c := &models.CaseComment{}
		rows.Scan(&c.ID, &c.CaseID, &c.MemberID, &c.Content, &c.CreatedAt, &c.MemberName)
		comments = append(comments, c)
	}
	if comments == nil {
		comments = []*models.CaseComment{}
	}
	return comments, nil
}

func (s *CaseService) AddCaseComment(comment *models.CaseComment) (*models.CaseComment, error) {
	result, err := s.db.Exec(`INSERT INTO case_comments (case_id, member_id, content) VALUES (?,?,?)`,
		comment.CaseID, comment.MemberID, comment.Content)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	c := &models.CaseComment{}
	s.db.QueryRow(`SELECT cc.id, cc.case_id, cc.member_id, cc.content, cc.created_at, COALESCE(m.name,'') FROM case_comments cc LEFT JOIN members m ON cc.member_id=m.id WHERE cc.id=?`, id).
		Scan(&c.ID, &c.CaseID, &c.MemberID, &c.Content, &c.CreatedAt, &c.MemberName)
	return c, nil
}
