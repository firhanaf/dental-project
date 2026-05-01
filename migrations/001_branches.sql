-- ============================================================
-- Migration 001 — Branches
-- ============================================================

CREATE TABLE branches (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL,
    address    TEXT,
    phone      VARCHAR(20),
    is_active  BOOLEAN     NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed 2 cabang awal
INSERT INTO branches (id, name, address, phone) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Cabang Utama',  'Jl. ________________', '021-xxxxxxx'),
    ('22222222-2222-2222-2222-222222222222', 'Cabang 2',      'Jl. ________________', '021-xxxxxxx');
