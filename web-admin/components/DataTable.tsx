interface Column<T> {
  key: keyof T | string;
  label: string;
  render?: (row: T) => React.ReactNode;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  rows: T[];
  keyField: keyof T;
  loading?: boolean;
  emptyText?: string;
}

export default function DataTable<T>({
  columns,
  rows,
  keyField,
  loading,
  emptyText = "Tidak ada data",
}: DataTableProps<T>) {
  return (
    <div
      className="overflow-x-auto rounded-xl border"
      style={{
        background: "var(--surface)",
        borderColor: "var(--border)",
        boxShadow: "var(--shadow-card)",
      }}
    >
      <table className="w-full text-sm">
        <thead>
          <tr style={{ borderBottom: "1px solid var(--border)" }}>
            {columns.map((col) => (
              <th
                key={String(col.key)}
                className="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wide"
                style={{ color: "var(--text-muted)" }}
              >
                {col.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {loading ? (
            <tr>
              <td
                colSpan={columns.length}
                className="text-center py-10 text-sm"
                style={{ color: "var(--text-muted)" }}
              >
                Memuat…
              </td>
            </tr>
          ) : rows.length === 0 ? (
            <tr>
              <td
                colSpan={columns.length}
                className="text-center py-10 text-sm"
                style={{ color: "var(--text-muted)" }}
              >
                {emptyText}
              </td>
            </tr>
          ) : (
            rows.map((row, idx) => (
              <tr
                key={String(row[keyField])}
                style={{
                  borderBottom:
                    idx < rows.length - 1
                      ? "1px solid var(--border)"
                      : "none",
                  background:
                    idx % 2 === 1 ? "rgba(244,246,251,0.5)" : "transparent",
                }}
              >
                {columns.map((col) => (
                  <td
                    key={String(col.key)}
                    className="px-4 py-3 align-top"
                    style={{ color: "var(--text)" }}
                  >
                    {col.render
                      ? col.render(row)
                      : String((row as Record<string, unknown>)[String(col.key)] ?? "—")}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
