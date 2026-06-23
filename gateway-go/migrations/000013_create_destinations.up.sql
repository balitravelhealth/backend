CREATE TABLE destinations (
    destination_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama_daerah    VARCHAR(200) NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
