import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' }
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export const authApi = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  me: () => api.get('/auth/me'),
  changePassword: (data) => api.put('/auth/change-password', data),
}

export const adminApi = {
  getUsers: () => api.get('/admin/users'),
  createUser: (data) => api.post('/admin/users', data),
  updateRole: (id, role) => api.put(`/admin/users/${id}/role`, { role }),
  deleteUser: (id) => api.delete(`/admin/users/${id}`),
  updateUserPassword: (id, password) => api.put(`/admin/users/${id}/password`, { password }),
}

export const tableApi = {
  // Public — คืน { tables, floor_items }
  getFloorPlan: (date) => api.get(`/tables${date ? '?date=' + date : ''}`),

  // User reservations
  uploadSlip: (file) => {
    const form = new FormData()
    form.append('slip', file)
    return api.post('/reservations/upload-slip', form, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
  createReservation: (data) => api.post('/reservations', data),
  getMyReservations: () => api.get('/reservations'),
  cancelReservation: (id) => api.delete(`/reservations/${id}`),

  // Admin — tables
  adminGetTables: () => api.get('/admin/tables'),
  adminCreateTable: (data) => api.post('/admin/tables', data),
  adminUpdateTable: (id, data) => api.put(`/admin/tables/${id}`, data),
  adminDeleteTable: (id) => api.delete(`/admin/tables/${id}`),

  // Admin — floor items
  adminGetFloorItems: () => api.get('/admin/floor-items'),
  adminCreateFloorItem: (data) => api.post('/admin/floor-items', data),
  adminUpdateFloorItem: (id, data) => api.put(`/admin/floor-items/${id}`, data),
  adminDeleteFloorItem: (id) => api.delete(`/admin/floor-items/${id}`),

  // Admin — reservations
  adminGetReservations: (date, status) => {
    const p = new URLSearchParams()
    if (date) p.append('date', date)
    if (status) p.append('status', status)
    return api.get(`/admin/reservations?${p}`)
  },
  adminGetSummary: () => api.get('/admin/reservations/summary'),
  adminUpdateStatus: (id, status, adminNote) =>
    api.put(`/admin/reservations/${id}/status`, { status, admin_note: adminNote }),
  adminDeleteReservation: (id) => api.delete(`/admin/reservations/${id}`),
  adminDeleteReservationsBulk: (ids) => api.post('/admin/reservations/delete-bulk', { ids }),
}
