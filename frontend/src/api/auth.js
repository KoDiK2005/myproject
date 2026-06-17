const BASE_URL = 'http://localhost:8080'

// Сохраняем токены в localStorage
export function saveTokens(access, refresh) {
  localStorage.setItem('access_token', access)
  localStorage.setItem('refresh_token', refresh)
}

export function getAccessToken() {
  return localStorage.getItem('access_token')
}

export function clearTokens() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
}

export async function login(email, password) {
  const res = await fetch(`${BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка входа')
  saveTokens(data.access_token, data.refresh_token)
  return data
}

export async function register(name, email, password) {
  const res = await fetch(`${BASE_URL}/api/v1/users`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email, password }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка регистрации')
  return data
}

export async function logout() {
  const refresh = localStorage.getItem('refresh_token')
  if (refresh) {
    await fetch(`${BASE_URL}/auth/logout`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refresh }),
    }).catch(() => {}) // бэк упал — не страшно, всё равно чистим
  }
  clearTokens()
}
