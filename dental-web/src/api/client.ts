import axios from 'axios'

const client = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

client.interceptors.response.use(
  (res) => {
    // Unwrap BE envelope: { success, data, meta? }
    // Paginated responses keep { data, meta } shape (matches PaginatedResponse<T>)
    // All other responses become the inner data directly
    if (res.data && typeof res.data === 'object' && 'success' in res.data) {
      if (res.data.meta) {
        res.data = { data: res.data.data, meta: res.data.meta }
      } else {
        res.data = res.data.data
      }
    }
    return res
  },
  (err) => {
    // Normalize error: { success, error: { code, message } } → { message }
    const msg = err.response?.data?.error?.message
    if (msg) err.response.data = { message: msg }

    // Jangan redirect jika 401 datang dari endpoint login itu sendiri
    // (artinya credentials salah, bukan token expired)
    const isLoginEndpoint = err.config?.url?.includes('/auth/login')
    if (err.response?.status === 401 && !isLoginEndpoint) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export default client
