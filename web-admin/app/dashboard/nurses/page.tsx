"use client";

import { useEffect, useState, useCallback } from "react";
import { nursesApi } from "@/lib/api";
import type { Nurse } from "@/lib/types";
import DataTable from "@/components/DataTable";
import Modal from "@/components/Modal";
import PageHeader, { AddButton } from "@/components/PageHeader";
import FormField, { Input } from "@/components/FormField";
import { useToast, Toaster } from "@/components/Toast";

const emptyForm = {
  email: "", password: "", nama_lengkap: "",
  nomor_lisensi: "", sertifikasi: "",
};

export default function NursesPage() {
  const [rows, setRows]       = useState<Nurse[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen]       = useState(false);
  const [form, setForm]       = useState(emptyForm);
  const [saving, setSaving]   = useState(false);
  const [error, setError]     = useState("");
  const { toasts, remove, toast } = useToast();

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await nursesApi.list();
      setRows(res.data ?? []);
    } catch {
      toast.error("Gagal memuat data perawat");
    } finally { setLoading(false); }
  }, [toast]);

  useEffect(() => { load(); }, [load]);

  async function handleToggle(row: Nurse) {
    try {
      await nursesApi.toggle(row.id);
      toast.success(row.aktif ? `${row.nama_lengkap} dinonaktifkan` : `${row.nama_lengkap} diaktifkan`);
      load();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal mengubah status");
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setSaving(true);
    try {
      await nursesApi.create({
        email: form.email,
        password: form.password,
        nama_lengkap: form.nama_lengkap,
        nomor_lisensi: form.nomor_lisensi,
        sertifikasi: form.sertifikasi || undefined,
      });
      toast.success(`Akun perawat ${form.nama_lengkap} dibuat`);
      setOpen(false);
      load();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Gagal menyimpan");
    } finally { setSaving(false); }
  }

  const columns = [
    { key: "id", label: "ID", render: (r: Nurse) => <span className="font-mono text-xs">{r.id}</span> },
    { key: "nama_lengkap", label: "Nama", render: (r: Nurse) => <span className="font-medium">{r.nama_lengkap}</span> },
    { key: "nomor_lisensi", label: "No. Lisensi STR" },
    { key: "sertifikasi", label: "Sertifikasi", render: (r: Nurse) => r.sertifikasi ?? "—" },
    {
      key: "aktif", label: "Status",
      render: (r: Nurse) => (
        <span
          className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold"
          style={{
            color: r.aktif ? "var(--success)" : "var(--text-muted)",
            background: r.aktif ? "#f0fdf4" : "var(--bg)",
          }}
        >
          <span className="w-1.5 h-1.5 rounded-full" style={{ background: r.aktif ? "var(--success)" : "#94a3b8" }} />
          {r.aktif ? "Aktif" : "Nonaktif"}
        </span>
      ),
    },
    {
      key: "actions", label: "",
      render: (r: Nurse) => (
        <button
          onClick={() => handleToggle(r)}
          className="px-2.5 py-1 rounded-md text-xs font-medium transition-colors"
          style={{
            color:  r.aktif ? "var(--warning)" : "var(--success)",
            border: `1px solid ${r.aktif ? "#fcd34d" : "#86efac"}`,
          }}
        >
          {r.aktif ? "Nonaktifkan" : "Aktifkan"}
        </button>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Kelola Perawat"
        description="Tambah akun perawat dan atur status aktif/nonaktif untuk akses sistem"
        action={<AddButton onClick={() => { setForm(emptyForm); setError(""); setOpen(true); }} label="Tambah Perawat" />}
      />

      <div className="reveal">
        <DataTable columns={columns} rows={rows} keyField="id" loading={loading} />
      </div>

      <Modal open={open} title="Tambah Akun Perawat" onClose={() => setOpen(false)}>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Panduan */}
          <div className="px-3 py-2.5 rounded-lg text-xs" style={{ background: "#eff6ff", color: "#1d4ed8", border: "1px solid #bfdbfe" }}>
            Akun perawat digunakan untuk mengakses fitur nursing care di aplikasi. Nomor lisensi STR dan sertifikasi harus sesuai dokumen resmi.
          </div>

          <FormField label="Nama Lengkap" required>
            <Input value={form.nama_lengkap}
              onChange={(e) => setForm({ ...form, nama_lengkap: e.target.value })}
              required placeholder="Nama lengkap sesuai KTP" />
          </FormField>
          <div className="grid grid-cols-2 gap-4">
            <FormField label="Email" required hint="Digunakan untuk login">
              <Input type="email" value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
                required placeholder="perawat@balihealth.me" />
            </FormField>
            <FormField label="Password" required hint="Minimal 8 karakter — sarankan gunakan kombinasi huruf & angka">
              <Input type="password" value={form.password}
                onChange={(e) => setForm({ ...form, password: e.target.value })}
                required minLength={8} placeholder="Min. 8 karakter" />
            </FormField>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <FormField label="Nomor Lisensi (STR)" required hint="Surat Tanda Registrasi perawat">
              <Input value={form.nomor_lisensi}
                onChange={(e) => setForm({ ...form, nomor_lisensi: e.target.value })}
                required placeholder="cth: STR-12345678" />
            </FormField>
            <FormField label="Sertifikasi" hint="Opsional — pisahkan dengan koma jika lebih dari satu">
              <Input value={form.sertifikasi}
                onChange={(e) => setForm({ ...form, sertifikasi: e.target.value })}
                placeholder="cth: BTLS, ACLS, PPGD" />
            </FormField>
          </div>

          {error && (
            <p className="text-sm px-3 py-2 rounded-lg"
               style={{ color: "var(--danger)", background: "#fef2f2", border: "1px solid #fecaca" }}>
              {error}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={() => setOpen(false)}
              className="px-4 py-2 rounded-lg text-sm font-medium"
              style={{ color: "var(--text-muted)", border: "1px solid var(--border)" }}>
              Batal
            </button>
            <button type="submit" disabled={saving}
              className="px-4 py-2 rounded-lg text-sm font-semibold text-white disabled:opacity-60"
              style={{ background: "var(--brand-grad)" }}>
              {saving ? "Membuat akun…" : "Buat Akun Perawat"}
            </button>
          </div>
        </form>
      </Modal>

      <Toaster toasts={toasts} remove={remove} />
    </div>
  );
}
