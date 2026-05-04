import client from './client'
import type { User } from '../types'

export interface CreateUserPayload {
  name: string
  email: string
  password: string
  role: string
  branch_id?: string | null
}

export interface UpdateUserPayload {
  name?: string
  email?: string
  password?: string
  role?: string
  branch_id?: string | null
  is_active?: boolean
}

export const getUsers = () =>
  client.get<User[]>('/users').then((r) => r.data)

export const createUser = (data: CreateUserPayload) =>
  client.post<User>('/users', data).then((r) => r.data)

export const updateUser = (id: string, data: UpdateUserPayload) =>
  client.put<User>(`/users/${id}`, data).then((r) => r.data)

export const activateUser = (id: string) =>
  client.post(`/users/${id}/activate`)

export const deactivateUser = (id: string) =>
  client.post(`/users/${id}/deactivate`)

export const deleteUser = (id: string) =>
  client.delete(`/users/${id}`)
