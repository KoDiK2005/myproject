import { apiFetch, getAccessToken } from './auth'

const BASE_URL = 'http://localhost:8080'

export async function getComments(postId, page = 1, limit = 20) {
  const token = getAccessToken()
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  const params = new URLSearchParams({ page, limit })
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/comments?${params}`, { headers })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки комментариев')
  return data // { data, page, total, limit }
}

export async function createComment(postId, body) {
  const res = await apiFetch(`${BASE_URL}/api/v1/posts/${postId}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ body }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка добавления комментария')
  return data
}

export async function deleteComment(commentId) {
  const res = await apiFetch(`${BASE_URL}/api/v1/comments/${commentId}`, {
    method: 'DELETE',
  })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка удаления комментария')
}
