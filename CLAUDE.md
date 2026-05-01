# CLAUDE.md — Dental Record System

Dokumen ini adalah panduan lengkap untuk Claude Code melanjutkan pengembangan sistem manajemen rekam medis klinik gigi. Baca seluruh dokumen ini sebelum menulis satu baris kode pun.

---

## Konteks Proyek

Sistem internal klinik gigi untuk 2 cabang. Diakses lewat browser oleh dokter dan suster. Bukan aplikasi publik — tidak perlu SEO, tidak perlu diindex Google.

**Stack yang sudah disepakati dan tidak boleh diganti:**
- Backend: Go 1.23, framework Chi, database driver pgx/v5
- Database: PostgreSQL 16
- Frontend: React + Vite (dikerjakan terpisah, belum dimulai)
- Deployment: Docker Compose di VPS Ubuntu 22.04
- File storage: lokal di VPS (`/uploads/`)

**Module name Go:** `github.com/yourusername/dental-api`
Ganti `yourusername` dengan username GitHub yang sebenarnya saat setup.

---

## Struktur Folder

```
dental-project/
├── CLAUDE.md                        ← file ini
├── migrations/                      ← SQL files, SUDAH LENGKAP
│   ├── 000_init_enums.sql           ← enums + extensions
│   ├── 001_branches.sql             ← tabel branches + seed
│   ├── 002_users.sql                ← tabel users + seed superadmin
│   ├── 003_patients.sql             ← tabel patients + sequence no_rm + trigger
│   ├── 004_visits.sql               ← tabel visits + trigger status pasien
│   ├── 005_attachments.sql          ← tabel attachments
│   ├── 006_fullpackage.sql          ← audit_logs, control_reminders, export_logs
│   └── 007_views_indexes.sql        ← view v_patient_list + composite indexes
└── dental-api/                      ← Go REST API
    ├── cmd/server/main.go           ← SUDAH LENGKAP — entry point + routing
    ├── config/config.go             ← SUDAH LENGKAP — env config loader
    ├── go.mod                       ← SUDAH LENGKAP — semua dependencies
    ├── Makefile                     ← SUDAH LENGKAP — run/build/migrate helpers
    ├── .env.example                 ← SUDAH LENGKAP — template env vars
    ├── internal/
    │   ├── model/model.go           ← SUDAH LENGKAP — semua domain structs
    │   ├── middleware/auth.go       ← SUDAH LENGKAP — JWT auth middleware
    │   ├── handler/
    │   │   └── auth.go              ← SUDAH LENGKAP — Login, Me, Logout
    │   ├── service/
    │   │   └── auth.go              ← SUDAH LENGKAP — AuthService.Login
    │   ├── repository/
    │   │   ├── user.go              ← SUDAH LENGKAP — FindByEmail, UpdateLastLogin
    │   │   └── patient.go           ← SUDAH LENGKAP — List, GetByID, Create, Update, SoftDelete
    │   └── dto/                     ← Data Transfer Objects (perlu dibuat)
    └── pkg/
        ├── jwt/jwt.go               ← SUDAH LENGKAP — Generate, Verify
        ├── response/response.go     ← SUDAH LENGKAP — JSON helpers
        ├── storage/                 ← PERLU DIBUAT — local file upload
        └── excel/                   ← PERLU DIBUAT — export xlsx
```

---

## Status Implementasi

### ✅ Sudah selesai — JANGAN diubah tanpa alasan kuat

