import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { login, resetPassword } from '../api/auth'
import { useAuth } from '../hooks/useAuth'
import { Spinner } from '../components/ui'

function EyeIcon({ open }: { open: boolean }) {
  return open ? (
    <svg width="15" height="15" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      <path strokeLinecap="round" strokeLinejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
    </svg>
  ) : (
    <svg width="15" height="15" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
    </svg>
  )
}

export default function LoginPage() {
  const { signIn } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [successName, setSuccessName] = useState('')

  // Forgot password state
  const [mode, setMode] = useState<'login' | 'forgot'>('login')
  const [resetForm, setResetForm] = useState({ email: '', token: '', newPassword: '', confirm: '' })
  const [resetLoading, setResetLoading] = useState(false)
  const [resetError, setResetError] = useState('')
  const [resetSuccess, setResetSuccess] = useState(false)

  const setReset = (field: string, value: string) =>
    setResetForm((p) => ({ ...p, [field]: value }))

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await login(email, password)
      signIn(res.token, res.user)
      setSuccessName(res.user.name)
      setTimeout(() => navigate('/'), 1500)
    } catch (err: any) {
      setError(err.response?.data?.message ?? 'Email atau password salah')
      setLoading(false)
    }
  }

  const handleResetSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setResetError('')
    if (resetForm.newPassword !== resetForm.confirm) {
      setResetError('Konfirmasi password tidak cocok')
      return
    }
    if (resetForm.newPassword.length < 8) {
      setResetError('Password baru minimal 8 karakter')
      return
    }
    setResetLoading(true)
    try {
      await resetPassword(resetForm.email, resetForm.token, resetForm.newPassword)
      setResetSuccess(true)
    } catch (err: any) {
      setResetError(err.response?.data?.message ?? 'Kode reset tidak valid atau sudah kedaluwarsa')
    } finally {
      setResetLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4"
      style={{ background: 'var(--bg)' }}>
      <div className="w-full" style={{ maxWidth: 380 }}>
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="w-12 h-12 rounded-2xl mx-auto mb-4 flex items-center justify-center text-2xl"
            style={{ background: 'var(--teal)' }}>
            🦷
          </div>
          <h1 className="text-[20px] font-semibold tracking-tight" style={{ color: 'var(--text)' }}>
            Klinik Gigi Sehat
          </h1>
          <p className="text-[13px] mt-1" style={{ color: 'var(--text3)' }}>
            Sistem Rekam Medis Internal
          </p>
        </div>

        <div className="card" style={{ padding: 28 }}>
          {successName ? (
            <div className="flex flex-col items-center gap-3 py-4 text-center">
              <div className="w-10 h-10 rounded-full flex items-center justify-center text-xl"
                style={{ background: 'var(--teal-l)' }}>
                ✓
              </div>
              <div>
                <p className="text-[14px] font-semibold" style={{ color: 'var(--teal-d)' }}>
                  Selamat datang, {successName}!
                </p>
                <p className="text-[12px] mt-1" style={{ color: 'var(--text3)' }}>
                  Mengarahkan ke dashboard...
                </p>
              </div>
              <Spinner size="sm" />
            </div>
          ) : mode === 'login' ? (
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="rounded-[var(--radius)] px-4 py-3 text-[13px]"
                style={{ background: 'var(--danger)', color: 'var(--danger-t)', border: '1px solid #FECACA' }}>
                {error}
              </div>
            )}

            <div>
              <label className="form-label">Email</label>
              <input
                className="form-input"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="nama@klinik.local"
                required
                autoFocus
              />
            </div>

            <div>
              <label className="form-label">Password</label>
              <div className="relative">
                <input
                  className="form-input"
                  style={{ paddingRight: 38 }}
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                />
                <button
                  type="button"
                  tabIndex={-1}
                  onClick={() => setShowPassword(v => !v)}
                  className="absolute inset-y-0 right-0 flex items-center px-3"
                  style={{ color: 'var(--text3)' }}
                >
                  <EyeIcon open={showPassword} />
                </button>
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary w-full justify-center"
              style={{ height: 40, fontSize: 14 }}
            >
              {loading ? <Spinner size="sm" /> : null}
              Masuk
            </button>

            <div className="text-center pt-1">
              <button
                type="button"
                onClick={() => { setMode('forgot'); setError('') }}
                className="text-[12px]"
                style={{ color: 'var(--teal)', background: 'none', border: 'none', cursor: 'pointer' }}
              >
                Lupa Password?
              </button>
            </div>
          </form>
          ) : resetSuccess ? (
            <div className="flex flex-col items-center gap-3 py-4 text-center">
              <div className="w-10 h-10 rounded-full flex items-center justify-center text-xl"
                style={{ background: 'var(--teal-l)' }}>
                ✓
              </div>
              <div>
                <p className="text-[14px] font-semibold" style={{ color: 'var(--teal-d)' }}>
                  Password berhasil direset!
                </p>
                <p className="text-[12px] mt-1" style={{ color: 'var(--text3)' }}>
                  Silakan login dengan password baru.
                </p>
              </div>
              <button
                className="btn btn-primary"
                style={{ fontSize: 13 }}
                onClick={() => { setMode('login'); setResetSuccess(false); setResetForm({ email: '', token: '', newPassword: '', confirm: '' }) }}
              >
                Kembali ke Login
              </button>
            </div>
          ) : (
          <form onSubmit={handleResetSubmit} className="space-y-4">
            <div>
              <p className="text-[13px] font-semibold mb-1" style={{ color: 'var(--text)' }}>Reset Password</p>
              <p className="text-[12px]" style={{ color: 'var(--text3)' }}>
                Masukkan email Anda dan kode reset yang diberikan admin.
              </p>
            </div>

            {resetError && (
              <div className="rounded-[var(--radius)] px-4 py-3 text-[13px]"
                style={{ background: 'var(--danger)', color: 'var(--danger-t)', border: '1px solid #FECACA' }}>
                {resetError}
              </div>
            )}

            <div>
              <label className="form-label">Email</label>
              <input
                className="form-input"
                type="email"
                value={resetForm.email}
                onChange={(e) => setReset('email', e.target.value)}
                placeholder="nama@klinik.local"
                required
                autoFocus
              />
            </div>

            <div>
              <label className="form-label">Kode Reset (dari Admin)</label>
              <input
                className="form-input"
                type="text"
                value={resetForm.token}
                onChange={(e) => setReset('token', e.target.value.toUpperCase())}
                placeholder="Contoh: ABC12DEF"
                maxLength={8}
                required
                style={{ letterSpacing: '0.15em', textTransform: 'uppercase' }}
              />
            </div>

            <div>
              <label className="form-label">Password Baru</label>
              <input
                className="form-input"
                type="password"
                value={resetForm.newPassword}
                onChange={(e) => setReset('newPassword', e.target.value)}
                placeholder="Min. 8 karakter"
                minLength={8}
                required
              />
            </div>

            <div>
              <label className="form-label">Konfirmasi Password Baru</label>
              <input
                className="form-input"
                type="password"
                value={resetForm.confirm}
                onChange={(e) => setReset('confirm', e.target.value)}
                placeholder="Ulangi password baru"
                minLength={8}
                required
              />
            </div>

            <button
              type="submit"
              disabled={resetLoading}
              className="btn btn-primary w-full justify-center"
              style={{ height: 40, fontSize: 14 }}
            >
              {resetLoading ? <Spinner size="sm" /> : null}
              Reset Password
            </button>

            <div className="text-center pt-1">
              <button
                type="button"
                onClick={() => { setMode('login'); setResetError('') }}
                className="text-[12px]"
                style={{ color: 'var(--text3)', background: 'none', border: 'none', cursor: 'pointer' }}
              >
                ← Kembali ke Login
              </button>
            </div>
          </form>
          )}
        </div>
      </div>
    </div>
  )
}
