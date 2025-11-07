-- Create timeslot_availabilities table: stores per-court/timeslot/date availability
CREATE TABLE IF NOT EXISTS timeslot_availabilities (
    id SERIAL PRIMARY KEY,
    court_id INTEGER NOT NULL REFERENCES courts(id) ON DELETE CASCADE,
    timeslot_id INTEGER NOT NULL REFERENCES timeslots(id) ON DELETE CASCADE,
    booking_date VARCHAR(10) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Unique constraint so one row per court/timeslot/date
CREATE UNIQUE INDEX IF NOT EXISTS idx_timeslot_avail_unique ON timeslot_availabilities(court_id, timeslot_id, booking_date);

-- Trigger to update updated_at on row modification
CREATE TRIGGER update_timeslot_availabilities_updated_at
    BEFORE UPDATE ON timeslot_availabilities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE timeslot_availabilities IS 'Per-date availability flags for timeslots per court';
