import { useState, useEffect } from 'react'
import { getIncomingRequests } from '../api/friends'
import { getAccessToken } from '../api/auth'

// каждые 30 сек проверяем входящие заявки
const POLL_INTERVAL = 30_000

export function useFriendBadge() {
  const [count, setCount] = useState(0)

  useEffect(() => {
    if (!getAccessToken()) return

    async function check() {
      try {
        const incoming = await getIncomingRequests()
        setCount(incoming?.length ?? 0)
      } catch {
        // тихо — не мешаем работе приложения
      }
    }

    check()
    const id = setInterval(check, POLL_INTERVAL)
    return () => clearInterval(id)
  }, [])

  return count
}
