import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 自动添加 token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器 - 统一处理响应格式
api.interceptors.response.use(
  response => {
    // 如果响应包含 code 字段，检查是否成功
    if (response.data && response.data.code !== undefined) {
      if (response.data.code !== 0) {
        return Promise.reject({ 
          message: response.data.message || 'Request failed',
          code: response.data.code 
        })
      }
      // 返回 data 字段的内容
      return response.data.data || response.data
    }
    return response
  },
  error => {
    // 处理 401 错误
    if (error.response?.status === 401 || error.code === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
