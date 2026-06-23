import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { getAccessToken } from './api/auth'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import VerifyEmailPage from './pages/VerifyEmailPage'
import HomePage from './pages/HomePage'
import ProfilePage from './pages/ProfilePage'
import PostPage from './pages/PostPage'
import FriendsPage from './pages/FriendsPage'
import UserPage from './pages/UserPage'
import PeoplePage from './pages/PeoplePage'
import MessagesPage from './pages/MessagesPage'
import ChatPage from './pages/ChatPage'
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
        <Route path="/verify-email" element={<VerifyEmailPage />} />
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
        <Route
          path="/people"
          element={
            <PrivateRoute>
              <PeoplePage />
            </PrivateRoute>
          }
        />
        <Route
          path="/messages"
          element={<PrivateRoute><MessagesPage /></PrivateRoute>}
        />
        <Route
          path="/messages/:id"
          element={<PrivateRoute><ChatPage /></PrivateRoute>}
        />
        {/* профиль другого юзера — доступен и без авторизации */}
        <Route path="/users/:id" element={<UserPage />} />
      </Routes>
    </BrowserRouter>
  )
}
