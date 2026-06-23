DELETE FROM medical_facilities
WHERE nama IN (
    'RSUP Prof. Dr. I.G.N.G. Ngoerah',
    'RS Mata Bali Mandara',
    'RSUD Bali Mandara',
    'RS Prima Medika',
    'RS Kasih Ibu',
    'RS Surya Husadha',
    'RS Wangaya (RSUD Wangaya)',
    'BIMC Hospital Kuta',
    'BIMC Hospital Nusa Dua',
    'Siloam Hospitals Bali',
    'RS Dharma Yadnya',
    'Klinik Bali Med (International SOS)',
    'RSUD Sanjiwani Gianyar',
    'RS Bumi Waras',
    'Klinik Seminyak Square Medical',
    'RS Bali Royal (BROS)',
    'Klinik Legian Medical Clinic'
);

DELETE FROM destinations
WHERE nama_daerah IN (
    'Kota Denpasar',
    'Kabupaten Badung',
    'Kabupaten Gianyar',
    'Kabupaten Tabanan',
    'Kabupaten Buleleng',
    'Kabupaten Karangasem',
    'Kabupaten Klungkung',
    'Kabupaten Bangli',
    'Kabupaten Jembrana'
);

DELETE FROM roles WHERE nama_role IN ('traveler', 'nurse', 'admin');
