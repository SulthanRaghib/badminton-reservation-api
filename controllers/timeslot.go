package controllers

import (
	"badminton-reservation-api/models"
	"badminton-reservation-api/utils"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type TimeslotController struct {
	web.Controller
}

// GetAvailableTimeslots returns timeslots for a given court and booking_date with availability flag
// GetAvailableTimeslots godoc
// @Summary Get timeslots for a court and date (includes availability flag)
// @Description Returns all globally active timeslots and an `available` boolean per timeslot for the specified court and booking_date. `available=false` means the slot is already booked/unavailable for that date and court.
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

	// Return all globally active timeslots along with an `available` flag per timeslot
	allSlots, err := models.GetAllTimeslots()
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error retrieving timeslots", err.Error())
		return
	}

	type SlotWithAvailability struct {
		Id        int    `json:"id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		IsActive  bool   `json:"is_active"`
		Available bool   `json:"available"`
	}

	var result []SlotWithAvailability
	for _, s := range allSlots {
		// default: only consider globally active timeslots
		if !s.IsActive {
			continue
		}
		// Check timeslot_availabilities: if an explicit unavailable row exists -> available=false
		o := orm.NewOrm()
		cnt, err := o.QueryTable("timeslot_availabilities").Filter("court_id", courtId).Filter("timeslot_id", s.Id).Filter("booking_date", bookingDate).Filter("is_active", false).Count()

		available := true
		if err == nil && cnt > 0 {
			available = false
		} else {
			// fallback to checking reservations table
			ok, err := models.CheckAvailability(courtId, s.Id, bookingDate)
			if err != nil {
				utils.SendInternalError(&c.Controller, "Error checking availability", err.Error())
				return
			}
			available = ok
		}

		result = append(result, SlotWithAvailability{
			Id:        s.Id,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
			IsActive:  s.IsActive,
			Available: available,
		})
	}

	utils.SendSuccess(&c.Controller, "Timeslots retrieved successfully", result)
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
