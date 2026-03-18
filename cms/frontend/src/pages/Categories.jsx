import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Typography, Popconfirm, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../api/axios'

const { Title } = Typography

function Categories() {
  const [categories, setCategories] = useState([])
  const [loading, setLoading] = useState(true)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [editingId, setEditingId] = useState(null)

  useEffect(() => {
    loadCategories()
  }, [])

  const loadCategories = async () => {
    try {
      const response = await api.get('/categories')
      setCategories(response || [])
    } catch (error) {
      console.error('Failed to load categories:', error)
      message.error('加载分类失败')
    } finally {
      setLoading(false)
    }
  }

  const handleOpenModal = (record = null) => {
    if (record) {
      setEditingId(record.id)
      form.setFieldsValue(record)
    } else {
      setEditingId(null)
      form.resetFields()
    }
    setModalVisible(true)
  }

  const handleSubmit = async (values) => {
    try {
      if (editingId) {
        await api.put(`/categories/${editingId}`, values)
        message.success('分类更新成功')
      } else {
        await api.post('/categories', values)
        message.success('分类创建成功')
      }
      setModalVisible(false)
      loadCategories()
    } catch (error) {
      console.error('Failed to save category:', error)
      message.error(error.message || '保存分类失败')
    }
  }

  const handleDelete = async (id) => {
    try {
      await api.delete(`/categories/${id}`)
      message.success('分类已删除')
      loadCategories()
    } catch (error) {
      console.error('Failed to delete category:', error)
      message.error('删除分类失败')
    }
  }

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'URL 标识',
      dataIndex: 'slug',
      key: 'slug',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '排序',
      dataIndex: 'sort_order',
      key: 'sort_order',
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
          <Button 
            type="link" 
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个分类吗？"
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
        <Title level={2} style={{ margin: 0 }}>分类管理</Title>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          新建分类
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={categories}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        title={editingId ? '编辑分类' : '新建分类'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入分类名称' }]}
          >
            <Input placeholder="分类名称" />
          </Form.Item>

          <Form.Item
            name="slug"
            label="URL 标识"
            rules={[{ required: true, message: '请输入 URL 标识' }]}
          >
            <Input placeholder="例如：tech-news" />
          </Form.Item>

          <Form.Item
            name="description"
            label="描述"
          >
            <TextArea rows={3} placeholder="分类描述" />
          </Form.Item>

          <Form.Item
            name="sort_order"
            label="排序"
            initialValue={0}
          >
            <Input type="number" placeholder="数字越小越靠前" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingId ? '更新' : '创建'}
              </Button>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Categories
