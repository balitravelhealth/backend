# BaliTravelHealth Platform Backend

Backend, expert-system service, database migrations, and admin dashboard for a Bali travel health platform.

This repository contains a Go API gateway, a Python/FastAPI expert-system service, PostgreSQL schema migrations, Docker infrastructure, and a Next.js web admin dashboard.

> Medical disclaimer: this system is intended as an informational and early triage support tool only. Diagnosis output, recommendations, emergency guide content, and health knowledge base data must be reviewed and validated by qualified healthcare professionals before production use. This software does not replace doctors, nurses, licensed medical services, emergency services, or direct clinical evaluation.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Repository Structure](#repository-structure)
- [Requirements](#requirements)
- [Environment Configuration](#environment-configuration)
- [Running with Docker Compose](#running-with-docker-compose)
- [Running Locally](#running-locally)
- [Database Migrations](#database-migrations)
- [API Endpoints](#api-endpoints)
- [Expert System](#expert-system)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Security and Open Source Checklist](#security-and-open-source-checklist)
- [Open Source Dependency Licenses](#open-source-dependency-licenses)
- [Project License](#project-license)

## Features

- Google ID token authentication for traveler users.
- Email/password authentication for admin and nurse users.
- JWT access tokens with a 15-minute lifetime.
- Refresh tokens with a 30-day lifetime.
- Refresh token rotation and token reuse detection.
- Role-based access control for `traveler`, `nurse`, and `admin`.
- Traveler health profile management.
- Traveler profile and emergency contact management.
- Vaccination record management.
- Health assessment submission and history.
- Python expert-system integration for diagnosis support.
- Forward Chaining and Certainty Factor inference.
- Admin-managed knowledge base for symptoms, diseases, and expert rules.
- `pre_travel` and `post_travel` diagnosis categories.
- Bali location classification based on latitude and longitude.
- Destination, regional health risk, and nearby medical facility APIs.
- Step-based emergency guides.
- Decision-tree emergency guide flows.
- Nurse management and nursing care records.
- Web admin dashboard for assessments, facilities, destinations, emergency guides, nurses, and expert knowledge base data.
- Dockerfile for each service.
- Docker Compose stack for local development.

## Architecture

```text
Mobile Client / Web Admin
        |
        v
Go API Gateway (Gin) :8080
        |
        |-- PostgreSQL :5432
        |
        |-- Python Expert Service (FastAPI) :8001
                |
                |-- reads published expert rules from PostgreSQL
```

Diagnosis flow:

1. The client submits symptoms to `POST /assessment`.
2. The Go gateway loads the authenticated user's health profile from PostgreSQL.
3. The gateway forwards symptoms, diagnosis category, and optional user profile data to the Python expert service.
4. The expert service loads `published` rules from PostgreSQL.
5. The inference engine matches rule premises against the submitted symptoms using Forward Chaining.
6. The engine calculates Certainty Factor scores, ranks diagnosis candidates, and maps the top result to a risk level.
7. The gateway stores the final assessment result in `health_assessments`.

## Repository Structure

```text
.
|-- gateway-go/
|   |-- cmd/server/main.go          # API gateway entrypoint
|   |-- internal/database/          # PostgreSQL connection
|   |-- internal/handlers/          # HTTP handlers
|   |-- internal/middleware/        # auth, CORS, RBAC
|   |-- internal/models/            # response and domain models
|   |-- internal/repository/        # PostgreSQL queries
|   |-- internal/services/          # business logic
|   |-- migrations/                 # SQL migrations
|   |-- Dockerfile
|   |-- go.mod
|   `-- go.sum
|
|-- expert-py/
|   |-- main.py                     # FastAPI entrypoint
|   |-- app/database.py             # PostgreSQL connection
|   |-- app/input_layer/            # Pydantic schemas
|   |-- app/knowledge_base/         # rule loader
|   |-- app/logic_engine/           # Forward Chaining and CF logic
|   |-- app/output_layer/           # risk classifier and recommendations
|   |-- tests/
|   |-- requirements.txt
|   `-- Dockerfile
|
|-- web-admin/
|   |-- app/                        # Next.js app router
|   |-- components/
|   |-- lib/api.ts                  # API client for the gateway
|   |-- package.json
|   |-- package-lock.json
|   `-- Dockerfile
|
|-- infra/docker/docker-compose.yml
`-- .gitignore
```

## Requirements

### Docker Setup

- Docker Engine
- Docker Compose

Docker images used by the stack:

- `postgres:16-alpine`
- `golang:1.25.10-alpine3.23` for the gateway build stage
- `python:3.12-slim` for the expert service
- `node:22-alpine` for the admin dashboard

### Local Setup Without Docker

- Go 1.25+
- Python 3.12+
- Node.js 22+
- npm
- PostgreSQL 16+
- A database migration tool. `golang-migrate` is recommended, but the SQL files can also be applied manually with `psql`.

### Main Dependencies

Go gateway:

- Gin HTTP framework
- gin-contrib/cors
- pgx PostgreSQL driver
- golang-jwt/jwt
- Google API ID token validator
- bcrypt from `golang.org/x/crypto`

Python expert service:

- FastAPI
- Uvicorn
- Pydantic
- psycopg
- pytest

Web admin:

- Next.js 15
- React 19
- Recharts
- Tailwind CSS
- TypeScript

## Environment Configuration

Create `infra/docker/.env` for Docker Compose.

```env
POSTGRES_USER=balitravelhealthdb
POSTGRES_PASSWORD=change_this_password
POSTGRES_DB=balitravelhealth
DB_PORT=5432

JWT_SECRET=change_this_to_a_long_random_secret
GOOGLE_OAUTH_CLIENT_IDS=your-google-client-id.apps.googleusercontent.com

NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Gateway environment variables:

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `DB_HOST` | Yes | - | PostgreSQL host. In Docker Compose: `db`. |
| `DB_PORT` | Yes | - | PostgreSQL port, usually `5432`. |
| `POSTGRES_USER` | Yes | - | Database username. |
| `POSTGRES_PASSWORD` | Yes | - | Database password. |
| `POSTGRES_DB` | Yes | - | Database name. |
| `JWT_SECRET` | Yes | - | HS256 secret for access tokens. Use a long random value. |
| `GOOGLE_OAUTH_CLIENT_IDS` | Required for Google login | - | One or more Google OAuth client IDs, separated by commas. |
| `EXPERT_SERVICE_URL` | Yes | - | Expert service URL. In Docker Compose: `http://expert:8001`. |
| `PORT` | No | `8080` | Gateway HTTP port. |
| `GIN_MODE` | No | `debug` | Use `release` in production. |
| `CORS_ALLOWED_ORIGINS` | No | `*` | Comma-separated allowed origins. If needed in Docker Compose, add it to the gateway service environment. |

Expert service environment variables:

| Variable | Required | Code Default | Description |
| --- | --- | --- | --- |
| `DB_HOST` | No | `db` | PostgreSQL host. |
| `DB_PORT` | No | `5432` | PostgreSQL port. |
| `POSTGRES_DB` | No | `balitravelhealth` | Database name. |
| `POSTGRES_USER` | No | `balitravelhealthdb` | Database username. |
| `POSTGRES_PASSWORD` | Yes | empty | Database password. |

Web admin environment variables:

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `NEXT_PUBLIC_API_BASE_URL` | No | `http://localhost:8080` | Gateway base URL used by the browser. |

## Running with Docker Compose

1. Open the Docker Compose directory:

```bash
cd infra/docker
```

2. Create `.env` using the values from [Environment Configuration](#environment-configuration).

3. Start the stack:

```bash
docker compose --env-file .env up -d --build
```

4. Run database migrations:

```bash
migrate -path ../../gateway-go/migrations \
  -database "postgres://balitravelhealthdb:change_this_password@localhost:5432/balitravelhealth?sslmode=disable" \
  up
```

5. Check the services:

```bash
curl http://localhost:8080/health
docker compose exec expert python -c "import urllib.request; print(urllib.request.urlopen('http://localhost:8001/health').read().decode())"
```

Notes:

- In the current `docker-compose.yml`, the gateway `8080`, PostgreSQL `5432`, and web admin `3000` ports are published only to `127.0.0.1`.
- The expert service is not published to the host by default. It is accessed internally by the gateway through `http://expert:8001`.
- The web admin dashboard is available at `http://localhost:3000`.

## Running Locally

### 1. Prepare PostgreSQL

Create the database and user:

```sql
CREATE DATABASE balitravelhealth;
CREATE USER balitravelhealthdb WITH PASSWORD 'change_this_password';
GRANT ALL PRIVILEGES ON DATABASE balitravelhealth TO balitravelhealthdb;
```

Then apply all migrations in `gateway-go/migrations`.

### 2. Run the Expert Service

```bash
cd expert-py
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

export DB_HOST=localhost
export DB_PORT=5432
export POSTGRES_USER=balitravelhealthdb
export POSTGRES_PASSWORD=change_this_password
export POSTGRES_DB=balitravelhealth

uvicorn main:app --reload --host 0.0.0.0 --port 8001
```

PowerShell:

```powershell
cd expert-py
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt

$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:POSTGRES_USER="balitravelhealthdb"
$env:POSTGRES_PASSWORD="change_this_password"
$env:POSTGRES_DB="balitravelhealth"

uvicorn main:app --reload --host 0.0.0.0 --port 8001
```

### 3. Run the Go Gateway

```bash
cd gateway-go
go mod download

export DB_HOST=localhost
export DB_PORT=5432
export POSTGRES_USER=balitravelhealthdb
export POSTGRES_PASSWORD=change_this_password
export POSTGRES_DB=balitravelhealth
export JWT_SECRET=change_this_to_a_long_random_secret
export GOOGLE_OAUTH_CLIENT_IDS=your-google-client-id.apps.googleusercontent.com
export EXPERT_SERVICE_URL=http://localhost:8001

go run ./cmd/server
```

PowerShell:

```powershell
cd gateway-go
go mod download

$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:POSTGRES_USER="balitravelhealthdb"
$env:POSTGRES_PASSWORD="change_this_password"
$env:POSTGRES_DB="balitravelhealth"
$env:JWT_SECRET="change_this_to_a_long_random_secret"
$env:GOOGLE_OAUTH_CLIENT_IDS="your-google-client-id.apps.googleusercontent.com"
$env:EXPERT_SERVICE_URL="http://localhost:8001"

go run ./cmd/server
```

### 4. Run the Web Admin

```bash
cd web-admin
npm install
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run dev
```

PowerShell:

```powershell
cd web-admin
npm install
$env:NEXT_PUBLIC_API_BASE_URL="http://localhost:8080"
npm run dev
```

## Database Migrations

The `gateway-go/migrations` directory contains 25 `up` and `down` migrations.

The migrations cover:

- Users and authentication provider enum.
- Health profiles.
- Health assessments.
- Refresh tokens.
- Roles, permissions, user roles, and role permissions.
- Traveler profiles.
- Nurse profiles.
- Nursing care records.
- Vaccination records.
- Bali destinations.
- Medical facilities.
- AOI locations.
- Health risks.
- Emergency guides and emergency guide flows.
- Expert symptoms.
- Expert diseases.
- Expert rules.
- Initial seed data, expert knowledge base, emergency guide data, and SOP knowledge base data.

Using `golang-migrate`:

```bash
migrate -path gateway-go/migrations \
  -database "postgres://balitravelhealthdb:change_this_password@localhost:5432/balitravelhealth?sslmode=disable" \
  up
```

Rollback one migration:

```bash
migrate -path gateway-go/migrations \
  -database "postgres://balitravelhealthdb:change_this_password@localhost:5432/balitravelhealth?sslmode=disable" \
  down 1
```

If you do not use a migration tool, apply all `*.up.sql` files manually in filename order.

## API Endpoints

Default gateway base URL: `http://localhost:8080`.

Protected endpoints require:

```http
Authorization: Bearer <access_token>
```

### Health

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| GET | `/health` | No | Gateway and database health check. |
| GET | `/uploads/*` | No | Static files uploaded from admin. |

### Auth

| Method | Endpoint | Auth | Body |
| --- | --- | --- | --- |
| POST | `/auth/google` | No | `{ "id_token": "...", "device_info": "Android" }` |
| POST | `/auth/refresh` | No | `{ "refresh_token": "...", "device_info": "Android" }` |
| POST | `/auth/logout` | No | `{ "refresh_token": "..." }` |
| POST | `/admin/auth/login` | No | `{ "email": "...", "password": "...", "device_info": "browser" }` |
| POST | `/admin/bootstrap` | No | `{ "email": "...", "password": "min-8-char" }` |

`/admin/bootstrap` only succeeds when no admin account exists yet.

### Public Content

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/location/classify?lat=-8.65&lng=115.21` | Checks whether the coordinates are inside Bali and returns an approximate region. |
| GET | `/facilities/nearby?lat=-8.65&lng=115.21&radius_km=20&limit=10` | Returns nearby medical facilities. |
| GET | `/destinations` | Lists destinations or Bali regions. |
| GET | `/destinations/:id/health-risks` | Lists health risks for a destination. |
| GET | `/emergency-guides` | Lists step-based emergency guides. |
| GET | `/emergency-guide-flows` | Lists decision-tree emergency guides. |
| GET | `/emergency-guide-flows/:id` | Returns one decision-tree emergency guide. |
| GET | `/expert/symptoms?kategori=pre_travel` | Lists public symptoms for a diagnosis category. |

### Traveler Protected Endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/health-profile` | Get the authenticated user's health profile. |
| POST | `/health-profile` | Create a health profile. |
| PUT | `/health-profile` | Update a health profile. |
| GET | `/traveler-profile` | Get the traveler profile. |
| POST | `/traveler-profile` | Create a traveler profile. |
| PUT | `/traveler-profile` | Update a traveler profile. |
| POST | `/assessment` | Submit symptoms and run diagnosis. |
| GET | `/assessments?page=1&limit=10` | List the authenticated user's assessment history. |
| GET | `/vaccinations` | List vaccination records. |
| POST | `/vaccinations` | Create a vaccination record. |
| DELETE | `/vaccinations/:id` | Delete a vaccination record. |
| GET | `/nurses` | List active nurses. |
| POST | `/nursing/appointments` | Create a nurse appointment. |
| GET | `/nursing/my-records` | List the traveler's own nursing care records. |
| GET | `/nursing/nurse-records` | List records assigned to the logged-in nurse. |
| PUT | `/nursing/records/:id` | Update a nursing care record assigned to the logged-in nurse. |

Example assessment request:

```bash
curl -X POST http://localhost:8080/assessment \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symptoms": [1, 2, 3],
    "kategori": "pre_travel"
  }'
```

### Admin Protected Endpoints

`/admin/*` endpoints require a user with the `admin` or `nurse` role, except where service-level business rules further restrict behavior.

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/admin/upload` | Upload an image and return an `/uploads/...` URL. |
| GET | `/admin/facilities` | List medical facilities. |
| POST | `/admin/facilities` | Create a medical facility. |
| PUT | `/admin/facilities/:id` | Update a medical facility. |
| DELETE | `/admin/facilities/:id` | Delete a medical facility. |
| POST | `/admin/destinations` | Create a destination. |
| PUT | `/admin/destinations/:id` | Update a destination. |
| DELETE | `/admin/destinations/:id` | Delete a destination. |
| POST | `/admin/health-risks` | Create a health risk. |
| PUT | `/admin/health-risks/:id` | Update a health risk. |
| DELETE | `/admin/health-risks/:id` | Delete a health risk. |
| POST | `/admin/emergency-guides` | Create a step-based emergency guide. |
| PUT | `/admin/emergency-guides/:id` | Update a step-based emergency guide. |
| DELETE | `/admin/emergency-guides/:id` | Delete a step-based emergency guide. |
| GET | `/admin/emergency-guide-flows` | List emergency guide flows for admin. |
| POST | `/admin/emergency-guide-flows` | Create an emergency guide flow. |
| PUT | `/admin/emergency-guide-flows/:id` | Update an emergency guide flow. |
| DELETE | `/admin/emergency-guide-flows/:id` | Delete an emergency guide flow. |
| GET | `/admin/nurses` | List all nurses. |
| POST | `/admin/nurses` | Create a nurse account. |
| PUT | `/admin/nurses/:id/toggle` | Activate or deactivate a nurse. |
| GET | `/admin/assessments?page=1&limit=20` | List all assessments. |
| GET | `/admin/expert/symptoms` | List master symptoms. |
| POST | `/admin/expert/symptoms` | Create a master symptom. |
| PUT | `/admin/expert/symptoms/:id` | Update a master symptom. |
| DELETE | `/admin/expert/symptoms/:id` | Delete a master symptom. |
| GET | `/admin/expert/diseases` | List master diseases. |
| POST | `/admin/expert/diseases` | Create a master disease. |
| PUT | `/admin/expert/diseases/:id` | Update a master disease. |
| DELETE | `/admin/expert/diseases/:id` | Delete a master disease. |
| GET | `/admin/expert/rules` | List expert-system rules. |
| POST | `/admin/expert/rules` | Create a draft rule. |
| PUT | `/admin/expert/rules/:id` | Update a rule. |
| DELETE | `/admin/expert/rules/:id` | Delete a rule. |
| POST | `/admin/expert/rules/:id/publish` | Publish a rule. |
| POST | `/admin/expert/rules/:id/unpublish` | Unpublish a rule. |

Bootstrap the first admin:

```bash
curl -X POST http://localhost:8080/admin/bootstrap \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "change_this_password"
  }'
```

Admin login:

```bash
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "change_this_password",
    "device_info": "browser"
  }'
```

## Expert System

The expert service is located in `expert-py`.

Endpoints:

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/health` | Basic health check. |
| GET | `/health/db` | Health check with database/rule loading verification. |
| POST | `/diagnose` | Diagnose based on symptoms and optional user profile. |

`/diagnose` request:

```json
{
  "symptoms": [1, 2, 3],
  "kategori": "pre_travel",
  "user_profile": {
    "age": 30,
    "gender": "male",
    "weight_kg": 70
  }
}
```

Response:

```json
{
  "diagnosis": "Top diagnosis name",
  "confidence_score": 0.88,
  "risk_level": "Darurat",
  "recommendation": "Recommended action",
  "all_results": [
    {
      "disease_id": 1,
      "disease_nama": "Disease name",
      "confidence_score": 0.88,
      "risk_level": "Darurat",
      "recommendation": "Recommended action"
    }
  ]
}
```

Core logic:

- Only rules with `status = 'published'` are used.
- Rules can be filtered by `pre_travel` or `post_travel`.
- Forward Chaining matches the entire rule `premis` against input symptoms.
- A single rule CF is calculated as `MB - MD`.
- Multiple matching rules for the same disease are combined with:

```text
CFcombine = CFold + CFnew * (1 - CFold)
```

- User profile data can boost risk scores for selected conditions, for example age `>= 60` for cardiovascular or heat-related diagnoses.
- Risk level mapping. The API currently returns risk labels in Indonesian because these values are part of the stored domain model:

| CF | Risk Level |
| --- | --- |
| `>= 0.8` | `Darurat` |
| `>= 0.6` | `Tinggi` |
| `>= 0.4` | `Sedang` |
| `< 0.4` | `Rendah` |

## Database Schema

Main tables:

| Table | Purpose |
| --- | --- |
| `users` | Accounts for travelers, admins, and nurses. |
| `refresh_tokens` | Refresh tokens stored as hashes. |
| `roles` | Role master table: `traveler`, `nurse`, `admin`. |
| `permissions` | Permission master table. |
| `user_roles` | User-to-role relation. |
| `role_permissions` | Role-to-permission relation. |
| `health_profiles` | User health profile data. |
| `travelers` | Traveler profile data. |
| `health_assessments` | Assessment history and diagnosis results. |
| `vaccination_records` | Vaccination history. |
| `nurses` | Nurse profiles. |
| `nursing_care_records` | Nursing care records. |
| `destinations` | Bali destination/region list. |
| `medical_facilities` | Medical facility data. |
| `aoi_locations` | Area of interest location data. |
| `health_risks` | Destination-specific health risks. |
| `emergency_guides` | Step-based emergency guides. |
| `emergency_guide_flows` | JSON decision-tree emergency guides. |
| `expert_symptoms` | Expert-system symptom master data. |
| `expert_diseases` | Expert-system disease/diagnosis master data. |
| `expert_rules` | IF-THEN rules with Certainty Factor values. |

Seed data includes:

- Initial roles: `traveler`, `nurse`, and `admin`.
- 9 Bali regencies/cities.
- Sample medical facilities in Bali.
- Expert-system symptoms, diseases, and rules.
- Emergency guides and emergency guide flows.
- Additional SOP knowledge base data.

## Testing

### Expert Service

```bash
cd expert-py
pytest
```

Existing tests cover:

- CF formula.
- CF combination.
- Forward Chaining symptom matcher.
- CF aggregation by disease.
- Risk level classification.

### Go Gateway

```bash
cd gateway-go
go test ./...
```

Notes:

- Some gateway tests are integration tests and require PostgreSQL to be running with migrations applied.
- If the database is unavailable, RBAC tests are skipped.
- Make sure database environment variables and `JWT_SECRET` are available when running the full test suite.

### Web Admin

```bash
cd web-admin
npm run build
```

The `npm run lint` script exists in `package.json`, but verify the lint setup against the Next.js version used by this project.

## Security and Open Source Checklist

Before publishing this repository:

- Add a `LICENSE` file for the project itself.
- Do not commit `.env`, private keys, service accounts, OAuth secrets, or production credentials.
- Replace all example passwords and secrets.
- Use a long, random, environment-specific `JWT_SECRET`.
- Set `GIN_MODE=release` in production.
- Restrict CORS with `CORS_ALLOWED_ORIGINS`; do not use wildcard CORS for production.
- Run the gateway behind HTTPS or a trusted reverse proxy.
- Review all seed data and medical content before public or production use.
- Add privacy policy and data retention policy documentation because the system processes health-related data.
- Back up the database regularly.
- Use least-privilege database accounts in production.
- Consider rate limiting for authentication and assessment endpoints.
- Add production observability such as structured logs, metrics, and alerts.
- Run dependency audits and vulnerability scans before each release.

Suggested audit commands:

```bash
# Go
go test ./...
go list -m all

# Python
pip install pip-audit pip-licenses
pip-audit
pip-licenses

# Node
npm audit
npx license-checker --summary
```

## Open Source Dependency Licenses

> Important: the list below is based on the dependency files currently present in this repository. For an official release, rerun a license audit against the final locked dependency versions, especially because `expert-py/requirements.txt` does not currently pin Python package versions.

### Summary

- The Go gateway mainly uses permissive open source dependencies such as MIT, BSD-3-Clause, and Apache-2.0 licensed packages.
- The Python expert service uses permissive packages such as FastAPI, Pydantic, and pytest; Uvicorn is BSD-3-Clause; `psycopg` uses LGPL-3.0.
- The web admin, based on `package-lock.json`, is mostly MIT, Apache-2.0, ISC, and BSD licensed. Transitive dependencies also include LGPL-3.0-or-later through `sharp/libvips`, MPL-2.0 through `axe-core`, and data licenses such as CC-BY-4.0 and CC0-1.0.

### Main Go Dependencies

| Package | Version in `go.mod` | Common upstream license |
| --- | --- | --- |
| `github.com/gin-gonic/gin` | `v1.12.0` | MIT |
| `github.com/gin-contrib/cors` | `v1.7.7` | MIT |
| `github.com/golang-jwt/jwt/v5` | `v5.3.1` | MIT |
| `github.com/jackc/pgx/v5` | `v5.9.2` | MIT |
| `golang.org/x/crypto` | `v0.51.0` | BSD-3-Clause |
| `google.golang.org/api` | `v0.280.0` | BSD-3-Clause |
| `go.opentelemetry.io/otel` | `v1.43.0` | Apache-2.0 |
| `go.mongodb.org/mongo-driver/v2` | `v2.5.0` | Apache-2.0 |
| `github.com/goccy/go-json` | `v0.10.5` | MIT |
| `github.com/json-iterator/go` | `v1.1.12` | MIT |

Note: `go.mod` currently marks all dependencies as indirect. Run `go mod tidy` once a Go toolchain is available so direct and transitive dependencies are easier to audit.

### Main Python Dependencies

| Package | Version | Common upstream license | Notes |
| --- | --- | --- | --- |
| `fastapi` | Unpinned | MIT | API framework. |
| `uvicorn[standard]` | Unpinned | BSD-3-Clause | ASGI server. |
| `pydantic` | Unpinned | MIT | Schema validation. |
| `pytest` | Unpinned | MIT | Testing. |
| `psycopg[binary]` | Unpinned | LGPL-3.0-only | PostgreSQL driver. Review license obligations when distributing binaries or containers. |

Recommendation: pin Python dependency versions before release, for example with `pip-tools` or another lockfile workflow.

### Direct Node/Web Admin Dependencies

Based on `web-admin/package-lock.json`:

| Package | Locked Version | License |
| --- | --- | --- |
| `autoprefixer` | `10.5.0` | MIT |
| `next` | `15.5.18` | MIT |
| `react` | `19.2.6` | MIT |
| `react-dom` | `19.2.6` | MIT |
| `recharts` | `3.8.1` | MIT |
| `@types/node` | `20.19.41` | MIT |
| `@types/react` | `19.2.15` | MIT |
| `@types/react-dom` | `19.2.3` | MIT |
| `eslint` | `9.39.4` | MIT |
| `eslint-config-next` | `15.5.18` | MIT |
| `postcss` | `8.5.15` | MIT |
| `tailwindcss` | `3.4.19` | MIT |
| `typescript` | `5.9.3` | Apache-2.0 |

Transitive license summary from `package-lock.json`:

| License | Package Count |
| --- | ---: |
| MIT | 358 |
| Apache-2.0 | 34 |
| ISC | 28 |
| LGPL-3.0-or-later | 10 |
| BSD-2-Clause | 7 |
| BSD-3-Clause | 3 |
| Apache-2.0 AND LGPL-3.0-or-later | 3 |
| Not declared in lockfile | 2 |
| Other | 8 |

Transitive packages that deserve extra review:

| Package | Version | License |
| --- | --- | --- |
| `@img/sharp-libvips-*` | `1.2.4` | LGPL-3.0-or-later |
| `@img/sharp-*` | `0.34.5` | Apache-2.0 AND LGPL-3.0-or-later |
| `axe-core` | `4.11.4` | MPL-2.0 |
| `argparse` | `2.0.1` | Python-2.0 |
| `caniuse-lite` | `1.0.30001793` | CC-BY-4.0 |
| `language-subtag-registry` | `0.3.23` | CC0-1.0 |
| `busboy` | `1.6.0` | Not declared in lockfile |
| `streamsearch` | `1.1.0` | Not declared in lockfile |
