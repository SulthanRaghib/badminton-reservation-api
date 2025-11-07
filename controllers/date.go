package controllers

import (
	"badminton-reservation-api/utils"
	"os"
	"strconv"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type DateController struct {
	web.Controller
}

type DateResponse struct {
	Date      string `json:"date"`
	DayName   string `json:"day_name"`
	IsWeekend bool   `json:"is_weekend"`
}

// GetAvailableDates godoc
// @Summary Get available dates for booking
// @Description Get list of available dates for the next N days
// @Tags dates
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/dates [get]
func (c *DateController) GetAvailableDates() {
	// Get max booking days from env, default to 30
	maxDaysStr := os.Getenv("MAX_BOOKING_DAYS_AHEAD")
	maxDays := 30
	if maxDaysStr != "" {
		if days, err := strconv.Atoi(maxDaysStr); err == nil {
			maxDays = days
		}
	}

	var dates []DateResponse
	today := time.Now()

	// Generate dates for the next N days
	for i := 0; i <= maxDays; i++ {
		date := today.AddDate(0, 0, i)

		// Format date
		dateStr := date.Format("2006-01-02")
		dayName := date.Format("Monday")
		isWeekend := date.Weekday() == time.Saturday || date.Weekday() == time.Sunday

		dates = append(dates, DateResponse{
			Date:      dateStr,
			DayName:   dayName,
			IsWeekend: isWeekend,
		})
	}

	utils.SendSuccess(&c.Controller, "Available dates retrieved successfully", dates)
}
