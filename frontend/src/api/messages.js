import { apiFetch, getAccessToken } from './auth'

const BASE = 'http://localhost:8080/api/v1'
const WS_URL = 'ws://localhost:8080/ws'

// список всех переписок
export async function getConversations() {
  const res = await apiFetch(`${BASE}/messages`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data ?? []
}

// история сообщений с конкретным юзером
export async function getHistory(userId, page = 1, limit = 50) {
  const params = new URLSearchParams({ page, limit })
  const res = await apiFetch(`${BASE}/messages/${userId}?${params}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data ?? []
}

// отправить сообщение
export async function sendMessage(userId, content) {
  const res = await apiFetch(`${BASE}/messages/${userId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data
}

// WebSocket соединение — передаём токен в query param (браузер не даёт header для WS)
export function connectWS(onMessage) {
  const token = getAccessToken()
  if (!token) return null

  const ws = new WebSocket(`${WS_URL}?token=${token}`)

  ws.onmessage = (e) => {
    try {
      const data = JSON.parse(e.data)
      if (data.type === 'message' && data.payload) {
        onMessage(data.payload)
      }
    } catch {}
  }

  ws.onerror = () => {} // тихо
  return ws
}
