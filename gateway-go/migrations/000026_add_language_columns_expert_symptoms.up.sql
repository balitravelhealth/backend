-- Add language support to expert_symptoms
-- expert_symptoms already has label_id and label_en, just add descriptions and updated_at

ALTER TABLE expert_symptoms
  ADD COLUMN IF NOT EXISTS deskripsi_id TEXT,
  ADD COLUMN IF NOT EXISTS deskripsi_en TEXT,
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Ensure both languages are present for label
ALTER TABLE expert_symptoms
  ADD CONSTRAINT expert_symptoms_label_en_not_null CHECK (label_en IS NOT NULL);
