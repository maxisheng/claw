import { Outlet, Navigate, Link, useLocation } from 'react-router-dom'
import { useEffect, useState } from 'react'

function App() {
  const [user, setUser] = useState(null)
  const location = useLocation()

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')
    if (token && userData) {
      setUser(JSON.parse(userData))
    }
  }, [])

  if (!user && location.pathname !== '/login') {
    return <Navigate to="/login" replace />
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
  }

  return (
    <div className="app">
      <aside className="sidebar">
        <h1>CMS Admin</h1>
        <nav>
          <Link to="/">📊 Dashboard</Link>
          <Link to="/articles">📝 Articles</Link>
          <Link to="/categories">📁 Categories</Link>
        </nav>
        <div className="user-info">
          {user && <span>{user.username}</span>}
          <button onClick={handleLogout}>Logout</button>
        </div>
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  )
}

export default App
