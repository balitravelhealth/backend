const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("admin_token") : null;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });

  if (res.status === 401) {
    if (typeof window !== "undefined") {
      localStorage.removeItem("admin_token");
      if (!window.location.pathname.startsWith("/login")) {
        window.location.href = "/login";
      }
    }
    let msg = "Sesi habis, silakan login ulang";
    try {
      const body = await res.clone().json();
      if (body.error) msg = body.error;
    } catch {}
    throw new ApiError(401, msg);
  }

  if (!res.ok) {
    let msg = res.statusText;
    try {
      const body = await res.json();
      msg = body.error ?? body.message ?? msg;
    } catch {}
    throw new ApiError(res.status, msg);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "POST", body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body) }),
  del: (path: string) => request<void>(path, { method: "DELETE" }),
};

export async function uploadImage(file: File): Promise<string> {
  const token = typeof window !== "undefined" ? localStorage.getItem("admin_token") : null;
  const form = new FormData();
  form.append("file", file);
  const res = await fetch(`${BASE_URL}/admin/upload`, {
    method: "POST",
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: form,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new ApiError(res.status, body.error ?? "Upload gagal");
  }
  const data = await res.json();
  return data.url as string;
}

// ── Auth ──────────────────────────────────────────────────────────────────────

export interface AdminLoginResponse {
  access_token: string;
  admin: { id: number; email: string; nama: string };
}

export function adminLogin(email: string, password: string) {
  return api.post<AdminLoginResponse>("/admin/auth/login", { email, password });
}

export function saveToken(token: string) {
  localStorage.setItem("admin_token", token);
}

export function clearToken() {
  localStorage.removeItem("admin_token");
}

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("admin_token");
}

// ── Facilities (GO-24) ────────────────────────────────────────────────────────

import type {
  Facility, Destination, HealthRisk, EmergencyGuide, EmergencyGuideFlow, GuideNode,
  Nurse, AssessmentPage, Symptom, Disease, Rule,
} from "./types";

export const facilitiesApi = {
  list: () => api.get<{ data: Facility[] }>("/admin/facilities"),
  create: (body: Omit<Facility, "id" | "created_at" | "updated_at">) =>
    api.post<Facility>("/admin/facilities", body),
  update: (id: number, body: Partial<Facility>) =>
    api.put<Facility>(`/admin/facilities/${id}`, body),
  del: (id: number) => api.del(`/admin/facilities/${id}`),
};

// ── Destinations & health risks (GO-25) ───────────────────────────────────────

export const destinationsApi = {
  list: () => api.get<{ data: Destination[] }>("/destinations"),
  create: (nama_daerah: string) =>
    api.post<Destination>("/admin/destinations", { nama_daerah }),
  update: (id: number, nama_daerah: string) =>
    api.put<Destination>(`/admin/destinations/${id}`, { nama_daerah }),
  del: (id: number) => api.del(`/admin/destinations/${id}`),
};

export const healthRisksApi = {
  list: (destinationId: number) =>
    api.get<{ data: HealthRisk[] }>(`/destinations/${destinationId}/health-risks`),
  create: (body: Omit<HealthRisk, "id" | "created_at" | "updated_at">) =>
    api.post<HealthRisk>("/admin/health-risks", body),
  update: (id: number, body: Partial<HealthRisk>) =>
    api.put<HealthRisk>(`/admin/health-risks/${id}`, body),
  del: (id: number) => api.del(`/admin/health-risks/${id}`),
};

export const emergencyGuidesApi = {
  list: () => api.get<{ data: EmergencyGuide[] }>("/emergency-guides"),
  create: (body: Omit<EmergencyGuide, "id" | "created_at" | "updated_at">) =>
    api.post<EmergencyGuide>("/admin/emergency-guides", body),
  update: (id: number, body: Partial<EmergencyGuide>) =>
    api.put<EmergencyGuide>(`/admin/emergency-guides/${id}`, body),
  del: (id: number) => api.del(`/admin/emergency-guides/${id}`),
};

