"use client";

import { useEffect, useState, useCallback } from "react";
import { facilitiesApi, destinationsApi } from "@/lib/api";
import type { Facility, Destination } from "@/lib/types";
import DataTable from "@/components/DataTable";
import Modal from "@/components/Modal";
import PageHeader, { AddButton, ActionBtn } from "@/components/PageHeader";
import FormField, { Input, Select, Textarea } from "@/components/FormField";
import { useToast, Toaster } from "@/components/Toast";

const KATEGORI = ["Rumah Sakit", "Klinik", "Puskesmas", "Apotek", "Dokter Umum", "Lainnya"];

const emptyForm = {
  destination_id: "",
  nama: "",
  kategori: "",
  alamat: "",
  latitude: "",
  longitude: "",
  kontak: "",
  jam_operasional: "",
};

export default function FacilitiesPage() {
  const [rows, setRows]         = useState<Facility[]>([]);
  const [destinations, setDest] = useState<Destination[]>([]);
  const [loading, setLoading]   = useState(true);
  const [open, setOpen]         = useState(false);
  const [editing, setEditing]   = useState<Facility | null>(null);
  const [form, setForm]         = useState(emptyForm);
  const [saving, setSaving]     = useState(false);
  const [error, setError]       = useState("");
  const { toasts, remove, toast } = useToast();

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [fRes, dRes] = await Promise.all([
        facilitiesApi.list(),
        destinationsApi.list(),
      ]);
      setRows(fRes.data ?? []);
      setDest(dRes.data ?? []);
    } catch {
      toast.error("Gagal memuat data");
    } finally {
      setLoading(false);
    }
  }, [toast]);

  useEffect(() => { load(); }, [load]);

  function openAdd() {
    setEditing(null);
    setForm(emptyForm);
    setError("");
    setOpen(true);
  }

  function openEdit(row: Facility) {
    setEditing(row);
    setForm({
      destination_id: String(row.destination_id),
      nama: row.nama,
      kategori: row.kategori ?? "",
      alamat: row.alamat ?? "",
      latitude: String(row.latitude ?? ""),
      longitude: String(row.longitude ?? ""),
      kontak: row.kontak ?? "",
      jam_operasional: row.jam_operasional ?? "",
    });
    setError("");
    setOpen(true);
  }

  async function handleDelete(row: Facility) {
    if (!confirm(`Hapus fasilitas "${row.nama}"?`)) return;
    try {
      await facilitiesApi.del(row.id);
      toast.success(`"${row.nama}" dihapus`);
      load();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal menghapus");
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    if (!form.destination_id) {
      setError("Pilih destinasi terlebih dahulu");
      return;
    }
    setSaving(true);
    try {
      const payload = {
        destination_id: Number(form.destination_id),
        nama: form.nama,
        kategori: form.kategori || undefined,
        alamat: form.alamat || undefined,
        latitude: form.latitude ? Number(form.latitude) : undefined,
        longitude: form.longitude ? Number(form.longitude) : undefined,
        kontak: form.kontak || undefined,
        jam_operasional: form.jam_operasional || undefined,
      };
      if (editing) {
        await facilitiesApi.update(editing.id, payload);
        toast.success("Fasilitas diperbarui");
      } else {
        await facilitiesApi.create(payload as Parameters<typeof facilitiesApi.create>[0]);
        toast.success("Fasilitas ditambahkan");
      }
      setOpen(false);
      load();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Gagal menyimpan");
    } finally {
      setSaving(false);
    }
  }

  const columns = [
    { key: "id", label: "ID", render: (r: Facility) => <span className="font-mono text-xs">{r.id}</span> },
    { key: "nama", label: "Nama", render: (r: Facility) => <span className="font-medium">{r.nama}</span> },
    { key: "kategori", label: "Kategori", render: (r: Facility) => r.kategori ?? "—" },
    {
      key: "destination_id", label: "Destinasi",
      render: (r: Facility) =>
        destinations.find((d) => d.id === r.destination_id)?.nama_daerah ?? `#${r.destination_id}`,
    },
    { key: "kontak", label: "Kontak", render: (r: Facility) => r.kontak ?? "—" },
    {
      key: "jam_operasional", label: "Jam Buka",
      render: (r: Facility) => <span className="text-xs">{r.jam_operasional ?? "—"}</span>,
    },
    {
      key: "actions", label: "",
      render: (r: Facility) => (
        <div className="flex items-center gap-1">
          <ActionBtn variant="edit"   onClick={() => openEdit(r)} />
          <ActionBtn variant="delete" onClick={() => handleDelete(r)} />
        </div>
      ),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Fasilitas Kesehatan"
        description="Kelola data fasilitas kesehatan di destinasi wisata Bali"
        action={<AddButton onClick={openAdd} label="Tambah Fasilitas" />}
      />

      <div className="reveal">
        <DataTable columns={columns} rows={rows} keyField="id" loading={loading} />
      </div>

      <Modal
        open={open}
        title={editing ? `Edit — ${editing.nama}` : "Tambah Fasilitas Kesehatan"}
        onClose={() => setOpen(false)}
        size="lg"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Panduan singkat */}
          <div className="px-3 py-2.5 rounded-lg text-xs" style={{ background: "#eff6ff", color: "#1d4ed8", border: "1px solid #bfdbfe" }}>
            Isi nama dan pilih destinasi (wajib). Koordinat latitude/longitude digunakan untuk tampilkan peta di aplikasi wisatawan.
          </div>

          <div className="grid grid-cols-2 gap-4">
            <FormField label="Nama Fasilitas" required>
              <Input
                value={form.nama}
                onChange={(e) => setForm({ ...form, nama: e.target.value })}
                required
                placeholder="cth: RSUP Sanglah"
              />
            </FormField>
            <FormField label="Destinasi" required hint="Kecamatan/kabupaten tempat fasilitas berada">
              <Select
                value={form.destination_id}
                onChange={(e) => setForm({ ...form, destination_id: e.target.value })}
                required
              >
                <option value="">— pilih destinasi —</option>
                {destinations.map((d) => (
                  <option key={d.id} value={d.id}>{d.nama_daerah}</option>
                ))}
              </Select>
            </FormField>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <FormField label="Kategori">
              <Select
                value={form.kategori}
                onChange={(e) => setForm({ ...form, kategori: e.target.value })}
              >
                <option value="">— pilih —</option>
                {KATEGORI.map((k) => <option key={k} value={k}>{k}</option>)}
              </Select>
            </FormField>
            <FormField label="Nomor Kontak" hint="Format: +62...">
              <Input
                value={form.kontak}
                onChange={(e) => setForm({ ...form, kontak: e.target.value })}
                placeholder="+62361-123456"
              />
            </FormField>
          </div>

          <FormField label="Alamat Lengkap">
            <Textarea
              value={form.alamat}
              onChange={(e) => setForm({ ...form, alamat: e.target.value })}
              placeholder="Jl. ..., Kelurahan ..., Kecamatan ..."
            />
          </FormField>

          <div className="grid grid-cols-2 gap-4">
            <FormField label="Latitude" hint="cth: -8.6705 (negatif untuk Bali)">
              <Input
                type="number" step="any"
                value={form.latitude}
                onChange={(e) => setForm({ ...form, latitude: e.target.value })}
                placeholder="-8.6705"
              />
            </FormField>
            <FormField label="Longitude" hint="cth: 115.2126">
              <Input
                type="number" step="any"
                value={form.longitude}
                onChange={(e) => setForm({ ...form, longitude: e.target.value })}
                placeholder="115.2126"
              />
            </FormField>
          </div>

          <FormField label="Jam Operasional">
            <Input
              value={form.jam_operasional}
              onChange={(e) => setForm({ ...form, jam_operasional: e.target.value })}
              placeholder="cth: Senin–Jumat 08.00–17.00 / 24 jam"
            />
          </FormField>

          {error && (
            <p className="text-sm px-3 py-2 rounded-lg"
               style={{ color: "var(--danger)", background: "#fef2f2", border: "1px solid #fecaca" }}>
              {error}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <button
              type="button" onClick={() => setOpen(false)}
              className="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
              style={{ color: "var(--text-muted)", border: "1px solid var(--border)" }}
            >
              Batal
            </button>
            <button
              type="submit" disabled={saving}
              className="px-4 py-2 rounded-lg text-sm font-semibold text-white transition-opacity disabled:opacity-60"
              style={{ background: "var(--brand-grad)" }}
            >
              {saving ? "Menyimpan…" : editing ? "Simpan Perubahan" : "Tambah Fasilitas"}
            </button>
          </div>
        </form>
      </Modal>

      <Toaster toasts={toasts} remove={remove} />
    </div>
  );
}
