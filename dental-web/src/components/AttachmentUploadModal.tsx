import { useState, type FormEvent } from 'react'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { uploadAttachment } from '../api/attachments'
import type { Visit } from '../types'
import { Button, FormField } from './ui'

interface Props {
  visits: Visit[]
  onClose: () => void
  onSuccess: () => void
}

export default function AttachmentUploadModal({ visits, onClose, onSuccess }: Props) {
  const [visitId, setVisitId] = useState(visits[0]?.id ?? '')
  const [file, setFile] = useState<File | null>(null)

  const mutation = useMutation({
    mutationFn: () => uploadAttachment(visitId, file!),
    onSuccess: () => {
      toast.success('Lampiran berhasil diupload')
      onSuccess()
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message ?? 'Gagal upload file', { duration: 3000 })
    },
  })

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (!file) { toast.error('Pilih file terlebih dahulu'); return }
    if (!visitId) { toast.error('Pilih kunjungan terlebih dahulu'); return }
    mutation.mutate()
  }

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h3>Upload Lampiran</h3>
          <button className="modal-close" onClick={onClose}>&times;</button>
        </div>

        <form onSubmit={handleSubmit} className="modal-body space-y-4">
          <FormField label="Kunjungan" required>
            <select
              className="form-select"
              value={visitId}
              onChange={(e) => setVisitId(e.target.value)}
              required
            >
              {visits.length === 0 && <option value="">Belum ada kunjungan</option>}
              {visits.map((v) => (
                <option key={v.id} value={v.id}>
                  {new Date(v.visit_date).toLocaleDateString('id-ID')} — {v.chief_complaint.slice(0, 40)}
                </option>
              ))}
            </select>
          </FormField>

          <FormField label="File" required>
            <div style={{
              border: '1px dashed var(--border2)',
              borderRadius: 'var(--radius)',
              padding: '14px 16px',
              background: 'var(--bg)',
            }}>
              <input
                type="file"
                accept=".pdf,image/jpeg,image/png,image/webp"
                style={{ fontSize: 13, color: 'var(--text2)', width: '100%' }}
                onChange={(e) => {
                  const f = e.target.files?.[0] ?? null
                  if (f && f.size > 20 * 1024 * 1024) {
                    toast.error('Ukuran file maksimal 20 MB', { duration: 3000 })
                    e.target.value = ''
                    setFile(null)
                    return
                  }
                  setFile(f)
                }}
                required
              />
              <p className="text-[11px] mt-2" style={{ color: 'var(--text3)' }}>
                PDF, JPEG, PNG, WebP · Maks. 20 MB
              </p>
            </div>
          </FormField>

          <div className="flex gap-3 pt-1">
            <Button type="submit" variant="primary" loading={mutation.isPending}>Upload</Button>
            <Button type="button" variant="secondary" onClick={onClose}>Batal</Button>
          </div>
        </form>
      </div>
    </div>
  )
}
