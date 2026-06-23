-- ── PY-4 / PY-4b: Seed expert knowledge base ────────────────────────────────
-- ⚠️  SEMUA nilai MB, MD, dan bobot_cf di bawah adalah PLACEHOLDER.
--    Wajib diverifikasi oleh narasumber perawat sebelum sistem digunakan
--    untuk diagnosis nyata (lihat CLAUDE.md §2 dan PY-17).
-- ─────────────────────────────────────────────────────────────────────────────

-- ── 1. Master Gejala (expert_symptoms) ───────────────────────────────────────
INSERT INTO expert_symptoms (kode, label_id, label_en) VALUES
  -- Gastrointestinal
  ('S_DIARE',         'Diare (BAB cair > 3x/hari)',            'Diarrhea (>3 loose stools/day)'),
  ('S_MUAL',          'Mual',                                  'Nausea'),
  ('S_MUNTAH',        'Muntah',                                'Vomiting'),
  ('S_KRAM_PERUT',    'Kram / nyeri perut',                   'Abdominal cramps/pain'),
  ('S_NYERI_RUQ',     'Nyeri perut kanan atas',               'Right upper quadrant pain'),
  ('S_SEMBELIT',      'Sembelit',                              'Constipation'),

  -- Demam & sistemik
  ('S_DEMAM_RINGAN',  'Demam ringan (37.5–38°C)',             'Low-grade fever (37.5–38°C)'),
  ('S_DEMAM_TINGGI',  'Demam tinggi mendadak (≥38.5°C)',      'Sudden high fever (≥38.5°C)'),
  ('S_DEMAM_PERIODIK','Demam periodik / berulang',             'Periodic / recurrent fever'),
  ('S_MENGGIGIL',     'Menggigil',                             'Chills / rigors'),
  ('S_BERKERINGAT',   'Berkeringat berlebihan',                'Excessive sweating'),
  ('S_KELEMAHAN',     'Badan lemas / kelelahan umum',          'Weakness / general fatigue'),
  ('S_BB_TURUN',      'Penurunan berat badan tidak disengaja', 'Unintentional weight loss'),

  -- Kepala & neurologis
  ('S_NYERI_KEPALA',  'Nyeri kepala',                          'Headache'),
  ('S_NYERI_MATA',    'Nyeri di belakang mata (retro-orbital)','Retro-orbital pain'),
  ('S_KAKU_LEHER',    'Kaku kuduk / leher',                   'Neck stiffness'),
  ('S_KEJANG',        'Kejang',                                'Seizure'),
  ('S_BINGUNG',       'Kebingungan / penurunan kesadaran',     'Confusion / altered consciousness'),
  ('S_HIDROFOBIA',    'Hidrofobia / aerofobia',                'Hydrophobia / aerophobia'),
  ('S_KESEMUTAN_GIGITAN', 'Kesemutan / nyeri di area gigitan', 'Tingling/pain at bite site'),

  -- Kulit & muskuloskeletal
  ('S_RUAM',          'Ruam kulit',                            'Skin rash'),
  ('S_RUAM_MERAH',    'Bercak merah (rose spots)',             'Rose spots'),
  ('S_NYERI_SENDI',   'Nyeri sendi dan otot',                 'Joint and muscle pain'),
  ('S_LUKA',          'Luka terbuka / lecet',                  'Open wound / abrasion'),
  ('S_INFEKSI_KULIT', 'Tanda infeksi kulit (merah, bengkak, bernanah)', 'Skin infection signs'),
  ('S_KULIT_PANAS',   'Kulit merah, panas, kering',           'Red, hot, dry skin'),
  ('S_PENDARAHAN',    'Perdarahan / bintik merah (petekie)',   'Bleeding / petechiae'),

  -- Pernapasan
  ('S_BATUK_KRONIK',  'Batuk kronik (> 2 minggu)',             'Chronic cough (>2 weeks)'),
  ('S_BATUK_DARAH',   'Batuk darah (hemoptisis)',              'Hemoptysis'),
  ('S_KERINGAT_MALAM','Keringat malam',                        'Night sweats'),

  -- Hati
  ('S_IKTERUS',       'Ikterus / kulit dan mata kuning',       'Jaundice'),
  ('S_URINE_GELAP',   'Urine berwarna gelap (teh)',            'Dark urine (tea-colored)'),

  -- Riwayat eksposur
  ('S_GIGITAN_HEWAN', 'Riwayat gigitan / cakaran hewan',       'History of animal bite/scratch'),
  ('S_RIWAYAT_ENDEMIK','Riwayat perjalanan ke daerah endemik malaria','History of travel to malaria-endemic area'),

  -- Heat
  ('S_PUSING',        'Pusing / kepala ringan',                'Dizziness / lightheadedness'),
  ('S_SUHU_TINGGI_NONINFEKSI', 'Suhu tubuh tinggi tanpa fokus infeksi jelas', 'Elevated body temperature without clear infection focus')
