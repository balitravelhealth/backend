CREATE TABLE roles (
    role_id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama_role VARCHAR(50) NOT NULL,
    deskripsi TEXT,
    CONSTRAINT roles_nama_unique UNIQUE (nama_role)
);
