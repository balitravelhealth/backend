CREATE TABLE permissions (
    permission_id   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama_permission VARCHAR(100) NOT NULL,
    CONSTRAINT permissions_nama_unique UNIQUE (nama_permission)
);
