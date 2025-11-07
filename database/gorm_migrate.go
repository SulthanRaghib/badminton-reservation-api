package database

import (
	"fmt"
	"os"
	"time"

	"badminton-reservation-api/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORM models used only for migrations (mirror of existing models)
type GormCourt struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	Name         string    `gorm:"column:name;size:100;not null" json:"name"`
	Description  string    `gorm:"column:description;type:text" json:"description"`
	PricePerHour float64   `gorm:"column:price_per_hour;type:numeric(10,2);not null" json:"price_per_hour"`
	Status       string    `gorm:"column:status;size:20;default:active" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (GormCourt) TableName() string { return "courts" }

type GormTimeslot struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	StartTime string    `gorm:"column:start_time;size:10;not null" json:"start_time"`
	EndTime   string    `gorm:"column:end_time;size:10;not null" json:"end_time"`
	IsActive  bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (GormTimeslot) TableName() string { return "timeslots" }

type GormReservation struct {
	Id            string    `gorm:"primaryKey;column:id;size:36" json:"id"`
	CourtId       uint      `gorm:"column:court_id;not null" json:"court_id"`
	TimeslotId    uint      `gorm:"column:timeslot_id;not null" json:"timeslot_id"`
	BookingDate   string    `gorm:"column:booking_date;size:10;not null" json:"booking_date"`
	CustomerName  string    `gorm:"column:customer_name;size:255;not null" json:"customer_name"`
	CustomerEmail string    `gorm:"column:customer_email;size:255;not null" json:"customer_email"`
	CustomerPhone string    `gorm:"column:customer_phone;size:50;not null" json:"customer_phone"`
	TotalPrice    float64   `gorm:"column:total_price;type:numeric(10,2);not null" json:"total_price"`
	Status        string    `gorm:"column:status;size:32;default:pending" json:"status"`
	Notes         string    `gorm:"column:notes;type:text" json:"notes"`
	ExpiredAt     time.Time `gorm:"column:expired_at" json:"expired_at"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (GormReservation) TableName() string { return "reservations" }

type GormPayment struct {
	Id             string    `gorm:"primaryKey;column:id;size:36" json:"id"`
	ReservationId  string    `gorm:"column:reservation_id;size:36;not null" json:"reservation_id"`
	OrderId        string    `gorm:"column:order_id;size:128" json:"order_id"`
	PaymentUrl     string    `gorm:"column:payment_url;type:text" json:"payment_url"`
	Amount         float64   `gorm:"column:amount;type:numeric(10,2);not null" json:"amount"`
	PaymentGateway string    `gorm:"column:payment_gateway;size:64;default:midtrans" json:"payment_gateway"`
	Status         string    `gorm:"column:status;size:32;default:pending" json:"status"`
	TransactionId  string    `gorm:"column:transaction_id;size:128" json:"transaction_id"`
	Notification   string    `gorm:"column:notification;type:text" json:"notification"`
	ExpiredAt      time.Time `gorm:"column:expired_at" json:"expired_at"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (GormPayment) TableName() string { return "payments" }

// RunGormMigrations connects using GORM and runs AutoMigrate for the migration models.
// It reads DB connection info from env vars: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
func RunGormMigrations() error {
	// Prefer full DB_URL if present (works with Neon). utils.GetDataSource trims quotes.
	dsn := utils.GetDataSource()
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "require"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			host, user, password, dbname, port, sslmode)
	}

	// Use a simple logger to avoid noisy output in production
	gormLogger := logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return fmt.Errorf("failed to connect via gorm: %w", err)
	}

	// Auto-migrate tables
	if err := db.AutoMigrate(&GormCourt{}, &GormTimeslot{}, &GormReservation{}, &GormPayment{}); err != nil {
		return fmt.Errorf("gorm automigrate error: %w", err)
	}

	return nil
}
