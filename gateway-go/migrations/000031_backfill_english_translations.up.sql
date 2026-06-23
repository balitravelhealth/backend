-- Backfill English translations
-- This migration creates initial values for English translations
-- They will be properly populated by the backfill script using DeepL

-- For expert_symptoms: copy label_id to label_en as placeholder if label_en is NULL
UPDATE expert_symptoms
SET label_en = label_id
WHERE label_en IS NULL;

-- For expert_diseases: copy nama_id to nama_en as placeholder if nama_en is NULL
UPDATE expert_diseases
SET nama_en = nama_id
WHERE nama_en IS NULL;

-- For health_risks: copy nama_risiko_id to nama_risiko_en as placeholder if nama_risiko_en is NULL
UPDATE health_risks
SET nama_risiko_en = nama_risiko_id
WHERE nama_risiko_en IS NULL;

-- For emergency_guides: copy isi_media_id to isi_media_en as placeholder if isi_media_en is NULL
UPDATE emergency_guides
SET isi_media_en = isi_media_id
WHERE isi_media_en IS NULL;

-- For emergency_guide_flows: copy title_id to title_en and nodes_id to nodes_en as placeholder
UPDATE emergency_guide_flows
SET title_en = title_id
WHERE title_en IS NULL;

UPDATE emergency_guide_flows
SET nodes_en = nodes_id
WHERE nodes_en IS NULL;
