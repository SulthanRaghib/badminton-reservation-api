-- Create timeslots table
CREATE TABLE IF NOT EXISTS timeslots (
	id SERIAL PRIMARY KEY,
	start_time VARCHAR(10) NOT NULL,
	end_time VARCHAR(10) NOT NULL,
	is_active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_timeslots_is_active ON timeslots(is_active);

-- Trigger to update updated_at on row modification (reuses function from courts migration)
CREATE TRIGGER update_timeslots_updated_at
	BEFORE UPDATE ON timeslots
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE timeslots IS 'Timeslot definitions (start_time, end_time)';

