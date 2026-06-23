-- Roles
INSERT INTO roles (nama_role, deskripsi) VALUES
    ('traveler', 'Wisatawan pengguna aplikasi mobile'),
    ('nurse',    'Perawat BMTA yang melayani wisatawan'),
    ('admin',    'Administrator sistem')
ON CONFLICT (nama_role) DO NOTHING;

-- Destinations (kabupaten/kota Bali)
INSERT INTO destinations (nama_daerah) VALUES
    ('Kota Denpasar'),
    ('Kabupaten Badung'),
    ('Kabupaten Gianyar'),
    ('Kabupaten Tabanan'),
    ('Kabupaten Buleleng'),
    ('Kabupaten Karangasem'),
    ('Kabupaten Klungkung'),
    ('Kabupaten Bangli'),
    ('Kabupaten Jembrana')
ON CONFLICT DO NOTHING;

-- Medical facilities (BMTA-affiliated, Bali)
-- destination_id references: 1=Denpasar, 2=Badung, 3=Gianyar
INSERT INTO medical_facilities
    (destination_id, nama, kategori, alamat, latitude, longitude, kontak, jam_operasional)
VALUES
    (
        1,
        'RSUP Prof. Dr. I.G.N.G. Ngoerah',
        'Rumah Sakit Umum Pemerintah',
        'Jl. Diponegoro No.1, Dauh Puri Klod, Denpasar Barat',
        -8.661700, 115.213700,
        '(0361) 227911',
        '24 Jam'
    ),
    (
        1,
        'RS Mata Bali Mandara',
        'Rumah Sakit Khusus Mata',
        'Jl. Angsoka No.8, Tonja, Denpasar Utara',
        -8.636400, 115.218300,
        '(0361) 261559',
        'Senin–Sabtu 07.00–21.00'
    ),
    (
        1,
        'RSUD Bali Mandara',
        'Rumah Sakit Umum Daerah',
        'Jl. By Pass Ngurah Rai No.548, Sanur Kaja, Denpasar Selatan',
        -8.720500, 115.249000,
        '(0361) 4710808',
        '24 Jam'
    ),
    (
        1,
        'RS Prima Medika',
        'Rumah Sakit Swasta',
        'Jl. Pulau Serangan No.9X, Tonja, Denpasar Utara',
        -8.634800, 115.220100,
        '(0361) 436388',
        '24 Jam'
    ),
    (
        1,
        'RS Kasih Ibu',
        'Rumah Sakit Swasta',
        'Jl. Teuku Umar No.120, Dauh Puri Klod, Denpasar Barat',
        -8.667200, 115.203100,
        '(0361) 223036',
        '24 Jam'
    ),
    (
        1,
        'RS Surya Husadha',
        'Rumah Sakit Swasta',
        'Jl. Pulau Serangan No.1-7, Tonja, Denpasar Utara',
        -8.635700, 115.219400,
        '(0361) 233787',
        '24 Jam'
    ),
    (
        1,
        'RS Wangaya (RSUD Wangaya)',
        'Rumah Sakit Umum Daerah',
        'Jl. Kartini No.133, Dauh Puri, Denpasar Barat',
        -8.654900, 115.210300,
        '(0361) 222141',
        '24 Jam'
    ),
    (
        2,
        'BIMC Hospital Kuta',
        'Klinik Internasional',
        'Jl. By Pass Ngurah Rai No.100X, Kuta',
        -8.721300, 115.177600,
        '(0361) 761263',
        '24 Jam'
    ),
    (
        2,
        'BIMC Hospital Nusa Dua',
        'Klinik Internasional',
        'Kawasan ITDC Nusa Dua Blok D, Benoa, Kuta Selatan',
        -8.802200, 115.225600,
        '(0361) 3000911',
        '24 Jam'
    ),
    (
        2,
        'Siloam Hospitals Bali',
        'Rumah Sakit Swasta Internasional',
        'Jl. Sunset Road No.818, Seminyak, Kuta',
        -8.696900, 115.163200,
        '(0361) 779900',
        '24 Jam'
    ),
    (
        2,
        'RS Dharma Yadnya',
        'Rumah Sakit Swasta',
        'Jl. By Pass Ngurah Rai No.19, Kesiman, Denpasar Timur',
        -8.680300, 115.252700,
        '(0361) 466588',
        '24 Jam'
    ),
    (
        2,
        'Klinik Bali Med (International SOS)',
        'Klinik Internasional',
        'Jl. Bypass Ngurah Rai No.505X, Jimbaran, Kuta Selatan',
        -8.770200, 115.173100,
        '(0361) 710505',
        '24 Jam'
    ),
    (
        3,
        'RSUD Sanjiwani Gianyar',
        'Rumah Sakit Umum Daerah',
        'Jl. Ciung Wanara No.2, Gianyar',
        -8.537600, 115.329200,
        '(0361) 943049',
        '24 Jam'
    ),
    (
        3,
        'RS Bumi Waras',
        'Rumah Sakit Swasta',
        'Jl. Raya Ubud No.36, Ubud, Gianyar',
        -8.507000, 115.262000,
        '(0361) 975255',
        'Senin–Sabtu 08.00–20.00'
    ),
    (
        2,
        'Klinik Seminyak Square Medical',
        'Klinik Swasta',
        'Jl. Raya Seminyak, Seminyak, Kuta',
        -8.693100, 115.162500,
        '(0361) 735600',
        '24 Jam'
    ),
    (
        1,
        'RS Bali Royal (BROS)',
        'Rumah Sakit Swasta',
        'Jl. Letda Tantular No.6, Renon, Denpasar Selatan',
        -8.677500, 115.234300,
        '(0361) 244545',
        '24 Jam'
    ),
    (
        2,
        'Klinik Legian Medical Clinic',
        'Klinik Swasta',
        'Jl. Padma Utara No.13, Legian, Kuta',
        -8.707400, 115.167500,
        '(0361) 758503',
        '24 Jam'
    );
