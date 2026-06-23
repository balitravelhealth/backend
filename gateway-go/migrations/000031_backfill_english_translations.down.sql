-- Rollback backfill of English translations
-- This just clears the placeholder values

UPDATE expert_symptoms SET label_en = NULL;
UPDATE expert_symptoms SET deskripsi_en = NULL;

UPDATE expert_diseases SET nama_en = NULL;
UPDATE expert_diseases SET deskripsi_en = NULL;
UPDATE expert_diseases SET rekomendasi_default_en = NULL;

UPDATE health_risks SET nama_risiko_en = NULL;
UPDATE health_risks SET saran_pencegahan_en = NULL;
UPDATE health_risks SET rekomendasi_vaksinasi_en = NULL;

UPDATE emergency_guides SET isi_media_en = NULL;

UPDATE emergency_guide_flows SET title_en = NULL;
UPDATE emergency_guide_flows SET nodes_en = NULL;
