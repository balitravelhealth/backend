"use client";

import { useEffect, useState, useCallback } from "react";
import { destinationsApi, healthRisksApi } from "@/lib/api";
import type { Destination, HealthRisk } from "@/lib/types";
import DataTable from "@/components/DataTable";
import Modal from "@/components/Modal";
import PageHeader, { AddButton, ActionBtn } from "@/components/PageHeader";
import FormField, { Input, Textarea } from "@/components/FormField";
import { useToast, Toaster } from "@/components/Toast";

export default function DestinationsPage() {
  const [destinations, setDest]   = useState<Destination[]>([]);
  const [loading, setLoading]     = useState(true);
  const [destModal, setDestModal] = useState(false);
  const [editingDest, setEditingDest] = useState<Destination | null>(null);
  const [destName, setDestName]   = useState("");
  const [destSaving, setDestSaving] = useState(false);
  const [destError, setDestError] = useState("");

  // health-risk panel
  const [selected, setSelected]   = useState<Destination | null>(null);
  const [risks, setRisks]         = useState<HealthRisk[]>([]);
  const [riskLoading, setRiskLoading] = useState(false);
  const [riskModal, setRiskModal] = useState(false);
  const [editingRisk, setEditingRisk] = useState<HealthRisk | null>(null);
  const [riskForm, setRiskForm]   = useState({ nama_risiko_id: "", nama_risiko_en: "", saran_pencegahan_id: "", saran_pencegahan_en: "", rekomendasi_vaksinasi_id: "", rekomendasi_vaksinasi_en: "" });
  const [riskSaving, setRiskSaving] = useState(false);
  const [riskError, setRiskError] = useState("");

  const { toasts, remove, toast } = useToast();

  const loadDest = useCallback(async () => {
    setLoading(true);
    try {
      const res = await destinationsApi.list();
      setDest(res.data ?? []);
    } catch {
      toast.error("Gagal memuat destinasi");
    } finally { setLoading(false); }
  }, [toast]);

  useEffect(() => { loadDest(); }, [loadDest]);

  async function loadRisks(dest: Destination) {
    setSelected(dest);
    setRiskLoading(true);
    try {
      const res = await healthRisksApi.list(dest.id);
      setRisks(res.data ?? []);
    } catch {
      toast.error("Gagal memuat risiko kesehatan");
    } finally { setRiskLoading(false); }
  }

  function openAddDest() {
    setEditingDest(null);
    setDestName("");
    setDestError("");
    setDestModal(true);
  }

  function openEditDest(d: Destination) {
    setEditingDest(d);
    setDestName(d.nama_daerah);
    setDestError("");
    setDestModal(true);
  }

  async function handleSaveDest(e: React.FormEvent) {
    e.preventDefault();
    setDestError("");
    setDestSaving(true);
    try {
      if (editingDest) {
        await destinationsApi.update(editingDest.id, destName);
        // refresh selected panel if currently viewing this destination
        if (selected?.id === editingDest.id) {
          setSelected((prev) => prev ? { ...prev, nama_daerah: destName } : prev);
        }
        toast.success("Nama destinasi diperbarui");
      } else {
        await destinationsApi.create(destName);
        toast.success("Destinasi ditambahkan");
      }
      setDestModal(false);
      setDestName("");
      loadDest();
    } catch (err: unknown) {
      setDestError(err instanceof Error ? err.message : "Gagal menyimpan");
    } finally { setDestSaving(false); }
  }

  async function handleDeleteDest(d: Destination) {
    if (!confirm(`Hapus destinasi "${d.nama_daerah}"?\nSemua risiko kesehatan terkait ikut dihapus.`)) return;
    try {
      await destinationsApi.del(d.id);
      if (selected?.id === d.id) setSelected(null);
      toast.success(`"${d.nama_daerah}" dihapus`);
      loadDest();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal menghapus destinasi");
    }
  }

  function openAddRisk() {
    setEditingRisk(null);
    setRiskForm({ nama_risiko_id: "", nama_risiko_en: "", saran_pencegahan_id: "", saran_pencegahan_en: "", rekomendasi_vaksinasi_id: "", rekomendasi_vaksinasi_en: "" });
    setRiskError("");
    setRiskModal(true);
  }

  function openEditRisk(r: HealthRisk) {
    setEditingRisk(r);
    setRiskForm({
      nama_risiko_id: r.nama_risiko_id,
      nama_risiko_en: r.nama_risiko_en,
      saran_pencegahan_id: r.saran_pencegahan_id ?? "",
      saran_pencegahan_en: r.saran_pencegahan_en ?? "",
      rekomendasi_vaksinasi_id: r.rekomendasi_vaksinasi_id ?? "",
      rekomendasi_vaksinasi_en: r.rekomendasi_vaksinasi_en ?? "",
    });
    setRiskError("");
    setRiskModal(true);
  }

  async function handleDeleteRisk(r: HealthRisk) {
    if (!confirm(`Hapus risiko "${r.nama_risiko_id}"?`)) return;
    try {
      await healthRisksApi.del(r.id);
      toast.success(`Risiko "${r.nama_risiko_id}" dihapus`);
      if (selected) loadRisks(selected);
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal menghapus risiko");
    }
  }

  async function handleSaveRisk(e: React.FormEvent) {
    e.preventDefault();
    setRiskError("");
    setRiskSaving(true);
    try {
      const payload = {
        destination_id: selected!.id,
        nama_risiko_id: riskForm.nama_risiko_id,
        nama_risiko_en: riskForm.nama_risiko_en,
        saran_pencegahan_id: riskForm.saran_pencegahan_id || undefined,
        saran_pencegahan_en: riskForm.saran_pencegahan_en || undefined,
        rekomendasi_vaksinasi_id: riskForm.rekomendasi_vaksinasi_id || undefined,
        rekomendasi_vaksinasi_en: riskForm.rekomendasi_vaksinasi_en || undefined,
      };
      if (editingRisk) {
        await healthRisksApi.update(editingRisk.id, payload);
        toast.success("Risiko diperbarui");
      } else {
        await healthRisksApi.create(payload);
        toast.success("Risiko ditambahkan");
      }
      setRiskModal(false);
      if (selected) loadRisks(selected);
    } catch (err: unknown) {
      setRiskError(err instanceof Error ? err.message : "Gagal menyimpan");
    } finally { setRiskSaving(false); }
  }

  const destColumns = [
    { key: "id", label: "ID", render: (d: Destination) => <span className="font-mono text-xs">{d.id}</span> },
    {
      key: "nama_daerah", label: "Nama Destinasi",
      render: (d: Destination) => (
        <button
          className="font-medium text-left hover:underline"
          style={{ color: "var(--brand)" }}
          onClick={() => loadRisks(d)}
        >
          {d.nama_daerah}
        </button>
      ),
    },
    {
      key: "actions", label: "",
      render: (d: Destination) => (
        <div className="flex items-center gap-1">
          <ActionBtn variant="edit" onClick={() => openEditDest(d)} />
          <ActionBtn variant="delete" onClick={() => handleDeleteDest(d)} />
        </div>
      ),
    },
  ];

  const riskColumns = [
    { key: "id", label: "ID", render: (r: HealthRisk) => <span className="font-mono text-xs">{r.id}</span> },
    { key: "nama_risiko_id", label: "Nama Risiko (ID)", render: (r: HealthRisk) => <span className="font-medium">{r.nama_risiko_id}</span> },
    { key: "nama_risiko_en", label: "Nama Risiko (EN)", render: (r: HealthRisk) => <span className="text-xs">{r.nama_risiko_en || "—"}</span> },
    {
      key: "saran_pencegahan_id", label: "Saran Pencegahan",
      render: (r: HealthRisk) => {
        const t = r.saran_pencegahan_id ?? "";
        return <span className="text-xs">{t.length > 60 ? t.slice(0, 60) + "…" : t || "—"}</span>;
      },
    },
    { key: "rekomendasi_vaksinasi_id", label: "Vaksinasi", render: (r: HealthRisk) => r.rekomendasi_vaksinasi_id ?? "—" },
    {
      key: "actions", label: "",
      render: (r: HealthRisk) => (
        <div className="flex items-center gap-1">
          <ActionBtn variant="edit"   onClick={() => openEditRisk(r)} />
          <ActionBtn variant="delete" onClick={() => handleDeleteRisk(r)} />
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Destinasi & Risiko Kesehatan"
        description="Klik nama destinasi untuk mengelola risiko kesehatannya"
        action={<AddButton onClick={openAddDest} label="Tambah Destinasi" />}
      />

      <div className="reveal">
        <DataTable columns={destColumns} rows={destinations} keyField="id" loading={loading} />
      </div>

      {selected && (
        <div className="reveal">
          <div className="flex items-center justify-between mb-3">
            <div>
              <h3 className="font-semibold text-base" style={{ color: "var(--text)" }}>
                Risiko Kesehatan — {selected.nama_daerah}
              </h3>
              <p className="text-xs mt-0.5" style={{ color: "var(--text-muted)" }}>
                Informasi ini ditampilkan kepada wisatawan saat mengecek risiko destinasi
              </p>
            </div>
            <AddButton onClick={openAddRisk} label="Tambah Risiko" />
          </div>
          <DataTable columns={riskColumns} rows={risks} keyField="id" loading={riskLoading} />
        </div>
      )}

      {/* Destination modal */}
      <Modal
        open={destModal}
        title={editingDest ? `Edit Destinasi — ${editingDest.nama_daerah}` : "Tambah Destinasi"}
        onClose={() => setDestModal(false)}
        size="sm"
      >
        <form onSubmit={handleSaveDest} className="space-y-4">
          <FormField label="Nama Daerah" required hint="Nama kabupaten/kecamatan di Bali">
            <Input
              value={destName}
              onChange={(e) => setDestName(e.target.value)}
              required
              placeholder="cth: Kabupaten Badung"
            />
          </FormField>
          {destError && (
            <p className="text-sm px-3 py-2 rounded-lg"
               style={{ color: "var(--danger)", background: "#fef2f2", border: "1px solid #fecaca" }}>
              {destError}
            </p>
          )}
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={() => setDestModal(false)}
              className="px-4 py-2 rounded-lg text-sm font-medium"
              style={{ color: "var(--text-muted)", border: "1px solid var(--border)" }}>
              Batal
            </button>
            <button type="submit" disabled={destSaving}
              className="px-4 py-2 rounded-lg text-sm font-semibold text-white disabled:opacity-60"
              style={{ background: "var(--brand-grad)" }}>
              {destSaving ? "Menyimpan…" : editingDest ? "Simpan Perubahan" : "Tambah"}
            </button>
          </div>
        </form>
      </Modal>

      {/* Risk modal */}
      <Modal
        open={riskModal}
        title={editingRisk
          ? `Edit Risiko — ${editingRisk.nama_risiko_id}`
          : `Tambah Risiko — ${selected?.nama_daerah ?? ""}`}
        onClose={() => setRiskModal(false)}
      >
        <form onSubmit={handleSaveRisk} className="space-y-4">
          <div className="grid grid-cols-2 gap-3">
            <FormField label="Nama Risiko (Indonesia)" required>
              <Input
                value={riskForm.nama_risiko_id}
                onChange={(e) => setRiskForm({ ...riskForm, nama_risiko_id: e.target.value })}
                required placeholder="cth: Demam Berdarah Dengue"
              />
            </FormField>
            <FormField label="Nama Risiko (English)" required>
              <Input
                value={riskForm.nama_risiko_en}
                onChange={(e) => setRiskForm({ ...riskForm, nama_risiko_en: e.target.value })}
                required placeholder="cth: Dengue Hemorrhagic Fever"
              />
            </FormField>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <FormField label="Saran Pencegahan (Indonesia)">
              <Textarea
                value={riskForm.saran_pencegahan_id}
                onChange={(e) => setRiskForm({ ...riskForm, saran_pencegahan_id: e.target.value })}
                placeholder="cth: Gunakan repelen nyamuk..."
              />
            </FormField>
            <FormField label="Saran Pencegahan (English)">
              <Textarea
                value={riskForm.saran_pencegahan_en}
                onChange={(e) => setRiskForm({ ...riskForm, saran_pencegahan_en: e.target.value })}
                placeholder="cth: Use mosquito repellent..."
              />
            </FormField>
          </div>
          <div className="grid grid-cols-2 gap-3">
            <FormField label="Rekomendasi Vaksinasi (Indonesia)">
              <Input
                value={riskForm.rekomendasi_vaksinasi_id}
                onChange={(e) => setRiskForm({ ...riskForm, rekomendasi_vaksinasi_id: e.target.value })}
                placeholder="cth: Vaksin Hepatitis A, Tifoid"
              />
            </FormField>
            <FormField label="Rekomendasi Vaksinasi (English)">
              <Input
                value={riskForm.rekomendasi_vaksinasi_en}
                onChange={(e) => setRiskForm({ ...riskForm, rekomendasi_vaksinasi_en: e.target.value })}
                placeholder="cth: Hepatitis A, Typhoid vaccine"
              />
            </FormField>
          </div>
          {riskError && (
            <p className="text-sm px-3 py-2 rounded-lg"
               style={{ color: "var(--danger)", background: "#fef2f2", border: "1px solid #fecaca" }}>
              {riskError}
            </p>
          )}
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={() => setRiskModal(false)}
              className="px-4 py-2 rounded-lg text-sm font-medium"
              style={{ color: "var(--text-muted)", border: "1px solid var(--border)" }}>
              Batal
            </button>
            <button type="submit" disabled={riskSaving}
              className="px-4 py-2 rounded-lg text-sm font-semibold text-white disabled:opacity-60"
              style={{ background: "var(--brand-grad)" }}>
              {riskSaving ? "Menyimpan…" : editingRisk ? "Simpan Perubahan" : "Tambah Risiko"}
            </button>
          </div>
        </form>
      </Modal>

      <Toaster toasts={toasts} remove={remove} />
    </div>
  );
}