ON CONFLICT (kode) DO NOTHING;

-- ── 2. Master Penyakit (expert_diseases) ─────────────────────────────────────
INSERT INTO expert_diseases (nama, deskripsi, rekomendasi_default) VALUES
  (
    'Bali Belly (Diare Wisatawan)',
    'Gastroenteritis akut akibat kontaminasi makanan/minuman. Penyakit paling umum pada wisatawan ke Bali.',
    '{"Rendah":"Minum cairan elektrolit, hindari makanan berisiko. Istirahat.","Sedang":"Pertimbangkan oralit. Konsultasi dokter jika tidak membaik dalam 48 jam.","Tinggi":"Konsultasi dokter segera. Pertimbangkan antibiotik oral jika ada tanda demam.","Darurat":"Segera ke IGD, risiko dehidrasi berat."}'
  ),
  (
    'Demam Berdarah Dengue (DBD)',
    'Infeksi virus dengue yang ditularkan nyamuk Aedes aegypti. Endemik di Bali.',
    '{"Rendah":"Pantau gejala, minum banyak cairan, kompres jika demam.","Sedang":"Konsultasi dokter dalam 24 jam, pemeriksaan trombosit.","Tinggi":"Konsultasi dokter segera atau ke UGD. Risiko kebocoran plasma.","Darurat":"Segera ke IGD. Risiko syok dengue yang mengancam jiwa."}'
  ),
  (
    'Rabies',
    'Infeksi virus rabies dari gigitan/cakaran hewan (anjing, kera, kelelawar). Bali merupakan daerah endemik rabies.',
    '{"Rendah":"Cuci luka segera dengan sabun dan air mengalir 15 menit. Segera ke faskes untuk vaksinasi post-exposure (PEP) — JANGAN tunggu gejala.","Sedang":"Segera ke IGD untuk PEP (vaksin + imunoglobulin). Waktu kritis.","Tinggi":"Darurat medis. Segera IGD.","Darurat":"Darurat medis mengancam jiwa. Hubungi 119 segera."}'
  ),
  (
    'Hepatitis A',
    'Infeksi virus hepatitis A melalui makanan/minuman terkontaminasi. Dapat dicegah dengan vaksinasi.',
    '{"Rendah":"Istirahat, hindari alkohol, diet ringan. Pantau selama 1–2 minggu.","Sedang":"Konsultasi dokter, cek fungsi hati (SGOT/SGPT).","Tinggi":"Rujuk ke dokter spesialis penyakit dalam.","Darurat":"Segera ke IGD, risiko gagal hati akut."}'
  ),
  (
    'Demam Tifoid (Tifus)',
    'Infeksi bakteri Salmonella typhi melalui makanan/minuman terkontaminasi.',
    '{"Rendah":"Istirahat, minum cairan cukup, diet lunak.","Sedang":"Konsultasi dokter untuk antibiotik. Pantau ketat.","Tinggi":"Konsultasi dokter segera, risiko komplikasi perforasi usus.","Darurat":"Segera ke IGD. Risiko perforasi atau sepsis."}'
  ),
  (
    'Japanese Encephalitis',
    'Infeksi virus JE yang ditularkan nyamuk di daerah pertanian/pedesaan. Dapat dicegah dengan vaksinasi.',
    '{"Rendah":"Konsultasi dokter segera.","Sedang":"Rujuk ke rumah sakit untuk pemeriksaan lebih lanjut.","Tinggi":"Segera ke IGD, risiko kerusakan otak permanen.","Darurat":"Darurat medis. Hubungi 119 segera."}'
  ),
  (
    'Heat-Related Illness (Penyakit Terkait Panas)',
    'Kondisi akibat paparan panas berlebih: heat exhaustion hingga heat stroke. Umum di Bali saat musim panas.',
    '{"Rendah":"Pindah ke tempat sejuk, minum air putih, kompres dingin.","Sedang":"Istirahat di tempat ber-AC, minum elektrolit, pantau ketat.","Tinggi":"Konsultasi dokter segera atau ke UGD.","Darurat":"Heat stroke — darurat medis. Hubungi 119 segera."}'
  ),
  (
    'Infeksi Luka / Gangguan Kulit',
    'Luka terbuka atau infeksi kulit akibat trauma, serangga, atau lingkungan lembap di daerah tropis.',
    '{"Rendah":"Bersihkan luka, tutup dengan perban bersih. Pantau tanda infeksi.","Sedang":"Konsultasi dokter, pertimbangkan antibiotik topikal/oral.","Tinggi":"Konsultasi dokter segera, risiko selulitis atau MRSA.","Darurat":"Segera ke IGD, risiko sepsis."}'
  ),
  (
    'Tuberkulosis (TBC) — Screening Post-Travel',
    'Infeksi bakteri Mycobacterium tuberculosis. Skrining diindikasikan untuk wisatawan yang tinggal lama di daerah endemik.',
    '{"Rendah":"Konsultasi dokter untuk tes Mantoux atau IGRA.","Sedang":"Rujuk ke poli paru atau spesialis penyakit dalam.","Tinggi":"Pemeriksaan BTA sputum dan foto toraks segera.","Darurat":"Segera ke fasilitas kesehatan untuk isolasi dan pemeriksaan."}'
  ),
  (
    'Malaria — Screening Post-Travel',
    'Infeksi parasit Plasmodium yang ditularkan nyamuk Anopheles. Skrining untuk wisatawan dari daerah endemik.',
    '{"Rendah":"Konsultasi dokter untuk pemeriksaan darah tepi / RDT malaria.","Sedang":"Pemeriksaan darah tepi SEGERA, terutama jika demam > 38°C.","Tinggi":"Segera ke faskes untuk diagnostik dan terapi antimalarial.","Darurat":"Segera ke IGD. Risiko malaria berat / cerebral malaria."}'
  )
