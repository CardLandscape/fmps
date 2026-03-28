package main

import (
	"net/http"
	"os"
	"path/filepath"

	"fmps/handlers"
	"fmps/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg Config) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Serve static files from ../web/dist (relative to executable)
	execPath, err := os.Executable()
	var distDir string
	if err == nil {
		distDir = filepath.Join(filepath.Dir(execPath), "..", "web", "dist")
	} else {
		distDir = filepath.Join("..", "web", "dist")
	}

	if _, err := os.Stat(distDir); err == nil {
		r.Static("/assets", filepath.Join(distDir, "assets"))
		r.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))
		r.NoRoute(func(c *gin.Context) {
			if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusNotFound, gin.H{"message": "接口不存在"})
				return
			}
			c.File(filepath.Join(distDir, "index.html"))
		})
	}

	// Init handlers
	authHandler := &handlers.AuthHandler{DB: db, JWTSecret: cfg.JWTSecret}
	memberHandler := &handlers.MemberHandler{DB: db}
	ruleHandler := &handlers.RuleHandler{DB: db}
	recordHandler := &handlers.RecordHandler{DB: db}
	settingHandler := &handlers.SettingHandler{DB: db}
	statsHandler := &handlers.StatsHandler{DB: db}

	// Public routes
	r.POST("/api/login", authHandler.Login)

	// Protected routes
	api := r.Group("/api", middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/members", memberHandler.List)
		api.POST("/members", memberHandler.Create)
		api.PUT("/members/:id", memberHandler.Update)
		api.DELETE("/members/:id", memberHandler.Delete)

		api.GET("/rules", ruleHandler.List)
		api.POST("/rules", ruleHandler.Create)
		api.PUT("/rules/:id", ruleHandler.Update)
		api.DELETE("/rules/:id", ruleHandler.Delete)

		api.GET("/records", recordHandler.List)
		api.POST("/records", recordHandler.Create)
		api.DELETE("/records/:id", recordHandler.Delete)

		api.GET("/settings", settingHandler.Get)
		api.PUT("/settings", settingHandler.Update)

		api.GET("/stats", statsHandler.Get)
	}

	return r
}
