# 🏸 Badminton Court Reservation API

A RESTful API for managing badminton court reservations with integrated payment gateway (Midtrans).

## 🚀 Features

- ✅ Fetch available dates for booking
- ✅ Fetch available timeslots based on selected date
- ✅ Fetch available courts based on date & timeslot
- ✅ Create reservations with customer details
- ✅ Integrated payment gateway (Midtrans)
- ✅ Automatic reservation expiration (30 minutes)
- ✅ Background job to cleanup expired reservations
- ✅ PostgreSQL database (Neon / cloud Postgres)

## 🛠️ Tech Stack

- **Framework**: Beego v2
- **ORM**: Beego ORM
- **Database**: PostgreSQL (Neon / cloud Postgres)
- **Payment Gateway**: Midtrans Snap
- **Language**: Go 1.21+

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL database (Neon account)
- Midtrans account (for payment integration)

## 🔧 Installation

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd badminton-reservation-api
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup environment variables

Copy `.env.example` to `.env` and fill in your credentials:

```bash
cp .env.example .env
```

Edit `.env` file:

```env
# App Configuration
APP_NAME=badminton-reservation-api
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# Neon / PostgreSQL Configuration
# DB connection: use the connection details from your Neon project (Host, Port, User, Password, DB name)
DB_HOST=<your-neon-host>
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_database_password
DB_NAME=postgres
DB_SSLMODE=require

# Midtrans Configuration
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxx
MIDTRANS_IS_PRODUCTION=false

# Reservation Settings
RESERVATION_TIMEOUT_MINUTES=30
MAX_BOOKING_DAYS_AHEAD=30
```

### 4. Run database migrations

Apply the SQL migrations to your Neon/Postgres database. You can use the Neon Console SQL editor or run them with psql. Execute in order:

1. `database/migrations/001_create_courts.sql`
2. `database/migrations/002_create_timeslots.sql`
3. `database/migrations/003_create_reservations.sql`
4. `database/migrations/004_create_payments.sql`

### 5. Seed the database

Run the seed file:

```sql
database/seeds/seed_data.sql
```

### 6. Run the application

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Endpoints

### Health Check

```http
GET /health
```

### 1. Get Available Dates

Get list of available dates for the next 30 days.

```http
GET /api/v1/dates
```

**Response:**

```json
{
  "success": true,
  "message": "Available dates retrieved successfully",
  "data": [
    {
      "date": "2025-11-05",
      "day_name": "Wednesday",
      "is_weekend": false
    }
  ]
}
```

### 2. Get Available Timeslots

Get available timeslots for a specific date.

```http
GET /api/v1/timeslots?date=2025-11-10
```

**Response:**

```json
{
  "success": true,
  "message": "Available timeslots retrieved successfully",
  "data": [
    {
      "id": 1,
      "start_time": "08:00:00",
      "end_time": "09:00:00",
      "is_active": true
    }
  ]
}
```

### 3. Get Available Courts

Get available courts for specific date and timeslot.

```http
GET /api/v1/courts?date=2025-11-10&timeslot_id=3
```

**Response:**

```json
{
  "success": true,
  "message": "Available courts retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Court A",
      "description": "Premium indoor court",
      "price_per_hour": 100000,
      "status": "active"
    }
  ]
}
```

### 4. Create Reservation

Create a new court reservation.

```http
POST /api/v1/reservations
Content-Type: application/json

{
  "court_id": 1,
  "timeslot_id": 3,
  "booking_date": "2025-11-10",
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "customer_phone": "08123456789",
  "notes": "Optional notes"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Reservation created successfully",
  "data": {
    "id": "uuid-here",
    "court_id": 1,
    "timeslot_id": 3,
    "booking_date": "2025-11-10",
    "customer_name": "John Doe",
    "total_price": 100000,
    "status": "pending",
    "expired_at": "2025-11-05T10:30:00Z"
  }
}
```

### 5. Process Payment

Initiate payment for a reservation.

```http
POST /api/v1/payments/process
Content-Type: application/json

{
  "reservation_id": "uuid-here"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Payment transaction created successfully",
  "data": {
    "payment_id": "payment-uuid",
    "token": "midtrans-snap-token",
    "redirect_url": "https://app.sandbox.midtrans.com/snap/v2/..."
  }
}
```

### 6. Payment Callback (Webhook)

This endpoint is called by Midtrans after payment.

```http
POST /api/v1/payments/callback
Content-Type: application/json

{
  "order_id": "RES-xxxxx",
  "transaction_status": "settlement",
  "gross_amount": "100000",
  ...
}
```

### 7. Get Reservation Details

Get reservation by ID.

```http
GET /api/v1/reservations/:id
```

### 8. Get Customer Reservations

Get all reservations for a customer email.

```http
GET /api/v1/reservations/customer?email=john@example.com
```

## 💳 Payment Flow

1. User creates a reservation (status: `pending`)
2. User initiates payment via `/api/v1/payments/process`
3. System creates Midtrans Snap transaction
4. User redirected to Midtrans payment page
5. User completes payment
6. Midtrans sends notification to `/api/v1/payments/callback`
7. System updates reservation status to `paid`

## 🔄 Reservation Status Flow

```
pending → waiting_payment → paid
   ↓
expired (after 30 minutes)
```

## 🧪 Testing with Midtrans Sandbox

Use these test cards in Midtrans Sandbox:

**Success Payment:**

- Card Number: `4811 1111 1111 1114`
- CVV: `123`
- Exp: Any future date

**Failed Payment:**

- Card Number: `4911 1111 1111 1113`

## 📁 Project Structure

```
badminton-reservation-api/
├── conf/              # Beego configuration
├── controllers/       # HTTP request handlers
├── models/           # Database models
├── routers/          # API routes
├── services/         # Business logic & integrations
├── middleware/       # HTTP middlewares
├── utils/            # Utility functions
├── database/         # SQL migrations & seeds
├── main.go          # Application entry point
└── .env             # Environment variables
```

## 🐛 Troubleshooting

### Database connection error

Make sure your Neon/Postgres credentials are correct and the database is accessible.

### Payment not working

1. Check if Midtrans credentials are correct
2. Ensure callback URL is accessible (use ngrok for local testing)
3. Check Midtrans dashboard for transaction logs

### Reservation expires too fast

Adjust `RESERVATION_TIMEOUT_MINUTES` in `.env` file.

## 📝 License

MIT License

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first.

## 📞 Support

For questions or issues, please open an issue in the repository.
