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
	_, err := o.QueryTable(new(Timeslot)).Filter("is_active", true).All(&slots)
	return slots, err
}
