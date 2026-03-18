import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import api from '../api/axios'

function Articles() {
  const [articles, setArticles] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadArticles()
  }, [])

  const loadArticles = async () => {
    try {
      const response = await api.get('/articles')
      setArticles(response.data)
    } catch (error) {
      console.error('Failed to load articles:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id) => {
    if (confirm('Are you sure you want to delete this article?')) {
      try {
        await api.delete(`/articles/${id}`)
        loadArticles()
      } catch (error) {
        console.error('Failed to delete article:', error)
      }
    }
  }

  if (loading) return <div>Loading...</div>

  return (
    <div className="articles-page">
      <div className="page-header">
        <h1>Articles</h1>
        <Link to="/articles/new" className="btn-primary">+ New Article</Link>
      </div>
      <table className="data-table">
        <thead>
          <tr>
            <th>Title</th>
            <th>Status</th>
            <th>Author</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {articles.map(article => (
            <tr key={article.id}>
              <td>{article.title}</td>
              <td><span className={`status-badge ${article.status}`}>{article.status}</span></td>
              <td>{article.author?.username || 'Unknown'}</td>
              <td>{new Date(article.created_at).toLocaleDateString()}</td>
              <td>
                <Link to={`/articles/${article.id}`} className="btn-sm">Edit</Link>
                <button onClick={() => handleDelete(article.id)} className="btn-sm btn-danger">Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Articles
