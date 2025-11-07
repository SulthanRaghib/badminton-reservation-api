package main

import (
	mp "badminton-reservation-api/services/payment"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Use a test server key (can be from .env). Must be non-empty for signature.
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		serverKey = "SB-MIDTEST-SERVER-KEY"
	}

	orderID := "TEST-ORDER-123"
	statusCode := "200"
	grossAmount := "75000"

	// compute signature: sha512(order_id+status_code+gross_amount+serverKey)
	h := sha512.Sum512([]byte(orderID + statusCode + grossAmount + serverKey))
	signature := hex.EncodeToString(h[:])

	notification := map[string]interface{}{
		"transaction_status": "settlement",
		"order_id":           orderID,
		"gross_amount":       grossAmount,
		"payment_type":       "bank_transfer",
		"transaction_id":     "TRX-123456",
		"fraud_status":       "accept",
		"status_code":        statusCode,
		"signature_key":      signature,
	}

	// set env so ParseNotification uses same serverKey
	os.Setenv("MIDTRANS_SERVER_KEY", serverKey)

	svc := mp.NewMidtransService()

	// Parse and verify
	resp, err := svc.ParseNotification(notification)
	if err != nil {
		fmt.Println("ParseNotification error:", err)
		return
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println("Parsed notification:", string(b))
}
