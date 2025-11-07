package payment

import (
	"badminton-reservation-api/models"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService struct {
	Client snap.Client
}

type MidtransResponse struct {
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}

// NewMidtransService creates a new Midtrans service instance
func NewMidtransService() *MidtransService {
	var client snap.Client
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// Set to Production if needed
	if os.Getenv("MIDTRANS_IS_PRODUCTION") == "true" {
		client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Production)
	}

	return &MidtransService{
		Client: client,
	}
}

// CreateTransaction creates a Snap transaction for payment
func (s *MidtransService) CreateTransaction(reservation *models.Reservation, payment *models.Payment) (*MidtransResponse, error) {
	// Support a mock mode for local testing without Midtrans API key
	if os.Getenv("MIDTRANS_MOCK") == "true" {
		orderId := fmt.Sprintf("MOCK-%s-%d", reservation.Id[:8], time.Now().Unix())
		// Use a deterministic mock token and redirect URL
		mockToken := fmt.Sprintf("mock-token-%s", reservation.Id[:8])
		mockUrl := fmt.Sprintf("https://mock-pay.example.com/redirect/%s", orderId)

		payment.OrderId = orderId
		payment.PaymentUrl = mockUrl

		response := &MidtransResponse{
			Token:       mockToken,
			RedirectUrl: mockUrl,
		}
		return response, nil
	}

	// Generate unique order ID
	orderId := fmt.Sprintf("RES-%s-%d", reservation.Id[:8], time.Now().Unix())

	// Create Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId,
			GrossAmt: int64(reservation.TotalPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: reservation.CustomerName,
			Email: reservation.CustomerEmail,
			Phone: reservation.CustomerPhone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    strconv.Itoa(reservation.CourtId),
				Name:  fmt.Sprintf("Court Booking - %s", reservation.BookingDate),
				Price: int64(reservation.TotalPrice),
				Qty:   1,
			},
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Callbacks: &snap.Callbacks{
			Finish: os.Getenv("APP_URL") + "/payment/finish",
		},
	}

	// Create transaction
	snapResp, err := s.Client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	// Update payment with order_id and payment_url
	payment.OrderId = orderId
	payment.PaymentUrl = snapResp.RedirectURL

	response := &MidtransResponse{
		Token:       snapResp.Token,
		RedirectUrl: snapResp.RedirectURL,
	}

	return response, nil
}

// VerifySignature verifies the signature from Midtrans notification
func (s *MidtransService) VerifySignature(orderId, statusCode, grossAmount, serverKey, signatureKey string) bool {
	// Midtrans signature: SHA512(order_id+status_code+gross_amount+ServerKey)
	h := sha512.Sum512([]byte(orderId + statusCode + grossAmount + serverKey))
	expectedSignature := hex.EncodeToString(h[:])
	return expectedSignature == signatureKey
}

// TransactionStatusResponse represents the response from Midtrans transaction status API
type TransactionStatusResponse struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionID     string `json:"transaction_id"`
	FraudStatus       string `json:"fraud_status"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
}

// ParseNotification parses and validates the notification from Midtrans
func (s *MidtransService) ParseNotification(notificationPayload map[string]interface{}) (*TransactionStatusResponse, error) {
	// Convert map to JSON
	jsonData, err := json.Marshal(notificationPayload)
	if err != nil {
		return nil, err
	}

	// Parse to struct
	var statusResp TransactionStatusResponse
	err = json.Unmarshal(jsonData, &statusResp)
	if err != nil {
		return nil, err
	}

	// Verify signature
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isValid := s.VerifySignature(
		statusResp.OrderID,
		statusResp.StatusCode,
		statusResp.GrossAmount,
		serverKey,
		statusResp.SignatureKey,
	)

	if !isValid {
		return nil, fmt.Errorf("invalid signature")
	}

	return &statusResp, nil
}

// GetPaymentStatus returns the standardized payment status
func GetPaymentStatus(transactionStatus string, fraudStatus string) string {
	if transactionStatus == "capture" {
		if fraudStatus == "accept" {
			return "success"
		}
		return "pending"
	} else if transactionStatus == "settlement" {
		return "success"
	} else if transactionStatus == "pending" {
		return "pending"
	} else if transactionStatus == "deny" || transactionStatus == "expire" || transactionStatus == "cancel" {
		return "failed"
	}
	return "pending"
}
