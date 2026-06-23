import { apiFetch, getAccessToken } from './auth'

const BASE_URL = 'http://localhost:8080'

export async function getComments(postId) {
  const token = getAccessToken()
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/comments`, { headers })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки комментариев')
  return data
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
