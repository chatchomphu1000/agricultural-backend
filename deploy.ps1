# AWS Deploy Script for Windows PowerShell

param(
    [string]$Action = "help"
)

function Build-Lambda {
    Write-Host "Building Lambda function..." -ForegroundColor Green
    $env:GOOS = "linux"
    $env:GOARCH = "amd64" 
    $env:CGO_ENABLED = "0"
    go build -o lambda ./cmd/lambda/main.go
    Write-Host "Build completed!" -ForegroundColor Green
}

function Clean-Build {
    Write-Host "Cleaning up..." -ForegroundColor Yellow
    if (Test-Path "lambda") {
        Remove-Item "lambda"
    }
    if (Test-Path "lambda.exe") {
        Remove-Item "lambda.exe" 
    }
    Write-Host "Cleanup completed!" -ForegroundColor Green
}

function Deploy-AWS {
    Write-Host "Deploying to AWS..." -ForegroundColor Green
    Build-Lambda
    sam deploy --guided
}

function Deploy-Fast {
    Write-Host "Fast deploying to AWS..." -ForegroundColor Green
    Build-Lambda
    sam deploy
}

function Start-Local {
    Write-Host "Starting local API..." -ForegroundColor Green
    Build-Lambda
    sam local start-api
}

function Validate-Template {
    Write-Host "Validating SAM template..." -ForegroundColor Green
    sam validate
}

switch ($Action) {
    "build" { Build-Lambda }
    "clean" { Clean-Build }
    "deploy" { Deploy-AWS }
    "deploy-fast" { Deploy-Fast }
    "local" { Start-Local }
    "validate" { Validate-Template }
    default {
        Write-Host "AWS Deploy Script" -ForegroundColor Cyan
        Write-Host "Usage: .\deploy.ps1 [action]" -ForegroundColor White
        Write-Host ""
        Write-Host "Actions:" -ForegroundColor Yellow
        Write-Host "  build        - Build Lambda function"
        Write-Host "  clean        - Clean build artifacts"
        Write-Host "  deploy       - Deploy to AWS (guided)"
        Write-Host "  deploy-fast  - Fast deploy (no prompts)"
        Write-Host "  local        - Start local API server"
        Write-Host "  validate     - Validate SAM template"
    }
}
