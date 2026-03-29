package handlers

import (
	"fmt"
	"net/http"
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

// PunishmentStep 解析后的惩罚步骤
type PunishmentStep struct {
	StartTime         string   `json:"start_time"`
	Duration          int      `json:"duration"`
	PunishmentDetails string   `json:"punishment_details"`
	Requirements      []string `json:"requirements"`
	DeductScoreRule   string   `json:"deduct_score_rule"`
	DeductScore       int      `json:"deduct_score"`
}

// ParsePunishmentProcess 解析惩罚过程文本，每行一条步骤
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

	steps := ParsePunishmentProcess(cas.PunishmentProcess)

	c.JSON(http.StatusOK, gin.H{
		"case":             cas,
		"punishment_steps": steps,
		"penalty_points":   penalties,
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

	cas.Status = "pending"
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

	if err := h.DB.Save(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, cas)
}

// CompletePunishment 完成惩罚
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

	cas.Status = "completed"
	if err := h.DB.Save(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, cas)
}
