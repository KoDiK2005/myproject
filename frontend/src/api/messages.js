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

// Единственное WS-соединение на вкладку (бэкенд держит один Client на userID в hub —
// второе соединение от того же юзера вытеснит первое). Сообщения/уведомления раздаются
// всем подписчикам, поэтому чат и колокольчик уведомлений могут слушать один и тот же сокет.
let ws = null
let reconnectDelay = 1000
const maxReconnectDelay = 30000
const messageListeners = new Set()
const notificationListeners = new Set()

function openWS() {
  const token = getAccessToken()
  if (!token) return

  ws = new WebSocket(`${WS_URL}?token=${token}`)

  ws.onopen = () => {
    reconnectDelay = 1000
  }

  ws.onmessage = (e) => {
    try {
      const data = JSON.parse(e.data)
      if (data.type === 'message' && data.payload) {
        messageListeners.forEach(fn => fn(data.payload))
      } else if (data.type === 'notification' && data.payload) {
        notificationListeners.forEach(fn => fn(data.payload))
      }
    } catch {}
  }

  ws.onerror = () => {} // тихо, обработку делает onclose

  ws.onclose = () => {
    if (messageListeners.size === 0 && notificationListeners.size === 0) return
    setTimeout(openWS, reconnectDelay)
    reconnectDelay = Math.min(reconnectDelay * 2, maxReconnectDelay)
  }
}

function ensureWSConnected() {
  if (!ws || ws.readyState === WebSocket.CLOSED) openWS()
}

// подписаться на входящие сообщения, возвращает функцию отписки
export function subscribeMessages(onMessage) {
  messageListeners.add(onMessage)
  ensureWSConnected()
  return () => messageListeners.delete(onMessage)
}

// подписаться на уведомления (лайки/комментарии/заявки в друзья), возвращает функцию отписки
export function subscribeNotifications(onNotification) {
  notificationListeners.add(onNotification)
  ensureWSConnected()
  return () => notificationListeners.delete(onNotification)
}
