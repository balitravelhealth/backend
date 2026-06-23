CREATE TABLE health_profiles (
    id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id        BIGINT       NOT NULL,
    tanggal_lahir  DATE,
    jenis_kelamin  VARCHAR(20),
    tinggi_cm      NUMERIC(5,1),
    berat_kg       NUMERIC(5,1),
    golongan_darah VARCHAR(5),
    riwayat_alergi TEXT,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT health_profiles_user_id_unique UNIQUE (user_id),
    CONSTRAINT fk_health_profiles_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
