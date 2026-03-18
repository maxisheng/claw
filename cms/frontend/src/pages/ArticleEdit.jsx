import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import api from '../api/axios'

function ArticleEdit() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [article, setArticle] = useState({
    title: '',
    content: '',
    summary: '',
    slug: '',
    status: 'draft',
    category_id: null
  })
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    api.get('/categories').then(res => setCategories(res.data))
    if (id) {
      api.get(`/articles/${id}`).then(res => setArticle(res.data))
    }
  }, [id])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    try {
      if (id) {
        await api.put(`/articles/${id}`, article)
      } else {
        await api.post('/articles', article)
      }
      navigate('/articles')
    } catch (error) {
      console.error('Failed to save article:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="article-edit">
      <h1>{id ? 'Edit Article' : 'New Article'}</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Title</label>
          <input
            type="text"
            value={article.title}
            onChange={e => setArticle({...article, title: e.target.value})}
            required
          />
        </div>
        <div className="form-group">
          <label>Slug</label>
          <input
            type="text"
            value={article.slug}
            onChange={e => setArticle({...article, slug: e.target.value})}
          />
        </div>
        <div className="form-group">
          <label>Summary</label>
          <textarea
            value={article.summary}
            onChange={e => setArticle({...article, summary: e.target.value})}
            rows="2"
          />
        </div>
        <div className="form-group">
          <label>Category</label>
          <select
            value={article.category_id || ''}
            onChange={e => setArticle({...article, category_id: parseInt(e.target.value)})}
          >
            <option value="">Select Category</option>
            {categories.map(cat => (
              <option key={cat.id} value={cat.id}>{cat.name}</option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label>Status</label>
          <select
            value={article.status}
            onChange={e => setArticle({...article, status: e.target.value})}
          >
            <option value="draft">Draft</option>
            <option value="published">Published</option>
            <option value="archived">Archived</option>
          </select>
        </div>
        <div className="form-group">
          <label>Content</label>
          <textarea
            value={article.content}
            onChange={e => setArticle({...article, content: e.target.value})}
            rows="15"
            required
          />
        </div>
        <div className="form-actions">
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Saving...' : 'Save'}
          </button>
          <button type="button" onClick={() => navigate('/articles')} className="btn-secondary">
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}

export default ArticleEdit