| File | Isi |
|------|-----|
| `migrations/*.sql` | Semua tabel, enum, trigger, view, index |
| `cmd/server/main.go` | Routing lengkap semua endpoint, dependency injection |
| `config/config.go` | Config loader dari env vars |
| `internal/model/model.go` | Semua domain structs (Branch, User, Patient, Visit, Attachment, AuditLog, PatientListRow) |
| `internal/middleware/auth.go` | Authenticate, RequireWrite, RequireSuperAdmin, GetClaims |
| `internal/handler/auth.go` | Login, Me, Logout |
| `internal/service/auth.go` | AuthService dengan Login |
| `internal/repository/user.go` | UserRepo dengan FindByEmail, UpdateLastLogin |
| `internal/repository/patient.go` | PatientRepo dengan List, GetByID, Create, Update, SoftDelete |
| `pkg/jwt/jwt.go` | JWTManager dengan Generate, Verify |
| `pkg/response/response.go` | JSON, WithMeta, Error, Unauthorized, Forbidden, NotFound, BadRequest, InternalError |

### 🔲 Perlu dibuat — target utama

Ikuti urutan ini persis karena ada dependency antar modul:

1. `pkg/storage/local.go` — file upload handler
2. `internal/repository/branch.go` — BranchRepo
3. `internal/repository/visit.go` — VisitRepo
4. `internal/repository/attachment.go` — AttachmentRepo
5. `internal/service/branch.go` — BranchService
6. `internal/service/patient.go` — PatientService
7. `internal/service/visit.go` — VisitService
8. `internal/service/attachment.go` — AttachmentService
9. `internal/service/export.go` — ExportService
10. `internal/service/user_mgmt.go` — UserMgmtService
11. `internal/handler/branch.go`
12. `internal/handler/patient.go`
13. `internal/handler/visit.go`
14. `internal/handler/attachment.go`
15. `internal/handler/export.go`
16. `internal/handler/user_mgmt.go`
17. `pkg/excel/excel.go` — export xlsx
18. `docker-compose.yml` — lengkap dengan semua service
19. `Dockerfile` — multi-stage build

---

## Arsitektur & Pola Kode

### Layered architecture — WAJIB diikuti

```
HTTP Request
    ↓
Handler        → decode request, validate input, call service, encode response
    ↓
Service        → business logic, orchestrate repositories
    ↓
Repository     → SQL queries ke PostgreSQL, TIDAK ada business logic di sini
    ↓
PostgreSQL
```

Aturan keras:
- Handler TIDAK boleh query langsung ke DB
- Repository TIDAK boleh punya business logic
- Service TIDAK boleh encode/decode JSON
- Semua error dari repository di-wrap dengan `fmt.Errorf("context: %w", err)`

### Pola constructor — ikuti ini untuk semua layer

```go
// Repository
type BranchRepo struct { db *pgxpool.Pool }
func NewBranchRepo(db *pgxpool.Pool) *BranchRepo { return &BranchRepo{db: db} }

// Service
type BranchService struct { repo *repository.BranchRepo }
func NewBranchService(repo *repository.BranchRepo) *BranchService { return &BranchService{repo: repo} }

// Handler
type BranchHandler struct { svc *service.BranchService }
func NewBranchHandler(svc *service.BranchService) *BranchHandler { return &BranchHandler{svc: svc} }
```

### Pola response — WAJIB pakai `pkg/response`

```go
// Sukses
response.JSON(w, http.StatusOK, data)
response.JSON(w, http.StatusCreated, data)
response.WithMeta(w, http.StatusOK, data, &response.Meta{Page: 1, Limit: 20, Total: 100, HasNext: true})

// Error
response.BadRequest(w, "Nama wajib diisi")
response.NotFound(w)
response.Forbidden(w)
response.Unauthorized(w)
response.InternalError(w)
response.Error(w, 422, "VALIDATION_ERROR", "Format tanggal tidak valid")
```

JANGAN pernah menulis `json.NewEncoder(w).Encode(...)` langsung di handler selain untuk kasus yang sangat spesifik.

### Mengambil claims dari context

```go
func (h *SomeHandler) SomeMethod(w http.ResponseWriter, r *http.Request) {
    claims := middleware.GetClaims(r)
    // claims.UserID   — uuid.UUID
    // claims.Role     — model.UserRole ("superadmin", "write", "readonly")
    // claims.BranchID — *uuid.UUID (nil jika superadmin)
    // claims.Name     — string
}
```

