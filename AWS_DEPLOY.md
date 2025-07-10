# Deploy to AWS Lambda Guide

## Prerequisites

1. **AWS CLI** - ติดตั้งและ configure AWS credentials
2. **AWS SAM CLI** - สำหรับ deploy serverless applications  
3. **MongoDB Atlas** - สร้าง cluster บน MongoDB Atlas
4. **Go 1.21+** - Go programming language

## Installation

### 1. Install AWS CLI
```powershell
# Windows (using Chocolatey)
choco install awscli

# หรือ download จาก: https://aws.amazon.com/cli/
```

### 2. Install AWS SAM CLI
```powershell
# Windows (using Chocolatey)
choco install aws-sam-cli

# หรือ download จาก: https://aws.amazon.com/serverless/sam/
```

### 3. Configure AWS Credentials
```powershell
aws configure
# AWS Access Key ID: YOUR_ACCESS_KEY
# AWS Secret Access Key: YOUR_SECRET_KEY  
# Default region name: us-east-1
# Default output format: json
```

## Setup Environment Variables

ใน AWS Systems Manager Parameter Store ตั้งค่า:

```powershell
# MongoDB Atlas connection string
aws ssm put-parameter --name "/agricultural/mongodb-uri" --value "mongodb+srv://username:password@cluster.mongodb.net/" --type "SecureString"

# JWT Secret Key
aws ssm put-parameter --name "/agricultural/jwt-secret" --value "your-super-secure-jwt-secret-key" --type "SecureString"

# Admin Password
aws ssm put-parameter --name "/agricultural/admin-password" --value "your-secure-admin-password" --type "SecureString"
```

## Deploy Steps

### Method 1: Using PowerShell Script (Windows)

```powershell
# Build และ deploy ครั้งแรก (guided setup)
.\deploy.ps1 deploy

# Deploy รอบถัดไป (fast)
.\deploy.ps1 deploy-fast

# Test locally
.\deploy.ps1 local

# Build only
.\deploy.ps1 build

# Clean artifacts
.\deploy.ps1 clean
```

### Method 2: Using Makefile (Linux/Mac)

```bash
# Build และ deploy ครั้งแรก
make deploy

# Deploy รอบถัดไป
make deploy-fast

# Test locally  
make local

# Build only
make build

# Clean artifacts
make clean
```

### Method 3: Manual Commands

```powershell
# 1. Build Lambda function
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o lambda ./cmd/lambda/main.go

# 2. Deploy with SAM
sam deploy --guided
```

## Configuration

### First Deploy (Guided)
SAM จะถามคำถาม:
- **Stack Name**: agricultural-api
- **AWS Region**: us-east-1 (หรือ region ที่ต้องการ)
- **Confirm changes**: Y
- **Allow SAM CLI IAM role creation**: Y
- **Save parameters**: Y

## API Endpoints

หลังจาก deploy สำเร็จ คุณจะได้ URL เช่น:
```
https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/Prod/
```

### Available Endpoints:
- **Health Check**: `GET /health`
- **Swagger Docs**: `GET /swagger/index.html`
- **Authentication**: `POST /api/auth/login`
- **Products**: `GET /api/products`
- **All APIs**: ตาม swagger documentation

## Environment Variables

Lambda function จะใช้ environment variables:
```env
MONGODB_URI=<from Parameter Store>
MONGODB_DATABASE=agricultural
JWT_SECRET=<from Parameter Store>
FRONTEND_URL=*
ADMIN_EMAIL=admin@agricultural.com
ADMIN_PASSWORD=<from Parameter Store>
```

## Local Testing

```powershell
# Start local API Gateway
.\deploy.ps1 local

# API จะทำงานที่ http://localhost:3000
```

## Monitoring

### CloudWatch Logs
- ดู logs ใน AWS CloudWatch
- Log group: `/aws/lambda/agricultural-api-AgriculturalAPI-xxxxx`

### API Gateway Logs
- ดู API access logs ใน CloudWatch

## Cost Optimization

### AWS Lambda Pricing:
- **Free Tier**: 1M requests/month + 400,000 GB-seconds
- **Pay per use**: ~$0.20 per 1M requests
- **Memory**: 512MB (configurable)
- **Timeout**: 30 seconds

### MongoDB Atlas:
- **Free Tier**: M0 cluster (512MB storage)
- **Shared Clusters**: Starting $9/month

## Troubleshooting

### 1. Build Issues
```powershell
# ตรวจสอบ Go version
go version

# Clean และ build ใหม่
.\deploy.ps1 clean
.\deploy.ps1 build
```

### 2. Deploy Issues  
```powershell
# Validate template
sam validate

# Check AWS credentials
aws sts get-caller-identity
```

### 3. Database Connection
- ตรวจสอบ MongoDB Atlas IP whitelist (เพิ่ม 0.0.0.0/0 สำหรับ Lambda)
- ตรวจสอบ connection string ใน Parameter Store

### 4. Environment Variables
```powershell
# ตรวจสอบ parameters
aws ssm get-parameter --name "/agricultural/mongodb-uri"
aws ssm get-parameter --name "/agricultural/jwt-secret"
```

## Security Notes

1. **Parameter Store**: ใช้ SecureString สำหรับ sensitive data
2. **IAM Roles**: Lambda จะสร้าง IAM role อัตโนมัติ
3. **API Gateway**: รองรับ HTTPS เท่านั้น
4. **CORS**: ตั้งค่าให้รองรับ cross-origin requests

## Updates

```powershell
# Update code และ deploy
git add .
git commit -m "Update Lambda function"
.\deploy.ps1 deploy-fast
```

## Cleanup

```powershell
# Delete CloudFormation stack
aws cloudformation delete-stack --stack-name agricultural-api

# Delete parameters
aws ssm delete-parameter --name "/agricultural/mongodb-uri"
aws ssm delete-parameter --name "/agricultural/jwt-secret" 
aws ssm delete-parameter --name "/agricultural/admin-password"
```
