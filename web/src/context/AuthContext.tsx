import { createContext, useContext, useState, ReactNode } from 'react'

interface AuthState {
  token: string | null
  userId: number | null
}

interface AuthContextType extends AuthState {
  login: (token: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

function decodeUserId(token: string): number | null {
  try {
    const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')))
    const sub = payload.sub
    const id = typeof sub === 'string' ? parseInt(sub, 10) : sub
    return isNaN(id) ? null : id
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [auth, setAuth] = useState<AuthState>(() => {
    const token = localStorage.getItem('token')
    return { token, userId: token ? decodeUserId(token) : null }
  })

  const login = (token: string) => {
    localStorage.setItem('token', token)
    setAuth({ token, userId: decodeUserId(token) })
  }

  const logout = () => {
    localStorage.removeItem('token')
    setAuth({ token: null, userId: null })
  }

  return (
    <AuthContext.Provider value={{ ...auth, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
