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
	// When credentials are allowed, do NOT set Access-Control-Allow-Origin to '*'.
	// Instead echo back the request Origin if allowed.
	if allowedOrigins == "*" {
		if origin != "" {
			ctx.Output.Header("Access-Control-Allow-Origin", origin)
		} else {
			ctx.Output.Header("Access-Control-Allow-Origin", "*")
		}
	} else {
		origins := strings.Split(allowedOrigins, ",")
		for _, allowedOrigin := range origins {
			if strings.TrimSpace(allowedOrigin) == origin {
				ctx.Output.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}
	}

	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	// Allow common request headers and the Access-Control request headers
	ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
	ctx.Output.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Max-Age", "86400")

	// Handle preflight requests: respond with 200 OK and the CORS headers above.
	if strings.ToUpper(ctx.Input.Method()) == "OPTIONS" {
		// Some browsers expect 200; return 200 with empty body.
		ctx.Output.SetStatus(200)
		ctx.Output.Body([]byte(""))
		return
	}
}
