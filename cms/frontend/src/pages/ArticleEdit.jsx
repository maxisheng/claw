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
    api.get('/categories').then(res => setCategories(res || []))
    if (id) {
      setLoading(true)
      api.get(`/articles/${id}`)
        .then(res => {
          form.setFieldsValue(res)
        })
        .finally(() => setLoading(false))
    }
  }, [id, form])

  const handleSubmit = async (values) => {
    setSubmitting(true)
    try {
      if (id) {
        await api.put(`/articles/${id}`, values)
        message.success('文章更新成功')
      } else {
        await api.post('/articles', values)
        message.success('文章创建成功')
      }
      navigate('/articles')
    } catch (error) {
      console.error('Failed to save article:', error)
      message.error(error.message || '保存文章失败')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Title level={2}>{id ? '编辑文章' : '新建文章'}</Title>
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
            label="标题"
            rules={[{ required: true, message: '请输入文章标题' }]}
          >
            <Input placeholder="请输入文章标题" />
          </Form.Item>

          <Form.Item
            name="slug"
            label="URL 标识"
            rules={[{ required: true, message: '请输入 URL 标识' }]}
          >
            <Input placeholder="例如：my-first-article" />
          </Form.Item>

          <Form.Item
            name="summary"
            label="摘要"
          >
            <TextArea rows={2} placeholder="文章摘要简介" />
          </Form.Item>

          <Form.Item
            name="category_id"
            label="分类"
          >
            <Select placeholder="请选择分类" allowClear>
              {categories.map(cat => (
                <Select.Option key={cat.id} value={cat.id}>{cat.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="status"
            label="状态"
            initialValue="draft"
          >
            <Select>
              <Select.Option value="draft">草稿</Select.Option>
              <Select.Option value="published">已发布</Select.Option>
              <Select.Option value="archived">已归档</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入文章内容' }]}
          >
            <TextArea rows={15} placeholder="请输入文章内容..." />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button 
                type="primary" 
                htmlType="submit" 
                icon={<SaveOutlined />}
                loading={submitting}
              >
                保存
              </Button>
              <Button 
                icon={<UndoOutlined />} 
                onClick={() => navigate('/articles')}
              >
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default ArticleEdit
