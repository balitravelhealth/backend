"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import { adminLogin, saveToken, ApiError } from "@/lib/api";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail]       = useState("");
  const [password, setPassword] = useState("");
  const [error, setError]       = useState("");
  const [loading, setLoading]   = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const res = await adminLogin(email, password);
      saveToken(res.access_token);
      router.replace("/dashboard");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Terjadi kesalahan, coba lagi");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div
      className="min-h-screen flex items-center justify-center"
      style={{ background: "var(--bg)" }}
    >
      {/* Brand strip on left (decorative) */}
      <div
        className="fixed inset-y-0 left-0 w-1.5 hidden md:block"
        style={{ background: "var(--brand-grad)" }}
      />

      <div
        className="w-full max-w-sm reveal"
        style={{
          background: "var(--surface)",
          borderRadius: "var(--radius)",
          boxShadow: "var(--shadow-hover)",
          padding: "2.5rem",
          border: "1px solid var(--border)",
        }}
      >
        {/* Logo / wordmark */}
        <div className="mb-8">
          <div
            className="inline-flex items-center justify-center w-10 h-10 rounded-xl mb-3"
            style={{ background: "var(--brand-grad)" }}
          >
            <svg className="w-5 h-5 text-white fill-current" viewBox="0 0 20 20">
              <path d="M10 2a8 8 0 100 16A8 8 0 0010 2zm0 3a1 1 0 011 1v3.586l2.707 2.707a1 1 0 01-1.414 1.414L9.586 11H6a1 1 0 110-2h2.586V6a1 1 0 011-1z" />
            </svg>
          </div>
          <h1 className="text-xl font-semibold" style={{ color: "var(--text)" }}>
            BaliTravelHealth
          </h1>
          <p className="text-sm mt-0.5" style={{ color: "var(--text-muted)" }}>
            Masuk ke panel admin
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label
              className="block text-sm font-medium mb-1.5"
              style={{ color: "var(--text)" }}
            >
              Email
            </label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="admin@balihealth.me"
              className="w-full px-3 py-2.5 text-sm rounded-lg outline-none transition-all"
              style={{
                border: "1px solid var(--border)",
                color: "var(--text)",
                background: "var(--surface)",
              }}
              onFocus={(e) => (e.target.style.borderColor = "var(--brand)")}
              onBlur={(e) => (e.target.style.borderColor = "var(--border)")}
            />
          </div>

          <div>
            <label
              className="block text-sm font-medium mb-1.5"
              style={{ color: "var(--text)" }}
            >
              Password
            </label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2.5 text-sm rounded-lg outline-none transition-all"
              style={{
                border: "1px solid var(--border)",
                color: "var(--text)",
                background: "var(--surface)",
              }}
              onFocus={(e) => (e.target.style.borderColor = "var(--brand)")}
              onBlur={(e) => (e.target.style.borderColor = "var(--border)")}
            />
          </div>

          {error && (
            <p
              className="text-sm px-3 py-2 rounded-lg"
              style={{
                color: "var(--danger)",
                background: "#fef2f2",
                border: "1px solid #fecaca",
              }}
            >
              {error}
            </p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full py-2.5 rounded-lg text-sm font-semibold text-white transition-opacity disabled:opacity-60"
            style={{ background: "var(--brand-grad)" }}
          >
            {loading ? "Memproses..." : "Masuk"}
          </button>
        </form>
      </div>
    </div>
  );
}
