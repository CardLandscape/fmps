package handlers

import (
	"net/http"
	"strconv"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MemberHandler struct {
	DB *gorm.DB
}

func (h *MemberHandler) List(c *gin.Context) {
	var members []models.Member
	if err := h.DB.Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) Create(c *gin.Context) {
	var member models.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}
	if err := h.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (h *MemberHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	var member models.Member
	if err := h.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "成员不存在"})
		return
	}

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求格式错误"})
		return
	}
	member.ID = uint(id)

	if err := h.DB.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (h *MemberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的ID"})
		return
	}

	if err := h.DB.Delete(&models.Member{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
