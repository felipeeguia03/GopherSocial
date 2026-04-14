import { useState, FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { API_URL } from '../App'

export function RegisterPage() {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password }),
      })
      const json = await res.json()
      if (!res.ok) {
        setError(json.error || 'Error al registrarse')
        return
      }
      setSuccess(true)
    } catch {
      setError('Error de conexión. ¿Está el servidor activo?')
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="auth-page">
        <div className="auth-card">
          <div className="auth-header">
            <h1>Revisa tu correo</h1>
            <p>Enviamos un enlace de activación a <strong>{email}</strong></p>
          </div>
          <div className="success-msg">
            ¡Cuenta creada! Haz clic en el enlace del correo para activar tu cuenta.
          </div>
          <div className="auth-footer">
            ¿Ya la activaste? <Link to="/login">Iniciar sesión</Link>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <div className="auth-header">
          <h1>Crear cuenta</h1>
          <p>Únete a Gopher Social hoy</p>
        </div>

        {error && <div className="error-msg">{error}</div>}

        <form className="auth-form" onSubmit={handleSubmit}>
          <div className="form-field">
            <label htmlFor="username">Nombre de usuario</label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={e => setUsername(e.target.value)}
              placeholder="johndoe"
              required
              autoFocus
            />
          </div>
          <div className="form-field">
            <label htmlFor="email">Correo electrónico</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
            />
          </div>
          <div className="form-field">
            <label htmlFor="password">Contraseña</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              placeholder="••••••••"
              minLength={3}
              required
            />
          </div>
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? <span className="spinner" /> : 'Crear cuenta'}
          </button>
        </form>

        <div className="auth-footer">
          ¿Ya tienes cuenta? <Link to="/login">Iniciar sesión</Link>
        </div>
      </div>
    </div>
  )
}
