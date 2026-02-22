-- Migration: Add price column to tables table
-- เพิ่มคอลัมน์ price สำหรับเก็บราคาต่อโต๊ะ

-- เพิ่มคอลัมน์ price ถ้ายังไม่มี
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'tables' AND column_name = 'price'
    ) THEN
        ALTER TABLE tables ADD COLUMN price DOUBLE PRECISION DEFAULT 0;
        RAISE NOTICE 'Added price column to tables table';
    ELSE
        RAISE NOTICE 'Column price already exists in tables table';
    END IF;
END $$;