ON CONFLICT (nama) DO NOTHING;

-- ── 3. Expert Rules ───────────────────────────────────────────────────────────
-- Gunakan DO block karena expert_rules.created_by wajib merujuk users.id
-- Buat system user sementara; ID-nya dipakai sebagai created_by pada rules seed.
-- ⚠️  MB/MD di bawah adalah PLACEHOLDER — wajib verifikasi narasumber perawat.
DO $$
DECLARE
    sys_uid   BIGINT;
    d_bb      BIGINT; -- Bali Belly
    d_dbd     BIGINT; -- DBD
    d_rabies  BIGINT;
    d_hepa    BIGINT; -- Hepatitis A
    d_tifus   BIGINT;
    d_je      BIGINT; -- Japanese Encephalitis
    d_heat    BIGINT;
    d_luka    BIGINT; -- Infeksi Luka
    d_tbc     BIGINT;
    d_malaria BIGINT;
    -- symptom IDs
    s_diare          BIGINT; s_mual           BIGINT; s_muntah       BIGINT;
    s_kram           BIGINT; s_nyeri_ruq      BIGINT; s_sembelit     BIGINT;
    s_demam_ringan   BIGINT; s_demam_tinggi   BIGINT; s_demam_periodik BIGINT;
    s_menggigil      BIGINT; s_berkeringat    BIGINT; s_kelemahan    BIGINT;
    s_bb_turun       BIGINT; s_nyeri_kepala   BIGINT; s_nyeri_mata   BIGINT;
    s_kaku_leher     BIGINT; s_kejang         BIGINT; s_bingung      BIGINT;
    s_hidrofobia     BIGINT; s_kesemutan_gigitan BIGINT;
    s_ruam           BIGINT; s_ruam_merah     BIGINT; s_nyeri_sendi  BIGINT;
    s_luka           BIGINT; s_infeksi_kulit  BIGINT; s_kulit_panas  BIGINT;
    s_pendarahan     BIGINT; s_batuk_kronik   BIGINT; s_batuk_darah  BIGINT;
    s_keringat_malam BIGINT; s_ikterus        BIGINT; s_urine_gelap  BIGINT;
    s_gigitan_hewan  BIGINT; s_riwayat_endemik BIGINT;
    s_pusing         BIGINT; s_suhu_tinggi    BIGINT;
