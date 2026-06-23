-- Add language support to emergency_guides
-- Rename existing column to clarify it stores Indonesian

ALTER TABLE emergency_guides
  RENAME COLUMN isi_media TO isi_media_id;

-- Add English language column (initially nullable)
ALTER TABLE emergency_guides
  ADD COLUMN IF NOT EXISTS isi_media_en JSONB;

-- Populate EN column with ID values as placeholder
UPDATE emergency_guides SET isi_media_en = isi_media_id WHERE isi_media_en IS NULL;

-- Ensure EN content exists (after populating)
ALTER TABLE emergency_guides
  ADD CONSTRAINT emergency_guides_isi_media_en_not_null CHECK (isi_media_en IS NOT NULL);
