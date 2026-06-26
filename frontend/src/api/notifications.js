import { apiFetch } from './auth'

const BASE = 'http://localhost:8080/api/v1'

export async function getNotifications(page = 1, limit = 20) {
  const params = new URLSearchParams({ page, limit })
  const res = await apiFetch(`${BASE}/notifications?${params}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data
}

export async function getUnreadCount() {
  const res = await apiFetch(`${BASE}/notifications/unread-count`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data.count ?? 0
}

export async function markNotificationRead(id) {
  const res = await apiFetch(`${BASE}/notifications/${id}/read`, { method: 'POST' })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}

export async function markAllNotificationsRead() {
  const res = await apiFetch(`${BASE}/notifications/read-all`, { method: 'POST' })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}
