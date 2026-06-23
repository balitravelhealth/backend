-- Remove seeded rules, diseases, symptoms, and system seed user
DELETE FROM expert_rules   WHERE created_by IN (SELECT id FROM users WHERE email = 'system.seed@balihealth.internal');
DELETE FROM expert_diseases WHERE nama IN (
    'Bali Belly (Diare Wisatawan)',
    'Demam Berdarah Dengue (DBD)',
    'Rabies',
    'Hepatitis A',
    'Demam Tifoid (Tifus)',
    'Japanese Encephalitis',
    'Heat-Related Illness (Penyakit Terkait Panas)',
    'Infeksi Luka / Gangguan Kulit',
    'Tuberkulosis (TBC) — Screening Post-Travel',
    'Malaria — Screening Post-Travel'
);
DELETE FROM expert_symptoms WHERE kode LIKE 'S_%';
DELETE FROM users WHERE email = 'system.seed@balihealth.internal';
