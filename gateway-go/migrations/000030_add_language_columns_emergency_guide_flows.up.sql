-- Add language support to emergency_guide_flows
-- Rename existing columns to clarify they store Indonesian

ALTER TABLE emergency_guide_flows
  RENAME COLUMN title TO title_id;

ALTER TABLE emergency_guide_flows
  RENAME COLUMN nodes TO nodes_id;

-- Add English language columns (initially nullable)
ALTER TABLE emergency_guide_flows
  ADD COLUMN IF NOT EXISTS title_en VARCHAR(200),
  ADD COLUMN IF NOT EXISTS nodes_en JSONB;

-- Populate EN columns with ID values as placeholder
UPDATE emergency_guide_flows SET title_en = title_id WHERE title_en IS NULL;
UPDATE emergency_guide_flows SET nodes_en = nodes_id WHERE nodes_en IS NULL;

-- Ensure EN translations exist (after populating)
ALTER TABLE emergency_guide_flows
  ADD CONSTRAINT emergency_guide_flows_title_en_not_null CHECK (title_en IS NOT NULL);

ALTER TABLE emergency_guide_flows
  ADD CONSTRAINT emergency_guide_flows_nodes_en_not_null CHECK (nodes_en IS NOT NULL);
