package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")

	var allowedOrigins []string
	allowCredentials := false

	if originsEnv != "" {
		allowedOrigins = strings.Split(originsEnv, ",")
		allowCredentials = true
	} else {
		allowedOrigins = []string{"*"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	})
}