### Mengambil URL parameter

```go
import "github.com/go-chi/chi/v5"

id := chi.URLParam(r, "id")
patientID, err := uuid.Parse(id)
if err != nil {
    response.BadRequest(w, "ID tidak valid")
    return
}
```

---

## Database

### Konvensi query

Selalu filter `deleted_at IS NULL` untuk tabel yang punya soft delete (patients, visits, attachments).

```go
// BENAR
"SELECT ... FROM patients WHERE id = $1 AND deleted_at IS NULL"

// SALAH — tidak filter soft delete
"SELECT ... FROM patients WHERE id = $1"
```

Selalu gunakan `pgxpool.Pool`, bukan koneksi langsung. Pool sudah di-setup di `main.go`.

### Soft delete pattern

```go
func (r *SomeRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
    result, err := r.db.Exec(ctx,
        "UPDATE table SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL",
        id,
    )
    if err != nil {
        return fmt.Errorf("soft delete: %w", err)
    }
    if result.RowsAffected() == 0 {
        return fmt.Errorf("record not found")
    }
    return nil
}
```

### Pagination pattern

Semua endpoint list harus support pagination dengan query params `?page=1&limit=20`.

```go
// Di handler
page, _ := strconv.Atoi(r.URL.Query().Get("page"))
limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
if page < 1 { page = 1 }
if limit < 1 || limit > 100 { limit = 20 }
```

---

## API Endpoints — Referensi Lengkap

Semua route sudah didefinisikan di `main.go`. Implementasikan method-method berikut:

### Auth (sudah selesai)
| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/api/v1/auth/login` | `AuthHandler.Login` | Public |
| GET | `/api/v1/auth/me` | `AuthHandler.Me` | Any |
| POST | `/api/v1/auth/logout` | `AuthHandler.Logout` | Any |

### Branches
| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/branches` | `BranchHandler.List` | Any |

Response: array of `{id, name, address, phone, is_active}`

### Patients
| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/patients` | `PatientHandler.List` | Any |
| GET | `/api/v1/patients/{id}` | `PatientHandler.GetByID` | Any |
| POST | `/api/v1/patients` | `PatientHandler.Create` | Write |
| PUT | `/api/v1/patients/{id}` | `PatientHandler.Update` | Write |
| DELETE | `/api/v1/patients/{id}` | `PatientHandler.Delete` | Write |

**GET /patients — query params:**
- `page` int, default 1
- `limit` int, default 20, max 100
- `search` string — cari di nama, no_rm, phone
- `status` string — `new` | `active` | `needs_control`
- `branch_id` uuid — filter per cabang

**POST /patients — request body:**
```json
{
  "branch_id": "uuid",
  "name": "string (required)",
  "nik": "string (optional)",
  "date_of_birth": "2000-01-15 (required, format YYYY-MM-DD)",
  "gender": "male | female (required)",
  "phone": "string (required)",
  "address": "string (optional)",
  "occupation": "string (optional)",
  "allergy_notes": "string (optional)"
}
```
Response: patient object lengkap + `no_rm` yang baru digenerate.

**Logika branch_id pada Create:**
- Jika role `write`: `branch_id` diambil dari `claims.BranchID`, request body boleh tidak menyertakan atau harus match
- Jika role `superadmin`: `branch_id` wajib ada di request body

### Visits
| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/patients/{id}/visits` | `VisitHandler.ListByPatient` | Any |
| GET | `/api/v1/visits/{id}` | `VisitHandler.GetByID` | Any |
| POST | `/api/v1/visits` | `VisitHandler.Create` | Write |
| PUT | `/api/v1/visits/{id}` | `VisitHandler.Update` | Write |
| DELETE | `/api/v1/visits/{id}` | `VisitHandler.Delete` | Write |

