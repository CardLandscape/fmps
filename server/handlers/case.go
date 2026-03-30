package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CaseHandler struct {
	DB *gorm.DB
}

// PunishmentStep 解析后的惩罚步骤（旧格式，保留兼容）
type PunishmentStep struct {
	StartTime         string   `json:"start_time"`
	Duration          int      `json:"duration"`
	PunishmentDetails string   `json:"punishment_details"`
	Requirements      []string `json:"requirements"`
	DeductScoreRule   string   `json:"deduct_score_rule"`
	DeductScore       int      `json:"deduct_score"`
}

// ParsePunishmentProcess 解析惩罚过程文本，每行一条步骤（旧格式兼容）
// 格式：starttime|duration|punishmentdetails|req1|req2|req3|req4|req5|deductscorerule|deductscore
func ParsePunishmentProcess(text string) []PunishmentStep {
	var steps []PunishmentStep
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}
		step := PunishmentStep{}
		step.StartTime = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			step.Duration, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
		}
		if len(parts) > 2 {
			step.PunishmentDetails = strings.TrimSpace(parts[2])
		}
		// requirements: parts[3..7]
		for i := 3; i <= 7 && i < len(parts); i++ {
			req := strings.TrimSpace(parts[i])
			if req != "" {
				step.Requirements = append(step.Requirements, req)
			}
		}
		if len(parts) > 8 {
			step.DeductScoreRule = strings.TrimSpace(parts[8])
		}
		if len(parts) > 9 {
			step.DeductScore, _ = strconv.Atoi(strings.TrimSpace(parts[9]))
		}
		steps = append(steps, step)
	}
	return steps
}

// timeLineRegex matches lines that start with a time like "21:00" or "7:00"
var timeLineRegex = regexp.MustCompile(`^\d{1,2}:\d{2}`)

// levelSectionRegex matches exact level section headers like "A级", "B级", etc.
var levelSectionRegex = regexp.MustCompile(`^([A-D])级$`)

// ParseTxtByLevel parses a TXT file and extracts (prepItems, steps) for the given level (A/B/C/D).
// Preparation items are taken from a "惩罚工具" section (if present) and from non-timestamped
// lines before the first timestamped step within the level's section.
// Execution steps are the timestamped lines within the level's "惩罚流程" section.
func ParseTxtByLevel(content, level string) (prepItems []string, steps []string) {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	// --- Phase 1: extract global preparation items from "惩罚工具" section ---
	inToolSection := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Detect section headers
		if strings.Contains(trimmed, "惩罚工具") {
			inToolSection = true
			continue
		}
		// End tool section on other major headers
		if inToolSection && (strings.Contains(trimmed, "惩罚成绩") ||
			strings.Contains(trimmed, "惩罚姿势") ||
			strings.Contains(trimmed, "惩罚流程") ||
			levelSectionRegex.MatchString(trimmed)) {
			inToolSection = false
			continue
		}
		if inToolSection {
			// Split by common Chinese list separators
			subItems := splitPrepItems(trimmed)
			prepItems = append(prepItems, subItems...)
		}
	}

	// --- Phase 2: find the level's "惩罚流程" section and extract timestamped steps ---
	// We look for the level header inside the "惩罚流程" block
	inFlowSection := false
	inLevelSection := false
	targetHeader := level + "级"

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Detect "惩罚流程" section
		if strings.Contains(trimmed, "惩罚流程") && !inFlowSection {
			inFlowSection = true
			continue
		}
		if !inFlowSection {
			continue
		}
		// Detect target level header within flow section
		if trimmed == targetHeader || strings.HasPrefix(trimmed, targetHeader) {
			inLevelSection = true
			continue
		}
		// End of this level's section: next level header or a non-content separator
		if inLevelSection {
			m := levelSectionRegex.FindStringSubmatch(trimmed)
			if m != nil && m[1] != level {
				break
			}
			if timeLineRegex.MatchString(trimmed) {
				steps = append(steps, trimmed)
			}
		}
	}

	// De-duplicate and clean prep items
	seen := map[string]bool{}
	var cleanPrep []string
	for _, item := range prepItems {
		item = strings.TrimSpace(item)
		if item != "" && !seen[item] {
			seen[item] = true
			cleanPrep = append(cleanPrep, item)
		}
	}
	prepItems = cleanPrep

	return prepItems, steps
}

// splitPrepItems splits a preparation line by common separators (、，,，句末。) into individual items
func splitPrepItems(line string) []string {
	// Remove trailing 。or period
	line = strings.TrimRight(line, "。.")
	// Split by 、and ，
	re := regexp.MustCompile(`[、，,]+`)
	parts := re.Split(line, -1)
	var items []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		// Remove trailing 等
		p = strings.TrimRight(p, "等")
		p = strings.TrimSpace(p)
		if p != "" {
			items = append(items, p)
		}
	}
	return items
}

