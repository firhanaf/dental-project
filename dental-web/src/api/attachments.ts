import client from './client'
import type { Attachment } from '../types'

export const getAttachmentsByPatient = (patientId: string) =>
  client.get<Attachment[]>(`/patients/${patientId}/attachments`).then((r) => r.data)

export const uploadAttachment = (visitId: string, file: File) => {
  const form = new FormData()
  form.append('visit_id', visitId)
  form.append('file', file)
  return client
    .post<Attachment>('/attachments', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    .then((r) => r.data)
}

export const deleteAttachment = (id: string) =>
  client.delete(`/attachments/${id}`)

export const getDownloadUrl = (id: string) =>
  `/api/v1/attachments/${id}/download`

export const fetchBlobUrl = async (id: string): Promise<string> => {
  const res = await client.get(`/attachments/${id}/download`, { responseType: 'blob' })
  return URL.createObjectURL(res.data)
}
