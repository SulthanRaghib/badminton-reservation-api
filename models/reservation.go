package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Reservation struct {
	Id            string    `orm:"column(id);pk" json:"id"`
	CourtId       int       `orm:"column(court_id)" json:"court_id"`
	TimeslotId    int       `orm:"column(timeslot_id)" json:"timeslot_id"`
	BookingDate   string    `orm:"column(booking_date);size(10)" json:"booking_date"`
	CustomerName  string    `orm:"column(customer_name);size(255)" json:"customer_name"`
	CustomerEmail string    `orm:"column(customer_email);size(255)" json:"customer_email"`
	CustomerPhone string    `orm:"column(customer_phone);size(50)" json:"customer_phone"`
	TotalPrice    float64   `orm:"column(total_price);digits(10);decimals(2)" json:"total_price"`
	Status        string    `orm:"column(status);size(32)" json:"status"`
	Notes         string    `orm:"column(notes);type(text);null" json:"notes"`
	ExpiredAt     time.Time `orm:"column(expired_at);type(datetime);null" json:"expired_at"`
	CreatedAt     time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt     time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (r *Reservation) TableName() string {
	return "reservations"
}

func init() {
	orm.RegisterModel(new(Reservation))
}

// CreateReservation inserts a new reservation record
func CreateReservation(r *Reservation) error {
	o := orm.NewOrm()
	// Use raw insert to avoid drivers that do not support LastInsertId for Postgres
	_, err := o.Raw(`INSERT INTO reservations (id, court_id, timeslot_id, booking_date, customer_name, customer_email, customer_phone, total_price, status, notes, expired_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), now())`, r.Id, r.CourtId, r.TimeslotId, r.BookingDate, r.CustomerName, r.CustomerEmail, r.CustomerPhone, r.TotalPrice, r.Status, r.Notes, r.ExpiredAt).Exec()
	return err
}

// GetReservationById returns reservation by id
func GetReservationById(id string) (*Reservation, error) {
	o := orm.NewOrm()
	res := &Reservation{Id: id}
	err := o.Read(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetReservationsByEmail returns all reservations for a given email
func GetReservationsByEmail(email string) ([]*Reservation, error) {
	o := orm.NewOrm()
	var list []*Reservation
	_, err := o.QueryTable(new(Reservation)).Filter("customer_email", email).All(&list)
	return list, err
}

// UpdateReservationStatus updates status for a reservation
func UpdateReservationStatus(id string, status string) error {
	o := orm.NewOrm()
	r := &Reservation{Id: id}
	if err := o.Read(r); err != nil {
		return err
	}
	r.Status = status
	_, err := o.Update(r)
	return err
}

// CheckAvailability checks if a court is available for a given timeslot and date
func CheckAvailability(courtId int, timeslotId int, bookingDate string) (bool, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(new(Reservation)).Filter("court_id", courtId).Filter("timeslot_id", timeslotId).Filter("booking_date", bookingDate).Filter("status__in", "pending", "waiting_payment", "paid").Count()
	if err != nil {
		return false, err
	}
	return cnt == 0, nil
}

// ExpireOldReservations will mark pending reservations whose ExpiredAt is before now as expired
func ExpireOldReservations() error {
	o := orm.NewOrm()
	now := time.Now()
	// Raw update: set status to 'expired' where status = 'pending' and expired_at < now
	_, err := o.Raw("UPDATE reservations SET status = ? WHERE status = ? AND expired_at < ?", "expired", "pending", now).Exec()
	return err
}
