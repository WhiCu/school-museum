# Запуск проекта школьного музея (Docker)

Write-Host '==================================' -ForegroundColor Cyan
Write-Host '  School Museum - Start' -ForegroundColor Cyan
Write-Host '==================================' -ForegroundColor Cyan

Write-Host ''
Write-Host 'Starting all services...' -ForegroundColor Yellow
docker compose up -d --build

Write-Host ''
Write-Host 'Waiting for backend (may take ~60s on first start)...' -ForegroundColor Yellow
for ($i = 1; $i -le 90; $i++) {
    try {
        $null = Invoke-WebRequest -Uri 'http://localhost:8081/ping?message=ping' -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
        Write-Host 'Backend ready!' -ForegroundColor Green
        break
    } catch {
        if ($i -eq 90) {
            Write-Host 'Backend did not respond in 90s. Check: docker compose logs server' -ForegroundColor Red
        }
        Start-Sleep -Seconds 1
    }
}

Write-Host ''
Write-Host '==================================' -ForegroundColor Cyan
Write-Host '  All started!' -ForegroundColor Green
Write-Host '==================================' -ForegroundColor Cyan
Write-Host ''
Write-Host '  Frontend:  http://localhost' -ForegroundColor White
Write-Host '  API:       http://localhost:8081' -ForegroundColor White
Write-Host ''
Write-Host '  To stop: run stop.ps1' -ForegroundColor Yellow
Write-Host '==================================' -ForegroundColor Cyan