package controllers

import (
	"badminton-reservation-api/models"
	"badminton-reservation-api/utils"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
)

type ReservationController struct {
	web.Controller
}

type CreateReservationRequest struct {
	CourtId       int    `json:"court_id"`
	TimeslotId    int    `json:"timeslot_id"`
	BookingDate   string `json:"booking_date"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	Notes         string `json:"notes"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// CreateReservation godoc
// @Summary Create a new reservation
// @Description Create a new court reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body CreateReservationRequest true "Reservation details"
// @Success 201 {object} utils.Response
// @Router /api/v1/reservations [post]
func (c *ReservationController) CreateReservation() {
	var req CreateReservationRequest

	// Parse request body
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendBadRequest(&c.Controller, "Invalid request body", err.Error())
		return
	}

	// Validate required fields
	missingFields := utils.ValidateRequired(map[string]string{
		"customer_name":  req.CustomerName,
		"customer_email": req.CustomerEmail,
		"customer_phone": req.CustomerPhone,
		"booking_date":   req.BookingDate,
	})

	if len(missingFields) > 0 {
		utils.SendBadRequest(&c.Controller, "Missing required fields", missingFields)
		return
	}

	// Validate email format
	if !utils.ValidateEmail(req.CustomerEmail) {
		utils.SendBadRequest(&c.Controller, "Invalid email format", nil)
		return
	}

	// Validate phone format
	if !utils.ValidatePhone(req.CustomerPhone) {
		utils.SendBadRequest(&c.Controller, "Invalid phone format", nil)
		return
	}

	// Validate date format
	if !utils.ValidateDate(req.BookingDate) {
		utils.SendBadRequest(&c.Controller, "Invalid date format. Use YYYY-MM-DD", nil)
		return
	}

	// Validate date range
	isValid, err := utils.ValidateDateRange(req.BookingDate, 30)
	if err != nil || !isValid {
		utils.SendBadRequest(&c.Controller, "Date must be within the next 30 days and not in the past", nil)
		return
	}

	// Verify court exists and is active
	court, err := models.GetCourtById(req.CourtId)
	if err != nil {
		utils.SendNotFound(&c.Controller, "Court not found")
		return
	}
	if court.Status != "active" {
		utils.SendBadRequest(&c.Controller, "Court is not available", nil)
		return
	}

	// Verify timeslot exists and is active
	timeslot, err := models.GetTimeslotById(req.TimeslotId)
	if err != nil {
		utils.SendNotFound(&c.Controller, "Timeslot not found")
		return
	}
	if !timeslot.IsActive {
		utils.SendBadRequest(&c.Controller, "Timeslot is not available", nil)
		return
	}

	// Check availability
	isAvailable, err := models.CheckAvailability(req.CourtId, req.TimeslotId, req.BookingDate)
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error checking availability", err.Error())
		return
	}
	if !isAvailable {
		utils.SendConflict(&c.Controller, "This court is already booked for the selected date and timeslot", nil)
		return
	}

	// Calculate expiration time (default 30 minutes)
	timeoutMinutes := 30
	if timeoutStr := os.Getenv("RESERVATION_TIMEOUT_MINUTES"); timeoutStr != "" {
		if minutes, err := strconv.Atoi(timeoutStr); err == nil {
			timeoutMinutes = minutes
		}
	}
	expiredAt := time.Now().Add(time.Duration(timeoutMinutes) * time.Minute)

	// Create reservation
	reservation := &models.Reservation{
		Id:            uuid.New().String(),
		CourtId:       req.CourtId,
		TimeslotId:    req.TimeslotId,
		BookingDate:   req.BookingDate,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		TotalPrice:    court.PricePerHour,
		Status:        "pending",
		Notes:         req.Notes,
		ExpiredAt:     expiredAt,
	}

	err = models.CreateReservation(reservation)
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error creating reservation", err.Error())
		return
	}

	// Get full reservation with relations
	fullReservation, _ := models.GetReservationById(reservation.Id)

	c.Ctx.Output.SetStatus(201)
	utils.SendSuccess(&c.Controller, "Reservation created successfully. Please complete payment within 30 minutes.", fullReservation)
}

// GetReservationById godoc
// @Summary Get reservation by ID
// @Description Get reservation details by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/reservations/{id} [get]
func (c *ReservationController) GetReservationById() {
	id := c.Ctx.Input.Param(":id")

	reservation, err := models.GetReservationById(id)
	if err != nil {
		utils.SendNotFound(&c.Controller, "Reservation not found")
		return
	}

	utils.SendSuccess(&c.Controller, "Reservation retrieved successfully", reservation)
}

// GetReservationsByEmail godoc
// @Summary Get reservations by email
// @Description Get all reservations for a customer email
// @Tags reservations
// @Accept json
// @Produce json
// @Param email query string true "Customer email"
// @Success 200 {object} utils.Response
// @Router /api/v1/reservations/customer [get]
func (c *ReservationController) GetReservationsByEmail() {
	email := c.GetString("email")

	if email == "" {
		utils.SendBadRequest(&c.Controller, "Email parameter is required", nil)
		return
	}

	if !utils.ValidateEmail(email) {
		utils.SendBadRequest(&c.Controller, "Invalid email format", nil)
		return
	}

	reservations, err := models.GetReservationsByEmail(email)
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error fetching reservations", err.Error())
		return
	}

	utils.SendSuccess(&c.Controller, "Reservations retrieved successfully", reservations)
}

// UpdateStatus updates reservation status (admin/test endpoint)
// @Summary Update reservation status
// @Description Update reservation status (e.g., expired, cancelled, paid)
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID"
// @Param body body UpdateStatusRequest true "Status payload"
// @Success 200 {object} utils.Response
// @Router /api/v1/reservations/{id}/status [post]
func (c *ReservationController) UpdateStatus() {
	id := c.Ctx.Input.Param(":id")
	var req UpdateStatusRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendBadRequest(&c.Controller, "Invalid request body", err.Error())
		return
	}

	if req.Status == "" {
		utils.SendBadRequest(&c.Controller, "status is required", nil)
		return
	}

	if err := models.UpdateReservationStatus(id, req.Status); err != nil {
		utils.SendInternalError(&c.Controller, "Error updating reservation status", err.Error())
		return
	}

	utils.SendSuccess(&c.Controller, "Reservation status updated", map[string]string{"id": id, "status": req.Status})
}
