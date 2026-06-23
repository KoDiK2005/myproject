import { useState, useEffect, useRef } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { getCurrentUserId, logout, logoutAll } from '../api/auth'
import { getUser, updateUser, uploadAvatar } from '../api/users'
import { useToast, ToastContainer } from '../components/Toast'

const BASE_URL = 'http://localhost:8080'

export default function ProfilePage() {
  const navigate = useNavigate()
  const fileInputRef = useRef(null)
  const { toasts, showToast } = useToast()

  const [user, setUser] = useState(null)
  const [form, setForm] = useState({ name: '', email: '' })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [uploadingAvatar, setUploadingAvatar] = useState(false)

  const userId = getCurrentUserId()

  useEffect(() => {
    if (!userId) { navigate('/login'); return }
    getUser(userId)
      .then(data => {
        setUser(data)
        setForm({ name: data.name, email: data.email })
      })
      .catch(err => showToast(err.message))
      .finally(() => setLoading(false))
  }, [userId])

  function handleChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function handleSave(e) {
    e.preventDefault()
    setSaving(true)
    try {
      const updated = await updateUser(userId, form.name, form.email)
      setUser(updated)
      showToast('Профиль обновлён', 'success')
    } catch (err) {
      showToast(err.message)
    } finally {
      setSaving(false)
    }
  }

  async function handleAvatarChange(e) {
    const file = e.target.files[0]
    if (!file) return
    setUploadingAvatar(true)
    try {
      const data = await uploadAvatar(userId, file)
      setUser(prev => ({ ...prev, avatar: data.avatar }))
      showToast('Аватар обновлён', 'success')
    } catch (err) {
      showToast(err.message)
    } finally {
      setUploadingAvatar(false)
    }
  }

  async function handleLogout() {
    await logout()
    navigate('/login')
  }

  async function handleLogoutAll() {
    if (!confirm('Завершить все сессии на всех устройствах?')) return
    await logoutAll()
    navigate('/login')
  }

  if (loading) return <div className="auth-container"><p>Загружаем профиль...</p></div>

  return (
    <div className="profile-page">
      <ToastContainer toasts={toasts} />

      <nav className="profile-nav">
        <Link to="/">← Лента</Link>
        <button onClick={handleLogout} className="logout-btn">Выйти</button>
      </nav>

      <div className="profile-container">
        <h1>Профиль</h1>

        <div className="avatar-section">
          <div className="avatar-wrapper" onClick={() => fileInputRef.current.click()}>
            {user?.avatar
              ? <img src={`${BASE_URL}${user.avatar}`} alt="avatar" className="avatar-img" />
              : <div className="avatar-placeholder">{user?.name?.[0]?.toUpperCase()}</div>
            }
            <div className="avatar-overlay">
              {uploadingAvatar ? '...' : '📷'}
            </div>
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept=".jpg,.jpeg,.png"
            style={{ display: 'none' }}
            onChange={handleAvatarChange}
          />
          <p className="avatar-hint">Нажми чтобы сменить аватар (jpg/png, макс 2MB)</p>
        </div>

        <form onSubmit={handleSave} className="profile-form">
          <div className="field">
            <label>Имя</label>
            <input
              type="text"
              name="name"
              value={form.name}
              onChange={handleChange}
              minLength={2}
              maxLength={50}
              required
            />
          </div>
          <div className="field">
            <label>Email</label>
            <input
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
              required
            />
          </div>

          <button type="submit" disabled={saving}>
            {saving ? 'Сохраняем...' : 'Сохранить'}
          </button>
        </form>

        <div className="profile-security">
          <h2>Безопасность</h2>
          <p className="feed-status">Если подозреваешь, что аккаунт мог быть скомпрометирован — завершит все активные сессии на всех устройствах.</p>
          <button type="button" className="btn-danger" onClick={handleLogoutAll}>
            Выйти со всех устройств
          </button>
        </div>
      </div>
    </div>
  )
}
