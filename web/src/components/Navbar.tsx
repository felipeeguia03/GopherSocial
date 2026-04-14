import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function Navbar() {
  const { token, userId, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <nav className="navbar">
      <div className="navbar-inner">
        <NavLink to="/" className="navbar-brand">
          Gopher<span>Social</span>
        </NavLink>

        <ul className="navbar-links">
          {token ? (
            <>
              <li>
                <NavLink
                  to="/"
                  className={({ isActive }) => isActive ? 'nav-active' : ''}
                >
                  Inicio
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/search"
                  className={({ isActive }) => isActive ? 'nav-active' : ''}
                >
                  Buscar
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/posts/new"
                  className={({ isActive }) => `nav-cta${isActive ? ' nav-active' : ''}`}
                >
                  + Nueva publicación
                </NavLink>
              </li>
              {userId && (
                <li>
                  <NavLink
                    to={`/users/${userId}`}
                    className={({ isActive }) => isActive ? 'nav-active' : ''}
                  >
                    Perfil
                  </NavLink>
                </li>
              )}
              <li>
                <button onClick={handleLogout} className="btn-ghost">
                  Cerrar sesión
                </button>
              </li>
            </>
          ) : (
            <>
              <li>
                <NavLink to="/login">Iniciar sesión</NavLink>
              </li>
              <li>
                <NavLink to="/register" className="nav-cta">
                  Registrarse
                </NavLink>
              </li>
            </>
          )}
        </ul>
      </div>
    </nav>
  )
}
