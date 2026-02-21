@echo off
echo ========================================
echo   Starting Frontend (Vue 3 + Vite)
echo ========================================
echo.

cd /d "%~dp0frontend"

echo Installing packages...
npm install

echo.
echo App: http://localhost:5173
echo Press Ctrl+C to stop
echo.

npm run dev
