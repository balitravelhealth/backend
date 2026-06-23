CREATE TYPE rule_kategori AS ENUM ('pre_travel', 'post_travel');
CREATE TYPE rule_status   AS ENUM ('draft', 'published');

CREATE TABLE expert_rules (
    rule_id    BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama       VARCHAR(255)  NOT NULL,
    premis     JSONB         NOT NULL,
    disease_id BIGINT        NOT NULL,
    bobot_cf   NUMERIC(4,3)  NOT NULL CHECK (bobot_cf >= 0 AND bobot_cf <= 1),
    mb         NUMERIC(4,3)  NOT NULL CHECK (mb >= 0 AND mb <= 1),
    md         NUMERIC(4,3)  NOT NULL CHECK (md >= 0 AND md <= 1),
    kategori   rule_kategori NOT NULL,
    status     rule_status   NOT NULL DEFAULT 'draft',
    created_by BIGINT        NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_expert_rules_disease
        FOREIGN KEY (disease_id) REFERENCES expert_diseases (id) ON DELETE RESTRICT,
    CONSTRAINT fk_expert_rules_creator
        FOREIGN KEY (created_by) REFERENCES users           (id) ON DELETE RESTRICT
);

CREATE INDEX idx_expert_rules_disease_id ON expert_rules (disease_id);
CREATE INDEX idx_expert_rules_status     ON expert_rules (status);
