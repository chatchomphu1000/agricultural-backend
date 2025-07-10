# 🚀 AWS Deploy Guide (ไม่ต้องใช้ SAM CLI)

## ✅ Prerequisites ที่ต้องมี:
1. **AWS CLI** - มีแล้ว ✅
2. **AWS Credentials** - ต้อง configure
3. **MongoDB Atlas** - สร้าง cluster (Free tier ได้)

## 🔧 ขั้นตอนการ Deploy:

### 1. Configure AWS Credentials
```powershell
# ใส่ AWS credentials (ต้องมี Access Key และ Secret Key)
"C:\Program Files\Amazon\AWSCLIV2\aws.exe" configure

# AWS Access Key ID: YOUR_ACCESS_KEY
# AWS Secret Access Key: YOUR_SECRET_KEY  
# Default region name: us-east-1
# Default output format: json
```

### 2. Deploy ขึ้น AWS (วิธีง่าย)
```powershell
# Deploy แบบง่าย (ใช้ค่า default)
.\aws-deploy.ps1 deploy

# หรือ Deploy พร้อม MongoDB และ JWT Secret
.\aws-deploy.ps1 deploy -MongoURI "mongodb+srv://username:password@cluster.mongodb.net/" -JWTSecret "your-secret-key" -AdminPassword "admin123"
```

### 3. ดูผลลัพธ์ หลัง Deploy เสร็จ
```powershell
# ดู API URL และ outputs
.\aws-deploy.ps1 outputs
```

## 📋 Commands ที่ใช้ได้:

```powershell
# Deploy ครั้งแรก
.\aws-deploy.ps1 deploy

# ดูสถานะ stack
.\aws-deploy.ps1 status

# ดู API URL และ endpoints
.\aws-deploy.ps1 outputs

# Build Lambda เท่านั้น
.\aws-deploy.ps1 build

# ลบทุกอย่างใน AWS
.\aws-deploy.ps1 delete

# ทำความสะอาดไฟล์ local
.\aws-deploy.ps1 clean

# ดูคำสั่งทั้งหมด
.\aws-deploy.ps1 help
```

## 🌍 MongoDB Atlas Setup (ถ้ายังไม่มี):

1. **สร้างบัญชี**: https://cloud.mongodb.com/
2. **สร้าง Free Cluster** (M0 - Free tier)
3. **สร้าง Database User**
4. **เพิ่ม IP Whitelist**: `0.0.0.0/0` (สำหรับ Lambda)
5. **คัดลอก Connection String**

## 🎯 หลังจาก Deploy สำเร็จ:

คุณจะได้:
- **API URL**: `https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/prod`
- **Swagger Docs**: `https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/prod/swagger/index.html`
- **Health Check**: `https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/prod/health`

## 💰 ค่าใช้จ่าย AWS:
- **Lambda**: Free tier 1M requests/month
- **API Gateway**: Free tier 1M requests/month  
- **CloudWatch Logs**: Free tier 5GB/month
- **รวม**: ~$0 สำหรับการใช้งานทดสอบ

## ❗ หากมีปัญหา:

```powershell
# ตรวจสอบ AWS credentials
"C:\Program Files\Amazon\AWSCLIV2\aws.exe" sts get-caller-identity

# ดู logs หาก deploy ล้มเหลว
"C:\Program Files\Amazon\AWSCLIV2\aws.exe" cloudformation describe-stack-events --stack-name agricultural-api

# ลบทุกอย่างและเริ่มใหม่
.\aws-deploy.ps1 delete
.\aws-deploy.ps1 clean
```
