package controllers

import (
	"badminton-reservation-api/models"
	"badminton-reservation-api/services/payment"
	"badminton-reservation-api/utils"
	"encoding/json"
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
)

type PaymentController struct {
	web.Controller
}

type ProcessPaymentRequest struct {
	ReservationId string `json:"reservation_id"`
}

type ProcessPaymentResponse struct {
	PaymentId   string `json:"payment_id"`
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}

// ProcessPayment godoc
// @Summary Process payment for a reservation
// @Description Create payment transaction with Midtrans
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body ProcessPaymentRequest true "Payment details"
// @Success 200 {object} utils.Response
// @Router /api/v1/payments/process [post]
func (c *PaymentController) ProcessPayment() {
	var req ProcessPaymentRequest

	// Parse request body
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		utils.SendBadRequest(&c.Controller, "Invalid request body", err.Error())
		return
	}

	// Validate reservation_id
	if req.ReservationId == "" {
		utils.SendBadRequest(&c.Controller, "Reservation ID is required", nil)
		return
	}

	// Get reservation
	reservation, err := models.GetReservationById(req.ReservationId)
	if err != nil {
		utils.SendNotFound(&c.Controller, "Reservation not found")
		return
	}

	// Check if reservation is still pending
	if reservation.Status != "pending" {
		utils.SendBadRequest(&c.Controller, "Reservation is not in pending status", nil)
		return
	}

	// Check if reservation has expired
	if time.Now().After(reservation.ExpiredAt) {
		// Update status to expired
		_ = models.UpdateReservationStatus(reservation.Id, "expired")
		utils.SendBadRequest(&c.Controller, "Reservation has expired", nil)
		return
	}

	// Check if payment already exists
	existingPayment, err := models.GetPaymentByReservationId(reservation.Id)
	if err == nil && existingPayment != nil {
		// Payment already exists, return existing payment details
		response := ProcessPaymentResponse{
			PaymentId:   existingPayment.Id,
			Token:       "", // Token might be expired
			RedirectUrl: existingPayment.PaymentUrl,
		}
		utils.SendSuccess(&c.Controller, "Payment already exists", response)
		return
	}

	// Create payment record
	paymentRecord := &models.Payment{
		Id:             uuid.New().String(),
		ReservationId:  reservation.Id,
		Amount:         reservation.TotalPrice,
		PaymentGateway: "midtrans",
		Status:         "pending",
		ExpiredAt:      reservation.ExpiredAt,
	}

	// Create Midtrans transaction
	midtransService := payment.NewMidtransService()
	_, err = midtransService.CreateTransaction(reservation, paymentRecord)
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error creating payment transaction", err.Error())
		return
	}

	// Save payment record
	// Ensure paymentRecord has OrderId and PaymentUrl (set by midtrans service)
	err = models.CreatePayment(paymentRecord)
	if err != nil {
		utils.SendInternalError(&c.Controller, "Error saving payment record", err.Error())
		return
	}

	// Update reservation status to waiting_payment
	err = models.UpdateReservationStatus(reservation.Id, "waiting_payment")
	if err != nil {
		logs.Error("Error updating reservation status:", err)
	}

	// Return the full payment record so client immediately gets id + order/payment link
	utils.SendSuccess(&c.Controller, "Payment transaction created successfully", paymentRecord)
}

// PaymentCallback godoc
// @Summary Handle payment callback from Midtrans
// @Description Webhook endpoint for Midtrans payment notification
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/payments/callback [post]
func (c *PaymentController) PaymentCallback() {
	var notification map[string]interface{}

	// Log raw body for debugging
	raw := c.Ctx.Input.RequestBody
	logs.Info("PaymentCallback raw body:", string(raw))

	// If body is empty, try to parse form values (some webhook send form-encoded)
	if len(raw) == 0 {
		// Try ParseForm and extract values
		if err := c.Ctx.Request.ParseForm(); err == nil && len(c.Ctx.Request.Form) > 0 {
			notification = make(map[string]interface{})
			for k, v := range c.Ctx.Request.Form {
				if len(v) > 0 {
					notification[k] = v[0]
				}
			}
			logs.Info("Parsed notification from form values", notification)
		} else {
			logs.Error("Empty notification body and no form values")
			utils.SendBadRequest(&c.Controller, "Invalid notification", "empty request body")
			return
		}
	} else {
		// Try JSON first
		if err := json.Unmarshal(raw, &notification); err != nil {
			// If JSON unmarshal fails, attempt to parse as url-encoded form in the raw body
			logs.Warn("Error parsing notification JSON:", err, " — attempting url-encoded parse")
			if vals, perr := url.ParseQuery(string(raw)); perr == nil && len(vals) > 0 {
				notification = make(map[string]interface{})
				for k, v := range vals {
					if len(v) > 0 {
						notification[k] = v[0]
					}
				}
				logs.Info("Parsed notification from url-encoded body", notification)
			} else {
				logs.Error("Error parsing notification:", err)
				utils.SendBadRequest(&c.Controller, "Invalid notification", err.Error())
				return
			}
		} else {
			logs.Info("Received payment notification:", notification)
		}
	}

	// Initialize Midtrans service
	midtransService := payment.NewMidtransService()

	// Parse and verify notification
	statusResp, err := midtransService.ParseNotification(notification)
	if err != nil {
		logs.Error("Error parsing notification:", err)
		utils.SendBadRequest(&c.Controller, "Invalid notification signature", err.Error())
		return
	}

	// Get payment by order_id
	paymentRecord, err := models.GetPaymentByOrderId(statusResp.OrderID)
	if err != nil {
		logs.Error("Payment not found for order:", statusResp.OrderID)
		utils.SendNotFound(&c.Controller, "Payment not found")
		return
	}

	// Determine payment status
	paymentStatus := payment.GetPaymentStatus(statusResp.TransactionStatus, statusResp.FraudStatus)

	// Convert notification to JSON string
	notificationJSON, _ := json.Marshal(notification)

	// Update payment status
	err = models.UpdatePaymentStatus(
		paymentRecord.Id,
		paymentStatus,
		statusResp.TransactionID,
		string(notificationJSON),
	)
	if err != nil {
		logs.Error("Error updating payment status:", err)
		utils.SendInternalError(&c.Controller, "Error updating payment status", err.Error())
		return
	}

	// Update reservation status based on payment status
	var reservationStatus string
	switch paymentStatus {
	case "success":
		reservationStatus = "paid"
	case "failed":
		reservationStatus = "cancelled"
	default:
		reservationStatus = "waiting_payment"
	}

	err = models.UpdateReservationStatus(paymentRecord.ReservationId, reservationStatus)
	if err != nil {
		logs.Error("Error updating reservation status:", err)
	}

	utils.SendSuccess(&c.Controller, "Payment notification processed successfully", map[string]string{
		"status": paymentStatus,
	})
}

// GetPaymentStatus godoc
// @Summary Get payment status
// @Description Get payment status by payment ID or reservation ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID or Reservation ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/payments/{id} [get]
func (c *PaymentController) GetPaymentStatus() {
	id := c.Ctx.Input.Param(":id")

	// Try to get payment by payment ID first
	paymentRecord, err := models.GetPaymentById(id)
	if err != nil {
		// If not found, try by reservation ID
		paymentRecord, err = models.GetPaymentByReservationId(id)
		if err != nil {
			utils.SendNotFound(&c.Controller, "Payment not found")
			return
		}
	}

	utils.SendSuccess(&c.Controller, "Payment status retrieved successfully", paymentRecord)
}
