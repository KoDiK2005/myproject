const BASE_URL = 'http://localhost:8080'

// refresh_token больше не трогаем из JS — он лежит в httpOnly cookie,
// браузер сам шлёт её на /auth/* при credentials: 'include'.
export function saveAccessToken(access) {
  localStorage.setItem('access_token', access)
}

export function getAccessToken() {
  return localStorage.getItem('access_token')
}

export function clearTokens() {
  localStorage.removeItem('access_token')
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

// Тихо обновляет access token через refresh-cookie.
// Возвращает новый access token или null если refresh истёк/отозван.
async function silentRefresh() {
  try {
    const res = await fetch(`${BASE_URL}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    })
    if (!res.ok) {
      clearTokens() // refresh протух — выкидываем юзера
      return null
    }
    const data = await res.json()
    saveAccessToken(data.access_token)
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
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка входа')
  saveAccessToken(data.access_token)
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
  await fetch(`${BASE_URL}/auth/logout`, {
    method: 'POST',
    credentials: 'include',
  }).catch(() => {})
  clearTokens()
}

// выйти со всех устройств — отзывает все refresh-токены (например при компрометации аккаунта)
export async function logoutAll() {
  await apiFetch(`${BASE_URL}/auth/logout-all`, { method: 'POST', credentials: 'include' }).catch(() => {})
  clearTokens()
}

export async function verifyEmail(token) {
  const res = await fetch(`${BASE_URL}/auth/verify-email?token=${encodeURIComponent(token)}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Не удалось подтвердить email')
  return data
}

export async function resendVerification() {
  const res = await apiFetch(`${BASE_URL}/auth/resend-verification`, { method: 'POST' })
  if (res.status === 204) return
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Не удалось отправить письмо')
}
