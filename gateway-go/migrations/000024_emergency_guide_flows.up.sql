-- Decision-tree emergency guides (replaces sequential emergency_guides for flow-based UX)
-- Each flow contains an array of "nodes" in JSONB.
-- Node schema:
--   {
--     "id":          "check_responsiveness",   -- unique key within flow
--     "title":       "Check for responsiveness",
--     "instruction": "Tap shoulder and shout...",
--     "image_url":   "https://...",            -- optional illustration
--     "is_entry":    true,                     -- exactly one node per flow
--     "choices": [                             -- null = end of guide
--       {"label": "Responsive",   "next_id": "call_ambulance", "variant": "yes"},
--       {"label": "Unresponsive", "next_id": "call_119",       "variant": "no"}
--     ]
--   }

CREATE TABLE IF NOT EXISTS emergency_guide_flows (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    kategori    VARCHAR(100) NOT NULL,
    deskripsi   TEXT,
    nodes       JSONB NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_egf_kategori ON emergency_guide_flows(kategori);

-- ── Seed: Basic Life Support (BLS) ───────────────────────────────────────────
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Basic Life Support (BLS)',
  'BLS',
  'Panduan pertolongan pertama untuk korban yang tidak responsif. Ikuti langkah-langkah berikut.',
  '[
    {
      "id": "check_responsiveness",
      "title": "Periksa Respons Korban",
      "instruction": "Pastikan lokasi aman. Tepuk bahu korban dan panggil keras: \"Hei! Kamu baik-baik saja?\"\n\nLangkah:\n1) Pastikan area sekitar aman bagi Anda dan korban\n2) Berlutut di samping korban\n3) Tepuk kedua bahu korban dengan tegas\n4) Panggil keras di dekat telinga korban",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Responsif", "next_id": "call_ambulance", "variant": "yes"},
        {"label": "Tidak Responsif", "next_id": "call_119", "variant": "no"}
      ]
    },
    {
      "id": "call_ambulance",
      "title": "Hubungi Ambulans",
      "instruction": "Korban responsif — tetap pantau kondisinya.\n\n1) Tanyakan: \"Apakah kamu baik-baik saja?\"\n2) Minta orang terdekat menghubungi ambulans, atau hubungi 119 sendiri\n3) Tetap dampingi korban hingga bantuan tiba\n4) Ikuti petunjuk petugas saat tiba",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"},
        {"label": "Ada Tersedak? Lanjut →", "next_id": null, "variant": "yes"}
      ]
    },
    {
      "id": "call_119",
      "title": "Panggil 119 & Ambil AED",
      "instruction": "Korban TIDAK responsif — tindak segera!\n\n1) Teriak minta tolong di sekitar Anda\n2) Tugaskan satu orang menghubungi 119\n3) Tugaskan orang lain mencari AED terdekat\n4) Siapkan diri untuk melakukan CPR",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Cari AED di Peta", "next_id": "aed_map", "variant": "neutral"},
        {"label": "Lanjut Periksa Napas", "next_id": "check_breathing", "variant": "yes"}
      ]
    },
    {
      "id": "aed_map",
      "title": "Fasilitas Medis & AED Terdekat",
      "instruction": "Gunakan fitur pencarian fasilitas di aplikasi untuk menemukan AED atau rumah sakit terdekat.\n\nNomor darurat:\n• 119 — Ambulans\n• 110 — Polisi\n• 112 — Darurat universal",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Lanjut Periksa Napas", "next_id": "check_breathing", "variant": "yes"}
      ]
    },
    {
      "id": "check_breathing",
      "title": "Periksa Pernapasan",
      "instruction": "Lihat apakah dada dan perut naik-turun selama minimal 5 detik (tidak lebih dari 10 detik).\n\nPerhatikan:\n• Jika TIDAK bernapas atau bernapas abnormal (terengah-engah/gasping) → tanda henti jantung\n• Gasping bukan napas normal",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Bernapas Normal", "next_id": "recovery_position", "variant": "yes"},
        {"label": "Tidak Bernapas", "next_id": "cpr_adult", "variant": "no"}
      ]
    },
    {
      "id": "recovery_position",
      "title": "Posisi Pemulihan (Recovery Position)",
      "instruction": "Korban bernapas — jangan biarkan telentang!\n\n1) JANGAN biarkan korban telentang\n2) Angkat rahang bawah dan gulingkan ke samping\n3) Posisikan agar saluran napas tetap terbuka\n4) Pantau pernapasan terus-menerus\n5) Ikuti petunjuk EMT/EMS saat tiba",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"}
      ]
    },
    {
      "id": "cpr_adult",
      "title": "CPR Dewasa",
      "instruction": "Tekan 5–6 cm sedalam mungkin dengan ritme beat.\n\n1) Baringkan korban telentang di permukaan keras\n2) Jika tidak bernapas — letakkan kedua telapak tangan di tengah dada (sternum)\n3) Lakukan kompresi dada 30 kali dengan ritme beat (100–120x/menit)\n4) Berikan 2 napas bantuan: tutup hidung, angkat dagu, tiup hingga dada terangkat\n5) Ulangi siklus 30 kompresi + 2 napas\n6) Saat AED tiba, nyalakan dan ikuti instruksinya\n7) Lanjutkan hingga korban pulih atau petugas medis mengambil alih",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"}
      ]
    }
  ]'
);

