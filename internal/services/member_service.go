package services

import (
	"database/sql"
	"errors"
	"fmt"
	"fmps/internal/models"
	"time"
)

type MemberService struct {
	db *sql.DB
}

func NewMemberService(db *sql.DB) *MemberService {
	return &MemberService{db: db}
}

func (s *MemberService) GetAllMembers() ([]*models.Member, error) {
	rows, err := s.db.Query(`SELECT id, name, role, COALESCE(avatar,''), is_protected, protected_until, created_at, updated_at FROM members ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*models.Member
	for rows.Next() {
		m := &models.Member{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Role, &m.Avatar, &m.IsProtected, &m.ProtectedUntil, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if m.ProtectedUntil != nil {
			until, _ := time.Parse("2006-01-02 15:04:05", *m.ProtectedUntil)
			if time.Now().After(until) {
				m.IsProtected = false
			}
		}
		members = append(members, m)
	}
	if members == nil {
		members = []*models.Member{}
	}
	return members, nil
}

func (s *MemberService) GetMember(id int64) (*models.Member, error) {
	m := &models.Member{}
	err := s.db.QueryRow(`SELECT id, name, role, COALESCE(avatar,''), is_protected, protected_until, created_at, updated_at FROM members WHERE id=?`, id).
		Scan(&m.ID, &m.Name, &m.Role, &m.Avatar, &m.IsProtected, &m.ProtectedUntil, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MemberService) CreateMember(member *models.Member) (*models.Member, error) {
	result, err := s.db.Exec(`INSERT INTO members (name, role, avatar) VALUES (?,?,?)`,
		member.Name, member.Role, member.Avatar)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return s.GetMember(id)
}

func (s *MemberService) UpdateMember(member *models.Member) (*models.Member, error) {
	_, err := s.db.Exec(`UPDATE members SET name=?, role=?, avatar=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		member.Name, member.Role, member.Avatar, member.ID)
	if err != nil {
		return nil, err
	}
	return s.GetMember(member.ID)
}

func (s *MemberService) DeleteMember(id int64) error {
	if s.CheckProtection(id) {
		return errors.New("该成员处于保护期内，无法删除")
	}
	var count int64
	s.db.QueryRow(`SELECT COUNT(*) FROM cases WHERE member_id=? AND status IN ('pending','in_progress')`, id).Scan(&count)
	if count > 0 {
		return fmt.Errorf("该成员有 %d 个进行中的案件，无法删除", count)
	}
	_, err := s.db.Exec(`DELETE FROM members WHERE id=?`, id)
	return err
}

func (s *MemberService) SetProtection(memberId int64, durationHours int) error {
	until := time.Now().Add(time.Duration(durationHours) * time.Hour).Format("2006-01-02 15:04:05")
	_, err := s.db.Exec(`UPDATE members SET is_protected=1, protected_until=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, until, memberId)
	return err
}

func (s *MemberService) CheckProtection(memberId int64) bool {
	var isProtected bool
	var protectedUntil *string
	s.db.QueryRow(`SELECT is_protected, protected_until FROM members WHERE id=?`, memberId).
		Scan(&isProtected, &protectedUntil)
	if !isProtected || protectedUntil == nil {
		return false
	}
	until, err := time.Parse("2006-01-02 15:04:05", *protectedUntil)
	if err != nil {
		return false
	}
	return time.Now().Before(until)
}
