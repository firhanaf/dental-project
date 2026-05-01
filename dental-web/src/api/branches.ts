import client from './client'
import type { Branch } from '../types'

export const getBranches = () =>
  client.get<Branch[]>('/branches').then((r) => r.data)
