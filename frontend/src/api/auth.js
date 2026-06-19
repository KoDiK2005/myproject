const BASE_URL = 'http://localhost:8080'

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

// Декодируем JWT без библиотек — нам нужен только user_id
export function getCurrentUserId() {
  const token = getAccessToken()
  if (!token) return null
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.user_id
  } catch {
    return null
  }
}

// Тихо обновляет access token через refresh token.
// Возвращает новый access token или null если refresh истёк.
async function silentRefresh() {
  const refresh = localStorage.getItem('refresh_token')
  if (!refresh) return null
  try {
    const res = await fetch(`${BASE_URL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refresh }),
    })
    if (!res.ok) {
      clearTokens() // refresh протух — выкидываем юзера
      return null
    }
    const data = await res.json()
    saveTokens(data.access_token, data.refresh_token)
    return data.access_token
  } catch {
    return null
  }
}

// Универсальная fetch-обёртка с авто-рефрешем.
// Если ответ 401 — пробуем обновить токен и повторяем запрос один раз.
export async function apiFetch(url, options = {}) {
  const makeRequest = (token) =>
    fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
    })

  let res = await makeRequest(getAccessToken())

  if (res.status === 401) {
    const newToken = await silentRefresh()
    if (!newToken) {
      // refresh тоже сдох — отправляем на логин
      window.location.href = '/login'
      throw new Error('Сессия истекла')
    }
    res = await makeRequest(newToken) // повторяем с новым токеном
  }

  return res
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
    }).catch(() => {})
  }
  clearTokens()
}
