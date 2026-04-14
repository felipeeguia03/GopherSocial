import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider, Outlet } from 'react-router-dom'
import './index.css'
import { AuthProvider } from './context/AuthContext'
import { Navbar } from './components/Navbar'
import { ProtectedRoute } from './components/ProtectedRoute'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { FeedPage } from './pages/FeedPage'
import { PostPage } from './pages/PostPage'
import { NewPostPage } from './pages/NewPostPage'
import { ProfilePage } from './pages/ProfilePage'
import { ConfirmationPage } from './ConfirmationPage'
import { SearchPage } from './pages/SearchPage'

function Layout() {
  return (
    <div className="layout">
      <Navbar />
      <Outlet />
    </div>
  )
}

const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      {
        path: '/',
        element: <ProtectedRoute><FeedPage /></ProtectedRoute>,
      },
      {
        path: '/posts/new',
        element: <ProtectedRoute><NewPostPage /></ProtectedRoute>,
      },
      {
        path: '/posts/:postID',
        element: <ProtectedRoute><PostPage /></ProtectedRoute>,
      },
      {
        path: '/users/:userID',
        element: <ProtectedRoute><ProfilePage /></ProtectedRoute>,
      },
      {
        path: '/search',
        element: <ProtectedRoute><SearchPage /></ProtectedRoute>,
      },
      {
        path: '/confirm/:token',
        element: <ConfirmationPage />,
      },
      {
        path: '/login',
        element: <LoginPage />,
      },
      {
        path: '/register',
        element: <RegisterPage />,
      },
    ],
  },
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  </StrictMode>,
)
