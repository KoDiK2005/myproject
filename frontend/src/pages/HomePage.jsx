import { useNavigate } from 'react-router-dom'
import { logout } from '../api/auth'

export default function HomePage() {
  const navigate = useNavigate()

  async function handleLogout() {
    await logout()
    navigate('/login')
  }

  return (
    <div style={{ padding: '2rem' }}>
      <h1>Главная</h1>
      <p>Ты залогинен 🎉 Скоро тут будет лента постов.</p>
      <button onClick={handleLogout}>Выйти</button>
    </div>
  )
}