**POST /visits — request body:**
```json
{
  "patient_id": "uuid (required)",
  "doctor_id": "uuid (required)",
  "visit_date": "2025-04-30 (required)",
  "chief_complaint": "string (required)",
  "diagnosis": "string (optional)",
  "treatment": "string (optional)",
  "teeth_involved": "string (optional, contoh: 16,17,36)",
  "cost": 350000,
  "next_control_date": "2025-07-30 (optional)",
  "notes": "string (optional)"
}
```
`branch_id` diisi otomatis dari `claims.BranchID` di service layer, bukan dari request body.

### Attachments
| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/patients/{id}/attachments` | `AttachmentHandler.ListByPatient` | Any |
| GET | `/api/v1/attachments/{id}` | `AttachmentHandler.GetByID` | Any |
| GET | `/api/v1/attachments/{id}/download` | `AttachmentHandler.Download` | Any |
| POST | `/api/v1/attachments` | `AttachmentHandler.Upload` | Write |
| DELETE | `/api/v1/attachments/{id}` | `AttachmentHandler.Delete` | Write |

**POST /attachments — multipart/form-data:**
- `visit_id` (form field, uuid)
- `file` (file field, PDF atau image)

**GET /attachments/{id}/download:**
- Stream file langsung ke response
- Set header `Content-Disposition: attachment; filename="nama_asli.pdf"`
- Set header `Content-Type` sesuai `mime_type` dari DB
- JANGAN expose path fisik file ke client — selalu serve via endpoint ini

### Export
| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/export/patients` | `ExportHandler.ExportPatients` | Any |
| GET | `/api/v1/export/visits` | `ExportHandler.ExportVisits` | Any |

**Query params untuk kedua export:**
- `branch_id` uuid (optional)
- `date_from` string YYYY-MM-DD (optional, untuk visits)
- `date_to` string YYYY-MM-DD (optional, untuk visits)

Response: file `.xlsx` langsung (bukan JSON).
Headers: `Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
         `Content-Disposition: attachment; filename="export_pasien_2025-04-30.xlsx"`

### User Management (superadmin only)
| Method | Path | Handler | Auth |
|--------|------|---------|------|
| GET | `/api/v1/users` | `UserMgmtHandler.List` | SuperAdmin |
| POST | `/api/v1/users` | `UserMgmtHandler.Create` | SuperAdmin |
| PUT | `/api/v1/users/{id}` | `UserMgmtHandler.Update` | SuperAdmin |
| DELETE | `/api/v1/users/{id}` | `UserMgmtHandler.Deactivate` | SuperAdmin |

**POST /users — request body:**
```json
{
  "branch_id": "uuid (null untuk superadmin)",
  "name": "string (required)",
  "email": "string (required)",
  "password": "string (required, min 8 karakter)",
  "role": "superadmin | write | readonly (required)"
}
```
Password di-hash dengan bcrypt cost 12 sebelum disimpan. JANGAN simpan plaintext.

---

## pkg/storage — Yang Perlu Dibuat

```go
// pkg/storage/local.go
package storage

type Local struct {
    baseDir     string
    maxSizeMB   int64
}

func NewLocal(baseDir string, maxSizeMB int64) *Local

// Save menyimpan file ke disk
// path: /{patientID}/{visitID}/{uuid}.{ext}
// Returns: storedName (uuid.ext), filePath (relatif), error
func (s *Local) Save(file multipart.File, header *multipart.FileHeader, patientID, visitID string) (storedName, filePath string, err error)

// Delete menghapus file dari disk
func (s *Local) Delete(filePath string) error

