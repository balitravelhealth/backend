-- Rollback language support from emergency_guides

ALTER TABLE emergency_guides
  DROP CONSTRAINT IF EXISTS emergency_guides_isi_media_en_not_null;

ALTER TABLE emergency_guides
  DROP COLUMN IF EXISTS isi_media_en;

-- Rename column back to original name
ALTER TABLE emergency_guides
  RENAME COLUMN isi_media_id TO isi_media;
