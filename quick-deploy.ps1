param(
    [string]$Action = "help"
)

function Build-Lambda {
    Write-Host "Building Lambda function..." -ForegroundColor Yellow
    
    $env:GOOS = "linux"
    $env:GOARCH = "amd64" 
    $env:CGO_ENABLED = "0"
    
    go build -o lambda ./cmd/lambda/main.go
    
    if (Test-Path "lambda") {
        Write-Host "Lambda built successfully" -ForegroundColor Green
        
        if (Test-Path "lambda.zip") {
            Remove-Item "lambda.zip"
        }
        Compress-Archive -Path "lambda" -DestinationPath "lambda.zip"
        Write-Host "Lambda packaged as lambda.zip" -ForegroundColor Green
    } else {
        Write-Host "Build failed" -ForegroundColor Red
        exit 1
    }
}

function Test-AWS {
    try {
        $result = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" sts get-caller-identity
        if ($LASTEXITCODE -eq 0) {
            Write-Host "AWS credentials OK" -ForegroundColor Green
            return $true
        }
    }
    catch {
        Write-Host "AWS not configured. Run: aws configure" -ForegroundColor Red
        return $false
    }
    return $false
}

function Deploy-ToAWS {
    Write-Host "Deploying to AWS..." -ForegroundColor Yellow
    
    if (-not (Test-AWS)) {
        exit 1
    }
    
    Build-Lambda
    
    # Create S3 bucket
    $bucketName = "agricultural-deploy-$(Get-Random)"
    Write-Host "Creating S3 bucket: $bucketName" -ForegroundColor Blue
    
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" s3 mb "s3://$bucketName"
    
    # Upload lambda
    Write-Host "Uploading lambda.zip..." -ForegroundColor Blue
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" s3 cp lambda.zip "s3://$bucketName/lambda.zip"
    
    # Deploy stack
    Write-Host "Deploying CloudFormation..." -ForegroundColor Blue
    $stackName = "agricultural-api"
    
    # Simple deploy command
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation create-stack --stack-name $stackName --template-body file://cloudformation.yaml --capabilities CAPABILITY_IAM --parameters ParameterKey=MongoDBURI,ParameterValue=mongodb://localhost:27017 ParameterKey=JWTSecret,ParameterValue=jwt-secret ParameterKey=AdminPassword,ParameterValue=admin123
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Deployment initiated. Waiting..." -ForegroundColor Yellow
        & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation wait stack-create-complete --stack-name $stackName
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Deployment successful!" -ForegroundColor Green
            Get-Outputs
        }
    }
}

function Get-Outputs {
    $stackName = "agricultural-api"
    Write-Host "Getting stack outputs..." -ForegroundColor Blue
    
    $apiUrl = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stacks --stack-name $stackName --query "Stacks[0].Outputs[?OutputKey=='ApiUrl'].OutputValue" --output text
    
    if ($apiUrl) {
        Write-Host ""
        Write-Host "API URL: $apiUrl" -ForegroundColor Cyan
        Write-Host "Swagger: $apiUrl/swagger/index.html" -ForegroundColor Cyan
        Write-Host "Health: $apiUrl/health" -ForegroundColor Cyan
    }
}

function Clean-Files {
    Write-Host "Cleaning artifacts..." -ForegroundColor Yellow
    
    $files = @("lambda", "lambda.exe", "lambda.zip")
    foreach ($file in $files) {
        if (Test-Path $file) {
            Remove-Item $file
            Write-Host "Removed $file" -ForegroundColor Gray
        }
    }
    Write-Host "Cleanup complete" -ForegroundColor Green
}

function Show-Help {
    Write-Host "AWS Deploy Script" -ForegroundColor Cyan
    Write-Host "=================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Yellow
    Write-Host "  .\quick-deploy.ps1 build     - Build Lambda function"
    Write-Host "  .\quick-deploy.ps1 deploy    - Deploy to AWS"
    Write-Host "  .\quick-deploy.ps1 outputs   - Show API URLs"
    Write-Host "  .\quick-deploy.ps1 clean     - Clean artifacts"
    Write-Host ""
    Write-Host "Before deploy, run: aws configure" -ForegroundColor Green
}

# Main
switch ($Action.ToLower()) {
    "build" { Build-Lambda }
    "deploy" { Deploy-ToAWS }
    "outputs" { Get-Outputs }
    "clean" { Clean-Files }
    default { Show-Help }
}
