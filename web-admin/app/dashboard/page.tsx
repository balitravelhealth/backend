"use client";

import { useEffect, useState, useCallback } from "react";
import type { ReactNode } from "react";
import {
  PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, Tooltip,
  ResponsiveContainer, CartesianGrid, Label,
} from "recharts";
import {
  facilitiesApi, nursesApi, assessmentsApi, rulesApi, symptomsApi, diseasesApi,
} from "@/lib/api";
import type { Assessment } from "@/lib/types";

// ── colour palette ────────────────────────────────────────────────────────────

const RISK_COLORS: Record<string, string> = {
  Darurat: "#ef4444",
  Tinggi:  "#f97316",
  Sedang:  "#eab308",
  Rendah:  "#22c55e",
};

const CAT_COLOR = "#6366f1";

// ── tiny helpers ──────────────────────────────────────────────────────────────

function shortCat(cat: string): string {
  return cat
    .replace("Rumah Sakit Umum Daerah", "RSUD")
    .replace("Rumah Sakit Umum Pemerintah", "RSUP")
    .replace("Rumah Sakit Swasta Internasional", "RS Intl")
    .replace("Rumah Sakit Swasta", "RS Swasta")
    .replace("Rumah Sakit Khusus", "RS Khusus")
    .replace("Klinik Internasional", "Klinik Intl")
    .replace("Klinik Swasta", "Klinik");
}

// ── StatCard ──────────────────────────────────────────────────────────────────

function StatCard({
  label, value, sub, icon, color, loading,
}: {
  label: string;
  value: string | number;
  sub?: string;
  icon: ReactNode;
  color: string;
  loading: boolean;
}) {
  return (
    <div
      className="reveal p-5 flex items-center gap-4"
      style={{
        background: "var(--surface)",
        borderRadius: "var(--radius)",
        boxShadow: "var(--shadow-card)",
        border: "1px solid var(--border)",
      }}
    >
      <div
        className="w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0"
        style={{ background: `${color}18`, color }}
      >
        {icon}
      </div>
      <div className="min-w-0">
        {loading ? (
          <div className="h-8 w-16 rounded-lg animate-pulse" style={{ background: "var(--border)" }} />
        ) : (
          <p className="text-3xl font-bold leading-none" style={{ color: "var(--text)" }}>
            {value}
          </p>
        )}
        <p className="text-sm mt-1 truncate" style={{ color: "var(--text-muted)" }}>{label}</p>
        {sub && <p className="text-xs mt-0.5" style={{ color: "var(--text-muted)" }}>{sub}</p>}
      </div>
    </div>
  );
}

// ── ChartCard ─────────────────────────────────────────────────────────────────

function ChartCard({
  title, subtitle, children,
}: {
  title: string;
  subtitle?: string;
  children: ReactNode;
}) {
  return (
    <div
      className="p-5"
      style={{
        background: "var(--surface)",
        borderRadius: "var(--radius)",
        boxShadow: "var(--shadow-card)",
        border: "1px solid var(--border)",
      }}
    >
      <div className="mb-4">
        <h3 className="font-semibold text-sm" style={{ color: "var(--text)" }}>{title}</h3>
        {subtitle && <p className="text-xs mt-0.5" style={{ color: "var(--text-muted)" }}>{subtitle}</p>}
      </div>
      {children}
    </div>
  );
}

// ── custom donut label ────────────────────────────────────────────────────────

function DonutLabel({ viewBox, total }: { viewBox?: { cx: number; cy: number }; total: number }) {
  const cx = viewBox?.cx ?? 0;
  const cy = viewBox?.cy ?? 0;
  return (
    <text x={cx} y={cy} textAnchor="middle" dominantBaseline="middle">
      <tspan x={cx} dy="-6" fontSize={22} fontWeight={700} fill="var(--text)">{total}</tspan>
      <tspan x={cx} dy={20} fontSize={11} fill="var(--text-muted)">asesmen</tspan>
    </text>
  );
}

// ── custom tooltip ────────────────────────────────────────────────────────────

