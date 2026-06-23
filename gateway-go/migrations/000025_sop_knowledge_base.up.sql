-- Migration 000025: SOP Knowledge Base dari Eka Swedarma / Udayana
-- Menambahkan gejala, penyakit, rule, dan alur panduan darurat berdasarkan SOP

-- ─────────────────────────────────────────────────────────────────────────────
-- SECTION 1: New Symptoms
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO expert_symptoms (kode, label_id, label_en) VALUES
  ('S_TINJA_BERDARAH',       'BAB berdarah / berlendir',                        'Bloody / mucoid stool'),
  ('S_DEHIDRASI',            'Tanda dehidrasi (mulut kering, kulit tidak elastis, urine sedikit)', 'Signs of dehydration'),
  ('S_RUAM_URTIKARIA',       'Ruam biduran / urtikaria (bentol merah gatal)',    'Urticaria / hives'),
  ('S_MATA_MERAH_GATAL',     'Mata merah dan gatal (konjungtivitis alergi)',     'Red itchy eyes'),
  ('S_BENGKAK_BIBIR_WAJAH',  'Bengkak pada bibir, wajah, atau lidah',           'Swelling of lips, face or tongue'),
  ('S_SESAK_NAPAS',          'Sesak napas / napas berbunyi (wheezing)',          'Shortness of breath / wheezing'),
  ('S_KULIT_MERAH_BAKAR',    'Kulit merah akibat paparan matahari (sunburn)',    'Sunburned skin (red)'),
  ('S_LEPUHAN_KULIT',        'Lepuhan pada kulit akibat bakar matahari',         'Skin blisters from sunburn'),
  ('S_PAPARAN_MATAHARI_LAMA','Riwayat paparan matahari berkepanjangan',          'History of prolonged sun exposure'),
  ('S_MUAL_PERJALANAN',      'Mual saat bepergian (mabuk darat/laut/udara)',    'Nausea during travel'),
  ('S_KERINGAT_DINGIN',      'Keringat dingin',                                 'Cold sweating'),
  ('S_PUSING_PERJALANAN',    'Pusing / vertigo saat berkendara atau berlayar',   'Dizziness during travel'),
  ('S_GANGGUAN_PENGLIHATAN', 'Gangguan penglihatan (kabur / kembar)',            'Visual disturbance (blurred/double)'),
  ('S_SEMPOYONGAN',          'Sempoyongan / kehilangan keseimbangan',            'Unsteady gait / loss of balance'),
  ('S_KONSUMSI_ALKOHOL_LOKAL','Riwayat konsumsi alkohol lokal (arak/tuak) di Bali', 'History of consuming local alcohol in Bali'),
  ('S_GIGITAN_ULAR',         'Riwayat gigitan ular',                            'History of snake bite'),
  ('S_SENGATAN_UBUR',        'Sengatan ubur-ubur atau makhluk laut beracun',    'Jellyfish or marine creature sting'),
  ('S_TUSUKAN_DURI_LAUT',    'Tertusuk duri bulu babi atau benda tajam di laut','Sea urchin spine or sharp marine object puncture'),
  ('S_PINGSAN_SEMENTARA',    'Pingsan sementara / kehilangan kesadaran singkat', 'Brief loss of consciousness / syncope'),
  ('S_PUCAT',                'Wajah / kulit pucat',                             'Pale face / skin'),
  ('S_NADI_LEMAH',           'Denyut nadi lemah atau tidak teratur',            'Weak or irregular pulse'),
  ('S_PERDARAHAN_AKTIF',     'Perdarahan aktif yang sulit berhenti',            'Active bleeding that is hard to stop'),
  ('S_LUKA_KECIL',           'Luka kecil / lecet permukaan',                   'Minor cut or surface abrasion'),
  ('S_BENGKAK_SENDI_TRAUMA', 'Bengkak pada sendi akibat benturan / jatuh',     'Joint swelling due to trauma/fall'),
  ('S_TIDAK_BISA_GERAK',     'Anggota badan tidak bisa digerakkan normal',      'Limb cannot move normally'),
  ('S_SENGATAN_LEBAH',       'Riwayat sengatan lebah atau tawon',               'History of bee or wasp sting'),
  ('S_BENGKAK_LOKAL_SENGATAN','Bengkak dan nyeri lokal di area sengatan',       'Local swelling and pain at sting site'),
  ('S_LEMAS_DI_PANAS',       'Lemas saat berada di lingkungan panas terik',     'Weakness in hot environment'),
  ('S_TIDAK_BERKERINGAT',    'Tidak berkeringat meski suhu tubuh tinggi (anhidrosis)', 'Absence of sweating despite high body temperature'),
  ('S_NYERI_SENDI_TRAUMA',   'Nyeri sendi akibat benturan atau cedera',         'Joint pain from trauma or injury'),
  ('S_DETAK_CEPAT',          'Detak jantung cepat (takikardia)',                'Rapid heart rate (tachycardia)'),
  ('S_BENGKAK_TENGGOROKAN',  'Rasa sesak / bengkak di tenggorokan',             'Throat tightness / swelling'),
  ('S_NYERI_DI_AREA_SENGATAN','Nyeri menjalar dari area sengatan',              'Radiating pain from sting area')
