import { useState, useEffect, type FormEvent } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { getVisit, createVisit, updateVisit } from '../api/visits'
import { getUsers } from '../api/users'
import { useAuth } from '../hooks/useAuth'
import { Button, FormField, Spinner } from './ui'
import DateInput from './DateInput'

interface Props {
  patientId: string
  visitId?: string
  onClose: () => void
  onSuccess: () => void
}

export default function VisitFormModal({ patientId, visitId, onClose, onSuccess }: Props) {
  const isEdit = !!visitId
  const { user, isSuperAdmin } = useAuth()

  const { data: visitData, isLoading: visitLoading } = useQuery({
    queryKey: ['visit', visitId],
    queryFn: () => getVisit(visitId!),
    enabled: isEdit,
  })

  const { data: users } = useQuery({
    queryKey: ['users'],
    queryFn: getUsers,
    enabled: isSuperAdmin,
  })

  const doctors = isSuperAdmin
    ? (users ?? []).filter((u) => u.role === 'write' && u.is_active)
    : null

  const [form, setForm] = useState({
    visit_date: new Date().toISOString().slice(0, 10),
    doctor_id: isSuperAdmin ? '' : (user?.id ?? ''),
    chief_complaint: '',
    diagnosis: '',
    treatment: '',
    teeth_involved: '',
    cost: '0',
    next_control_date: '',
    notes: '',
  })

  useEffect(() => {
    if (visitData) {
      setForm({
        visit_date: visitData.visit_date.slice(0, 10),
        doctor_id: visitData.doctor_id,
        chief_complaint: visitData.chief_complaint,
        diagnosis: visitData.diagnosis ?? '',
        treatment: visitData.treatment ?? '',
        teeth_involved: visitData.teeth_involved ?? '',
        cost: String(visitData.cost),
        next_control_date: visitData.next_control_date?.slice(0, 10) ?? '',
        notes: visitData.notes ?? '',
      })
    }
  }, [visitData])

  const mutation = useMutation({
    mutationFn: () => {
      const payload = {
        patient_id: patientId,
        doctor_id: form.doctor_id || user!.id,
        visit_date: form.visit_date,
        chief_complaint: form.chief_complaint,
        diagnosis: form.diagnosis || undefined,
        treatment: form.treatment || undefined,
        teeth_involved: form.teeth_involved || undefined,
        cost: parseFloat(form.cost) || 0,
        next_control_date: form.next_control_date || undefined,
        notes: form.notes || undefined,
      }
      return isEdit ? updateVisit(visitId!, payload) : createVisit(payload)
    },
    onSuccess: () => {
      toast.success(isEdit ? 'Kunjungan berhasil diperbarui' : 'Kunjungan berhasil ditambahkan')
      onSuccess()
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message ?? 'Terjadi kesalahan', { duration: 3000 })
    },
  })

  const set = (field: string, value: string) => setForm((p) => ({ ...p, [field]: value }))

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    mutation.mutate()
  }

  return (
    <div className="modal-overlay">
      <div className="modal" style={{ maxWidth: 520 }}>
        <div className="modal-header">
          <h3>{isEdit ? 'Edit Kunjungan' : 'Tambah Kunjungan'}</h3>
          <button className="modal-close" onClick={onClose}>&times;</button>
        </div>

        {isEdit && visitLoading ? (
          <div className="flex justify-center py-12"><Spinner /></div>
        ) : (
          <form onSubmit={handleSubmit} className="modal-body space-y-4">
            <div className="grid grid-cols-2 gap-3">
              <FormField label="Tanggal Kunjungan" required>
                <DateInput
                  value={form.visit_date}
                  onChange={(v) => set('visit_date', v)}
                  required
                />
              </FormField>
              <FormField label="Biaya (Rp)" required>
                <input
                  className="form-input"
                  type="number"
                  min="0"
                  step="1000"
                  value={form.cost}
                  onChange={(e) => set('cost', e.target.value)}
                  required
                />
              </FormField>
            </div>

            {isSuperAdmin && doctors && (
              <FormField label="Dokter" required>
                <select
                  className="form-select"
                  value={form.doctor_id}
                  onChange={(e) => set('doctor_id', e.target.value)}
                  required
                >
                  <option value="">Pilih dokter...</option>
                  {doctors.map((d) => (
                    <option key={d.id} value={d.id}>{d.name}</option>
                  ))}
                </select>
              </FormField>
            )}

            <FormField label="Keluhan Utama" required>
              <input
                className="form-input"
                value={form.chief_complaint}
                onChange={(e) => set('chief_complaint', e.target.value)}
                required
              />
            </FormField>

            <FormField label="Diagnosis">
              <input
                className="form-input"
                value={form.diagnosis}
                onChange={(e) => set('diagnosis', e.target.value)}
              />
            </FormField>

            <FormField label="Tindakan">
              <input
                className="form-input"
                value={form.treatment}
                onChange={(e) => set('treatment', e.target.value)}
              />
            </FormField>

            <div className="grid grid-cols-2 gap-3">
              <FormField label="Gigi yang Terlibat">
                <input
                  className="form-input"
                  value={form.teeth_involved}
                  onChange={(e) => set('teeth_involved', e.target.value)}
                  placeholder="Contoh: 16,17,36"
                />
              </FormField>
              <FormField label="Tanggal Kontrol">
                <DateInput
                  value={form.next_control_date}
                  onChange={(v) => set('next_control_date', v)}
                />
              </FormField>
            </div>

            <FormField label="Catatan">
              <input
                className="form-input"
                value={form.notes}
                onChange={(e) => set('notes', e.target.value)}
              />
            </FormField>

            <div className="flex gap-3 pt-1">
              <Button type="submit" variant="primary" loading={mutation.isPending}>
                {isEdit ? 'Simpan' : 'Tambah Kunjungan'}
              </Button>
              <Button type="button" variant="secondary" onClick={onClose}>Batal</Button>
            </div>
          </form>
        )}
      </div>
    </div>
  )
}
