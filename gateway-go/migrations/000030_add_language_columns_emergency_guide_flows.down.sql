-- Rollback language support from emergency_guide_flows

ALTER TABLE emergency_guide_flows
  DROP CONSTRAINT IF EXISTS emergency_guide_flows_title_en_not_null;

ALTER TABLE emergency_guide_flows
  DROP CONSTRAINT IF EXISTS emergency_guide_flows_nodes_en_not_null;

ALTER TABLE emergency_guide_flows
  DROP COLUMN IF EXISTS title_en,
  DROP COLUMN IF EXISTS nodes_en;

-- Rename columns back to original names
ALTER TABLE emergency_guide_flows
  RENAME COLUMN title_id TO title;

ALTER TABLE emergency_guide_flows
  RENAME COLUMN nodes_id TO nodes;
