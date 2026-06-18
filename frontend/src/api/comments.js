import { getAccessToken } from './auth'

const BASE_URL = 'http://localhost:8080'

function authHeaders(extra = {}) {
  return { Authorization: `Bearer ${getAccessToken()}`, ...extra }
}

export async function getComments(postId) {
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/comments`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки комментариев')
  return data
}

export async function createComment(postId, body) {
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/comments`, {
    method: 'POST',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ body }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка добавления комментария')
  return data
}

export async function deleteComment(commentId) {
  const res = await fetch(`${BASE_URL}/api/v1/comments/${commentId}`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка удаления комментария')
}
