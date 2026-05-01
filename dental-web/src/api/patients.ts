import client from './client'
import type { Patient, PatientListRow, PaginatedResponse } from '../types'

export interface PatientListParams {
  page?: number
  limit?: number
  search?: string
  status?: string
  branch_id?: string
}

export interface CreatePatientPayload {
  name: string
  date_of_birth: string
  gender: string
  phone: string
  nik?: string
  address?: string
  occupation?: string
  allergy_notes?: string
  branch_id?: string
}

export type UpdatePatientPayload = Partial<CreatePatientPayload>

export const getPatients = (params: PatientListParams) =>
  client.get<PaginatedResponse<PatientListRow>>('/patients', { params }).then((r) => r.data)

export const getPatient = (id: string) =>
  client.get<Patient>(`/patients/${id}`).then((r) => r.data)

export const createPatient = (data: CreatePatientPayload) =>
  client.post<Patient>('/patients', data).then((r) => r.data)

export const updatePatient = (id: string, data: UpdatePatientPayload) =>
  client.put<Patient>(`/patients/${id}`, data).then((r) => r.data)

export const deletePatient = (id: string) =>
  client.delete(`/patients/${id}`)
