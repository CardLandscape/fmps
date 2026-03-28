package handlers

import (
	"net/http"
	"time"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret string
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}

	var usernameSetting, passwordSetting models.Setting
	if err := h.DB.Where("key = ?", "admin_username").First(&usernameSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器错误"})
		return
	}
	if err := h.DB.Where("key = ?", "admin_password").First(&passwordSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器错误"})
		return
	}

	if req.Username != usernameSetting.Value || req.Password != passwordSetting.Value {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenStr,
		"message": "登录成功",
	})
}
