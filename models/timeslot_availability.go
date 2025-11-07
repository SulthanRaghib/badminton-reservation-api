package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type TimeslotAvailability struct {
	Id          int       `orm:"column(id);auto;pk" json:"id"`
	CourtId     int       `orm:"column(court_id)" json:"court_id"`
	TimeslotId  int       `orm:"column(timeslot_id)" json:"timeslot_id"`
	BookingDate string    `orm:"column(booking_date);size(10)" json:"booking_date"`
	IsActive    bool      `orm:"column(is_active);default(true)" json:"is_active"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);auto_now_add" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);auto_now" json:"updated_at"`
}

func (t *TimeslotAvailability) TableName() string {
	return "timeslot_availabilities"
}

func init() {
	orm.RegisterModel(new(TimeslotAvailability))
}

// MarkTimeslotUnavailable marks a timeslot as unavailable (is_active=false) for a given court and date
func MarkTimeslotUnavailable(courtId int, timeslotId int, bookingDate string) error {
	o := orm.NewOrm()
	// Upsert: insert or update to set is_active = false
	_, err := o.Raw(`INSERT INTO timeslot_availabilities (court_id, timeslot_id, booking_date, is_active, created_at, updated_at)
        VALUES (?, ?, ?, false, now(), now())
        ON CONFLICT (court_id, timeslot_id, booking_date) DO UPDATE SET is_active = false, updated_at = now()`, courtId, timeslotId, bookingDate).Exec()
	return err
}

// MarkTimeslotAvailable marks timeslot as available (is_active=true) for a given court and date
func MarkTimeslotAvailable(courtId int, timeslotId int, bookingDate string) error {
	o := orm.NewOrm()
	// Upsert to set is_active = true
	_, err := o.Raw(`INSERT INTO timeslot_availabilities (court_id, timeslot_id, booking_date, is_active, created_at, updated_at)
		VALUES (?, ?, ?, true, now(), now())
		ON CONFLICT (court_id, timeslot_id, booking_date) DO UPDATE SET is_active = true, updated_at = now()`, courtId, timeslotId, bookingDate).Exec()
	return err
}

// RemoveAvailabilityRow deletes the availability row (optional) for cleanliness
func RemoveAvailabilityRow(courtId int, timeslotId int, bookingDate string) error {
	o := orm.NewOrm()
	_, err := o.Raw(`DELETE FROM timeslot_availabilities WHERE court_id = ? AND timeslot_id = ? AND booking_date = ?`, courtId, timeslotId, bookingDate).Exec()
	return err
}
