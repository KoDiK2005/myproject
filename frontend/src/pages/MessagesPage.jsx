import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getConversations } from '../api/messages'
import { useToast, ToastContainer } from '../components/Toast'
import Navbar from '../components/Navbar'

const BASE_IMG = 'http://localhost:8080'

export default function MessagesPage() {
  const navigate = useNavigate()
  const { toasts, showToast } = useToast()
  const [convs, setConvs] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      try {
        const data = await getConversations()
        setConvs(data)
      } catch (err) {
        showToast(err.message)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  function formatTime(iso) {
    const d = new Date(iso)
    const now = new Date()
    if (d.toDateString() === now.toDateString()) {
      return d.toLocaleTimeString('ru', { hour: '2-digit', minute: '2-digit' })
    }
    return d.toLocaleDateString('ru', { day: 'numeric', month: 'short' })
  }

  return (
    <div className="messages-page">
      <ToastContainer toasts={toasts} />

      <Navbar onBack={() => navigate('/')} title="Сообщения" />

      <div className="convs-list">
        {loading && <p className="feed-status">Загружаем...</p>}
        {!loading && convs.length === 0 && (
          <p className="feed-status">Нет переписок. Напиши кому-нибудь!</p>
        )}
        {convs.map(conv => (
          <Link
            key={conv.user_id}
            to={`/messages/${conv.user_id}`}
            className="conv-card"
          >
            <div className="conv-avatar">
              {conv.user_avatar
                ? <img src={`${BASE_IMG}${conv.user_avatar}`} alt={conv.user_name} />
                : <span className="avatar-placeholder">{conv.user_name?.[0]?.toUpperCase()}</span>
              }
              {conv.unread > 0 && <span className="conv-unread">{conv.unread}</span>}
            </div>
            <div className="conv-info">
              <div className="conv-name">{conv.user_name}</div>
              <div className="conv-last">{conv.last_message}</div>
            </div>
            <div className="conv-time">{formatTime(conv.last_at)}</div>
          </Link>
        ))}
      </div>
    </div>
  )
}
