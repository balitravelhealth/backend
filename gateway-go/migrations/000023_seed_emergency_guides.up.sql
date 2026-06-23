-- ── PY-15 / PY-16: Seed konten BLS & Emergency Guide ────────────────────────
-- ⚠️  Konten ini adalah DRAFT medis — WAJIB diverifikasi narasumber perawat
--    sebelum digunakan (PY-17). Lihat CLAUDE.md §2 "Keselamatan dulu".
-- Kategori: CPR_DEWASA | CPR_ANAK | CEK_NAPAS | AED | TERSEDAK_DEWASA |
--           TERSEDAK_ANAK | ACCIDENTAL_INGESTION | LUKA | ALERGI | DARURAT
-- ─────────────────────────────────────────────────────────────────────────────

INSERT INTO emergency_guides (kategori, langkah, isi_media) VALUES

-- ── CEK_NAPAS (Cek respons & napas) ──────────────────────────────────────────
('CEK_NAPAS', 1, '{"judul":"Pastikan Lingkungan Aman","teks":"Sebelum mendekati korban, pastikan area sekitar aman untuk Anda dan korban.","ikon":"shield-check"}'),
('CEK_NAPAS', 2, '{"judul":"Cek Respons Korban","teks":"Tepuk bahu korban dengan kuat dan tanya keras: \"Halo, kamu baik-baik saja?\" Amati reaksinya.","ikon":"hand-tap"}'),
('CEK_NAPAS', 3, '{"judul":"Panggil Bantuan","teks":"Jika tidak ada respons, segera minta orang di sekitar untuk menghubungi 119 dan mengambil AED.","ikon":"phone-call","nomor_darurat":"119"}'),
('CEK_NAPAS', 4, '{"judul":"Cek Napas","teks":"Tengadahkan kepala, angkat dagu (Head-tilt Chin-lift). Lihat, dengar, dan rasakan napas selama maksimal 10 detik.","ikon":"lungs"}'),
('CEK_NAPAS', 5, '{"judul":"Mulai CPR jika Tidak Bernafas","teks":"Jika korban tidak bernapas normal (gasping bukan napas normal), segera mulai CPR. Lanjut ke panduan CPR.","ikon":"arrow-right"}'),

-- ── CPR_DEWASA ────────────────────────────────────────────────────────────────
('CPR_DEWASA', 1, '{"judul":"Posisi Tangan","teks":"Letakkan tumit telapak tangan di tengah dada (setengah bawah tulang dada). Tumpuk tangan kedua di atasnya, jari-jari saling kait.","ikon":"hands"}'),
('CPR_DEWASA', 2, '{"judul":"Posisi Tubuh","teks":"Posisikan diri berlutut tegak di samping korban. Pastikan lengan Anda lurus dan bahu sejajar dengan tangan.","ikon":"person-kneeling"}'),
('CPR_DEWASA', 3, '{"judul":"Kompresi Dada","teks":"Tekan dada sedalam 5–6 cm dengan kecepatan 100–120 kali per menit. Biarkan dada mengembang penuh setiap kali.","ikon":"heart-pulse","ritme":"100-120/menit","kedalaman":"5-6 cm"}'),
('CPR_DEWASA', 4, '{"judul":"Ventilasi (Napas Buatan)","teks":"Setelah 30 kompresi, berikan 2 napas buatan (tiap 1 detik). Tutup hidung, tiup mulut hingga dada terangkat. Rasio 30:2.","ikon":"wind","rasio":"30:2"}'),
('CPR_DEWASA', 5, '{"judul":"Lanjutkan Siklus","teks":"Ulangi siklus 30 kompresi + 2 napas. Jangan berhenti kecuali korban mulai bergerak/bernapas normal, AED siap dipasang, atau tenaga medis tiba.","ikon":"repeat"}'),
('CPR_DEWASA', 6, '{"judul":"Pasang AED jika Tersedia","teks":"Segera gunakan AED begitu tersedia. Pasang pad dan ikuti instruksi suara AED. Jangan hentikan CPR saat AED sedang dianalisis.","ikon":"aed"}'),

