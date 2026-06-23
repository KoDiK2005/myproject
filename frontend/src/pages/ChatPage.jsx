import { useState, useEffect, useRef, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getHistory, sendMessage, connectWS } from '../api/messages'
import { getUser } from '../api/users'
import { getCurrentUserId } from '../api/auth'
import { useToast, ToastContainer } from '../components/Toast'

export default function ChatPage() {
  const { id: partnerId } = useParams()
  const navigate = useNavigate()
  const { toasts, showToast } = useToast()
  const myId = getCurrentUserId()

  const [partner, setPartner] = useState(null)
  const [messages, setMessages] = useState([])
  const [text, setText] = useState('')
  const [sending, setSending] = useState(false)
  const [loading, setLoading] = useState(true)

  const bottomRef = useRef(null)
  const wsRef = useRef(null)
  const inputRef = useRef(null)

  // загружаем историю и данные собеседника
  useEffect(() => {
    async function load() {
      try {
        const [hist, user] = await Promise.all([
          getHistory(partnerId),
          getUser(partnerId),
        ])
        setMessages(hist)
        setPartner(user)
      } catch (err) {
        showToast(err.message)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [partnerId])

  // WebSocket — подключаемся и получаем входящие сообщения реалтайм
  useEffect(() => {
    const ws = connectWS((msg) => {
      // принимаем только сообщения от нашего собеседника
      if (msg.sender_id === parseInt(partnerId)) {
        setMessages(prev => [...prev, msg])
      }
    })
    wsRef.current = ws

    return () => {
      ws?.close()
    }
  }, [partnerId])

  // скролл вниз при новых сообщениях
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = useCallback(async (e) => {
    e?.preventDefault()
    const content = text.trim()
    if (!content || sending) return

    setSending(true)
    setText('')
    try {
      const msg = await sendMessage(partnerId, content)
      setMessages(prev => [...prev, msg])
    } catch (err) {
      showToast(err.message)
      setText(content) // возвращаем текст если не отправилось
    } finally {
      setSending(false)
      inputRef.current?.focus()
    }
  }, [text, sending, partnerId])

  // отправка по Enter (Shift+Enter = новая строка)
  function handleKeyDown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  function formatTime(iso) {
    return new Date(iso).toLocaleTimeString('ru', { hour: '2-digit', minute: '2-digit' })
  }

  function formatDate(iso) {
    return new Date(iso).toLocaleDateString('ru', { day: 'numeric', month: 'long' })
  }

  // группируем сообщения по дате
  function groupByDate(msgs) {
    const groups = []
    let lastDate = null
    for (const msg of msgs) {
      const d = new Date(msg.created_at).toDateString()
      if (d !== lastDate) {
        groups.push({ type: 'date', label: formatDate(msg.created_at) })
        lastDate = d
      }
      groups.push({ type: 'msg', msg })
    }
    return groups
  }

  return (
    <div className="chat-page">
      <ToastContainer toasts={toasts} />

      <nav className="navbar chat-navbar">
        <button onClick={() => navigate('/messages')} className="back-btn">← Назад</button>
        {partner && (
          <div className="chat-partner">
            <div className="chat-partner-avatar">
              {partner.avatar
                ? <img src={`http://localhost:8080${partner.avatar}`} alt={partner.name} />
                : <span>{partner.name?.[0]?.toUpperCase()}</span>
              }
            </div>
            <span className="chat-partner-name">{partner.name}</span>
          </div>
        )}
        <div /> {/* spacer */}
      </nav>

      <div className="chat-messages">
        {loading && <p className="feed-status">Загружаем историю...</p>}

        {!loading && messages.length === 0 && (
          <p className="feed-status">Начни переписку!</p>
        )}

        {groupByDate(messages).map((item, i) =>
          item.type === 'date'
            ? <div key={`date-${i}`} className="chat-date-divider">{item.label}</div>
            : (
              <div
                key={item.msg.id}
                className={`chat-bubble-wrap ${item.msg.sender_id === myId ? 'mine' : 'theirs'}`}
              >
                <div className="chat-bubble">
                  <p className="chat-bubble-text">{item.msg.content}</p>
                  <span className="chat-bubble-time">{formatTime(item.msg.created_at)}</span>
                </div>
              </div>
            )
        )}
        <div ref={bottomRef} />
      </div>

      <form onSubmit={handleSend} className="chat-input-form">
        <textarea
          ref={inputRef}
          className="chat-input"
          placeholder="Напиши сообщение..."
          value={text}
          onChange={e => setText(e.target.value)}
          onKeyDown={handleKeyDown}
          rows={1}
          maxLength={5000}
        />
        <button
          type="submit"
          className="chat-send-btn"
          disabled={sending || !text.trim()}
        >
          ↑
        </button>
      </form>
    </div>
  )
}
