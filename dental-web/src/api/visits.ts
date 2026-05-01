import client from './client'
import type { Visit } from '../types'

export interface CreateVisitPayload {
  patient_id: string
  doctor_id: string
  visit_date: string
  chief_complaint: string
  diagnosis?: string
  treatment?: string
  teeth_involved?: string
  cost: number
  next_control_date?: string
  notes?: string
}

export type UpdateVisitPayload = Partial<CreateVisitPayload>

export const getVisitsByPatient = (patientId: string) =>
  client.get<Visit[]>(`/patients/${patientId}/visits`).then((r) => r.data)

export const getVisit = (id: string) =>
  client.get<Visit>(`/visits/${id}`).then((r) => r.data)

export const createVisit = (data: CreateVisitPayload) =>
  client.post<Visit>('/visits', data).then((r) => r.data)

export const updateVisit = (id: string, data: UpdateVisitPayload) =>
  client.put<Visit>(`/visits/${id}`, data).then((r) => r.data)

export const deleteVisit = (id: string) =>
  client.delete(`/visits/${id}`)
