-- ============================================================
-- Migration 008 — Testing Seed Data
-- HANYA untuk environment testing/staging
-- Jangan dijalankan di production yang sudah memiliki data nyata
-- ============================================================

BEGIN;

-- ── 1. Update nama cabang ─────────────────────────────────────
UPDATE branches SET
    name       = 'Klinik Gigi Sehat Utama',
    address    = 'Jl. Sudirman No. 45, Jakarta Pusat 10220',
    phone      = '021-5551234',
    updated_at = NOW()
WHERE id = '11111111-1111-1111-1111-111111111111';

UPDATE branches SET
    name       = 'Klinik Gigi Sehat Selatan',
    address    = 'Jl. TB Simatupang No. 12, Jakarta Selatan 12560',
    phone      = '021-7891234',
    updated_at = NOW()
WHERE id = '22222222-2222-2222-2222-222222222222';

-- ── 2. Fix password superadmin (Admin@1234) ───────────────────
-- crypt() dari pgcrypto menghasilkan bcrypt hash yang kompatibel dengan Go bcrypt
UPDATE users SET
    password_hash = crypt('Admin@1234', gen_salt('bf', 12)),
    updated_at    = NOW()
WHERE email = 'admin@klinik.local';

-- ── 3. Insert users test ──────────────────────────────────────
-- Password semua dokter: Dokter@1234
-- Password semua suster: Suster@1234

-- Dokter Cabang Utama
INSERT INTO users (id, branch_id, name, email, password_hash, role) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
     '11111111-1111-1111-1111-111111111111',
     'drg. Budi Santoso', 'budi@klinik.local',
     crypt('Dokter@1234', gen_salt('bf', 12)), 'write'),

    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
     '11111111-1111-1111-1111-111111111111',
     'drg. Ahmad Fauzi', 'ahmad@klinik.local',
     crypt('Dokter@1234', gen_salt('bf', 12)), 'write');

-- Dokter Cabang Selatan
INSERT INTO users (id, branch_id, name, email, password_hash, role) VALUES
    ('cccccccc-cccc-cccc-cccc-cccccccccccc',
     '22222222-2222-2222-2222-222222222222',
     'drg. Siti Rahayu', 'siti@klinik.local',
     crypt('Dokter@1234', gen_salt('bf', 12)), 'write');

-- Suster / Readonly
INSERT INTO users (id, branch_id, name, email, password_hash, role) VALUES
    ('dddddddd-dddd-dddd-dddd-dddddddddddd',
     '11111111-1111-1111-1111-111111111111',
     'Ani Susanti', 'ani@klinik.local',
     crypt('Suster@1234', gen_salt('bf', 12)), 'readonly'),

    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
     '22222222-2222-2222-2222-222222222222',
     'Rini Permata', 'rini@klinik.local',
     crypt('Suster@1234', gen_salt('bf', 12)), 'readonly');

-- ── 4. Insert patients ────────────────────────────────────────
-- Status awal 'new' — akan diupdate otomatis oleh trigger saat visit diinsert
-- Pasien tanpa visit tetap berstatus 'new'

-- Cabang Utama (5 pasien)
INSERT INTO patients
    (id, branch_id, created_by, no_rm, name, nik, date_of_birth, gender, phone, address, occupation, allergy_notes)
VALUES
    -- p1: tidak ada visit → status tetap 'new'
    ('a0000001-0000-0000-0000-000000000001',
     '11111111-1111-1111-1111-111111111111',
     'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
     generate_no_rm(),
     'Dewi Hartati', '3174056508900001', '1990-08-25', 'female',
     '08123456001', 'Jl. Kebon Jeruk No. 10, Jakarta Barat', 'Karyawan Swasta', NULL),

    -- p2: akan menjadi 'active' (kunjungan terakhir < 90 hari)
    ('a0000001-0000-0000-0000-000000000002',
     '11111111-1111-1111-1111-111111111111',
     'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
     generate_no_rm(),
     'Rudi Hermawan', '3174051204850002', '1985-04-12', 'male',
     '08123456002', 'Jl. Puri Kencana No. 5, Jakarta Barat', 'Wiraswasta', 'Alergi penisilin'),

    -- p3: akan menjadi 'active' (kunjungan terakhir < 90 hari, ada 2 kunjungan)
    ('a0000001-0000-0000-0000-000000000003',
     '11111111-1111-1111-1111-111111111111',
     'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
     generate_no_rm(),
     'Maya Anggraini', '3174052309920003', '1992-09-23', 'female',
     '08123456003', 'Jl. Daan Mogot No. 77, Jakarta Barat', 'Ibu Rumah Tangga', NULL),

    -- p4: akan menjadi 'needs_control' (kunjungan terakhir > 90 hari)
    ('a0000001-0000-0000-0000-000000000004',
     '11111111-1111-1111-1111-111111111111',
     'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
     generate_no_rm(),
     'Santoso Widjaja', '3174050307780004', '1978-07-03', 'male',
     '08123456004', 'Jl. Palmerah Barat No. 20, Jakarta Barat', 'PNS', NULL),

    -- p5: akan menjadi 'needs_control' (kunjungan terakhir > 90 hari)
    ('a0000001-0000-0000-0000-000000000005',
     '11111111-1111-1111-1111-111111111111',
     'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
     generate_no_rm(),
     'Fitriani Sari', '3174052801950005', '1995-01-28', 'female',
     '08123456005', NULL, 'Mahasiswi', NULL);

