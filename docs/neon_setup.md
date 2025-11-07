# Neon (cloud Postgres) setup

This document explains how to connect this app to a Neon (Postgres) database and apply migrations.

1. Create a Neon project

   - Go to https://neon.tech/ and sign in or sign up.
   - Create a new project and note the connection details (host, port, user, password, db name).

2. Set environment variables

   - Copy `.env.example` to `.env` and fill the `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` with the Neon credentials.

   Example (Windows CMD):

```bat
set DB_HOST=<your-neon-host>
set DB_PORT=5432
set DB_USER=<your-db-user>
set DB_PASSWORD=<your-db-password>
set DB_NAME=<your-db-name>
set DB_SSLMODE=require
```

3. Test connection with psql (optional)
   - If you have `psql` installed, run:

```bat
psql "postgres://%DB_USER%:%DB_PASSWORD%@%DB_HOST%:%DB_PORT%/%DB_NAME%?sslmode=require" -c "SELECT 1;"
```

Replace environment variables or use the literal connection string from Neon Console.

4. Apply SQL migrations

   - Option A (recommended for repeatable deploys): run SQL files in `database/migrations/` using the Neon Console SQL editor or psql in order:

     1. `001_create_courts.sql`
     2. `002_create_timeslots.sql`
     3. `003_create_reservations.sql`
     4. `004_create_payments.sql`

   - Option B (convenience): run GORM AutoMigrate (will create/alter tables automatically):

```bat
go run ./cmd/migrate
```

This will use the `.env` values and run `RunGormMigrations()`.

5. Seed data (optional)

   - If `database/seeds/seed_data.sql` exists, run it in the Neon Console or with psql to insert sample courts and timeslots.

6. Run the application

```bat
go run main.go
```

Or build and run:

```bat
go build -o badminton-api.exe
badminton-api.exe
```

Notes

- Neon is Postgres-compatible; the app uses beego/orm and GORM for migrations. Use whichever migration strategy fits your deployment workflow.
- Ensure `DB_SSLMODE=require` when connecting to Neon from outside its private network.