-- ── CPR_ANAK (1–8 tahun) ──────────────────────────────────────────────────────
('CPR_ANAK', 1, '{"judul":"Cek Respons Anak","teks":"Tepuk telapak kaki anak, panggil namanya. Jika tidak respons, panggil bantuan 119 dan mulai CPR.","ikon":"child","nomor_darurat":"119"}'),
('CPR_ANAK', 2, '{"judul":"Posisi Tangan","teks":"Gunakan SATU tangan (atau dua jari untuk bayi < 1 tahun) di tengah dada, tepat di bawah puting. Untuk anak 1–8 tahun, satu telapak tangan cukup.","ikon":"hand-one"}'),
('CPR_ANAK', 3, '{"judul":"Kompresi Dada","teks":"Tekan sedalam sekitar 5 cm (bayi 4 cm). Kecepatan 100–120 kali per menit. Biarkan dada mengembang penuh.","ikon":"heart-pulse","ritme":"100-120/menit","kedalaman":"5 cm (bayi 4 cm)"}'),
('CPR_ANAK', 4, '{"judul":"Napas Buatan","teks":"Setelah 30 kompresi berikan 2 napas ringan. Untuk bayi, tutup mulut DAN hidung sekaligus. Tiup pelan hingga dada terlihat naik. Rasio 30:2.","ikon":"wind","rasio":"30:2"}'),
('CPR_ANAK', 5, '{"judul":"Lanjutkan & Gunakan AED","teks":"Lanjutkan 30:2 hingga bantuan tiba. Untuk anak ≥ 1 tahun, AED anak (dosis pediatrik) boleh digunakan.","ikon":"repeat"}'),

-- ── AED (Automated External Defibrillator) ────────────────────────────────────
('AED', 1, '{"judul":"Nyalakan AED","teks":"Tekan tombol ON atau buka tutup AED. Ikuti semua instruksi suara dan visual dari perangkat.","ikon":"power"}'),
('AED', 2, '{"judul":"Pasang Pad","teks":"Pasang pad sesuai gambar di AED: satu di kanan atas dada (bawah tulang selangka), satu di sisi kiri bawah dada (di atas tulang rusuk bawah). Pastikan kulit kering.","ikon":"pads"}'),
('AED', 3, '{"judul":"Jauhi Korban saat Analisis","teks":"AED akan menganalisis irama jantung. Pastikan TIDAK ADA yang menyentuh korban saat analisis berlangsung.","ikon":"warning"}'),
('AED', 4, '{"judul":"Tekan Tombol Shock","teks":"Jika AED merekomendasikan shock, pastikan tidak ada yang menyentuh korban, lalu tekan tombol shock. Segera lanjutkan CPR setelah shock.","ikon":"bolt"}'),
('AED', 5, '{"judul":"Lanjutkan CPR","teks":"Setelah shock, langsung lanjutkan CPR 30:2 selama 2 menit sebelum AED menganalisis ulang. Jangan berhenti tanpa instruksi AED.","ikon":"repeat"}'),

-- ── TERSEDAK_DEWASA (Heimlich Maneuver) ──────────────────────────────────────
('TERSEDAK_DEWASA', 1, '{"judul":"Konfirmasi Tersedak","teks":"Tanya: \"Kamu tersedak?\" Jika tidak bisa berbicara, batuk, atau bernapas dan memegang leher — ini tersedak berat. Segera bertindak.","ikon":"question"}'),
('TERSEDAK_DEWASA', 2, '{"judul":"Berdiri di Belakang Korban","teks":"Berdiri sedikit di samping belakang korban. Minta korban membungkuk ke depan. Dukung dada dengan satu tangan.","ikon":"person-behind"}'),
('TERSEDAK_DEWASA', 3, '{"judul":"5 Tepukan Punggung","teks":"Berikan 5 tepukan keras di antara kedua tulang belikat menggunakan tumit telapak tangan. Cek setelah setiap tepukan apakah benda sudah keluar.","ikon":"hand-back","jumlah":5}'),
('TERSEDAK_DEWASA', 4, '{"judul":"5 Dorongan Perut (Heimlich)","teks":"Berdiri di belakang korban, kepalkan satu tangan setinggi pusar (di atas pusar, di bawah tulang dada). Genggam dengan tangan lain. Dorong ke dalam dan ke atas 5 kali dengan kuat.","ikon":"fist","jumlah":5}'),
('TERSEDAK_DEWASA', 5, '{"judul":"Ulangi & Hubungi 119","teks":"Bergantian 5 tepukan punggung + 5 dorongan perut. Jika korban tidak sadar, baringkan dan mulai CPR. Hubungi 119.","ikon":"repeat","nomor_darurat":"119"}'),

