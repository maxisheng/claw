import { useEffect, useState } from 'react'
import api from '../api/axios'

function Categories() {
  const [categories, setCategories] = useState([])
  const [newCategory, setNewCategory] = useState({ name: '', slug: '', description: '' })
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    loadCategories()
  }, [])

  const loadCategories = async () => {
    const response = await api.get('/categories')
    setCategories(response.data)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    await api.post('/categories', newCategory)
    setNewCategory({ name: '', slug: '', description: '' })
    setShowForm(false)
    loadCategories()
  }

  const handleDelete = async (id) => {
    if (confirm('Are you sure?')) {
      await api.delete(`/categories/${id}`)
      loadCategories()
    }
  }

  return (
    <div className="categories-page">
      <div className="page-header">
        <h1>Categories</h1>
        <button onClick={() => setShowForm(!showForm)} className="btn-primary">
          {showForm ? 'Cancel' : '+ New Category'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="category-form">
          <div className="form-group">
            <label>Name</label>
            <input
              type="text"
              value={newCategory.name}
              onChange={e => setNewCategory({...newCategory, name: e.target.value})}
              required
            />
          </div>
          <div className="form-group">
            <label>Slug</label>
            <input
              type="text"
              value={newCategory.slug}
              onChange={e => setNewCategory({...newCategory, slug: e.target.value})}
            />
          </div>
          <div className="form-group">
            <label>Description</label>
            <textarea
              value={newCategory.description}
              onChange={e => setNewCategory({...newCategory, description: e.target.value})}
              rows="2"
            />
          </div>
          <button type="submit" className="btn-primary">Create</button>
        </form>
      )}

      <table className="data-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Slug</th>
            <th>Description</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {categories.map(cat => (
            <tr key={cat.id}>
              <td>{cat.name}</td>
              <td>{cat.slug}</td>
              <td>{cat.description}</td>
              <td>
                <button onClick={() => handleDelete(cat.id)} className="btn-sm btn-danger">Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Categories
