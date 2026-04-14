import { useState, useEffect, FormEvent } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { API_URL } from '../App'

interface Comment {
  id: number
  post_id: number
  user_id: number
  content: string
  created_at: string
  user: { id: number; username: string }
}

interface Post {
  id: number
  title: string
  content: string
  user_id: number
  tags: string[]
  comments: Comment[]
  version: number
  created_at: string
  updated_at: string
  user: { id: number; username: string }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('es-ES', {
    month: 'long', day: 'numeric', year: 'numeric',
  })
}

export function PostPage() {
  const { postID } = useParams<{ postID: string }>()
  const { token, userId } = useAuth()
  const navigate = useNavigate()

  const [post, setPost] = useState<Post | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [editing, setEditing] = useState(false)
  const [editTitle, setEditTitle] = useState('')
  const [editContent, setEditContent] = useState('')
  const [editLoading, setEditLoading] = useState(false)
  const [editError, setEditError] = useState('')

  const [comment, setComment] = useState('')
  const [commentLoading, setCommentLoading] = useState(false)
  const [commentError, setCommentError] = useState('')

  const authHeaders = {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${token}`,
  }

  useEffect(() => {
    if (!postID) return
    setLoading(true)
    fetch(`${API_URL}/posts/${postID}`, { headers: authHeaders })
      .then(r => r.json())
      .then(json => {
        if (json.error) setError(json.error)
        else setPost(json.data)
      })
      .catch(() => setError('Error de conexión'))
      .finally(() => setLoading(false))
  }, [postID])

  const startEdit = () => {
    if (!post) return
    setEditTitle(post.title)
    setEditContent(post.content)
    setEditing(true)
    setEditError('')
  }

  const handleEdit = async (e: FormEvent) => {
    e.preventDefault()
    if (!post) return
    setEditLoading(true)
    setEditError('')
    try {
      const res = await fetch(`${API_URL}/posts/${post.id}`, {
        method: 'PATCH',
        headers: authHeaders,
        body: JSON.stringify({ title: editTitle, content: editContent }),
      })
      const json = await res.json()
      if (!res.ok) { setEditError(json.error || 'Error al actualizar'); return }
      setPost(json.data)
      setEditing(false)
    } catch {
      setEditError('Error de conexión')
    } finally {
      setEditLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!post || !confirm('¿Eliminar esta publicación?')) return
    try {
      const res = await fetch(`${API_URL}/posts/${post.id}`, {
        method: 'DELETE',
        headers: authHeaders,
      })
      if (res.ok || res.status === 204) navigate('/')
      else {
        const json = await res.json()
        alert(json.error || 'Error al eliminar')
      }
    } catch {
      alert('Error de conexión')
    }
  }

  const handleComment = async (e: FormEvent) => {
    e.preventDefault()
    if (!post || !comment.trim()) return
    setCommentLoading(true)
    setCommentError('')
    try {
      const res = await fetch(`${API_URL}/posts/${post.id}/comments`, {
        method: 'POST',
        headers: authHeaders,
        body: JSON.stringify({ content: comment }),
      })
      if (res.ok) {
        const json = await res.json()
        setPost(prev => prev ? {
          ...prev,
          comments: [...(prev.comments || []), json.data]
        } : prev)
        setComment('')
      } else {
        const json = await res.json()
        setCommentError(json.error || 'Error al publicar el comentario')
      }
    } catch {
      setCommentError('Error de conexión')
    } finally {
      setCommentLoading(false)
    }
  }

  if (loading) return <main className="page-content"><div className="page-loader"><div className="spinner" /></div></main>
  if (error) return <main className="page-content"><div className="error-msg">{error}</div></main>
  if (!post) return null

  const isOwner = userId === post.user_id

  return (
    <main className="page-content">
      <Link to="/" className="back-link">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M19 12H5M12 5l-7 7 7 7"/></svg>
        Volver al feed
      </Link>

      {editing ? (
        <form className="edit-form" onSubmit={handleEdit}>
          {editError && <div className="error-msg">{editError}</div>}
          <div className="form-field">
            <label>Título</label>
            <input value={editTitle} onChange={e => setEditTitle(e.target.value)} placeholder="Título de la publicación" />
          </div>
          <div className="form-field">
            <label>Contenido</label>
            <textarea
              value={editContent}
              onChange={e => setEditContent(e.target.value)}
              rows={8}
              placeholder="Contenido de la publicación"
            />
          </div>
          <div className="edit-form-actions">
            <button type="submit" className="btn-primary" disabled={editLoading}>
              {editLoading ? <span className="spinner" /> : 'Guardar cambios'}
            </button>
            <button type="button" className="btn-ghost" onClick={() => setEditing(false)}>
              Cancelar
            </button>
          </div>
        </form>
      ) : (
        <article>
          <div className="post-page-header">
            {post.tags?.length > 0 && (
              <div className="post-tags" style={{ marginBottom: '0.75rem' }}>
                {post.tags.map(tag => <span key={tag} className="tag">{tag}</span>)}
              </div>
            )}
            <h1 className="post-page-title">{post.title || '(Sin título)'}</h1>
            <div className="post-page-meta">
              <span className="post-meta-item" style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/></svg>
                <Link to={`/users/${post.user_id}`} style={{ color: 'inherit' }}>
                  {post.user?.username || `User #${post.user_id}`}
                </Link>
              </span>
              <span className="post-meta-item" style={{ fontSize: '0.875rem', color: 'var(--text-muted)' }}>
                {formatDate(post.created_at)}
              </span>
              {post.updated_at !== post.created_at && (
                <span className="post-meta-item" style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>
                  (editado {formatDate(post.updated_at)})
                </span>
              )}
            </div>
          </div>

          <div className="post-page-body">{post.content}</div>

          {isOwner && (
            <div className="post-actions">
              <button className="btn-secondary btn-sm" onClick={startEdit}>Editar</button>
              <button className="btn-danger btn-sm" onClick={handleDelete}>Eliminar</button>
            </div>
          )}
        </article>
      )}

      <div className="divider" />

      <section className="comments-section">
        <div className="comments-title">
          {(post.comments?.length || 0)} Comentario{post.comments?.length !== 1 ? 's' : ''}
        </div>

        {post.comments?.map(c => (
          <div key={c.id} className="comment">
            <div className="comment-author">
              <Link to={`/users/${c.user_id}`} className="comment-username">
                {c.user?.username || `User #${c.user_id}`}
              </Link>
              <span className="comment-date">{formatDate(c.created_at)}</span>
            </div>
            <div className="comment-body">{c.content}</div>
          </div>
        ))}

        <form className="comment-form" onSubmit={handleComment}>
          {commentError && <div className="error-msg">{commentError}</div>}
          <textarea
            value={comment}
            onChange={e => setComment(e.target.value)}
            rows={3}
            placeholder="Añadir un comentario…"
            required
          />
          <div>
            <button type="submit" className="btn-primary btn-sm" disabled={commentLoading}>
              {commentLoading ? <span className="spinner" /> : 'Comentar'}
            </button>
          </div>
        </form>
      </section>
    </main>
  )
}