export const emergencyFlowsApi = {
  list: () => api.get<{ data: EmergencyGuideFlow[] }>("/admin/emergency-guide-flows"),
  get: (id: number) => api.get<EmergencyGuideFlow>(`/emergency-guide-flows/${id}`),
  create: (body: { title_id: string; title_en: string; kategori: string; deskripsi?: string; nodes_id: GuideNode[]; nodes_en: GuideNode[] }) =>
    api.post<EmergencyGuideFlow>("/admin/emergency-guide-flows", body),
  update: (id: number, body: { title_id: string; title_en: string; kategori: string; deskripsi?: string; nodes_id: GuideNode[]; nodes_en: GuideNode[] }) =>
    api.put<EmergencyGuideFlow>(`/admin/emergency-guide-flows/${id}`, body),
  del: (id: number) => api.del(`/admin/emergency-guide-flows/${id}`),
};

// ── Nurses (GO-26) ────────────────────────────────────────────────────────────

export const nursesApi = {
  list: () => api.get<{ data: Nurse[] }>("/admin/nurses"),
  create: (body: {
    email: string; password: string; nama_lengkap: string;
    nomor_lisensi: string; sertifikasi?: string;
  }) => api.post<Nurse>("/admin/nurses", body),
  toggle: (id: number) => api.put<Nurse>(`/admin/nurses/${id}/toggle`, {}),
};

// ── Assessments (GO-27) ───────────────────────────────────────────────────────

export const assessmentsApi = {
  list: (page = 1, limit = 20) =>
    api.get<AssessmentPage>(`/admin/assessments?page=${page}&limit=${limit}`),
};

// ── Expert knowledge base (GO-28) ─────────────────────────────────────────────

export const symptomsApi = {
  list: () => api.get<{ data: Symptom[] }>("/admin/expert/symptoms"),
  create: (body: { kode: string; label_id: string; label_en: string; deskripsi_id?: string; deskripsi_en?: string }) =>
    api.post<Symptom>("/admin/expert/symptoms", body),
  update: (id: number, body: { kode: string; label_id: string; label_en: string; deskripsi_id?: string; deskripsi_en?: string }) =>
    api.put<Symptom>(`/admin/expert/symptoms/${id}`, body),
  del: (id: number) => api.del(`/admin/expert/symptoms/${id}`),
};

export interface RekDefault {
  Rendah:  string;
  Sedang:  string;
  Tinggi:  string;
  Darurat: string;
}

export const diseasesApi = {
  list: () => api.get<{ data: Disease[] }>("/admin/expert/diseases"),
  create: (body: { nama_id: string; nama_en: string; deskripsi_id?: string; deskripsi_en?: string; rekomendasi_default_id?: RekDefault; rekomendasi_default_en?: RekDefault }) =>
    api.post<Disease>("/admin/expert/diseases", body),
  update: (id: number, body: { nama_id: string; nama_en: string; deskripsi_id?: string; deskripsi_en?: string; rekomendasi_default_id?: RekDefault; rekomendasi_default_en?: RekDefault }) =>
    api.put<Disease>(`/admin/expert/diseases/${id}`, body),
  del: (id: number) => api.del(`/admin/expert/diseases/${id}`),
};

export const rulesApi = {
  list: () => api.get<{ data: Rule[] }>("/admin/expert/rules"),
  create: (body: {
    nama: string; premis: number[]; disease_id: number;
    bobot_cf: number; mb: number; md: number; kategori: string;
  }) => api.post<Rule>("/admin/expert/rules", body),
  update: (id: number, body: {
    nama: string; premis: number[]; disease_id: number;
    bobot_cf: number; mb: number; md: number; kategori: string;
  }) => api.put<Rule>(`/admin/expert/rules/${id}`, body),
  publish: (id: number) => api.post<Rule>(`/admin/expert/rules/${id}/publish`, {}),
  unpublish: (id: number) => api.post<Rule>(`/admin/expert/rules/${id}/unpublish`, {}),
  del: (id: number) => api.del(`/admin/expert/rules/${id}`),
};
