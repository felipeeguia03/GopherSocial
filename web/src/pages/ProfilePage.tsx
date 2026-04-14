import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { API_URL } from '../App'

interface User {
  id: number
  username: string
  email: string
  is_active: boolean
  created_at: string
  role: { name: string }
}

interface Post {
  id: number
  title: string
  content: string
  tags: string[]
  comments_count: number
  created_at: string
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('es-ES', {
    month: 'long', day: 'numeric', year: 'numeric',
  })
}

function formatShort(iso: string) {
  return new Date(iso).toLocaleDateString('es-ES', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
}

export function ProfilePage() {
  const { userID } = useParams<{ userID: string }>()
  const { token, userId: myId } = useAuth()
  const navigate = useNavigate()
  const [user, setUser] = useState<User | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [postsLoading, setPostsLoading] = useState(true)
  const [error, setError] = useState('')
  const [following, setFollowing] = useState(false)
  const [followLoading, setFollowLoading] = useState(false)
  const [followError, setFollowError] = useState('')

  const authHeaders = {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${token}`,
  }

  useEffect(() => {
    if (!userID) return
    setLoading(true)
    setPostsLoading(true)

    fetch(`${API_URL}/users/${userID}`, { headers: authHeaders })
      .then(r => r.json())
      .then(json => {
        if (json.error) setError(json.error)
        else setUser(json.data)
      })
      .catch(() => setError('Error de conexión'))
      .finally(() => setLoading(false))

    fetch(`${API_URL}/users/${userID}/posts`, { headers: authHeaders })
      .then(r => r.json())
      .then(json => setPosts(json.data ?? []))
      .catch(() => {})
      .finally(() => setPostsLoading(false))
  }, [userID])

  const handleFollow = async () => {
    if (!user || !myId) return
    setFollowLoading(true)
    setFollowError('')
    const action = following ? 'unfollow' : 'follow'
    try {
      const res = await fetch(`${API_URL}/users/${user.id}/${action}`, {
        method: 'PUT',
        headers: authHeaders,
      })
      if (res.ok || res.status === 204) {
        setFollowing(!following)
      } else {
        const json = await res.json()
        setFollowError(json.error || 'Error al realizar la acción')
      }
    } catch {
      setFollowError('Error de conexión')
    } finally {
      setFollowLoading(false)
    }
  }

  if (loading) return <main className="page-content"><div className="page-loader"><div className="spinner" /></div></main>
  if (error) return <main className="page-content"><div className="error-msg">{error}</div></main>
  if (!user) return null

  const isMe = myId === user.id
  const initials = user.username.slice(0, 2).toUpperCase()

  return (
    <main className="page-content">
      {/* Profile card */}
      <div className="profile-card">
        <div style={{ display: 'flex', alignItems: 'center', gap: '1.25rem', flex: 1 }}>
          <div className="profile-avatar">{initials}</div>
          <div className="profile-info">
            <div className="profile-username">{user.username}</div>
            <div className="profile-email">{user.email}</div>
            <div style={{ marginTop: '0.5rem', display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
              {user.is_active && (
                <span className="profile-badge">
                  <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                  Activo
                </span>
              )}
              {user.role?.name && <span className="tag">{user.role.name}</span>}
            </div>
          </div>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: '0.5rem' }}>
          <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>
            Miembro desde {formatDate(user.created_at)}
          </div>
          {!isMe && (
            <button
              className={following ? 'btn-secondary' : 'btn-primary'}
              onClick={handleFollow}
              disabled={followLoading}
            >
              {followLoading ? <span className="spinner" /> : following ? 'Dejar de seguir' : 'Seguir'}
            </button>
          )}
        </div>
      </div>

      {followError && <div className="error-msg">{followError}</div>}

      {/* Posts section */}
      <div style={{ marginTop: '2rem' }}>
        <h2 style={{ fontSize: '1.1rem', fontWeight: 600, marginBottom: '1rem', color: 'var(--text-secondary)' }}>
          {isMe ? 'Tus publicaciones' : `Publicaciones de ${user.username}`}
          {!postsLoading && (
            <span style={{ marginLeft: '0.5rem', fontSize: '0.85rem', color: 'var(--text-muted)' }}>
              ({posts.length})
            </span>
          )}
        </h2>

        {postsLoading ? (
          <div className="page-loader"><div className="spinner" /></div>
        ) : posts.length === 0 ? (
          <div className="empty-state">
            <h3>{isMe ? 'Aún no hay publicaciones' : 'Este usuario no tiene publicaciones'}</h3>
            {isMe && <p><a href="/posts/new">Crea tu primera publicación</a></p>}
          </div>
        ) : (
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
                    {post.tags.map(tag => <span key={tag} className="tag">{tag}</span>)}
                  </div>
                )}
                <div className="post-card-title">{post.title || '(Sin título)'}</div>
                {post.content && <div className="post-card-content">{post.content}</div>}
                <div className="post-card-meta">
                  <span className="post-meta-item">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
                    {post.comments_count ?? 0}
                  </span>
                  <span className="post-meta-item">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                    {formatShort(post.created_at)}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </main>
  )
}
