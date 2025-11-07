package middleware

import (
	"os"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS(ctx *context.Context) {
	// Get allowed origins from environment variable
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	origin := ctx.Input.Header("Origin")

	// Check if origin is allowed
	if allowedOrigins == "*" {
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
	} else {
		origins := strings.Split(allowedOrigins, ",")
		for _, allowedOrigin := range origins {
			if strings.TrimSpace(allowedOrigin) == origin {
				ctx.Output.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}
	}

	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Max-Age", "86400")

	// Handle preflight requests
	if ctx.Input.Method() == "OPTIONS" {
		ctx.Output.SetStatus(204)
		return
	}
}
