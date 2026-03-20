package seeds

import (
	"database/sql"
	"fmps/internal/services"
	"log"
)

func SeedDefaults(db *sql.DB, settingService *services.SettingService) {
	if settingService.IsInitialized() {
		return
	}

	log.Println("Initializing default data...")

	clauses := []struct {
		code, title, desc, category string
		severity                    int
	}{
		{"C001", "不按时完成作业", "未能在规定时间内完成学校布置的作业", "学习", 2},
		{"C002", "说脏话", "使用不文明、粗鲁的语言", "行为", 3},
		{"C003", "不整理房间", "房间脏乱，未能保持整洁", "家务", 1},
		{"C004", "顶嘴/态度不好", "对家长态度恶劣，顶嘴不听从管教", "态度", 3},
		{"C005", "沉迷电子设备", "过度使用手机、游戏等电子设备，超出规定时间", "行为", 4},
		{"C006", "不按时睡觉", "不遵守作息时间，超时不睡觉", "行为", 2},
		{"C007", "浪费食物", "故意浪费食物，不珍惜粮食", "行为", 2},
		{"C008", "打架/欺负他人", "与他人发生肢体冲突或言语欺凌", "行为", 5},
		{"C009", "不做家务", "被分配的家务任务未完成", "家务", 2},
		{"C010", "考试不及格", "学科考试成绩不及格", "学习", 4},
	}

	for _, c := range clauses {
		db.Exec(`INSERT OR IGNORE INTO clauses (code, title, description, severity, category) VALUES (?,?,?,?,?)`,
			c.code, c.title, c.desc, c.severity, c.category)
	}

	templates := []struct {
		name, ptype, desc string
		duration          int
	}{
		{"口头警告", "口头警告", "口头提醒并警告，适用于轻微违规行为", 0},
		{"书面检讨", "书面检讨", "需撰写不少于200字的检讨书，认真反思自己的错误", 30},
		{"限制娱乐", "限制娱乐", "禁止使用电子设备（手机/平板/电脑/电视）", 120},
	}

	for _, t := range templates {
		db.Exec(`INSERT OR IGNORE INTO punishment_templates (name, punishment_type, description, duration_minutes) VALUES (?,?,?,?)`,
			t.name, t.ptype, t.desc, t.duration)
	}

	settingService.MarkInitialized()
	log.Println("Default data initialized successfully")
}
