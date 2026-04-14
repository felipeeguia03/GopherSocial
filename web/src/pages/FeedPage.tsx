import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { API_URL } from '../App'

interface Post {
  id: number
  title: string
  content: string
  user_id: number
  tags: string[]
  comments_count: number
  created_at: string
  user: { id: number; username: string }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('es-ES', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

const LIMIT = 10

export function FeedPage() {
  const { token } = useAuth()
  const navigate = useNavigate()
  const [posts, setPosts] = useState<Post[]>([])
  const [offset, setOffset] = useState(0)
  const [hasMore, setHasMore] = useState(true)
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState('')
  const [search, setSearch] = useState('')
  const [searchInput, setSearchInput] = useState('')

  const fetchPosts = useCallback(async (currentOffset: number, currentSearch: string, replace: boolean) => {
    const setter = replace ? setLoading : setLoadingMore
    setter(true)
    setError('')
    try {
      const params = new URLSearchParams({
        limit: String(LIMIT),
        offset: String(currentOffset),
        sort: 'desc',
      })
      if (currentSearch) params.set('search', currentSearch)

      const res = await fetch(`${API_URL}/feed?${params}`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) {
        const json = await res.json()
        setError(json.error || 'Error al cargar el feed')
        return
      }
      const json = await res.json()
      const data: Post[] = json.data || []
      if (replace) {
        setPosts(data)
      } else {
        setPosts(prev => [...prev, ...data])
      }
      setHasMore(data.length === LIMIT)
    } catch {
      setError('Error de conexión')
    } finally {
      setter(false)
    }
  }, [token])

  useEffect(() => {
    setOffset(0)
    fetchPosts(0, search, true)
  }, [search, fetchPosts])

  const handleLoadMore = () => {
    const next = offset + LIMIT
    setOffset(next)
    fetchPosts(next, search, false)
  }

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setSearch(searchInput)
  }

  return (
    <main className="page-content">
      <div className="feed-header">
        <h1>Feed</h1>
      </div>

      <form className="feed-filters" onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Buscar publicaciones…"
          value={searchInput}
          onChange={e => setSearchInput(e.target.value)}
        />
        <button type="submit" className="btn-secondary">Buscar</button>
        {search && (
          <button
            type="button"
            className="btn-ghost"
            onClick={() => { setSearchInput(''); setSearch('') }}
          >
            Limpiar
          </button>
        )}
      </form>

      {error && <div className="error-msg">{error}</div>}

      {loading ? (
        <div className="page-loader"><div className="spinner" /></div>
      ) : posts.length === 0 ? (
        <div className="empty-state">
          <h3>{search ? 'No se encontraron publicaciones' : 'Tu feed está vacío'}</h3>
          <p>{search ? 'Prueba con otro término' : 'Sigue usuarios o crea tu primera publicación'}</p>
        </div>
      ) : (
        <>
          <div className="post-list">
            {posts.map((post, i) => (
              <div
                key={post.id}
                className="post-card"
                style={{ animationDelay: `${i * 0.05}s` }}
                onClick={() => navigate(`/posts/${post.id}`)}
              >
                {post.tags?.length > 0 && (
                  <div className="post-tags">
                    {post.tags.map(tag => (
                      <span key={tag} className="tag">{tag}</span>
                    ))}
                  </div>
                )}
                <div className="post-card-title">{post.title || '(Sin título)'}</div>
                {post.content && (
                  <div className="post-card-content">{post.content}</div>
                )}
                <div className="post-card-meta">
                  <span className="post-meta-item">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/></svg>
                    {post.user?.username || `User #${post.user_id}`}
                  </span>
                  <span className="post-meta-item">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
                    {post.comments_count ?? 0}
                  </span>
                  <span className="post-meta-item">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                    {formatDate(post.created_at)}
                  </span>
                </div>
              </div>
            ))}
          </div>

          {hasMore && (
            <div className="load-more">
              <button
                className="btn-secondary"
                onClick={handleLoadMore}
                disabled={loadingMore}
              >
                {loadingMore ? <><span className="spinner" /> Cargando…</> : 'Cargar más'}
              </button>
            </div>
          )}
        </>
      )}
    </main>
  )
}
