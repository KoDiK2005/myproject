import { useState, useEffect, useCallback } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { logout } from '../api/auth'
import { getPosts, createPost, deletePost } from '../api/posts'
import PostCard from '../components/PostCard'

export default function HomePage() {
  const navigate = useNavigate()

  const [posts, setPosts] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [searchInput, setSearchInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // форма создания поста
  const [showForm, setShowForm] = useState(false)
  const [newPost, setNewPost] = useState({ title: '', body: '' })
  const [creating, setCreating] = useState(false)

  const limit = 10
  const totalPages = Math.ceil(total / limit)

  const fetchPosts = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const res = await getPosts({ page, limit, search })
      setPosts(res.data ?? [])
      setTotal(res.total)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }, [page, search])

  useEffect(() => {
    fetchPosts()
  }, [fetchPosts])

  function handleSearch(e) {
    e.preventDefault()
    setPage(1)
    setSearch(searchInput)
  }

  async function handleCreate(e) {
    e.preventDefault()
    if (!newPost.title.trim() || !newPost.body.trim()) return
    setCreating(true)
    try {
      const post = await createPost(newPost.title, newPost.body)
      setPosts(prev => [post, ...prev])
      setTotal(prev => prev + 1)
      setNewPost({ title: '', body: '' })
      setShowForm(false)
    } catch (err) {
      setError(err.message)
    } finally {
      setCreating(false)
    }
  }

  async function handleDelete(id) {
    try {
      await deletePost(id)
      setPosts(prev => prev.filter(p => p.id !== id))
      setTotal(prev => prev - 1)
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleLogout() {
    await logout()
    navigate('/login')
  }

  return (
    <div className="home-page">
      {/* Navbar */}
      <nav className="navbar">
        <span className="navbar-brand">MyProject</span>
        <div className="navbar-links">
          <Link to="/profile">Профиль</Link>
          <button onClick={handleLogout} className="logout-btn">Выйти</button>
        </div>
      </nav>

      <div className="feed">
        {/* Поиск */}
        <form onSubmit={handleSearch} className="search-form">
          <input
            type="text"
            placeholder="Поиск постов..."
            value={searchInput}
            onChange={e => setSearchInput(e.target.value)}
            className="search-input"
          />
          <button type="submit" className="search-btn">Найти</button>
          {search && (
            <button type="button" className="search-btn" onClick={() => { setSearch(''); setSearchInput(''); setPage(1) }}>
              Сбросить
            </button>
          )}
        </form>

        {/* Кнопка создать пост */}
        <button className="create-btn" onClick={() => setShowForm(v => !v)}>
          {showForm ? 'Отмена' : '+ Новый пост'}
        </button>

        {/* Форма создания */}
        {showForm && (
          <form onSubmit={handleCreate} className="create-form">
            <input
              type="text"
              placeholder="Заголовок"
              value={newPost.title}
              onChange={e => setNewPost({ ...newPost, title: e.target.value })}
              required
              maxLength={200}
            />
            <textarea
              placeholder="Текст поста..."
              value={newPost.body}
              onChange={e => setNewPost({ ...newPost, body: e.target.value })}
              required
              rows={4}
            />
            <button type="submit" disabled={creating}>
              {creating ? 'Публикуем...' : 'Опубликовать'}
            </button>
          </form>
        )}

        {error && <p className="error">{error}</p>}

        {/* Лента */}
        {loading ? (
          <p className="feed-status">Загружаем посты...</p>
        ) : posts.length === 0 ? (
          <p className="feed-status">Постов пока нет. Будь первым!</p>
        ) : (
          posts.map(post => (
            <PostCard key={post.id} post={post} onDelete={handleDelete} />
          ))
        )}

        {/* Пагинация */}
        {totalPages > 1 && (
          <div className="pagination">
            <button onClick={() => setPage(p => p - 1)} disabled={page === 1}>← Назад</button>
            <span>{page} / {totalPages}</span>
            <button onClick={() => setPage(p => p + 1)} disabled={page === totalPages}>Вперёд →</button>
          </div>
        )}
      </div>
    </div>
  )
}