BEGIN
    -- Buat system user untuk seeding; gunakan ON CONFLICT agar idempoten
    INSERT INTO users (email, provider)
    VALUES ('system.seed@balihealth.internal', 'email')
    ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
    RETURNING id INTO sys_uid;

    -- Load disease IDs
    SELECT id INTO d_bb     FROM expert_diseases WHERE nama = 'Bali Belly (Diare Wisatawan)';
    SELECT id INTO d_dbd    FROM expert_diseases WHERE nama = 'Demam Berdarah Dengue (DBD)';
    SELECT id INTO d_rabies FROM expert_diseases WHERE nama = 'Rabies';
    SELECT id INTO d_hepa   FROM expert_diseases WHERE nama = 'Hepatitis A';
    SELECT id INTO d_tifus  FROM expert_diseases WHERE nama = 'Demam Tifoid (Tifus)';
    SELECT id INTO d_je     FROM expert_diseases WHERE nama = 'Japanese Encephalitis';
    SELECT id INTO d_heat   FROM expert_diseases WHERE nama = 'Heat-Related Illness (Penyakit Terkait Panas)';
    SELECT id INTO d_luka   FROM expert_diseases WHERE nama = 'Infeksi Luka / Gangguan Kulit';
    SELECT id INTO d_tbc    FROM expert_diseases WHERE nama = 'Tuberkulosis (TBC) — Screening Post-Travel';
    SELECT id INTO d_malaria FROM expert_diseases WHERE nama = 'Malaria — Screening Post-Travel';

    -- Load symptom IDs
    SELECT symptom_id INTO s_diare            FROM expert_symptoms WHERE kode = 'S_DIARE';
    SELECT symptom_id INTO s_mual             FROM expert_symptoms WHERE kode = 'S_MUAL';
    SELECT symptom_id INTO s_muntah           FROM expert_symptoms WHERE kode = 'S_MUNTAH';
    SELECT symptom_id INTO s_kram             FROM expert_symptoms WHERE kode = 'S_KRAM_PERUT';
    SELECT symptom_id INTO s_nyeri_ruq        FROM expert_symptoms WHERE kode = 'S_NYERI_RUQ';
    SELECT symptom_id INTO s_sembelit         FROM expert_symptoms WHERE kode = 'S_SEMBELIT';
    SELECT symptom_id INTO s_demam_ringan     FROM expert_symptoms WHERE kode = 'S_DEMAM_RINGAN';
    SELECT symptom_id INTO s_demam_tinggi     FROM expert_symptoms WHERE kode = 'S_DEMAM_TINGGI';
    SELECT symptom_id INTO s_demam_periodik   FROM expert_symptoms WHERE kode = 'S_DEMAM_PERIODIK';
    SELECT symptom_id INTO s_menggigil        FROM expert_symptoms WHERE kode = 'S_MENGGIGIL';
    SELECT symptom_id INTO s_berkeringat      FROM expert_symptoms WHERE kode = 'S_BERKERINGAT';
    SELECT symptom_id INTO s_kelemahan        FROM expert_symptoms WHERE kode = 'S_KELEMAHAN';
    SELECT symptom_id INTO s_bb_turun         FROM expert_symptoms WHERE kode = 'S_BB_TURUN';
    SELECT symptom_id INTO s_nyeri_kepala     FROM expert_symptoms WHERE kode = 'S_NYERI_KEPALA';
    SELECT symptom_id INTO s_nyeri_mata       FROM expert_symptoms WHERE kode = 'S_NYERI_MATA';
    SELECT symptom_id INTO s_kaku_leher       FROM expert_symptoms WHERE kode = 'S_KAKU_LEHER';
    SELECT symptom_id INTO s_kejang           FROM expert_symptoms WHERE kode = 'S_KEJANG';
    SELECT symptom_id INTO s_bingung          FROM expert_symptoms WHERE kode = 'S_BINGUNG';
    SELECT symptom_id INTO s_hidrofobia       FROM expert_symptoms WHERE kode = 'S_HIDROFOBIA';
    SELECT symptom_id INTO s_kesemutan_gigitan FROM expert_symptoms WHERE kode = 'S_KESEMUTAN_GIGITAN';
    SELECT symptom_id INTO s_ruam             FROM expert_symptoms WHERE kode = 'S_RUAM';
    SELECT symptom_id INTO s_ruam_merah       FROM expert_symptoms WHERE kode = 'S_RUAM_MERAH';
    SELECT symptom_id INTO s_nyeri_sendi      FROM expert_symptoms WHERE kode = 'S_NYERI_SENDI';
    SELECT symptom_id INTO s_luka             FROM expert_symptoms WHERE kode = 'S_LUKA';
    SELECT symptom_id INTO s_infeksi_kulit    FROM expert_symptoms WHERE kode = 'S_INFEKSI_KULIT';
    SELECT symptom_id INTO s_kulit_panas      FROM expert_symptoms WHERE kode = 'S_KULIT_PANAS';
    SELECT symptom_id INTO s_pendarahan       FROM expert_symptoms WHERE kode = 'S_PENDARAHAN';
    SELECT symptom_id INTO s_batuk_kronik     FROM expert_symptoms WHERE kode = 'S_BATUK_KRONIK';
    SELECT symptom_id INTO s_batuk_darah      FROM expert_symptoms WHERE kode = 'S_BATUK_DARAH';
    SELECT symptom_id INTO s_keringat_malam   FROM expert_symptoms WHERE kode = 'S_KERINGAT_MALAM';
    SELECT symptom_id INTO s_ikterus          FROM expert_symptoms WHERE kode = 'S_IKTERUS';
    SELECT symptom_id INTO s_urine_gelap      FROM expert_symptoms WHERE kode = 'S_URINE_GELAP';
    SELECT symptom_id INTO s_gigitan_hewan    FROM expert_symptoms WHERE kode = 'S_GIGITAN_HEWAN';
    SELECT symptom_id INTO s_riwayat_endemik  FROM expert_symptoms WHERE kode = 'S_RIWAYAT_ENDEMIK';
    SELECT symptom_id INTO s_pusing           FROM expert_symptoms WHERE kode = 'S_PUSING';
    SELECT symptom_id INTO s_suhu_tinggi      FROM expert_symptoms WHERE kode = 'S_SUHU_TINGGI_NONINFEKSI';

    -- ── PRE-TRAVEL RULES ─────────────────────────────────────────────────────
    -- Bali Belly: diare + kram (gejala utama)
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Bali Belly — diare + kram perut',
        jsonb_build_array(s_diare, s_kram),
        d_bb, 0.700, 0.800, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Bali Belly: diare + mual + muntah
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Bali Belly — diare + mual + muntah',
        jsonb_build_array(s_diare, s_mual, s_muntah),
        d_bb, 0.750, 0.850, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Bali Belly: diare + demam ringan + kram
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Bali Belly — diare + demam ringan + kram',
        jsonb_build_array(s_diare, s_demam_ringan, s_kram),
        d_bb, 0.700, 0.800, 0.100, 'pre_travel', 'published', sys_uid
    );

    -- DBD: demam tinggi + nyeri sendi + nyeri kepala
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'DBD — demam tinggi + nyeri sendi + nyeri kepala',
        jsonb_build_array(s_demam_tinggi, s_nyeri_sendi, s_nyeri_kepala),
        d_dbd, 0.700, 0.800, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- DBD: demam tinggi + ruam + nyeri mata
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'DBD — demam tinggi + ruam + nyeri mata',
        jsonb_build_array(s_demam_tinggi, s_ruam, s_nyeri_mata),
        d_dbd, 0.750, 0.850, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- DBD: demam tinggi + perdarahan (tanda berat)
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'DBD berat — demam tinggi + perdarahan',
        jsonb_build_array(s_demam_tinggi, s_pendarahan),
        d_dbd, 0.800, 0.900, 0.100, 'pre_travel', 'published', sys_uid
    );

    -- Rabies: riwayat gigitan + kesemutan di area gigitan
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Rabies — gigitan hewan + kesemutan area gigitan',
        jsonb_build_array(s_gigitan_hewan, s_kesemutan_gigitan),
        d_rabies, 0.800, 0.900, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Rabies: gigitan hewan + hidrofobia (stadium lanjut — DARURAT)
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Rabies stadium lanjut — hidrofobia + kebingungan',
        jsonb_build_array(s_gigitan_hewan, s_hidrofobia, s_bingung),
        d_rabies, 0.900, 0.950, 0.050, 'pre_travel', 'published', sys_uid
    );
    -- Rabies: gigitan hewan (eksposur saja — perlu PEP segera)
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Eksposur rabies — riwayat gigitan hewan',
        jsonb_build_array(s_gigitan_hewan),
        d_rabies, 0.600, 0.700, 0.100, 'pre_travel', 'published', sys_uid
    );

    -- Hepatitis A: ikterus + mual + kelemahan
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Hepatitis A — ikterus + mual + kelemahan',
        jsonb_build_array(s_ikterus, s_mual, s_kelemahan),
        d_hepa, 0.750, 0.850, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Hepatitis A: ikterus + urine gelap + nyeri RUQ
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Hepatitis A — ikterus + urine gelap + nyeri RUQ',
        jsonb_build_array(s_ikterus, s_urine_gelap, s_nyeri_ruq),
        d_hepa, 0.800, 0.900, 0.100, 'pre_travel', 'published', sys_uid
    );

    -- Tifus: demam tinggi + kelemahan + nyeri kepala
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Tifus — demam tinggi persisten + kelemahan + nyeri kepala',
        jsonb_build_array(s_demam_tinggi, s_kelemahan, s_nyeri_kepala),
        d_tifus, 0.650, 0.750, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Tifus: demam + sembelit/diare + ruam merah
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Tifus — demam + rose spots + sembelit',
        jsonb_build_array(s_demam_tinggi, s_ruam_merah, s_sembelit),
        d_tifus, 0.750, 0.850, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Tifus: demam + diare + nyeri perut
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Tifus — demam + diare + nyeri perut',
        jsonb_build_array(s_demam_tinggi, s_diare, s_kram),
        d_tifus, 0.600, 0.700, 0.100, 'pre_travel', 'published', sys_uid
    );

    -- Japanese Encephalitis: demam tinggi + kaku leher + nyeri kepala parah
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Japanese Encephalitis — demam tinggi + kaku leher + nyeri kepala',
        jsonb_build_array(s_demam_tinggi, s_kaku_leher, s_nyeri_kepala),
        d_je, 0.800, 0.900, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- JE: demam + kejang + kebingungan (stadium berat)
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Japanese Encephalitis berat — demam + kejang + kebingungan',
        jsonb_build_array(s_demam_tinggi, s_kejang, s_bingung),
        d_je, 0.900, 0.950, 0.050, 'pre_travel', 'published', sys_uid
    );

    -- Heat exhaustion: pusing + kelemahan + berkeringat + kulit panas
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Heat Exhaustion — pusing + kelemahan + berkeringat',
        jsonb_build_array(s_pusing, s_kelemahan, s_berkeringat),
        d_heat, 0.650, 0.750, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Heat stroke: kulit panas kering + kebingungan + suhu tinggi (DARURAT)
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Heat Stroke — kulit panas kering + kebingungan + suhu tinggi',
        jsonb_build_array(s_kulit_panas, s_bingung, s_suhu_tinggi),
        d_heat, 0.900, 0.950, 0.050, 'pre_travel', 'published', sys_uid
    );

    -- Infeksi luka: luka + tanda infeksi kulit
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Infeksi Luka — luka terbuka + tanda infeksi',
        jsonb_build_array(s_luka, s_infeksi_kulit),
        d_luka, 0.750, 0.850, 0.100, 'pre_travel', 'published', sys_uid
    );
    -- Infeksi luka dengan demam (selulitis / sepsis dini)
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Infeksi Luka Berat — luka + infeksi + demam',
        jsonb_build_array(s_luka, s_infeksi_kulit, s_demam_tinggi),
        d_luka, 0.800, 0.900, 0.100, 'pre_travel', 'published', sys_uid
    );

    -- ── POST-TRAVEL RULES ────────────────────────────────────────────────────
    -- TBC: batuk kronik + keringat malam + bb turun
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Skrining TBC — batuk kronik + keringat malam + penurunan BB',
        jsonb_build_array(s_batuk_kronik, s_keringat_malam, s_bb_turun),
        d_tbc, 0.800, 0.900, 0.100, 'post_travel', 'published', sys_uid
    );
    -- TBC: batuk darah + batuk kronik
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Skrining TBC — batuk kronik + batuk darah',
        jsonb_build_array(s_batuk_kronik, s_batuk_darah),
        d_tbc, 0.850, 0.900, 0.050, 'post_travel', 'published', sys_uid
    );
    -- TBC: batuk kronik + demam ringan
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Skrining TBC — batuk kronik + demam ringan sore hari',
        jsonb_build_array(s_batuk_kronik, s_demam_ringan, s_kelemahan),
        d_tbc, 0.700, 0.800, 0.100, 'post_travel', 'published', sys_uid
    );

    -- Malaria: demam periodik + menggigil + berkeringat + riwayat endemik
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Skrining Malaria — demam periodik + menggigil + riwayat endemik',
        jsonb_build_array(s_demam_periodik, s_menggigil, s_riwayat_endemik),
        d_malaria, 0.800, 0.900, 0.100, 'post_travel', 'published', sys_uid
    );
    -- Malaria: demam periodik + berkeringat + nyeri kepala + riwayat endemik
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Skrining Malaria — demam periodik + berkeringat + nyeri kepala',
        jsonb_build_array(s_demam_periodik, s_berkeringat, s_nyeri_kepala, s_riwayat_endemik),
        d_malaria, 0.750, 0.850, 0.100, 'post_travel', 'published', sys_uid
    );
    -- Malaria berat: demam + menggigil + kebingungan
    INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
    VALUES (
        'Malaria Berat — demam + menggigil + kebingungan',
        jsonb_build_array(s_demam_periodik, s_menggigil, s_bingung),
        d_malaria, 0.900, 0.950, 0.050, 'post_travel', 'published', sys_uid
    );

END $$;