ON CONFLICT (kode) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- SECTION 2: New Diseases (IDs will be 11–24 sequentially)
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO expert_diseases (nama, deskripsi, rekomendasi_default) VALUES
  (
    'Reaksi Alergi Ringan',
    'Respons imun terhadap alergen berupa urtikaria, gatal, atau mata merah tanpa gejala anafilaksis. Umum pada wisatawan yang terpapar makanan, sengatan serangga, atau lingkungan baru.',
    '{"Rendah":"Konsumsi antihistamin oral (cetirizine/loratadine). Hindari pemicu alergen. Kompres dingin pada area gatal.","Sedang":"Antihistamin oral + krim kortikosteroid topikal. Monitor gejala 24 jam. Cari apotek terdekat.","Tinggi":"Kunjungi klinik/dokter segera. Mungkin perlu antihistamin injeksi atau kortikosteroid oral.","Darurat":"Jika muncul sesak napas atau bengkak wajah/lidah, segera ke IGD — risiko anafilaksis."}'
  ),
  (
    'Anafilaksis',
    'Reaksi alergi berat dan mengancam jiwa ditandai bengkak wajah/tenggorokan, sesak napas, dan penurunan tekanan darah. Membutuhkan epinefrin segera dan bantuan medis darurat.',
    '{"Rendah":"N/A — anafilaksis selalu darurat.","Sedang":"Hubungi ambulans segera. Jika tersedia, injeksi epinefrin auto-injector (EpiPen).","Tinggi":"Posisi berbaring kaki ditinggikan. Bebaskan jalan napas. Siapkan CPR jika tidak responsif.","Darurat":"DARURAT: Hubungi 119/ambulans SEKARANG. Injeksi epinefrin. Jaga jalan napas. Awasi hingga tim medis tiba."}'
  ),
  (
    'Sunburn (Luka Bakar Matahari)',
    'Kerusakan kulit akibat paparan sinar UV berlebih, ditandai kemerahan, nyeri, dan lepuhan. Di Bali, risiko tinggi bagi wisatawan yang tidak menggunakan tabir surya.',
    '{"Rendah":"Pindah ke tempat teduh. Kompres dingin dan oleskan aloe vera / pelembap. Minum banyak air.","Sedang":"Kompres dingin 15–20 menit. Ibuprofen/paracetamol untuk nyeri. Hindari paparan matahari ulang.","Tinggi":"Kunjungi klinik jika ada lepuhan luas atau demam. Jangan pecahkan lepuhan.","Darurat":"Jika demam tinggi, menggigil, atau dehidrasi berat, segera ke IGD."}'
  ),
  (
    'Mabuk Perjalanan (Motion Sickness)',
    'Gangguan keseimbangan sensorik akibat gerakan kendaraan, perahu, atau pesawat. Umum pada wisatawan saat perjalanan darat atau laut di Bali.',
    '{"Rendah":"Duduk di tempat dengan gerakan minimal (tengah kapal/bus). Pandang cakrawala. Hindari membaca.","Sedang":"Konsumsi obat antimabuk (dimenhydrinate/meclizine) 30 menit sebelum perjalanan.","Tinggi":"Hentikan perjalanan jika memungkinkan. Istirahat di udara segar. Minum air sedikit-sedikit.","Darurat":"Jika muntah terus-menerus dan tidak bisa minum cairan, cari bantuan medis untuk rehidrasi IV."}'
  ),
  (
    'Keracunan Metanol (Arak Oplosan)',
    'Keracunan akibat konsumsi minuman beralkohol yang mengandung metanol, sering terjadi pada wisatawan yang mengonsumsi arak oplosan di Bali. Dapat menyebabkan kebutaan dan kematian.',
    '{"Rendah":"Hentikan konsumsi alkohol. Minum air banyak. Monitor gejala visual.","Sedang":"Segera cari bantuan medis. Informasikan jenis alkohol yang dikonsumsi.","Tinggi":"Ke IGD segera. Penanganan memerlukan etanol sebagai antidot atau fomepizole.","Darurat":"DARURAT MEDIS: Gangguan penglihatan akibat metanol bisa permanen. Hubungi 119 segera."}'
  ),
  (
    'Syok (Shock)',
    'Kondisi kegagalan sirkulasi yang menyebabkan organ vital kekurangan oksigen. Dapat disebabkan perdarahan masif, dehidrasi berat, atau reaksi alergi berat (anafilaksis).',
    '{"Rendah":"N/A — syok selalu darurat.","Sedang":"Baringkan korban, tinggikan kaki 30 cm. Jaga kehangatan. Hubungi bantuan medis.","Tinggi":"Hentikan perdarahan jika ada. Jangan beri makan/minum. Monitor kesadaran dan napas.","Darurat":"DARURAT: Hubungi 119 segera. Mulai CPR jika tidak ada denyut nadi dan napas."}'
  ),
  (
    'Sinkop (Pingsan)',
    'Kehilangan kesadaran sementara akibat berkurangnya aliran darah ke otak, sering disebabkan panas terik, dehidrasi, atau berdiri terlalu lama.',
    '{"Rendah":"Baringkan di tempat sejuk, tinggikan kaki. Longgarkan pakaian ketat. Beri minum saat sadar.","Sedang":"Monitor hingga benar-benar pulih. Hindari berdiri tiba-tiba. Minum cairan elektrolit.","Tinggi":"Kunjungi dokter jika sering pingsan atau tidak pulih dalam 5 menit.","Darurat":"Jika tidak sadar > 5 menit atau tidak ada napas, hubungi 119 dan mulai DRSABC."}'
  ),
  (
    'Gigitan Ular',
    'Cedera akibat gigitan ular, berpotensi berbahaya jika ular berbisa. Di Bali terdapat beberapa spesies ular berbisa. Membutuhkan penanganan segera di fasilitas kesehatan.',
    '{"Rendah":"Tenangkan korban, imobilisasi area gigitan. Catat ciri-ciri ular jika aman.","Sedang":"Pasang perban tekan (bukan tourniquet ketat). Jangan isap luka. Segera ke rumah sakit.","Tinggi":"Ke IGD dengan serum antivenom tersedia. Bawa foto/ciri ular jika ada.","Darurat":"DARURAT: Jangan tunda — antivenom harus diberikan sesegera mungkin. Hubungi 119."}'
  ),
  (
    'Sengatan Hewan Laut',
    'Cedera akibat sengatan ubur-ubur, bulu babi, atau hewan laut beracun lain. Umum pada wisatawan yang berenang atau snorkeling di perairan Bali.',
    '{"Rendah":"Sengatan ubur-ubur: siram cuka, rendam air panas 40–45°C selama 20 menit. Jangan digosok.","Sedang":"Bulu babi: rendam air panas 30 menit, duri kecil akan larut. Duri besar: ke klinik.","Tinggi":"Jika muncul reaksi alergi (bengkak, sesak), segera ke IGD.","Darurat":"Reaksi anafilaksis akibat sengatan laut: hubungi 119 segera."}'
  ),
  (
    'Reaksi Gigitan Serangga',
    'Reaksi lokal atau sistemik akibat gigitan atau sengatan serangga (nyamuk, lebah, tawon, semut api). Di Bali umum ditemui wisatawan di alam terbuka.',
    '{"Rendah":"Bersihkan area gigitan. Kompres dingin. Oleskan krim antihistamin atau hidrokortison topikal.","Sedang":"Antihistamin oral untuk gatal sistemik. Monitor selama 30 menit untuk reaksi alergi.","Tinggi":"Jika bengkak menyebar atau sesak napas muncul, ke klinik segera.","Darurat":"Sengatan lebah/tawon multipel atau reaksi anafilaksis: hubungi 119 segera."}'
  ),
  (
    'Perdarahan Aktif',
    'Perdarahan dari luka yang sulit dihentikan. Dapat berasal dari luka kecil hingga trauma berat. Memerlukan penanganan pertolongan pertama segera untuk mencegah kehilangan darah masif.',
    '{"Rendah":"Tekan luka bersih dengan kain bersih 10 menit. Tinggikan bagian yang berdarah.","Sedang":"Tekan kuat dan konsisten. Jangan angkat kain penekan. Tambah kain jika perlu.","Tinggi":"Luka dalam atau tidak berhenti dalam 15 menit: ke klinik/IGD segera.","Darurat":"Perdarahan masif atau syok: hubungi 119 segera. Tekan luka sepanjang perjalanan."}'
  ),
  (
    'Keseleo / Cedera Ligamen',
    'Peregangan atau robekan ligamen akibat gerakan tiba-tiba atau jatuh. Umum pada wisatawan yang beraktivitas fisik di Bali (hiking, surfing).',
    '{"Rendah":"Terapkan RICE: Rest (istirahat), Ice (kompres es 20 menit), Compression (perban elastis), Elevation (tinggikan).","Sedang":"Hindari aktivitas berat 48 jam. Ibuprofen untuk nyeri dan bengkak.","Tinggi":"Jika tidak bisa berjalan atau bengkak sangat besar, ke klinik untuk X-ray.","Darurat":"Jika tulang tampak tidak normal atau sangat nyeri saat ditekan, mungkin patah — ke IGD."}'
  ),
  (
    'Patah Tulang / Fraktur',
    'Retakan atau patahan tulang akibat trauma. Memerlukan imobilisasi segera dan penanganan medis. Wisatawan berisiko saat kecelakaan lalu lintas atau aktivitas ekstrem.',
    '{"Rendah":"Imobilisasi area yang dicurigai patah dengan papan/bidai darurat. Jangan paksa luruskan.","Sedang":"Gunakan kain sebagai sling untuk lengan. Dukung area fraktur. Ke klinik/IGD segera.","Tinggi":"Patah tulang terbuka (tulang tembus kulit): tutupi dengan kain bersih, JANGAN tekan tulang. IGD segera.","Darurat":"Trauma kepala/tulang belakang: JANGAN gerakkan korban. Hubungi 119 segera."}'
  ),
  (
    'Sengatan Lebah / Tawon',
    'Sengatan lebah atau tawon menghasilkan racun yang menyebabkan nyeri lokal, bengkak, dan berpotensi reaksi alergi sistemik atau anafilaksis pada individu sensitif.',
    '{"Rendah":"Keluarkan sengat dengan cara dikerik (jangan dipencet). Kompres dingin. Antihistamin oral.","Sedang":"Pantau 30 menit untuk reaksi alergi. Ibuprofen untuk nyeri.","Tinggi":"Jika bengkak menyebar dari area lokal, ke klinik segera.","Darurat":"Sengatan multipel, sesak napas, atau anafilaksis: hubungi 119 SEGERA."}'
  )
