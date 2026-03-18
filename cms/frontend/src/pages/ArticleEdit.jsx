import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Form, Input, Select, Button, Space, Card, message, Typography } from 'antd'
import { SaveOutlined, UndoOutlined } from '@ant-design/icons'
import api from '../api/axios'

const { Title } = Typography
const { TextArea } = Input

function ArticleEdit() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    api.get('/categories').then(res => setCategories(res.data))
    if (id) {
      setLoading(true)
      api.get(`/articles/${id}`)
        .then(res => {
          form.setFieldsValue(res.data)
        })
        .finally(() => setLoading(false))
    }
  }, [id, form])

  const handleSubmit = async (values) => {
    setSubmitting(true)
    try {
      if (id) {
        await api.put(`/articles/${id}`, values)
        message.success('Article updated successfully')
      } else {
        await api.post('/articles', values)
        message.success('Article created successfully')
      }
      navigate('/articles')
    } catch (error) {
      console.error('Failed to save article:', error)
      message.error('Failed to save article')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Title level={2}>{id ? 'Edit Article' : 'New Article'}</Title>
      </div>
      <Card loading={loading}>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          size="large"
        >
          <Form.Item
            name="title"
            label="Title"
            rules={[{ required: true, message: 'Please input the title!' }]}
          >
            <Input placeholder="Enter article title" />
          </Form.Item>

          <Form.Item
            name="slug"
            label="Slug"
            rules={[{ required: true, message: 'Please input the slug!' }]}
          >
            <Input placeholder="e.g., my-first-article" />
          </Form.Item>

          <Form.Item
            name="summary"
            label="Summary"
          >
            <TextArea rows={2} placeholder="Brief summary of the article" />
          </Form.Item>

          <Form.Item
            name="category_id"
            label="Category"
          >
            <Select placeholder="Select a category" allowClear>
              {categories.map(cat => (
                <Select.Option key={cat.id} value={cat.id}>{cat.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="status"
            label="Status"
            initialValue="draft"
          >
            <Select>
              <Select.Option value="draft">Draft</Select.Option>
              <Select.Option value="published">Published</Select.Option>
              <Select.Option value="archived">Archived</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="content"
            label="Content"
            rules={[{ required: true, message: 'Please input the content!' }]}
          >
            <TextArea rows={15} placeholder="Write your article content here..." />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button 
                type="primary" 
                htmlType="submit" 
                icon={<SaveOutlined />}
                loading={submitting}
              >
                Save
              </Button>
              <Button 
                icon={<UndoOutlined />} 
                onClick={() => navigate('/articles')}
              >
                Cancel
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ArticleEdit