-- Cabang Selatan (5 pasien)
INSERT INTO patients
    (id, branch_id, created_by, no_rm, name, nik, date_of_birth, gender, phone, address, occupation, allergy_notes)
VALUES
    -- p6: tidak ada visit → status tetap 'new'
    ('b0000002-0000-0000-0000-000000000006',
     '22222222-2222-2222-2222-222222222222',
     'cccccccc-cccc-cccc-cccc-cccccccccccc',
     generate_no_rm(),
     'Hendra Kusuma', '3171020911800006', '1980-11-09', 'male',
     '08198765001', 'Jl. Fatmawati No. 33, Jakarta Selatan', 'Dokter', NULL),

    -- p7: akan menjadi 'active'
    ('b0000002-0000-0000-0000-000000000007',
     '22222222-2222-2222-2222-222222222222',
     'cccccccc-cccc-cccc-cccc-cccccccccccc',
     generate_no_rm(),
     'Lestari Putri', '3171021507870007', '1987-07-15', 'female',
     '08198765002', 'Jl. Sisingamangaraja No. 8, Jakarta Selatan', 'Guru', 'Alergi latex'),

    -- p8: akan menjadi 'active'
    ('b0000002-0000-0000-0000-000000000008',
     '22222222-2222-2222-2222-222222222222',
     'cccccccc-cccc-cccc-cccc-cccccccccccc',
     generate_no_rm(),
     'Agus Prasetyo', '3171021802930008', '1993-02-18', 'male',
     '08198765003', NULL, 'Karyawan Swasta', NULL),

    -- p9: akan menjadi 'needs_control'
    ('b0000002-0000-0000-0000-000000000009',
     '22222222-2222-2222-2222-222222222222',
     'cccccccc-cccc-cccc-cccc-cccccccccccc',
     generate_no_rm(),
     'Nurhayati Hakim', '3171021310760009', '1976-10-13', 'female',
     '08198765004', 'Jl. Wolter Monginsidi No. 15, Jakarta Selatan', 'Wiraswasta', NULL),

    -- p10: tidak ada visit → status tetap 'new'
    ('b0000002-0000-0000-0000-000000000010',
     '22222222-2222-2222-2222-222222222222',
     'cccccccc-cccc-cccc-cccc-cccccccccccc',
     generate_no_rm(),
     'Bayu Setiawan', '3171020506010010', '2001-06-05', 'male',
     '08198765005', NULL, 'Pelajar', NULL);

-- ── 5. Insert visits ──────────────────────────────────────────
-- Trigger trg_update_patient_status akan otomatis update status pasien
-- active         = MAX(visit_date) >= NOW() - 90 hari  (>= 2026-01-31)
-- needs_control  = MAX(visit_date) <  NOW() - 90 hari  (< 2026-01-31)

-- [p2] Rudi Hermawan — 1 kunjungan terbaru → active
INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost, next_control_date)
VALUES (
    'c1000001-0000-0000-0000-000000000001',
    'a0000001-0000-0000-0000-000000000002',
    '11111111-1111-1111-1111-111111111111',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    '2026-04-10',
    'Sakit gigi geraham kanan bawah', 'Karies profunda gigi 46',
    'Penambalan GIC diikuti komposit', '46',
    350000, '2026-07-10'
);

-- [p3] Maya Anggraini — 2 kunjungan, terakhir terbaru → active
INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost)
VALUES (
    'c1000001-0000-0000-0000-000000000002',
    'a0000001-0000-0000-0000-000000000003',
    '11111111-1111-1111-1111-111111111111',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '2025-11-15',
    'Pembersihan karang gigi rutin', NULL,
    'Scaling ultrasonik seluruh rahang', NULL,
    200000
);

INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost, next_control_date)
VALUES (
    'c1000001-0000-0000-0000-000000000003',
    'a0000001-0000-0000-0000-000000000003',
    '11111111-1111-1111-1111-111111111111',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '2026-04-22',
    'Kontrol pasca scaling, gusi bengkak berkurang', 'Gingiva sehat',
    'Polishing, instruksi oral hygiene', NULL,
    150000, '2026-10-22'
);

-- [p4] Santoso Widjaja — kunjungan lama → needs_control
INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost)
VALUES (
    'c1000001-0000-0000-0000-000000000004',
    'a0000001-0000-0000-0000-000000000004',
    '11111111-1111-1111-1111-111111111111',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    '2025-09-20',
    'Sakit gigi depan atas, sensitif terhadap dingin',
    'Karies media gigi 11 dan 12', 'Penambalan komposit A2',
    '11,12', 500000
);

-- [p5] Fitriani Sari — kunjungan lama → needs_control
INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost)
VALUES (
    'c1000001-0000-0000-0000-000000000005',
    'a0000001-0000-0000-0000-000000000005',
    '11111111-1111-1111-1111-111111111111',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    '2025-10-05',
    'Gigi bungsu kiri bawah nyeri dan bengkak',
    'Perikoronitis, impaksi mesioangular gigi 38', 'Ekstraksi gigi 38',
    '38', 800000
);

-- [p7] Lestari Putri (cabang selatan) — kunjungan terbaru → active
INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost, next_control_date)
VALUES (
    'c2000002-0000-0000-0000-000000000006',
    'b0000002-0000-0000-0000-000000000007',
    '22222222-2222-2222-2222-222222222222',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    '2026-04-05',
    'Mahkota gigi patah, ingin dibuatkan yang baru',
    'Mahkota lama rusak gigi 36', 'Preparasi mahkota porselen',
    '36', 2500000, '2026-05-15'
);

-- [p8] Agus Prasetyo (cabang selatan) — kunjungan terbaru → active
INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost)
VALUES (
    'c2000002-0000-0000-0000-000000000007',
    'b0000002-0000-0000-0000-000000000008',
    '22222222-2222-2222-2222-222222222222',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    '2026-03-18',
    'Sakit gigi geraham kiri atas sejak 3 hari',
    'Pulpitis reversibel gigi 26', 'Penambalan sementara, evaluasi 1 minggu',
    '26', 250000
);

-- [p9] Nurhayati Hakim (cabang selatan) — kunjungan lama → needs_control
INSERT INTO visits
    (id, patient_id, branch_id, doctor_id, created_by, visit_date,
     chief_complaint, diagnosis, treatment, teeth_involved, cost)
VALUES (
    'c2000002-0000-0000-0000-000000000008',
    'b0000002-0000-0000-0000-000000000009',
    '22222222-2222-2222-2222-222222222222',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    '2025-08-10',
    'Gigi hilang beberapa, ingin dibuatkan gigi palsu',
    'Edentulous parsial rahang bawah', 'Protesa lepasan akrilik 3 elemen',
    '34,35,36', 4500000
);

-- ── 6. Sesuaikan sequence no_rm ───────────────────────────────
-- 10 pasien sudah diinsert, set ke 10 agar nextval() berikutnya = 11
SELECT setval('no_rm_seq', 10, true);

COMMIT;

-- ── Ringkasan akun untuk testing ─────────────────────────────
-- Superadmin   : admin@klinik.local    / Admin@1234
-- Dokter utama : budi@klinik.local     / Dokter@1234  (write, cabang utama)
-- Dokter utama : ahmad@klinik.local    / Dokter@1234  (write, cabang utama)
-- Dokter selatan: siti@klinik.local    / Dokter@1234  (write, cabang selatan)
-- Suster utama : ani@klinik.local      / Suster@1234  (readonly, cabang utama)
-- Suster selatan: rini@klinik.local    / Suster@1234  (readonly, cabang selatan)
--
-- Status pasien setelah seed:
-- Cabang Utama : Dewi (new), Rudi (active), Maya (active, 2 visit),
--                Santoso (needs_control), Fitriani (needs_control)
-- Cabang Selatan: Hendra (new), Lestari (active), Agus (active),
--                Nurhayati (needs_control), Bayu (new)
