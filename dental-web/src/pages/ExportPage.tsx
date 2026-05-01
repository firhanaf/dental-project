import { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { exportPatients, exportVisits } from '../api/export'
import { getBranches } from '../api/branches'
import { Button, FormField, PageHeader, ErrorMessage } from '../components/ui'

export default function ExportPage() {
  const [branchId, setBranchId] = useState('')
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')
  const [error, setError] = useState('')

  const { data: branches } = useQuery({ queryKey: ['branches'], queryFn: getBranches })

  const exportPatientsMutation = useMutation({
    mutationFn: () => exportPatients({ branch_id: branchId || undefined }),
    onError: () => setError('Gagal export data pasien'),
  })

  const exportVisitsMutation = useMutation({
    mutationFn: () =>
      exportVisits({
        branch_id: branchId || undefined,
        date_from: dateFrom || undefined,
        date_to: dateTo || undefined,
      }),
    onError: () => setError('Gagal export data kunjungan'),
  })

  return (
    <div style={{ maxWidth: 480 }}>
      <PageHeader title="Laporan" />

      <div className="card" style={{ padding: 24 }}>
        {error && <ErrorMessage message={error} />}

        <FormField label="Filter Cabang">
          <select
            className="form-select"
            value={branchId}
            onChange={(e) => setBranchId(e.target.value)}
          >
            <option value="">Semua Cabang</option>
            {branches?.map((b) => (
              <option key={b.id} value={b.id}>{b.name}</option>
            ))}
          </select>
        </FormField>

        {/* Export Pasien */}
        <div style={{ borderTop: '1px solid var(--border)', marginTop: 20, paddingTop: 20 }}>
          <p className="text-[13px] font-semibold" style={{ color: 'var(--text)', marginBottom: 4 }}>
            Export Pasien
          </p>
          <p className="text-[12px] mb-4" style={{ color: 'var(--text3)' }}>
            Semua data pasien beserta info cabang dan status terkini.
          </p>
          <Button
            variant="primary"
            onClick={() => { setError(''); exportPatientsMutation.mutate() }}
            loading={exportPatientsMutation.isPending}
          >
            Download Excel — Pasien
          </Button>
        </div>

        {/* Export Kunjungan */}
        <div style={{ borderTop: '1px solid var(--border)', marginTop: 20, paddingTop: 20 }}>
          <p className="text-[13px] font-semibold" style={{ color: 'var(--text)', marginBottom: 4 }}>
            Export Kunjungan
          </p>
          <p className="text-[12px] mb-4" style={{ color: 'var(--text3)' }}>
            Semua data kunjungan beserta diagnosis, tindakan, dan biaya.
          </p>
          <div className="grid grid-cols-2 gap-3 mb-4">
            <FormField label="Dari Tanggal">
              <input
                className="form-input"
                type="date"
                value={dateFrom}
                onChange={(e) => setDateFrom(e.target.value)}
              />
            </FormField>
            <FormField label="Sampai Tanggal">
              <input
                className="form-input"
                type="date"
                value={dateTo}
                onChange={(e) => setDateTo(e.target.value)}
              />
            </FormField>
          </div>
          <Button
            variant="primary"
            onClick={() => { setError(''); exportVisitsMutation.mutate() }}
            loading={exportVisitsMutation.isPending}
          >
            Download Excel — Kunjungan
          </Button>
        </div>
      </div>
    </div>
  )
}
