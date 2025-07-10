# AWS CloudFormation Direct Deployment Script

param(
    [string]$Action = "help",
    [string]$StackName = "agricultural-api",
    [string]$Region = "us-east-1",
    [string]$MongoURI = "",
    [string]$JWTSecret = "",
    [string]$AdminPassword = ""
)

$ErrorActionPreference = "Stop"

function Test-AWSCreds {
    try {
        $result = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" sts get-caller-identity 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ AWS credentials configured" -ForegroundColor Green
            return $true
        }
    }
    catch {
        Write-Host "❌ AWS credentials not configured" -ForegroundColor Red
        return $false
    }
    return $false
}
}

function Build-Lambda {
    Write-Host "🔨 Building Lambda function..." -ForegroundColor Yellow
    
    # Build for Linux
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"
    
    # Build the binary
    go build -o lambda ./cmd/lambda/main.go
    
    if (-not (Test-Path "lambda")) {
        throw "Failed to build lambda binary"
    }
    
    # Create ZIP file
    if (Test-Path "lambda.zip") {
        Remove-Item "lambda.zip"
    }
    
    # Use PowerShell's Compress-Archive
    Compress-Archive -Path "lambda" -DestinationPath "lambda.zip"
    
    Write-Host "✅ Lambda function built and packaged" -ForegroundColor Green
}

function Deploy-Stack {
    Write-Host "🚀 Deploying CloudFormation stack..." -ForegroundColor Yellow
    
    # Check if stack exists
    $stackExists = $false
    try {
        $result = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stacks --stack-name $StackName --region $Region 2>$null
        if ($LASTEXITCODE -eq 0) {
            $stackExists = $true
            Write-Host "📦 Stack exists, updating..." -ForegroundColor Blue
        }
    }
    catch {
        Write-Host "📦 Creating new stack..." -ForegroundColor Blue
    }
    
    # Get or create S3 bucket name
    $bucketName = "agricultural-lambda-deploy-$(Get-Random)"
    
    # Create S3 bucket if not exists
    Write-Host "📁 Creating S3 bucket: $bucketName" -ForegroundColor Blue
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" s3 mb "s3://$bucketName" --region $Region
    
    # Upload lambda.zip to S3
    Write-Host "📤 Uploading Lambda package to S3..." -ForegroundColor Blue
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" s3 cp lambda.zip "s3://$bucketName/lambda.zip"
    
    # Prepare parameters
    $parameters = @()
    if ($MongoURI) { $parameters += "ParameterKey=MongoDBURI,ParameterValue=$MongoURI" }
    if ($JWTSecret) { $parameters += "ParameterKey=JWTSecret,ParameterValue=$JWTSecret" }
    if ($AdminPassword) { $parameters += "ParameterKey=AdminPassword,ParameterValue=$AdminPassword" }
    
    # Update template to use the created bucket
    $templateContent = Get-Content "cloudformation.yaml" -Raw
    $templateContent = $templateContent -replace "S3Bucket: !Ref LambdaDeploymentBucket", "S3Bucket: $bucketName"
    $templateContent | Set-Content "cloudformation-deploy.yaml"
    
    # Deploy stack
    if ($stackExists) {
        $cmd = "update-stack"
    } else {
        $cmd = "create-stack"
    }
    
    $deployArgs = @(
        "cloudformation", $cmd,
        "--stack-name", $StackName,
        "--template-body", "file://cloudformation-deploy.yaml",
        "--capabilities", "CAPABILITY_IAM",
        "--region", $Region
    )
    
    if ($parameters.Count -gt 0) {
        $deployArgs += "--parameters"
        $deployArgs += $parameters
    }
    
    Write-Host "🔄 Executing: aws $($deployArgs -join ' ')" -ForegroundColor Gray
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" @deployArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Stack deployment initiated" -ForegroundColor Green
        
        # Wait for stack completion
        Write-Host "⏳ Waiting for stack completion..." -ForegroundColor Yellow
        & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation wait "stack-$cmd-complete" --stack-name $StackName --region $Region
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "🎉 Stack deployment completed!" -ForegroundColor Green
            Get-StackOutputs
        } else {
            Write-Host "❌ Stack deployment failed" -ForegroundColor Red
        }
    } else {
        Write-Host "❌ Stack deployment failed" -ForegroundColor Red
    }
    
    # Cleanup
    if (Test-Path "cloudformation-deploy.yaml") {
        Remove-Item "cloudformation-deploy.yaml"
    }
}

