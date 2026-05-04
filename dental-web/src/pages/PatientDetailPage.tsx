import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { getPatient, deletePatient } from '../api/patients'
import { getVisitsByPatient, deleteVisit } from '../api/visits'
import { getAttachmentsByPatient, deleteAttachment } from '../api/attachments'
import type { Attachment } from '../types'
import { useAuth } from '../hooks/useAuth'
import {
  Spinner, Button, Badge, EmptyState, PageHeader,
  formatDate, formatRupiah, formatGender,
} from '../components/ui'
import { getBranches } from '../api/branches'
import VisitFormModal from '../components/VisitFormModal'
import AttachmentUploadModal from '../components/AttachmentUploadModal'
import PreviewModal from '../components/PreviewModal'

export default function PatientDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const { isWrite } = useAuth()

  const [visitModal, setVisitModal] = useState<{ open: boolean; visitId?: string }>({ open: false })
  const [uploadModal, setUploadModal] = useState(false)
  const [previewAttachment, setPreviewAttachment] = useState<Attachment | null>(null)

  const { data: patient, isLoading } = useQuery({
    queryKey: ['patient', id],
    queryFn: () => getPatient(id!),
  })

  const { data: visits } = useQuery({
    queryKey: ['visits', id],
    queryFn: () => getVisitsByPatient(id!),
    enabled: !!id,
  })

  const { data: attachments } = useQuery({
    queryKey: ['attachments', id],
    queryFn: () => getAttachmentsByPatient(id!),
    enabled: !!id,
  })

  const { data: branches } = useQuery({ queryKey: ['branches'], queryFn: getBranches })

  const deleteMutation = useMutation({
    mutationFn: () => deletePatient(id!),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['patients'] })
      toast.success('Pasien berhasil dihapus')
      setTimeout(() => navigate('/patients'), 1000)
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message ?? 'Gagal menghapus pasien', { duration: 3000 })
    },
  })

  const deleteVisitMutation = useMutation({
    mutationFn: (visitId: string) => deleteVisit(visitId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['visits', id] })
      toast.success('Kunjungan berhasil dihapus')
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message ?? 'Gagal menghapus kunjungan', { duration: 3000 })
    },
  })

  const deleteAttachmentMutation = useMutation({
    mutationFn: (attId: string) => deleteAttachment(attId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['attachments', id] })
      toast.success('Lampiran berhasil dihapus')
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message ?? 'Gagal menghapus lampiran', { duration: 3000 })
    },
  })

  if (isLoading) {
    return <div className="flex justify-center py-20"><Spinner size="lg" /></div>
  }

  if (!patient) {
    return <div className="empty-state">Pasien tidak ditemukan</div>
  }

  const handleDeletePatient = () => {
    if (confirm(`Hapus pasien ${patient.name}? Tindakan ini tidak dapat dibatalkan.`)) {
      deleteMutation.mutate()
    }
  }

  return (
    <div className="space-y-4">
      <PageHeader
        title={patient.name}
        action={
          isWrite ? (
            <div className="flex gap-2">
              <Button variant="secondary" size="sm" onClick={() => navigate(`/patients/${id}/edit`)}>
                Edit
              </Button>
              <Button variant="danger" size="sm" onClick={handleDeletePatient} loading={deleteMutation.isPending}>
                Hapus
              </Button>
            </div>
          ) : undefined
        }
      />

      {/* Info pasien */}
      <div className="card" style={{ padding: '18px 22px' }}>
        <div className="card-header" style={{ margin: '-18px -22px 16px', padding: '14px 22px' }}>
          <h2>Info Pasien</h2>
          <Badge status={patient.status} />
        </div>
        <div className="grid grid-cols-2 gap-x-8 gap-y-4" style={{ gridTemplateColumns: 'repeat(3, 1fr)' }}>
          <InfoRow label="No. RM">
            <span className="mono font-semibold" style={{ color: 'var(--teal)' }}>{patient.no_rm}</span>
          </InfoRow>
          <InfoRow label="Tanggal Lahir">{formatDate(patient.date_of_birth)}</InfoRow>
          <InfoRow label="Jenis Kelamin">{formatGender(patient.gender)}</InfoRow>
          <InfoRow label="Telepon">{patient.phone}</InfoRow>
          <InfoRow label="NIK">{patient.nik ?? '—'}</InfoRow>
          <InfoRow label="Cabang">
            {branches?.find((b) => b.id === patient.branch_id)?.name ?? '—'}
          </InfoRow>
          <InfoRow label="Alamat">{patient.address ?? '—'}</InfoRow>
          <InfoRow label="Pekerjaan">{patient.occupation ?? '—'}</InfoRow>
          <InfoRow label="Catatan Alergi">
            {patient.allergy_notes ? (
              <span style={{ color: 'var(--amber-t)', fontWeight: 500 }}>{patient.allergy_notes}</span>
            ) : '—'}
          </InfoRow>
          <InfoRow label="Terdaftar">{formatDate(patient.created_at)}</InfoRow>
        </div>
      </div>

      {/* Riwayat Kunjungan */}
      <div className="card">
        <div className="card-header">
          <h2>Riwayat Kunjungan ({visits?.length ?? 0})</h2>
          {isWrite && (
            <Button size="sm" variant="primary" onClick={() => setVisitModal({ open: true })}>
              + Tambah Kunjungan
            </Button>
          )}
        </div>
        {(visits?.length ?? 0) === 0 ? (
          <EmptyState message="Belum ada kunjungan" />
        ) : (
          <div className="space-y-3" style={{ padding: '16px 20px' }}>
            {visits?.map((v) => (
              <div key={v.id} className="visit-item">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex flex-wrap items-center gap-2 mb-2">
                      <span className="mono text-[13px] font-semibold" style={{ color: 'var(--teal)' }}>
                        {formatDate(v.visit_date)}
                      </span>
                      {v.teeth_involved && (
                        <span className="badge badge-readonly" style={{ fontFamily: 'DM Mono, monospace' }}>
                          Gigi {v.teeth_involved}
                        </span>
                      )}
                      {v.next_control_date && (
                        <span className="badge badge-needs_control">
                          Kontrol {formatDate(v.next_control_date)}
                        </span>
                      )}
                    </div>
                    <p className="text-[13px]" style={{ color: 'var(--text)' }}>
                      <span style={{ fontWeight: 500 }}>Keluhan:</span> {v.chief_complaint}
                    </p>
                    {v.diagnosis && (
                      <p className="text-[13px] mt-1" style={{ color: 'var(--text2)' }}>
                        <span style={{ fontWeight: 500 }}>Diagnosis:</span> {v.diagnosis}
                      </p>
                    )}
                    {v.treatment && (
                      <p className="text-[13px] mt-1" style={{ color: 'var(--text2)' }}>
                        <span style={{ fontWeight: 500 }}>Tindakan:</span> {v.treatment}
                      </p>
                    )}
                    {v.notes && (
                      <p className="text-[13px] mt-1" style={{ color: 'var(--text3)', fontStyle: 'italic' }}>{v.notes}</p>
                    )}
                    <p className="text-[13px] font-semibold mt-2" style={{ color: 'var(--text)' }}>
                      {formatRupiah(v.cost)}
                    </p>
                  </div>
                  {isWrite && (
                    <div className="flex gap-1.5 shrink-0">
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={() => setVisitModal({ open: true, visitId: v.id })}
                      >
                        Edit
                      </button>
                      <button
                        className="btn btn-ghost btn-sm"
                        style={{ color: 'var(--danger-t)' }}
                        onClick={() => {
                          if (confirm('Hapus kunjungan ini?')) deleteVisitMutation.mutate(v.id)
                        }}
                      >
                        Hapus
                      </button>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Lampiran */}
      <div className="card">
        <div className="card-header">
          <h2>Lampiran ({attachments?.length ?? 0})</h2>
          {isWrite && (
            visits && visits.length > 0 ? (
              <Button size="sm" onClick={() => setUploadModal(true)}>+ Upload</Button>
            ) : (
              <span className="text-[12px]" style={{ color: 'var(--text3)', fontStyle: 'italic' }}>
                Tambah kunjungan dulu
              </span>
            )
          )}
        </div>
        {(attachments?.length ?? 0) === 0 ? (
          <EmptyState message="Belum ada lampiran" />
        ) : (
          <div style={{ borderTop: '1px solid var(--border)' }}>
            {attachments?.map((a) => (
              <div key={a.id}
                className="flex items-center justify-between"
                style={{ padding: '10px 20px', borderBottom: '1px solid var(--border)' }}
              >
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center text-base"
                    style={{ background: 'var(--teal-l)' }}>
                    {a.file_type === 'pdf' ? '📄' : '🖼️'}
                  </div>
                  <div>
                    <p className="text-[13px] font-medium" style={{ color: 'var(--text)' }}>{a.original_name}</p>
                    <p className="text-[11px]" style={{ color: 'var(--text3)' }}>
                      {(a.size_bytes / 1024).toFixed(0)} KB · {formatDate(a.created_at)}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <button
                    className="btn btn-ghost btn-sm"
                    style={{ color: 'var(--teal)' }}
                    onClick={() => setPreviewAttachment(a)}
                  >
                    Preview
                  </button>
                  {isWrite && (
                    <button
                      className="btn btn-ghost btn-sm"
                      style={{ color: 'var(--danger-t)' }}
                      onClick={() => {
                        if (confirm('Hapus lampiran ini?')) deleteAttachmentMutation.mutate(a.id)
                      }}
                    >
                      Hapus
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {visitModal.open && (
        <VisitFormModal
          patientId={id!}
          visitId={visitModal.visitId}
          onClose={() => setVisitModal({ open: false })}
          onSuccess={() => {
            qc.invalidateQueries({ queryKey: ['visits', id] })
            qc.invalidateQueries({ queryKey: ['patient', id] })
            setVisitModal({ open: false })
          }}
        />
      )}

      {uploadModal && (
        <AttachmentUploadModal
          visits={visits ?? []}
          onClose={() => setUploadModal(false)}
          onSuccess={() => {
            qc.invalidateQueries({ queryKey: ['attachments', id] })
            setUploadModal(false)
          }}
        />
      )}

      {previewAttachment && (
        <PreviewModal
          attachment={previewAttachment}
          onClose={() => setPreviewAttachment(null)}
        />
      )}
    </div>
  )
}

function InfoRow({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <p className="form-label" style={{ marginBottom: 3 }}>{label}</p>
      <p className="text-[13px]" style={{ color: 'var(--text)' }}>{children}</p>
    </div>
  )
}