ON CONFLICT (nama) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- SECTION 3: New Expert Rules
-- Menggunakan subquery untuk mendapatkan symptom_id dari kode
-- dan disease_id berdasarkan nama
-- ─────────────────────────────────────────────────────────────────────────────

-- Reaksi Alergi Ringan (disease 11)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Alergi Ringan — urtikaria + mata merah gatal',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_RUAM_URTIKARIA'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_MATA_MERAH_GATAL')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Reaksi Alergi Ringan'),
  0.750, 0.750, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Alergi Ringan — urtikaria + keringat dingin',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_RUAM_URTIKARIA'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KERINGAT_DINGIN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Reaksi Alergi Ringan'),
  0.700, 0.700, 0.150, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Alergi — sengatan serangga + bengkak lokal',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SENGATAN_LEBAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_LOKAL_SENGATAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Reaksi Alergi Ringan'),
  0.700, 0.750, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Anafilaksis (disease 12)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Anafilaksis — bengkak wajah + sesak napas',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_BIBIR_WAJAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SESAK_NAPAS')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Anafilaksis'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Anafilaksis — bengkak tenggorokan + sesak + urtikaria',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_TENGGOROKAN'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SESAK_NAPAS'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_RUAM_URTIKARIA')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Anafilaksis'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Anafilaksis — sesak + nadi lemah + pucat',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SESAK_NAPAS'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NADI_LEMAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUCAT')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Anafilaksis'),
  0.900, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Sunburn (disease 13)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Sunburn ringan — kulit merah + paparan matahari lama',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KULIT_MERAH_BAKAR'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PAPARAN_MATAHARI_LAMA')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sunburn (Luka Bakar Matahari)'),
  0.800, 0.800, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Sunburn berat — kulit merah + lepuhan',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KULIT_MERAH_BAKAR'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_LEPUHAN_KULIT')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sunburn (Luka Bakar Matahari)'),
  0.900, 0.900, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Mabuk Perjalanan (disease 14)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Mabuk perjalanan — mual + pusing saat berkendara',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_MUAL_PERJALANAN'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUSING_PERJALANAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Mabuk Perjalanan (Motion Sickness)'),
  0.800, 0.800, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Mabuk perjalanan — mual + muntah + keringat dingin',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_MUAL_PERJALANAN'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_MUNTAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KERINGAT_DINGIN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Mabuk Perjalanan (Motion Sickness)'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal')
  WHERE EXISTS (SELECT 1 FROM expert_symptoms WHERE kode='S_MUNTAH');

-- Keracunan Metanol (disease 15)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Keracunan metanol — konsumsi alkohol lokal + gangguan penglihatan',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KONSUMSI_ALKOHOL_LOKAL'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_GANGGUAN_PENGLIHATAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Keracunan Metanol (Arak Oplosan)'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Keracunan metanol — konsumsi alkohol lokal + sempoyongan + mual',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KONSUMSI_ALKOHOL_LOKAL'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SEMPOYONGAN'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_MUAL')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Keracunan Metanol (Arak Oplosan)'),
  0.900, 0.900, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal')
  WHERE EXISTS (SELECT 1 FROM expert_symptoms WHERE kode='S_MUAL');

-- Syok (disease 16)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Syok — pucat + nadi lemah + keringat dingin',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUCAT'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NADI_LEMAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KERINGAT_DINGIN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Syok (Shock)'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Syok — perdarahan aktif + pucat + detak cepat',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PERDARAHAN_AKTIF'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUCAT'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_DETAK_CEPAT')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Syok (Shock)'),
  0.900, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Sinkop (disease 17)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Sinkop — pingsan sementara + pucat',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PINGSAN_SEMENTARA'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUCAT')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sinkop (Pingsan)'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Sinkop — pingsan + lemas di panas + keringat dingin',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PINGSAN_SEMENTARA'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_LEMAS_DI_PANAS'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KERINGAT_DINGIN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sinkop (Pingsan)'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Gigitan Ular (disease 18)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Gigitan ular — riwayat gigitan ular + nyeri area sengatan',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_GIGITAN_ULAR'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NYERI_DI_AREA_SENGATAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Gigitan Ular'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Gigitan ular berbisa — gigitan + bengkak lokal + pucat',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_GIGITAN_ULAR'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_LOKAL_SENGATAN'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUCAT')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Gigitan Ular'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Sengatan Hewan Laut (disease 19)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Sengatan hewan laut — sengatan ubur-ubur + nyeri lokal',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SENGATAN_UBUR'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NYERI_DI_AREA_SENGATAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sengatan Hewan Laut'),
  0.900, 0.900, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Tusukan bulu babi — tusukan duri laut + nyeri lokal',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_TUSUKAN_DURI_LAUT'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NYERI_DI_AREA_SENGATAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sengatan Hewan Laut'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Reaksi Gigitan Serangga (disease 20)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Reaksi serangga — sengatan lebah + bengkak lokal',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SENGATAN_LEBAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_LOKAL_SENGATAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Reaksi Gigitan Serangga'),
  0.800, 0.800, 0.150, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Reaksi serangga sistemik — urtikaria + detak cepat setelah gigitan',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_RUAM_URTIKARIA'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_DETAK_CEPAT'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SENGATAN_LEBAH')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Reaksi Gigitan Serangga'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Perdarahan Aktif (disease 21)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Perdarahan aktif — luka + perdarahan tidak berhenti',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_LUKA_KECIL'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PERDARAHAN_AKTIF')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Perdarahan Aktif'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Perdarahan berat — perdarahan aktif + pucat + nadi lemah',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PERDARAHAN_AKTIF'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUCAT'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NADI_LEMAH')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Perdarahan Aktif'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Keseleo (disease 22)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Keseleo — bengkak sendi trauma + nyeri sendi',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_SENDI_TRAUMA'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NYERI_SENDI_TRAUMA')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Keseleo / Cedera Ligamen'),
  0.800, 0.800, 0.150, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Patah Tulang (disease 23)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Fraktur — tidak bisa gerak + nyeri sendi trauma parah',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_TIDAK_BISA_GERAK'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NYERI_SENDI_TRAUMA'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_SENDI_TRAUMA')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Patah Tulang / Fraktur'),
  0.900, 0.900, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Sengatan Lebah/Tawon (disease 24)
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Sengatan lebah — riwayat sengatan + nyeri + bengkak lokal',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SENGATAN_LEBAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_NYERI_DI_AREA_SENGATAN'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_LOKAL_SENGATAN')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sengatan Lebah / Tawon'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Sengatan lebah sistemik — sengatan + sesak napas + bengkak wajah',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SENGATAN_LEBAH'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_SESAK_NAPAS'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BENGKAK_BIBIR_WAJAH')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Sengatan Lebah / Tawon'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal');

