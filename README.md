# Dental Record System — Project Skeleton

## Struktur Folder
migrations/     → SQL files, jalankan urut 000-007
dental-api/     → Go REST API skeleton

## Sudah implementasi penuh:
- Semua migration SQL (enums, tables, triggers, views, indexes)
- Config loader, domain models, JWT manager, response helper
- Auth middleware (Authenticate / RequireWrite / RequireSuperAdmin)
- Auth handler + service (login, me, logout)
- Patient repository (List, GetByID, Create, Update, SoftDelete)
- main.go dengan routing SEMUA endpoint sudah terdefinisi

## Perlu dilanjutkan (pola sama, tinggal copy-adapt):
- handler & service: patient, visit, attachment, export, user_mgmt, branch
- repository: user, branch, visit, attachment
- pkg/storage: local file upload handler
- pkg/excel: export .xlsx dengan excelize

## Quick start:
  cp .env.example .env && nano .env
  make up && make migrate && make run
