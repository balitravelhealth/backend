-- Rollback language support from expert_symptoms

ALTER TABLE expert_symptoms
  DROP CONSTRAINT IF EXISTS expert_symptoms_label_en_not_null;

ALTER TABLE expert_symptoms
  DROP COLUMN IF EXISTS deskripsi_id,
  DROP COLUMN IF EXISTS deskripsi_en,
  DROP COLUMN IF EXISTS updated_at;