// ServeFile membaca file dan menulis ke http.ResponseWriter
func (s *Local) ServeFile(w http.ResponseWriter, filePath, mimeType, originalName string) error
```

Validasi yang harus ada di `Save`:
- Cek ukuran file tidak melebihi `maxSizeMB`
- Cek MIME type hanya `application/pdf`, `image/jpeg`, `image/png`, `image/webp`
- Generate `storedName` dengan `uuid.New().String() + ext`
- Buat folder jika belum ada (`os.MkdirAll`)

---

## pkg/excel — Yang Perlu Dibuat

Gunakan library `github.com/xuri/excelize/v2`.

```go
// pkg/excel/excel.go
package excel

// ExportPatients membuat file xlsx dari list pasien
// Kolom: No. RM, Nama, Tgl Lahir, Usia, JK, Telepon, Alamat, Alergi, Status, Cabang, Tgl Input
func ExportPatients(patients []model.PatientListRow) (*bytes.Buffer, error)

// ExportVisits membuat file xlsx dari list kunjungan
// Kolom: No. RM, Nama Pasien, Tgl Kunjungan, Dokter, Keluhan, Diagnosis, Tindakan, Gigi, Biaya, Cabang
func ExportVisits(visits []VisitExportRow) (*bytes.Buffer, error)
```

Style header: background teal (#0F6E56), text putih, bold. Kolom biaya format Rupiah. Freeze row pertama.

---

## Docker Compose — Yang Perlu Dibuat

File `docker-compose.yml` di root `dental-project/`:

```yaml
services:
  nginx:        # reverse proxy + serve React build
  go-api:       # backend, port 8080 internal
  postgres:     # database, port 5432 internal
  backup-cron:  # pg_dump harian, image: prodrigestivill/postgres-backup-local

volumes:
  pg-data:
  uploads:
  ssl-certs:
  backups:

networks:
  dental-net:
```

Aturan penting:
- Hanya `nginx` yang expose port ke luar (80 dan 443)
- `postgres` dan `go-api` hanya di `dental-net`, tidak expose ke host
- `go-api` depends_on `postgres` dengan healthcheck
- Semua secret dari file `.env` via `env_file: .env`

---

## Dockerfile — Yang Perlu Dibuat

```dockerfile
# Multi-stage build
# Stage 1: builder — golang:1.23-alpine
# Stage 2: runtime — alpine:3.20
# Binary output: /app/dental-api
# Expose: 8080
# Non-root user: appuser
```

Target image size < 20MB.

---

## Hal-hal Penting yang Jangan Sampai Salah

### Keamanan
1. `PasswordHash` di struct `User` punya json tag `json:"-"` — JANGAN ubah ini. Password hash tidak boleh pernah dikirim ke client.
2. `StoredName` dan `FilePath` di struct `Attachment` juga `json:"-"` — path fisik file tidak boleh exposed ke client.
3. Download file selalu lewat endpoint `/attachments/{id}/download` dengan auth check, bukan direct URL ke disk.
4. Semua input dari user harus divalidasi sebelum masuk ke repository.
5. Gunakan parameterized query (`$1, $2, ...`) SELALU — jangan pernah string concatenation untuk SQL.

### No. RM
Format saat ini: `RM-{YYYY}-{4digit}` contoh `RM-2025-0001`.
Ini **mungkin akan berubah** setelah konfirmasi dari client (mereka mungkin punya format sendiri).
Function `generate_no_rm()` ada di PostgreSQL (`migrations/003_patients.sql`), dipanggil dari `PatientRepo.Create`.
Jika format berubah, cukup update function tersebut di migration baru — tidak perlu ubah Go code.

### Role & branch_id
- `superadmin`: `branch_id = NULL`, bisa akses semua cabang
- `write` (dokter): `branch_id` terisi, hanya akses cabangnya
- `readonly` (suster): `branch_id` terisi, hanya baca data cabangnya

Untuk semua query yang melibatkan data per cabang, filter harus dinamis:
```go
if claims.Role != model.RoleSuperAdmin && claims.BranchID != nil {
    // tambahkan filter branch_id = claims.BranchID
}
```

### Soft delete
Tabel `patients`, `visits`, `attachments` punya kolom `deleted_at`.
Semua query SELECT harus filter `WHERE deleted_at IS NULL`.
View `v_patient_list` sudah include filter ini.

Saat delete attachment, selain set `deleted_at` di DB, juga hapus file fisik dari disk via `storage.Delete(filePath)`.

### Timezone
Semua timestamp disimpan sebagai `TIMESTAMPTZ` di PostgreSQL.
DSN sudah set `timezone=Asia/Jakarta` di `config.go`.
Di Go, semua `time.Time` otomatis dalam WIB.

---

## Cara Menjalankan Development

```bash
# 1. Clone dan masuk ke folder
cd dental-project/dental-api