-- ── TERSEDAK_ANAK (< 1 tahun: back blows + chest thrusts) ───────────────────
('TERSEDAK_ANAK', 1, '{"judul":"Bayi < 1 Tahun: Posisi Tengkurap","teks":"Pegang bayi tengkurap di lengan bawah Anda, kepala lebih rendah dari dada. Sangga kepala dengan jari.","ikon":"baby"}'),
('TERSEDAK_ANAK', 2, '{"judul":"5 Tepukan Punggung Bayi","teks":"Berikan 5 tepukan keras di punggung bayi (antara kedua tulang belikat) menggunakan tumit telapak tangan.","ikon":"hand-back","jumlah":5}'),
('TERSEDAK_ANAK', 3, '{"judul":"5 Dorongan Dada Bayi","teks":"Balikkan bayi telentang di lengan. Gunakan 2 jari di tengah dada, tepat di bawah puting. Dorong ke bawah 5 kali.","ikon":"two-fingers","jumlah":5}'),
('TERSEDAK_ANAK', 4, '{"judul":"Anak 1–8 Tahun","teks":"Gunakan teknik yang sama seperti dewasa namun dengan kekuatan lebih ringan. Jangan gunakan dorongan perut (Heimlich) pada bayi.","ikon":"child"}'),
('TERSEDAK_ANAK', 5, '{"judul":"Jika Tidak Sadar","teks":"Jika anak tidak sadar, baringkan, cek mulut (keluarkan benda jika terlihat), dan mulai CPR anak. Hubungi 119.","ikon":"alert","nomor_darurat":"119"}'),

-- ── ACCIDENTAL_INGESTION ─────────────────────────────────────────────────────
('ACCIDENTAL_INGESTION', 1, '{"judul":"Jangan Picu Muntah","teks":"JANGAN memicu muntah kecuali diperintahkan oleh tenaga medis atau Poison Control. Muntah bisa memperburuk kondisi pada beberapa zat.","ikon":"warning-red"}'),
('ACCIDENTAL_INGESTION', 2, '{"judul":"Identifikasi Zat","teks":"Catat nama zat, jumlah yang tertelan, dan waktu kejadian. Simpan kemasan/botol untuk ditunjukkan ke dokter.","ikon":"bottle"}'),
('ACCIDENTAL_INGESTION', 3, '{"judul":"Hubungi Layanan Darurat","teks":"Segera hubungi 119 atau IGD terdekat. Berikan informasi zat yang tertelan secara lengkap.","ikon":"phone-call","nomor_darurat":"119"}'),
('ACCIDENTAL_INGESTION', 4, '{"judul":"Pertolongan Awal","teks":"Jika zat terkena kulit/mata, bilas dengan air mengalir minimal 15 menit. Untuk zat yang tertelan, beri air minum (50–100 ml) jika korban sadar dan bisa menelan.","ikon":"water-drop"}'),
('ACCIDENTAL_INGESTION', 5, '{"judul":"Awasi Kondisi","teks":"Pantau kesadaran, pernapasan, dan kondisi umum. Jika kejang atau tidak sadar, baringkan miring (recovery position) dan hubungi 119.","ikon":"eye","nomor_darurat":"119"}'),

