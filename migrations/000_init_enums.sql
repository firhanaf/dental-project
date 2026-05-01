-- ============================================================
-- Migration 000 — Enums & Extensions
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- untuk full-text search nama pasien

CREATE TYPE user_role     AS ENUM ('superadmin', 'write', 'readonly');
CREATE TYPE gender_type   AS ENUM ('male', 'female');
CREATE TYPE patient_status AS ENUM ('new', 'active', 'needs_control');
CREATE TYPE file_type     AS ENUM ('pdf', 'image');
CREATE TYPE reminder_status AS ENUM ('pending', 'dismissed', 'done');
CREATE TYPE export_type   AS ENUM ('patients', 'visits');
CREATE TYPE audit_action  AS ENUM ('create', 'update', 'delete', 'restore');