# 2. Copy env dan isi dengan nilai yang benar
cp .env.example .env

# 3. Start semua service via Docker
make up

# 4. Jalankan migration (hanya sekali, atau saat ada migration baru)
make migrate

# 5. Jalankan server Go
make run

# 6. Test health check
curl http://localhost:8080/health
# {"status":"ok"}

# 7. Test login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@klinik.local","password":"Admin@1234"}'
```

Setelah login berhasil, gunakan token dari response untuk request selanjutnya:
```bash
curl -H "Authorization: Bearer {token}" http://localhost:8080/api/v1/auth/me
```

---

## Urutan Prioritas Pengerjaan

Ikuti urutan ini — setiap langkah harus `go build ./...` tanpa error sebelum lanjut ke langkah berikutnya:

**Fase 1 — Fondasi (bisa langsung ditest)**
1. `pkg/storage/local.go`
2. `internal/repository/branch.go` + `internal/service/branch.go` + `internal/handler/branch.go`
3. `internal/service/patient.go` + `internal/handler/patient.go`
4. Verifikasi: `GET /api/v1/patients` dan `POST /api/v1/patients` berjalan

**Fase 2 — Kunjungan & Lampiran**
5. `internal/repository/visit.go` + `internal/service/visit.go` + `internal/handler/visit.go`
6. `internal/repository/attachment.go` + `internal/service/attachment.go` + `internal/handler/attachment.go`
7. Verifikasi: full flow input pasien → input kunjungan → upload lampiran → download

**Fase 3 — Export & User Management**
8. `pkg/excel/excel.go`
9. `internal/service/export.go` + `internal/handler/export.go`
10. `internal/service/user_mgmt.go` + `internal/handler/user_mgmt.go`
11. Verifikasi: export xlsx bisa dibuka di Excel, user bisa dibuat dan login

**Fase 4 — Infrastruktur**
12. `Dockerfile`
13. `docker-compose.yml`
14. Verifikasi: `docker compose up` berjalan penuh, semua endpoint accessible

---

## Pertanyaan yang Masih Pending dari Client

Beberapa hal belum dikonfirmasi client — implementasi sementara sudah menggunakan asumsi yang masuk akal. Catat ini untuk direvisi setelah jawaban client masuk:

| # | Pertanyaan | Asumsi saat ini | File yang perlu diubah jika beda |
|---|-----------|-----------------|----------------------------------|
| 1 | Format No. RM | `RM-{YYYY}-{4digit}` global | `migrations/003_patients.sql` function `generate_no_rm()` |
| 2 | Ada data pasien lama? | Tidak ada, mulai dari nol | Perlu buat script import jika ada |
| 3 | No. RM per cabang atau global? | Global | `migrations/003_patients.sql` sequence |
| 4 | Definisi role "admin" | `superadmin` = kelola user saja, `write` = dokter | `migrations/000_init_enums.sql` |
| 5 | Akses suster: semua cabang atau hanya cabangnya? | Hanya cabangnya | Filter query di semua repository |

---

*Dibuat berdasarkan sesi perencanaan lengkap di Claude.ai — termasuk ERD, arsitektur Docker, breakdown fitur 28 items, dan SPK freelance.*
