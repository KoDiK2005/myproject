import { apiFetch } from './auth'

const BASE = 'http://localhost:8080/api/v1'

export async function blockUser(userId) {
  const res = await apiFetch(`${BASE}/blocks/${userId}`, { method: 'POST' })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}

export async function unblockUser(userId) {
  const res = await apiFetch(`${BASE}/blocks/${userId}`, { method: 'DELETE' })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}

export async function getBlockedUsers() {
  const res = await apiFetch(`${BASE}/blocks`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data ?? []
}
