import { useCallback } from 'react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { useAuth } from '../hooks/useAuth'
import { useIdleTimeout } from '../hooks/useIdleTimeout'
import { logout } from '../api/auth'
import IdleWarningModal from './IdleWarningModal'

const navItems = [
  { to: '/',        label: 'Dashboard',       end: true  },
  { to: '/patients', label: 'Data Pasien',    end: false },
  { to: '/export',   label: 'Laporan',        end: false },
  { to: '/users',    label: 'Manajemen User', superAdminOnly: true, end: false },
]

export default function Layout() {
  const { user, signOut, isSuperAdmin } = useAuth()
  const navigate = useNavigate()

  const handleLogout = useCallback(async () => {
    const name = user?.name ?? ''
    try { await logout() } catch {}
    toast.success(`Sampai jumpa, ${name}! Anda telah keluar.`)
    signOut()
    navigate('/login')
  }, [user?.name, signOut, navigate])

  const handleIdleLogout = useCallback(async () => {
    try { await logout() } catch {}
    toast('Sesi berakhir otomatis karena tidak ada aktivitas.', {
      icon: '🔒',
      duration: 4000,
    })
    signOut()
    navigate('/login')
  }, [signOut, navigate])

  const { showWarning, countdown, stayLoggedIn } = useIdleTimeout({
    idleMinutes: 1ok ,
    warningSeconds: 60,
    onLogout: handleIdleLogout,
  })

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg)' }}>
      {/* Topbar */}
      <header className="topbar">
        {/* Logo */}
        <div className="flex items-center gap-2.5 shrink-0">
          <div className="w-8 h-8 rounded-[9px] flex items-center justify-center text-white text-base"
            style={{ background: 'var(--teal)' }}>
            🦷
          </div>
          <span className="font-semibold text-[14px] tracking-tight" style={{ color: 'var(--text)' }}>
            Klinik Gigi Sehat
          </span>
        </div>

        {/* Nav */}
        <nav className="flex items-center gap-0.5">
          {navItems.map((item) => {
            if (item.superAdminOnly && !isSuperAdmin) return null
            return (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.end}
                className={({ isActive }) =>
                  `btn btn-ghost btn-sm ${isActive ? '!bg-[var(--teal-l)] !text-[var(--teal-d)] !font-semibold' : ''}`
                }
              >
                {item.label}
              </NavLink>
            )
          })}
        </nav>

        <div className="flex-1" />

        {/* User chip */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            <div className="w-7 h-7 rounded-full flex items-center justify-center text-[11px] font-semibold text-white"
              style={{ background: 'var(--teal-m)' }}>
              {user?.name?.charAt(0).toUpperCase()}
            </div>
            <div className="leading-none">
              <p className="text-[13px] font-medium" style={{ color: 'var(--text)' }}>{user?.name}</p>
              <p className="text-[11px]" style={{ color: 'var(--text3)' }}>
                {user?.role === 'superadmin' ? 'Super Admin' : user?.role === 'write' ? 'Dokter' : 'Suster'}
              </p>
            </div>
          </div>
          <div style={{ width: 1, height: 20, background: 'var(--border)' }} />
          <button
            onClick={handleLogout}
            className="btn btn-ghost btn-sm"
            style={{ color: 'var(--text3)' }}
          >
            Keluar
          </button>
        </div>
      </header>

      {/* Content */}
      <main className="mx-auto page-enter" style={{ maxWidth: 1200, padding: '24px 24px' }}>
        <Outlet />
      </main>

      {showWarning && (
        <IdleWarningModal
          countdown={countdown}
          onStay={stayLoggedIn}
          onLogout={handleIdleLogout}
        />
      )}
    </div>
  )
}
