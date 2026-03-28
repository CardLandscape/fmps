package main

import (
	"fmt"
	"log"

	"fmps/db"
	"fmps/seeds"
)

func main() {
	cfg := LoadConfig()

	database, err := db.Init(cfg.DBPath)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	seeds.Run(database)

	r := SetupRouter(database, cfg)

	fmt.Printf("FMPS 服务启动，访问 http://localhost:%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
