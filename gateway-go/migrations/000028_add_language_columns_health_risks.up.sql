-- Add language support to health_risks
-- Rename existing columns to clarify they store Indonesian

ALTER TABLE health_risks
  RENAME COLUMN nama_risiko TO nama_risiko_id;

ALTER TABLE health_risks
  RENAME COLUMN saran_pencegahan TO saran_pencegahan_id;

ALTER TABLE health_risks
  RENAME COLUMN rekomendasi_vaksinasi TO rekomendasi_vaksinasi_id;

-- Add English language columns (initially nullable)
ALTER TABLE health_risks
  ADD COLUMN IF NOT EXISTS nama_risiko_en VARCHAR(200),
  ADD COLUMN IF NOT EXISTS saran_pencegahan_en TEXT,
  ADD COLUMN IF NOT EXISTS rekomendasi_vaksinasi_en TEXT;

-- Populate EN columns with ID values as placeholder
UPDATE health_risks SET nama_risiko_en = nama_risiko_id WHERE nama_risiko_en IS NULL;

-- Ensure EN translations exist (after populating)
ALTER TABLE health_risks
  ADD CONSTRAINT health_risks_nama_risiko_en_not_null CHECK (nama_risiko_en IS NOT NULL);