-- Extra rules: Bali Belly with new symptoms
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Bali Belly berat — diare + tinja berdarah',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_DIARE'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_TINJA_BERDARAH')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Bali Belly (Diare Wisatawan)'),
  0.900, 0.900, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal')
  WHERE EXISTS (SELECT 1 FROM expert_symptoms WHERE kode='S_DIARE');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Bali Belly + dehidrasi — diare + tanda dehidrasi',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_DIARE'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_DEHIDRASI')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Bali Belly (Diare Wisatawan)'),
  0.850, 0.900, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal')
  WHERE EXISTS (SELECT 1 FROM expert_symptoms WHERE kode='S_DIARE');

-- Heat-related + new symptoms
INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Heat Exhaustion — lemas di panas + keringat dingin + pusing',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_LEMAS_DI_PANAS'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KERINGAT_DINGIN'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_PUSING')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Heat-Related Illness (Penyakit Terkait Panas)'),
  0.850, 0.850, 0.100, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal')
  WHERE EXISTS (SELECT 1 FROM expert_symptoms WHERE kode='S_PUSING');

INSERT INTO expert_rules (nama, premis, disease_id, bobot_cf, mb, md, kategori, status, created_by)
SELECT
  'Heat Stroke — tidak berkeringat + kulit panas + kebingungan',
  jsonb_build_array(
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_TIDAK_BERKERINGAT'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_KULIT_PANAS'),
    (SELECT symptom_id FROM expert_symptoms WHERE kode='S_BINGUNG')
  ),
  (SELECT id FROM expert_diseases WHERE nama='Heat-Related Illness (Penyakit Terkait Panas)'),
  0.950, 0.950, 0.050, 'pre_travel', 'published',
  (SELECT id FROM users WHERE email='system.seed@balihealth.internal')
  WHERE EXISTS (SELECT 1 FROM expert_symptoms WHERE kode='S_KULIT_PANAS')
    AND EXISTS (SELECT 1 FROM expert_symptoms WHERE kode='S_BINGUNG');

-- ─────────────────────────────────────────────────────────────────────────────
-- SECTION 4: New Emergency Guide Flows
-- ─────────────────────────────────────────────────────────────────────────────

