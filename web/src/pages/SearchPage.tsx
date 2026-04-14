import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { API_URL } from '../App'

interface User {
  id: number
  username: string
  email: string
  is_active: boolean
  created_at: string
}

export function SearchPage() {
  const { token, userId: myId } = useAuth()
  const navigate = useNavigate()
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  // track follow state per user id
  const [followed, setFollowed] = useState<Record<number, boolean>>({})
  const [followLoading, setFollowLoading] = useState<Record<number, boolean>>({})
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const authHeaders = {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${token}`,
  }

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current)
    if (!query.trim()) {
      setResults([])
      return
    }
    debounceRef.current = setTimeout(async () => {
      setLoading(true)
      setError('')
      try {
        const res = await fetch(`${API_URL}/users/search?q=${encodeURIComponent(query)}`, {
          headers: authHeaders,
        })
        const json = await res.json()
        if (!res.ok) {
          setError(json.error || 'Error buscando usuarios')
          setResults([])
        } else {
          setResults(json.data ?? [])
        }
      } catch {
        setError('Error de conexión')
      } finally {
        setLoading(false)
      }
    }, 350)
  }, [query])

  const handleFollow = async (user: User) => {
    if (!myId) return
    const isFollowing = followed[user.id]
    const action = isFollowing ? 'unfollow' : 'follow'
    setFollowLoading(prev => ({ ...prev, [user.id]: true }))
    try {
      const res = await fetch(`${API_URL}/users/${user.id}/${action}`, {
        method: 'PUT',
        headers: authHeaders,
      })
      if (res.ok || res.status === 204) {
        setFollowed(prev => ({ ...prev, [user.id]: !isFollowing }))
      }
    } finally {
      setFollowLoading(prev => ({ ...prev, [user.id]: false }))
    }
  }

  const initials = (name: string) => name.slice(0, 2).toUpperCase()

  return (
    <main className="page-content">
      <div className="feed-header">
        <h1>Buscar usuarios</h1>
      </div>

      <div className="feed-filters">
        <input
          autoFocus
          type="text"
          placeholder="Buscar por nombre de usuario…"
          value={query}
          onChange={e => setQuery(e.target.value)}
        />
        {query && (
          <button className="btn-ghost" onClick={() => { setQuery(''); setResults([]) }}>
            Limpiar
          </button>
        )}
      </div>

      {error && <div className="error-msg">{error}</div>}

      {loading && (
        <div className="page-loader"><div className="spinner" /></div>
      )}

      {!loading && query && results.length === 0 && !error && (
        <div className="empty-state">
          <h3>No se encontraron usuarios</h3>
          <p>Prueba con otro nombre</p>
        </div>
      )}

      {results.length > 0 && (
        <div className="post-list">
          {results.map(user => {
            const isMe = user.id === myId
            const isFollowing = followed[user.id] ?? false
            return (
              <div key={user.id} className="profile-card" style={{ cursor: 'default' }}>
                <div
                  style={{ display: 'flex', alignItems: 'center', gap: '1rem', flex: 1, cursor: 'pointer' }}
                  onClick={() => navigate(`/users/${user.id}`)}
                >
                  <div className="profile-avatar">{initials(user.username)}</div>
                  <div className="profile-info">
                    <div className="profile-username">{user.username}</div>
                    <div className="profile-email">{user.email}</div>
                  </div>
                </div>

                {!isMe && (
                  <button
                    className={isFollowing ? 'btn-secondary' : 'btn-primary'}
                    onClick={() => handleFollow(user)}
                    disabled={followLoading[user.id]}
                    style={{ flexShrink: 0 }}
                  >
                    {followLoading[user.id]
                      ? <span className="spinner" />
                      : isFollowing ? 'Dejar de seguir' : 'Seguir'
                    }
                  </button>
                )}
              </div>
            )
          })}
        </div>
      )}
    </main>
  )
}
