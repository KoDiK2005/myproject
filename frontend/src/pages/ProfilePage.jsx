import { useState, useEffect, useRef } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { getCurrentUserId, logout } from '../api/auth'
import { getUser, updateUser, uploadAvatar } from '../api/users'

const BASE_URL = 'http://localhost:8080'

export default function ProfilePage() {
  const navigate = useNavigate()
  const fileInputRef = useRef(null)

  const [user, setUser] = useState(null)
  const [form, setForm] = useState({ name: '', email: '' })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [uploadingAvatar, setUploadingAvatar] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const userId = getCurrentUserId()

  useEffect(() => {
    if (!userId) { navigate('/login'); return }

    getUser(userId)
      .then(data => {
        setUser(data)
        setForm({ name: data.name, email: data.email })
      })
      .catch(err => setError(err.message))
      .finally(() => setLoading(false))
  }, [userId])

  function handleChange(e) {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  async function handleSave(e) {
    e.preventDefault()
    setError('')
    setSuccess('')
    setSaving(true)
    try {
      const updated = await updateUser(userId, form.name, form.email)
      setUser(updated)
      setSuccess('Профиль обновлён')
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  async function handleAvatarChange(e) {
    const file = e.target.files[0]
    if (!file) return

    setError('')
    setUploadingAvatar(true)
    try {
      const data = await uploadAvatar(userId, file)
      // обновляем аватар без перезагрузки всего профиля
      setUser(prev => ({ ...prev, avatar: data.avatar }))
      setSuccess('Аватар обновлён')
    } catch (err) {
      setError(err.message)
    } finally {
      setUploadingAvatar(false)
    }
  }

  async function handleLogout() {
    await logout()
    navigate('/login')
  }

  if (loading) return <div className="auth-container"><p>Загружаем профиль...</p></div>

  return (
    <div className="profile-page">
      <nav className="profile-nav">
        <Link to="/">← Лента</Link>
        <button onClick={handleLogout} className="logout-btn">Выйти</button>
      </nav>

      <div className="profile-container">
        <h1>Профиль</h1>

        {/* Аватар */}
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

        {/* Форма */}
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

          {error && <p className="error">{error}</p>}
          {success && <p className="success">{success}</p>}

          <button type="submit" disabled={saving}>
            {saving ? 'Сохраняем...' : 'Сохранить'}
          </button>
        </form>
      </div>
    </div>
  )
}
