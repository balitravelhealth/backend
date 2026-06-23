CREATE TABLE expert_symptoms (
    symptom_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    kode       VARCHAR(50)  NOT NULL,
    label_id   VARCHAR(255) NOT NULL,
    label_en   VARCHAR(255),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT expert_symptoms_kode_unique UNIQUE (kode)
);
