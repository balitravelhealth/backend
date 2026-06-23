-- Rollback Migration 000025: Remove SOP Knowledge Base data

-- Remove new emergency guide flows (by title)
DELETE FROM emergency_guide_flows WHERE title IN (
  'Reaksi Alergi & Anafilaksis',
  'Luka Bakar Matahari (Sunburn)',
  'Mabuk Perjalanan (Motion Sickness)',
  'Gigitan Ular',
  'Sengatan Hewan Laut',
  'Perdarahan & Penanganan Luka',
  'Syok & Sinkop (Pingsan)',
  'Keseleo & Patah Tulang',
  'Keracunan Arak Oplosan (Metanol)',
  'Sengatan Lebah & Tawon'
);

-- Remove rules associated with new diseases (IDs >= 27 from this migration)
DELETE FROM expert_rules WHERE disease_id IN (
  SELECT id FROM expert_diseases WHERE nama IN (
    'Reaksi Alergi Ringan', 'Anafilaksis', 'Sunburn (Luka Bakar Matahari)',
    'Mabuk Perjalanan (Motion Sickness)', 'Keracunan Metanol (Arak Oplosan)',
    'Syok (Shock)', 'Sinkop (Pingsan)', 'Gigitan Ular', 'Sengatan Hewan Laut',
    'Reaksi Gigitan Serangga', 'Perdarahan Aktif', 'Keseleo / Cedera Ligamen',
    'Patah Tulang / Fraktur', 'Sengatan Lebah / Tawon'
  )
);

-- Remove extra rules added for existing diseases using new symptoms
DELETE FROM expert_rules WHERE nama IN (
  'Bali Belly berat — diare + tinja berdarah',
  'Bali Belly + dehidrasi — diare + tanda dehidrasi',
  'Heat Exhaustion — lemas di panas + keringat dingin + pusing',
  'Heat Stroke — tidak berkeringat + kulit panas + kebingungan'
);

-- Remove new diseases
DELETE FROM expert_diseases WHERE nama IN (
  'Reaksi Alergi Ringan', 'Anafilaksis', 'Sunburn (Luka Bakar Matahari)',
  'Mabuk Perjalanan (Motion Sickness)', 'Keracunan Metanol (Arak Oplosan)',
  'Syok (Shock)', 'Sinkop (Pingsan)', 'Gigitan Ular', 'Sengatan Hewan Laut',
  'Reaksi Gigitan Serangga', 'Perdarahan Aktif', 'Keseleo / Cedera Ligamen',
  'Patah Tulang / Fraktur', 'Sengatan Lebah / Tawon'
);

-- Remove new symptoms
DELETE FROM expert_symptoms WHERE kode IN (
  'S_TINJA_BERDARAH', 'S_DEHIDRASI', 'S_RUAM_URTIKARIA', 'S_MATA_MERAH_GATAL',
  'S_BENGKAK_BIBIR_WAJAH', 'S_SESAK_NAPAS', 'S_KULIT_MERAH_BAKAR', 'S_LEPUHAN_KULIT',
  'S_PAPARAN_MATAHARI_LAMA', 'S_MUAL_PERJALANAN', 'S_KERINGAT_DINGIN',
  'S_PUSING_PERJALANAN', 'S_GANGGUAN_PENGLIHATAN', 'S_SEMPOYONGAN',
  'S_KONSUMSI_ALKOHOL_LOKAL', 'S_GIGITAN_ULAR', 'S_SENGATAN_UBUR',
  'S_TUSUKAN_DURI_LAUT', 'S_PINGSAN_SEMENTARA', 'S_PUCAT', 'S_NADI_LEMAH',
  'S_PERDARAHAN_AKTIF', 'S_LUKA_KECIL', 'S_BENGKAK_SENDI_TRAUMA',
  'S_TIDAK_BISA_GERAK', 'S_SENGATAN_LEBAH', 'S_BENGKAK_LOKAL_SENGATAN',
  'S_LEMAS_DI_PANAS', 'S_TIDAK_BERKERINGAT', 'S_NYERI_SENDI_TRAUMA',
  'S_DETAK_CEPAT', 'S_BENGKAK_TENGGOROKAN', 'S_NYERI_DI_AREA_SENGATAN'
);
