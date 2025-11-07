-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
	id VARCHAR(36) PRIMARY KEY,
	reservation_id VARCHAR(36) NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
	order_id VARCHAR(128),
	payment_url TEXT,
	amount NUMERIC(10,2) NOT NULL,
	payment_gateway VARCHAR(64) NOT NULL DEFAULT 'midtrans',
	status VARCHAR(32) NOT NULL DEFAULT 'pending',
	transaction_id VARCHAR(128),
	notification TEXT,
	expired_at TIMESTAMP NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payments_reservation_id ON payments(reservation_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

-- Trigger to keep updated_at current
CREATE TRIGGER update_payments_updated_at
	BEFORE UPDATE ON payments
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE payments IS 'Payment transactions related to reservations';

