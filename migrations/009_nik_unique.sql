-- ============================================================
-- Migration 009 — NIK Partial Unique Index
-- NIK unik per pasien aktif (NULL diperbolehkan)
-- ============================================================

-- Partial unique index: unik hanya jika NIK diisi dan pasien belum dihapus
CREATE UNIQUE INDEX idx_patients_nik_unique
    ON patients(nik)
    WHERE nik IS NOT NULL AND deleted_at IS NULL;
