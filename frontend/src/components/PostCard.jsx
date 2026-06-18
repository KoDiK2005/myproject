import { getCurrentUserId } from '../api/auth'

export default function PostCard({ post, onDelete }) {
  const currentUserId = getCurrentUserId()
  const isOwner = currentUserId === post.user_id

  const date = new Date(post.created_at).toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'long', year: 'numeric',
  })

  async function handleDelete() {
    if (!confirm('Удалить пост?')) return
    onDelete(post.id)
  }

  return (
    <div className="post-card">
      <div className="post-header">
        <h2 className="post-title">{post.title}</h2>
        {isOwner && (
          <button className="post-delete-btn" onClick={handleDelete} title="Удалить">✕</button>
        )}
      </div>
      <p className="post-body">{post.body}</p>
      <div className="post-meta">
        <span>Пользователь #{post.user_id}</span>
        <span>{date}</span>
      </div>
    </div>
  )
}
