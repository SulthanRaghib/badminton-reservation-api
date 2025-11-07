package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Court struct {
	Id           int       `orm:"column(id);auto;pk" json:"id"`
	Name         string    `orm:"column(name);size(100)" json:"name"`
	Description  string    `orm:"column(description);type(text);null" json:"description"`
	PricePerHour float64   `orm:"column(price_per_hour);digits(10);decimals(2)" json:"price_per_hour"`
	Status       string    `orm:"column(status);size(20);default(active)" json:"status"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (c *Court) TableName() string {
	return "courts"
}

func init() {
	orm.RegisterModel(new(Court))
}

// GetAllActiveCourts retrieves all active courts
func GetAllActiveCourts() ([]*Court, error) {
	o := orm.NewOrm()
	var courts []*Court
	_, err := o.QueryTable(new(Court)).Filter("status", "active").All(&courts)
	return courts, err
}

// GetCourtById retrieves a court by ID
func GetCourtById(id int) (*Court, error) {
	o := orm.NewOrm()
	court := &Court{Id: id}
	err := o.Read(court)
	if err != nil {
		return nil, err
	}
	return court, nil
}

// GetAvailableCourts retrieves courts available for a specific date and timeslot
func GetAvailableCourts(bookingDate string, timeslotId int) ([]*Court, error) {
	o := orm.NewOrm()
	var courts []*Court

	// Get all active courts that are not reserved for the given date and timeslot
	_, err := o.Raw(`
		SELECT c.* FROM courts c
		WHERE c.status = 'active'
		AND c.id NOT IN (
			SELECT r.court_id FROM reservations r
			WHERE r.booking_date = ?
			AND r.timeslot_id = ?
			AND r.status IN ('pending', 'waiting_payment', 'paid')
		)
		ORDER BY c.id
	`, bookingDate, timeslotId).QueryRows(&courts)

	return courts, err
}
