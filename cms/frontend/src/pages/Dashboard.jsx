import { useEffect, useState } from 'react'
import { Card, Row, Col, Statistic, Typography } from 'antd'
import { FileTextOutlined, FolderOutlined, UserOutlined } from '@ant-design/icons'
import api from '../api/axios'

const { Title } = Typography

function Dashboard() {
  const [stats, setStats] = useState({ articles: 0, categories: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.get('/articles'),
      api.get('/categories')
    ]).then(([articlesRes, categoriesRes]) => {
      setStats({
        articles: articlesRes.data.length,
        categories: categoriesRes.data.length
      })
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  return (
    <div>
      <Title level={2}>Dashboard</Title>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={8}>
          <Card loading={loading}>
            <Statistic
              title="Total Articles"
              value={stats.articles}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card loading={loading}>
            <Statistic
              title="Total Categories"
              value={stats.categories}
              prefix={<FolderOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card loading={loading}>
            <Statistic
              title="Current User"
              value={1}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
