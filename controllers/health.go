package controllers

import (
	"badminton-reservation-api/utils"

	"github.com/beego/beego/v2/server/web"
)

type HealthController struct {
	web.Controller
}

// Get returns a simple health response
// Get godoc
// @Summary Health check
// @Description Returns service health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /health [get]
func (h *HealthController) Get() {
	utils.SendSuccess(&h.Controller, "OK", map[string]string{"status": "ok"})
}