// computeGrade maps total deducted points to a grade string.
// totalDeducted should be the absolute value of total negative score_delta.
func computeGrade(totalDeducted int) string {
	if totalDeducted == 0 {
		return "满分"
	} else if totalDeducted <= 5 {
		return "优"
	} else if totalDeducted <= 15 {
		return "良"
	} else if totalDeducted <= 19 {
		return "达标"
	} else if totalDeducted <= 39 {
		return "不达标"
	}
	return "态度不端正"
}

func (h *CaseHandler) List(c *gin.Context) {
	var cases []models.Case
	if err := h.DB.Preload("Member").Preload("ParentMember").Preload("ChildMember").Find(&cases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, cases)
}

func (h *CaseHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var cas models.Case
	if err := h.DB.Preload("Member").Preload("ParentMember").Preload("ChildMember").First(&cas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "案件不存在"})
		return
	}

	var penalties []models.PenaltyPoint
	h.DB.Where("case_id = ?", id).Order("created_at desc").Find(&penalties)

	// Legacy steps (pipe format)
	steps := ParsePunishmentProcess(cas.PunishmentProcess)

	// New parsed steps
	var parsedSteps []string
	if cas.ParsedSteps != "" {
		_ = json.Unmarshal([]byte(cas.ParsedSteps), &parsedSteps)
	}

	// New prep items
	var prepItems []string
	if cas.PrepItems != "" {
		_ = json.Unmarshal([]byte(cas.PrepItems), &prepItems)
	}

	// Compute running total
	totalDeducted := 0
	for _, p := range penalties {
		if !p.Revoked && p.ScoreDelta < 0 {
			totalDeducted += -p.ScoreDelta
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"case":             cas,
		"punishment_steps": steps,
		"parsed_steps":     parsedSteps,
		"prep_items":       prepItems,
		"penalty_points":   penalties,
		"total_deducted":   totalDeducted,
		"current_grade":    computeGrade(totalDeducted),
	})
}

// validateCaseMembers checks that parent and child are set, exist, have the correct roles,
// and do not represent the same natural person (matched by name_cn + birth_date).
func (h *CaseHandler) validateCaseMembers(parentID, childID uint) error {
	if parentID == 0 {
		return fmt.Errorf("请选择家长成员")
	}
	if childID == 0 {
		return fmt.Errorf("请选择小孩成员")
	}
	if parentID == childID {
		return fmt.Errorf("案件中的家长和小孩不能是同一条记录")
	}

	var parent models.Member
	if err := h.DB.First(&parent, parentID).Error; err != nil {
		return fmt.Errorf("家长成员不存在")
	}
	if parent.Role != "parent" && parent.Role != "adult" {
		return fmt.Errorf("所选家长成员的类型不是「家长」")
	}

	var child models.Member
	if err := h.DB.First(&child, childID).Error; err != nil {
		return fmt.Errorf("小孩成员不存在")
	}
	if child.Role != "child" {
		return fmt.Errorf("所选小孩成员的类型不是「小孩」")
	}

	// Check they are not the same natural person (same name_cn + birth_date)
	pName := strings.TrimSpace(parent.NameCn)
	cName := strings.TrimSpace(child.NameCn)
	pBirth := strings.TrimSpace(parent.BirthDate)
	cBirth := strings.TrimSpace(child.BirthDate)
	if pName != "" && cName != "" && pBirth != "" && cBirth != "" {
		if pName == cName && pBirth == cBirth {
			return fmt.Errorf("案件中的家长和小孩不能是同一人")
		}
	}

	return nil
}

// ParseTxtRequest is the request body for the parse-txt endpoint
type ParseTxtRequest struct {
	Content     string `json:"content"`
	Level       string `json:"level"`
	TxtFilename string `json:"txt_filename"`
}

// ParseTxt parses a TXT file content and returns prep items + steps for the given level.
// POST /cases/parse-txt
func (h *CaseHandler) ParseTxt(c *gin.Context) {
	var req ParseTxtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}
	level := strings.ToUpper(strings.TrimSpace(req.Level))
	if level != "A" && level != "B" && level != "C" && level != "D" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "惩罚级别只能为 A/B/C/D"})
		return
	}
	prepItems, steps := ParseTxtByLevel(req.Content, level)
	c.JSON(http.StatusOK, gin.H{
		"prep_items": prepItems,
		"steps":      steps,
	})
}

