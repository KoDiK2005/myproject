import { apiFetch } from './auth'

const BASE_URL = 'http://localhost:8080'

export async function getLikeCount(postId) {
  const res = await fetch(`${BASE_URL}/api/v1/posts/${postId}/likes`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data // { count: N }
}

export async function likePost(postId) {
  const res = await apiFetch(`${BASE_URL}/api/v1/posts/${postId}/like`, {
    method: 'POST',
  })
  if (res.status === 204 || res.status === 200) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}

export async function unlikePost(postId) {
  const res = await apiFetch(`${BASE_URL}/api/v1/posts/${postId}/like`, {
    method: 'DELETE',
  })
  if (res.status === 204 || res.status === 200) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}
