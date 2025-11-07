package controllers

import (
	"badminton-reservation-api/models"
	"badminton-reservation-api/utils"

	"github.com/beego/beego/v2/server/web"
)

type CourtController struct {
	web.Controller
}

// GetAvailableCourts returns courts available for a given booking_date and timeslot_id
// GetAvailableCourts godoc
// @Summary Get available courts for a date and timeslot
// @Description Returns courts available for a given booking_date and timeslot_id
// @Tags courts
// @Accept json
// @Produce json
// @Param booking_date query string true "Booking date (YYYY-MM-DD)"
// @Param timeslot_id query int true "Timeslot ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/courts [get]
func (c *CourtController) GetAvailableCourts() {
	bookingDate := c.GetString("booking_date")
	timeslotId, _ := c.GetInt("timeslot_id", 0)

	if bookingDate == "" || timeslotId == 0 {
		utils.SendBadRequest(&c.Controller, "booking_date and timeslot_id are required", nil)
		return
	}

	courts, err := models.GetAvailableCourts(bookingDate, timeslotId)
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error retrieving courts", err.Error())
		return
	}

	utils.SendSuccess(&c.Controller, "Available courts retrieved successfully", courts)
}

// GetAllCourts returns all active courts
// GetAllCourts godoc
// @Summary Get all active courts
// @Description Returns all active courts
// @Tags courts
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/courts/all [get]
func (c *CourtController) GetAllCourts() {
	courts, err := models.GetAllActiveCourts()
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error retrieving courts", err.Error())
		return
	}
	utils.SendSuccess(&c.Controller, "Courts retrieved successfully", courts)
}
