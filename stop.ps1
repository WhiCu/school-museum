# Stop backend and frontend

Write-Host 'Stopping...' -ForegroundColor Yellow

if (Test-Path '.frontend.jobid') {
    $jobId = Get-Content '.frontend.jobid'
    Stop-Job -Id $jobId -ErrorAction SilentlyContinue
    Remove-Job -Id $jobId -Force -ErrorAction SilentlyContinue
    Remove-Item '.frontend.jobid'
    Write-Host 'Frontend stopped' -ForegroundColor Green
}

docker-compose down
Write-Host 'Backend stopped' -ForegroundColor Green
Write-Host 'Done!' -ForegroundColor Cyan