package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Simple tool to simulate a Midtrans webhook (callback) to the local backend.
// Usage examples:
//  go run tools/send_webhook/main.go -reservation <reservation-id>
//  go run tools/send_webhook/main.go -order RES-... -amount 100000

func main() {
	var (
		reservation = flag.String("reservation", "", "reservation id to fetch payment/order_id from backend")
		order       = flag.String("order", "", "order_id to use in simulated notification")
		amount      = flag.String("amount", "", "gross amount as string (e.g. 100000)")
		urlFlag     = flag.String("url", "http://localhost:8080/api/v1/payments/callback", "callback URL to POST to")
		serverKey   = flag.String("server-key", os.Getenv("MIDTRANS_SERVER_KEY"), "Midtrans server key (overrides MIDTRANS_SERVER_KEY env)")
		txnId       = flag.String("tx", "SIM-TX-001", "transaction id to include in notification")
	)
	flag.Parse()

	if *reservation == "" && *order == "" {
		fmt.Fprintln(os.Stderr, "Either -reservation or -order must be provided")
		flag.Usage()
		os.Exit(2)
	}
	// If reservation provided, try to fetch payment info from local backend
	if *reservation != "" {
		apiUrl := fmt.Sprintf("http://localhost:8080/api/v1/payments/%s", *reservation)
		resp, err := http.Get(apiUrl)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error fetching payment by reservation:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Fprintln(os.Stderr, "error parsing response:", err)
			fmt.Fprintln(os.Stderr, "body:", string(body))
			os.Exit(1)
		}
		// Expect shape: { success: true, data: { order_id:..., amount:... } }
		if d, ok := data["data"].(map[string]interface{}); ok {
			if o, ok := d["order_id"].(string); ok && *order == "" {
				*order = o
			}
			if a, ok := d["amount"]; ok && *amount == "" {
				// amount may be number
				switch v := a.(type) {
				case float64:
					*amount = fmt.Sprintf("%d", int64(v))
				case string:
					*amount = v
				default:
					*amount = fmt.Sprintf("%v", v)
				}
			}
		} else {
			fmt.Fprintln(os.Stderr, "unexpected response shape, body:", string(body))
			os.Exit(1)
		}
	}

	if *order == "" || *amount == "" {
		fmt.Fprintln(os.Stderr, "missing order or amount; got order=", *order, " amount=", *amount)
		os.Exit(1)
	}

	if *serverKey == "" {
		fmt.Fprintln(os.Stderr, "MIDTRANS server key is empty; set MIDTRANS_SERVER_KEY env or pass -server-key flag")
		os.Exit(1)
	}

	statusCode := "200"
	transactionStatus := "settlement"
	fraudStatus := "accept"

	// compute signature: sha512(order_id + status_code + gross_amount + serverKey)
	h := sha512.Sum512([]byte(*order + statusCode + *amount + *serverKey))
	signature := hex.EncodeToString(h[:])

	payload := map[string]interface{}{
		"transaction_status": transactionStatus,
		"order_id":           *order,
		"gross_amount":       *amount,
		"payment_type":       "mock",
		"transaction_id":     *txnId,
		"fraud_status":       fraudStatus,
		"status_code":        statusCode,
		"signature_key":      signature,
	}

	b, _ := json.Marshal(payload)
	fmt.Println("POST", *urlFlag)
	fmt.Println("payload:", string(b))

	resp, err := http.Post(*urlFlag, "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error posting webhook:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("response status:", resp.Status)
	fmt.Println("response body:", string(respBody))
}
