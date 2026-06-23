import { useState, useEffect } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import { verifyEmail } from '../api/auth'

export default function VerifyEmailPage() {
  const [params] = useSearchParams()
  const token = params.get('token')
  const [status, setStatus] = useState('loading') // loading | success | error
  const [error, setError] = useState('')

  useEffect(() => {
    if (!token) {
      setStatus('error')
      setError('Ссылка неполная — отсутствует токен')
      return
    }
    verifyEmail(token)
      .then(() => setStatus('success'))
      .catch(err => {
        setStatus('error')
        setError(err.message)
      })
  }, [token])

  return (
    <div className="auth-container">
      <div className="auth-box">
        <h1>Подтверждение email</h1>
        {status === 'loading' && <p>Проверяем ссылку...</p>}
        {status === 'success' && <p>Email подтверждён! Можешь пользоваться всеми функциями.</p>}
        {status === 'error' && <p className="error">{error}</p>}
        <p className="auth-link"><Link to="/">← На главную</Link></p>
      </div>
    </div>
  )
}
