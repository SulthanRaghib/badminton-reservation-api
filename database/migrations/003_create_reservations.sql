-- Create reservations table
-- Note: BookingDate is stored as text (YYYY-MM-DD) to match model's string field
CREATE TABLE IF NOT EXISTS reservations (
	id VARCHAR(36) PRIMARY KEY,
	court_id INTEGER NOT NULL REFERENCES courts(id) ON DELETE RESTRICT,
	timeslot_id INTEGER NOT NULL REFERENCES timeslots(id) ON DELETE RESTRICT,
	booking_date VARCHAR(10) NOT NULL,
	customer_name VARCHAR(255) NOT NULL,
	customer_email VARCHAR(255) NOT NULL,
	customer_phone VARCHAR(50) NOT NULL,
	total_price NUMERIC(10,2) NOT NULL,
	status VARCHAR(32) NOT NULL DEFAULT 'pending',
	notes TEXT,
	expired_at TIMESTAMP NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Helpful indexes for common queries
CREATE INDEX IF NOT EXISTS idx_reservations_booking_date ON reservations(booking_date);
CREATE INDEX IF NOT EXISTS idx_reservations_court_date_timeslot ON reservations(court_id, booking_date, timeslot_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);

-- Trigger to keep updated_at current
CREATE TRIGGER update_reservations_updated_at
	BEFORE UPDATE ON reservations
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE reservations IS 'Court reservations - one row per booking attempt/transaction';

