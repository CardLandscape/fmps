package handlers

import (
	"net/http"

	"fmps/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatsHandler struct {
	DB *gorm.DB
}

type MemberStat struct {
	MemberID    uint   `json:"member_id"`
	MemberName  string `json:"member_name"`
	RecordCount int64  `json:"record_count"`
	TotalPoints int64  `json:"total_points"`
}

func (h *StatsHandler) Get(c *gin.Context) {
	var totalMembers, totalRules, totalRecords int64
	var totalPoints int64

	h.DB.Model(&models.Member{}).Count(&totalMembers)
	h.DB.Model(&models.Rule{}).Count(&totalRules)
	h.DB.Model(&models.Record{}).Count(&totalRecords)
	h.DB.Model(&models.Record{}).Select("COALESCE(SUM(points), 0)").Scan(&totalPoints)

	// Single aggregated query to avoid N+1
	type memberAgg struct {
		MemberID    uint  `gorm:"column:member_id"`
		RecordCount int64 `gorm:"column:record_count"`
		TotalPoints int64 `gorm:"column:total_points"`
	}
	var aggs []memberAgg
	h.DB.Model(&models.Record{}).
		Select("member_id, COUNT(*) as record_count, COALESCE(SUM(points), 0) as total_points").
		Group("member_id").
		Scan(&aggs)

	aggMap := make(map[uint]memberAgg, len(aggs))
	for _, a := range aggs {
		aggMap[a.MemberID] = a
	}

	var members []models.Member
	h.DB.Find(&members)

	memberStats := make([]MemberStat, 0, len(members))
	for _, m := range members {
		a := aggMap[m.ID]
		memberStats = append(memberStats, MemberStat{
			MemberID:    m.ID,
			MemberName:  m.Name,
			RecordCount: a.RecordCount,
			TotalPoints: a.TotalPoints,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_members": totalMembers,
		"total_rules":   totalRules,
		"total_records": totalRecords,
		"total_points":  totalPoints,
		"member_stats":  memberStats,
	})
}
