CREATE TABLE aoi_locations (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    facility_id BIGINT       NOT NULL,
    deskripsi   TEXT,
    latitude    NUMERIC(9,6) NOT NULL,
    longitude   NUMERIC(9,6) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_aoi_facility
        FOREIGN KEY (facility_id) REFERENCES medical_facilities (id) ON DELETE CASCADE
);

CREATE INDEX idx_aoi_facility_id ON aoi_locations (facility_id);
