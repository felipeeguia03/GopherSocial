import { useState, FormEvent } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { API_URL } from '../App'

export function NewPostPage() {
  const { token, userId } = useAuth()
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [tagsInput, setTagsInput] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    const tags = tagsInput
      .split(',')
      .map(t => t.trim())
      .filter(Boolean)
    try {
      const res = await fetch(`${API_URL}/posts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ title, content, user_id: userId, tags }),
      })
      const json = await res.json()
      if (!res.ok) {
        setError(json.error || 'Error al crear la publicación')
        return
      }
      navigate(`/posts/${json.data.id}`)
    } catch {
      setError('Error de conexión')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page-content">
      <Link to="/" className="back-link">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M19 12H5M12 5l-7 7 7 7"/></svg>
        Volver al feed
      </Link>

      <div className="new-post-card">
        <h1 className="new-post-title">Nueva publicación</h1>

        {error && <div className="error-msg">{error}</div>}

        <form className="new-post-form" onSubmit={handleSubmit}>
          <div className="form-field">
            <label htmlFor="title">Título</label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder="Un título interesante…"
              autoFocus
            />
          </div>
          <div className="form-field">
            <label htmlFor="content">Contenido</label>
            <textarea
              id="content"
              value={content}
              onChange={e => setContent(e.target.value)}
              rows={10}
              placeholder="Escribe tu publicación aquí…"
            />
          </div>
          <div className="form-field">
            <label htmlFor="tags">Etiquetas</label>
            <input
              id="tags"
              type="text"
              value={tagsInput}
              onChange={e => setTagsInput(e.target.value)}
              placeholder="golang, api, tutorial"
            />
            <div className="tags-hint">Separa las etiquetas con comas</div>
          </div>
          <div>
            <button type="submit" className="btn-primary" disabled={loading}>
              {loading ? <span className="spinner" /> : 'Publicar'}
            </button>
          </div>
        </form>
      </div>
    </main>
  )
}
