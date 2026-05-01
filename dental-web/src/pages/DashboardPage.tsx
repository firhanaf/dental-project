import { useQueries, useQuery } from '@tanstack/react-query'
import { getPatients } from '../api/patients'
import { getBranches } from '../api/branches'
import { Spinner } from '../components/ui'
import type { PatientStatus } from '../types'

const statuses: { key: PatientStatus; label: string; color: string }[] = [
  { key: 'new',           label: 'Pasien Baru',    color: '#3B82F6' },
  { key: 'active',        label: 'Aktif',          color: '#22C55E' },
  { key: 'needs_control', label: 'Perlu Kontrol',  color: '#F59E0B' },
]

export default function DashboardPage() {
  const { data: totalData, isLoading } = useQuery({
    queryKey: ['patients', 'total'],
    queryFn: () => getPatients({ limit: 1 }),
  })

  const statusQueries = useQueries({
    queries: statuses.map(({ key }) => ({
      queryKey: ['patients', 'status', key],
      queryFn: () => getPatients({ limit: 1, status: key }),
    })),
  })

  const { data: branches } = useQuery({ queryKey: ['branches'], queryFn: getBranches })

  if (isLoading) {
    return <div className="flex justify-center py-20"><Spinner size="lg" /></div>
  }

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-[18px] font-semibold tracking-tight" style={{ color: 'var(--text)' }}>
          Dashboard
        </h1>
        <p className="text-[13px] mt-0.5" style={{ color: 'var(--text3)' }}>
          Ringkasan data klinik
        </p>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-4 gap-3.5">
        <div className="card stat-card">
          <div className="stat-label">Total Pasien</div>
          <div className="stat-value">{totalData?.meta.total ?? '—'}</div>
          <div className="text-[12px] mt-2" style={{ color: 'var(--text3)' }}>semua cabang</div>
        </div>
        {statuses.map(({ key, label, color }, i) => (
          <div key={key} className="card stat-card">
            <div className="stat-label">{label}</div>
            <div className="stat-value" style={{ color }}>{statusQueries[i].data?.meta.total ?? '—'}</div>
          </div>
        ))}
      </div>

      {/* Cabang */}
      <div className="card">
        <div className="card-header"><h2>Cabang</h2></div>
        <div className="divide-y" style={{ borderColor: 'var(--border)' }}>
          {branches?.length === 0 && (
            <div className="empty-state">Tidak ada cabang</div>
          )}
          {branches?.map((b) => (
            <div key={b.id} className="flex items-center justify-between px-5 py-4">
              <div>
                <p className="font-semibold text-[14px]" style={{ color: 'var(--text)' }}>{b.name}</p>
                <p className="text-[12px] mt-0.5" style={{ color: 'var(--text3)' }}>
                  {b.address ?? '—'} · {b.phone ?? '—'}
                </p>
              </div>
              <span className={`badge ${b.is_active ? 'badge-active' : 'badge-readonly'}`}>
                {b.is_active ? 'Aktif' : 'Nonaktif'}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
