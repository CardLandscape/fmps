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
export const deleteMemberWithAuth = (id, authPassword) =>
  api.post(`/members/${id}/delete`, { auth_password: authPassword })

export const getRules = () => api.get('/rules')
export const createRule = data => api.post('/rules', data)
export const updateRule = (id, data) => api.put(`/rules/${id}`, data)
export const deleteRule = id => api.delete(`/rules/${id}`)

export const getRecords = params => api.get('/records', { params })
export const createRecord = data => api.post('/records', data)
export const deleteRecord = id => api.delete(`/records/${id}`)

export const getSettings = () => api.get('/settings')
export const updateSettings = data => api.put('/settings', data)

export const getCases = () => api.get('/cases')
export const getCase = id => api.get(`/cases/${id}`)
export const createCase = data => api.post('/cases', data)
export const updateCase = (id, data) => api.put(`/cases/${id}`, data)
export const deleteCase = id => api.delete(`/cases/${id}`)
export const parseCaseTxt = data => api.post('/cases/parse-txt', data)
export const startPunishment = id => api.post(`/cases/${id}/start`)
export const completePunishment = id => api.post(`/cases/${id}/complete`)
export const completeStep = id => api.post(`/cases/${id}/complete-step`)
export const addPenalty = (caseId, data) => api.post(`/cases/${caseId}/penalty`, data)
export const revokePenalty = (penaltyId, data) => api.post(`/penalty/${penaltyId}/revoke`, data)

export default api
