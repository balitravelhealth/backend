CREATE TABLE nursing_care_records (
    id                     BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    traveler_id            BIGINT NOT NULL,
    nurse_id               BIGINT      NOT NULL,
    tanggal_kunjungan      TIMESTAMPTZ NOT NULL,
    nursing_assessment     TEXT,
    nursing_diagnosis      TEXT,
    nursing_planning       TEXT,
    nursing_implementation TEXT,
    nursing_evaluation     TEXT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_ncr_traveler
        FOREIGN KEY (traveler_id) REFERENCES travelers (id) ON DELETE CASCADE,
    CONSTRAINT fk_ncr_nurse
        FOREIGN KEY (nurse_id)    REFERENCES nurses    (id) ON DELETE CASCADE
);

CREATE INDEX idx_ncr_traveler_id ON nursing_care_records (traveler_id);
CREATE INDEX idx_ncr_nurse_id    ON nursing_care_records (nurse_id);
