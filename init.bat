@echo off
echo ===================================
echo   CareerForge Studio Initializer
echo ===================================
echo.

:: Check Node.js
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [INFO] Node.js is not installed. Attempting to auto-install via winget...
    winget install -e --id OpenJS.NodeJS --accept-package-agreements --accept-source-agreements
    echo.
    echo [SUCCESS] Node.js installed! PLEASE RESTART YOUR TERMINAL AND RUN THIS SCRIPT AGAIN.
    pause
    exit /b
)

:: Check Go
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [INFO] Go is not installed. Attempting to auto-install via winget...
    winget install -e --id GoLang.Go --accept-package-agreements --accept-source-agreements
    echo.
    echo [SUCCESS] Go installed! PLEASE RESTART YOUR TERMINAL AND RUN THIS SCRIPT AGAIN.
    pause
    exit /b
)

:: Check Antigravity CLI
where agy >nul 2>&1
if %errorlevel% neq 0 (
    echo [INFO] Antigravity CLI not found. Installing globally via npm...
    npm install -g @google/antigravity-cli
)

echo [INFO] Installing project dependencies...
call npm install

echo [INFO] Building Go Dashboard...
cd dashboard
go build .
cd ..

echo.
echo Setup complete! Launching CareerForge Studio...
call launch.bat
