import { getAccessToken } from './auth'

const BASE_URL = 'http://localhost:8080'

function authHeaders(extra = {}) {
  return {
    Authorization: `Bearer ${getAccessToken()}`,
    ...extra,
  }
}

export async function getPosts({ page = 1, limit = 10, search = '' } = {}) {
  const params = new URLSearchParams({ page, limit })
  if (search) params.append('search', search)

  const res = await fetch(`${BASE_URL}/api/v1/posts?${params}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки постов')
  return data // { data: [], total, page, limit }
}

export async function createPost(title, body) {
  const res = await fetch(`${BASE_URL}/api/v1/posts`, {
    method: 'POST',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ title, body }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка создания поста')
  return data
}

export async function deletePost(id) {
  const res = await fetch(`${BASE_URL}/api/v1/posts/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка удаления')
}
