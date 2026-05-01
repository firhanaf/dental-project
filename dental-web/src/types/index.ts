export type UserRole = 'superadmin' | 'write' | 'readonly'
export type Gender = 'male' | 'female'
export type PatientStatus = 'new' | 'active' | 'needs_control'
export type FileType = 'pdf' | 'image'

export interface Branch {
  id: string
  name: string
  address: string | null
  phone: string | null
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface User {
  id: string
  branch_id: string | null
  name: string
  email: string
  role: UserRole
  is_active: boolean
  last_login_at: string | null
  created_at: string
  updated_at: string
}

export interface Patient {
  id: string
  branch_id: string
  created_by: string
  no_rm: string
  name: string
  nik: string | null
  date_of_birth: string
  gender: Gender
  phone: string
  address: string | null
  occupation: string | null
  allergy_notes: string | null
  status: PatientStatus
  deleted_at: string | null
  created_at: string
  updated_at: string
}

export interface PatientListRow extends Patient {
  branch_name: string
  age: number
  last_visit_date: string | null
  last_diagnosis: string | null
  last_doctor: string | null
  total_visits: number
  total_cost: number
}

export interface Visit {
  id: string
  patient_id: string
  branch_id: string
  doctor_id: string
  created_by: string
  visit_date: string
  chief_complaint: string
  diagnosis: string | null
  treatment: string | null
  teeth_involved: string | null
  cost: number
  next_control_date: string | null
  notes: string | null
  deleted_at: string | null
  created_at: string
  updated_at: string
}

export interface Attachment {
  id: string
  visit_id: string
  uploaded_by: string
  original_name: string
  file_type: FileType
  mime_type: string
  size_bytes: number
  deleted_at: string | null
  created_at: string
}

export interface AuthClaims {
  user_id: string
  role: UserRole
  branch_id: string | null
  name: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface Meta {
  page: number
  limit: number
  total: number
  has_next: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  meta: Meta
}

export interface ApiError {
  code: string
  message: string
}
