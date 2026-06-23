CREATE TABLE expert_diseases (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama                VARCHAR(255) NOT NULL,
    deskripsi           TEXT,
    rekomendasi_default JSONB,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT expert_diseases_nama_unique UNIQUE (nama)
);
