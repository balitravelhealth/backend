CREATE TABLE nurses (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id       BIGINT       NOT NULL,
    nama_lengkap  VARCHAR(255) NOT NULL,
    nomor_lisensi VARCHAR(100) NOT NULL,
    sertifikasi   TEXT,
    aktif         BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT nurses_user_id_unique UNIQUE (user_id),
    CONSTRAINT nurses_lisensi_unique UNIQUE (nomor_lisensi),
    CONSTRAINT fk_nurses_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
