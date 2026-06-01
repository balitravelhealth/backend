# BaliTravelHealth Platform Backend

Backend dan dashboard admin untuk platform kesehatan wisatawan Bali. Repository ini berisi API gateway berbasis Go, service sistem pakar berbasis Python/FastAPI, database PostgreSQL, migrasi SQL, dan web admin berbasis Next.js.

> Catatan medis: sistem ini hanya alat bantu informasi dan triase awal. Hasil diagnosis, rekomendasi, dan panduan darurat wajib divalidasi oleh tenaga kesehatan berwenang sebelum digunakan di lingkungan produksi. Sistem ini tidak menggantikan konsultasi dokter, perawat, fasilitas kesehatan, atau layanan gawat darurat.

## Daftar Isi

- [Fitur Utama](#fitur-utama)
- [Arsitektur](#arsitektur)
- [Struktur Repository](#struktur-repository)
- [Requirements](#requirements)
- [Konfigurasi Environment](#konfigurasi-environment)
- [Menjalankan dengan Docker Compose](#menjalankan-dengan-docker-compose)
- [Menjalankan Secara Lokal](#menjalankan-secara-lokal)
- [Migrasi Database](#migrasi-database)
- [API Endpoint](#api-endpoint)
- [Sistem Pakar](#sistem-pakar)
- [Skema Database](#skema-database)
- [Testing](#testing)
- [Keamanan dan Checklist Open Source](#keamanan-dan-checklist-open-source)
- [Lisensi Dependency Open Source](#lisensi-dependency-open-source)
- [Lisensi Project](#lisensi-project)

## Fitur Utama

- Autentikasi pengguna wisatawan dengan Google ID token.
- Autentikasi admin dan perawat dengan email/password.
- JWT access token berdurasi 15 menit dan refresh token berdurasi 30 hari.
- Rotasi refresh token dan deteksi reuse token.
- RBAC untuk role `traveler`, `nurse`, dan `admin`.
- Profil kesehatan wisatawan.
- Profil traveler dan kontak darurat.
- Riwayat vaksinasi.
- Asesmen kesehatan wisatawan dengan integrasi sistem pakar.
- Sistem pakar Forward Chaining + Certainty Factor.
- Knowledge base gejala, penyakit, dan rule yang bisa dikelola admin.
- Kategori diagnosis `pre_travel` dan `post_travel`.
- Klasifikasi lokasi Bali berdasarkan koordinat.
- Daftar destinasi, risiko kesehatan daerah, fasilitas kesehatan, dan fasilitas terdekat.
- Panduan emergency guide berbasis langkah dan decision-tree flow.
- Manajemen perawat dan nursing care record.
- Web admin untuk dashboard, knowledge base, fasilitas, destinasi, panduan darurat, perawat, dan asesmen.
- Dockerfile untuk setiap service dan Docker Compose untuk stack lokal.

## Arsitektur

```text
Client Mobile / Web Admin
        |
        v
Go API Gateway (Gin) :8080
        |
        |-- PostgreSQL :5432
        |
        |-- Python Expert Service (FastAPI) :8001
                |
                |-- membaca published rules dari PostgreSQL
```

Alur diagnosis:

1. Client mengirim gejala ke `POST /assessment`.
2. Gateway membaca profil kesehatan user dari database.
3. Gateway meneruskan gejala, kategori, dan profil user ke `expert-py`.
4. Expert service memuat rule berstatus `published` dari PostgreSQL.
5. Engine mencocokkan premis rule dengan gejala input menggunakan Forward Chaining.
6. Engine menghitung Certainty Factor, mengurutkan kandidat diagnosis, lalu menentukan level risiko.
7. Gateway menyimpan hasil ke tabel `health_assessments`.

## Struktur Repository

```text
.
|-- gateway-go/
|   |-- cmd/server/main.go          # entrypoint API gateway
|   |-- internal/database/          # koneksi PostgreSQL
|   |-- internal/handlers/          # HTTP handlers
|   |-- internal/middleware/        # auth, CORS, RBAC
|   |-- internal/models/            # model response/domain
|   |-- internal/repository/        # query PostgreSQL
|   |-- internal/services/          # business logic
|   |-- migrations/                 # 25 file migrasi SQL
|   |-- Dockerfile
|   |-- go.mod
|   `-- go.sum
|
|-- expert-py/
|   |-- main.py                     # entrypoint FastAPI
|   |-- app/database.py             # koneksi PostgreSQL
|   |-- app/input_layer/            # Pydantic schema
|   |-- app/knowledge_base/         # loader rule dari database
|   |-- app/logic_engine/           # Forward Chaining + CF
|   |-- app/output_layer/           # risk classifier + rekomendasi
|   |-- tests/
|   |-- requirements.txt
|   `-- Dockerfile
|
|-- web-admin/
|   |-- app/                        # Next.js app router
|   |-- components/
|   |-- lib/api.ts                  # API client ke gateway
|   |-- package.json
|   |-- package-lock.json
|   `-- Dockerfile
|
|-- infra/docker/docker-compose.yml
`-- .gitignore
```

## Requirements

### Untuk Docker

- Docker Engine
- Docker Compose

Image yang digunakan:

- `postgres:16-alpine`
- `golang:1.25.10-alpine3.23` untuk build gateway
- `python:3.12-slim` untuk expert service
- `node:22-alpine` untuk web admin

### Untuk menjalankan lokal tanpa Docker

- Go 1.25+
- Python 3.12+
- Node.js 22+
- npm
- PostgreSQL 16+
- Tool migrasi database. Direkomendasikan `golang-migrate`, tetapi file SQL juga bisa dijalankan manual dengan `psql`.

### Dependency utama

Gateway Go:

- Gin HTTP framework
- gin-contrib/cors
- pgx PostgreSQL driver
- golang-jwt/jwt
- Google API ID token validator
- bcrypt dari `golang.org/x/crypto`

Expert Python:

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

## Konfigurasi Environment

Buat file `.env` di `infra/docker/.env` untuk Docker Compose.

```env
POSTGRES_USER=balitravelhealthdb
POSTGRES_PASSWORD=change_this_password
POSTGRES_DB=balitravelhealth
DB_PORT=5432

JWT_SECRET=change_this_to_a_long_random_secret
GOOGLE_OAUTH_CLIENT_IDS=your-google-client-id.apps.googleusercontent.com

NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Variable yang dikenali gateway:

| Variable | Wajib | Default | Keterangan |
| --- | --- | --- | --- |
| `DB_HOST` | Ya | - | Host PostgreSQL. Di Compose: `db`. |
| `DB_PORT` | Ya | - | Port PostgreSQL. Umumnya `5432`. |
| `POSTGRES_USER` | Ya | - | Username database. |
| `POSTGRES_PASSWORD` | Ya | - | Password database. |
| `POSTGRES_DB` | Ya | - | Nama database. |
| `JWT_SECRET` | Ya | - | Secret HS256 untuk access token. Gunakan string panjang dan acak. |
| `GOOGLE_OAUTH_CLIENT_IDS` | Ya untuk login Google | - | Satu atau beberapa Google OAuth client ID, dipisah koma. |
| `EXPERT_SERVICE_URL` | Ya | - | URL service expert. Di Compose: `http://expert:8001`. |
| `PORT` | Tidak | `8080` | Port gateway. |
| `GIN_MODE` | Tidak | `debug` | Gunakan `release` di produksi. |
| `CORS_ALLOWED_ORIGINS` | Tidak | `*` | Origin yang diizinkan, dipisah koma. Jika dipakai di Compose, tambahkan variable ini ke environment service gateway. |

Variable expert service:

| Variable | Wajib | Default di kode | Keterangan |
| --- | --- | --- | --- |
| `DB_HOST` | Tidak | `db` | Host PostgreSQL. |
| `DB_PORT` | Tidak | `5432` | Port PostgreSQL. |
| `POSTGRES_DB` | Tidak | `balitravelhealth` | Nama database. |
| `POSTGRES_USER` | Tidak | `balitravelhealthdb` | Username database. |
| `POSTGRES_PASSWORD` | Ya | kosong | Password database. |

Variable web admin:

| Variable | Wajib | Default | Keterangan |
| --- | --- | --- | --- |
| `NEXT_PUBLIC_API_BASE_URL` | Tidak | `http://localhost:8080` | Base URL gateway yang dipanggil browser. |

## Menjalankan dengan Docker Compose

1. Masuk ke folder Compose:

```bash
cd infra/docker
```

2. Buat file environment `.env` dengan isi pada bagian [Konfigurasi Environment](#konfigurasi-environment).

3. Jalankan stack:

```bash
docker compose --env-file .env up -d --build
```

4. Jalankan migrasi database:

```bash
migrate -path ../../gateway-go/migrations \
  -database "postgres://balitravelhealthdb:change_this_password@localhost:5432/balitravelhealth?sslmode=disable" \
  up
```

5. Cek service:

```bash
curl http://localhost:8080/health
docker compose exec expert python -c "import urllib.request; print(urllib.request.urlopen('http://localhost:8001/health').read().decode())"
```

Catatan:

- Pada `docker-compose.yml` saat ini, port gateway `8080`, PostgreSQL `5432`, dan web admin `3000` dipublish hanya ke `127.0.0.1`.
- Expert service tidak dipublish ke host oleh Compose. Service ini diakses internal oleh gateway melalui `http://expert:8001`.
- Web admin tersedia di `http://localhost:3000`.

## Menjalankan Secara Lokal

### 1. Siapkan PostgreSQL

Buat database:

```sql
CREATE DATABASE balitravelhealth;
CREATE USER balitravelhealthdb WITH PASSWORD 'change_this_password';
GRANT ALL PRIVILEGES ON DATABASE balitravelhealth TO balitravelhealthdb;
```

Lalu jalankan semua migrasi pada folder `gateway-go/migrations`.

### 2. Jalankan expert service

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

Untuk PowerShell:

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

### 3. Jalankan gateway Go

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

Untuk PowerShell:

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

### 4. Jalankan web admin

```bash
cd web-admin
npm install
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run dev
```

Untuk PowerShell:

```powershell
cd web-admin
npm install
$env:NEXT_PUBLIC_API_BASE_URL="http://localhost:8080"
npm run dev
```

## Migrasi Database

Folder `gateway-go/migrations` berisi 25 migrasi `up` dan `down`.

Isi migrasi mencakup:

- User dan auth provider.
- Profil kesehatan.
- Asesmen kesehatan.
- Refresh token.
- Role, permission, user role, dan role permission.
- Traveler.
- Nurse.
- Nursing care record.
- Vaccination record.
- Destinasi Bali.
- Fasilitas kesehatan.
- AOI location.
- Health risk.
- Emergency guide dan emergency guide flow.
- Expert symptoms.
- Expert diseases.
- Expert rules.
- Seed data awal, knowledge base, emergency guide, dan SOP knowledge base.

Dengan `golang-migrate`:

```bash
migrate -path gateway-go/migrations \
  -database "postgres://balitravelhealthdb:change_this_password@localhost:5432/balitravelhealth?sslmode=disable" \
  up
```

Rollback satu langkah:

```bash
migrate -path gateway-go/migrations \
  -database "postgres://balitravelhealthdb:change_this_password@localhost:5432/balitravelhealth?sslmode=disable" \
  down 1
```

Jika tidak memakai tool migrasi, jalankan file `*.up.sql` secara berurutan berdasarkan nama file.

## API Endpoint

Base URL gateway default: `http://localhost:8080`.

Endpoint protected membutuhkan header:

```http
Authorization: Bearer <access_token>
```

### Health

| Method | Endpoint | Auth | Keterangan |
| --- | --- | --- | --- |
| GET | `/health` | Tidak | Health check gateway dan koneksi database. |
| GET | `/uploads/*` | Tidak | Static file hasil upload admin. |

### Auth

| Method | Endpoint | Auth | Body |
| --- | --- | --- | --- |
| POST | `/auth/google` | Tidak | `{ "id_token": "...", "device_info": "Android" }` |
| POST | `/auth/refresh` | Tidak | `{ "refresh_token": "...", "device_info": "Android" }` |
| POST | `/auth/logout` | Tidak | `{ "refresh_token": "..." }` |
| POST | `/admin/auth/login` | Tidak | `{ "email": "...", "password": "...", "device_info": "browser" }` |
| POST | `/admin/bootstrap` | Tidak | `{ "email": "...", "password": "min-8-char" }` |

`/admin/bootstrap` hanya berhasil jika belum ada admin di database.

### Public Content

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| GET | `/location/classify?lat=-8.65&lng=115.21` | Cek apakah koordinat berada di Bali dan region perkiraan. |
| GET | `/facilities/nearby?lat=-8.65&lng=115.21&radius_km=20&limit=10` | Fasilitas kesehatan terdekat. |
| GET | `/destinations` | Daftar destinasi/kabupaten/kota. |
| GET | `/destinations/:id/health-risks` | Risiko kesehatan per destinasi. |
| GET | `/emergency-guides` | Panduan emergency berbasis langkah. |
| GET | `/emergency-guide-flows` | Daftar panduan emergency decision-tree. |
| GET | `/emergency-guide-flows/:id` | Detail panduan emergency decision-tree. |
| GET | `/expert/symptoms?kategori=pre_travel` | Daftar gejala publik untuk kategori tertentu. |

### Traveler Protected

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| GET | `/health-profile` | Ambil profil kesehatan user. |
| POST | `/health-profile` | Buat profil kesehatan. |
| PUT | `/health-profile` | Update profil kesehatan. |
| GET | `/traveler-profile` | Ambil profil traveler. |
| POST | `/traveler-profile` | Buat profil traveler. |
| PUT | `/traveler-profile` | Update profil traveler. |
| POST | `/assessment` | Submit gejala dan jalankan diagnosis. |
| GET | `/assessments?page=1&limit=10` | Riwayat asesmen user. |
| GET | `/vaccinations` | Riwayat vaksin user. |
| POST | `/vaccinations` | Tambah riwayat vaksin. |
| DELETE | `/vaccinations/:id` | Hapus riwayat vaksin. |
| GET | `/nurses` | Daftar perawat aktif. |
| POST | `/nursing/appointments` | Buat appointment perawat. |
| GET | `/nursing/my-records` | Riwayat nursing care milik traveler. |
| GET | `/nursing/nurse-records` | Record yang ditugaskan ke perawat login. |
| PUT | `/nursing/records/:id` | Update nursing care record oleh perawat terkait. |

Contoh submit assessment:

```bash
curl -X POST http://localhost:8080/assessment \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "symptoms": [1, 2, 3],
    "kategori": "pre_travel"
  }'
```

### Admin Protected

Endpoint `/admin/*` membutuhkan user dengan role `admin` atau `nurse`, kecuali operasi tertentu dibatasi secara service logic.

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| POST | `/admin/upload` | Upload gambar, mengembalikan URL `/uploads/...`. |
| GET | `/admin/facilities` | List fasilitas. |
| POST | `/admin/facilities` | Buat fasilitas. |
| PUT | `/admin/facilities/:id` | Update fasilitas. |
| DELETE | `/admin/facilities/:id` | Hapus fasilitas. |
| POST | `/admin/destinations` | Buat destinasi. |
| PUT | `/admin/destinations/:id` | Update destinasi. |
| DELETE | `/admin/destinations/:id` | Hapus destinasi. |
| POST | `/admin/health-risks` | Buat risiko kesehatan. |
| PUT | `/admin/health-risks/:id` | Update risiko kesehatan. |
| DELETE | `/admin/health-risks/:id` | Hapus risiko kesehatan. |
| POST | `/admin/emergency-guides` | Buat emergency guide langkah. |
| PUT | `/admin/emergency-guides/:id` | Update emergency guide langkah. |
| DELETE | `/admin/emergency-guides/:id` | Hapus emergency guide langkah. |
| GET | `/admin/emergency-guide-flows` | List emergency flow untuk admin. |
| POST | `/admin/emergency-guide-flows` | Buat emergency flow. |
| PUT | `/admin/emergency-guide-flows/:id` | Update emergency flow. |
| DELETE | `/admin/emergency-guide-flows/:id` | Hapus emergency flow. |
| GET | `/admin/nurses` | List semua perawat. |
| POST | `/admin/nurses` | Buat akun perawat. |
| PUT | `/admin/nurses/:id/toggle` | Aktif/nonaktifkan perawat. |
| GET | `/admin/assessments?page=1&limit=20` | List semua asesmen. |
| GET | `/admin/expert/symptoms` | List gejala master. |
| POST | `/admin/expert/symptoms` | Buat gejala master. |
| PUT | `/admin/expert/symptoms/:id` | Update gejala master. |
| DELETE | `/admin/expert/symptoms/:id` | Hapus gejala master. |
| GET | `/admin/expert/diseases` | List penyakit master. |
| POST | `/admin/expert/diseases` | Buat penyakit master. |
| PUT | `/admin/expert/diseases/:id` | Update penyakit master. |
| DELETE | `/admin/expert/diseases/:id` | Hapus penyakit master. |
| GET | `/admin/expert/rules` | List rule sistem pakar. |
| POST | `/admin/expert/rules` | Buat rule draft. |
| PUT | `/admin/expert/rules/:id` | Update rule. |
| DELETE | `/admin/expert/rules/:id` | Hapus rule. |
| POST | `/admin/expert/rules/:id/publish` | Publish rule. |
| POST | `/admin/expert/rules/:id/unpublish` | Unpublish rule. |

Contoh bootstrap admin:

```bash
curl -X POST http://localhost:8080/admin/bootstrap \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "change_this_password"
  }'
```

Contoh login admin:

```bash
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "change_this_password",
    "device_info": "browser"
  }'
```

## Sistem Pakar

Expert service tersedia di `expert-py`.

Endpoint:

| Method | Endpoint | Keterangan |
| --- | --- | --- |
| GET | `/health` | Health check sederhana. |
| GET | `/health/db` | Health check plus cek rule published di database. |
| POST | `/diagnose` | Diagnosis berdasarkan gejala dan profil user. |

Request `/diagnose`:

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
  "diagnosis": "Nama diagnosis teratas",
  "confidence_score": 0.88,
  "risk_level": "Darurat",
  "recommendation": "Rekomendasi tindakan",
  "all_results": [
    {
      "disease_id": 1,
      "disease_nama": "Nama penyakit",
      "confidence_score": 0.88,
      "risk_level": "Darurat",
      "recommendation": "Rekomendasi tindakan"
    }
  ]
}
```

Logika inti:

- Rule hanya digunakan jika `status = 'published'`.
- Rule dapat difilter dengan kategori `pre_travel` atau `post_travel`.
- Forward Chaining mencocokkan seluruh `premis` rule terhadap gejala input.
- CF satu rule dihitung dengan `MB - MD`.
- Beberapa rule untuk penyakit yang sama digabung dengan:

```text
CFcombine = CFold + CFnew * (1 - CFold)
```

- User profile dapat menaikkan skor untuk risiko tertentu, misalnya usia `>= 60` pada diagnosis terkait kardiovaskular/heat-related.
- Mapping risk level:

| CF | Risk Level |
| --- | --- |
| `>= 0.8` | `Darurat` |
| `>= 0.6` | `Tinggi` |
| `>= 0.4` | `Sedang` |
| `< 0.4` | `Rendah` |

## Skema Database

Tabel utama:

| Tabel | Fungsi |
| --- | --- |
| `users` | Akun user, admin, dan perawat. |
| `refresh_tokens` | Refresh token yang disimpan sebagai hash. |
| `roles` | Role `traveler`, `nurse`, `admin`. |
| `permissions` | Master permission. |
| `user_roles` | Relasi user ke role. |
| `role_permissions` | Relasi role ke permission. |
| `health_profiles` | Profil kesehatan user. |
| `travelers` | Profil traveler. |
| `health_assessments` | Riwayat asesmen dan hasil diagnosis. |
| `vaccination_records` | Riwayat vaksinasi. |
| `nurses` | Profil perawat. |
| `nursing_care_records` | Catatan layanan keperawatan. |
| `destinations` | Daftar daerah/destinasi Bali. |
| `medical_facilities` | Fasilitas kesehatan. |
| `aoi_locations` | Area of interest lokasi. |
| `health_risks` | Risiko kesehatan per destinasi. |
| `emergency_guides` | Panduan emergency berbasis langkah. |
| `emergency_guide_flows` | Panduan emergency berbasis decision-tree JSON. |
| `expert_symptoms` | Master gejala sistem pakar. |
| `expert_diseases` | Master penyakit/diagnosis. |
| `expert_rules` | Rule IF-THEN + Certainty Factor. |

Seed data mencakup:

- Role awal: `traveler`, `nurse`, `admin`.
- 9 kabupaten/kota di Bali.
- Contoh fasilitas kesehatan di Bali.
- Knowledge base gejala, penyakit, dan rule.
- Emergency guide dan emergency guide flow.
- SOP knowledge base tambahan.

## Testing

### Expert service

```bash
cd expert-py
pytest
```

Test yang tersedia mengecek:

- Rumus CF.
- Kombinasi CF.
- Forward Chaining symptom matcher.
- Agregasi CF per penyakit.
- Klasifikasi level risiko.

### Gateway Go

```bash
cd gateway-go
go test ./...
```

Catatan:

- Beberapa test gateway bersifat integration test dan membutuhkan PostgreSQL yang sudah hidup serta sudah dimigrasi.
- Jika database tidak tersedia, test RBAC akan melakukan skip.
- Pastikan environment database dan `JWT_SECRET` tersedia saat menjalankan test penuh.

### Web admin

```bash
cd web-admin
npm run build
```

Script `npm run lint` tersedia di `package.json`, tetapi pastikan konfigurasi lint Next.js sesuai versi Next.js yang digunakan.

## Keamanan dan Checklist Open Source

Sebelum repository dipublish:

- Tambahkan file `LICENSE` untuk lisensi project sendiri.
- Jangan commit file `.env`, private key, service account, OAuth secret, atau credential lain.
- Ganti semua contoh password dan secret.
- Gunakan `JWT_SECRET` panjang, acak, dan berbeda per environment.
- Set `GIN_MODE=release` di production.
- Batasi CORS dengan `CORS_ALLOWED_ORIGINS`, jangan gunakan wildcard untuk production.
- Jalankan gateway di belakang HTTPS/reverse proxy.
- Review ulang seed data dan konten medis sebelum publik.
- Tambahkan privacy policy dan data retention policy karena project memproses data kesehatan.
- Backup database secara berkala.
- Gunakan akun database dengan privilege minimum untuk production.
- Pertimbangkan rate limiting untuk endpoint auth dan assessment.
- Tambahkan observability production: structured logging, metrics, dan alerting.
- Jalankan audit dependency dan vulnerability scan sebelum release.

Contoh audit dependency:

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

## Lisensi Dependency Open Source

> Penting: daftar di bawah disusun dari file dependency yang ada di repository. Untuk rilis resmi, jalankan audit lisensi ulang pada versi dependency final, terutama karena `expert-py/requirements.txt` belum mengunci versi package Python.

### Ringkasan

- Go gateway memakai dependency open source dengan lisensi permissive seperti MIT, BSD-3-Clause, dan Apache-2.0 pada dependency utama.
- Python expert service memakai FastAPI/Pydantic/pytest yang permissive, Uvicorn BSD-3-Clause, dan `psycopg` yang menggunakan LGPL-3.0.
- Web admin berdasarkan `package-lock.json` memakai mayoritas MIT, Apache-2.0, ISC, dan BSD. Transitive dependency juga memuat LGPL-3.0-or-later melalui paket `sharp/libvips`, MPL-2.0 melalui `axe-core`, serta beberapa lisensi data seperti CC-BY-4.0 dan CC0-1.0.

### Dependency utama Go

| Package | Versi di `go.mod` | Lisensi umum upstream |
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

Catatan: `go.mod` saat ini menandai semua dependency sebagai indirect. Rapikan dengan `go mod tidy` setelah environment Go tersedia agar dependency langsung dan transitive lebih mudah diaudit.

### Dependency utama Python

| Package | Versi | Lisensi umum upstream | Keterangan |
| --- | --- | --- | --- |
| `fastapi` | Tidak dipin | MIT | Framework API. |
| `uvicorn[standard]` | Tidak dipin | BSD-3-Clause | ASGI server. |
| `pydantic` | Tidak dipin | MIT | Validasi schema. |
| `pytest` | Tidak dipin | MIT | Testing. |
| `psycopg[binary]` | Tidak dipin | LGPL-3.0-only | Driver PostgreSQL. Perhatikan kewajiban lisensi saat mendistribusikan binary/container. |

Rekomendasi: pin versi Python dependency sebelum rilis, misalnya dengan `pip-tools` atau lockfile lain.

### Dependency langsung Node/Web Admin

Berdasarkan `web-admin/package-lock.json`:

| Package | Versi terkunci | Lisensi |
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

Ringkasan lisensi transitive dari `package-lock.json`:

| Lisensi | Jumlah package |
| --- | ---: |
| MIT | 358 |
| Apache-2.0 | 34 |
| ISC | 28 |
| LGPL-3.0-or-later | 10 |
| BSD-2-Clause | 7 |
| BSD-3-Clause | 3 |
| Apache-2.0 AND LGPL-3.0-or-later | 3 |
| Tidak mendeklarasikan lisensi di lockfile | 2 |
| Lainnya | 8 |

Package transitive yang perlu diperhatikan:

| Package | Versi | Lisensi |
| --- | --- | --- |
| `@img/sharp-libvips-*` | `1.2.4` | LGPL-3.0-or-later |
| `@img/sharp-*` | `0.34.5` | Apache-2.0 AND LGPL-3.0-or-later |
| `axe-core` | `4.11.4` | MPL-2.0 |
| `argparse` | `2.0.1` | Python-2.0 |
| `caniuse-lite` | `1.0.30001793` | CC-BY-4.0 |
| `language-subtag-registry` | `0.3.23` | CC0-1.0 |
| `busboy` | `1.6.0` | Tidak dideklarasikan di lockfile |
| `streamsearch` | `1.1.0` | Tidak dideklarasikan di lockfile |

## Lisensi Project

Repository ini belum memiliki file `LICENSE`. Jika tujuan Anda adalah merilis sebagai open source, pilih dan tambahkan lisensi project sebelum publikasi.

Pilihan umum:

- MIT: sederhana dan permissive.
- Apache-2.0: permissive, dengan grant paten eksplisit.
- GPL-3.0: copyleft kuat.

Pastikan lisensi project kompatibel dengan dependency yang dipakai dan kebijakan distribusi Anda, terutama untuk dependency LGPL pada Python/Node transitive.
