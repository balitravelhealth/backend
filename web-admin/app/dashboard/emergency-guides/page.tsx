"use client";

import { useEffect, useState, useCallback, useRef } from "react";
import { emergencyFlowsApi, uploadImage } from "@/lib/api";
import type { EmergencyGuideFlow, GuideNode, GuideChoice } from "@/lib/types";

// ── constants ─────────────────────────────────────────────────────────────────

const KATEGORI = ["BLS", "CPR_ANAK", "TERSEDAK", "LUKA", "ALERGI", "DARURAT", "LAINNYA"];

const VARIANT_CFG = {
  yes:     { label: "Ya ✓",     cls: "bg-green-100 text-green-700 border-green-300 ring-green-200" },
  no:      { label: "Tidak ✗",  cls: "bg-red-100   text-red-700   border-red-300   ring-red-200" },
  neutral: { label: "Lainnya →", cls: "bg-gray-100  text-gray-700  border-gray-300  ring-gray-200" },
} as const;

// ── tiny helpers ──────────────────────────────────────────────────────────────

function uid() {
  return Math.random().toString(36).slice(2, 8);
}

function makeNode(isEntry = false): GuideNode {
  return { id: uid(), title: "", instruction: "", image_url: "", is_entry: isEntry, choices: [] };
}

// ── ImageUploader ─────────────────────────────────────────────────────────────

function ImageUploader({
  value,
  onChange,
}: {
  value: string;
  onChange: (url: string) => void;
}) {
  const [uploading, setUploading] = useState(false);
  const [err, setErr]             = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  async function handleFile(file: File) {
    if (!file) return;
    setErr("");
    setUploading(true);
    try {
      const url = await uploadImage(file);
      onChange(url);
    } catch (e: unknown) {
      setErr(e instanceof Error ? e.message : "Upload gagal");
    } finally {
      setUploading(false);
    }
  }

  function onDrop(e: React.DragEvent) {
    e.preventDefault();
    const file = e.dataTransfer.files[0];
    if (file) handleFile(file);
  }

  function onPaste(e: React.ClipboardEvent) {
    const file = Array.from(e.clipboardData.files).find(f => f.type.startsWith("image/"));
    if (file) handleFile(file);
  }

  return (
    <div className="space-y-2">
      {/* drop zone */}
      <div
        onDrop={onDrop}
        onDragOver={e => e.preventDefault()}
        onPaste={onPaste}
        onClick={() => !uploading && inputRef.current?.click()}
        className="relative border-2 border-dashed border-[var(--border)] rounded-xl overflow-hidden cursor-pointer hover:border-[var(--brand)] transition-colors group"
      >
        {value ? (
          /* preview */
          <div className="relative">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={value}
              alt="preview"
              className="w-full max-h-48 object-cover"
              onError={e => { (e.target as HTMLImageElement).style.display = "none"; }}
            />
            <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity">
              <span className="text-white text-sm font-medium">Klik untuk ganti gambar</span>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center gap-2 py-8 px-4 text-center">
            <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center text-2xl">🖼️</div>
            <p className="text-sm font-medium text-[var(--text)]">
              {uploading ? "Mengunggah…" : "Pilih gambar atau seret ke sini"}
            </p>
            <p className="text-xs text-[var(--text-muted)]">JPEG, PNG, WebP, GIF · Maks 5 MB</p>
            <p className="text-xs text-[var(--text-muted)]">Bisa juga paste gambar (Ctrl+V)</p>
          </div>
        )}
        {uploading && (
          <div className="absolute inset-0 bg-white/70 flex items-center justify-center">
            <div className="flex items-center gap-2 text-sm font-medium text-[var(--brand)]">
              <svg className="animate-spin w-5 h-5" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" strokeDasharray="30" strokeDashoffset="10" />
              </svg>
              Mengunggah…
            </div>
          </div>
        )}
      </div>

      {/* URL manual */}
      <div className="flex gap-2 items-center">
        <span className="text-xs text-[var(--text-muted)] shrink-0">atau URL:</span>
        <input
          value={value}
          onChange={e => onChange(e.target.value)}
          placeholder="https://..."
          className="flex-1 border border-[var(--border)] rounded-lg px-3 py-1.5 text-xs focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
        />
        {value && (
          <button
            onClick={() => onChange("")}
            className="text-xs text-red-400 hover:text-red-600 px-2 py-1 border border-red-200 rounded-lg"
          >hapus</button>
        )}
      </div>

      {err && <p className="text-xs text-red-600">⚠ {err}</p>}

      <input
        ref={inputRef}
        type="file"
        accept="image/jpeg,image/png,image/webp,image/gif"
        className="hidden"
        onChange={e => { const f = e.target.files?.[0]; if (f) handleFile(f); e.target.value = ""; }}
      />
    </div>
  );
}

