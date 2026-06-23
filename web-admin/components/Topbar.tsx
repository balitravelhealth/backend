"use client";

import { usePathname } from "next/navigation";

const pageTitles: Record<string, string> = {
  "/dashboard": "Dashboard",
  "/dashboard/facilities": "Fasilitas Kesehatan",
  "/dashboard/destinations": "Destinasi & Risiko Kesehatan",
  "/dashboard/emergency-guides": "Panduan Darurat",
  "/dashboard/nurses": "Kelola Perawat",
  "/dashboard/assessments": "Data Asesmen",
  "/dashboard/knowledge-base": "Knowledge Base Sistem Pakar",
};

export default function Topbar() {
  const pathname = usePathname();

  const title =
    Object.entries(pageTitles)
      .filter(([key]) => pathname.startsWith(key))
      .sort((a, b) => b[0].length - a[0].length)[0]?.[1] ?? "Admin";

  return (
    <header
      className="flex items-center justify-between px-6 py-3 border-b"
      style={{
        background: "var(--surface)",
        borderColor: "var(--border)",
        minHeight: "56px",
      }}
    >
      <h1 className="text-base font-semibold" style={{ color: "var(--text)" }}>
        {title}
      </h1>

      <div className="flex items-center gap-2">
        <div
          className="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold"
          style={{ background: "var(--brand-grad)" }}
        >
          A
        </div>
        <span className="text-sm font-medium hidden sm:block" style={{ color: "var(--text-muted)" }}>
          Admin
        </span>
      </div>
    </header>
  );
}
