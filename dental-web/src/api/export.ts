import client from './client'

export interface ExportParams {
  branch_id?: string
  date_from?: string
  date_to?: string
}

export const exportPatients = async (params: ExportParams) => {
  const res = await client.get('/export/patients', {
    params,
    responseType: 'blob',
  })
  const url = URL.createObjectURL(res.data)
  const a = document.createElement('a')
  a.href = url
  a.download = `export_pasien_${new Date().toISOString().slice(0, 10)}.xlsx`
  a.click()
  URL.revokeObjectURL(url)
}

export const exportVisits = async (params: ExportParams) => {
  const res = await client.get('/export/visits', {
    params,
    responseType: 'blob',
  })
  const url = URL.createObjectURL(res.data)
  const a = document.createElement('a')
  a.href = url
  a.download = `export_kunjungan_${new Date().toISOString().slice(0, 10)}.xlsx`
  a.click()
  URL.revokeObjectURL(url)
}