// ── ChoiceRow ─────────────────────────────────────────────────────────────────

function ChoiceRow({
  choice, index, allNodes, nodeId,
  onChange, onDelete,
}: {
  choice: GuideChoice;
  index: number;
  allNodes: GuideNode[];
  nodeId: string;
  onChange: (c: GuideChoice) => void;
  onDelete: () => void;
}) {
  const targets = allNodes.filter(n => n.id !== nodeId);
  const cfg = VARIANT_CFG[choice.variant];

  return (
    <div className="flex items-center gap-2 p-2.5 bg-[var(--bg)] rounded-xl border border-[var(--border)]">
      {/* variant pill */}
      <select
        value={choice.variant}
        onChange={e => onChange({ ...choice, variant: e.target.value as GuideChoice["variant"] })}
        className={`text-xs font-medium rounded-lg px-2 py-1.5 border ${cfg.cls} focus:outline-none focus:ring-2`}
      >
        {(Object.keys(VARIANT_CFG) as Array<keyof typeof VARIANT_CFG>).map(v => (
          <option key={v} value={v}>{VARIANT_CFG[v].label}</option>
        ))}
      </select>

      {/* label */}
      <input
        value={choice.label}
        onChange={e => onChange({ ...choice, label: e.target.value })}
        placeholder={`Teks tombol pilihan ${index + 1}`}
        className="flex-1 text-sm border border-[var(--border)] rounded-lg px-3 py-1.5 focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
      />

      {/* arrow */}
      <span className="text-[var(--text-muted)] text-sm shrink-0">→</span>

      {/* next_id */}
      <select
        value={choice.next_id ?? ""}
        onChange={e => onChange({ ...choice, next_id: e.target.value || null })}
        className="text-xs border border-[var(--border)] rounded-lg px-2 py-1.5 max-w-[160px] focus:outline-none focus:ring-2 focus:ring-[var(--brand)] bg-white"
      >
        <option value="">⛔ Akhiri panduan</option>
        {targets.map(n => (
          <option key={n.id} value={n.id}>
            {n.title ? `${n.title}` : n.id}
          </option>
        ))}
      </select>

      <button
        onClick={onDelete}
        className="shrink-0 w-6 h-6 rounded-full flex items-center justify-center text-red-400 hover:bg-red-50 hover:text-red-600 transition-colors text-base"
      >×</button>
    </div>
  );
}

// ── NodeEditor (single node, full panel) ─────────────────────────────────────

