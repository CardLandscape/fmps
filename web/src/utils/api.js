import axios from 'axios'
import router from '@/router'

const api = axios.create({
  baseURL: '/api'
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('fmps_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('fmps_token')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export const login = (username, password) =>
  api.post('/login', { username, password })

export const getStats = () => api.get('/stats')

export const getMembers = () => api.get('/members')
export const createMember = data => api.post('/members', data)
export const updateMember = (id, data) => api.put(`/members/${id}`, data)
export const deleteMember = id => api.delete(`/members/${id}`)

export const getRules = () => api.get('/rules')
export const createRule = data => api.post('/rules', data)
export const updateRule = (id, data) => api.put(`/rules/${id}`, data)
export const deleteRule = id => api.delete(`/rules/${id}`)

export const getRecords = params => api.get('/records', { params })
export const createRecord = data => api.post('/records', data)
export const deleteRecord = id => api.delete(`/records/${id}`)

export const getSettings = () => api.get('/settings')
export const updateSettings = data => api.put('/settings', data)

export default api
