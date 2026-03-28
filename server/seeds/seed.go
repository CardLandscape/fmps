package seeds

import (
	"fmps/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Run(db *gorm.DB) {
	// Upsert default credentials (only insert if not exists)
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Setting{Key: "admin_username", Value: "admin"})
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Setting{Key: "admin_password", Value: "123456"})

	// Seed sample rules if table is empty
	var count int64
	db.Model(&models.Rule{}).Count(&count)
	if count == 0 {
		rules := []models.Rule{
			{Name: "不按时完成作业", Description: "没有按时完成学校作业", Points: 5, Category: "学习"},
			{Name: "顶嘴", Description: "对长辈顶嘴不礼貌", Points: 3, Category: "礼仪"},
			{Name: "打架", Description: "与他人发生肢体冲突", Points: 10, Category: "行为"},
			{Name: "说谎", Description: "故意欺骗他人", Points: 4, Category: "诚信"},
			{Name: "不整理房间", Description: "不保持房间整洁", Points: 2, Category: "习惯"},
		}
		db.Create(&rules)
	}
}
