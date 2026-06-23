interface FormFieldProps {
  label: string;
  required?: boolean;
  children: React.ReactNode;
  hint?: string;
}

export default function FormField({ label, required, children, hint }: FormFieldProps) {
  return (
    <div>
      <label className="block text-sm font-medium mb-1.5" style={{ color: "var(--text)" }}>
        {label}
        {required && <span className="ml-0.5" style={{ color: "var(--danger)" }}>*</span>}
      </label>
      {children}
      {hint && (
        <p className="text-xs mt-1" style={{ color: "var(--text-muted)" }}>
          {hint}
        </p>
      )}
    </div>
  );
}

export function Input({
  ...props
}: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className="w-full px-3 py-2 text-sm rounded-lg outline-none transition-all"
      style={{
        border: "1px solid var(--border)",
        color: "var(--text)",
        background: "var(--surface)",
        ...((props.style as React.CSSProperties) ?? {}),
      }}
      onFocus={(e) => {
        e.target.style.borderColor = "var(--brand)";
        props.onFocus?.(e);
      }}
      onBlur={(e) => {
        e.target.style.borderColor = "var(--border)";
        props.onBlur?.(e);
      }}
    />
  );
}

export function Textarea({
  ...props
}: React.TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      rows={3}
      {...props}
      className="w-full px-3 py-2 text-sm rounded-lg outline-none transition-all resize-none"
      style={{
        border: "1px solid var(--border)",
        color: "var(--text)",
        background: "var(--surface)",
        ...((props.style as React.CSSProperties) ?? {}),
      }}
      onFocus={(e) => {
        e.target.style.borderColor = "var(--brand)";
        props.onFocus?.(e);
      }}
      onBlur={(e) => {
        e.target.style.borderColor = "var(--border)";
        props.onBlur?.(e);
      }}
    />
  );
}

export function Select({
  ...props
}: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      {...props}
      className="w-full px-3 py-2 text-sm rounded-lg outline-none transition-all"
      style={{
        border: "1px solid var(--border)",
        color: "var(--text)",
        background: "var(--surface)",
        ...((props.style as React.CSSProperties) ?? {}),
      }}
      onFocus={(e) => {
        e.target.style.borderColor = "var(--brand)";
        props.onFocus?.(e);
      }}
      onBlur={(e) => {
        e.target.style.borderColor = "var(--border)";
        props.onBlur?.(e);
      }}
    />
  );
}
