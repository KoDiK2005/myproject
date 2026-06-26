import { useState, useEffect, useRef, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  getNotifications,
  getUnreadCount,
  markNotificationRead,
  markAllNotificationsRead,
} from '../api/notifications'
import { subscribeNotifications } from '../api/messages'
import { getAccessToken } from '../api/auth'

const LABELS = {
  friend_request: actor => `${actor} хочет дружить с тобой`,
  friend_accept: actor => `${actor} принял(а) твою заявку в друзья`,
  like: actor => `${actor} лайкнул(а) твой пост`,
  comment: actor => `${actor} прокомментировал(а) твой пост`,
}

export default function NotificationBell() {
  const navigate = useNavigate()
  const [open, setOpen] = useState(false)
  const [unread, setUnread] = useState(0)
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(false)
  const rootRef = useRef(null)

  useEffect(() => {
    if (!getAccessToken()) return
    getUnreadCount().then(setUnread).catch(() => {})

    const unsubscribe = subscribeNotifications((n) => {
      setUnread(prev => prev + 1)
      setItems(prev => [n, ...prev])
    })
    return unsubscribe
  }, [])

  // закрытие по клику снаружи
  useEffect(() => {
    function onClickOutside(e) {
      if (rootRef.current && !rootRef.current.contains(e.target)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', onClickOutside)
    return () => document.removeEventListener('mousedown', onClickOutside)
  }, [])

  const toggle = useCallback(async () => {
    const next = !open
    setOpen(next)
    if (!next) return

    setLoading(true)
    try {
      const resp = await getNotifications()
      setItems(resp.data ?? [])
      if (unread > 0) {
        await markAllNotificationsRead()
        setUnread(0)
      }
    } catch {
      // тихо
    } finally {
      setLoading(false)
    }
  }, [open, unread])

  async function handleClick(n) {
    setOpen(false)
    if (!n.read) {
      try { await markNotificationRead(n.id) } catch {}
    }
    if (n.type === 'like' || n.type === 'comment') {
      if (n.post_id) navigate(`/posts/${n.post_id}`)
    } else if (n.type === 'friend_request') {
      navigate('/friends')
    } else if (n.type === 'friend_accept') {
      navigate(`/users/${n.actor_id}`)
    }
  }

  if (!getAccessToken()) return null

  return (
    <div className="notification-bell" ref={rootRef}>
      <button className="notification-bell-btn" onClick={toggle} title="Уведомления">
        🔔{unread > 0 && <span className="navbar-badge">{unread}</span>}
      </button>

      {open && (
        <div className="notification-dropdown">
          {loading ? (
            <p className="notification-empty">Загружаем...</p>
          ) : items.length === 0 ? (
            <p className="notification-empty">Уведомлений пока нет</p>
          ) : (
            items.map(n => (
              <button
                key={n.id}
                className={`notification-item ${n.read ? '' : 'unread'}`}
                onClick={() => handleClick(n)}
              >
                {LABELS[n.type]?.(n.actor_name) ?? n.type}
                <span className="notification-item-date">
                  {new Date(n.created_at).toLocaleString('ru-RU')}
                </span>
              </button>
            ))
          )}
        </div>
      )}
    </div>
  )
}
