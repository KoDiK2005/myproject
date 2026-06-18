import { getAccessToken } from './auth'

const BASE_URL = 'http://localhost:8080'

function authHeaders() {
  return { Authorization: `Bearer ${getAccessToken()}` }
}

export async function getLikeCount(postId) {
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/likes`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data // { count: N }
}

export async function likePost(postId) {
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/like`, {
    method: 'POST',
    headers: authHeaders(),
  })
  if (res.status === 204 || res.status === 200) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}

export async function unlikePost(postId) {
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/like`, {
    method: 'DELETE',
    headers: authHeaders(),
  })
  if (res.status === 204 || res.status === 200) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}
