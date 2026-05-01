-- ============================================================
-- Migration 003 — Patients
-- ============================================================

-- Sequence No. RM global (RM-2025-0001)
-- Dibuat sebagai sequence PostgreSQL agar atomic & tidak bisa duplikat
CREATE SEQUENCE no_rm_seq START 1 INCREMENT 1;

CREATE TABLE patients (
    id            UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id     UUID           NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    created_by    UUID           NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    -- No. RM: auto-generate via function, format RM-{YYYY}-{4digit}
    -- CATATAN: jika client punya format sendiri, ubah fungsi generate_no_rm()
    no_rm         VARCHAR(20)    NOT NULL,

    name          VARCHAR(150)   NOT NULL,
    nik           VARCHAR(16),
    date_of_birth DATE           NOT NULL,
    gender        gender_type    NOT NULL,
    phone         VARCHAR(20)    NOT NULL,
    address       TEXT,
    occupation    VARCHAR(100),
    allergy_notes TEXT,

    -- Status otomatis diupdate oleh fungsi update_patient_status()
    -- 'new'          = belum punya kunjungan
    -- 'active'       = punya kunjungan, kunjungan terakhir < 90 hari
    -- 'needs_control' = kunjungan terakhir >= 90 hari
    status        patient_status NOT NULL DEFAULT 'new',

    -- Soft delete — data tidak benar-benar hilang
    deleted_at    TIMESTAMPTZ,

    created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT patients_no_rm_unique UNIQUE (no_rm)
);

-- Index untuk search performa tinggi
CREATE INDEX idx_patients_branch_id  ON patients(branch_id);
CREATE INDEX idx_patients_no_rm      ON patients(no_rm);
CREATE INDEX idx_patients_phone      ON patients(phone);
CREATE INDEX idx_patients_status     ON patients(status);
CREATE INDEX idx_patients_deleted_at ON patients(deleted_at) WHERE deleted_at IS NULL;

-- GIN index untuk search nama (partial match, case-insensitive)
CREATE INDEX idx_patients_name_trgm ON patients USING GIN (name gin_trgm_ops);

-- ── FUNCTION: generate No. RM ──────────────────────────────
-- Dipanggil saat INSERT pasien baru, bukan DEFAULT karena butuh tahun berjalan
CREATE OR REPLACE FUNCTION generate_no_rm()
RETURNS VARCHAR(20) AS $$
DECLARE
    seq_val  BIGINT;
    year_str TEXT;
BEGIN
    seq_val  := nextval('no_rm_seq');
    year_str := TO_CHAR(NOW(), 'YYYY');
    -- Format: RM-2025-0001
    -- Jika client punya format lain, ubah baris di bawah ini saja
    RETURN 'RM-' || year_str || '-' || LPAD(seq_val::TEXT, 4, '0');
END;
$$ LANGUAGE plpgsql;

-- ── FUNCTION: update status pasien otomatis ──────────────────
CREATE OR REPLACE FUNCTION update_patient_status()
RETURNS TRIGGER AS $$
BEGIN
    -- Dipanggil setelah INSERT/UPDATE pada tabel visits
    UPDATE patients SET
        status = CASE
            WHEN (
                SELECT MAX(visit_date) FROM visits
                WHERE patient_id = NEW.patient_id AND deleted_at IS NULL
            ) >= NOW() - INTERVAL '90 days'
            THEN 'active'::patient_status
            ELSE 'needs_control'::patient_status
        END,
        updated_at = NOW()
    WHERE id = NEW.patient_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
