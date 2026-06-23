"use client";

import { useEffect, useState, useCallback } from "react";
import { assessmentsApi } from "@/lib/api";
import type { Assessment } from "@/lib/types";
import DataTable from "@/components/DataTable";
import PageHeader from "@/components/PageHeader";

const RISK_COLORS: Record<string, { color: string; bg: string }> = {
  Darurat: { color: "var(--danger)",  bg: "#fef2f2" },
  Tinggi:  { color: "#ea580c",        bg: "#fff7ed" },
  Sedang:  { color: "var(--warning)", bg: "#fffbeb" },
  Rendah:  { color: "var(--success)", bg: "#f0fdf4" },
};

function RiskBadge({ level }: { level?: string }) {
  if (!level) return <span style={{ color: "var(--text-muted)" }}>—</span>;
  const s = RISK_COLORS[level] ?? { color: "var(--text-muted)", bg: "var(--bg)" };
  return (
    <span
      className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold"
      style={{ color: s.color, background: s.bg }}
    >
      {level}
    </span>
  );
}

export default function AssessmentsPage() {
  const [rows, setRows]       = useState<Assessment[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage]       = useState(1);
  const [total, setTotal]     = useState(0);
  const limit = 20;

  const load = useCallback(async (p: number) => {
    setLoading(true);
    try {
      const res = await assessmentsApi.list(p, limit);
      setRows(res.data ?? []);
      setTotal(res.total ?? 0);
    } finally { setLoading(false); }
  }, []);

  useEffect(() => { load(page); }, [load, page]);

  const totalPages = Math.ceil(total / limit);

  const columns = [
    { key: "id",      label: "ID",       render: (r: Assessment) => <span className="font-mono text-xs">{r.id}</span> },
    { key: "user_id", label: "User ID",  render: (r: Assessment) => <span className="font-mono text-xs">{r.user_id}</span> },
    { key: "diagnosis", label: "Diagnosis", render: (r: Assessment) => r.diagnosis ?? "—" },
    { key: "confidence_score", label: "CF",
      render: (r: Assessment) =>
        r.confidence_score != null ? (
          <span className="font-mono text-xs">{r.confidence_score.toFixed(3)}</span>
        ) : "—",
    },
    { key: "risk_level", label: "Risiko", render: (r: Assessment) => <RiskBadge level={r.risk_level} /> },
    {
      key: "created_at", label: "Tanggal",
      render: (r: Assessment) =>
        new Date(r.created_at).toLocaleDateString("id-ID", {
          day: "2-digit", month: "short", year: "numeric",
        }),
    },
  ];

  return (
    <div>
      <PageHeader
        title="Data Asesmen"
        description={`Total ${total} asesmen (read-only)`}
      />

      <div className="reveal">
        <DataTable columns={columns} rows={rows} keyField="id" loading={loading} />
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-end gap-2 mt-4">
          <button
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
            className="px-3 py-1.5 rounded-lg text-sm font-medium disabled:opacity-40 transition-colors"
            style={{ border: "1px solid var(--border)", color: "var(--text-muted)" }}
          >
            ← Sebelumnya
          </button>
          <span className="text-sm" style={{ color: "var(--text-muted)" }}>
            {page} / {totalPages}
          </span>
          <button
            disabled={page >= totalPages}
            onClick={() => setPage((p) => p + 1)}
            className="px-3 py-1.5 rounded-lg text-sm font-medium disabled:opacity-40 transition-colors"
            style={{ border: "1px solid var(--border)", color: "var(--text-muted)" }}
          >
            Berikutnya →
          </button>
        </div>
      )}
    </div>
  );
}
