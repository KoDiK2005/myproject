import { useNavigate, Link } from 'react-router-dom'
import { logout } from '../api/auth'

export default function HomePage() {
  const navigate = useNavigate()

  async function handleLogout() {
    await logout()
    navigate('/login')
  }

  return (
    <div style={{ padding: '2rem', maxWidth: '600px', margin: '0 auto' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1>Главная</h1>
        <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
          <Link to="/profile" style={{ color: '#a78bfa', textDecoration: 'none' }}>Профиль</Link>
          <button onClick={handleLogout} className="logout-btn">Выйти</button>
        </div>
      </div>
      <p style={{ color: '#aaa' }}>Скоро тут будет лента постов.</p>
    </div>
  )
}
