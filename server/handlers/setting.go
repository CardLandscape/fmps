package handlers

import (
	"net/http"
	"strings"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingHandler struct {
	DB *gorm.DB
}

func (h *SettingHandler) Get(c *gin.Context) {
	var settings []models.Setting
	if err := h.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}

	result := make(map[string]string)
	for _, s := range settings {
		if strings.Contains(strings.ToLower(s.Key), "password") {
			result[s.Key] = "***"
		} else {
			result[s.Key] = s.Value
		}
	}
	c.JSON(http.StatusOK, result)
}

func (h *SettingHandler) Update(c *gin.Context) {
	var updates map[string]string
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	for key, value := range updates {
		// Hash password before storing
		if key == "admin_password" {
			hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "密码处理失败"})
				return
			}
			value = string(hashed)
		}
		setting := models.Setting{Key: key, Value: value}
		h.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&setting)
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置已更新"})
}
