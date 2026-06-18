import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { getCurrentUserId } from '../api/auth'
import { getPosts } from '../api/posts'
import { getComments, createComment, deleteComment } from '../api/comments'
import { getLikeCount, likePost, unlikePost } from '../api/likes'

export default function PostPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const currentUserId = getCurrentUserId()
  const postId = Number(id)

  const [post, setPost] = useState(null)
  const [comments, setComments] = useState([])
  const [likeCount, setLikeCount] = useState(0)
  const [liked, setLiked] = useState(false)
  const [commentBody, setCommentBody] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    async function load() {
      try {
        // getPosts не имеет getById — используем search по id через список
        // но у нас есть эндпоинт /posts/:id — делаем прямой fetch
        const [postRes, commentsData, likesData] = await Promise.all([
          fetch(`http://localhost:8080/api/v1/posts/${postId}`).then(r => r.json()),
          getComments(postId),
          getLikeCount(postId),
        ])
        setPost(postRes)
        setComments(commentsData ?? [])
        setLikeCount(likesData.count ?? 0)
      } catch (err) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [postId])

  async function handleLike() {
    try {
      if (liked) {
        await unlikePost(postId)
        setLikeCount(c => c - 1)
      } else {
        await likePost(postId)
        setLikeCount(c => c + 1)
      }
      setLiked(v => !v)
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleComment(e) {
    e.preventDefault()
    if (!commentBody.trim()) return
    setSubmitting(true)
    try {
      const comment = await createComment(postId, commentBody)
      setComments(prev => [...prev, comment])
      setCommentBody('')
    } catch (err) {
      setError(err.message)
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDeleteComment(commentId) {
    if (!confirm('Удалить комментарий?')) return
    try {
      await deleteComment(commentId)
      setComments(prev => prev.filter(c => c.id !== commentId))
    } catch (err) {
      setError(err.message)
    }
  }

  if (loading) return <div className="auth-container"><p>Загрузка...</p></div>
  if (!post || post.error) return <div className="auth-container"><p>Пост не найден</p></div>

  return (
    <div className="post-page">
      <nav className="navbar">
        <button className="back-btn" onClick={() => navigate(-1)}>← Назад</button>
        <Link to="/profile" className="navbar-links" style={{ color: '#aaa', textDecoration: 'none', fontSize: '0.9rem' }}>Профиль</Link>
      </nav>

      <div className="post-full">
        {/* Пост */}
        <h1 className="post-full-title">{post.title}</h1>
        <p className="post-full-meta">Автор #{post.user_id}</p>
        <p className="post-full-body">{post.body}</p>

        {/* Лайк */}
        <div className="like-section">
          <button
            className={`like-btn ${liked ? 'liked' : ''}`}
            onClick={handleLike}
          >
            {liked ? '❤️' : '🤍'} {likeCount}
          </button>
        </div>

        {error && <p className="error">{error}</p>}

        {/* Комментарии */}
        <div className="comments-section">
          <h2>Комментарии ({comments.length})</h2>

          <form onSubmit={handleComment} className="comment-form">
            <textarea
              placeholder="Написать комментарий..."
              value={commentBody}
              onChange={e => setCommentBody(e.target.value)}
              rows={3}
              required
            />
            <button type="submit" disabled={submitting}>
              {submitting ? 'Отправляем...' : 'Отправить'}
            </button>
          </form>

          <div className="comments-list">
            {comments.length === 0 ? (
              <p className="feed-status">Комментариев пока нет</p>
            ) : (
              comments.map(c => (
                <div key={c.id} className="comment-card">
                  <div className="comment-header">
                    <span className="comment-author">Пользователь #{c.user_id}</span>
                    {c.user_id === currentUserId && (
                      <button className="post-delete-btn" onClick={() => handleDeleteComment(c.id)}>✕</button>
                    )}
                  </div>
                  <p className="comment-body">{c.body}</p>
                  <span className="comment-date">
                    {new Date(c.created_at).toLocaleDateString('ru-RU')}
                  </span>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
