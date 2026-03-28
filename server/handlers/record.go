package handlers

import (
	"net/http"
	"strconv"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecordHandler struct {
	DB *gorm.DB
}

func (h *RecordHandler) List(c *gin.Context) {
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	query := h.DB.Preload("Member").Preload("Rule").Order("id DESC").Limit(limit)

	if memberID := c.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}

	var records []models.Record
	if err := query.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *RecordHandler) Create(c *gin.Context) {
	var record models.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	// If Points is 0, copy from Rule.Points
	if record.Points == 0 && record.RuleID != 0 {
		var rule models.Rule
		if err := h.DB.First(&rule, record.RuleID).Error; err == nil {
			record.Points = rule.Points
		}
	}

	if err := h.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
		return
	}

	// Reload with associations
	h.DB.Preload("Member").Preload("Rule").First(&record, record.ID)
	c.JSON(http.StatusCreated, record)
}

func (h *RecordHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	if err := h.DB.Delete(&models.Record{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
