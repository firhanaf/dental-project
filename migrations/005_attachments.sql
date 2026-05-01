-- ============================================================
-- Migration 005 — Attachments
-- ============================================================

CREATE TABLE attachments (
    id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    visit_id      UUID        NOT NULL REFERENCES visits(id) ON DELETE RESTRICT,
    uploaded_by   UUID        NOT NULL REFERENCES users(id)  ON DELETE RESTRICT,

    -- original_name: nama asli file yang diupload user (untuk tampilan UI)
    -- stored_name:   nama file di disk (UUID-based, anti-collision & path traversal)
    original_name VARCHAR(255) NOT NULL,
    stored_name   VARCHAR(255) NOT NULL,
    -- file_path: relatif dari UPLOAD_DIR, contoh: "/{patient_id}/{visit_id}/abc123.pdf"
    file_path     TEXT         NOT NULL,

    file_type     file_type    NOT NULL,  -- 'pdf' | 'image'
    mime_type     VARCHAR(100) NOT NULL,  -- 'application/pdf', 'image/jpeg', dst
    size_bytes    BIGINT       NOT NULL,

    -- Soft delete: file fisik di disk dihapus saat deleted_at di-set
    deleted_at    TIMESTAMPTZ,

    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT attachments_stored_name_unique UNIQUE (stored_name)
);

CREATE INDEX idx_attachments_visit_id    ON attachments(visit_id);
CREATE INDEX idx_attachments_deleted_at  ON attachments(deleted_at) WHERE deleted_at IS NULL;
