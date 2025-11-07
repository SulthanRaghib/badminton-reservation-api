package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type Timeslot struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	StartTime string `orm:"column(start_time);size(10)" json:"start_time"`
	EndTime   string `orm:"column(end_time);size(10)" json:"end_time"`
	IsActive  bool   `orm:"column(is_active);default(true)" json:"is_active"`
}

func (t *Timeslot) TableName() string {
	return "timeslots"
}

func init() {
	orm.RegisterModel(new(Timeslot))
}

// GetTimeslotById returns a timeslot by id
func GetTimeslotById(id int) (*Timeslot, error) {
	o := orm.NewOrm()
	t := &Timeslot{Id: id}
	if err := o.Read(t); err != nil {
		return nil, err
	}
	return t, nil
}

// GetAllTimeslots returns all timeslots
func GetAllTimeslots() ([]*Timeslot, error) {
	o := orm.NewOrm()
	var list []*Timeslot
	_, err := o.QueryTable(new(Timeslot)).OrderBy("id").All(&list)
	return list, err
}

// GetAvailableTimeslots returns timeslots that are active and not booked for the given court/date
// For simplicity this function only returns active timeslots; controllers may filter further
func GetAvailableTimeslots(courtId int, bookingDate string) ([]*Timeslot, error) {
	o := orm.NewOrm()
	var slots []*Timeslot
	// Return timeslots that are globally active and NOT marked unavailable for this court/date
	_, err := o.Raw(`SELECT t.id, t.start_time, t.end_time, t.is_active
		FROM timeslots t
		WHERE t.is_active = true
		AND t.id NOT IN (
			SELECT timeslot_id FROM timeslot_availabilities
			WHERE court_id = ? AND booking_date = ? AND is_active = false
		)
		ORDER BY t.id`, courtId, bookingDate).QueryRows(&slots)
	return slots, err
}
