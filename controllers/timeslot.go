package controllers

import (
	"badminton-reservation-api/models"
	"badminton-reservation-api/utils"

	"github.com/beego/beego/v2/server/web"
)

type TimeslotController struct {
	web.Controller
}

// GetAvailableTimeslots returns timeslots available for a given court and booking_date
// GetAvailableTimeslots godoc
// @Summary Get available timeslots for a court and date
// @Description Returns timeslots available for a given court and booking_date
// @Tags timeslots
// @Accept json
// @Produce json
// @Param booking_date query string true "Booking date (YYYY-MM-DD)"
// @Param court_id query int true "Court ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/timeslots [get]
func (c *TimeslotController) GetAvailableTimeslots() {
	bookingDate := c.GetString("booking_date")
	courtId, _ := c.GetInt("court_id", 0)

	if bookingDate == "" || courtId == 0 {
		utils.SendBadRequest(&c.Controller, "booking_date and court_id are required", nil)
		return
	}

	slots, err := models.GetAvailableTimeslots(courtId, bookingDate)
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error retrieving timeslots", err.Error())
		return
	}

	utils.SendSuccess(&c.Controller, "Available timeslots retrieved successfully", slots)
}

// GetAllTimeslots returns all timeslots
// GetAllTimeslots godoc
// @Summary Get all timeslots
// @Description Returns all defined timeslots
// @Tags timeslots
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/timeslots/all [get]
func (c *TimeslotController) GetAllTimeslots() {
	slots, err := models.GetAllTimeslots()
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error retrieving timeslots", err.Error())
		return
	}
	utils.SendSuccess(&c.Controller, "Timeslots retrieved successfully", slots)
}
