# 🏸 API Reservasi Lapangan Badminton

Ini adalah backend RESTful API yang dibangun menggunakan framework **Beego (Golang)** untuk mengelola sistem reservasi lapangan badminton. API ini mencakup pengelolaan lapangan, slot waktu, pembuatan reservasi, dan integrasi gateway pembayaran dengan **Midtrans**.

Proyek ini berfungsi sebagai backend untuk aplikasi frontend Next.js berikut:

➡️ **Frontend Repo:** [https://github.com/SulthanRaghib/badminton-reservation](https://www.google.com/search?q=https://github.com/SulthanRaghib/badminton-reservation)

---

## ✨ Fitur Utama

Ringkasan fitur utama yang disediakan oleh API ini:

- **🏸 Manajemen Lapangan & Slot Waktu**

  - GET `/api/v1/dates` — daftar tanggal yang tersedia untuk pemesanan
  - GET `/api/v1/courts/all` — daftar semua lapangan
  - GET `/api/v1/timeslots` — daftar slot waktu dan ketersediaannya (query: `booking_date`, `court_id`)

- **🧾 Manajemen Reservasi**

  - POST `/api/v1/reservations` — buat reservasi baru
  - GET `/api/v1/reservations/:id` — ambil detail reservasi berdasarkan ID
  - GET `/api/v1/reservations/customer` — cari riwayat reservasi berdasarkan email (query: `email`)

- **💳 Integrasi Pembayaran (Midtrans)**

  - POST `/api/v1/payments/process` — inisiasi transaksi pembayaran untuk reservasi
  - POST `/api/v1/payments/callback` — webhook callback dari Midtrans untuk memperbarui status pembayaran

- **⌛ Reservasi Berbatas Waktu**

  - Reservasi awal berstatus `pending` dan akan otomatis `expired` jika tidak dibayar dalam 30 menit (nilai dapat diubah lewat `.env`)

- **🔄 Ketersediaan Slot Dinamis**

  - Saat reservasi dibuat, slot untuk `court_id` + `timeslot_id` pada `booking_date` akan ditandai `unavailable`.
  - Jika reservasi `expired` atau `cancelled`, slot akan dikembalikan menjadi `available`.

- **🐳 Dukungan Docker & Otomatisasi**

  - `Dockerfile` + `docker-entrypoint.sh` otomatis menjalankan migrasi dan seeding saat container dijalankan.
  - Tersedia target `make` dan beberapa CLI (dalam `cmd/`) untuk migrasi, seed, dan utilitas DB.

- **📖 Dokumentasi API**
  - Swagger UI tersedia di `/swagger` (self-hosted) untuk dokumentasi interaktif.

---

## 🛠️ Tumpukan Teknologi (Tech Stack)

- **Bahasa:** **Go** (v1.24+)
- **Framework:** **Beego** (v2.3.8)
- **Database:** **PostgreSQL** (Sangat direkomendasikan menggunakan [Neon](https://neon.tech/))
- **ORM:** **GORM** (untuk migrasi) & **Beego ORM** (untuk _query_ model)
- **Gateway Pembayaran:** **Midtrans**
- **Dokumentasi:** **Swagger** (via `swaggo`)
- **Konfigurasi:** `godotenv` untuk manajemen _environment variable_
- **Kontainerisasi:** **Docker**
- **Peralatan (Tooling):** **Make**

---

## 🚀 Instalasi & Menjalankan

### Prasyarat

1.  **Go** (versi 1.24 atau lebih baru).
2.  **PostgreSQL Database:** Direkomendasikan mendaftar di [Neon](https://neon.tech/) untuk mendapatkan _connection string_ Postgres gratis.
3.  **Akun Midtrans:** Diperlukan akun Sandbox Midtrans untuk mendapatkan _Server Key_.

### Langkah-langkah Menjalankan (Lokal)

1.  **Clone Repositori**

    ```bash
    git clone https://github.com/SulthanRaghib/badminton-reservation-api.git
    cd badminton-reservation-api
    ```

2.  **Instal Dependensi**

    ```bash
    go mod download
    ```

3.  **Konfigurasi Environment**
    Salin file `.env.example` menjadi `.env`.

    ```bash
    cp .env.example .env
    ```

    Buka file `.env` dan isi variabel berikut dengan kredensial Anda (terutama dari Neon dan Midtrans):

    ```env
    # Database (Neon / Postgres)
    DB_HOST=<your-neon-host>
    DB_PORT=5432
    DB_USER=<your-db-user>
    DB_PASSWORD=<your-db-password>
    DB_NAME=<your-db-name>
    DB_SSLMODE=require

    # Midtrans
    MIDTRANS_SERVER_KEY=<your-midtrans-server-key>
    MIDTRANS_CLIENT_KEY=<your-midtrans-client-key>
    MIDTRANS_IS_PRODUCTION=false
    ```

4.  **Migrasi & Seed Database**
    Proyek ini menggunakan `make` untuk mempermudah. Cukup jalankan:

    ```bash
    # Menjalankan migrasi GORM untuk membuat tabel
    make gorm-migrate
    ```

    ```bash
    # Mengisi data awal (lapangan & slot waktu)
    make db-seed-run
    ```

5.  **Jalankan Server API**

    ```bash
    go run main.go
    ```

    Server API akan berjalan di `http://localhost:8080` (atau port yang ditentukan di `.env`).

---

## 🐳 Menjalankan dengan Docker

Jika Anda memiliki Docker, Anda dapat menggunakan `Dockerfile` yang disediakan.

1.  **Pastikan file `.env` Anda sudah terisi** sesuai langkah instalasi di atas.

2.  **Build Docker Image:**

    ```bash
    make docker-build
    ```

    _(Atau: `docker build -t badminton-api:latest .`)_

3.  **Jalankan Docker Container:**

    ```bash
    docker run --rm -p 8080:8080 --env-file .env badminton-api:latest
    ```

    Skrip `docker-entrypoint.sh` akan secara otomatis menjalankan migrasi (`./migrate`) dan _seeding_ (`./seed`) sebelum memulai aplikasi utama (`./main`).

---

## 📁 Struktur Proyek

Berikut adalah gambaran umum struktur direktori utama:

```
badminton-reservation-api/
├── cmd/                # Aplikasi CLI pendukung
│   ├── dropdb/         # Skrip untuk membersihkan database
│   ├── migrate/        # Skrip migrasi GORM
│   └── seed/           # Skrip untuk seeding data
├── controllers/        # Handler HTTP (Logika Beego)
│   ├── reservation.go  # Logika untuk membuat & mengambil reservasi
│   ├── payment.go      # Logika untuk memproses pembayaran & callback
│   ├── court.go        # Logika untuk mengambil data lapangan
│   ├── timeslot.go     # Logika untuk mengambil data slot waktu
│   ├── date.go         # Logika untuk mengambil data tanggal
│   └── swagger_ui.go   # Handler untuk menyajikan Swagger UI
├── database/           # Skema & migrasi SQL
│   ├── migrations/     # File .sql untuk struktur tabel
│   └── seeds/          # File .sql untuk data awal
├── docs/               # File dokumentasi Swagger (generated)
├── middleware/         # Middleware HTTP
│   └── cors.go         # Konfigurasi CORS
├── models/             # Model data (structs) dan query ORM
│   ├── reservation.go  # Model & logika database Reservasi
│   ├── payment.go      # Model & logika database Pembayaran
│   └── ...
├── routers/            # Definisi rute API
│   └── route.go        # Mendaftarkan semua endpoint controller
├── services/           # Logika bisnis eksternal
│   └── payment/        # Logika integrasi Midtrans
├── utils/              # Fungsi helper
│   ├── database.go     # Koneksi DB
│   ├── response.go     # Standar respon JSON
│   └── validator.go    # Validasi email, telepon, dll.
├── .env.example        # Template konfigurasi
├── Dockerfile          #
├── go.mod              # Dependensi Go
├── main.go             # Entrypoint aplikasi
└── Makefile            # Skrip helper
```

---

## 📚 Dokumentasi API (Endpoints)

Untuk dokumentasi interaktif, jalankan server dan kunjungi:
**➡️ `http://localhost:8080/swagger`**

Endpoint utama yang tersedia di bawah `/api/v1`:

| Method | Endpoint                        | Deskripsi                                                                       |
| :----- | :------------------------------ | :------------------------------------------------------------------------------ |
| `GET`  | `/health`                       | Memeriksa status kesehatan API.                                                 |
| `GET`  | `/api/v1/dates`                 | Mendapatkan daftar tanggal yang tersedia untuk pemesanan.                       |
| `GET`  | `/api/v1/courts/all`            | Mendapatkan daftar semua lapangan yang terdaftar (aktif).                       |
| `GET`  | `/api/v1/courts`                | Mendapatkan lapangan yang _tersedia_ (Query: `booking_date`, `timeslot_id`).    |
| `GET`  | `/api/v1/timeslots/all`         | Mendapatkan daftar semua slot waktu global.                                     |
| `GET`  | `/api/v1/timeslots`             | Mendapatkan slot waktu & ketersediaannya (Query: `booking_date`, `court_id`).   |
| `POST` | `/api/v1/reservations`          | Membuat reservasi baru.                                                         |
| `GET`  | `/api/v1/reservations/:id`      | Mengambil detail reservasi berdasarkan ID-nya.                                  |
| `GET`  | `/api/v1/reservations/customer` | Mencari semua reservasi berdasarkan email (Query: `email`).                     |
| `POST` | `/api/v1/payments/process`      | Memulai proses pembayaran untuk reservasi (Body: `reservation_id`).             |
| `GET`  | `/api/v1/payments/:id`          | Mendapatkan status pembayaran (ID bisa berupa ID Reservasi atau ID Pembayaran). |
| `POST` | `/api/v1/payments/callback`     | **[WEBHOOK]** Endpoint internal untuk menerima notifikasi dari Midtrans.        |
