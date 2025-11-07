package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Payment struct {
	Id             string    `orm:"column(id);pk" json:"id"`
	ReservationId  string    `orm:"column(reservation_id);size(64)" json:"reservation_id"`
	OrderId        string    `orm:"column(order_id);size(128);null" json:"order_id"`
	PaymentUrl     string    `orm:"column(payment_url);type(text);null" json:"payment_url"`
	Amount         float64   `orm:"column(amount);digits(10);decimals(2)" json:"amount"`
	PaymentGateway string    `orm:"column(payment_gateway);size(64)" json:"payment_gateway"`
	Status         string    `orm:"column(status);size(32)" json:"status"`
	TransactionId  string    `orm:"column(transaction_id);size(128);null" json:"transaction_id"`
	Notification   string    `orm:"column(notification);type(text);null" json:"notification"`
	ExpiredAt      time.Time `orm:"column(expired_at);type(datetime);null" json:"expired_at"`
	CreatedAt      time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt      time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (p *Payment) TableName() string {
	return "payments"
}

func init() {
	orm.RegisterModel(new(Payment))
}

// CreatePayment inserts a new payment record
func CreatePayment(p *Payment) error {
	o := orm.NewOrm()
	// Use raw insert to avoid LastInsertId issues on Postgres drivers
	_, err := o.Raw(`INSERT INTO payments (id, reservation_id, order_id, payment_url, amount, payment_gateway, status, transaction_id, notification, expired_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), now())`, p.Id, p.ReservationId, p.OrderId, p.PaymentUrl, p.Amount, p.PaymentGateway, p.Status, p.TransactionId, p.Notification, p.ExpiredAt).Exec()
	return err
}

// GetPaymentByReservationId returns payment by reservation id
func GetPaymentByReservationId(reservationId string) (*Payment, error) {
	o := orm.NewOrm()
	payment := &Payment{}
	err := o.QueryTable(new(Payment)).Filter("reservation_id", reservationId).One(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// GetPaymentByOrderId returns payment by order id
func GetPaymentByOrderId(orderId string) (*Payment, error) {
	o := orm.NewOrm()
	payment := &Payment{}
	err := o.QueryTable(new(Payment)).Filter("order_id", orderId).One(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// GetPaymentById returns payment by id
func GetPaymentById(id string) (*Payment, error) {
	o := orm.NewOrm()
	payment := &Payment{Id: id}
	err := o.Read(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePaymentStatus updates payment status and other related fields
func UpdatePaymentStatus(id string, status string, transactionId string, notification string) error {
	o := orm.NewOrm()
	p := &Payment{Id: id}
	if err := o.Read(p); err != nil {
		return err
	}
	p.Status = status
	if transactionId != "" {
		p.TransactionId = transactionId
	}
	if notification != "" {
		p.Notification = notification
	}
	_, err := o.Update(p)
	return err
}
