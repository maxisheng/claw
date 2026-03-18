import { useEffect, useState } from 'react'
import { Card, Row, Col, Statistic, Typography, Table } from 'antd'
import { FileTextOutlined, FolderOutlined, UserOutlined, EyeOutlined } from '@ant-design/icons'
import api from '../api/axios'

const { Title } = Typography

function Dashboard() {
  const [stats, setStats] = useState({
    articles: 0,
    categories: 0,
    totalViews: 0
  })
  const [recentArticles, setRecentArticles] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get('/articles'),
      api.get('/categories')
    ]).then(([articlesRes, categoriesRes]) => {
      const articles = articlesRes || []
      const categories = categoriesRes || []
      
      setStats({
        articles: articles.length,
        categories: categories.length,
        totalViews: articles.reduce((sum, a) => sum + (a.view_count || 0), 0)
      })
      setRecentArticles(articles.slice(0, 5))
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  const articleColumns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => {
        const colors = { draft: 'default', published: 'success', archived: 'default' }
        return <span style={{ color: colors[status] || '#666' }}>{status}</span>
      },
    },
    {
      title: '浏览',
      dataIndex: 'view_count',
      key: 'view_count',
    },
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => new Date(date).toLocaleDateString('zh-CN'),
    },
  ]

  return (
    <div>
      <Title level={2}>控制台</Title>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={loading}>
            <Statistic
              title="文章总数"
              value={stats.articles}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={loading}>
            <Statistic
              title="分类数量"
              value={stats.categories}
              prefix={<FolderOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={loading}>
            <Statistic
              title="总浏览量"
              value={stats.totalViews}
              prefix={<EyeOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card loading={loading}>
            <Statistic
              title="管理员"
              value={1}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Title level={4} style={{ marginTop: 24 }}>最近文章</Title>
      <Table
        columns={articleColumns}
        dataSource={recentArticles}
        rowKey="id"
        pagination={false}
        loading={loading}
      />
    </div>
  )
}

export default Dashboard
