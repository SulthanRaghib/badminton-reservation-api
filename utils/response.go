package utils

import (
	"github.com/beego/beego/v2/server/web"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SendSuccess sends a successful JSON response
func SendSuccess(c *web.Controller, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.Data["json"] = response
	c.ServeJSON()
}

// SendError sends an error JSON response
func SendError(c *web.Controller, statusCode int, message string, err interface{}) {
	response := Response{
		Success: false,
		Message: message,
		Error:   err,
	}
	c.Ctx.Output.SetStatus(statusCode)
	c.Data["json"] = response
	c.ServeJSON()
}

// SendBadRequest sends a 400 Bad Request response
func SendBadRequest(c *web.Controller, message string, err interface{}) {
	SendError(c, 400, message, err)
}

// SendNotFound sends a 404 Not Found response
func SendNotFound(c *web.Controller, message string) {
	SendError(c, 404, message, nil)
}

// SendInternalError sends a 500 Internal Server Error response
func SendInternalError(c *web.Controller, message string, err interface{}) {
	SendError(c, 500, message, err)
}

// SendUnauthorized sends a 401 Unauthorized response
func SendUnauthorized(c *web.Controller, message string) {
	SendError(c, 401, message, nil)
}

// SendConflict sends a 409 Conflict response
func SendConflict(c *web.Controller, message string, err interface{}) {
	SendError(c, 409, message, err)
}
