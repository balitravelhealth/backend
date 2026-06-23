"use client";

import React, { useState } from "react";
import type { Disease } from "@/lib/types";
import type { RekDefault } from "@/lib/api";
import FormField, { Input, Textarea } from "./FormField";
import LanguageTabs, { type Language } from "./LanguageTabs";

interface BilingualDiseaseEditorProps {
  disease: Disease | null;
  onSave: (disease: Omit<Disease, "id" | "created_at" | "updated_at">) => Promise<void>;
  loading: boolean;
  error: string;
}

const emptyRek: RekDefault = { Rendah: "", Sedang: "", Tinggi: "", Darurat: "" };
const RISK_LEVELS: (keyof RekDefault)[] = ["Rendah", "Sedang", "Tinggi", "Darurat"];

export default function BilingualDiseaseEditor({
  disease,
  onSave,
  loading,
  error,
}: BilingualDiseaseEditorProps) {
  const [lang, setLang] = React.useState<Language>("id");
  const [form, setForm] = useState({
    nama_id: disease?.nama_id ?? "",
    nama_en: disease?.nama_en ?? "",
    deskripsi_id: disease?.deskripsi_id ?? "",
    deskripsi_en: disease?.deskripsi_en ?? "",
    rekomendasi_default_id: disease?.rekomendasi_default_id ?? emptyRek,
    rekomendasi_default_en: disease?.rekomendasi_default_en ?? emptyRek,
  });

  const handleSave = async () => {
    if (!form.nama_id.trim() || !form.nama_en.trim()) {
      alert("Nama penyakit (ID & EN) wajib diisi");
      return;
    }
    await onSave({
      nama_id: form.nama_id,
      nama_en: form.nama_en,
      deskripsi_id: form.deskripsi_id || undefined,
      deskripsi_en: form.deskripsi_en || undefined,
      rekomendasi_default_id: form.rekomendasi_default_id,
      rekomendasi_default_en: form.rekomendasi_default_en,
    } as Omit<Disease, "id" | "created_at" | "updated_at">);
  };

  const handleRekChange = (
    lang: Language,
    level: keyof RekDefault,
    value: string,
  ) => {
    if (lang === "id") {
      setForm({
        ...form,
        rekomendasi_default_id: { ...form.rekomendasi_default_id, [level]: value },
      });
    } else {
      setForm({
        ...form,
        rekomendasi_default_en: { ...form.rekomendasi_default_en, [level]: value },
      });
    }
  };

  const rek = lang === "id" ? form.rekomendasi_default_id : form.rekomendasi_default_en;

  return (
    <div>
      <LanguageTabs value={lang} onChange={setLang} />

      {lang === "id" ? (
        <>
          <FormField label="Nama Penyakit (Bahasa Indonesia)" required>
            <Input
              value={form.nama_id}
              onChange={(e) => setForm({ ...form, nama_id: e.target.value })}
              placeholder="e.g., Bali Belly (Diare Wisatawan)"
            />
          </FormField>
          <FormField label="Deskripsi (Bahasa Indonesia)">
            <Textarea
              value={form.deskripsi_id}
              onChange={(e) => setForm({ ...form, deskripsi_id: e.target.value })}
              placeholder="Deskripsi penyakit..."
            />
          </FormField>
        </>
      ) : (
        <>
          <FormField label="Disease Name (English)" required>
            <Input
              value={form.nama_en}
              onChange={(e) => setForm({ ...form, nama_en: e.target.value })}
              placeholder="e.g., Bali Belly (Traveler's Diarrhea)"
            />
          </FormField>
          <FormField label="Description (English)">
            <Textarea
              value={form.deskripsi_en}
              onChange={(e) => setForm({ ...form, deskripsi_en: e.target.value })}
              placeholder="Disease description..."
            />
          </FormField>
        </>
      )}

      <div className="mt-6 pt-4 border-t">
        <h3 style={{ fontSize: "0.875rem", fontWeight: "600", marginBottom: "1rem" }}>
          Rekomendasi {lang === "id" ? "(Bahasa Indonesia)" : "(English)"}
        </h3>
        {RISK_LEVELS.map((level) => (
          <FormField key={level} label={`Rekomendasi - ${level}`}>
            <Textarea
              value={rek[level] ?? ""}
              onChange={(e) => handleRekChange(lang, level, e.target.value)}
              placeholder={`Masukkan rekomendasi untuk risiko ${level}...`}
              rows={2}
            />
          </FormField>
        ))}
      </div>

      {error && <div style={{ color: "var(--danger)", fontSize: "0.875rem", marginTop: "0.5rem" }}>{error}</div>}

      <div className="flex gap-2 mt-6">
        <button
          onClick={handleSave}
          disabled={loading}
          style={{
            padding: "0.5rem 1rem",
            background: "var(--primary)",
            color: "white",
            border: "none",
            borderRadius: "0.375rem",
            cursor: loading ? "not-allowed" : "pointer",
            opacity: loading ? 0.6 : 1,
          }}
        >
          {loading ? "Menyimpan..." : "Simpan"}
        </button>
      </div>
    </div>
  );
}
