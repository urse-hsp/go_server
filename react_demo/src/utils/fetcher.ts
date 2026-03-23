// utils/fetcher.ts
const NEXT_PUBLIC_API_BASE_URL = ''

type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

interface FetchOptions extends RequestInit {
  method?: RequestMethod
  body?: any // 自动 JSON 序列化
}

const DEFAULT_HEADERS = {
  'Content-Type': 'application/json',
}

// 获取认证 Token（可根据项目调整）
const getAuthToken = (): string | null => {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('token')
}

// 核心 fetch 封装
const request = async <T>(
  url: string,
  options: FetchOptions = {},
): Promise<T> => {
  const { method = 'GET', body, headers, ...rest } = options

  // 合并 headers
  const finalHeaders: any = {
    ...DEFAULT_HEADERS,
    ...headers,
  }

  const token = getAuthToken()
  if (token) {
    finalHeaders['Authorization'] = `Bearer ${token}`
  }

  // 如果是 GET/DELETE，不应有 body
  const hasBody = method !== 'GET' && method !== 'DELETE'

  const config: RequestInit = {
    method,
    headers: finalHeaders,
    ...rest,
  }

  if (hasBody && body !== undefined) {
    config.body = typeof body === 'string' ? body : JSON.stringify(body)
  }

  const fullUrl = `${NEXT_PUBLIC_API_BASE_URL || ''}${url}`

  try {
    const res = await fetch(fullUrl, config)

    if (!res.ok) {
      // 统一错误格式
      let errorMessage = `HTTP Error: ${res.status}`
      try {
        const errorData = await res.json()
        errorMessage = errorData.message || errorMessage
      } catch {
        // 忽略 JSON 解析失败
      }
      throw new Error(errorMessage)
    }

    // 如果响应体为空（如 204 No Content），返回 null
    const contentType = res.headers.get('content-type')
    if (!contentType || !contentType.includes('application/json')) {
      return null as T
    }

    const data: T = await res.json()
    return data
  } catch (err) {
    console.error('Fetch request failed:', err)
    throw err // 抛出错误供 SWR 捕获
  }
}

// SWR 默认使用 GET，所以提供一个只用于 GET 的 fetcher
export const fetcher = <T>(url: string): Promise<T> => {
  return request<T>(url, { method: 'GET' })
}

// 导出通用 request，可用于 POST/PUT 等（配合 mutate 使用）
export default request
