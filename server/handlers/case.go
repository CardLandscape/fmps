package handlers

import (
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
	if err := h.DB.Preload("Member").Find(&cases).Error; err != nil {
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
	if err := h.DB.Preload("Member").First(&cas, id).Error; err != nil {
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

func (h *CaseHandler) Create(c *gin.Context) {
	var cas models.Case
	if err := c.ShouldBindJSON(&cas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}
	cas.Status = "pending"
	if err := h.DB.Create(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
		return
	}
	h.DB.Preload("Member").First(&cas, cas.ID)
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

	if err := h.DB.Save(&cas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	h.DB.Preload("Member").First(&cas, cas.ID)
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
