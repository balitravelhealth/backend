CREATE TABLE emergency_guides (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    kategori   VARCHAR(100) NOT NULL,
    langkah    INTEGER      NOT NULL,
    isi_media  JSONB        NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT emergency_guides_kategori_langkah_unique UNIQUE (kategori, langkah)
);

CREATE INDEX idx_emergency_guides_kategori ON emergency_guides (kategori);
