package handlers

import (
	"net/http"
	"strconv"
	"time"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenaltyHandler struct {
	DB *gorm.DB
}

type CreatePenaltyRequest struct {
	RuleText   string `json:"rule_text"`
	ScoreDelta int    `json:"score_delta"`
	Reason     string `json:"reason"`
}

type RevokePenaltyRequest struct {
	Password string `json:"password"`
	Reason   string `json:"reason"`
}

// AddPenalty 扣分
func (h *PenaltyHandler) AddPenalty(c *gin.Context) {
	caseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的案件ID"})
		return
	}

	var cas models.Case
	if err := h.DB.First(&cas, caseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "案件不存在"})
		return
	}

	var req CreatePenaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	penalty := models.PenaltyPoint{
		CaseID:     uint(caseID),
		MemberID:   cas.MemberID,
		RuleText:   req.RuleText,
		ScoreDelta: req.ScoreDelta,
		Reason:     req.Reason,
	}

	if err := h.DB.Create(&penalty).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "扣分失败"})
		return
	}

	c.JSON(http.StatusCreated, penalty)
}

// RevokePenalty 撤回扣分
func (h *PenaltyHandler) RevokePenalty(c *gin.Context) {
	penaltyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的扣分ID"})
		return
	}

	var req RevokePenaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	// 查询授权密码
	var setting models.Setting
	if err := h.DB.Where("key = ?", "authorization_password").First(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统配置错误"})
		return
	}

	// 查询扣分记录
	var penalty models.PenaltyPoint
	if err := h.DB.First(&penalty, penaltyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "扣分记录不存在"})
		return
	}

	if penalty.Revoked {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该扣分记录已撤回"})
		return
	}

	// 验证密码
	if req.Password != setting.Value {
		// 密码错误：额外扣 100 分，原因为"作弊"
		cheatPenalty := models.PenaltyPoint{
			CaseID:     penalty.CaseID,
			MemberID:   penalty.MemberID,
			RuleText:   "授权密码错误",
			ScoreDelta: -100,
			Reason:     "作弊",
		}
		h.DB.Create(&cheatPenalty)
		c.JSON(http.StatusForbidden, gin.H{"message": "授权密码错误，已额外扣除100分"})
		return
	}

	// 密码正确：撤回
	now := time.Now()
	penalty.Revoked = true
	penalty.RevokedAt = &now
	penalty.RevokeReason = req.Reason
	if err := h.DB.Save(&penalty).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "撤回失败"})
		return
	}

	c.JSON(http.StatusOK, penalty)
}
