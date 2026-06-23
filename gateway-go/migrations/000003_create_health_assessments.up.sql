CREATE TABLE health_assessments (
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id          BIGINT      NOT NULL,
    symptoms         JSONB       NOT NULL,
    ai_analysis_raw  JSONB,
    diagnosis        TEXT,
    confidence_score NUMERIC(5,4),
    risk_level       VARCHAR(20),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_health_assessments_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_health_assessments_user_id ON health_assessments (user_id);
