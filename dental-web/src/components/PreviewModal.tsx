import { useEffect, useState } from 'react'
import { fetchBlobUrl } from '../api/attachments'
import type { Attachment } from '../types'
import { Spinner } from './ui'

interface Props {
  attachment: Attachment
  onClose: () => void
}

export default function PreviewModal({ attachment, onClose }: Props) {
  const [blobUrl, setBlobUrl] = useState<string | null>(null)
  const [error, setError] = useState(false)

  useEffect(() => {
    let url: string
    fetchBlobUrl(attachment.id)
      .then((u) => {
        url = u
        if (attachment.file_type === 'pdf') {
          window.open(u, '_blank')
          onClose()
        } else {
          setBlobUrl(u)
        }
      })
      .catch(() => setError(true))

    return () => {
      if (url) URL.revokeObjectURL(url)
    }
  }, [attachment.id])

  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [onClose])

  if (attachment.file_type === 'pdf' && !error) return null

  return (
    <div
      style={{
        position: 'fixed', inset: 0, zIndex: 300,
        background: 'rgba(0,0,0,.85)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        padding: 16,
      }}
      onClick={(e) => { if (e.target === e.currentTarget) onClose() }}
    >
      <div style={{ position: 'relative', maxWidth: 900, maxHeight: '90vh', width: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <button
          onClick={onClose}
          style={{
            position: 'absolute', top: -40, right: 0,
            background: 'none', border: 'none', cursor: 'pointer',
            color: 'rgba(255,255,255,.7)', fontSize: 32, lineHeight: 1, padding: 0,
          }}
        >
          &times;
        </button>

        {error ? (
          <div className="card" style={{ padding: '24px 32px', textAlign: 'center', color: 'var(--danger-t)' }}>
            Gagal memuat file
          </div>
        ) : blobUrl ? (
          <img
            src={blobUrl}
            alt={attachment.original_name}
            style={{ maxWidth: '100%', maxHeight: '85vh', objectFit: 'contain', borderRadius: 8, boxShadow: '0 8px 40px rgba(0,0,0,.4)' }}
          />
        ) : (
          <Spinner size="lg" />
        )}

        {blobUrl && (
          <p style={{
            position: 'absolute', bottom: -28,
            left: 0, right: 0, textAlign: 'center',
            color: 'rgba(255,255,255,.5)', fontSize: 12,
            overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
          }}>
            {attachment.original_name}
          </p>
        )}
      </div>
    </div>
  )
}