-- ── LUKA (Penanganan luka terbuka) ───────────────────────────────────────────
('LUKA', 1, '{"judul":"Cuci Tangan","teks":"Cuci tangan dengan sabun dan air mengalir sebelum menangani luka. Gunakan sarung tangan jika tersedia.","ikon":"hand-wash"}'),
('LUKA', 2, '{"judul":"Hentikan Perdarahan","teks":"Tekan luka dengan kain bersih atau kasa steril selama 10–15 menit tanpa mengangkatnya. Tinggikan anggota badan jika memungkinkan.","ikon":"bandage"}'),
('LUKA', 3, '{"judul":"Bersihkan Luka","teks":"Setelah perdarahan berhenti, bilas luka dengan air bersih mengalir selama 5–10 menit. Keluarkan kotoran yang terlihat dengan hati-hati.","ikon":"water-drop"}'),
('LUKA', 4, '{"judul":"Tutup Luka","teks":"Tutup dengan perban steril atau kain bersih. Ganti perban setiap hari atau jika basah/kotor. Pantau tanda infeksi: merah, bengkak, bernanah, bau.","ikon":"bandage-check"}'),
('LUKA', 5, '{"judul":"Kapan ke Dokter","teks":"Segera ke dokter jika: luka dalam/menganga, perdarahan tidak berhenti, terkena benda berkarat, di wajah/sendi, atau ada tanda infeksi. Pastikan status vaksinasi tetanus.","ikon":"hospital","nomor_darurat":"119"}'),

-- ── ALERGI (Reaksi alergi / Anafilaksis) ────────────────────────────────────
('ALERGI', 1, '{"judul":"Kenali Tanda Anafilaksis","teks":"Gejala berat: sesak napas, bengkak bibir/lidah/tenggorokan, ruam seluruh tubuh, pusing, tekanan darah turun, kolaps. Ini DARURAT MEDIS.","ikon":"alert-red"}'),
('ALERGI', 2, '{"judul":"Gunakan EpiPen jika Ada","teks":"Jika korban memiliki EpiPen (epinefrin auto-injector), bantu gunakan segera di paha luar. Tahan 10 detik. Segera hubungi 119 meski sudah diberi EpiPen.","ikon":"injection","nomor_darurat":"119"}'),
('ALERGI', 3, '{"judul":"Posisi Berbaring","teks":"Baringkan korban dengan kaki ditinggikan (kecuali ada kesulitan bernapas — dudukkan). Jangan biarkan berdiri/berjalan.","ikon":"person-lying"}'),
('ALERGI', 4, '{"judul":"Hubungi 119 Segera","teks":"Anafilaksis adalah darurat medis. Hubungi 119 segera. Beritahu: nama korban, gejala, kemungkinan pemicu, apakah EpiPen sudah diberikan.","ikon":"phone-call","nomor_darurat":"119"}'),
('ALERGI', 5, '{"judul":"Pantau & Siap CPR","teks":"Pantau pernapasan dan kesadaran. Jika korban berhenti bernapas, mulai CPR. Jangan tinggalkan korban sendirian sampai bantuan tiba.","ikon":"heart-monitor"}'),

-- ── DARURAT (Info nomor & layanan darurat Bali) ───────────────────────────────
('DARURAT', 1, '{"judul":"Nomor Darurat Indonesia","teks":"119 — Ambulans & Kegawatdaruratan Medis Nasional (24 jam). Gratis dari semua operator.","ikon":"phone-red","nomor":"119"}'),
('DARURAT', 2, '{"judul":"Polisi & Pemadam","teks":"110 — Polisi. 113 — Pemadam Kebakaran. 112 — Nomor darurat universal (semua operator).","ikon":"badge","nomor_polisi":"110","nomor_damkar":"113","nomor_universal":"112"}'),
('DARURAT', 3, '{"judul":"IGD Terdekat di Bali","teks":"RSUP Prof Ngoerah Denpasar: (0361) 227911. BIMC Hospital Kuta: (0361) 761263. BIMC Nusa Dua: (0361) 3000911. Siloam Bali: (0361) 779900.","ikon":"hospital"}'),
('DARURAT', 4, '{"judul":"BPBD Bali (Bencana)","teks":"BPBD Provinsi Bali: (0361) 255373. Untuk situasi bencana alam, tsunami, atau evakuasi massal.","ikon":"emergency-mgmt"}'),
('DARURAT', 5, '{"judul":"Tips Menghubungi 119","teks":"Sampaikan: (1) Lokasi lengkap/landmark, (2) Kondisi korban, (3) Jumlah korban, (4) Nama & nomor Anda. Tetap di telepon sampai diminta menutup.","ikon":"checklist"}')

ON CONFLICT (kategori, langkah) DO NOTHING;
