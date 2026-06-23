export default function NotFound() {
  return (
    <div style={{ display: "flex", alignItems: "center", justifyContent: "center", height: "100vh", flexDirection: "column", gap: "1rem", fontFamily: "sans-serif" }}>
      <h1 style={{ fontSize: "2rem", fontWeight: "bold" }}>404</h1>
      <p style={{ color: "#64748b" }}>Halaman tidak ditemukan</p>
      <a href="/dashboard" style={{ color: "#2563eb", textDecoration: "underline" }}>Kembali ke Dashboard</a>
    </div>
  );
}
