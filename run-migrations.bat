@echo off
echo ========================================
echo   Run Database Migrations
echo ========================================
echo.

set PGPASSWORD=password

echo Running migrations...
echo.

echo [1/2] Adding price column to tables...
psql -U postgres -d authapp -f backend\migrations\add_price_to_tables.sql

if %errorlevel% == 0 (
    echo [OK] Price column migration completed!
) else (
    echo [WARNING] Price column migration may have failed or column already exists.
)

echo.
echo [2/2] Dropping name column from tables...
psql -U postgres -d authapp -f backend\migrations\drop_name_from_tables.sql

if %errorlevel% == 0 (
    echo [OK] Drop name column migration completed!
) else (
    echo [WARNING] Drop name column migration may have failed or column does not exist.
)

echo.
echo ========================================
echo   Migration Complete!
echo ========================================
echo.
pause
