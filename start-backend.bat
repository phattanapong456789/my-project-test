@echo off
echo ========================================
echo   Starting Backend (Go + Gin + JWT)
echo ========================================
echo.

cd /d "%~dp0backend"

echo Downloading dependencies...
go mod tidy

echo.
echo Server: http://localhost:8080
echo Press Ctrl+C to stop
echo.

go run main.go
