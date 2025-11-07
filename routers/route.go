package routers

import (
	"badminton-reservation-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// API v1 namespace
	ns := web.NewNamespace("/api/v1",
		// Date routes
		web.NSRouter("/dates", &controllers.DateController{}, "get:GetAvailableDates"),

		// Timeslot routes
		web.NSRouter("/timeslots", &controllers.TimeslotController{}, "get:GetAvailableTimeslots"),
		web.NSRouter("/timeslots/all", &controllers.TimeslotController{}, "get:GetAllTimeslots"),

		// Court routes
		web.NSRouter("/courts", &controllers.CourtController{}, "get:GetAvailableCourts"),
		web.NSRouter("/courts/all", &controllers.CourtController{}, "get:GetAllCourts"),

		// Reservation routes
		web.NSRouter("/reservations", &controllers.ReservationController{}, "post:CreateReservation"),
		web.NSRouter("/reservations/:id", &controllers.ReservationController{}, "get:GetReservationById"),
		web.NSRouter("/reservations/customer", &controllers.ReservationController{}, "get:GetReservationsByEmail"),

		// Payment routes
		web.NSRouter("/payments/process", &controllers.PaymentController{}, "post:ProcessPayment"),
		web.NSRouter("/payments/callback", &controllers.PaymentController{}, "post:PaymentCallback"),
		web.NSRouter("/payments/:id", &controllers.PaymentController{}, "get:GetPaymentStatus"),
	)

	web.AddNamespace(ns)

	// Health check endpoint
	web.Router("/health", &controllers.HealthController{}, "get:Get")

	// Swagger UI (simple embedded UI)
	web.Router("/swagger", &controllers.SwaggerUIController{}, "get:UI")
	web.Router("/swagger/doc.json", &controllers.SwaggerUIController{}, "get:Doc")
}
