import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { getCurrentUserId } from '../api/auth'
import { getLikeCount, likePost, unlikePost } from '../api/likes'

export default function PostCard({ post, onDelete }) {
  const currentUserId = getCurrentUserId()
  const isOwner = currentUserId === post.user_id

  const [likeCount, setLikeCount] = useState(null)
  const [liked, setLiked] = useState(false)

  useEffect(() => {
    let cancelled = false
    getLikeCount(post.id).then(data => {
      if (cancelled) return
      setLikeCount(data.count ?? 0)
      setLiked(data.liked ?? false)
    }).catch(() => {})
    return () => { cancelled = true }
  }, [post.id])

  async function handleDelete(e) {
    e.preventDefault()
    if (!confirm('Удалить пост?')) return
    onDelete(post.id)
  }

  async function handleLike(e) {
    e.preventDefault()
    e.stopPropagation()
    if (!currentUserId) return
    try {
      if (liked) {
        await unlikePost(post.id)
        setLikeCount(c => c - 1)
      } else {
        await likePost(post.id)
        setLikeCount(c => c + 1)
      }
      setLiked(v => !v)
    } catch {}
  }

  return (
    <Link to={`/posts/${post.id}`} className="post-card-link">
      <div className="post-card">
        <div className="post-header">
          <h2 className="post-title">{post.title}</h2>
          <div className="post-header-right">
            {post.visibility === 'friends' && (
              <span className="visibility-badge" title="Только для друзей">🔒</span>
            )}
            {isOwner && (
              <button className="post-delete-btn" onClick={handleDelete} title="Удалить">✕</button>
            )}
          </div>
        </div>
        <p className="post-body post-body-clamp">{post.body}</p>
        <div className="post-meta">
          <Link
            to={`/users/${post.user_id}`}
            className="post-author-link"
            onClick={e => e.stopPropagation()}
          >
            Автор #{post.user_id}
          </Link>
          {currentUserId && (
            <button className={`like-btn-card ${liked ? 'liked' : ''}`} onClick={handleLike}>
              {liked ? '❤️' : '🤍'} {likeCount ?? ''}
            </button>
          )}
        </div>
      </div>
    </Link>
  )
}
