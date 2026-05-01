-- ============================================================
-- Migration 006 — Full Package Tables
-- (audit_logs, control_reminders, export_logs)
-- Untuk Paket MVP: jalankan migration ini tapi tidak expose ke API
-- Untuk Paket Full: aktifkan endpoint yang menggunakan tabel ini
-- ============================================================

-- ── AUDIT LOGS ─────────────────────────────────────────────
CREATE TABLE audit_logs (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    action      audit_action NOT NULL,
    table_name  VARCHAR(50)  NOT NULL,
    record_id   UUID         NOT NULL,
    old_data    JSONB,        -- snapshot data sebelum perubahan
    new_data    JSONB,        -- snapshot data sesudah perubahan
    ip_address  VARCHAR(45),  -- IPv4 atau IPv6
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_user_id    ON audit_logs(user_id);
CREATE INDEX idx_audit_table_rec  ON audit_logs(table_name, record_id);
CREATE INDEX idx_audit_created_at ON audit_logs(created_at DESC);

-- ── CONTROL REMINDERS ──────────────────────────────────────
-- Di-generate otomatis saat next_control_date diisi pada tabel visits
CREATE TABLE control_reminders (
    id          UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    visit_id    UUID            NOT NULL REFERENCES visits(id)    ON DELETE CASCADE,
    patient_id  UUID            NOT NULL REFERENCES patients(id)  ON DELETE CASCADE,
    due_date    DATE            NOT NULL,
    status      reminder_status NOT NULL DEFAULT 'pending',
    dismissed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reminders_due_date   ON control_reminders(due_date) WHERE status = 'pending';
CREATE INDEX idx_reminders_patient_id ON control_reminders(patient_id);

-- Trigger: buat reminder otomatis saat next_control_date diisi
CREATE OR REPLACE FUNCTION create_control_reminder()
RETURNS TRIGGER AS $$
BEGIN
    -- Hanya buat reminder jika next_control_date baru diisi atau berubah
    IF NEW.next_control_date IS NOT NULL AND
       (OLD.next_control_date IS NULL OR OLD.next_control_date <> NEW.next_control_date) THEN

        -- Hapus reminder lama untuk kunjungan ini jika ada
        DELETE FROM control_reminders WHERE visit_id = NEW.id;

        -- Buat reminder baru
        INSERT INTO control_reminders (visit_id, patient_id, due_date)
        VALUES (NEW.id, NEW.patient_id, NEW.next_control_date);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_control_reminder
    AFTER INSERT OR UPDATE ON visits
    FOR EACH ROW
    EXECUTE FUNCTION create_control_reminder();

-- ── EXPORT LOGS ────────────────────────────────────────────
CREATE TABLE export_logs (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    exported_by UUID        NOT NULL REFERENCES users(id)     ON DELETE RESTRICT,
    branch_id   UUID        REFERENCES branches(id) ON DELETE SET NULL,
    export_type export_type NOT NULL,
    date_from   DATE,
    date_to     DATE,
    row_count   INTEGER     NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_export_logs_user ON export_logs(exported_by);
CREATE INDEX idx_export_logs_date ON export_logs(created_at DESC);
