import { getAccessToken } from './auth'

const BASE_URL = 'http://localhost:8080'

// Хелпер — добавляет Authorization header к каждому запросу
function authHeaders(extra = {}) {
  return {
    Authorization: `Bearer ${getAccessToken()}`,
    ...extra,
  }
}

export async function getUser(id) {
  const res = await fetch(`${BASE_URL}/api/v1/users/${id}`, {
    headers: authHeaders(),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки профиля')
  return data
}

export async function updateUser(id, name, email) {
  const res = await fetch(`${BASE_URL}/api/v1/users/${id}`, {
    method: 'PUT',
    headers: authHeaders({ 'Content-Type': 'application/json' }),
    body: JSON.stringify({ name, email }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка обновления профиля')
  return data
}

export async function uploadAvatar(id, file) {
  const form = new FormData()
  form.append('avatar', file)

  const res = await fetch(`${BASE_URL}/api/v1/users/${id}/avatar`, {
    method: 'POST',
    headers: authHeaders(), // Content-Type НЕ ставим — браузер сам ставит multipart/form-data с boundary
    body: form,
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || 'Ошибка загрузки аватара')
  return data
}
