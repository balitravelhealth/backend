CREATE TABLE travelers (
    id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id        BIGINT       NOT NULL,
    nama_lengkap   VARCHAR(255) NOT NULL,
    tanggal_lahir  DATE,
    kontak_darurat TEXT,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT travelers_user_id_unique UNIQUE (user_id),
    CONSTRAINT fk_travelers_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
