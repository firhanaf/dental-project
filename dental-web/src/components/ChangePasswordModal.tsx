import { useState, type FormEvent } from 'react'
import { toast } from 'sonner'
import { changePassword } from '../api/auth'
import { Button, FormField } from './ui'

interface Props {
  onClose: () => void
}

export default function ChangePasswordModal({ onClose }: Props) {
  const [form, setForm] = useState({ current: '', next: '', confirm: '' })
  const [loading, setLoading] = useState(false)

  const set = (field: string, value: string) => setForm((p) => ({ ...p, [field]: value }))

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (form.next !== form.confirm) {
      toast.error('Konfirmasi password tidak cocok')
      return
    }
    if (form.next.length < 8) {
      toast.error('Password baru minimal 8 karakter')
      return
    }
    setLoading(true)
    try {
      await changePassword(form.current, form.next)
      toast.success('Password berhasil diubah')
      onClose()
    } catch (err: any) {
      toast.error(err.response?.data?.message ?? 'Gagal mengubah password')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay">
      <div className="modal" style={{ maxWidth: 400 }}>
        <div className="modal-header">
          <h3>Ganti Password</h3>
          <button className="modal-close" onClick={onClose}>&times;</button>
        </div>
        <form onSubmit={handleSubmit} className="modal-body space-y-4">
          <FormField label="Password Saat Ini" required>
            <input
              className="form-input"
              type="password"
              value={form.current}
              onChange={(e) => set('current', e.target.value)}
              required
              autoFocus
            />
          </FormField>
          <FormField label="Password Baru" required>
            <input
              className="form-input"
              type="password"
              value={form.next}
              onChange={(e) => set('next', e.target.value)}
              minLength={8}
              placeholder="Min. 8 karakter"
              required
            />
          </FormField>
          <FormField label="Konfirmasi Password Baru" required>
            <input
              className="form-input"
              type="password"
              value={form.confirm}
              onChange={(e) => set('confirm', e.target.value)}
              minLength={8}
              placeholder="Ulangi password baru"
              required
            />
          </FormField>
          <div className="flex gap-3 pt-1">
            <Button type="submit" variant="primary" loading={loading}>Simpan</Button>
            <Button type="button" variant="secondary" onClick={onClose}>Batal</Button>
          </div>
        </form>
      </div>
    </div>
  )
}
