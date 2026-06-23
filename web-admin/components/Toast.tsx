"use client";

import { useEffect, useCallback, useState, useMemo } from "react";

export interface ToastMsg {
  id: number;
  type: "success" | "error";
  text: string;
}

let _id = 0;

export function useToast() {
  const [toasts, setToasts] = useState<ToastMsg[]>([]);

  const show = useCallback((type: ToastMsg["type"], text: string) => {
    const id = ++_id;
    setToasts((prev) => [...prev, { id, type, text }]);
  }, []);

  const remove = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const toast = useMemo(() => ({
    success: (text: string) => show("success", text),
    error:   (text: string) => show("error", text),
  }), [show]);

  return { toasts, remove, toast };
}

function ToastItem({ t, onRemove }: { t: ToastMsg; onRemove: () => void }) {
  useEffect(() => {
    const timer = setTimeout(onRemove, 4000);
    return () => clearTimeout(timer);
  }, [onRemove]);

  const s =
    t.type === "success"
      ? { color: "#15803d", bg: "#f0fdf4", border: "#86efac" }
      : { color: "#dc2626", bg: "#fef2f2", border: "#fecaca" };

  return (
    <div
      className="flex items-center gap-2.5 px-4 py-3 rounded-xl text-sm font-medium shadow-lg animate-slideUp pointer-events-auto"
      style={{
        color: s.color, background: s.bg,
        border: `1px solid ${s.border}`,
        minWidth: 260, maxWidth: 400,
      }}
    >
      {t.type === "success" ? (
        <svg className="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      ) : (
        <svg className="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      )}
      <span className="flex-1">{t.text}</span>
      <button onClick={onRemove} className="opacity-50 hover:opacity-100 transition-opacity ml-1">
        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  );
}

export function Toaster({
  toasts,
  remove,
}: {
  toasts: ToastMsg[];
  remove: (id: number) => void;
}) {
  if (toasts.length === 0) return null;
  return (
    <div className="fixed bottom-6 right-6 z-[200] flex flex-col gap-2 pointer-events-none">
      {toasts.map((t) => (
        <ToastItem key={t.id} t={t} onRemove={() => remove(t.id)} />
      ))}
    </div>
  );
}
