# Остановка проекта школьного музея (Docker)

Write-Host 'Stopping...' -ForegroundColor Yellow
docker compose down
Write-Host 'All services stopped' -ForegroundColor Green
Write-Host 'Done!' -ForegroundColor Cyan