-- ── Seed: Accidental Ingestion / Tersedak ─────────────────────────────────────
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'Tersedak (Accidental Ingestion)',
  'TERSEDAK',
  'Panduan penanganan tersedak untuk dewasa, anak, dan bayi.',
  '[
    {
      "id": "confirm_age",
      "title": "Konfirmasi Usia Korban",
      "instruction": "Pilih kelompok usia korban untuk mendapatkan panduan yang tepat.",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Dewasa (>10 tahun)", "next_id": "heimlich_adult", "variant": "neutral"},
        {"label": "Anak (< 10 tahun)", "next_id": "heimlich_child", "variant": "neutral"},
        {"label": "Bayi (< 1 tahun)", "next_id": "heimlich_infant", "variant": "neutral"}
      ]
    },
    {
      "id": "heimlich_adult",
      "title": "Pertolongan Tersedak — Dewasa",
      "instruction": "1) Tanya: \"Kamu baik-baik saja?\" lalu \"Kamu tersedak?\"\n2) Perkenalkan diri dan katakan Anda akan membantu\n3) Kepalkan satu tangan dan letakkan di atas pusar, di bawah tulang rusuk\n4) Pegang kepalan tangan dengan tangan lain\n5) Hentakkan ke dalam dan ke atas dengan kuat\n6) Ulangi hingga benda asing keluar dan korban bisa bernapas\n\n*Untuk anak: lakukan hal yang sama namun berlutut terlebih dahulu",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"}
      ]
    },
    {
      "id": "heimlich_child",
      "title": "Pertolongan Tersedak — Anak",
      "instruction": "1) Berlutut di belakang anak\n2) Tanya: \"Kamu tersedak?\"\n3) Kepalkan satu tangan di atas pusar, di bawah tulang rusuk\n4) Pegang kepalan tangan dengan tangan lain\n5) Hentakkan ke dalam dan ke atas\n6) Ulangi hingga benda asing keluar\n\n*Selalu berlutut agar sejajar tinggi badan anak",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"}
      ]
    },
    {
      "id": "heimlich_infant",
      "title": "Pertolongan Tersedak — Bayi",
      "instruction": "JANGAN lakukan Heimlich pada bayi < 1 tahun!\n\n1) Tengkurapkan bayi di lengan Anda, kepala lebih rendah dari dada\n2) Berikan 5 tepukan keras di punggung bayi (antara kedua tulang belikat)\n3) Balikkan bayi, berikan 5 tekanan dada dengan 2 jari di tengah dada\n4) Ulangi hingga benda keluar atau bayi bisa bernapas\n5) Jika bayi tidak sadar → hubungi 119 segera",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"}
      ]
    }
  ]'
);

-- ── Seed: CPR Anak ────────────────────────────────────────────────────────────
INSERT INTO emergency_guide_flows (title, kategori, deskripsi, nodes) VALUES (
  'CPR Anak & Bayi',
  'CPR_ANAK',
  'Panduan CPR khusus untuk anak (1–8 tahun) dan bayi (< 1 tahun).',
  '[
    {
      "id": "select_age",
      "title": "Pilih Kelompok Usia",
      "instruction": "Teknik CPR berbeda untuk anak dan bayi. Pilih usia korban.",
      "image_url": "",
      "is_entry": true,
      "choices": [
        {"label": "Anak (1–8 tahun)", "next_id": "cpr_child", "variant": "neutral"},
        {"label": "Bayi (< 1 tahun)", "next_id": "cpr_infant", "variant": "neutral"}
      ]
    },
    {
      "id": "cpr_child",
      "title": "CPR Anak (1–8 tahun)",
      "instruction": "1) Baringkan anak telentang di permukaan keras\n2) Gunakan 1 atau 2 tangan (sesuai ukuran anak) di tengah dada\n3) Tekan sedalam ±5 cm, 30 kali kompresi\n4) Buka jalan napas: angkat dagu, tengadahkan kepala\n5) Berikan 2 napas bantuan (tiup pelan ±1 detik per napas)\n6) Ulangi siklus 30:2\n7) Jika AED tersedia → nyalakan dan ikuti instruksi\n8) Lanjutkan hingga korban responsif atau petugas medis tiba",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"}
      ]
    },
    {
      "id": "cpr_infant",
      "title": "CPR Bayi (< 1 tahun)",
      "instruction": "1) Baringkan bayi telentang di permukaan keras\n2) Gunakan 2 jari (telunjuk + tengah) di tengah dada, tepat di bawah garis puting\n3) Tekan sedalam ±4 cm, 30 kali kompresi\n4) Buka jalan napas: angkat dagu sedikit (jangan terlalu tengadah)\n5) Tutup mulut DAN hidung bayi, tiup pelan 2 kali\n6) Ulangi siklus 30:2\n7) Lanjutkan hingga bayi responsif atau bantuan tiba",
      "image_url": "",
      "is_entry": false,
      "choices": [
        {"label": "Selesai", "next_id": null, "variant": "neutral"}
      ]
    }
  ]'
);
