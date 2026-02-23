-- Migration: Drop name column from tables table
-- ลบคอลัมน์ name เพราะมีค่าเหมือนกับ number

-- ลบคอลัมน์ name ถ้ามีอยู่
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'tables' AND column_name = 'name'
    ) THEN
        ALTER TABLE tables DROP COLUMN name;
        RAISE NOTICE 'Dropped name column from tables table';
    ELSE
        RAISE NOTICE 'Column name does not exist in tables table';
    END IF;
END $$;
