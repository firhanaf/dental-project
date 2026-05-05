import { useState, useEffect, type FormEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { getPatient, createPatient, updatePatient } from '../api/patients'
import { getBranches } from '../api/branches'
import { useAuth } from '../hooks/useAuth'
import { Button, FormField, PageHeader, Spinner } from '../components/ui'
import DateInput from '../components/DateInput'

export default function PatientFormPage() {
  const { id } = useParams<{ id: string }>()
  const isEdit = !!id
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { isSuperAdmin } = useAuth()

  const { data: patient, isLoading: patientLoading } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => getPatient(id!),
    enabled: isEdit,
  })

  const { data: branches } = useQuery({
    queryKey: ['branches'],
    queryFn: getBranches,
    enabled: isSuperAdmin,
  })

  const [form, setForm] = useState({
    name: '',
    nik: '',
    date_of_birth: '',
    gender: 'male',
    phone: '',
    address: '',
    occupation: '',
    allergy_notes: '',
    branch_id: '',
  })
  useEffect(() => {
    if (patient) {
      setForm({
        name: patient.name,
        nik: patient.nik ?? '',
        date_of_birth: patient.date_of_birth.slice(0, 10),
        gender: patient.gender,
        phone: patient.phone,
        address: patient.address ?? '',
        occupation: patient.occupation ?? '',
        allergy_notes: patient.allergy_notes ?? '',
        branch_id: patient.branch_id,
      })
    }
  }, [patient])

  const mutation = useMutation({
    mutationFn: (data: typeof form) => {
      const payload = {
        ...data,
        nik: data.nik || undefined,
        address: data.address || undefined,
        occupation: data.occupation || undefined,
        allergy_notes: data.allergy_notes || undefined,
        branch_id: isSuperAdmin ? (data.branch_id || undefined) : undefined,
      }
      return isEdit ? updatePatient(id!, payload) : createPatient(payload)
    },
    onSuccess: (saved) => {
      qc.invalidateQueries({ queryKey: ['patients'] })
      qc.invalidateQueries({ queryKey: ['patient', saved.id] })
      toast.success(isEdit ? 'Data pasien berhasil diperbarui' : 'Pasien berhasil didaftarkan')
      setTimeout(() => navigate(`/patients/${saved.id}`), 1000)
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message ?? 'Terjadi kesalahan', { duration: 3000 })
    },
  })

  const set = (field: string, value: string) =>
    setForm((prev) => ({ ...prev, [field]: value }))

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (form.nik && form.nik.length !== 16) {
      toast.error('NIK harus 16 digit angka', { duration: 3000 })
      return
    }
    mutation.mutate(form)
  }

  if (isEdit && patientLoading) {
    return <div className="flex justify-center py-20"><Spinner size="lg" /></div>
  }

  return (
    <div style={{ maxWidth: 560 }}>
      <PageHeader title={isEdit ? 'Edit Pasien' : 'Tambah Pasien Baru'} />

      <div className="card" style={{ padding: 24 }}>
        <form onSubmit={handleSubmit} className="space-y-4">
          <FormField label="Nama Lengkap" required>
            <input
              className="form-input"
              value={form.name}
              onChange={(e) => set('name', e.target.value)}
              required
            />
          </FormField>

          <div className="grid grid-cols-2 gap-3">
            <FormField label="Tanggal Lahir" required>
              <DateInput
                value={form.date_of_birth}
                onChange={(v) => set('date_of_birth', v)}
                required
              />
            </FormField>
            <FormField label="Jenis Kelamin" required>
              <select
                className="form-select"
                value={form.gender}
                onChange={(e) => set('gender', e.target.value)}
              >
                <option value="male">Laki-laki</option>
                <option value="female">Perempuan</option>
              </select>
            </FormField>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <FormField label="No. Telepon" required>
              <input
                className="form-input"
                value={form.phone}
                onChange={(e) => set('phone', e.target.value)}
                placeholder="08xxxxxxxxxx"
                required
              />
            </FormField>
            <FormField label="NIK">
              <input
                className={`form-input ${form.nik && form.nik.length !== 16 ? 'invalid' : ''}`}
                value={form.nik}
                onChange={(e) => {
                  const val = e.target.value.replace(/\D/g, '')
                  set('nik', val)
                }}
                placeholder="16 digit (opsional)"
                maxLength={16}
                inputMode="numeric"
              />
              {form.nik && form.nik.length > 0 && form.nik.length !== 16 && (
                <p className="text-[11px] mt-1" style={{ color: 'var(--amber-t)' }}>
                  {form.nik.length}/16 digit
                </p>
              )}
            </FormField>
          </div>

          <FormField label="Alamat">
            <input
              className="form-input"
              value={form.address}
              onChange={(e) => set('address', e.target.value)}
            />
          </FormField>

          <FormField label="Pekerjaan">
            <input
              className="form-input"
              value={form.occupation}
              onChange={(e) => set('occupation', e.target.value)}
            />
          </FormField>

          <FormField label="Catatan Alergi">
            <input
              className="form-input"
              value={form.allergy_notes}
              onChange={(e) => set('allergy_notes', e.target.value)}
              placeholder="Contoh: Alergi penisilin"
            />
          </FormField>

          {isSuperAdmin && (
            <FormField label="Cabang" required>
              <select
                className="form-select"
                value={form.branch_id}
                onChange={(e) => set('branch_id', e.target.value)}
                required
              >
                <option value="">Pilih cabang...</option>
                {branches?.map((b) => (
                  <option key={b.id} value={b.id}>{b.name}</option>
                ))}
              </select>
            </FormField>
          )}

          {!isSuperAdmin && (
            <p className="text-[12px]" style={{ color: 'var(--text3)' }}>
              Pasien akan didaftarkan ke cabang Anda.
            </p>
          )}

          <div className="flex gap-3 pt-1">
            <Button type="submit" variant="primary" loading={mutation.isPending}>
              {isEdit ? 'Simpan Perubahan' : 'Daftarkan Pasien'}
            </Button>
            <Button type="button" variant="secondary" onClick={() => navigate(-1)}>
              Batal
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
