package main

import (
	"os"
	"time"

	"badminton-reservation-api/database"
	"badminton-reservation-api/middleware"
	"badminton-reservation-api/models"
	_ "badminton-reservation-api/routers"
	"badminton-reservation-api/utils"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var dbInitialized bool

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logs.Warning("No .env file found, using environment variables")
	}

	// Optionally skip DB initialization (dev convenience). Set SKIP_DB=true to start the
	// server without attempting to connect to the database (useful when you only need
	// to view Swagger UI or work on non-DB endpoints).
	if os.Getenv("SKIP_DB") == "true" {
		logs.Warning("SKIP_DB=true, skipping database initialization and background jobs")
	} else {
		// Initialize database connection (do not hard-panic on failure — allow dev to continue)
		if err := initDatabase(); err != nil {
			logs.Warning("Database initialization failed, continuing with SKIP_DB behavior:", err)
			dbInitialized = false
		} else {
			dbInitialized = true
			// Start background job to expire old reservations only when DB is initialized
			go expireReservationsJob()
		}
	}

	// Setup CORS middleware
	web.InsertFilter("*", web.BeforeRouter, middleware.CORS)
}

// initDatabase initializes DB and returns an error instead of panicking so the
// caller can decide how to proceed when DB is unreachable.
func initDatabase() error {
	// Derive data source: prefer DB_URL; fall back to individual env vars
	dataSource := utils.GetDataSource()

	// Register database driver
	if err := orm.RegisterDriver("postgres", orm.DRPostgres); err != nil {
		logs.Error("Failed to register database driver:", err)
		return err
	}

	// Register database
	if err := orm.RegisterDataBase("default", "postgres", dataSource); err != nil {
		logs.Error("Failed to register database:", err)
		return err
	}

	// Set database parameters
	orm.SetMaxIdleConns("default", 10)
	orm.SetMaxOpenConns("default", 100)

	// Test database connection
	if err := utils.TestDBConnection(dataSource, 5*time.Second); err != nil {
		logs.Error("Failed to connect to database:", err)
		return err
	}

	logs.Info("Database connected successfully")

	// Enable debug mode in development
	if os.Getenv("APP_ENV") == "development" {
		orm.Debug = true
	}

	// Optionally run GORM migrations to apply schema (useful for cloud Postgres / Neon migrations)
	// Enable by setting MIGRATE_WITH_GORM=true in environment
	if os.Getenv("MIGRATE_WITH_GORM") == "true" {
		logs.Info("MIGRATE_WITH_GORM=true, running GORM AutoMigrate...")
		if err := database.RunGormMigrations(); err != nil {
			logs.Error("GORM migration failed:", err)
			// Do not treat migration failure as fatal for startup
		} else {
			logs.Info("GORM migrations applied successfully")
		}
	}
	return nil
}

// expireReservationsJob runs periodically to expire old pending reservations
func expireReservationsJob() {
	ticker := time.NewTicker(5 * time.Minute) // Run every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		logs.Info("Running reservation expiration job...")
		err := models.ExpireOldReservations()
		if err != nil {
			logs.Error("Error expiring reservations:", err)
		} else {
			logs.Info("Reservation expiration job completed")
		}
	}
}

func main() {
	// Get port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Set Beego configurations
	web.BConfig.CopyRequestBody = true
	web.BConfig.RunMode = os.Getenv("APP_ENV")
	if web.BConfig.RunMode == "" {
		web.BConfig.RunMode = "dev"
	}

	// Print startup information
	logs.Info("========================================")
	logs.Info("🏸 Badminton Reservation API")
	logs.Info("========================================")
	logs.Info("Environment:", web.BConfig.RunMode)
	logs.Info("Port:", port)
	logs.Info("API Base URL:", os.Getenv("APP_URL"))
	logs.Info("========================================")
	logs.Info("API Endpoints:")
	logs.Info("  GET  /health")
	logs.Info("  GET  /api/v1/dates")
	logs.Info("      - Query: none (returns next N available dates, see MAX_BOOKING_DAYS_AHEAD env)")
	logs.Info("  GET  /api/v1/timeslots?booking_date=YYYY-MM-DD&court_id=X")
	logs.Info("      - Params: booking_date (required), court_id (required)")
	logs.Info("      - Response: returns globally active timeslots and an 'available' boolean per timeslot; available=false means already booked for that date and court")
	logs.Info("  GET  /api/v1/timeslots/all")
	logs.Info("      - Returns all defined timeslots")
	logs.Info("  GET  /api/v1/courts?booking_date=YYYY-MM-DD&timeslot_id=X")
	logs.Info("      - Params: booking_date (required), timeslot_id (required)")
	logs.Info("  GET  /api/v1/courts/all")
	logs.Info("      - Returns all active courts")
	logs.Info("  POST /api/v1/reservations")
	logs.Info("      - Body: {court_id,timeslot_id,booking_date,customer_name,customer_email,customer_phone,notes}")
	logs.Info("  GET  /api/v1/reservations/:id")
	logs.Info("  GET  /api/v1/reservations/customer?email=you@example.com")
	logs.Info("      - Query: email (required)")
	logs.Info("  POST /api/v1/payments/process")
	logs.Info("      - Body: {reservation_id}")
	logs.Info("  POST /api/v1/payments/callback")
	logs.Info("      - Webhook endpoint for payment notifications from the gateway")
	logs.Info("  GET  /api/v1/payments/:id")
	logs.Info("      - id can be a payment ID or reservation ID")
	logs.Info("========================================")

	// Run the application
	web.Run(":" + port)
}