func (h *CaseHandler) Create(c *gin.Context) {
	var cas models.Case
	if err := c.ShouldBindJSON(&cas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	// Validate parent and child members
	if err := h.validateCaseMembers(cas.ParentMemberID, cas.ChildMemberID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Set legacy MemberID to ChildMemberID for backward compat
	if cas.MemberID == 0 && cas.ChildMemberID > 0 {
		cas.MemberID = cas.ChildMemberID
	}

	// Validate punishment level if provided
	if cas.PunishmentLevel != "" {
		lvl := strings.ToUpper(cas.PunishmentLevel)
		if lvl != "A" && lvl != "B" && lvl != "C" && lvl != "D" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "惩罚级别只能为 A/B/C/D"})
			return
		}
		cas.PunishmentLevel = lvl
	}

	cas.Status = "pending"
	cas.CurrentStepIndex = -1
	if err := h.DB.Create(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
		return
	}
	h.DB.Preload("Member").Preload("ParentMember").Preload("ChildMember").First(&cas, cas.ID)
	c.JSON(http.StatusCreated, cas)
}

func (h *CaseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var cas models.Case
	if err := h.DB.First(&cas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "案件不存在"})
		return
	}

	if err := c.ShouldBindJSON(&cas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}
	cas.ID = uint(id)

	// Validate parent and child members if both are set
	if cas.ParentMemberID > 0 && cas.ChildMemberID > 0 {
		if err := h.validateCaseMembers(cas.ParentMemberID, cas.ChildMemberID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		// Keep legacy MemberID in sync
		if cas.MemberID == 0 {
			cas.MemberID = cas.ChildMemberID
		}
	}

	// Validate punishment level if provided
	if cas.PunishmentLevel != "" {
		lvl := strings.ToUpper(cas.PunishmentLevel)
		if lvl != "A" && lvl != "B" && lvl != "C" && lvl != "D" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "惩罚级别只能为 A/B/C/D"})
			return
		}
		cas.PunishmentLevel = lvl
	}

	if err := h.DB.Save(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	h.DB.Preload("Member").Preload("ParentMember").Preload("ChildMember").First(&cas, cas.ID)
	c.JSON(http.StatusOK, cas)
}

func (h *CaseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	if err := h.DB.Delete(&models.Case{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// StartPunishment 开始惩罚，将案件状态变为 active 并记录开始时间
// This transitions from pending (prep phase) to active (execution phase).
func (h *CaseHandler) StartPunishment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var cas models.Case
	if err := h.DB.First(&cas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "案件不存在"})
		return
	}

	now := time.Now()
	cas.Status = "active"
	cas.StartTime = &now
	// Start from step 0 if we have parsed steps, otherwise -1 (legacy mode)
	var parsedSteps []string
	if cas.ParsedSteps != "" {
		_ = json.Unmarshal([]byte(cas.ParsedSteps), &parsedSteps)
	}
	if len(parsedSteps) > 0 {
		cas.CurrentStepIndex = 0
	} else {
		cas.CurrentStepIndex = -1
	}

	if err := h.DB.Save(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, cas)
}

// CompleteStep marks the current execution step as done and advances to the next.
// If all steps are done, it finalizes the case with a computed grade.
// POST /cases/:id/complete-step
func (h *CaseHandler) CompleteStep(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var cas models.Case
	if err := h.DB.First(&cas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "案件不存在"})
		return
	}

	if cas.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "案件未处于执行状态"})
		return
	}

	var parsedSteps []string
	if cas.ParsedSteps != "" {
		_ = json.Unmarshal([]byte(cas.ParsedSteps), &parsedSteps)
	}

	nextIndex := cas.CurrentStepIndex + 1
	if len(parsedSteps) > 0 && nextIndex < len(parsedSteps) {
		// Advance to next step
		cas.CurrentStepIndex = nextIndex
		if err := h.DB.Save(&cas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"case": cas, "finished": false, "current_step_index": cas.CurrentStepIndex})
		return
	}

	// All steps done – complete the case
	cas.Status = "completed"
	grade := h.computeCaseGrade(cas.ID)
	cas.FinalGrade = grade
	if err := h.DB.Save(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"case": cas, "finished": true, "final_grade": grade})
}

// CompletePunishment 完成惩罚（直接结束，计算成绩）
func (h *CaseHandler) CompletePunishment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var cas models.Case
	if err := h.DB.First(&cas, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "案件不存在"})
		return
	}

	grade := h.computeCaseGrade(cas.ID)
	cas.Status = "completed"
	cas.FinalGrade = grade
	if err := h.DB.Save(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"case": cas, "final_grade": grade})
}

// computeCaseGrade queries all non-revoked penalties for a case and returns the grade string.
func (h *CaseHandler) computeCaseGrade(caseID uint) string {
	var penalties []models.PenaltyPoint
	h.DB.Where("case_id = ? AND revoked = ?", caseID, false).Find(&penalties)
	totalDeducted := 0
	for _, p := range penalties {
		if p.ScoreDelta < 0 {
			totalDeducted += -p.ScoreDelta
		}
	}
	return computeGrade(totalDeducted)
}