function CustomTooltip({ active, payload, label }: {
  active?: boolean; payload?: { name: string; value: number; fill: string }[]; label?: string;
}) {
  if (!active || !payload?.length) return null;
  return (
    <div
      className="px-3 py-2 rounded-lg text-xs shadow-lg"
      style={{ background: "var(--surface)", border: "1px solid var(--border)" }}
    >
      {label && <p className="font-semibold mb-1" style={{ color: "var(--text)" }}>{label}</p>}
      {payload.map((p, i) => (
        <p key={i} style={{ color: p.fill ?? "var(--text)" }}>
          {p.name}: <span className="font-semibold">{p.value}</span>
        </p>
      ))}
    </div>
  );
}

// ── Page ──────────────────────────────────────────────────────────────────────

interface Stats {
  facilities: number;
  activeNurses: number;
  totalNurses: number;
  totalAssessments: number;
  publishedRules: number;
  totalRules: number;
  totalSymptoms: number;
  totalDiseases: number;
}

export default function DashboardPage() {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<Stats>({
    facilities: 0, activeNurses: 0, totalNurses: 0,
    totalAssessments: 0, publishedRules: 0, totalRules: 0,
    totalSymptoms: 0, totalDiseases: 0,
  });
  const [riskData, setRiskData]         = useState<{ name: string; value: number; color: string }[]>([]);
  const [facByCat, setFacByCat]         = useState<{ name: string; count: number }[]>([]);
  const [topDiag, setTopDiag]           = useState<{ name: string; count: number }[]>([]);
  const [recentAss, setRecentAss]       = useState<Assessment[]>([]);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [facRes, nurRes, assRes, rulRes, symRes, disRes] = await Promise.all([
        facilitiesApi.list(),
        nursesApi.list(),
        assessmentsApi.list(1, 50),
        rulesApi.list(),
        symptomsApi.list(),
        diseasesApi.list(),
      ]);

      const facs    = facRes.data  ?? [];
      const nurses  = nurRes.data  ?? [];
      const asses   = assRes.data  ?? [];
      const rules   = rulRes.data  ?? [];
      const syms    = symRes.data  ?? [];
      const dises   = disRes.data  ?? [];

      // ── stats ──
      setStats({
        facilities:       facs.length,
        activeNurses:     nurses.filter((n) => n.aktif).length,
        totalNurses:      nurses.length,
        totalAssessments: assRes.total ?? asses.length,
        publishedRules:   rules.filter((r) => r.status === "published").length,
        totalRules:       rules.length,
        totalSymptoms:    syms.length,
        totalDiseases:    dises.length,
      });

      // ── risk distribution ──
      const riskMap: Record<string, number> = {};
      asses.forEach((a) => {
        const lvl = a.risk_level || "Tidak diketahui";
        riskMap[lvl] = (riskMap[lvl] ?? 0) + 1;
      });
      const riskOrder = ["Darurat", "Tinggi", "Sedang", "Rendah", "Tidak diketahui"];
      setRiskData(
        riskOrder
          .filter((k) => riskMap[k])
          .map((k) => ({
            name: k,
            value: riskMap[k],
            color: RISK_COLORS[k] ?? "#94a3b8",
          })),
      );

      // ── facilities by category ──
      const catMap: Record<string, number> = {};
      facs.forEach((f) => {
        const cat = shortCat(f.kategori ?? "Lainnya");
        catMap[cat] = (catMap[cat] ?? 0) + 1;
      });
      setFacByCat(
        Object.entries(catMap)
          .sort((a, b) => b[1] - a[1])
          .map(([name, count]) => ({ name, count })),
      );

      // ── top 5 diagnoses ──
      const diagMap: Record<string, number> = {};
      asses.forEach((a) => {
        if (a.diagnosis) diagMap[a.diagnosis] = (diagMap[a.diagnosis] ?? 0) + 1;
      });
      setTopDiag(
        Object.entries(diagMap)
          .sort((a, b) => b[1] - a[1])
          .slice(0, 5)
          .map(([name, count]) => ({ name: name.length > 28 ? name.slice(0, 26) + "…" : name, count })),
      );

      // ── recent 5 assessments ──
      setRecentAss(asses.slice(0, 5));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  const kbData = [
    { name: "Gejala", count: stats.totalSymptoms,   fill: "#6366f1" },
    { name: "Penyakit", count: stats.totalDiseases, fill: "#8b5cf6" },
    { name: "Rules", count: stats.totalRules,       fill: "#06b6d4" },
  ];

  return (
    <div className="space-y-6">

      {/* ── welcome ── */}
      <div>
        <h2 className="text-xl font-semibold" style={{ color: "var(--text)" }}>
          Selamat datang
        </h2>
        <p className="text-sm mt-0.5" style={{ color: "var(--text-muted)" }}>
          Ringkasan data BaliTravelHealth — diperbarui setiap kali halaman dibuka
        </p>
      </div>

      {/* ── stat cards ── */}
      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
        <StatCard loading={loading} label="Fasilitas Kesehatan" value={stats.facilities}
          sub="klinik & rumah sakit terdaftar"
          color="var(--brand)"
          icon={
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          }
        />
        <StatCard loading={loading} label="Perawat Aktif" value={stats.activeNurses}
          sub={`dari ${stats.totalNurses} perawat terdaftar`}
          color="var(--success)"
          icon={
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          }
        />
        <StatCard loading={loading} label="Total Asesmen" value={stats.totalAssessments}
          sub="diagnosa wisatawan tersimpan"
          color="var(--warning)"
          icon={
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          }
        />
        <StatCard loading={loading} label="Rules Aktif" value={stats.publishedRules}
          sub={`dari ${stats.totalRules} rule — ${stats.totalDiseases} penyakit, ${stats.totalSymptoms} gejala`}
          color="#8b5cf6"
          icon={
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
            </svg>
          }
        />
      </div>

      {/* ── charts row 1 ── */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">

        {/* Donut: risk distribution */}
        <ChartCard title="Distribusi Tingkat Risiko" subtitle="Berdasarkan hasil asesmen wisatawan">
          {riskData.length === 0 && !loading ? (
            <div className="flex items-center justify-center h-40 text-sm" style={{ color: "var(--text-muted)" }}>
              Belum ada data asesmen
            </div>
          ) : (
            <>
              <ResponsiveContainer width="100%" height={180}>
                <PieChart>
                  <Pie
                    data={riskData}
                    cx="50%" cy="50%"
                    innerRadius={55} outerRadius={80}
                    paddingAngle={3}
                    dataKey="value"
                  >
                    {riskData.map((entry, i) => (
                      <Cell key={i} fill={entry.color} strokeWidth={0} />
                    ))}
                    <Label
                      content={(props) => <DonutLabel viewBox={props.viewBox as { cx: number; cy: number }} total={riskData.reduce((s, d) => s + d.value, 0)} />}
                      position="center"
                    />
                  </Pie>
                  <Tooltip content={<CustomTooltip />} />
                </PieChart>
              </ResponsiveContainer>
              <div className="flex flex-wrap gap-x-4 gap-y-1 mt-1">
                {riskData.map((d) => (
                  <div key={d.name} className="flex items-center gap-1.5 text-xs" style={{ color: "var(--text-muted)" }}>
                    <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: d.color }} />
                    {d.name} ({d.value})
                  </div>
                ))}
              </div>
            </>
          )}
        </ChartCard>

        {/* Bar: facilities by category */}
        <ChartCard title="Fasilitas per Kategori" subtitle="Distribusi jenis fasilitas kesehatan">
          {facByCat.length === 0 && !loading ? (
            <div className="flex items-center justify-center h-40 text-sm" style={{ color: "var(--text-muted)" }}>
              Belum ada fasilitas
            </div>
          ) : (
            <ResponsiveContainer width="100%" height={200}>
              <BarChart data={facByCat} layout="vertical" margin={{ left: 8, right: 16 }}>
                <CartesianGrid strokeDasharray="3 3" horizontal={false} stroke="var(--border)" />
                <XAxis type="number" tick={{ fontSize: 11, fill: "var(--text-muted)" }} allowDecimals={false} />
                <YAxis type="category" dataKey="name" width={80} tick={{ fontSize: 10, fill: "var(--text-muted)" }} />
                <Tooltip content={<CustomTooltip />} />
                <Bar dataKey="count" name="Fasilitas" fill={CAT_COLOR} radius={[0, 4, 4, 0]} barSize={14} />
              </BarChart>
            </ResponsiveContainer>
          )}
        </ChartCard>

        {/* Bar: knowledge base */}
        <ChartCard title="Knowledge Base" subtitle="Jumlah komponen sistem pakar">
          <ResponsiveContainer width="100%" height={200}>
            <BarChart data={kbData} margin={{ left: 0, right: 8 }}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="var(--border)" />
              <XAxis dataKey="name" tick={{ fontSize: 12, fill: "var(--text-muted)" }} />
              <YAxis tick={{ fontSize: 11, fill: "var(--text-muted)" }} allowDecimals={false} />
              <Tooltip content={<CustomTooltip />} />
              <Bar dataKey="count" name="Jumlah" radius={[4, 4, 0, 0]} barSize={40}>
                {kbData.map((entry, i) => (
                  <Cell key={i} fill={entry.fill} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>

          {/* mini stats row */}
          <div className="grid grid-cols-3 gap-2 mt-3 pt-3 border-t" style={{ borderColor: "var(--border)" }}>
            {kbData.map((d) => (
              <div key={d.name} className="text-center">
                <p className="text-xl font-bold" style={{ color: d.fill }}>{loading ? "—" : d.count}</p>
                <p className="text-xs" style={{ color: "var(--text-muted)" }}>{d.name}</p>
              </div>
            ))}
          </div>
        </ChartCard>
      </div>

      {/* ── charts row 2 ── */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">

        {/* Top diagnoses */}
        <ChartCard title="Top 5 Diagnosa Terbanyak" subtitle="Berdasarkan data asesmen wisatawan">
          {topDiag.length === 0 && !loading ? (
            <div className="flex items-center justify-center h-36 text-sm" style={{ color: "var(--text-muted)" }}>
              Belum ada data diagnosa
            </div>
          ) : (
            <ResponsiveContainer width="100%" height={180}>
              <BarChart data={topDiag} layout="vertical" margin={{ left: 8, right: 16 }}>
                <CartesianGrid strokeDasharray="3 3" horizontal={false} stroke="var(--border)" />
                <XAxis type="number" allowDecimals={false} tick={{ fontSize: 11, fill: "var(--text-muted)" }} />
                <YAxis type="category" dataKey="name" width={130} tick={{ fontSize: 10, fill: "var(--text-muted)" }} />
                <Tooltip content={<CustomTooltip />} />
                <Bar dataKey="count" name="Asesmen" fill="#06b6d4" radius={[0, 4, 4, 0]} barSize={16} />
              </BarChart>
            </ResponsiveContainer>
          )}
        </ChartCard>

        {/* Recent assessments */}
        <ChartCard title="Asesmen Terbaru" subtitle="5 asesmen wisatawan paling baru">
          {recentAss.length === 0 && !loading ? (
            <div className="flex items-center justify-center h-36 text-sm" style={{ color: "var(--text-muted)" }}>
              Belum ada asesmen
            </div>
          ) : (
            <div className="divide-y" style={{ borderColor: "var(--border)" }}>
              {recentAss.map((a) => {
                const riskColor = RISK_COLORS[a.risk_level ?? ""] ?? "#94a3b8";
                return (
                  <div key={a.id} className="flex items-center gap-3 py-2.5">
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate" style={{ color: "var(--text)" }}>
                        {a.diagnosis ?? "Tidak terdiagnosis"}
                      </p>
                      <p className="text-xs" style={{ color: "var(--text-muted)" }}>
                        User #{a.user_id} · {new Date(a.created_at).toLocaleDateString("id-ID", {
                          day: "2-digit", month: "short", year: "numeric",
                        })}
                      </p>
                    </div>
                    <div className="flex items-center gap-2 flex-shrink-0">
                      {a.confidence_score != null && (
                        <span className="font-mono text-xs" style={{ color: "var(--text-muted)" }}>
                          CF {a.confidence_score.toFixed(2)}
                        </span>
                      )}
                      {a.risk_level && (
                        <span
                          className="text-xs font-semibold px-2 py-0.5 rounded-full"
                          style={{ color: riskColor, background: `${riskColor}18` }}
                        >
                          {a.risk_level}
                        </span>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </ChartCard>
      </div>
    </div>
  );
}
