-- Rollback language support from expert_diseases

ALTER TABLE expert_diseases
  DROP CONSTRAINT IF EXISTS expert_diseases_nama_en_not_null;

ALTER TABLE expert_diseases
  DROP CONSTRAINT IF EXISTS expert_diseases_nama_id_unique;

ALTER TABLE expert_diseases
  DROP COLUMN IF EXISTS nama_en,
  DROP COLUMN IF EXISTS deskripsi_en,
  DROP COLUMN IF EXISTS rekomendasi_default_en;

-- Rename columns back to original names
ALTER TABLE expert_diseases
  RENAME COLUMN nama_id TO nama;

ALTER TABLE expert_diseases
  RENAME COLUMN deskripsi_id TO deskripsi;

ALTER TABLE expert_diseases
  RENAME COLUMN rekomendasi_default_id TO rekomendasi_default;

-- Restore original constraint
ALTER TABLE expert_diseases
  ADD CONSTRAINT expert_diseases_nama_unique UNIQUE (nama);
