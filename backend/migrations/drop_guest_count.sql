-- Migration: Drop guest_count column from reservations table
-- เพราะโมเดล Go ไม่มีฟิลด์นี้

ALTER TABLE reservations DROP COLUMN guest_count;
