package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fmps/handlers"
	"fmps/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg Config) *gin.Engine {
	r := gin.Default()

	// CORS - allow localhost origins for development; restrict as needed via FMPS_CORS_ORIGIN env
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:8080", "http://127.0.0.1:8080"}
	if envOrigin := os.Getenv("FMPS_CORS_ORIGIN"); envOrigin != "" {
		allowedOrigins = append(allowedOrigins, envOrigin)
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Serve static files: use FMPS_WEB_DIR if set, otherwise fall back to ../web/dist relative to executable
	var distDir string
	if envDir := os.Getenv("FMPS_WEB_DIR"); envDir != "" {
		distDir = envDir
	} else {
		execPath, err := os.Executable()
		if err == nil {
			distDir = filepath.Join(filepath.Dir(execPath), "..", "web", "dist")
		} else {
			distDir = filepath.Join("..", "web", "dist")
		}
	}

	if _, err := os.Stat(distDir); err == nil {
		r.Static("/assets", filepath.Join(distDir, "assets"))
		r.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))
		r.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
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
	caseHandler := &handlers.CaseHandler{DB: db}
	penaltyHandler := &handlers.PenaltyHandler{DB: db}

	// Public routes
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/api/login", authHandler.Login)

	// Protected routes
	api := r.Group("/api", middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/members", memberHandler.List)
		api.POST("/members", memberHandler.Create)
		api.PUT("/members/:id", memberHandler.Update)
		api.POST("/members/:id/delete", memberHandler.DeleteWithAuth)
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

		api.GET("/cases", caseHandler.List)
		api.POST("/cases", caseHandler.Create)
		api.POST("/cases/parse-txt", caseHandler.ParseTxt)
		api.POST("/cases/parse-txt-levels", caseHandler.ParseTxtLevelsHandler)
		api.GET("/cases/:id", caseHandler.Get)
		api.PUT("/cases/:id", caseHandler.Update)
		api.DELETE("/cases/:id", caseHandler.Delete)
		api.POST("/cases/:id/start", caseHandler.StartPunishment)
		api.POST("/cases/:id/complete", caseHandler.CompletePunishment)
		api.POST("/cases/:id/complete-step", caseHandler.CompleteStep)
		api.POST("/cases/:id/penalty", penaltyHandler.AddPenalty)

		api.POST("/penalty/:id/revoke", penaltyHandler.RevokePenalty)
	}

	return r
}
