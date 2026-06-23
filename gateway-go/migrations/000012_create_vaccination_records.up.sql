CREATE TABLE vaccination_records (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id      BIGINT       NOT NULL,
    jenis_vaksin VARCHAR(100) NOT NULL,
    tanggal      DATE         NOT NULL,
    dosis        VARCHAR(50),
    catatan      TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_vaccination_records_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_vaccination_records_user_id ON vaccination_records (user_id);
