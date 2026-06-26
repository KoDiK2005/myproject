import { useState, useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { apiFetch } from '../api/auth'
import { useToast, ToastContainer } from '../components/Toast'
import { sendFriendRequest, getFriendStatus } from '../api/friends'
import Navbar from '../components/Navbar'

const BASE = 'http://localhost:8080/api/v1'

async function searchUsers(query, page = 1, limit = 20) {
  const params = new URLSearchParams({ search: query, page, limit })
  const res = await apiFetch(`${BASE}/users?${params}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка')
  return data
}

export default function PeoplePage() {
  const navigate = useNavigate()
  const { toasts, showToast } = useToast()

  const [query, setQuery] = useState('')
  const [results, setResults] = useState([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [searched, setSearched] = useState(false)

  // статусы дружбы: { userId: 'none' | 'pending_sent' | ... }
  const [statuses, setStatuses] = useState({})
  const [sending, setSending] = useState({})

  const handleSearch = useCallback(async (e) => {
    e?.preventDefault()
    if (!query.trim()) return
    setLoading(true)
    setSearched(true)
    try {
      const data = await searchUsers(query.trim())
      setResults(data.data ?? [])
      setTotal(data.total ?? 0)

      // загружаем статусы дружбы для всех найденных
      const statusMap = {}
      await Promise.all(
        (data.data ?? []).map(async user => {
          try {
            const s = await getFriendStatus(user.id)
            statusMap[user.id] = s.status
          } catch {
            statusMap[user.id] = 'none'
          }
        })
      )
      setStatuses(statusMap)
    } catch (err) {
      showToast(err.message)
    } finally {
      setLoading(false)
    }
  }, [query])

  async function handleAdd(userId) {
    setSending(prev => ({ ...prev, [userId]: true }))
    try {
      await sendFriendRequest(userId)
      setStatuses(prev => ({ ...prev, [userId]: 'pending_sent' }))
      showToast('Заявка отправлена!', 'success')
    } catch (err) {
      showToast(err.message)
    } finally {
      setSending(prev => ({ ...prev, [userId]: false }))
    }
  }

  function FriendBtn({ user }) {
    const status = statuses[user.id]
    const busy = sending[user.id]

    if (status === 'accepted') return <span className="friend-badge">👥 Друг</span>
    if (status === 'pending_sent') return <button className="btn-secondary" disabled>Заявка отправлена</button>
    if (status === 'pending_received') return <span className="status-label">Хочет дружить с тобой</span>
    return (
      <button
        className="btn-primary"
        onClick={() => handleAdd(user.id)}
        disabled={busy}
      >
        + Добавить
      </button>
    )
  }

  return (
    <div className="people-page">
      <ToastContainer toasts={toasts} />

      <Navbar
        onBack={() => navigate('/')}
        title="Поиск людей"
        links={[{ to: '/friends', label: 'Друзья' }]}
      />

      <div className="people-content">
        <form onSubmit={handleSearch} className="people-search-form">
          <input
            type="text"
            placeholder="Имя пользователя..."
            value={query}
            onChange={e => setQuery(e.target.value)}
            className="search-input"
            autoFocus
          />
          <button type="submit" className="btn-primary" disabled={loading || !query.trim()}>
            Найти
          </button>
        </form>

        {loading && <p className="feed-status">Ищем...</p>}

        {!loading && searched && results.length === 0 && (
          <p className="feed-status">Никого не нашли. Попробуй другое имя.</p>
        )}

        {!loading && results.length > 0 && (
          <>
            <p className="people-count">Найдено: {total}</p>
            <div className="people-list">
              {results.map(user => (
                <div key={user.id} className="friend-card">
                  <Link to={`/users/${user.id}`} className="friend-info">
                    <div className="friend-avatar">
                      {user.avatar
                        ? <img src={`http://localhost:8080${user.avatar}`} alt={user.name} />
                        : <span className="avatar-placeholder">{user.name?.[0]?.toUpperCase()}</span>
                      }
                    </div>
                    <div>
                      <strong>{user.name}</strong>
                      <p className="friend-email">{user.email}</p>
                    </div>
                  </Link>
                  <div className="friend-actions">
                    <FriendBtn user={user} />
                  </div>
                </div>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  )
}
