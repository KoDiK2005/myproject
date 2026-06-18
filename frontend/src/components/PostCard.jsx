import { Link } from 'react-router-dom'
import { getCurrentUserId } from '../api/auth'

export default function PostCard({ post, onDelete }) {
  const currentUserId = getCurrentUserId()
  const isOwner = currentUserId === post.user_id

  async function handleDelete(e) {
    e.preventDefault() // чтобы не переходить по ссылке при клике на крестик
    if (!confirm('Удалить пост?')) return
    onDelete(post.id)
  }

  return (
    <Link to={`/posts/${post.id}`} className="post-card-link">
      <div className="post-card">
        <div className="post-header">
          <h2 className="post-title">{post.title}</h2>
          {isOwner && (
            <button className="post-delete-btn" onClick={handleDelete} title="Удалить">✕</button>
          )}
        </div>
        <p className="post-body post-body-clamp">{post.body}</p>
        <div className="post-meta">
          <span>Автор #{post.user_id}</span>
        </div>
      </div>
    </Link>
  )
}
