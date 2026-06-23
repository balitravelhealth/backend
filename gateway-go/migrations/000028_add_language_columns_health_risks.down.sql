-- Rollback language support from health_risks

ALTER TABLE health_risks
  DROP CONSTRAINT IF EXISTS health_risks_nama_risiko_en_not_null;

ALTER TABLE health_risks
  DROP COLUMN IF EXISTS nama_risiko_en,
  DROP COLUMN IF EXISTS saran_pencegahan_en,
  DROP COLUMN IF EXISTS rekomendasi_vaksinasi_en;

-- Rename columns back to original names
ALTER TABLE health_risks
  RENAME COLUMN nama_risiko_id TO nama_risiko;

ALTER TABLE health_risks
  RENAME COLUMN saran_pencegahan_id TO saran_pencegahan;

ALTER TABLE health_risks
  RENAME COLUMN rekomendasi_vaksinasi_id TO rekomendasi_vaksinasi;
