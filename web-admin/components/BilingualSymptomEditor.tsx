"use client";

import React, { useState } from "react";
import type { Symptom } from "@/lib/types";
import FormField, { Input, Textarea } from "./FormField";
import LanguageTabs, { type Language } from "./LanguageTabs";

interface BilingualSymptomEditorProps {
  symptom: Symptom | null;
  onSave: (symptom: Omit<Symptom, "symptom_id" | "created_at" | "updated_at">) => Promise<void>;
  loading: boolean;
  error: string;
}

export default function BilingualSymptomEditor({
  symptom,
  onSave,
  loading,
  error,
}: BilingualSymptomEditorProps) {
  const [lang, setLang] = React.useState<Language>("id");
  const [form, setForm] = useState({
    kode: symptom?.kode ?? "",
    label_id: symptom?.label_id ?? "",
    label_en: symptom?.label_en ?? "",
    deskripsi_id: symptom?.deskripsi_id ?? "",
    deskripsi_en: symptom?.deskripsi_en ?? "",
  });

  const handleSave = async () => {
    if (!form.kode.trim() || !form.label_id.trim() || !form.label_en.trim()) {
      alert("Kode dan label (ID & EN) wajib diisi");
      return;
    }
    await onSave({
      kode: form.kode,
      label_id: form.label_id,
      label_en: form.label_en,
      deskripsi_id: form.deskripsi_id || undefined,
      deskripsi_en: form.deskripsi_en || undefined,
    } as Omit<Symptom, "symptom_id" | "created_at" | "updated_at">);
  };

  return (
    <div>
      <FormField label="Kode Gejala">
        <Input
          value={form.kode}
          onChange={(e) => setForm({ ...form, kode: e.target.value })}
          placeholder="e.g., S_DIARE"
          disabled={!!symptom}
        />
      </FormField>

      <LanguageTabs value={lang} onChange={setLang} />

      {lang === "id" ? (
        <>
          <FormField label="Label (Bahasa Indonesia)" required>
            <Input
              value={form.label_id}
              onChange={(e) => setForm({ ...form, label_id: e.target.value })}
              placeholder="Diare (BAB cair > 3x/hari)"
            />
          </FormField>
          <FormField label="Deskripsi (Bahasa Indonesia)">
            <Textarea
              value={form.deskripsi_id}
              onChange={(e) => setForm({ ...form, deskripsi_id: e.target.value })}
              placeholder="Deskripsi gejala dalam bahasa Indonesia..."
            />
          </FormField>
        </>
      ) : (
        <>
          <FormField label="Label (English)" required>
            <Input
              value={form.label_en}
              onChange={(e) => setForm({ ...form, label_en: e.target.value })}
              placeholder="Diarrhea (>3 loose stools/day)"
            />
          </FormField>
          <FormField label="Description (English)">
            <Textarea
              value={form.deskripsi_en}
              onChange={(e) => setForm({ ...form, deskripsi_en: e.target.value })}
              placeholder="Description of symptom in English..."
            />
          </FormField>
        </>
      )}

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
