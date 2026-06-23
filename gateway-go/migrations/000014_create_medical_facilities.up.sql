CREATE TABLE medical_facilities (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    destination_id  BIGINT       NOT NULL,
    nama            VARCHAR(255) NOT NULL,
    kategori        VARCHAR(100),
    alamat          TEXT,
    latitude        NUMERIC(9,6),
    longitude       NUMERIC(9,6),
    kontak          VARCHAR(100),
    jam_operasional TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_facilities_destination
        FOREIGN KEY (destination_id) REFERENCES destinations (destination_id) ON DELETE CASCADE
);

CREATE INDEX idx_facilities_destination_id ON medical_facilities (destination_id);
CREATE INDEX idx_facilities_koordinat      ON medical_facilities (latitude, longitude);
