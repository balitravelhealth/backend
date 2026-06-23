CREATE TABLE health_risks (
    id                    BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    destination_id        BIGINT       NOT NULL,
    nama_risiko           VARCHAR(200) NOT NULL,
    saran_pencegahan      TEXT,
    rekomendasi_vaksinasi TEXT,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_health_risks_destination
        FOREIGN KEY (destination_id) REFERENCES destinations (destination_id) ON DELETE CASCADE
);

CREATE INDEX idx_health_risks_destination_id ON health_risks (destination_id);
