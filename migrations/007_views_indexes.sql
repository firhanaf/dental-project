-- ============================================================
-- Migration 007 — Useful Views & Composite Indexes
-- ============================================================

-- ── VIEW: patient list (dipakai di halaman daftar pasien) ───
-- Menggabungkan info pasien + kunjungan terakhir + total kunjungan
CREATE VIEW v_patient_list AS
SELECT
    p.id,
    p.branch_id,
    b.name               AS branch_name,
    p.no_rm,
    p.name,
    p.date_of_birth,
    EXTRACT(YEAR FROM AGE(p.date_of_birth))::INT AS age,
    p.gender,
    p.phone,
    p.allergy_notes,
    p.status,
    p.created_at,
    p.updated_at,
    -- Kunjungan terakhir
    lv.visit_date        AS last_visit_date,
    lv.diagnosis         AS last_diagnosis,
    lv.doctor_name       AS last_doctor,
    -- Statistik
    COALESCE(vs.total_visits, 0)  AS total_visits,
    COALESCE(vs.total_cost, 0)    AS total_cost
FROM patients p
JOIN branches b ON b.id = p.branch_id
-- Subquery kunjungan terakhir
LEFT JOIN LATERAL (
    SELECT
        v.visit_date,
        v.diagnosis,
        u.name AS doctor_name
    FROM visits v
    JOIN users u ON u.id = v.doctor_id
    WHERE v.patient_id = p.id AND v.deleted_at IS NULL
    ORDER BY v.visit_date DESC
    LIMIT 1
) lv ON true
-- Subquery statistik total
LEFT JOIN LATERAL (
    SELECT
        COUNT(*)       AS total_visits,
        SUM(v.cost)    AS total_cost
    FROM visits v
    WHERE v.patient_id = p.id AND v.deleted_at IS NULL
) vs ON true
WHERE p.deleted_at IS NULL;

-- ── VIEW: reminders yang aktif (untuk dashboard notif) ──────
CREATE VIEW v_pending_reminders AS
SELECT
    cr.id,
    cr.due_date,
    cr.status,
    p.id         AS patient_id,
    p.no_rm,
    p.name       AS patient_name,
    p.phone      AS patient_phone,
    b.name       AS branch_name,
    CURRENT_DATE - cr.due_date AS days_overdue
FROM control_reminders cr
JOIN patients p ON p.id = cr.patient_id
JOIN branches b ON b.id = p.branch_id
WHERE cr.status = 'pending'
  AND p.deleted_at IS NULL
ORDER BY cr.due_date ASC;

-- ── Composite indexes untuk query umum ─────────────────────

-- Search pasien: filter branch + status + deleted
CREATE INDEX idx_patients_branch_status
    ON patients(branch_id, status)
    WHERE deleted_at IS NULL;

-- Visits per patient diurut tanggal (untuk riwayat kunjungan)
CREATE INDEX idx_visits_patient_date
    ON visits(patient_id, visit_date DESC)
    WHERE deleted_at IS NULL;

-- Export: filter visit by branch + date range
CREATE INDEX idx_visits_branch_date
    ON visits(branch_id, visit_date)
    WHERE deleted_at IS NULL;
