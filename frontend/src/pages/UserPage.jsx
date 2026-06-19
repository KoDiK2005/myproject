import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { getUser } from '../api/users'
import { getPostsByUser } from '../api/posts'
import {
  getFriendStatus,
  sendFriendRequest,
  acceptFriendRequest,
  rejectFriendRequest,
  removeFriend,
} from '../api/friends'
import { getCurrentUserId } from '../api/auth'
import PostCard from '../components/PostCard'
import { useToast, ToastContainer } from '../components/Toast'

export default function UserPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { toasts, showToast } = useToast()

  const currentUserId = getCurrentUserId()
  const isOwnProfile = currentUserId === parseInt(id)

  const [user, setUser] = useState(null)
  const [posts, setPosts] = useState([])
  const [friendStatus, setFriendStatus] = useState(null) // FriendshipStatus
  const [loading, setLoading] = useState(true)
  const [actionLoading, setActionLoading] = useState(false)

  useEffect(() => {
    if (isOwnProfile) {
      navigate('/profile', { replace: true })
      return
    }
    loadData()
  }, [id])

  async function loadData() {
    setLoading(true)
    try {
      const [userData, userPosts] = await Promise.all([
        getUser(id),
        getPostsByUser(id),
        // статус дружбы только если залогинен
        currentUserId
          ? getFriendStatus(id).then(s => setFriendStatus(s)).catch(() => {})
          : Promise.resolve(),
      ])
      setUser(userData)
      setPosts(userPosts ?? [])
    } catch (err) {
      showToast(err.message)
    } finally {
      setLoading(false)
    }
  }

  // После действия — перезагружаем статус
  async function refreshStatus() {
    try {
      const s = await getFriendStatus(id)
      setFriendStatus(s)
    } catch {}
  }

  async function handleSendRequest() {
    setActionLoading(true)
    try {
      await sendFriendRequest(id)
      await refreshStatus()
      showToast('Заявка отправлена!', 'success')
    } catch (err) {
      showToast(err.message)
    } finally {
      setActionLoading(false)
    }
  }

  async function handleAccept() {
    setActionLoading(true)
    try {
      await acceptFriendRequest(id)
      await refreshStatus()
      showToast('Теперь вы друзья!', 'success')
    } catch (err) {
      showToast(err.message)
    } finally {
      setActionLoading(false)
    }
  }

  async function handleReject() {
    setActionLoading(true)
    try {
      await rejectFriendRequest(id)
      await refreshStatus()
      showToast('Заявка отклонена', 'info')
    } catch (err) {
      showToast(err.message)
    } finally {
      setActionLoading(false)
    }
  }

  async function handleRemove() {
    if (!confirm('Удалить из друзей?')) return
    setActionLoading(true)
    try {
      await removeFriend(id)
      await refreshStatus()
      showToast('Удалён из друзей', 'info')
    } catch (err) {
      showToast(err.message)
    } finally {
      setActionLoading(false)
    }
  }

  // кнопка дружбы в зависимости от статуса
  function FriendButton() {
    if (!currentUserId) return null // гость — кнопки нет
    if (!friendStatus) return null

    const { status } = friendStatus
    const disabled = actionLoading

    if (status === 'none') return (
      <button className="btn-primary" onClick={handleSendRequest} disabled={disabled}>
        + Добавить в друзья
      </button>
    )
    if (status === 'pending_sent') return (
      <button className="btn-secondary" disabled>Заявка отправлена</button>
    )
    if (status === 'pending_received') return (
      <div className="friend-action-group">
        <button className="btn-primary" onClick={handleAccept} disabled={disabled}>Принять заявку</button>
        <button className="btn-secondary" onClick={handleReject} disabled={disabled}>Отклонить</button>
      </div>
    )
    if (status === 'accepted') return (
      <button className="btn-danger" onClick={handleRemove} disabled={disabled}>
        Удалить из друзей
      </button>
    )
    return null
  }

  if (loading) return <p className="feed-status">Загружаем профиль...</p>
  if (!user) return <p className="feed-status">Пользователь не найден</p>

  return (
    <div className="profile-page">
      <ToastContainer toasts={toasts} />

      <nav className="navbar">
        <button onClick={() => navigate(-1)} className="back-btn">← Назад</button>
        <Link to="/" className="navbar-brand">MyProject</Link>
      </nav>

      <div className="profile-card">
        <div className="profile-avatar">
          {user.avatar_url
            ? <img src={`http://localhost:8080${user.avatar_url}`} alt={user.name} className="avatar-img" />
            : <div className="avatar-placeholder-lg">{user.name?.[0]?.toUpperCase()}</div>
          }
        </div>
        <div className="profile-info">
          <h1>{user.name}</h1>
          <p className="profile-email">{user.email}</p>
          {friendStatus?.status === 'accepted' && (
            <span className="friend-badge">👥 Друг</span>
          )}
        </div>
        <div className="profile-actions">
          <FriendButton />
        </div>
      </div>

      <div className="profile-posts">
        <h2>Посты</h2>
        {posts.length === 0
          ? <p className="feed-status">Постов нет</p>
          : posts.map(post => <PostCard key={post.id} post={post} onDelete={() => {}} />)
        }
      </div>
    </div>
  )
}
