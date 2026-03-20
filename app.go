package main

import (
	"context"
	"fmps/internal/db"
	"fmps/internal/models"
	"fmps/internal/seeds"
	"fmps/internal/services"
)

type App struct {
	ctx             context.Context
	memberService   *services.MemberService
	clauseService   *services.ClauseService
	templateService *services.TemplateService
	caseService     *services.CaseService
	appealService   *services.AppealService
	settingService  *services.SettingService
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	database := db.InitDB()
	a.memberService = services.NewMemberService(database)
	a.clauseService = services.NewClauseService(database)
	a.templateService = services.NewTemplateService(database)
	a.caseService = services.NewCaseService(database)
	a.appealService = services.NewAppealService(database)
	a.settingService = services.NewSettingService(database)
	seeds.SeedDefaults(database, a.settingService)
}

// Member methods
func (a *App) GetAllMembers() ([]*models.Member, error) { return a.memberService.GetAllMembers() }
func (a *App) GetMember(id int64) (*models.Member, error) { return a.memberService.GetMember(id) }
func (a *App) CreateMember(member *models.Member) (*models.Member, error) {
	return a.memberService.CreateMember(member)
}
func (a *App) UpdateMember(member *models.Member) (*models.Member, error) {
	return a.memberService.UpdateMember(member)
}
func (a *App) DeleteMember(id int64) error { return a.memberService.DeleteMember(id) }
func (a *App) SetProtection(memberId int64, durationHours int) error {
	return a.memberService.SetProtection(memberId, durationHours)
}
func (a *App) CheckProtection(memberId int64) bool { return a.memberService.CheckProtection(memberId) }

// Clause methods
func (a *App) GetAllClauses() ([]*models.Clause, error) { return a.clauseService.GetAllClauses() }
func (a *App) GetClause(id int64) (*models.Clause, error) { return a.clauseService.GetClause(id) }
func (a *App) CreateClause(clause *models.Clause) (*models.Clause, error) {
	return a.clauseService.CreateClause(clause)
}
func (a *App) UpdateClause(clause *models.Clause) (*models.Clause, error) {
	return a.clauseService.UpdateClause(clause)
}
func (a *App) DeleteClause(id int64) error { return a.clauseService.DeleteClause(id) }
func (a *App) GetClausesByCategory(category string) ([]*models.Clause, error) {
	return a.clauseService.GetClausesByCategory(category)
}

// Template methods
func (a *App) GetAllTemplates() ([]*models.PunishmentTemplate, error) {
	return a.templateService.GetAllTemplates()
}
func (a *App) GetTemplate(id int64) (*models.PunishmentTemplate, error) {
	return a.templateService.GetTemplate(id)
}
func (a *App) CreateTemplate(template *models.PunishmentTemplate) (*models.PunishmentTemplate, error) {
	return a.templateService.CreateTemplate(template)
}
func (a *App) UpdateTemplate(template *models.PunishmentTemplate) (*models.PunishmentTemplate, error) {
	return a.templateService.UpdateTemplate(template)
}
func (a *App) DeleteTemplate(id int64) error { return a.templateService.DeleteTemplate(id) }

// Case methods
func (a *App) GetAllCases() ([]*models.Case, error) { return a.caseService.GetAllCases() }
func (a *App) GetCase(id int64) (*models.Case, error) { return a.caseService.GetCase(id) }
func (a *App) CreateCase(c *models.Case) (*models.Case, error) { return a.caseService.CreateCase(c) }
func (a *App) UpdateCaseStatus(id int64, status string) error {
	return a.caseService.UpdateCaseStatus(id, status)
}
func (a *App) GetCasesByMember(memberId int64) ([]*models.Case, error) {
	return a.caseService.GetCasesByMember(memberId)
}
func (a *App) GetCaseStats() (*models.CaseStats, error) { return a.caseService.GetCaseStats() }
func (a *App) StartPunishment(caseId int64) error      { return a.caseService.StartPunishment(caseId) }
func (a *App) CompletePunishment(caseId int64) error   { return a.caseService.CompletePunishment(caseId) }
func (a *App) GetCaseComments(caseId int64) ([]*models.CaseComment, error) {
	return a.caseService.GetCaseComments(caseId)
}
func (a *App) AddCaseComment(comment *models.CaseComment) (*models.CaseComment, error) {
	return a.caseService.AddCaseComment(comment)
}
func (a *App) GetRecentCases(limit int) ([]*models.Case, error) {
	return a.caseService.GetRecentCases(limit)
}

// Appeal methods
func (a *App) GetAppealsByCase(caseId int64) ([]*models.Appeal, error) {
	return a.appealService.GetAppealsByCase(caseId)
}
func (a *App) CreateAppeal(appeal *models.Appeal) (*models.Appeal, error) {
	return a.appealService.CreateAppeal(appeal)
}
func (a *App) ReviewAppeal(appealId int64, reviewerId int64, approved bool, comment string) error {
	return a.appealService.ReviewAppeal(appealId, reviewerId, approved, comment)
}

// Setting methods
func (a *App) GetSetting(key string) (string, error) { return a.settingService.GetSetting(key) }
func (a *App) SetSetting(key, value string) error    { return a.settingService.SetSetting(key, value) }
func (a *App) GetAllSettings() (map[string]string, error) {
	return a.settingService.GetAllSettings()
}
