@echo off
echo ========================================
echo   Setup PostgreSQL Database
echo ========================================
echo.

set PGPASSWORD=password

echo Creating database "authapp"...
psql -U postgres -c "CREATE DATABASE authapp;" 2>nul

if %errorlevel% == 0 (
    echo [OK] Database "authapp" created!
) else (
    echo [INFO] Database may already exist, that's OK.
)

echo.
echo ========================================
echo   Done! Connection Info:
echo   Host    : localhost
echo   Port    : 5432
echo   DB      : authapp
echo   User    : postgres
echo   Password: password
echo ========================================
echo.
echo If password is different, edit: backend\.env
echo.
pause
