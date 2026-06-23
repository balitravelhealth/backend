"use client";

import { useEffect, useState, useCallback, useMemo } from "react";
import { rulesApi, symptomsApi, diseasesApi } from "@/lib/api";
import type { RekDefault } from "@/lib/api";
import type { Rule, Symptom, Disease } from "@/lib/types";
import DataTable from "@/components/DataTable";
import Modal from "@/components/Modal";
import PageHeader, { AddButton, ActionBtn } from "@/components/PageHeader";
import FormField, { Input, Select, Textarea } from "@/components/FormField";
import { useToast, Toaster } from "@/components/Toast";

type Tab = "rules" | "symptoms" | "diseases";

const KATEGORI = ["pre_travel", "post_travel"];
const RISK_LEVELS: (keyof RekDefault)[] = ["Rendah", "Sedang", "Tinggi", "Darurat"];
const RISK_COLORS: Record<keyof RekDefault, string> = {
  Rendah:  "var(--success)",
  Sedang:  "var(--warning)",
  Tinggi:  "#ea580c",
  Darurat: "var(--danger)",
};

// ── Sub-components ────────────────────────────────────────────────────────────

function CfSlider({
  label, value, onChange,
}: { label: string; value: number; onChange: (v: number) => void }) {
  return (
    <div>
      <div className="flex items-center justify-between mb-1">
        <label className="text-xs font-medium" style={{ color: "var(--text-muted)" }}>
          {label}
        </label>
        <input
          type="number" min={0} max={1} step={0.01}
          value={value}
          onChange={(e) =>
            onChange(Math.min(1, Math.max(0, parseFloat(e.target.value) || 0)))
          }
          className="w-16 px-2 py-1 text-xs rounded text-center font-mono"
          style={{ border: "1px solid var(--border)", color: "var(--text)" }}
        />
      </div>
      <input
        type="range" min={0} max={1} step={0.01}
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        className="w-full accent-blue-600"
      />
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
  const published = status === "published";
  return (
    <span
      className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold"
      style={{
        color:      published ? "var(--success)" : "var(--text-muted)",
        background: published ? "#f0fdf4" : "var(--bg)",
      }}
    >
      <span className="w-1.5 h-1.5 rounded-full" style={{ background: published ? "var(--success)" : "#94a3b8" }} />
      {published ? "Published" : "Draft"}
    </span>
  );
}

// ── Main page ─────────────────────────────────────────────────────────────────

const emptyRule = {
  nama: "", disease_id: 0, kategori: "pre_travel",
  bobot_cf: 0.7, mb: 0.8, md: 0.1,
  premis: [] as number[],
};

const emptyRek: RekDefault = { Rendah: "", Sedang: "", Tinggi: "", Darurat: "" };

export default function KnowledgeBasePage() {
  const [tab, setTab]           = useState<Tab>("rules");
  const [rules, setRules]       = useState<Rule[]>([]);
  const [symptoms, setSymptoms] = useState<Symptom[]>([]);
  const [diseases, setDiseases] = useState<Disease[]>([]);
  const [loading, setLoading]   = useState(true);
  const { toasts, remove, toast } = useToast();

  // ── rule state
  const [ruleOpen, setRuleOpen]       = useState(false);
  const [editingRule, setEditingRule] = useState<Rule | null>(null);
  const [ruleForm, setRuleForm]       = useState(emptyRule);
  const [symSearch, setSymSearch]     = useState("");
  const [ruleSaving, setRuleSaving]   = useState(false);
  const [ruleError, setRuleError]     = useState("");

  // ── symptom state
  const [symOpen, setSymOpen]         = useState(false);
  const [editingSym, setEditingSym]   = useState<Symptom | null>(null);
  const [symForm, setSymForm]         = useState({ kode: "", label_id: "", label_en: "" });
  const [symSaving, setSymSaving]     = useState(false);
  const [symError, setSymError]       = useState("");

  // ── disease state
  const [disOpen, setDisOpen]       = useState(false);
  const [editingDis, setEditingDis] = useState<Disease | null>(null);
  const [disForm, setDisForm]       = useState({ nama_id: "", nama_en: "", deskripsi_id: "", deskripsi_en: "" });
  const [rekForm, setRekForm]       = useState<RekDefault>(emptyRek);
  const [rekFormEN, setRekFormEN]   = useState<RekDefault>(emptyRek);
  const [disSaving, setDisSaving]   = useState(false);
  const [disError, setDisError]     = useState("");

  const loadAll = useCallback(async () => {
    setLoading(true);
    try {
      const [rRes, sRes, dRes] = await Promise.all([
        rulesApi.list(), symptomsApi.list(), diseasesApi.list(),
      ]);
      setRules(rRes.data ?? []);
      setSymptoms(sRes.data ?? []);
      setDiseases(dRes.data ?? []);
    } catch {
      toast.error("Gagal memuat knowledge base");
    } finally { setLoading(false); }
  }, [toast]);

  useEffect(() => { loadAll(); }, [loadAll]);

  const filteredSymptoms = useMemo(() =>
    symSearch.trim()
      ? symptoms.filter((s) =>
          s.label_id.toLowerCase().includes(symSearch.toLowerCase()) ||
          s.kode.toLowerCase().includes(symSearch.toLowerCase()),
        )
      : symptoms,
    [symptoms, symSearch],
  );

  // ── Rule handlers ──────────────────────────────────────────────────────────

  function openAddRule() {
    setEditingRule(null);
    setRuleForm(emptyRule);
    setSymSearch("");
    setRuleError("");
    setRuleOpen(true);
  }

  function openEditRule(r: Rule) {
    setEditingRule(r);
    setRuleForm({
      nama: r.nama,
      disease_id: r.disease_id,
      kategori: r.kategori,
      bobot_cf: r.bobot_cf,
      mb: r.mb,
      md: r.md,
      premis: Array.isArray(r.premis) ? r.premis : [],
    });
    setSymSearch("");
    setRuleError("");
    setRuleOpen(true);
  }

  function togglePremis(id: number) {
    setRuleForm((f) => ({
      ...f,
      premis: f.premis.includes(id)
        ? f.premis.filter((x) => x !== id)
        : [...f.premis, id],
    }));
  }

  async function handleSaveRule(e: React.FormEvent) {
    e.preventDefault();
    setRuleError("");

    // Validate premis
    if (ruleForm.premis.length === 0) {
      setRuleError("Pilih minimal 1 gejala sebagai premis rule");
      return;
    }
    // Validate disease selection
    if (!ruleForm.disease_id || ruleForm.disease_id === 0) {
      setRuleError("Pilih penyakit konklusi");
      return;
    }
    // Validate CF values
    if (ruleForm.bobot_cf < 0 || ruleForm.bobot_cf > 1) {
      setRuleError(`Bobot CF harus 0–1 (saat ini: ${ruleForm.bobot_cf.toFixed(3)}). Sesuaikan nilai MB dan MD.`);
      return;
    }

    setRuleSaving(true);
    try {
      const payload = {
        nama: ruleForm.nama, premis: ruleForm.premis,
        disease_id: ruleForm.disease_id, bobot_cf: ruleForm.bobot_cf,
        mb: ruleForm.mb, md: ruleForm.md, kategori: ruleForm.kategori,
      };
      if (editingRule) {
        await rulesApi.update(editingRule.rule_id, payload);
        toast.success("Rule diperbarui (status tetap)");
      } else {
        await rulesApi.create(payload);
        toast.success("Rule disimpan sebagai Draft — publish jika sudah diverifikasi");
      }
      setRuleOpen(false);
      loadAll();
    } catch (err: unknown) {
      setRuleError(err instanceof Error ? err.message : "Gagal menyimpan");
    } finally { setRuleSaving(false); }
  }

  async function handleTogglePublish(r: Rule) {
    try {
      if (r.status === "published") {
        await rulesApi.unpublish(r.rule_id);
        toast.success(`"${r.nama}" kembali ke Draft`);
      } else {
        await rulesApi.publish(r.rule_id);
        toast.success(`"${r.nama}" dipublish — aktif di sistem pakar`);
      }
      loadAll();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal mengubah status publish");
    }
  }

  async function handleDeleteRule(r: Rule) {
    if (!confirm(`Hapus rule "${r.nama}"?`)) return;
    try {
      await rulesApi.del(r.rule_id);
      toast.success(`Rule "${r.nama}" dihapus`);
      loadAll();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal menghapus rule");
    }
  }

  // ── Symptom handlers ───────────────────────────────────────────────────────

  function openAddSym() {
    setEditingSym(null);
    setSymForm({ kode: "", label_id: "", label_en: "" });
    setSymError("");
    setSymOpen(true);
  }

  function openEditSym(s: Symptom) {
    setEditingSym(s);
    setSymForm({ kode: s.kode, label_id: s.label_id, label_en: s.label_en ?? "" });
    setSymError("");
    setSymOpen(true);
  }

  async function handleSaveSym(e: React.FormEvent) {
    e.preventDefault();
    setSymError("");
    setSymSaving(true);
    try {
      if (editingSym) {
        await symptomsApi.update(editingSym.symptom_id, {
          kode: symForm.kode,
          label_id: symForm.label_id,
          label_en: symForm.label_en,
        });
        toast.success("Gejala diperbarui");
      } else {
        await symptomsApi.create({
          kode: symForm.kode,
          label_id: symForm.label_id,
          label_en: symForm.label_en,
        });
        toast.success("Gejala ditambahkan");
      }
      setSymOpen(false);
      loadAll();
    } catch (err: unknown) {
      setSymError(err instanceof Error ? err.message : "Gagal menyimpan");
    } finally { setSymSaving(false); }
  }

  async function handleDeleteSym(s: Symptom) {
    const rulesUsing = rules.filter((r) =>
      Array.isArray(r.premis) && r.premis.includes(s.symptom_id),
    );
    if (rulesUsing.length > 0) {
      alert(
        `Gejala ini dipakai oleh ${rulesUsing.length} rule:\n` +
          rulesUsing.map((r) => `• ${r.nama}`).join("\n") +
          "\n\nHapus atau edit rule tersebut dulu.",
      );
      return;
    }
    if (!confirm(`Hapus gejala "${s.label_id}"?`)) return;
    try {
      await symptomsApi.del(s.symptom_id);
      toast.success(`Gejala "${s.label_id}" dihapus`);
      loadAll();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal menghapus gejala");
    }
  }

  // ── Disease handlers ───────────────────────────────────────────────────────

  function openAddDis() {
    setEditingDis(null);
    setDisForm({ nama_id: "", nama_en: "", deskripsi_id: "", deskripsi_en: "" });
    setRekForm(emptyRek);
    setRekFormEN(emptyRek);
    setDisError("");
    setDisOpen(true);
  }

  function openEditDis(d: Disease) {
    setEditingDis(d);
    setDisForm({ nama_id: d.nama_id, nama_en: d.nama_en, deskripsi_id: d.deskripsi_id ?? "", deskripsi_en: d.deskripsi_en ?? "" });
    const parseRek = (raw: unknown): RekDefault => {
      try {
        const r = typeof raw === "string" ? JSON.parse(raw) : raw as Record<string, string>;
        return { Rendah: r.Rendah ?? "", Sedang: r.Sedang ?? "", Tinggi: r.Tinggi ?? "", Darurat: r.Darurat ?? "" };
      } catch { return emptyRek; }
    };
    setRekForm(d.rekomendasi_default_id ? parseRek(d.rekomendasi_default_id) : emptyRek);
    setRekFormEN(d.rekomendasi_default_en ? parseRek(d.rekomendasi_default_en) : emptyRek);
    setDisError("");
    setDisOpen(true);
  }

  async function handleSaveDis(e: React.FormEvent) {
    e.preventDefault();
    setDisError("");
    setDisSaving(true);
    try {
      const hasRekID = Object.values(rekForm).some((v) => v.trim());
      const hasRekEN = Object.values(rekFormEN).some((v) => v.trim());
      const body = {
        nama_id: disForm.nama_id,
        nama_en: disForm.nama_en,
        deskripsi_id: disForm.deskripsi_id || undefined,
        deskripsi_en: disForm.deskripsi_en || undefined,
        rekomendasi_default_id: hasRekID ? rekForm : undefined,
        rekomendasi_default_en: hasRekEN ? rekFormEN : undefined,
      };
      if (editingDis) {
        await diseasesApi.update(editingDis.id, body);
        toast.success("Penyakit diperbarui");
      } else {
        await diseasesApi.create(body);
        toast.success("Penyakit ditambahkan");
      }
      setDisOpen(false);
      loadAll();
    } catch (err: unknown) {
      setDisError(err instanceof Error ? err.message : "Gagal menyimpan");
    } finally { setDisSaving(false); }
  }

  async function handleDeleteDis(d: Disease) {
    const rulesUsing = rules.filter((r) => r.disease_id === d.id);
    if (rulesUsing.length > 0) {
      alert(
        `Penyakit ini dipakai oleh ${rulesUsing.length} rule:\n` +
          rulesUsing.map((r) => `• ${r.nama}`).join("\n") +
          "\n\nHapus rule tersebut dulu.",
      );
      return;
    }
    if (!confirm(`Hapus penyakit "${d.nama_id}"?`)) return;
    try {
      await diseasesApi.del(d.id);
      toast.success(`"${d.nama_id}" dihapus`);
      loadAll();
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "Gagal menghapus penyakit");
    }
  }

  // ── Table column defs ──────────────────────────────────────────────────────

  const ruleColumns = [
    { key: "rule_id", label: "ID",
      render: (r: Rule) => <span className="font-mono text-xs text-gray-400">{r.rule_id}</span> },
    { key: "nama", label: "Nama Rule",
      render: (r: Rule) => <span className="font-medium text-sm">{r.nama}</span> },
    { key: "disease_id", label: "Konklusi",
      render: (r: Rule) => (
        <span className="text-xs font-medium" style={{ color: "var(--brand)" }}>
          {diseases.find((d) => d.id === r.disease_id)?.nama_id ?? `#${r.disease_id}`}
        </span>
      ),
    },
    { key: "premis", label: "Premis (Gejala)",
      render: (r: Rule) => {
        const ids: number[] = Array.isArray(r.premis) ? r.premis : [];
        const names = ids.map((id) => symptoms.find((s) => s.symptom_id === id)?.label_id ?? `#${id}`);
        return (
          <div className="flex flex-wrap gap-1 max-w-xs">
            {names.slice(0, 3).map((n, i) => (
              <span key={i} className="text-xs px-1.5 py-0.5 rounded"
                style={{ background: "var(--bg)", color: "var(--text-muted)" }}>{n}</span>
            ))}
            {names.length > 3 && (
              <span className="text-xs px-1.5 py-0.5 rounded"
                style={{ background: "var(--bg)", color: "var(--text-muted)" }}>
                +{names.length - 3} lagi
              </span>
            )}
          </div>
        );
      },
    },
    { key: "kategori", label: "Kategori",
      render: (r: Rule) => (
        <span className="text-xs px-2 py-0.5 rounded-full"
          style={{ background: "var(--bg)", color: "var(--text-muted)" }}>
          {r.kategori}
        </span>
      ),
    },
    { key: "bobot_cf", label: "CF",
      render: (r: Rule) => <span className="font-mono text-xs">{r.bobot_cf.toFixed(2)}</span> },
    { key: "status", label: "Status",
      render: (r: Rule) => <StatusBadge status={r.status} /> },
    { key: "actions", label: "",
      render: (r: Rule) => (
        <div className="flex items-center gap-1">
          <button onClick={() => handleTogglePublish(r)}
            className="px-2 py-1 rounded-md text-xs font-medium transition-colors"
            style={{
              color:  r.status === "published" ? "var(--warning)" : "var(--success)",
              border: `1px solid ${r.status === "published" ? "#fcd34d" : "#86efac"}`,
            }}>
            {r.status === "published" ? "Unpublish" : "Publish"}
          </button>
          <ActionBtn variant="edit"   onClick={() => openEditRule(r)} />
          <ActionBtn variant="delete" onClick={() => handleDeleteRule(r)} />
        </div>
      ),
    },
  ];

  const symColumns = [
    { key: "symptom_id", label: "ID",
      render: (s: Symptom) => <span className="font-mono text-xs text-gray-400">{s.symptom_id}</span> },
    { key: "kode", label: "Kode",
      render: (s: Symptom) => (
        <code className="text-xs px-1.5 py-0.5 rounded"
          style={{ background: "var(--bg)", color: "var(--brand)" }}>
          {s.kode}
        </code>
      ),
    },
    { key: "label_id", label: "Nama (Indonesia)",
      render: (s: Symptom) => <span className="font-medium">{s.label_id}</span> },
    { key: "label_en", label: "Nama (English)",
      render: (s: Symptom) => s.label_en ?? <span style={{ color: "var(--text-muted)" }}>—</span> },
    { key: "used", label: "Dipakai di",
      render: (s: Symptom) => {
        const n = rules.filter((r) =>
          Array.isArray(r.premis) && r.premis.includes(s.symptom_id),
        ).length;
        return (
          <span className="text-xs" style={{ color: n > 0 ? "var(--brand)" : "var(--text-muted)" }}>
            {n > 0 ? `${n} rule` : "belum dipakai"}
          </span>
        );
      },
    },
    { key: "actions", label: "",
      render: (s: Symptom) => (
        <div className="flex items-center gap-1">
          <ActionBtn variant="edit" onClick={() => openEditSym(s)} />
          <ActionBtn variant="delete" onClick={() => handleDeleteSym(s)} />
        </div>
      ),
    },
  ];

  const disColumns = [
    { key: "id", label: "ID",
      render: (d: Disease) => <span className="font-mono text-xs text-gray-400">{d.id}</span> },
    { key: "nama_id", label: "Nama (Indonesia)",
      render: (d: Disease) => <span className="font-medium">{d.nama_id}</span> },
    { key: "nama_en", label: "Nama (English)",
      render: (d: Disease) => <span className="text-xs">{d.nama_en || "—"}</span> },
    { key: "deskripsi_id", label: "Deskripsi",
      render: (d: Disease) => {
        const t = d.deskripsi_id ?? "";
        return (
          <span className="text-xs" style={{ color: t ? "var(--text)" : "var(--text-muted)" }}>
            {t.length > 55 ? t.slice(0, 55) + "…" : t || "—"}
          </span>
        );
      },
    },
    { key: "rek", label: "Rekomendasi",
      render: (d: Disease) => {
        const hasRek = d.rekomendasi_default_id != null;
        return (
          <span className="text-xs" style={{ color: hasRek ? "var(--success)" : "var(--warning)" }}>
            {hasRek ? "✓ Lengkap" : "⚠ Belum diisi"}
          </span>
        );
      },
    },
    { key: "used", label: "Dipakai di",
      render: (d: Disease) => {
        const n = rules.filter((r) => r.disease_id === d.id).length;
        return (
          <span className="text-xs" style={{ color: n > 0 ? "var(--brand)" : "var(--text-muted)" }}>
            {n > 0 ? `${n} rule` : "belum dipakai"}
          </span>
        );
      },
    },
    { key: "actions", label: "",
      render: (d: Disease) => (
        <div className="flex items-center gap-1">
          <ActionBtn variant="edit"   onClick={() => openEditDis(d)} />
          <ActionBtn variant="delete" onClick={() => handleDeleteDis(d)} />
        </div>
      ),
    },
  ];

  const addActions: Record<Tab, () => void> = {
    rules: openAddRule,
    symptoms: openAddSym,
    diseases: openAddDis,
  };
  const addLabels: Record<Tab, string> = {
    rules: "Tambah Rule", symptoms: "Tambah Gejala", diseases: "Tambah Penyakit",
  };

  const tabStyle = (t: Tab) => ({
    color:        tab === t ? "var(--brand)" : "var(--text-muted)",
    borderBottom: tab === t ? "2px solid var(--brand)" : "2px solid transparent",
    paddingBottom: "10px",
    fontWeight:   tab === t ? 600 : 400,
  });

  // CF live values
  const cfLive   = (ruleForm.mb - ruleForm.md).toFixed(3);
  const cfOk     = ruleForm.mb >= 0 && ruleForm.mb <= 1 &&
                   ruleForm.md >= 0 && ruleForm.md <= 1 &&
                   ruleForm.bobot_cf >= 0 && ruleForm.bobot_cf <= 1;
  const cfNeg    = ruleForm.bobot_cf < 0;

  const publishedCount = rules.filter((r) => r.status === "published").length;

  // ── Render ─────────────────────────────────────────────────────────────────

  return (
    <div>
      <PageHeader
        title="Knowledge Base Sistem Pakar"
        description="Rule IF-THEN + Certainty Factor — admin dapat menambah/mengubah sesuai perkembangan ilmu kesehatan"
        action={<AddButton onClick={addActions[tab]} label={addLabels[tab]} />}
      />

      {/* Tab bar */}
      <div className="flex gap-6 mb-5 border-b" style={{ borderColor: "var(--border)" }}>
        {(["rules", "symptoms", "diseases"] as Tab[]).map((t) => (
          <button key={t} onClick={() => setTab(t)} className="text-sm transition-colors" style={tabStyle(t)}>
            {t === "rules"
              ? `Rules (${publishedCount} published / ${rules.length} total)`
              : t === "symptoms"
              ? `Gejala Master (${symptoms.length})`
              : `Penyakit Master (${diseases.length})`}
          </button>
        ))}
      </div>

      <div className="reveal">
        {tab === "rules"    && <DataTable columns={ruleColumns} rows={rules}    keyField="rule_id"   loading={loading} />}
        {tab === "symptoms" && <DataTable columns={symColumns}  rows={symptoms} keyField="symptom_id" loading={loading} />}
        {tab === "diseases" && <DataTable columns={disColumns}  rows={diseases} keyField="id"         loading={loading} />}
      </div>

      {/* ═══════════════════════ RULE MODAL ═══════════════════════ */}
      <Modal open={ruleOpen} title={editingRule ? `Edit Rule — ${editingRule.nama}` : "Tambah Rule Baru"} onClose={() => setRuleOpen(false)} size="lg">
        <form onSubmit={handleSaveRule} className="space-y-5">

          {/* Panduan */}
          <div className="px-3 py-2.5 rounded-lg text-xs leading-relaxed" style={{ background: "#eff6ff", color: "#1d4ed8", border: "1px solid #bfdbfe" }}>
            <strong>Panduan Rule:</strong> Pilih gejala (premis) yang menjadi kondisi, lalu tentukan penyakit yang didiagnosis (konklusi).
            Nilai MB = keyakinan hipotesis benar, MD = keyakinan hipotesis salah. CF = MB − MD.
            Rule baru disimpan sebagai <em>Draft</em> — publish setelah diverifikasi perawat.
          </div>

          <div className="grid grid-cols-2 gap-4">
            <FormField label="Nama Rule" required hint="Deskripsi singkat kondisi yang ditangkap rule ini">
              <Input value={ruleForm.nama}
                onChange={(e) => setRuleForm({ ...ruleForm, nama: e.target.value })}
                required placeholder="cth: DBD dengan gejala berat" />
            </FormField>
            <FormField label="Kategori" required hint="pre_travel = saat wisata, post_travel = setelah kembali">
              <Select value={ruleForm.kategori}
                onChange={(e) => setRuleForm({ ...ruleForm, kategori: e.target.value })} required>
                {KATEGORI.map((k) => <option key={k} value={k}>{k}</option>)}
              </Select>
            </FormField>
          </div>

          {/* Konklusi */}
          <FormField label="Konklusi — Penyakit yang Didiagnosis" required hint="Pilih dari master penyakit — tidak bisa ketik bebas">
            <Select
              value={ruleForm.disease_id || ""}
              onChange={(e) => setRuleForm({ ...ruleForm, disease_id: Number(e.target.value) })}
              required
            >
              <option value="">— pilih penyakit konklusi —</option>
              {diseases.map((d) => <option key={d.id} value={d.id}>{d.nama_id}</option>)}
            </Select>
          </FormField>

          {/* Premis */}
          <FormField
            label="Premis — Gejala yang Harus Muncul"
            required
            hint={ruleForm.premis.length === 0
              ? "⚠ Pilih minimal 1 gejala"
              : `${ruleForm.premis.length} gejala dipilih`}
          >
            <div className="rounded-xl overflow-hidden" style={{ border: `1px solid ${ruleForm.premis.length === 0 ? "var(--warning)" : "var(--border)"}` }}>
              {/* Search */}
              <div className="px-3 py-2 border-b" style={{ borderColor: "var(--border)", background: "var(--bg)" }}>
                <input
                  type="text"
                  placeholder="Cari gejala…"
                  value={symSearch}
                  onChange={(e) => setSymSearch(e.target.value)}
                  className="w-full text-sm bg-transparent outline-none"
                  style={{ color: "var(--text)" }}
                />
              </div>
              {/* Selected chips */}
              {ruleForm.premis.length > 0 && (
                <div className="flex flex-wrap gap-1.5 px-3 py-2 border-b" style={{ borderColor: "var(--border)", background: "#eff6ff" }}>
                  {ruleForm.premis.map((id) => {
                    const s = symptoms.find((x) => x.symptom_id === id);
                    return (
                      <span key={id}
                        className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium"
                        style={{ background: "var(--brand)", color: "#fff" }}>
                        {s?.label_id ?? `#${id}`}
                        <button type="button" onClick={() => togglePremis(id)} className="hover:opacity-70">×</button>
                      </span>
                    );
                  })}
                </div>
              )}
              {/* Symptom grid */}
              <div className="grid grid-cols-2 gap-0 max-h-44 overflow-y-auto">
                {filteredSymptoms.length === 0 ? (
                  <div className="col-span-2 py-6 text-center text-sm" style={{ color: "var(--text-muted)" }}>
                    Tidak ada gejala yang cocok
                  </div>
                ) : (
                  filteredSymptoms.map((s) => {
                    const checked = ruleForm.premis.includes(s.symptom_id);
                    return (
                      <label key={s.symptom_id}
                        className="flex items-center gap-2 px-3 py-2 text-sm cursor-pointer transition-colors"
                        style={{ background: checked ? "#eff6ff" : "transparent" }}
                        onMouseEnter={(e) => { if (!checked) (e.currentTarget as HTMLElement).style.background = "var(--bg)"; }}
                        onMouseLeave={(e) => { if (!checked) (e.currentTarget as HTMLElement).style.background = "transparent"; }}
                      >
                        <input type="checkbox" checked={checked}
                          onChange={() => togglePremis(s.symptom_id)}
                          className="accent-blue-600 flex-shrink-0" />
                        <span style={{ color: "var(--text)" }}>{s.label_id}</span>
                      </label>
                    );
                  })
                )}
              </div>
            </div>
          </FormField>

          {/* CF panel */}
          <div className="rounded-xl p-4 space-y-3"
            style={{ background: "var(--bg)", border: `1px solid ${cfNeg ? "var(--danger)" : "var(--border)"}` }}>
            <div className="flex items-center justify-between">
              <p className="text-xs font-semibold uppercase tracking-wide" style={{ color: "var(--text-muted)" }}>
                Certainty Factor
              </p>
              <div className="flex items-center gap-2 text-xs font-mono">
                <span style={{ color: "var(--text-muted)" }}>MB − MD =</span>
                <span
                  className="font-bold px-2 py-0.5 rounded"
                  style={{
                    color:      cfOk ? "var(--brand)" : "var(--danger)",
                    background: cfOk ? "#eff6ff" : "#fef2f2",
                  }}
                >
                  {cfLive}
                </span>
                {cfNeg && (
                  <span style={{ color: "var(--danger)" }}>⚠ CF negatif, tidak valid!</span>
                )}
              </div>
            </div>

            <CfSlider label="MB — Measure of Belief (keyakinan hipotesis BENAR, 0–1)"
              value={ruleForm.mb}
              onChange={(v) => setRuleForm({ ...ruleForm, mb: v, bobot_cf: parseFloat((v - ruleForm.md).toFixed(3)) })} />
            <CfSlider label="MD — Measure of Disbelief (keyakinan hipotesis SALAH, 0–1)"
              value={ruleForm.md}
              onChange={(v) => setRuleForm({ ...ruleForm, md: v, bobot_cf: parseFloat((ruleForm.mb - v).toFixed(3)) })} />
            <CfSlider label="Bobot CF akhir (dihitung otomatis dari MB − MD, bisa override)"
              value={Math.max(0, Math.min(1, ruleForm.bobot_cf))}
              onChange={(v) => setRuleForm({ ...ruleForm, bobot_cf: v })} />

            <div className="text-xs pt-1 border-t space-y-0.5" style={{ borderColor: "var(--border)", color: "var(--text-muted)" }}>
              <p>CF ≥0.8 → Darurat · ≥0.6 → Tinggi · ≥0.4 → Sedang · &lt;0.4 → Rendah</p>
              <p>Nilai CF harus 0–1. Server akan menolak nilai di luar batas.</p>
            </div>
          </div>

          {ruleError && (
            <p className="text-sm px-3 py-2 rounded-lg"
              style={{ color: "var(--danger)", background: "#fef2f2", border: "1px solid #fecaca" }}>
              {ruleError}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-1">
            <button type="button" onClick={() => setRuleOpen(false)}
              className="px-4 py-2 rounded-lg text-sm font-medium"
              style={{ color: "var(--text-muted)", border: "1px solid var(--border)" }}>
              Batal
            </button>
            <button type="submit" disabled={ruleSaving || !cfOk}
              className="px-4 py-2 rounded-lg text-sm font-semibold text-white disabled:opacity-60"
              style={{ background: "var(--brand-grad)" }}>
              {ruleSaving
                ? "Menyimpan…"
                : editingRule
                ? "Simpan Perubahan"
                : "Simpan sebagai Draft"}
            </button>
          </div>
        </form>
      </Modal>

      {/* ═══════════════════════ SYMPTOM MODAL ═══════════════════════ */}
      <Modal
        open={symOpen}
        title={editingSym ? `Edit Gejala — ${editingSym.label_id}` : "Tambah Gejala Baru"}
        onClose={() => setSymOpen(false)}
        size="sm"
      >
        <form onSubmit={handleSaveSym} className="space-y-4">
          <FormField label="Kode Unik" required hint="Format S_NAMA — huruf besar dan underscore. cth: S_DEMAM_TINGGI">
            <Input value={symForm.kode}
              onChange={(e) => setSymForm({ ...symForm, kode: e.target.value.toUpperCase().replace(/[^A-Z0-9_]/g, "") })}
              required placeholder="S_DEMAM_TINGGI"
              disabled={editingSym !== null} />
          </FormField>
          {editingSym && (
            <p className="text-xs -mt-2" style={{ color: "var(--text-muted)" }}>
              Kode tidak bisa diubah karena sudah dipakai sebagai referensi di rule.
            </p>
          )}
          <FormField label="Nama Gejala (Indonesia)" required>
            <Input value={symForm.label_id}
              onChange={(e) => setSymForm({ ...symForm, label_id: e.target.value })}
              required placeholder="cth: Demam tinggi mendadak (≥38.5°C)" />
          </FormField>
          <FormField label="Nama Gejala (English)" hint="Opsional — untuk internasionalisasi">
            <Input value={symForm.label_en}
              onChange={(e) => setSymForm({ ...symForm, label_en: e.target.value })}
              placeholder="cth: High fever (≥38.5°C)" />
          </FormField>
          {symError && (
            <p className="text-sm px-3 py-2 rounded-lg"
              style={{ color: "var(--danger)", background: "#fef2f2", border: "1px solid #fecaca" }}>
              {symError}
            </p>
          )}
          <div className="flex justify-end gap-2 pt-1">
            <button type="button" onClick={() => setSymOpen(false)}
              className="px-4 py-2 rounded-lg text-sm font-medium"
              style={{ color: "var(--text-muted)", border: "1px solid var(--border)" }}>
              Batal
            </button>
            <button type="submit" disabled={symSaving}
              className="px-4 py-2 rounded-lg text-sm font-semibold text-white disabled:opacity-60"
              style={{ background: "var(--brand-grad)" }}>
              {symSaving ? "Menyimpan…" : editingSym ? "Simpan Perubahan" : "Tambah Gejala"}
            </button>
          </div>
        </form>
      </Modal>

      {/* ═══════════════════════ DISEASE MODAL ═══════════════════════ */}
      <Modal
        open={disOpen}
        title={editingDis ? `Edit Penyakit — ${editingDis.nama_id}` : "Tambah Penyakit Baru"}
        onClose={() => setDisOpen(false)}
        size="lg"
      >
        <form onSubmit={handleSaveDis} className="space-y-5">
          <div className="grid grid-cols-2 gap-4">
            <FormField label="Nama Penyakit (Indonesia)" required>
              <Input value={disForm.nama_id}
                onChange={(e) => setDisForm({ ...disForm, nama_id: e.target.value })}
                required placeholder="cth: Demam Berdarah Dengue (DBD)" />
            </FormField>
            <FormField label="Nama Penyakit (English)" required>
              <Input value={disForm.nama_en}
                onChange={(e) => setDisForm({ ...disForm, nama_en: e.target.value })}
                required placeholder="cth: Dengue Hemorrhagic Fever (DHF)" />
            </FormField>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <FormField label="Deskripsi (Indonesia)">
              <Input value={disForm.deskripsi_id}
                onChange={(e) => setDisForm({ ...disForm, deskripsi_id: e.target.value })}
                placeholder="cth: Infeksi virus dengue yang ditularkan nyamuk Aedes" />
            </FormField>
            <FormField label="Deskripsi (English)">
              <Input value={disForm.deskripsi_en}
                onChange={(e) => setDisForm({ ...disForm, deskripsi_en: e.target.value })}
                placeholder="cth: Dengue virus infection transmitted by Aedes mosquito" />
            </FormField>
          </div>

          {/* Rekomendasi per risk level — Indonesian */}
          <div className="rounded-xl overflow-hidden" style={{ border: "1px solid var(--border)" }}>
            <div className="px-4 py-2.5 border-b flex items-center justify-between"
              style={{ background: "var(--bg)", borderColor: "var(--border)" }}>
              <p className="text-xs font-semibold uppercase tracking-wide" style={{ color: "var(--text-muted)" }}>
                Rekomendasi per Level Risiko — Indonesia
              </p>
            </div>
            <div className="divide-y" style={{ borderColor: "var(--border)" }}>
              {RISK_LEVELS.map((level) => (
                <div key={level} className="px-4 py-3">
                  <div className="flex items-center gap-2 mb-1.5">
                    <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: RISK_COLORS[level] }} />
                    <label className="text-sm font-semibold" style={{ color: RISK_COLORS[level] }}>{level}</label>
                  </div>
                  <Textarea rows={2} value={rekForm[level]}
                    onChange={(e) => setRekForm({ ...rekForm, [level]: e.target.value })}
                    placeholder="cth: Istirahat, minum cairan cukup..." />
                </div>
              ))}
            </div>
          </div>

          {/* Rekomendasi per risk level — English */}
          <div className="rounded-xl overflow-hidden" style={{ border: "1px solid var(--border)" }}>
            <div className="px-4 py-2.5 border-b flex items-center justify-between"
              style={{ background: "var(--bg)", borderColor: "var(--border)" }}>
              <p className="text-xs font-semibold uppercase tracking-wide" style={{ color: "var(--text-muted)" }}>
                Rekomendasi per Level Risiko — English
              </p>
            </div>
            <div className="divide-y" style={{ borderColor: "var(--border)" }}>
              {RISK_LEVELS.map((level) => (
                <div key={level} className="px-4 py-3">
                  <div className="flex items-center gap-2 mb-1.5">
                    <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: RISK_COLORS[level] }} />
                    <label className="text-sm font-semibold" style={{ color: RISK_COLORS[level] }}>{level}</label>
                  </div>
                  <Textarea rows={2} value={rekFormEN[level]}
                    onChange={(e) => setRekFormEN({ ...rekFormEN, [level]: e.target.value })}
                    placeholder="cth: Rest, drink enough fluids..." />
                </div>
              ))}
            </div>
            <div className="px-4 py-2.5 border-t" style={{ borderColor: "var(--border)", background: "var(--bg)" }}>
              <p className="text-xs" style={{ color: "var(--warning)" }}>
                ⚠ Konten medis wajib diverifikasi oleh perawat berlisensi sebelum digunakan di produksi.
              </p>
            </div>
          </div>

          {disError && (
            <p className="text-sm px-3 py-2 rounded-lg"
              style={{ color: "var(--danger)", background: "#fef2f2", border: "1px solid #fecaca" }}>
              {disError}
            </p>
          )}

          <div className="flex justify-end gap-2 pt-1">
            <button type="button" onClick={() => setDisOpen(false)}
              className="px-4 py-2 rounded-lg text-sm font-medium"
              style={{ color: "var(--text-muted)", border: "1px solid var(--border)" }}>
              Batal
            </button>
            <button type="submit" disabled={disSaving}
              className="px-4 py-2 rounded-lg text-sm font-semibold text-white disabled:opacity-60"
              style={{ background: "var(--brand-grad)" }}>
              {disSaving ? "Menyimpan…" : editingDis ? "Simpan Perubahan" : "Tambah Penyakit"}
            </button>
          </div>
        </form>
      </Modal>

      <Toaster toasts={toasts} remove={remove} />
    </div>
  );
}
