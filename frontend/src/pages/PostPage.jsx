import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { getCurrentUserId } from '../api/auth'
import { getComments, createComment, deleteComment } from '../api/comments'
import { getLikeCount, likePost, unlikePost } from '../api/likes'
import { getPost, updatePost } from '../api/posts'
import { useToast, ToastContainer } from '../components/Toast'

export default function PostPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { toasts, showToast } = useToast()
  const currentUserId = getCurrentUserId()
  const postId = Number(id)

  const [post, setPost] = useState(null)
  const [comments, setComments] = useState([])
  const [likeCount, setLikeCount] = useState(0)
  const [liked, setLiked] = useState(false)
  const [commentBody, setCommentBody] = useState('')
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  const [editing, setEditing] = useState(false)
  const [editForm, setEditForm] = useState({ title: '', body: '', visibility: 'public' })
  const [savingEdit, setSavingEdit] = useState(false)

  useEffect(() => {
    async function load() {
      try {
        const [postRes, commentsData, likesData] = await Promise.all([
          getPost(postId),
          getComments(postId),
          getLikeCount(postId),
        ])
        setPost(postRes)
        setComments(commentsData ?? [])
        setLikeCount(likesData.count ?? 0)
        setLiked(likesData.liked ?? false)
      } catch (err) {
        showToast(err.message)
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
      showToast(err.message)
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
      showToast(err.message)
    } finally {
      setSubmitting(false)
    }
  }

  function startEditing() {
    setEditForm({ title: post.title, body: post.body, visibility: post.visibility })
    setEditing(true)
  }

  async function handleSaveEdit(e) {
    e.preventDefault()
    if (!editForm.title.trim() || !editForm.body.trim()) return
    setSavingEdit(true)
    try {
      const updated = await updatePost(postId, editForm.title, editForm.body, editForm.visibility)
      setPost(updated)
      setEditing(false)
      showToast('Пост обновлён', 'success')
    } catch (err) {
      showToast(err.message)
    } finally {
      setSavingEdit(false)
    }
  }

  async function handleDeleteComment(commentId) {
    try {
      await deleteComment(commentId)
      setComments(prev => prev.filter(c => c.id !== commentId))
      showToast('Комментарий удалён', 'info')
    } catch (err) {
      showToast(err.message)
    }
  }

  if (loading) return <div className="auth-container"><p>Загрузка...</p></div>
  if (!post || post.error) return <div className="auth-container"><p>Пост не найден</p></div>

  return (
    <div className="post-page">
      <ToastContainer toasts={toasts} />

      <nav className="navbar">
        <button className="back-btn" onClick={() => navigate(-1)}>← Назад</button>
        <Link to="/profile" style={{ color: '#aaa', textDecoration: 'none', fontSize: '0.9rem' }}>Профиль</Link>
      </nav>

      <div className="post-full">
        {editing ? (
          <form onSubmit={handleSaveEdit} className="create-form">
            <input
              type="text"
              value={editForm.title}
              onChange={e => setEditForm({ ...editForm, title: e.target.value })}
              required
              maxLength={200}
            />
            <textarea
              value={editForm.body}
              onChange={e => setEditForm({ ...editForm, body: e.target.value })}
              required
              rows={6}
            />
            <div className="visibility-toggle">
              <button
                type="button"
                className={`vis-btn ${editForm.visibility === 'public' ? 'active' : ''}`}
                onClick={() => setEditForm({ ...editForm, visibility: 'public' })}
              >
                🌍 Публично
              </button>
              <button
                type="button"
                className={`vis-btn ${editForm.visibility === 'friends' ? 'active' : ''}`}
                onClick={() => setEditForm({ ...editForm, visibility: 'friends' })}
              >
                🔒 Только друзья
              </button>
            </div>
            <div className="friend-action-group">
              <button type="submit" disabled={savingEdit}>
                {savingEdit ? 'Сохраняем...' : 'Сохранить'}
              </button>
              <button type="button" className="btn-secondary" onClick={() => setEditing(false)}>
                Отмена
              </button>
            </div>
          </form>
        ) : (
          <>
            <div className="post-header">
              <h1 className="post-full-title">{post.title}</h1>
              {post.user_id === currentUserId && (
                <button className="post-delete-btn" onClick={startEditing} title="Редактировать">✎</button>
              )}
            </div>
            <p className="post-full-meta">{post.author_name || `Автор #${post.user_id}`}</p>
            <p className="post-full-body">{post.body}</p>
          </>
        )}

        <div className="like-section">
          <button className={`like-btn ${liked ? 'liked' : ''}`} onClick={handleLike}>
            {liked ? '❤️' : '🤍'} {likeCount}
          </button>
        </div>

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
                    <span className="comment-author">{c.author_name || `Пользователь #${c.user_id}`}</span>
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
