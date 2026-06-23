interface StatCardProps {
  label: string;
  value: string | number;
  icon: React.ReactNode;
  color: string;
  index?: number;
}

export default function StatCard({ label, value, icon, color, index = 0 }: StatCardProps) {
  return (
    <div
      className="reveal card-hover p-5 flex items-center gap-4"
      style={{
        background: "var(--surface)",
        borderRadius: "var(--radius)",
        boxShadow: "var(--shadow-card)",
        border: "1px solid var(--border)",
        animationDelay: `${(index + 1) * 0.05}s`,
      }}
    >
      <div
        className="w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0"
        style={{ background: `${color}18`, color }}
      >
        {icon}
      </div>
      <div>
        <p className="text-3xl font-bold leading-none" style={{ color: "var(--text)" }}>
          {value}
        </p>
        <p className="text-sm mt-1" style={{ color: "var(--text-muted)" }}>
          {label}
        </p>
      </div>
    </div>
  );
}
