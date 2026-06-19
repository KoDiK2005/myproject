import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { getAccessToken } from './api/auth'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import HomePage from './pages/HomePage'
import ProfilePage from './pages/ProfilePage'
import PostPage from './pages/PostPage'
import FriendsPage from './pages/FriendsPage'
import UserPage from './pages/UserPage'
import './index.css'

// Защищённый роут — если нет токена, шлём на /login
function PrivateRoute({ children }) {
  return getAccessToken() ? children : <Navigate to="/login" replace />
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <HomePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <PrivateRoute>
              <ProfilePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/posts/:id"
          element={
            <PrivateRoute>
              <PostPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/friends"
          element={
            <PrivateRoute>
              <FriendsPage />
            </PrivateRoute>
          }
        />
        {/* профиль другого юзера — доступен и без авторизации */}
        <Route path="/users/:id" element={<UserPage />} />
      </Routes>
    </BrowserRouter>
  )
}
