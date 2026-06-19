import { apiFetch } from './auth'

const BASE = 'http://localhost:8080/api/v1'

export async function sendFriendRequest(userId) {
  const res = await apiFetch(`${BASE}/friends/request/${userId}`, { method: 'POST' })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data
}

export async function acceptFriendRequest(userId) {
  const res = await apiFetch(`${BASE}/friends/accept/${userId}`, { method: 'POST' })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data
}

export async function rejectFriendRequest(userId) {
  const res = await apiFetch(`${BASE}/friends/reject/${userId}`, { method: 'POST' })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data
}

export async function removeFriend(userId) {
  const res = await apiFetch(`${BASE}/friends/${userId}`, { method: 'DELETE' })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
}

export async function getFriends() {
  const res = await apiFetch(`${BASE}/friends`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data ?? []
}

export async function getIncomingRequests() {
  const res = await apiFetch(`${BASE}/friends/requests/incoming`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data ?? []
}

export async function getOutgoingRequests() {
  const res = await apiFetch(`${BASE}/friends/requests/outgoing`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data ?? []
}

export async function getFriendStatus(userId) {
  const res = await apiFetch(`${BASE}/friends/status/${userId}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data // { status: 'none' | 'pending_sent' | 'pending_received' | 'accepted' }
}
