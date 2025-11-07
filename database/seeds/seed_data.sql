-- Seed data for Badminton Reservation API

-- Courts (idempotent inserts)
INSERT INTO courts (name, description, price_per_hour, status)
SELECT 'Court A', 'Indoor premium court', 100000.00, 'active'
WHERE NOT EXISTS (SELECT 1 FROM courts WHERE name = 'Court A');

INSERT INTO courts (name, description, price_per_hour, status)
SELECT 'Court B', 'Standard indoor court', 75000.00, 'active'
WHERE NOT EXISTS (SELECT 1 FROM courts WHERE name = 'Court B');

INSERT INTO courts (name, description, price_per_hour, status)
SELECT 'Court C', 'Outdoor court', 50000.00, 'active'
WHERE NOT EXISTS (SELECT 1 FROM courts WHERE name = 'Court C');

-- Timeslots (idempotent inserts by start_time & end_time)
INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '07:00:00', '08:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='07:00:00' AND end_time='08:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '08:00:00', '09:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='08:00:00' AND end_time='09:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '09:00:00', '10:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='09:00:00' AND end_time='10:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '10:00:00', '11:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='10:00:00' AND end_time='11:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '11:00:00', '12:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='11:00:00' AND end_time='12:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '13:00:00', '14:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='13:00:00' AND end_time='14:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '14:00:00', '15:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='14:00:00' AND end_time='15:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '15:00:00', '16:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='15:00:00' AND end_time='16:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '16:00:00', '17:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='16:00:00' AND end_time='17:00:00');

INSERT INTO timeslots (start_time, end_time, is_active)
SELECT '18:00:00', '19:00:00', true
WHERE NOT EXISTS (SELECT 1 FROM timeslots WHERE start_time='18:00:00' AND end_time='19:00:00');
-- Verify data (counts)
SELECT 'Courts total:' as info, COUNT(*) as count FROM courts;
SELECT 'Timeslots total:' as info, COUNT(*) as count FROM timeslots;