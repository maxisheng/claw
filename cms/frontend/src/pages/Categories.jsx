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
      setCategories(response.data)
    } catch (error) {
      console.error('Failed to load categories:', error)
      message.error('Failed to load categories')
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
        message.success('Category updated successfully')
      } else {
        await api.post('/categories', values)
        message.success('Category created successfully')
      }
      setModalVisible(false)
      loadCategories()
    } catch (error) {
      console.error('Failed to save category:', error)
      message.error('Failed to save category')
    }
  }

  const handleDelete = async (id) => {
    try {
      await api.delete(`/categories/${id}`)
      message.success('Category deleted successfully')
      loadCategories()
    } catch (error) {
      console.error('Failed to delete category:', error)
      message.error('Failed to delete category')
    }
  }

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Slug',
      dataIndex: 'slug',
      key: 'slug',
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
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
          <Button 
            type="link" 
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
          >
            Edit
          </Button>
          <Popconfirm
            title="Are you sure to delete this category?"
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
        <Title level={2} style={{ margin: 0 }}>Categories</Title>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
        >
          New Category
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
        title={editingId ? 'Edit Category' : 'New Category'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{ status: 'draft' }}
        >
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: 'Please input the name!' }]}
          >
            <Input placeholder="Category name" />
          </Form.Item>

          <Form.Item
            name="slug"
            label="Slug"
            rules={[{ required: true, message: 'Please input the slug!' }]}
          >
            <Input placeholder="e.g., tech-news" />
          </Form.Item>

          <Form.Item
            name="description"
            label="Description"
          >
            <TextArea rows={3} placeholder="Category description" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingId ? 'Update' : 'Create'}
              </Button>
              <Button onClick={() => setModalVisible(false)}>
                Cancel
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Categories
