import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  getFriends,
  getIncomingRequests,
  getOutgoingRequests,
  acceptFriendRequest,
  rejectFriendRequest,
  removeFriend,
} from '../api/friends'
import { useToast, ToastContainer } from '../components/Toast'

// вкладки: friends | incoming | outgoing
export default function FriendsPage() {
  const navigate = useNavigate()
  const { toasts, showToast } = useToast()
  const [tab, setTab] = useState('friends')

  const [friends, setFriends] = useState([])
  const [incoming, setIncoming] = useState([])
  const [outgoing, setOutgoing] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadAll()
  }, [])

  async function loadAll() {
    setLoading(true)
    try {
      const [f, inc, out] = await Promise.all([
        getFriends(),
        getIncomingRequests(),
        getOutgoingRequests(),
      ])
      setFriends(f ?? [])
      setIncoming(inc ?? [])
      setOutgoing(out ?? [])
    } catch (err) {
      showToast(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function handleAccept(userId) {
    try {
      await acceptFriendRequest(userId)
      // перемещаем из incoming во friends (упрощённо — перезагружаем)
      await loadAll()
      showToast('Теперь вы друзья!', 'success')
    } catch (err) {
      showToast(err.message)
    }
  }

  async function handleReject(userId) {
    try {
      await rejectFriendRequest(userId)
      setIncoming(prev => prev.filter(u => u.id !== userId))
      showToast('Заявка отклонена', 'info')
    } catch (err) {
      showToast(err.message)
    }
  }

  async function handleRemove(userId) {
    if (!confirm('Удалить из друзей?')) return
    try {
      await removeFriend(userId)
      setFriends(prev => prev.filter(u => u.id !== userId))
      showToast('Удалён из друзей', 'info')
    } catch (err) {
      showToast(err.message)
    }
  }

  function UserCard({ user, actions }) {
    return (
      <div className="friend-card">
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
        <div className="friend-actions">{actions}</div>
      </div>
    )
  }

  return (
    <div className="friends-page">
      <ToastContainer toasts={toasts} />

      <nav className="navbar">
        <button onClick={() => navigate('/')} className="back-btn">← Назад</button>
        <span className="navbar-brand">Друзья</span>
      </nav>

      <div className="friends-tabs">
        <button
          className={`tab-btn ${tab === 'friends' ? 'active' : ''}`}
          onClick={() => setTab('friends')}
        >
          Друзья {friends.length > 0 && <span className="badge">{friends.length}</span>}
        </button>
        <button
          className={`tab-btn ${tab === 'incoming' ? 'active' : ''}`}
          onClick={() => setTab('incoming')}
        >
          Входящие {incoming.length > 0 && <span className="badge badge-red">{incoming.length}</span>}
        </button>
        <button
          className={`tab-btn ${tab === 'outgoing' ? 'active' : ''}`}
          onClick={() => setTab('outgoing')}
        >
          Исходящие {outgoing.length > 0 && <span className="badge">{outgoing.length}</span>}
        </button>
      </div>

      <div className="friends-content">
        {loading ? (
          <p className="feed-status">Загружаем...</p>
        ) : tab === 'friends' ? (
          friends.length === 0
            ? <p className="feed-status">Друзей пока нет. Найди людей!</p>
            : friends.map(u => (
              <UserCard key={u.id} user={u} actions={
                <button className="btn-danger" onClick={() => handleRemove(u.id)}>
                  Удалить
                </button>
              } />
            ))
        ) : tab === 'incoming' ? (
          incoming.length === 0
            ? <p className="feed-status">Нет входящих заявок</p>
            : incoming.map(u => (
              <UserCard key={u.id} user={u} actions={
                <>
                  <button className="btn-primary" onClick={() => handleAccept(u.id)}>Принять</button>
                  <button className="btn-secondary" onClick={() => handleReject(u.id)}>Отклонить</button>
                </>
              } />
            ))
        ) : (
          outgoing.length === 0
            ? <p className="feed-status">Нет исходящих заявок</p>
            : outgoing.map(u => (
              <UserCard key={u.id} user={u} actions={
                <span className="status-label">Ожидает ответа...</span>
              } />
            ))
        )}
      </div>
    </div>
  )
}
