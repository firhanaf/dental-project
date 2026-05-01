import type { ReactNode, ButtonHTMLAttributes, InputHTMLAttributes, SelectHTMLAttributes } from 'react'

/* ── Spinner ─────────────────────────────────────────────── */
export function Spinner({ size = 'md' }: { size?: 'sm' | 'md' | 'lg' }) {
  const s = size === 'sm' ? '14px' : size === 'lg' ? '36px' : '20px'
  return (
    <div style={{ width: s, height: s }}
      className="rounded-full border-2 border-[var(--border2)] border-t-[var(--teal-m)] animate-spin shrink-0" />
  )
}

/* ── Card ────────────────────────────────────────────────── */
export function Card({ children, className = '' }: { children: ReactNode; className?: string }) {
  return <div className={`card ${className}`}>{children}</div>
}

/* ── Button ──────────────────────────────────────────────── */
type BtnVariant = 'primary' | 'secondary' | 'danger' | 'ghost'
type BtnSize = 'sm' | 'md'

export function Button({
  variant = 'secondary',
  size = 'md',
  loading = false,
  className = '',
  children,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: BtnVariant
  size?: BtnSize
  loading?: boolean
}) {
  const variantClass = variant === 'primary' ? 'btn-primary'
    : variant === 'danger' ? 'btn-danger'
    : variant === 'ghost' ? 'btn-ghost'
    : ''
  return (
    <button
      {...props}
      disabled={props.disabled || loading}
      className={`btn ${variantClass} ${size === 'sm' ? 'btn-sm' : ''} ${className}`}
    >
      {loading && <Spinner size="sm" />}
      {children}
    </button>
  )
}

/* ── Input ───────────────────────────────────────────────── */
export function Input({ className = '', ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return <input {...props} className={`form-input ${className}`} />
}

/* ── Select ──────────────────────────────────────────────── */
export function Select({
  className = '', children, ...props
}: SelectHTMLAttributes<HTMLSelectElement> & { children: ReactNode }) {
  return (
    <select {...props} className={`form-select ${className}`}>
      {children}
    </select>
  )
}

/* ── FormField ───────────────────────────────────────────── */
export function FormField({
  label, required, children, hint,
}: { label: string; required?: boolean; children: ReactNode; hint?: string }) {
  return (
    <div>
      <label className="form-label">
        {label}{required && <span style={{ color: 'var(--danger-t)' }} className="ml-0.5">*</span>}
      </label>
      {children}
      {hint && <p className="text-[11px] mt-1" style={{ color: 'var(--text3)' }}>{hint}</p>}
    </div>
  )
}

/* ── Badge ───────────────────────────────────────────────── */
const badgeLabel: Record<string, string> = {
  new: 'Baru',
  active: 'Aktif',
  needs_control: 'Perlu Kontrol',
  superadmin: 'Super Admin',
  write: 'Dokter',
  readonly: 'Suster',
}

export function Badge({ status }: { status: string }) {
  return (
    <span className={`badge badge-${status}`}>
      {badgeLabel[status] ?? status}
    </span>
  )
}

/* ── PageHeader ──────────────────────────────────────────── */
export function PageHeader({ title, action }: { title: string; action?: ReactNode }) {
  return (
    <div className="flex items-center justify-between mb-5">
      <h1 className="text-[18px] font-semibold tracking-tight" style={{ color: 'var(--text)' }}>
        {title}
      </h1>
      {action}
    </div>
  )
}

/* ── EmptyState ──────────────────────────────────────────── */
export function EmptyState({ message }: { message: string }) {
  return <div className="empty-state">{message}</div>
}

/* ── ErrorMessage ────────────────────────────────────────── */
export function ErrorMessage({ message }: { message: string }) {
  return (
    <div className="flex items-start gap-2 rounded-[var(--radius)] px-4 py-3 text-[13px]"
      style={{ background: 'var(--danger)', color: 'var(--danger-t)', border: '1px solid #FECACA' }}>
      {message}
    </div>
  )
}

/* ── Formatters ──────────────────────────────────────────── */
export function formatDate(iso: string | null | undefined) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

export function formatRupiah(n: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency', currency: 'IDR', maximumFractionDigits: 0,
  }).format(n)
}

export function formatGender(g: string) {
  return g === 'male' ? 'Laki-laki' : 'Perempuan'
}