-- Flow: Reaksi Alergi dan Anafilaksis
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Reaksi Alergi & Anafilaksis',
  'Alergi',
  'Panduan penanganan pertama reaksi alergi dari ringan hingga anafilaksis berat.',
  '[
    {
      "id": "allergy_start",
      "title": "Reaksi Alergi",
      "instruction": "Apakah ada tanda-tanda ANAFILAKSIS? (sesak napas, bengkak wajah/lidah, kehilangan kesadaran, tekanan darah turun drastis)",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Ya — ada tanda anafilaksis", "next_id": "anaphylaxis_severe", "variant": "yes"},
        {"label": "Tidak — hanya gatal / urtikaria", "next_id": "allergy_mild", "variant": "no"}
      ]
    },
    {
      "id": "anaphylaxis_severe",
      "title": "ANAFILAKSIS — Darurat!",
      "instruction": "1. Hubungi 119/ambulans SEGERA.\n2. Baringkan korban, tinggikan kaki (kecuali jika sesak napas).\n3. Jika tersedia, gunakan EpiPen / epinefrin auto-injector di paha luar.\n4. Jika tidak ada denyut nadi dan napas, mulai CPR.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah sudah ada respons setelah epinefrin?", "next_id": "anaphylaxis_monitor", "variant": "neutral"}
      ]
    },
    {
      "id": "anaphylaxis_monitor",
      "title": "Monitor Pasca-Epinefrin",
      "instruction": "Epinefrin bekerja 10–20 menit. Tetap awasi korban hingga ambulans tiba. Posisi berbaring dengan kaki terangkat. Jangan biarkan berdiri atau berjalan.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Kondisi membaik", "next_id": "anaphylaxis_hospital", "variant": "yes"},
        {"label": "Kondisi memburuk / tidak ada napas", "next_id": "anaphylaxis_cpr", "variant": "no"}
      ]
    },
    {
      "id": "anaphylaxis_hospital",
      "title": "Tetap ke Rumah Sakit",
      "instruction": "Meski sudah membaik, pasien HARUS tetap ke rumah sakit untuk observasi minimal 4–8 jam. Reaksi anafilaksis bisa kambuh (biphasic reaction).",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "anaphylaxis_cpr",
      "title": "Mulai CPR",
      "instruction": "Tidak ada napas/denyut nadi? Mulai CPR:\n• 30 kompresi dada (kuat, cepat, tengah dada)\n• 2 napas buatan\n• Ulangi hingga ambulans tiba",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "allergy_mild",
      "title": "Alergi Ringan",
      "instruction": "Apakah ada bengkak pada wajah, bibir, atau tenggorokan?",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — ada bengkak wajah/bibir", "next_id": "allergy_escalate", "variant": "yes"},
        {"label": "Tidak — hanya gatal/ruam", "next_id": "allergy_treat", "variant": "no"}
      ]
    },
    {
      "id": "allergy_escalate",
      "title": "Waspada Anafilaksis",
      "instruction": "Bengkak wajah/bibir adalah tanda awal anafilaksis. Hubungi dokter atau pergi ke klinik segera. Jika sesak napas muncul, ikuti alur anafilaksis.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "allergy_treat",
      "title": "Penanganan Alergi Ringan",
      "instruction": "1. Berikan antihistamin oral (cetirizine 10mg atau loratadine 10mg).\n2. Kompres dingin pada area yang gatal.\n3. Hindari pemicu alergen.\n4. Awasi selama 30 menit — jika memburuk, segera ke klinik.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Luka Bakar Matahari (Sunburn)
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Luka Bakar Matahari (Sunburn)',
  'Kulit',
  'Panduan penanganan sunburn dari ringan hingga berat setelah paparan matahari di Bali.',
  '[
    {
      "id": "sunburn_start",
      "title": "Luka Bakar Matahari",
      "instruction": "Apakah ada lepuhan (blister) pada kulit yang terbakar?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Ya — ada lepuhan", "next_id": "sunburn_severe", "variant": "yes"},
        {"label": "Tidak — hanya merah dan nyeri", "next_id": "sunburn_mild", "variant": "no"}
      ]
    },
    {
      "id": "sunburn_mild",
      "title": "Sunburn Ringan",
      "instruction": "1. Segera pindah ke tempat teduh / ruangan ber-AC.\n2. Mandi air dingin (bukan es) selama 10–15 menit.\n3. Oleskan aloe vera gel atau pelembap yang mengandung aloe.\n4. Minum banyak air untuk mencegah dehidrasi.\n5. Konsumsi ibuprofen/paracetamol untuk nyeri.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah demam muncul?", "next_id": "sunburn_fever", "variant": "neutral"}
      ]
    },
    {
      "id": "sunburn_severe",
      "title": "Sunburn Berat — Ada Lepuhan",
      "instruction": "1. JANGAN pecahkan lepuhan — risiko infeksi.\n2. Tutupi dengan kain bersih yang lembab.\n3. Kunjungi klinik/dokter untuk penanganan luka yang tepat.\n4. Minum banyak air.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ada demam > 38°C atau pusing berat", "next_id": "sunburn_emergency", "variant": "yes"},
        {"label": "Tidak ada demam", "next_id": "sunburn_clinic", "variant": "no"}
      ]
    },
    {
      "id": "sunburn_fever",
      "title": "Cek Demam",
      "instruction": "Apakah suhu tubuh > 38°C atau ada pusing berat dan mual?",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — ada demam/pusing berat", "next_id": "sunburn_emergency", "variant": "yes"},
        {"label": "Tidak", "next_id": "sunburn_recovery", "variant": "no"}
      ]
    },
    {
      "id": "sunburn_emergency",
      "title": "Ke IGD Segera",
      "instruction": "Demam + sunburn berat bisa menyebabkan heat stroke. Segera ke IGD terdekat. Jaga tubuh tetap dingin selama perjalanan (kompres basah).",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "sunburn_clinic",
      "title": "Ke Klinik",
      "instruction": "Pergi ke klinik untuk perawatan luka yang tepat dan resep obat topikal. Hindari matahari langsung selama 1–2 minggu.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "sunburn_recovery",
      "title": "Pemulihan di Tempat",
      "instruction": "Istirahat di ruangan sejuk. Oleskan aloe vera setiap 2–3 jam. Hindari sabun keras. Kulit akan mengelupas dalam 3–5 hari — biarkan mengelupas sendiri.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Mabuk Perjalanan
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Mabuk Perjalanan (Motion Sickness)',
  'Umum',
  'Panduan penanganan mabuk perjalanan saat berwisata darat, laut, atau udara di Bali.',
  '[
    {
      "id": "motion_start",
      "title": "Mabuk Perjalanan",
      "instruction": "Apakah korban masih di dalam kendaraan / kapal yang bergerak?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Ya — masih dalam kendaraan", "next_id": "motion_in_vehicle", "variant": "yes"},
        {"label": "Tidak — sudah berhenti", "next_id": "motion_stopped", "variant": "no"}
      ]
    },
    {
      "id": "motion_in_vehicle",
      "title": "Penanganan Dalam Kendaraan",
      "instruction": "1. Pindah ke kursi dengan gerakan minimal (tengah bus/kapal, dekat sayap di pesawat).\n2. Pandang ke cakrawala atau titik tetap di kejauhan.\n3. Hindari membaca atau menatap layar ponsel.\n4. Buka jendela untuk udara segar jika bisa.\n5. Bernapas dalam-dalam dan perlahan.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah membaik?", "next_id": "motion_better", "variant": "yes"},
        {"label": "Tidak membaik / makin mual", "next_id": "motion_worse", "variant": "no"}
      ]
    },
    {
      "id": "motion_stopped",
      "title": "Istirahat Setelah Berhenti",
      "instruction": "1. Duduk atau berbaring di tempat sejuk.\n2. Minum air putih sedikit-sedikit.\n3. Makan camilan ringan (biskuit) jika perut terasa kosong.\n4. Hindari bau menyengat.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Pulih dalam 30 menit", "next_id": "motion_recovered", "variant": "yes"},
        {"label": "Tidak pulih / muntah terus", "next_id": "motion_severe", "variant": "no"}
      ]
    },
    {
      "id": "motion_better",
      "title": "Kondisi Membaik",
      "instruction": "Teruskan perjalanan dengan posisi yang lebih nyaman. Pertimbangkan obat mabuk perjalanan (dimenhydrinate) untuk perjalanan selanjutnya — konsumsi 30 menit sebelum berangkat.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "motion_worse",
      "title": "Makin Parah",
      "instruction": "Hentikan kendaraan jika aman dan memungkinkan. Turun dan hirup udara segar. Berikan obat mabuk jika tersedia.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "motion_recovered",
      "title": "Pulih",
      "instruction": "Istirahat cukup sebelum melanjutkan perjalanan. Makan sesuatu yang ringan. Untuk perjalanan berikutnya, konsumsi obat antimabuk 30 menit sebelum berangkat.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "motion_severe",
      "title": "Mabuk Parah — Perlu Bantuan",
      "instruction": "Jika muntah terus-menerus dan tidak bisa minum cairan selama > 1 jam, cari fasilitas kesehatan terdekat untuk rehidrasi intravena (cairan infus).",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Gigitan Ular
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Gigitan Ular',
  'Darurat',
  'Panduan penanganan pertama gigitan ular di Bali. Anggap semua gigitan ular berbahaya hingga terbukti sebaliknya.',
  '[
    {
      "id": "snake_start",
      "title": "Gigitan Ular",
      "instruction": "PENTING: Jauhi ular! Jangan mencoba menangkap atau membunuh ular. Catat ciri-ciri ular (warna, ukuran, pola) dari jarak aman atau foto jika ada.",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Lanjutkan penanganan pertama", "next_id": "snake_first_aid", "variant": "neutral"}
      ]
    },
    {
      "id": "snake_first_aid",
      "title": "Penanganan Pertama",
      "instruction": "1. Tenangkan korban — panik mempercepat penyebaran racun.\n2. Imobilisasi anggota yang digigit — posisi di BAWAH level jantung.\n3. Lepaskan perhiasan / jam tangan dari area yang digigit (bengkak akan muncul).\n4. JANGAN isap luka, JANGAN iris luka, JANGAN tourniquet ketat.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ada tanda-tanda gigitan berbisa?", "next_id": "snake_check_venom", "variant": "neutral"}
      ]
    },
    {
      "id": "snake_check_venom",
      "title": "Cek Tanda Gigitan Berbisa",
      "instruction": "Apakah ada salah satu dari ini dalam 30 menit setelah gigitan?\n• Bengkak cepat menyebar\n• Nyeri hebat di area gigitan\n• Pucat, pusing, mual\n• Kesemutan / mati rasa",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — ada tanda berbisa", "next_id": "snake_venomous", "variant": "yes"},
        {"label": "Tidak — tidak ada tanda berbisa", "next_id": "snake_monitor", "variant": "no"}
      ]
    },
    {
      "id": "snake_venomous",
      "title": "DARURAT — Kemungkinan Berbisa",
      "instruction": "SEGERA ke Rumah Sakit yang memiliki antivenom:\n• RSUP Sanglah Denpasar: (0361) 227911\n• RS Surya Husadha\n\nSelama perjalanan: imobilisasi penuh, pasang perban tekan (compression bandage) — bukan tourniquet. Awasi pernapasan.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "snake_monitor",
      "title": "Monitor Selama 4 Jam",
      "instruction": "Meski tidak ada tanda langsung, tetap awasi selama minimal 4 jam. Racun beberapa ular muncul terlambat. Segera ke klinik untuk pemeriksaan dan pencucian luka.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Sengatan Hewan Laut
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Sengatan Hewan Laut',
  'Darurat',
  'Panduan penanganan sengatan ubur-ubur, tusukan bulu babi, dan hewan laut beracun lain di Bali.',
  '[
    {
      "id": "marine_start",
      "title": "Sengatan Hewan Laut",
      "instruction": "Jenis cedera apa yang terjadi?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Sengatan ubur-ubur (jellyfish)", "next_id": "jellyfish_sting", "variant": "yes"},
        {"label": "Tusukan duri bulu babi (sea urchin)", "next_id": "urchin_sting", "variant": "no"},
        {"label": "Hewan laut lain / tidak tahu", "next_id": "marine_unknown", "variant": "neutral"}
      ]
    },
    {
      "id": "jellyfish_sting",
      "title": "Sengatan Ubur-ubur",
      "instruction": "1. Keluar dari air dengan hati-hati.\n2. JANGAN digosok — racun bisa menyebar.\n3. Siram area sengatan dengan CUKA (vinegar) selama 30 detik.\n4. Rendam atau alirkan AIR PANAS (40–45°C, bukan mendidih) selama 20 menit.\n5. Cabut tentakel yang tertinggal dengan pinset atau kartu — JANGAN dengan tangan.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah ada reaksi alergi (sesak napas, bengkak)?", "next_id": "marine_allergy", "variant": "yes"},
        {"label": "Tidak ada reaksi sistemik", "next_id": "jellyfish_monitor", "variant": "no"}
      ]
    },
    {
      "id": "urchin_sting",
      "title": "Tusukan Duri Bulu Babi",
      "instruction": "1. Rendam area yang tertusuk dalam air panas (40–45°C) selama 30–90 menit — membantu melarutkan duri kecil dan mengurangi nyeri.\n2. Duri kecil biasanya larut sendiri — JANGAN coba cabut paksa.\n3. Duri besar atau di area sensitif: ke klinik untuk pencabutan medis.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Duri terlihat jelas dan bisa dicabut?", "next_id": "urchin_remove", "variant": "yes"},
        {"label": "Duri kecil/dalam atau di wajah/tangan", "next_id": "urchin_clinic", "variant": "no"}
      ]
    },
    {
      "id": "urchin_remove",
      "title": "Pencabutan Duri",
      "instruction": "Gunakan pinset steril (dibakar dulu ujungnya) untuk mencabut duri yang terlihat. Setelah dicabut, bersihkan dengan antiseptik dan tutup dengan perban.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "urchin_clinic",
      "title": "Ke Klinik",
      "instruction": "Pergi ke klinik untuk pencabutan duri yang aman. Duri yang tersisa bisa menyebabkan infeksi atau granuloma jika dibiarkan.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "marine_unknown",
      "title": "Hewan Laut Tidak Diketahui",
      "instruction": "1. Keluar dari air segera.\n2. Bilas area dengan air laut (bukan air tawar untuk sengatan ubur-ubur).\n3. Rendam dalam air panas 40–45°C selama 20 menit.\n4. Segera ke klinik untuk identifikasi dan penanganan.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "marine_allergy",
      "title": "DARURAT — Reaksi Anafilaksis",
      "instruction": "Sesak napas atau bengkak wajah setelah sengatan laut = ANAFILAKSIS.\nHubungi 119 segera. Posisi berbaring, kaki terangkat. Gunakan EpiPen jika tersedia.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "jellyfish_monitor",
      "title": "Monitor 30 Menit",
      "instruction": "Awasi reaksi alergi selama 30 menit. Konsumsi antihistamin oral jika tersedia. Kunjungi klinik jika nyeri tidak berkurang dalam 1 jam.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Perdarahan dan Luka
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Perdarahan & Penanganan Luka',
  'Darurat',
  'Panduan penanganan perdarahan aktif dan luka terbuka dari ringan hingga berat.',
  '[
    {
      "id": "bleed_start",
      "title": "Penilaian Perdarahan",
      "instruction": "Seberapa parah perdarahannya?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Luka kecil / lecet ringan", "next_id": "bleed_minor", "variant": "no"},
        {"label": "Perdarahan deras / tidak berhenti", "next_id": "bleed_major", "variant": "yes"}
      ]
    },
    {
      "id": "bleed_minor",
      "title": "Luka Kecil / Ringan",
      "instruction": "1. Cuci tangan sebelum menangani luka.\n2. Bilas luka dengan air bersih mengalir selama 5 menit.\n3. Tekan dengan kain bersih hingga darah berhenti.\n4. Oleskan antiseptik (Betadine).\n5. Tutup dengan plester / perban.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah ada tanda infeksi (merah menyebar, bernanah)?", "next_id": "bleed_infection", "variant": "yes"},
        {"label": "Tidak ada tanda infeksi", "next_id": "bleed_healed", "variant": "no"}
      ]
    },
    {
      "id": "bleed_major",
      "title": "Perdarahan Deras",
      "instruction": "1. Tekan KUAT dengan kain bersih / perban langsung pada luka.\n2. JANGAN angkat kain penekan — tambahkan di atasnya jika perlu.\n3. Tinggikan bagian yang berdarah di atas level jantung jika bisa.\n4. Hubungi 119 atau pergi ke IGD segera.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah perdarahan berhenti dalam 10 menit?", "next_id": "bleed_stop_check", "variant": "neutral"}
      ]
    },
    {
      "id": "bleed_stop_check",
      "title": "Cek Perdarahan",
      "instruction": "Apakah perdarahan sudah berhenti atau sangat berkurang setelah penekanan 10 menit?",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — sudah berhenti", "next_id": "bleed_stopped", "variant": "yes"},
        {"label": "Tidak — masih deras", "next_id": "bleed_emergency", "variant": "no"}
      ]
    },
    {
      "id": "bleed_stopped",
      "title": "Perdarahan Terhenti",
      "instruction": "Pertahankan tekanan dan perban. Jangan lepas perban. Pergi ke klinik untuk pembersihan luka dan kemungkinan jahitan jika luka dalam.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "bleed_emergency",
      "title": "DARURAT — Perdarahan Tidak Terhenti",
      "instruction": "Jika perdarahan tidak terhenti setelah 15 menit penekanan kuat:\n1. Pertahankan tekanan — jangan lepas.\n2. Hubungi 119 segera.\n3. Awasi tanda syok: pucat, nadi lemah, pusing berat.\n4. Jika syok: baringkan, tinggikan kaki.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "bleed_infection",
      "title": "Tanda Infeksi",
      "instruction": "Tanda infeksi pada luka memerlukan penanganan dokter. Pergi ke klinik untuk pembersihan luka, antibiotik, dan kemungkinan suntikan tetanus.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "bleed_healed",
      "title": "Luka Sembuh Normal",
      "instruction": "Ganti perban setiap hari atau saat kotor. Jaga area tetap bersih dan kering. Jika luka dalam > 1 cm, pertimbangkan ke klinik untuk jahitan.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Syok dan Sinkop
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Syok & Sinkop (Pingsan)',
  'Darurat',
  'Panduan penanganan syok dan kehilangan kesadaran sementara (sinkop).',
  '[
    {
      "id": "shock_start",
      "title": "Penilaian Awal",
      "instruction": "Apakah korban masih sadar?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Ya — sadar tapi sangat lemas/pucat", "next_id": "shock_conscious", "variant": "yes"},
        {"label": "Tidak — tidak sadar", "next_id": "shock_unconscious", "variant": "no"}
      ]
    },
    {
      "id": "shock_conscious",
      "title": "Korban Sadar — Tanda Syok",
      "instruction": "Apakah ada tanda syok? (pilih yang sesuai)\n• Kulit pucat, dingin, lembab\n• Nadi cepat dan lemah\n• Napas cepat dan dangkal\n• Sangat lemas / bingung",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ada tanda syok", "next_id": "shock_treat", "variant": "yes"},
        {"label": "Hanya lemas / hampir pingsan", "next_id": "syncope_treat", "variant": "no"}
      ]
    },
    {
      "id": "shock_treat",
      "title": "Penanganan Syok",
      "instruction": "1. Hubungi 119 SEGERA.\n2. Baringkan korban dengan kaki ditinggikan 20–30 cm (kecuali cedera kepala/tulang belakang).\n3. Jaga kehangatan dengan selimut.\n4. JANGAN beri makan/minum.\n5. Atasi penyebab: hentikan perdarahan jika ada.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "syncope_treat",
      "title": "Penanganan Sinkop / Hampir Pingsan",
      "instruction": "1. Baringkan dan tinggikan kaki 30 cm.\n2. Pindah ke tempat sejuk / teduh.\n3. Longgarkan pakaian ketat.\n4. Kipasi atau kompres dahi dengan air dingin.\n5. Beri minum perlahan saat sudah sadar.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Pulih dalam 5 menit?", "next_id": "syncope_check", "variant": "neutral"}
      ]
    },
    {
      "id": "syncope_check",
      "title": "Cek Pemulihan",
      "instruction": "Apakah korban pulih sepenuhnya dalam 5 menit?",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — sudah pulih", "next_id": "syncope_recovered", "variant": "yes"},
        {"label": "Tidak — masih tidak sadar / bingung", "next_id": "shock_unconscious", "variant": "no"}
      ]
    },
    {
      "id": "syncope_recovered",
      "title": "Pulih dari Sinkop",
      "instruction": "Tetap istirahat 15–30 menit. Minum cairan (air atau minuman elektrolit). Hindari berdiri tiba-tiba. Jika sinkop berulang, konsultasikan ke dokter.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "shock_unconscious",
      "title": "DARURAT — Tidak Sadar",
      "instruction": "Gunakan DRSABC:\n• D: Danger — pastikan area aman\n• R: Response — panggil nama, cubit bahu\n• S: Send — hubungi 119 SEKARANG\n• A: Airway — miringkan kepala, buka mulut\n• B: Breathing — cek napas (10 detik)\n• C: CPR — jika tidak ada napas, mulai 30 kompresi + 2 napas buatan",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ada napas — posisi pemulihan", "next_id": "shock_recovery_position", "variant": "yes"},
        {"label": "Tidak ada napas — mulai CPR", "next_id": "shock_cpr", "variant": "no"}
      ]
    },
    {
      "id": "shock_recovery_position",
      "title": "Posisi Pemulihan",
      "instruction": "Posisi lateral stabil (recovery position):\n1. Lutut ditekuk sebagai sandaran.\n2. Satu tangan di bawah kepala.\n3. Kepala sedikit menengadah agar jalan napas terbuka.\n4. Awasi napas hingga bantuan tiba.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "shock_cpr",
      "title": "Mulai CPR",
      "instruction": "Tidak ada napas? Mulai CPR:\n1. Posisi korban telentang di permukaan keras.\n2. 30 kompresi dada: tengah dada, kedalaman 5–6 cm, kecepatan 100–120x/menit.\n3. 2 napas buatan (miringkan kepala, angkat dagu, tiup perlahan 1 detik).\n4. Ulangi 30:2 hingga ambulans tiba atau korban bernapas kembali.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Keseleo dan Patah Tulang
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Keseleo & Patah Tulang',
  'Cedera',
  'Panduan penanganan pertama cedera muskuloskeletal akibat trauma fisik saat berwisata.',
  '[
    {
      "id": "musculo_start",
      "title": "Penilaian Cedera",
      "instruction": "Jenis cedera apa yang terjadi?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Keseleo / terkilir (masih bisa sedikit digerakkan)", "next_id": "sprain_treat", "variant": "no"},
        {"label": "Kemungkinan patah tulang (tidak bisa digerakkan / sangat nyeri)", "next_id": "fracture_check", "variant": "yes"}
      ]
    },
    {
      "id": "sprain_treat",
      "title": "Penanganan Keseleo — RICE",
      "instruction": "Terapkan metode RICE:\n• R — Rest: Istirahatkan, jangan dipaksa dipakai\n• I — Ice: Kompres es (bungkus kain) 20 menit, 3–4x sehari\n• C — Compression: Balut dengan perban elastis (tidak terlalu ketat)\n• E — Elevation: Tinggikan bagian yang cedera di atas jantung",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah nyeri sangat hebat atau tidak bisa berjalan sama sekali?", "next_id": "sprain_severe", "variant": "yes"},
        {"label": "Nyeri ringan-sedang, bisa berjalan pelan", "next_id": "sprain_monitor", "variant": "no"}
      ]
    },
    {
      "id": "sprain_severe",
      "title": "Kemungkinan Fraktur",
      "instruction": "Keseleo yang sangat nyeri dengan bengkak besar mungkin sebenarnya patah tulang. Pergi ke klinik/IGD untuk X-ray.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "sprain_monitor",
      "title": "Pemantauan Keseleo",
      "instruction": "Istirahat 48–72 jam. Hindari olahraga berat. Konsumsi ibuprofen untuk nyeri dan bengkak. Jika tidak membaik dalam 3 hari, periksa ke dokter.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "fracture_check",
      "title": "Tanda Patah Tulang",
      "instruction": "Apakah ada tanda ini?\n• Bentuk anggota tubuh tidak normal / bengkok\n• Tulang menembus kulit (fraktur terbuka)\n• Nyeri sangat hebat saat ditekan ringan\n• Tidak bisa digerakkan sama sekali",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — ada tanda patah tulang", "next_id": "fracture_treat", "variant": "yes"},
        {"label": "Tidak yakin — mungkin keseleo parah", "next_id": "sprain_severe", "variant": "no"}
      ]
    },
    {
      "id": "fracture_treat",
      "title": "Penanganan Patah Tulang",
      "instruction": "1. JANGAN paksa luruskan tulang.\n2. Imobilisasi: gunakan papan / majalah sebagai bidai, ikat di atas dan bawah fraktur.\n3. Fraktur terbuka (tulang tembus kulit): tutupi dengan kain bersih, JANGAN tekan tulang.\n4. Segera ke IGD untuk penanganan.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah ada cedera kepala atau tulang belakang?", "next_id": "fracture_spine", "variant": "yes"},
        {"label": "Tidak ada cedera kepala/tulang belakang", "next_id": "fracture_transport", "variant": "no"}
      ]
    },
    {
      "id": "fracture_spine",
      "title": "DARURAT — Curiga Cedera Tulang Belakang",
      "instruction": "JANGAN gerakkan korban! Cedera tulang belakang yang bergerak bisa menyebabkan kelumpuhan permanen.\nHubungi 119 segera. Stabilkan kepala dan leher dengan tangan hingga ambulans tiba.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "fracture_transport",
      "title": "Transport ke IGD",
      "instruction": "Dengan bidai terpasang, transportasikan ke IGD terdekat untuk X-ray dan penanganan lebih lanjut. Jaga area fraktur tetap imobil selama perjalanan.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Keracunan Metanol
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Keracunan Arak Oplosan (Metanol)',
  'Darurat',
  'Panduan darurat keracunan metanol akibat konsumsi arak oplosan. Sangat berbahaya — memerlukan penanganan medis segera.',
  '[
    {
      "id": "methanol_start",
      "title": "Dugaan Keracunan Metanol",
      "instruction": "Apakah korban baru saja mengonsumsi arak lokal (arak Bali) atau minuman alkohol tidak bermerek?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Ya — baru konsumsi arak lokal", "next_id": "methanol_check", "variant": "yes"},
        {"label": "Tidak yakin", "next_id": "methanol_symptoms", "variant": "no"}
      ]
    },
    {
      "id": "methanol_check",
      "title": "Cek Gejala Metanol",
      "instruction": "Apakah ada salah satu gejala ini?\n• Penglihatan kabur atau buta sebagian\n• Nyeri kepala hebat\n• Mual dan muntah berulang\n• Sempoyongan / bingung\n• Sakit perut hebat",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — ada gejala di atas", "next_id": "methanol_emergency", "variant": "yes"},
        {"label": "Tidak ada gejala", "next_id": "methanol_monitor", "variant": "no"}
      ]
    },
    {
      "id": "methanol_symptoms",
      "title": "Gejala Keracunan Metanol",
      "instruction": "Gejala keracunan metanol muncul 6–24 jam setelah konsumsi:\n• Fase awal: mirip mabuk alkohol biasa\n• Fase kritis: gangguan penglihatan, muntah hebat, nyeri perut",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ada gejala penglihatan / muntah hebat", "next_id": "methanol_emergency", "variant": "yes"},
        {"label": "Belum ada gejala kritis", "next_id": "methanol_monitor", "variant": "no"}
      ]
    },
    {
      "id": "methanol_emergency",
      "title": "DARURAT MEDIS",
      "instruction": "Keracunan metanol adalah DARURAT MEDIS yang mengancam jiwa dan penglihatan.\n\n1. Hubungi 119 SEGERA.\n2. Bawa ke IGD terdekat — informasikan bahwa korban mengonsumsi arak lokal.\n3. Antidot: etanol IV atau fomepizole — hanya tersedia di RS.\n4. JANGAN beri minum atau makan.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "methanol_monitor",
      "title": "Monitor Ketat",
      "instruction": "Meski belum ada gejala kritis, segera ke klinik/dokter untuk pemeriksaan darah. Gangguan penglihatan akibat metanol bisa permanen jika terlambat ditangani. Awasi selama 24 jam.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);

-- Flow: Sengatan Lebah / Tawon
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Sengatan Lebah & Tawon',
  'Alergi',
  'Panduan penanganan sengatan lebah atau tawon, dari reaksi lokal hingga anafilaksis.',
  '[
    {
      "id": "bee_start",
      "title": "Sengatan Lebah / Tawon",
      "instruction": "Apakah ada sengat yang tertinggal di kulit?",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Ya — ada sengat tertinggal (lebah madu)", "next_id": "bee_remove_stinger", "variant": "yes"},
        {"label": "Tidak — tidak ada sengat tertinggal (tawon)", "next_id": "bee_treat", "variant": "no"}
      ]
    },
    {
      "id": "bee_remove_stinger",
      "title": "Keluarkan Sengat",
      "instruction": "JANGAN pencet sengat — akan mempercepat penyebaran racun.\nKeluarkan sengat dengan cara DIKERIK menggunakan kuku, kartu, atau pinset (jepit di pangkal sengat, tarik cepat).",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Sengat sudah dikeluarkan — lanjutkan perawatan", "next_id": "bee_treat", "variant": "neutral"}
      ]
    },
    {
      "id": "bee_treat",
      "title": "Perawatan Sengatan",
      "instruction": "1. Cuci area dengan sabun dan air bersih.\n2. Kompres es (bungkus kain) selama 15–20 menit.\n3. Konsumsi antihistamin oral (cetirizine) untuk gatal.\n4. Ibuprofen untuk nyeri.",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Apakah ada reaksi sistemik dalam 30 menit?", "next_id": "bee_systemic", "variant": "neutral"}
      ]
    },
    {
      "id": "bee_systemic",
      "title": "Cek Reaksi Sistemik",
      "instruction": "Apakah ada salah satu gejala ini dalam 30 menit setelah sengatan?\n• Urtikaria / biduran menyebar\n• Bengkak di lokasi jauh dari sengatan\n• Sesak napas\n• Pusing atau hampir pingsan\n• Bengkak wajah atau tenggorokan",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Ya — ada reaksi sistemik / alergi", "next_id": "bee_anaphylaxis", "variant": "yes"},
        {"label": "Tidak — hanya reaksi lokal", "next_id": "bee_local", "variant": "no"}
      ]
    },
    {
      "id": "bee_anaphylaxis",
      "title": "DARURAT — Anafilaksis",
      "instruction": "Reaksi sistemik setelah sengatan = kemungkinan anafilaksis.\n1. Hubungi 119 SEGERA.\n2. Gunakan EpiPen jika tersedia.\n3. Posisi: berbaring dengan kaki terangkat (atau duduk tegak jika sesak napas).\n4. Mulai CPR jika tidak sadar dan tidak ada napas.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    },
    {
      "id": "bee_local",
      "title": "Reaksi Lokal Normal",
      "instruction": "Bengkak dan kemerahan lokal pada area sengatan adalah normal dan akan membaik dalam 24–48 jam. Hindari menggaruk. Terus kompres dingin dan antihistamin jika perlu.",
      "image_url": "",
      "is_entry": false,
      "choices": []
    }
  ]'::jsonb
);
