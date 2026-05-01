-- ============================================================
-- Migration 002 — Users
-- ============================================================

CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id     UUID        REFERENCES branches(id) ON DELETE SET NULL,
    -- NULL branch_id = superadmin (akses semua cabang)
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(150) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          user_role   NOT NULL DEFAULT 'readonly',
    is_active     BOOLEAN     NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE INDEX idx_users_email     ON users(email);
CREATE INDEX idx_users_branch_id ON users(branch_id);
CREATE INDEX idx_users_role      ON users(role);

-- Seed superadmin default (password: Admin@1234 — WAJIB GANTI SAAT LAUNCH)
-- hash bcrypt dari 'Admin@1234'
INSERT INTO users (name, email, password_hash, role, branch_id) VALUES
    ('Super Admin', 'admin@klinik.local',
     '$2a$12$placeholder_hash_ganti_sebelum_launch_xxxxx',
     'superadmin', NULL);