function Get-StackOutputs {
    Write-Host "📋 Getting stack outputs..." -ForegroundColor Blue
    
    $outputs = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stacks --stack-name $StackName --region $Region --query "Stacks[0].Outputs" --output table
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "🌐 Stack Outputs:" -ForegroundColor Green
        Write-Host $outputs
        
        # Get API URL specifically
        $apiUrl = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stacks --stack-name $StackName --region $Region --query "Stacks[0].Outputs[?OutputKey=='ApiUrl'].OutputValue" --output text
        
        if ($apiUrl) {
            Write-Host ""
            Write-Host "🚀 API URL: $apiUrl" -ForegroundColor Cyan
            Write-Host "📖 Swagger: $apiUrl/swagger/index.html" -ForegroundColor Cyan
            Write-Host "💓 Health: $apiUrl/health" -ForegroundColor Cyan
        }
    }
}

function Delete-Stack {
    Write-Host "🗑️ Deleting CloudFormation stack..." -ForegroundColor Yellow
    
    # Get S3 bucket name before deletion
    try {
        $bucketName = & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stacks --stack-name $StackName --region $Region --query "Stacks[0].Outputs[?OutputKey=='S3Bucket'].OutputValue" --output text
        
        if ($bucketName) {
            Write-Host "🗑️ Emptying S3 bucket: $bucketName" -ForegroundColor Yellow
            & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" s3 rm "s3://$bucketName" --recursive
        }
    }
    catch {
        Write-Host "⚠️ Could not clean S3 bucket" -ForegroundColor Yellow
    }
    
    & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation delete-stack --stack-name $StackName --region $Region
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Stack deletion initiated" -ForegroundColor Green
        Write-Host "⏳ Waiting for deletion to complete..." -ForegroundColor Yellow
        & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation wait stack-delete-complete --stack-name $StackName --region $Region
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "🎉 Stack deleted successfully!" -ForegroundColor Green
        }
    }
}

function Clean-Artifacts {
    Write-Host "🧹 Cleaning build artifacts..." -ForegroundColor Yellow
    
    @("lambda", "lambda.exe", "lambda.zip", "cloudformation-deploy.yaml") | ForEach-Object {
        if (Test-Path $_) {
            Remove-Item $_
            Write-Host "🗑️ Removed $_" -ForegroundColor Gray
        }
    }
    
    Write-Host "✅ Cleanup completed" -ForegroundColor Green
}

function Show-Help {
    Write-Host "AWS CloudFormation Direct Deploy Script" -ForegroundColor Cyan
    Write-Host "=======================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\aws-deploy.ps1 [action] [options]" -ForegroundColor White
    Write-Host ""
    Write-Host "Actions:" -ForegroundColor Yellow
    Write-Host "  deploy       - Build และ deploy ขึ้น AWS"
    Write-Host "  status       - ดูสถานะ stack"
    Write-Host "  outputs      - ดู stack outputs (API URL)"
    Write-Host "  delete       - ลบ stack ทั้งหมด"
    Write-Host "  build        - Build Lambda function เท่านั้น"
    Write-Host "  clean        - ลบ build artifacts"
    Write-Host ""
    Write-Host "Options:" -ForegroundColor Yellow
    Write-Host "  -StackName   - ชื่อ CloudFormation stack (default: agricultural-api)"
    Write-Host "  -Region      - AWS region (default: us-east-1)"
    Write-Host "  -MongoURI    - MongoDB connection string"
    Write-Host "  -JWTSecret   - JWT secret key"
    Write-Host "  -AdminPassword - Admin user password"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Green
    Write-Host "  .\aws-deploy.ps1 deploy -MongoURI 'mongodb+srv://...' -JWTSecret 'secret123'"
    Write-Host "  .\aws-deploy.ps1 status"
    Write-Host "  .\aws-deploy.ps1 outputs"
    Write-Host "  .\aws-deploy.ps1 delete"
}

# Main script logic
try {
    switch ($Action.ToLower()) {
        "deploy" {
            if (-not (Test-AWSCreds)) {
                Write-Host "Please configure AWS credentials first:" -ForegroundColor Red
                Write-Host "aws configure" -ForegroundColor Yellow
                exit 1
            }
            Build-Lambda
            Deploy-Stack
        }
        "status" {
            if (-not (Test-AWSCreds)) { exit 1 }
            & "C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stacks --stack-name $StackName --region $Region --query "Stacks[0].StackStatus" --output text
        }
        "outputs" {
            if (-not (Test-AWSCreds)) { exit 1 }
            Get-StackOutputs
        }
        "delete" {
            if (-not (Test-AWSCreds)) { exit 1 }
            Delete-Stack
        }
        "build" {
            Build-Lambda
        }
        "clean" {
            Clean-Artifacts
        }
        default {
            Show-Help
        }
    }
}
catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
