import { Link } from 'react-router-dom'
import { getCurrentUserId } from '../api/auth'

export default function PostCard({ post, onDelete }) {
  const currentUserId = getCurrentUserId()
  const isOwner = currentUserId === post.user_id

  async function handleDelete(e) {
    e.preventDefault()
    if (!confirm('Удалить пост?')) return
    onDelete(post.id)
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
        </div>
      </div>
    </Link>
  )
}
