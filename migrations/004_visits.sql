-- ============================================================
-- Migration 004 — Visits
-- ============================================================

CREATE TABLE visits (
    id               UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id       UUID           NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    branch_id        UUID           NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    doctor_id        UUID           NOT NULL REFERENCES users(id)    ON DELETE RESTRICT,
    created_by       UUID           NOT NULL REFERENCES users(id)    ON DELETE RESTRICT,

    visit_date       DATE           NOT NULL,
    chief_complaint  TEXT           NOT NULL,
    diagnosis        TEXT,
    treatment        TEXT,
    teeth_involved   VARCHAR(100),   -- contoh: "16,17,36" atau "semua"
    cost             NUMERIC(12,2)  NOT NULL DEFAULT 0,
    next_control_date DATE,          -- trigger notif jika lewat
    notes            TEXT,

    -- Soft delete
    deleted_at       TIMESTAMPTZ,

    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_visits_patient_id      ON visits(patient_id);
CREATE INDEX idx_visits_branch_id       ON visits(branch_id);
CREATE INDEX idx_visits_doctor_id       ON visits(doctor_id);
CREATE INDEX idx_visits_visit_date      ON visits(visit_date DESC);
CREATE INDEX idx_visits_deleted_at      ON visits(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_visits_next_control    ON visits(next_control_date) WHERE next_control_date IS NOT NULL AND deleted_at IS NULL;

-- Trigger: update status pasien setiap kali ada kunjungan baru/diubah
CREATE TRIGGER trg_update_patient_status
    AFTER INSERT OR UPDATE ON visits
    FOR EACH ROW
    WHEN (NEW.deleted_at IS NULL)
    EXECUTE FUNCTION update_patient_status();
