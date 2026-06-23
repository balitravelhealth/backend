"use client";

import React from "react";

export type Language = "id" | "en";

interface LanguageTabsProps {
  value: Language;
  onChange: (lang: Language) => void;
}

export default function LanguageTabs({ value, onChange }: LanguageTabsProps) {
  return (
    <div className="flex gap-2 border-b mb-4">
      <button
        onClick={() => onChange("id")}
        className="px-4 py-2 font-medium transition-colors"
        style={{
          color: value === "id" ? "var(--primary)" : "var(--text-muted)",
          borderBottom: value === "id" ? "2px solid var(--primary)" : "none",
        }}
      >
        🇮🇩 Bahasa Indonesia
      </button>
      <button
        onClick={() => onChange("en")}
        className="px-4 py-2 font-medium transition-colors"
        style={{
          color: value === "en" ? "var(--primary)" : "var(--text-muted)",
          borderBottom: value === "en" ? "2px solid var(--primary)" : "none",
        }}
      >
        🇬🇧 English
      </button>
    </div>
  );
}
