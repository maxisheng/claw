import { useEffect, useState } from 'react'
import api from '../api/axios'

function Dashboard() {
  const [stats, setStats] = useState({ articles: 0, categories: 0 })

  useEffect(() => {
    Promise.all([
      api.get('/articles'),
      api.get('/categories')
    ]).then(([articlesRes, categoriesRes]) => {
      setStats({
        articles: articlesRes.data.length,
        categories: categoriesRes.data.length
      })
    })
  }, [])

  return (
    <div className="dashboard">
      <h1>Dashboard</h1>
      <div className="stats-grid">
        <div className="stat-card">
          <h3>📝 Articles</h3>
          <p className="stat-number">{stats.articles}</p>
        </div>
        <div className="stat-card">
          <h3>📁 Categories</h3>
          <p className="stat-number">{stats.categories}</p>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
