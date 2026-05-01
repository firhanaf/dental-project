import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { getPatients } from '../api/patients'
import { getBranches } from '../api/branches'
import { useAuth } from '../hooks/useAuth'
import { Spinner, Button, EmptyState, PageHeader, formatDate, formatGender } from '../components/ui'

export default function PatientListPage() {
  const { isWrite } = useAuth()
  const navigate = useNavigate()

  const [search, setSearch] = useState('')
  const [status, setStatus] = useState('')
  const [branchId, setBranchId] = useState('')
  const [page, setPage] = useState(1)

  const { data, isLoading, isFetching } = useQuery({
    queryKey: ['patients', { search, status, branchId, page }],
    queryFn: () => getPatients({
      page, limit: 20,
      search: search || undefined,
      status: status || undefined,
      branch_id: branchId || undefined,
    }),
    placeholderData: (prev) => prev,
  })

  const { data: branches } = useQuery({ queryKey: ['branches'], queryFn: getBranches })

  const handleSearch = (v: string) => { setSearch(v); setPage(1) }

  return (
    <div>
      <PageHeader
        title="Data Pasien"
        action={isWrite ? (
          <Button variant="primary" onClick={() => navigate('/patients/new')}>
            + Tambah Pasien
          </Button>
        ) : undefined}
      />

      {/* Toolbar */}
      <div className="flex flex-wrap items-center gap-2.5 mb-4">
        {/* Search */}
        <div className="relative">
          <svg className="absolute left-2.5 top-1/2 -translate-y-1/2 pointer-events-none"
            width="15" height="15" fill="none" viewBox="0 0 24 24"
            stroke="var(--text3)" strokeWidth="2">
            <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
          </svg>
          <input
            className="form-input pl-8"
            style={{ width: 240 }}
            placeholder="Cari nama, No. RM, telepon..."
            value={search}
            onChange={(e) => handleSearch(e.target.value)}
          />
        </div>

        <select
          className="form-select"
          style={{ width: 160 }}
          value={status}
          onChange={(e) => { setStatus(e.target.value); setPage(1) }}
        >
          <option value="">Semua Status</option>
          <option value="new">Baru</option>
          <option value="active">Aktif</option>
          <option value="needs_control">Perlu Kontrol</option>
        </select>

        <select
          className="form-select"
          style={{ width: 200 }}
          value={branchId}
          onChange={(e) => { setBranchId(e.target.value); setPage(1) }}
        >
          <option value="">Semua Cabang</option>
          {branches?.map((b) => (
            <option key={b.id} value={b.id}>{b.name}</option>
          ))}
        </select>

        {(search || status || branchId) && (
          <button className="btn btn-ghost btn-sm"
            onClick={() => { setSearch(''); setStatus(''); setBranchId(''); setPage(1) }}>
            Reset filter
          </button>
        )}

        {isFetching && !isLoading && (
          <div className="ml-auto"><Spinner size="sm" /></div>
        )}
      </div>

      {/* Table */}
      <div className="card overflow-hidden">
        {isLoading ? (
          <div className="flex justify-center py-16"><Spinner /></div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>No. RM</th>
                    <th>Nama Pasien</th>
                    <th>Tgl Lahir</th>
                    <th>Jenis Kelamin</th>
                    <th>Telepon</th>
                    <th>Cabang</th>
                    <th>Status</th>
                    <th>Kunjungan Terakhir</th>
                    <th>Total Kunjungan</th>
                  </tr>
                </thead>
                <tbody>
                  {(data?.data ?? []).length === 0 && (
                    <tr>
                      <td colSpan={9} className="!p-0">
                        <EmptyState message="Tidak ada pasien ditemukan" />
                      </td>
                    </tr>
                  )}
                  {(data?.data ?? []).map((p) => (
                    <tr key={p.id} onClick={() => navigate(`/patients/${p.id}`)}>
                      <td>
                        <span className="mono font-medium" style={{ color: 'var(--teal)' }}>{p.no_rm}</span>
                      </td>
                      <td>
                        <div className="flex items-center gap-2.5">
                          <div className="w-7 h-7 rounded-full flex items-center justify-center text-[11px] font-semibold shrink-0"
                            style={{ background: 'var(--teal-l)', color: 'var(--teal-d)' }}>
                            {p.name.charAt(0).toUpperCase()}
                          </div>
                          <span className="font-medium" style={{ color: 'var(--text)' }}>{p.name}</span>
                        </div>
                      </td>
                      <td style={{ color: 'var(--text2)' }}>{formatDate(p.date_of_birth)}</td>
                      <td style={{ color: 'var(--text2)' }}>{formatGender(p.gender)}</td>
                      <td style={{ color: 'var(--text2)' }}>{p.phone}</td>
                      <td style={{ color: 'var(--text2)' }}>{p.branch_name}</td>
                      <td><span className={`badge badge-${p.status}`}>
                        {p.status === 'new' ? 'Baru' : p.status === 'active' ? 'Aktif' : 'Perlu Kontrol'}
                      </span></td>
                      <td style={{ color: 'var(--text2)' }}>{formatDate(p.last_visit_date)}</td>
                      <td style={{ color: 'var(--text2)' }}>{p.total_visits}x</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {data && (
              <div className="flex items-center justify-between px-4 py-3"
                style={{ borderTop: '1px solid var(--border)' }}>
                <span className="text-[12px]" style={{ color: 'var(--text3)' }}>
                  {data.meta.total} pasien
                  {data.meta.total > 0 && ` · halaman ${data.meta.page} dari ${Math.ceil(data.meta.total / data.meta.limit)}`}
                </span>
                <div className="flex gap-2">
                  <button className="btn btn-sm" disabled={page === 1} onClick={() => setPage(p => p - 1)}>
                    ‹ Sebelumnya
                  </button>
                  <button className="btn btn-sm" disabled={!data.meta.has_next} onClick={() => setPage(p => p + 1)}>
                    Berikutnya ›
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}
