-- Add language support to expert_diseases
-- Rename existing columns to clarify they store Indonesian

ALTER TABLE expert_diseases
  RENAME COLUMN nama TO nama_id;

ALTER TABLE expert_diseases
  RENAME COLUMN deskripsi TO deskripsi_id;

ALTER TABLE expert_diseases
  RENAME COLUMN rekomendasi_default TO rekomendasi_default_id;

-- Add English language columns (initially nullable)
ALTER TABLE expert_diseases
  ADD COLUMN IF NOT EXISTS nama_en VARCHAR(255),
  ADD COLUMN IF NOT EXISTS deskripsi_en TEXT,
  ADD COLUMN IF NOT EXISTS rekomendasi_default_en JSONB;

-- Populate EN columns with ID values as placeholder
UPDATE expert_diseases SET nama_en = nama_id WHERE nama_en IS NULL;

-- Update constraint name
ALTER TABLE expert_diseases
  DROP CONSTRAINT IF EXISTS expert_diseases_nama_unique;

ALTER TABLE expert_diseases
  ADD CONSTRAINT expert_diseases_nama_id_unique UNIQUE (nama_id);

-- Now ensure EN translations exist (after populating)
ALTER TABLE expert_diseases
  ADD CONSTRAINT expert_diseases_nama_en_not_null CHECK (nama_en IS NOT NULL);