function NodeEditor({
  node, allNodes, onChange,
}: {
  node: GuideNode;
  allNodes: GuideNode[];
  onChange: (n: GuideNode) => void;
}) {
  function set<K extends keyof GuideNode>(k: K, v: GuideNode[K]) {
    onChange({ ...node, [k]: v });
  }

  function setEntry(checked: boolean) {
    // when marking as entry, parent will unmark others — just set it
    onChange({ ...node, is_entry: checked });
  }

  function addChoice() {
    const nextVariant: GuideChoice["variant"] =
      node.choices.filter(c => c.variant !== "neutral").length === 0 ? "yes"
        : node.choices.filter(c => c.variant === "yes").length === 0 ? "yes"
        : node.choices.filter(c => c.variant === "no").length === 0 ? "no"
        : "neutral";
    onChange({ ...node, choices: [...node.choices, { label: "", next_id: null, variant: nextVariant }] });
  }

  function updateChoice(i: number, c: GuideChoice) {
    onChange({ ...node, choices: node.choices.map((old, idx) => idx === i ? c : old) });
  }

  function removeChoice(i: number) {
    onChange({ ...node, choices: node.choices.filter((_, idx) => idx !== i) });
  }

  return (
    <div className="flex flex-col h-full">
      {/* node meta bar */}
      <div className={`flex items-center gap-3 px-5 py-3 border-b border-[var(--border)] ${node.is_entry ? "bg-blue-50" : "bg-[var(--bg)]"}`}>
        <div className="flex-1 flex items-center gap-2">
          <span className="text-xs text-[var(--text-muted)]">ID:</span>
          <input
            value={node.id}
            onChange={e => set("id", e.target.value.toLowerCase().replace(/[^a-z0-9_]/g, "_"))}
            className="font-mono text-sm bg-white border border-[var(--border)] rounded-lg px-2 py-0.5 w-44 focus:outline-none focus:ring-2 focus:ring-[var(--brand)]"
            placeholder="node_id"
          />
        </div>
        <label className="flex items-center gap-2 cursor-pointer select-none">
          <div
            onClick={() => setEntry(!node.is_entry)}
            className={`w-10 h-5 rounded-full transition-colors relative ${node.is_entry ? "bg-blue-500" : "bg-gray-300"}`}
          >
            <div className={`absolute top-0.5 w-4 h-4 rounded-full bg-white shadow transition-transform ${node.is_entry ? "translate-x-5" : "translate-x-0.5"}`} />
          </div>
          <span className={`text-xs font-medium ${node.is_entry ? "text-blue-600" : "text-[var(--text-muted)]"}`}>
            {node.is_entry ? "🔵 Entry Point" : "Entry Point"}
          </span>
        </label>
      </div>

      <div className="flex-1 overflow-y-auto p-5 space-y-5">

        {/* JUDUL */}
        <div>
          <label className="block text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wide mb-1.5">
            Judul Node
          </label>
          <input
            value={node.title}
            onChange={e => set("title", e.target.value)}
            placeholder="Contoh: Periksa Respons Korban"
            className="w-full border border-[var(--border)] rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand)] bg-white"
          />
        </div>

        {/* INSTRUKSI */}
        <div>
          <label className="block text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wide mb-1.5">
            Instruksi / Langkah-langkah
          </label>
          <textarea
            value={node.instruction}
            onChange={e => set("instruction", e.target.value)}
            rows={6}
            placeholder={"Tulis langkah-langkah penanganan di sini.\n\nContoh:\n1) Pastikan area sekitar aman\n2) Tepuk bahu korban sambil memanggil namanya\n3) Perhatikan respon korban"}
            className="w-full border border-[var(--border)] rounded-xl px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--brand)] resize-y bg-white leading-relaxed"
          />
        </div>

        {/* GAMBAR */}
        <div>
          <label className="block text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wide mb-1.5">
            Gambar Ilustrasi (opsional)
          </label>
          <ImageUploader value={node.image_url} onChange={v => set("image_url", v)} />
        </div>

        {/* PILIHAN / CABANG */}
        <div>
          <div className="flex items-center justify-between mb-3">
            <div>
              <label className="block text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wide">
                Pilihan Jawaban
              </label>
              <p className="text-xs text-[var(--text-muted)] mt-0.5">
                {node.choices.length === 0
                  ? "Tidak ada pilihan → node akhir panduan"
                  : `${node.choices.length} pilihan — klik tombol untuk melanjutkan ke node berikutnya`}
              </p>
            </div>
            <button
              onClick={addChoice}
              className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-[var(--brand)] text-white rounded-lg hover:opacity-90 transition-opacity"
            >
              + Tambah Pilihan
            </button>
          </div>

          {node.choices.length === 0 ? (
            <div className="border-2 border-dashed border-[var(--border)] rounded-xl p-5 text-center">
              <p className="text-2xl mb-1">🏁</p>
              <p className="text-sm text-[var(--text-muted)]">Node ini adalah <strong>akhir panduan</strong></p>
              <p className="text-xs text-[var(--text-muted)] mt-0.5">Klik "+ Tambah Pilihan" jika ada percabangan</p>
            </div>
          ) : (
            <div className="space-y-2">
              {node.choices.map((ch, i) => (
                <ChoiceRow
                  key={i}
                  choice={ch}
                  index={i}
                  allNodes={allNodes}
                  nodeId={node.id}
                  onChange={c => updateChoice(i, c)}
                  onDelete={() => removeChoice(i)}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// ── NodeListItem ──────────────────────────────────────────────────────────────

function NodeListItem({
  node, isActive, index, total, onClick, onMoveUp, onMoveDown, onDelete,
}: {
  node: GuideNode;
  isActive: boolean;
  index: number;
  total: number;
  onClick: () => void;
  onMoveUp: () => void;
  onMoveDown: () => void;
  onDelete: () => void;
}) {
  return (
    <div
      className={`group flex items-stretch rounded-xl border transition-all cursor-pointer
        ${isActive
          ? "border-[var(--brand)] bg-blue-50 shadow-sm"
          : "border-[var(--border)] bg-[var(--surface)] hover:border-gray-300 hover:bg-[var(--bg)]"
        }`}
      onClick={onClick}
    >
      {/* colour strip */}
      <div className={`w-1 shrink-0 rounded-l-xl ${node.is_entry ? "bg-blue-500" : isActive ? "bg-blue-300" : "bg-gray-200"}`} />

      <div className="flex-1 py-2.5 px-3 min-w-0">
        <div className="flex items-center gap-1.5">
          {node.is_entry && <span className="text-[10px] font-bold text-blue-600 bg-blue-100 px-1.5 py-0.5 rounded">START</span>}
          <span className="text-sm font-medium text-[var(--text)] truncate">
            {node.title || <span className="italic text-[var(--text-muted)]">(belum ada judul)</span>}
          </span>
        </div>
        <p className="text-[11px] text-[var(--text-muted)] mt-0.5 truncate font-mono">{node.id}</p>
        <p className="text-[11px] text-[var(--text-muted)] mt-0.5">
          {node.choices.length === 0 ? "🏁 akhir" : `${node.choices.length} pilihan`}
          {node.image_url ? " · 🖼️" : ""}
        </p>
      </div>

      {/* controls */}
      <div className="flex flex-col justify-center gap-0.5 pr-2 opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          onClick={e => { e.stopPropagation(); onMoveUp(); }}
          disabled={index === 0}
          className="w-5 h-5 flex items-center justify-center rounded text-[var(--text-muted)] hover:bg-gray-200 disabled:opacity-20 text-xs"
        >▲</button>
        <button
          onClick={e => { e.stopPropagation(); onMoveDown(); }}
          disabled={index === total - 1}
          className="w-5 h-5 flex items-center justify-center rounded text-[var(--text-muted)] hover:bg-gray-200 disabled:opacity-20 text-xs"
        >▼</button>
        {total > 1 && (
          <button
            onClick={e => { e.stopPropagation(); onDelete(); }}
            className="w-5 h-5 flex items-center justify-center rounded text-red-400 hover:bg-red-50 text-xs"
          >×</button>
        )}
      </div>
    </div>
  );
}

// ── FlowMap ───────────────────────────────────────────────────────────────────

function FlowMap({ nodes }: { nodes: GuideNode[] }) {
  const entry = nodes.find(n => n.is_entry) ?? nodes[0];
  if (!entry) return <p className="text-sm text-[var(--text-muted)] p-4">Tidak ada node.</p>;

  function Tree({ nodeId, visited }: { nodeId: string; visited: Set<string> }) {
    if (visited.has(nodeId)) return <span className="text-xs text-orange-500 italic">(loop ke {nodeId})</span>;
    const n = nodes.find(x => x.id === nodeId);
    if (!n) return <span className="text-xs text-red-500">node "{nodeId}" tidak ditemukan</span>;
    const next = new Set(visited); next.add(nodeId);
    return (
      <div>
        <div className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-medium border shadow-sm
          ${n.is_entry ? "bg-blue-600 text-white border-blue-600" : "bg-white text-[var(--text)] border-[var(--border)]"}`}>
          {n.is_entry && "▶ "}
          {n.title || <em>{n.id}</em>}
          {n.image_url && " 🖼️"}
        </div>
        {n.choices.length > 0 && (
          <div className="ml-6 mt-2 space-y-2 border-l-2 border-dashed border-gray-200 pl-4">
            {n.choices.map((ch, i) => (
              <div key={i} className="flex items-start gap-2">
                <span className={`mt-1 shrink-0 text-[10px] font-bold px-1.5 py-0.5 rounded border
                  ${VARIANT_CFG[ch.variant].cls}`}>
                  {VARIANT_CFG[ch.variant].label}
                </span>
                <div className="min-w-0">
                  <span className="text-xs text-[var(--text-muted)]">{ch.label} → </span>
                  {ch.next_id
                    ? <Tree nodeId={ch.next_id} visited={next} />
                    : <span className="text-xs italic text-gray-400">selesai</span>
                  }
                </div>
              </div>
            ))}
          </div>
        )}
        {n.choices.length === 0 && (
          <div className="ml-6 mt-1 pl-4 border-l-2 border-dashed border-gray-200">
            <span className="text-xs italic text-gray-400">selesai</span>
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="p-5 bg-[var(--bg)] rounded-2xl border border-[var(--border)] overflow-auto">
      <Tree nodeId={entry.id} visited={new Set()} />
    </div>
  );
}

// ── FlowEditor ────────────────────────────────────────────────────────────────

function FlowEditor({
  flow, onSave, onDelete, saving,
}: {
  flow: EmergencyGuideFlow | "new";
  onSave: (d: { title_id: string; title_en: string; kategori: string; deskripsi: string; nodes_id: GuideNode[]; nodes_en: GuideNode[] }) => Promise<void>;
  onDelete: () => Promise<void>;
  saving: boolean;
}) {
  const isNew = flow === "new";
  const cloneNodes = (ns: GuideNode[]) => (ns ?? []).map(n => ({ ...n, choices: [...n.choices] }));
  const init = isNew
    ? { title_id: "", title_en: "", kategori: "BLS", deskripsi: "", nodes_id: [makeNode(true)], nodes_en: [makeNode(true)] }
    : { title_id: flow.title_id, title_en: flow.title_en, kategori: flow.kategori, deskripsi: flow.deskripsi ?? "",
        nodes_id: cloneNodes(flow.nodes_id), nodes_en: cloneNodes(flow.nodes_en) };

  const [titleID, setTitleID]       = useState(init.title_id);
  const [titleEN, setTitleEN]       = useState(init.title_en);
  const [kategori, setKategori]     = useState(init.kategori);
  const [deskripsi, setDeskripsi]   = useState(init.deskripsi);
  const [nodesID, setNodesID]       = useState<GuideNode[]>(init.nodes_id);
  const [nodesEN, setNodesEN]       = useState<GuideNode[]>(init.nodes_en);
  const [activeLang, setActiveLang] = useState<"id" | "en">("id");
  const [activeIdx, setActiveIdx]   = useState(0);
  const [activeTab, setActiveTab]   = useState<"edit" | "map">("edit");
  const [error, setError]           = useState("");
  const [confirmDel, setConfirmDel] = useState(false);

  const flowKey = isNew ? "new" : flow.id;
  const prevKey = useRef(flowKey);
  if (prevKey.current !== flowKey) {
    prevKey.current = flowKey;
    setTitleID(init.title_id);
    setTitleEN(init.title_en);
    setKategori(init.kategori);
    setDeskripsi(init.deskripsi);
    setNodesID(cloneNodes(init.nodes_id));
    setNodesEN(cloneNodes(init.nodes_en));
    setActiveLang("id");
    setActiveIdx(0);
    setActiveTab("edit");
    setError("");
    setConfirmDel(false);
  }

  const nodes = activeLang === "id" ? nodesID : nodesEN;
  const setNodes = activeLang === "id" ? setNodesID : setNodesEN;
  const activeNode = nodes[activeIdx] ?? nodes[0] ?? null;

  function addNode() {
    const n = makeNode(false);
    setNodes(prev => [...prev, n]);
    setActiveIdx(nodes.length);
  }

  function updateNode(updated: GuideNode) {
    if (updated.is_entry) {
      setNodes(prev => prev.map((n, i) => i === activeIdx ? updated : { ...n, is_entry: false }));
    } else {
      setNodes(prev => prev.map((n, i) => i === activeIdx ? updated : n));
    }
  }

  function deleteNode(i: number) {
    const delId = nodes[i].id;
    const next = nodes
      .filter((_, idx) => idx !== i)
      .map(n => ({ ...n, choices: n.choices.map(c => c.next_id === delId ? { ...c, next_id: null } : c) }));
    setNodes(next);
    setActiveIdx(Math.min(i, next.length - 1));
  }

  function moveNode(i: number, dir: -1 | 1) {
    const j = i + dir;
    if (j < 0 || j >= nodes.length) return;
    const copy = [...nodes];
    [copy[i], copy[j]] = [copy[j], copy[i]];
    setNodes(copy);
    setActiveIdx(j);
  }

  function validate(): string | null {
    if (!titleID.trim()) return "Judul alur (Indonesia) harus diisi";
    if (!titleEN.trim()) return "Judul alur (English) harus diisi";
    if (nodes.length === 0) return "Minimal 1 node diperlukan";
    const entryCount = nodes.filter(n => n.is_entry).length;
    if (entryCount !== 1) return "Tepat 1 node harus dijadikan Entry Point (START)";
    const seen = new Set<string>();
    const allIds = new Set(nodes.map(n => n.id));
    for (const n of nodes) {
      if (!n.id.trim()) return `Ada node dengan ID kosong`;
      if (!n.title.trim()) return `Node "${n.id}" belum diberi judul`;
      if (seen.has(n.id)) return `ID node "${n.id}" dipakai lebih dari satu kali`;
      seen.add(n.id);
      for (const ch of n.choices) {
        if (!ch.label.trim()) return `Node "${n.id}": teks pilihan tidak boleh kosong`;
        if (ch.next_id && !allIds.has(ch.next_id)) {
          return `Node "${n.id}": pilihan mengarah ke node "${ch.next_id}" yang tidak ada`;
        }
      }
    }
    return null;
  }

  async function handleSave() {
    const err = validate();
    if (err) { setError(err); return; }
    setError("");
    try {
      await onSave({ title_id: titleID, title_en: titleEN, kategori, deskripsi, nodes_id: nodesID, nodes_en: nodesEN });
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Gagal menyimpan");
    }
  }

  return (
    <div className="flex-1 flex flex-col min-w-0 overflow-hidden">

      {/* ── top bar ── */}
      <div className="flex items-start gap-3 px-5 py-3 border-b border-[var(--border)] bg-[var(--surface)]">
        <div className="flex-1 min-w-0 space-y-1">
          <div className="grid grid-cols-2 gap-2">
            <input
              value={titleID}
              onChange={e => setTitleID(e.target.value)}
              placeholder="Judul (Indonesia)…"
              className="text-sm font-semibold bg-transparent border-b border-transparent hover:border-[var(--border)] focus:border-[var(--brand)] focus:outline-none pb-0.5"
            />
            <input
              value={titleEN}
              onChange={e => setTitleEN(e.target.value)}
              placeholder="Title (English)…"
              className="text-sm font-semibold bg-transparent border-b border-transparent hover:border-[var(--border)] focus:border-[var(--brand)] focus:outline-none pb-0.5"
            />
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            <select
              value={kategori}
              onChange={e => setKategori(e.target.value)}
              className="text-xs border border-[var(--border)] rounded-lg px-2 py-1 focus:outline-none focus:ring-1 focus:ring-[var(--brand)]"
            >
              {KATEGORI.map(k => <option key={k}>{k}</option>)}
            </select>
            <input
              value={deskripsi}
              onChange={e => setDeskripsi(e.target.value)}
              placeholder="Deskripsi singkat (opsional)"
              className="flex-1 min-w-[160px] text-xs text-[var(--text-muted)] bg-transparent border-b border-transparent hover:border-[var(--border)] focus:border-[var(--brand)] focus:outline-none"
            />
          </div>
        </div>
        <div className="flex items-center gap-2 shrink-0">
          {!isNew && (
            confirmDel ? (
              <>
                <span className="text-xs text-red-600 font-medium">Hapus alur ini?</span>
                <button onClick={onDelete} disabled={saving} className="px-3 py-1.5 text-xs bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50">Hapus</button>
                <button onClick={() => setConfirmDel(false)} className="px-3 py-1.5 text-xs border border-[var(--border)] rounded-lg">Batal</button>
              </>
            ) : (
              <button onClick={() => setConfirmDel(true)} className="px-3 py-1.5 text-xs border border-red-200 text-red-500 rounded-lg hover:bg-red-50">Hapus</button>
            )
          )}
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-4 py-1.5 text-sm font-medium bg-[var(--brand)] text-white rounded-lg hover:opacity-90 disabled:opacity-50"
          >
            {saving ? "Menyimpan…" : isNew ? "Buat Alur" : "Simpan Perubahan"}
          </button>
        </div>
      </div>

      {error && (
        <div className="mx-5 mt-3 px-4 py-2.5 bg-red-50 border border-red-200 text-red-700 text-sm rounded-xl flex items-start gap-2">
          <span>⚠</span> <span>{error}</span>
        </div>
      )}

      {/* ── tabs ── */}
      <div className="flex gap-5 px-5 border-b border-[var(--border)]">
        {(["edit", "map"] as const).map(t => (
          <button
            key={t}
            onClick={() => setActiveTab(t)}
            className={`py-2.5 text-sm border-b-2 transition-colors ${
              activeTab === t
                ? "border-[var(--brand)] text-[var(--brand)] font-semibold"
                : "border-transparent text-[var(--text-muted)] hover:text-[var(--text)]"
            }`}
          >
            {t === "edit" ? `✏️ Editor Node (${nodes.length})` : "🗺️ Peta Alur"}
          </button>
        ))}
      </div>

      {/* ── language switcher ── */}
      <div className="flex gap-1 px-5 py-2 border-b border-[var(--border)] bg-[var(--bg)]">
        <span className="text-xs text-[var(--text-muted)] mr-2 self-center">Bahasa node:</span>
        {(["id", "en"] as const).map(lang => (
          <button
            key={lang}
            onClick={() => { setActiveLang(lang); setActiveIdx(0); }}
            className={`px-3 py-1 text-xs font-medium rounded-lg transition-colors border
              ${activeLang === lang
                ? "bg-[var(--brand)] text-white border-[var(--brand)]"
                : "border-[var(--border)] text-[var(--text-muted)] hover:text-[var(--text)]"}`}
          >
            {lang === "id" ? "🇮🇩 Indonesia" : "🇬🇧 English"}
          </button>
        ))}
      </div>

      {activeTab === "map" && (
        <div className="flex-1 overflow-y-auto p-5">
          <FlowMap nodes={nodes} />
        </div>
      )}

      {activeTab === "edit" && (
        <div className="flex-1 flex overflow-hidden min-h-0">

          {/* ── node list sidebar ── */}
          <div className="w-56 shrink-0 border-r border-[var(--border)] flex flex-col bg-[var(--surface)]">
            <div className="px-3 py-2.5 border-b border-[var(--border)]">
              <p className="text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wide">Node</p>
            </div>
            <div className="flex-1 overflow-y-auto p-2 space-y-1.5">
              {nodes.map((n, i) => (
                <NodeListItem
                  key={`${n.id}-${i}`}
                  node={n}
                  isActive={i === activeIdx}
                  index={i}
                  total={nodes.length}
                  onClick={() => setActiveIdx(i)}
                  onMoveUp={() => moveNode(i, -1)}
                  onMoveDown={() => moveNode(i, 1)}
                  onDelete={() => deleteNode(i)}
                />
              ))}
            </div>
            <div className="p-2 border-t border-[var(--border)]">
              <button
                onClick={addNode}
                className="w-full py-2 text-xs font-medium text-[var(--brand)] border border-dashed border-[var(--brand)] rounded-xl hover:bg-blue-50 transition-colors"
              >
                + Tambah Node
              </button>
            </div>
          </div>

          {/* ── node detail ── */}
          <div className="flex-1 overflow-hidden">
            {activeNode ? (
              <NodeEditor
                node={activeNode}
                allNodes={nodes}
                onChange={updateNode}
              />
            ) : (
              <div className="flex items-center justify-center h-full text-[var(--text-muted)] text-sm">
                Pilih node dari daftar kiri
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

// ── FlowList sidebar ──────────────────────────────────────────────────────────

function FlowList({
  flows, selected, onSelect, onNew,
}: {
  flows: EmergencyGuideFlow[];
  selected: number | null;
  onSelect: (id: number) => void;
  onNew: () => void;
}) {
  return (
    <aside className="w-56 shrink-0 flex flex-col border-r border-[var(--border)] bg-[var(--surface)] rounded-l-2xl overflow-hidden">
      <div className="flex items-center justify-between px-4 py-3 border-b border-[var(--border)]">
        <span className="font-semibold text-[var(--text)] text-sm">Alur Panduan</span>
        <button
          onClick={onNew}
          title="Buat alur baru"
          className="w-7 h-7 flex items-center justify-center rounded-lg bg-[var(--brand)] text-white hover:opacity-90 text-xl leading-none"
        >+</button>
      </div>
      <ul className="flex-1 overflow-y-auto py-1.5 px-1.5 space-y-1">
        {flows.length === 0 && (
          <li className="py-8 text-center text-xs text-[var(--text-muted)]">Belum ada alur panduan</li>
        )}
        {flows.map(f => (
          <li key={f.id}>
            <button
              onClick={() => onSelect(f.id)}
              className={`w-full text-left px-3 py-2.5 rounded-xl flex flex-col gap-0.5 transition-colors
                ${selected === f.id
                  ? "bg-blue-50 border border-[var(--brand)]"
                  : "hover:bg-[var(--bg)] border border-transparent"
                }`}
            >
              <span className="text-sm font-medium text-[var(--text)] leading-tight truncate">{f.title_id}</span>
              <span className="text-[11px] text-[var(--text-muted)]">{f.kategori} · {f.nodes_id?.length ?? 0} node</span>
            </button>
          </li>
        ))}
      </ul>
    </aside>
  );
}

// ── Page root ─────────────────────────────────────────────────────────────────

export default function EmergencyGuidesPage() {
  const [flows, setFlows]       = useState<EmergencyGuideFlow[]>([]);
  const [selected, setSelected] = useState<number | "new" | null>(null);
  const [loading, setLoading]   = useState(true);
  const [saving, setSaving]     = useState(false);
  const [toast, setToast]       = useState<{ msg: string; ok: boolean } | null>(null);

  const showToast = useCallback((msg: string, ok = true) => {
    setToast({ msg, ok });
    setTimeout(() => setToast(null), 3000);
  }, []);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await emergencyFlowsApi.list();
      setFlows(res.data ?? []);
    } finally { setLoading(false); }
  }, []);

  useEffect(() => { load(); }, [load]);

  const selectedFlow =
    selected === "new" ? "new"
    : selected !== null ? (flows.find(f => f.id === selected) ?? null)
    : null;

  async function handleSave(data: { title_id: string; title_en: string; kategori: string; deskripsi: string; nodes_id: GuideNode[]; nodes_en: GuideNode[] }) {
    setSaving(true);
    try {
      if (selected === "new") {
        const created = await emergencyFlowsApi.create(data);
        await load();
        setSelected(created.id);
        showToast("Alur berhasil dibuat ✓");
      } else if (typeof selected === "number") {
        await emergencyFlowsApi.update(selected, data);
        await load();
        showToast("Perubahan berhasil disimpan ✓");
      }
    } catch (e: unknown) {
      throw e;
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete() {
    if (typeof selected !== "number") return;
    setSaving(true);
    try {
      await emergencyFlowsApi.del(selected);
      setSelected(null);
      await load();
      showToast("Alur berhasil dihapus");
    } catch {
      showToast("Gagal menghapus alur", false);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="h-full flex flex-col">
      <div className="px-6 pt-6 pb-4">
        <h1 className="text-2xl font-bold text-[var(--text)]">Panduan Darurat</h1>
        <p className="text-sm text-[var(--text-muted)] mt-0.5">
          Editor alur keputusan Yes/No — format MySOS · Gambar bisa diunggah langsung
        </p>
      </div>

      <div className="flex-1 flex mx-6 mb-6 border border-[var(--border)] rounded-2xl overflow-hidden shadow-sm min-h-0">
        {loading ? (
          <div className="flex-1 flex items-center justify-center text-[var(--text-muted)]">
            <svg className="animate-spin w-6 h-6 mr-2 text-[var(--brand)]" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" strokeDasharray="30" strokeDashoffset="10" />
            </svg>
            Memuat panduan…
          </div>
        ) : (
          <>
            <FlowList
              flows={flows}
              selected={typeof selected === "number" ? selected : null}
              onSelect={id => setSelected(id)}
              onNew={() => setSelected("new")}
            />

            {selectedFlow !== null ? (
              <FlowEditor
                key={selected}
                flow={selectedFlow}
                onSave={handleSave}
                onDelete={handleDelete}
                saving={saving}
              />
            ) : (
              <div className="flex-1 flex flex-col items-center justify-center gap-4">
                <div className="w-20 h-20 rounded-3xl bg-blue-50 flex items-center justify-center text-4xl shadow-inner">🚑</div>
                <div className="text-center">
                  <p className="font-semibold text-[var(--text)]">Pilih alur panduan untuk diedit</p>
                  <p className="text-sm text-[var(--text-muted)] mt-1">atau klik <strong>+</strong> untuk membuat alur baru</p>
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {toast && (
        <div className={`fixed bottom-6 right-6 px-4 py-3 rounded-xl shadow-xl text-sm font-medium z-50 animate-slideUp
          ${toast.ok ? "bg-green-600 text-white" : "bg-red-600 text-white"}`}>
          {toast.msg}
        </div>
      )}
    </div>
  );
}
