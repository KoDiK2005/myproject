import { apiFetch, getAccessToken } from './auth'

const BASE_URL = 'http://localhost:8080'

export async function getPosts({ page = 1, limit = 10, search = '' } = {}) {
  const params = new URLSearchParams({ page, limit })
  if (search) params.append('search', search)

  // передаём токен если есть — бэк вернёт персональный фид
  const headers = {}
  const token = getAccessToken()
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${BASE_URL}/api/v1/posts?${params}`, { headers })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки постов')
  return data
}

export async function getPost(id) {
  const token = getAccessToken()
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  const res = await fetch(`${BASE_URL}/api/v1/posts/${id}`, { headers })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Пост не найден')
  return data
}

// visibility: 'public' | 'friends'
export async function createPost(title, body, visibility = 'public') {
  const res = await apiFetch(`${BASE_URL}/api/v1/posts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, body, visibility }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка создания поста')
  return data
}

export async function getPostsByUser(userId) {
  const token = getAccessToken()
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  const res = await fetch(`${BASE_URL}/api/v1/users/${userId}/posts`, { headers })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки постов')
  return data?.data ?? []
}

export async function deletePost(id) {
  const res = await apiFetch(`${BASE_URL}/api/v1/posts/${id}`, { method: 'DELETE' })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка удаления')
}
