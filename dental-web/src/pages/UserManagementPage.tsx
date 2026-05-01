import { useState, type FormEvent } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getUsers, createUser, updateUser, deactivateUser } from '../api/users'
import { getBranches } from '../api/branches'
import {
  Button, FormField, Badge, EmptyState, PageHeader,
  ErrorMessage, formatDate, Spinner,
} from '../components/ui'

type ModalMode = 'create' | 'edit' | null

export default function UserManagementPage() {
  const qc = useQueryClient()
  const [modal, setModal] = useState<ModalMode>(null)
  const [editId, setEditId] = useState<string | null>(null)
  const [error, setError] = useState('')

  const [form, setForm] = useState({
    name: '', email: '', password: '', role: 'write', branch_id: '',
  })

  const { data: users, isLoading } = useQuery({ queryKey: ['users'], queryFn: getUsers })
  const { data: branches } = useQuery({ queryKey: ['branches'], queryFn: getBranches })

  const saveMutation = useMutation({
    mutationFn: () => {
      const payload = {
        ...form,
        branch_id: form.role === 'superadmin' ? null : (form.branch_id || null),
        password: form.password || undefined,
      }
      return editId ? updateUser(editId, payload) : createUser({ ...payload, password: form.password })
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['users'] })
      closeModal()
    },
    onError: (err: any) => setError(err.response?.data?.message ?? 'Terjadi kesalahan'),
  })

  const deactivateMutation = useMutation({
    mutationFn: (id: string) => deactivateUser(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['users'] }),
  })

  const openCreate = () => {
    setForm({ name: '', email: '', password: '', role: 'write', branch_id: '' })
    setEditId(null)
    setError('')
    setModal('create')
  }

  const openEdit = (u: any) => {
    setForm({ name: u.name, email: u.email, password: '', role: u.role, branch_id: u.branch_id ?? '' })
    setEditId(u.id)
    setError('')
    setModal('edit')
  }

  const closeModal = () => { setModal(null); setEditId(null) }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    setError('')
    saveMutation.mutate()
  }

  const set = (field: string, value: string) => setForm((p) => ({ ...p, [field]: value }))

  return (
    <div>
      <PageHeader
        title="Manajemen User"
        action={
          <Button variant="primary" onClick={openCreate}>+ Tambah User</Button>
        }
      />

      <div className="card overflow-hidden">
        {isLoading ? (
          <div className="flex justify-center py-16"><Spinner /></div>
        ) : (
          <div className="overflow-x-auto">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Nama</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Cabang</th>
                  <th>Status</th>
                  <th>Login Terakhir</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {(users?.length ?? 0) === 0 && (
                  <tr>
                    <td colSpan={7} className="!p-0">
                      <EmptyState message="Tidak ada user" />
                    </td>
                  </tr>
                )}
                {users?.map((u) => (
                  <tr key={u.id}>
                    <td>
                      <div className="flex items-center gap-2.5">
                        <div className="w-7 h-7 rounded-full flex items-center justify-center text-[11px] font-semibold text-white shrink-0"
                          style={{ background: 'var(--teal-m)' }}>
                          {u.name.charAt(0).toUpperCase()}
                        </div>
                        <span className="font-medium" style={{ color: 'var(--text)' }}>{u.name}</span>
                      </div>
                    </td>
                    <td style={{ color: 'var(--text2)' }}>{u.email}</td>
                    <td><Badge status={u.role} /></td>
                    <td style={{ color: 'var(--text2)' }}>
                      {branches?.find((b) => b.id === u.branch_id)?.name ?? (u.branch_id ? '—' : 'Semua Cabang')}
                    </td>
                    <td>
                      <span className={`badge ${u.is_active ? 'badge-active' : 'badge-readonly'}`}>
                        {u.is_active ? 'Aktif' : 'Nonaktif'}
                      </span>
                    </td>
                    <td style={{ color: 'var(--text3)' }}>{formatDate(u.last_login_at)}</td>
                    <td>
                      <div className="flex gap-1.5">
                        <button className="btn btn-ghost btn-sm" onClick={() => openEdit(u)}>Edit</button>
                        {u.is_active && (
                          <button
                            className="btn btn-ghost btn-sm"
                            style={{ color: 'var(--danger-t)' }}
                            onClick={() => {
                              if (confirm(`Nonaktifkan user ${u.name}?`)) deactivateMutation.mutate(u.id)
                            }}
                          >
                            Nonaktifkan
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {modal && (
        <div className="modal-overlay" onClick={(e) => { if (e.target === e.currentTarget) closeModal() }}>
          <div className="modal">
            <div className="modal-header">
              <h3>{modal === 'create' ? 'Tambah User' : 'Edit User'}</h3>
              <button className="modal-close" onClick={closeModal}>&times;</button>
            </div>
            <form onSubmit={handleSubmit} className="modal-body space-y-4">
              {error && <ErrorMessage message={error} />}

              <FormField label="Nama" required>
                <input className="form-input" value={form.name} onChange={(e) => set('name', e.target.value)} required />
              </FormField>

              <FormField label="Email" required>
                <input className="form-input" type="email" value={form.email} onChange={(e) => set('email', e.target.value)} required />
              </FormField>

              <FormField
                label={modal === 'create' ? 'Password' : 'Password Baru (kosongkan jika tidak diubah)'}
                required={modal === 'create'}
              >
                <input
                  className="form-input"
                  type="password"
                  value={form.password}
                  onChange={(e) => set('password', e.target.value)}
                  minLength={8}
                  required={modal === 'create'}
                  placeholder="Min. 8 karakter"
                />
              </FormField>

              <div className="grid grid-cols-2 gap-3">
                <FormField label="Role" required>
                  <select className="form-select" value={form.role} onChange={(e) => set('role', e.target.value)}>
                    <option value="write">Dokter (Write)</option>
                    <option value="readonly">Suster (Readonly)</option>
                    <option value="superadmin">Super Admin</option>
                  </select>
                </FormField>

                {form.role !== 'superadmin' && (
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
              </div>

              <div className="flex gap-3 pt-1">
                <Button type="submit" variant="primary" loading={saveMutation.isPending}>
                  {modal === 'create' ? 'Buat User' : 'Simpan'}
                </Button>
                <Button type="button" variant="secondary" onClick={closeModal}>Batal</Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
