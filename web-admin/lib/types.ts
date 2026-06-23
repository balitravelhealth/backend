export interface Facility {
  id: number;
  destination_id: number;
  nama: string;
  kategori?: string;
  alamat?: string;
  latitude?: number;
  longitude?: number;
  kontak?: string;
  jam_operasional?: string;
  created_at: string;
  updated_at: string;
}

export interface Destination {
  id: number;
  nama_daerah: string;
  created_at: string;
}

export interface HealthRisk {
  id: number;
  destination_id: number;
  nama_risiko_id: string;
  nama_risiko_en: string;
  saran_pencegahan_id?: string;
  saran_pencegahan_en?: string;
  rekomendasi_vaksinasi_id?: string;
  rekomendasi_vaksinasi_en?: string;
  created_at: string;
  updated_at: string;
}

export interface EmergencyGuide {
  id: number;
  kategori: string;
  langkah: number;
  isi_media_id: Record<string, unknown>;
  isi_media_en: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface Nurse {
  id: number;
  user_id: number;
  nama_lengkap: string;
  nomor_lisensi: string;
  sertifikasi?: string;
  aktif: boolean;
  created_at: string;
  updated_at: string;
}

export interface Assessment {
  id: number;
  user_id: number;
  symptoms: unknown;
  diagnosis?: string;
  confidence_score?: number;
  risk_level?: string;
  created_at: string;
}

export interface AssessmentPage {
  data: Assessment[];
  total: number;
  page: number;
  limit: number;
}

export interface Symptom {
  symptom_id: number;
  kode: string;
  label_id: string;
  label_en: string;
  deskripsi_id?: string;
  deskripsi_en?: string;
  created_at: string;
  updated_at: string;
}

export interface Disease {
  id: number;
  nama_id: string;
  nama_en: string;
  deskripsi_id?: string;
  deskripsi_en?: string;
  rekomendasi_default_id?: Record<string, string>;
  rekomendasi_default_en?: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export interface GuideChoice {
  label: string;
  next_id: string | null;
  variant: "yes" | "no" | "neutral";
}

export interface GuideNode {
  id: string;
  title: string;
  instruction: string;
  image_url: string;
  is_entry: boolean;
  choices: GuideChoice[];
}

export interface EmergencyGuideFlow {
  id: number;
  title_id: string;
  title_en: string;
  kategori: string;
  deskripsi?: string;
  nodes_id: GuideNode[];
  nodes_en: GuideNode[];
  created_at: string;
  updated_at: string;
}

export interface Rule {
  rule_id: number;
  nama: string;
  premis: number[];
  disease_id: number;
  bobot_cf: number;
  mb: number;
  md: number;
  kategori: string;
  status: string;
  created_by: number;
  created_at: string;
  updated_at: string;
}
