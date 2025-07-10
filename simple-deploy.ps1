# Simple AWS Deploy Script
param(
    [string]$Action = "help",
    [string]$MongoURI = "mongodb://localhost:27017",
    [string]$JWTSecret = "your-jwt-secret",
    [string]$AdminPassword = "password123"
)

function Show-Help {
    Write-Host "🚀 AWS Deploy Script" -ForegroundColor Cyan
    Write-Host "===================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Yellow
    Write-Host "  .\simple-deploy.ps1 build     - Build Lambda function"
    Write-Host "  .\simple-deploy.ps1 deploy    - Deploy to AWS"
    Write-Host "  .\simple-deploy.ps1 clean     - Clean artifacts"
    Write-Host ""
    Write-Host "Example:" -ForegroundColor Green
    Write-Host '  .\simple-deploy.ps1 deploy'
}

function Build-Lambda {
    Write-Host "🔨 Building Lambda function..." -ForegroundColor Yellow
    
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"
    
    go build -o lambda ./cmd/lambda/main.go
    
    if (Test-Path "lambda") {
        Write-Host "✅ Lambda built successfully" -ForegroundColor Green
        
        # Create ZIP
        if (Test-Path "lambda.zip") {
            Remove-Item "lambda.zip"
        }
        Compress-Archive -Path "lambda" -DestinationPath "lambda.zip"
        Write-Host "✅ Lambda packaged as lambda.zip" -ForegroundColor Green
    } else {
        Write-Host "❌ Build failed" -ForegroundColor Red
        exit 1
    }
}

function Deploy-ToAWS {
    Write-Host "🚀 Deploying to AWS..." -ForegroundColor Yellow
    
    # Check AWS CLI
    try {
        $result = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" sts get-caller-identity
        Write-Host "✅ AWS credentials OK" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ AWS not configured. Run: aws configure" -ForegroundColor Red
        exit 1
    }
    
    # Build first
    Build-Lambda
    
    # Create unique bucket name
    $bucketName = "agricultural-deploy-$(Get-Random)"
    Write-Host "📁 Creating S3 bucket: $bucketName" -ForegroundColor Blue
    
    # Create bucket
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" s3 mb "s3://$bucketName"
    
    # Upload lambda
    Write-Host "📤 Uploading lambda.zip..." -ForegroundColor Blue
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" s3 cp lambda.zip "s3://$bucketName/lambda.zip"
    
    # Update CloudFormation template
    $template = Get-Content "cloudformation.yaml" -Raw
    $template = $template -replace "S3Bucket: !Ref LambdaDeploymentBucket", "S3Bucket: $bucketName"
    $template | Set-Content "deploy-template.yaml"
    
    # Deploy CloudFormation
    Write-Host "☁️ Deploying CloudFormation..." -ForegroundColor Blue
    $stackName = "agricultural-api"
    
    $params = @(
        "cloudformation", "create-stack",
        "--stack-name", $stackName,
        "--template-body", "file://deploy-template.yaml",
        "--capabilities", "CAPABILITY_IAM",
        "--parameters",
        "ParameterKey=MongoDBURI,ParameterValue=$MongoURI",
        "ParameterKey=JWTSecret,ParameterValue=$JWTSecret", 
        "ParameterKey=AdminPassword,ParameterValue=$AdminPassword"
    )
    
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" @params
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "⏳ Waiting for deployment..." -ForegroundColor Yellow
        & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation wait stack-create-complete --stack-name $stackName
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "🎉 Deployment successful!" -ForegroundColor Green
            
            # Get API URL
            $apiUrl = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stacks --stack-name $stackName --query "Stacks[0].Outputs[?OutputKey=='ApiUrl'].OutputValue" --output text
            
            Write-Host ""
            Write-Host "🌐 API URL: $apiUrl" -ForegroundColor Cyan
            Write-Host "📖 Swagger: $apiUrl/swagger/index.html" -ForegroundColor Cyan
            Write-Host "💓 Health: $apiUrl/health" -ForegroundColor Cyan
        } else {
            Write-Host "❌ Deployment failed" -ForegroundColor Red
        }
    } else {
        Write-Host "❌ CloudFormation failed" -ForegroundColor Red
    }
    
    # Cleanup
    if (Test-Path "deploy-template.yaml") {
        Remove-Item "deploy-template.yaml"
    }
}

function Clean-Files {
    Write-Host "🧹 Cleaning artifacts..." -ForegroundColor Yellow
    
    $files = @("lambda", "lambda.exe", "lambda.zip", "deploy-template.yaml")
    foreach ($file in $files) {
        if (Test-Path $file) {
            Remove-Item $file
            Write-Host "Removed $file" -ForegroundColor Gray
        }
    }
    Write-Host "✅ Cleanup complete" -ForegroundColor Green
}

# Main logic
switch ($Action.ToLower()) {
    "build" { Build-Lambda }
    "deploy" { Deploy-ToAWS }
    "clean" { Clean-Files }
    default { Show-Help }
}
