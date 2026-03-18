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
      setArticles(response)
    } catch (error) {
      console.error('Failed to load articles:', error)
      message.error('加载文章失败')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id) => {
    try {
      await api.delete(`/articles/${id}`)
      message.success('文章已删除')
      loadArticles()
    } catch (error) {
      console.error('Failed to delete article:', error)
      message.error('删除文章失败')
    }
  }

  const statusColors = {
    draft: 'default',
    published: 'success',
    archived: 'default',
  }

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (text, record) => (
        <Link to={`/articles/${record.id}`}>{text}</Link>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={statusColors[status]}>{status === 'published' ? '已发布' : status === 'draft' ? '草稿' : '已归档'}</Tag>
      ),
    },
    {
      title: '作者',
      dataIndex: ['author', 'username'],
      key: 'author',
    },
    {
      title: '分类',
      dataIndex: ['category', 'name'],
      key: 'category',
    },
    {
      title: '浏览量',
      dataIndex: 'view_count',
      key: 'view_count',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => new Date(date).toLocaleDateString('zh-CN'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record) => (
        <Space size="small">
          <Link to={`/articles/${record.id}`}>
            <Button type="link" icon={<EditOutlined />}>编辑</Button>
          </Link>
          <Popconfirm
            title="确定要删除这篇文章吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Title level={2} style={{ margin: 0 }}>文章管理</Title>
        <Link to="/articles/new">
          <Button type="primary" icon={<PlusOutlined />}>新建文章</Button>
        </Link>
      </div>
      <Table
        columns={columns}
        dataSource={articles}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10, showSizeChanger: true }}
      />
    </div>
  )
}

export default Articles
