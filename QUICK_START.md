# 🚀 Quick Start Guide

Panduan cepat untuk menjalankan Badminton Reservation API dalam 10 menit!

## 📋 Prerequisites

- ✅ Go 1.21+ installed
- ✅ Neon (cloud Postgres) account
- ✅ Midtrans Sandbox account (free)

## 🎯 Step-by-Step Setup

### Step 1: Create Neon Project (5 menit)

1. **Buka** https://neon.tech/ dan login or sign up
2. **Create** a new project (choose a region)
3. **Catat** credentials from the Neon Console (Connection info):

- Host: `<your-neon-host>`
- Port: `5432`
- User: `<your-db-user>`
- Password: `<your-db-password>`
- Database name: `<your-db-name>`

### Step 2: Setup Database (3 menit)

1. **Buka** Neon Console → SQL Editor (or connect with psql)
2. **Copy-paste dan Execute** migrations secara berurutan:

**Migration 1 - Courts Table:**

```sql
-- Paste isi dari database/migrations/001_create_courts.sql
```

**Migration 2 - Timeslots Table:**

```sql
-- Paste isi dari database/migrations/002_create_timeslots.sql
```

**Migration 3 - Reservations Table:**

```sql
-- Paste isi dari database/migrations/003_create_reservations.sql
```

**Migration 4 - Payments Table:**

```sql
-- Paste isi dari database/migrations/004_create_payments.sql
```

**Seed Data:**

```sql
-- Paste isi dari database/seeds/seed_data.sql
```

3. **Verify** data dengan query:

```sql
SELECT COUNT(*) FROM courts; -- Should return 5
SELECT COUNT(*) FROM timeslots; -- Should return 14
```

### Step 3: Setup Midtrans Account (2 menit)

1. **Daftar** di https://dashboard.sandbox.midtrans.com
2. **Login** ke Dashboard
3. **Copy** credentials:
   - Settings → Access Keys
   - Server Key: `SB-Mid-server-xxxxx`
   - Client Key: `SB-Mid-client-xxxxx`

### Step 4: Clone & Setup Project (2 menit)

```bash
# Clone project
git clone <your-repo>
cd badminton-reservation-api

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env
```

### Step 5: Configure Environment Variables (1 menit)

Edit file `.env` and fill Neon connection details (or use environment variables):

```env
# App Config
APP_NAME=badminton-reservation-api
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Database (Neon / Postgres)
DB_HOST=<your-neon-host>
DB_PORT=5432
DB_USER=<your-db-user>
DB_PASSWORD=<your-db-password>
DB_NAME=<your-db-name>
DB_SSLMODE=require

# Midtrans
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxx
MIDTRANS_IS_PRODUCTION=false

# Settings
RESERVATION_TIMEOUT_MINUTES=30
MAX_BOOKING_DAYS_AHEAD=30
```

### Step 6: Run the API! 🎉

```bash
go run main.go
```

Output yang diharapkan:

```
========================================
🏸 Badminton Reservation API
========================================
Environment: development
Port: 8080
API Base URL: http://localhost:8080
========================================
```

## 🧪 Testing API

### Test 1: Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "success": true,
  "message": "API is running",
  "data": {
    "status": "healthy",
    "database": "connected"
  }
}
```

### Test 2: Get Available Dates

```bash
curl http://localhost:8080/api/v1/dates
```

### Test 3: Create a Reservation

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -d '{
    "court_id": 1,
    "timeslot_id": 3,
    "booking_date": "2025-11-10",
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_phone": "08123456789"
  }'
```

Save the `id` from response for next step!

### Test 4: Process Payment

```bash
curl -X POST http://localhost:8080/api/v1/payments/process \
  -H "Content-Type: application/json" \
  -d '{
    "reservation_id": "YOUR_RESERVATION_ID_HERE"
  }'
```

You'll get a `redirect_url` - open it in browser to test payment with Midtrans!

## 🎨 Test Payment Cards (Midtrans Sandbox)

**Successful Payment:**

- Card: `4811 1111 1111 1114`
- CVV: `123`
- Expiry: Any future date
- OTP: `112233`

**Failed Payment:**

- Card: `4911 1111 1111 1113`

## 📱 Import Postman Collection

1. Open Postman
2. Import `Badminton_API.postman_collection.json`
3. Set `base_url` variable to `http://localhost:8080`
4. Test all endpoints!

## 🐛 Common Issues

### Issue: "Failed to connect to database"

**Solution:**

- Check DB_HOST, DB_PASSWORD in .env
- Make sure Neon project is active
- Verify DB_SSLMODE=require

### Issue: "Midtrans payment not working"

**Solution:**

- Check MIDTRANS_SERVER_KEY and MIDTRANS_CLIENT_KEY
- Make sure you're using SANDBOX keys (starts with SB-)
- For callback testing, use ngrok: `ngrok http 8080`

### Issue: "Reservation expired too fast"

**Solution:**

- Adjust RESERVATION_TIMEOUT_MINUTES in .env
- Default is 30 minutes

## 🔥 Next Steps

1. **Test Full Flow:**

   - Create reservation → Process payment → Complete payment → Check status

2. **Customize:**

- Add more courts in your database (Neon)
- Modify timeslots
- Adjust pricing

3. **Deploy:**
   - Deploy to Railway/Render/Heroku
   - Update Midtrans callback URL
   - Switch to production keys when ready

## 📞 Need Help?

- Check README.md for detailed documentation
- Review API endpoints in Postman collection
- Test with Postman first before integrating with frontend

## ✅ Checklist

- [ ] Neon project created
- [ ] Database migrations executed
- [ ] Seed data inserted
- [ ] Midtrans account created
- [ ] Environment variables configured
- [ ] API running on localhost:8080
- [ ] Health check returns "healthy"
- [ ] Can create reservation
- [ ] Can process payment
- [ ] Payment callback works

**Congratulations! 🎉 Your Badminton Reservation API is ready!**
