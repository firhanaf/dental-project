import client from './client'
import type { LoginResponse, User } from '../types'

export const login = (email: string, password: string) =>
  client.post<LoginResponse>('/auth/login', { email, password }).then((r) => r.data)

export const me = () =>
  client.get<User>('/auth/me').then((r) => r.data)

export const logout = () =>
  client.post('/auth/logout')

export const changePassword = (currentPassword: string, newPassword: string) =>
  client.put('/auth/change-password', { current_password: currentPassword, new_password: newPassword })

export const resetPassword = (email: string, token: string, newPassword: string) =>
  client.post('/auth/reset-password', { email, token, new_password: newPassword })
