import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { Table, Button, Space, Tag, Typography, Popconfirm, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../api/axios'

const { Title } = Typography

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
      message.error('Failed to load articles')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id) => {
    try {
      await api.delete(`/articles/${id}`)
      message.success('Article deleted successfully')
      loadArticles()
    } catch (error) {
      console.error('Failed to delete article:', error)
      message.error('Failed to delete article')
    }
  }

  const statusColors = {
    draft: 'default',
    published: 'success',
    archived: 'default',
  }

  const columns = [
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      render: (text, record) => (
        <Link to={`/articles/${record.id}`}>{text}</Link>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={statusColors[status]}>{status.toUpperCase()}</Tag>
      ),
    },
    {
      title: 'Author',
      dataIndex: ['author', 'username'],
      key: 'author',
    },
    {
      title: 'Category',
      dataIndex: ['category', 'name'],
      key: 'category',
    },
    {
      title: 'Created',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => new Date(date).toLocaleDateString(),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space size="small">
          <Link to={`/articles/${record.id}`}>
            <Button type="link" icon={<EditOutlined />}>Edit</Button>
          </Link>
          <Popconfirm
            title="Are you sure to delete this article?"
            onConfirm={() => handleDelete(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>Delete</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={2} style={{ margin: 0 }}>Articles</Title>
        <Link to="/articles/new">
          <Button type="primary" icon={<PlusOutlined />}>New Article</Button>
        </Link>
      </div>
      <Table
        columns={columns}
        dataSource={articles}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10 }}
      />
    </div>
  )
}

export default Articles
