# Zapusk bekenda i frontenda

Write-Host '==================================' -ForegroundColor Cyan
Write-Host '  School Museum - Start' -ForegroundColor Cyan
Write-Host '==================================' -ForegroundColor Cyan

Write-Host ''
Write-Host '[1/2] Starting backend...' -ForegroundColor Yellow
docker-compose up -d --build

Write-Host ''
Write-Host 'Waiting for server (may take ~60s on first start)...' -ForegroundColor Yellow
for ($i = 1; $i -le 90; $i++) {
    try {
        $r = Invoke-WebRequest -Uri 'http://localhost:8081/museum/exhibitions' -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
        Write-Host 'Backend ready!' -ForegroundColor Green
        break
    } catch {
        if ($i -eq 90) {
            Write-Host 'Backend did not respond in 90s. Check: docker-compose logs server' -ForegroundColor Red
        }
        Start-Sleep -Seconds 1
    }
}

Write-Host ''
Write-Host '[2/2] Starting frontend on port 5500...' -ForegroundColor Yellow
$sd = $PSScriptRoot
$frontendJob = Start-Job -ScriptBlock {
    Set-Location (Join-Path $using:sd 'frontend')
    py -m http.server 5500
}

$frontendJob.Id | Out-File -FilePath '.frontend.jobid'

Write-Host ''
Write-Host '==================================' -ForegroundColor Cyan
Write-Host '  All started!' -ForegroundColor Green
Write-Host '==================================' -ForegroundColor Cyan
Write-Host ''
Write-Host '  Frontend:  http://localhost:5500' -ForegroundColor White
Write-Host '  API:       http://localhost:8081' -ForegroundColor White
Write-Host '  Jaeger:    http://localhost:16686' -ForegroundColor White
Write-Host '  Umami:     http://localhost:3000' -ForegroundColor White
Write-Host ''
Write-Host '  To stop: run stop.ps1' -ForegroundColor Yellow
Write-Host '==================================' -ForegroundColor Cyan