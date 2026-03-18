import { Outlet, Navigate, Link, useLocation, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { Layout, Menu, Button, Dropdown, Avatar, theme, Modal, Form, Input, message } from 'antd'
import {
  DashboardOutlined,
  FileTextOutlined,
  FolderOutlined,
  UserOutlined,
  LockOutlined,
  LogoutOutlined,
  SettingOutlined
} from '@ant-design/icons'
import api from '../api/axios'

const { Header, Sider, Content } = Layout
const { TextArea } = Input

function App() {
  const [user, setUser] = useState(null)
  const [collapsed, setCollapsed] = useState(false)
  const [passwordModalVisible, setPasswordModalVisible] = useState(false)
  const [passwordForm] = Form.useForm()
  const [passwordLoading, setPasswordLoading] = useState(false)
  const location = useLocation()
  const navigate = useNavigate()
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken()

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')
    if (token && userData) {
      setUser(JSON.parse(userData))
    }
  }, [])

  if (!user && location.pathname !== '/login') {
    return <Navigate to="/login" replace />
  }

  const handleLogout = async () => {
    try {
      await api.post('/admin/logout')
    } catch (e) {}
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
    navigate('/login')
    message.success('已退出登录')
  }

  const handleChangePassword = async (values) => {
    setPasswordLoading(true)
    try {
      await api.put('/admin/change-password', {
        old_password: values.old_password,
        new_password: values.new_password
      })
      message.success('密码修改成功！')
      setPasswordModalVisible(false)
      passwordForm.resetFields()
    } catch (err) {
      message.error(err.message || '修改密码失败')
    } finally {
      setPasswordLoading(false)
    }
  }

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">控制台</Link>,
    },
    {
      key: '/articles',
      icon: <FileTextOutlined />,
      label: <Link to="/articles">文章管理</Link>,
    },
    {
      key: '/categories',
      icon: <FolderOutlined />,
      label: <Link to="/categories">分类管理</Link>,
    },
  ]

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
      onClick: () => message.info('功能开发中...'),
    },
    {
      key: 'password',
      icon: <LockOutlined />,
      label: '修改密码',
      onClick: () => setPasswordModalVisible(true),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
        <div style={{ 
          height: 32, 
          margin: 16, 
          background: 'rgba(255, 255, 255, 0.2)',
          borderRadius: 6,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: 'white',
          fontWeight: 'bold'
        }}>
          {collapsed ? 'CMS' : 'CMS 管理系统'}
        </div>
        <Menu 
          theme="dark" 
          mode="inline" 
          selectedKeys={[location.pathname]}
          items={menuItems}
        />
      </Sider>
      <Layout>
        <Header style={{ 
          padding: '0 16px', 
          background: colorBgContainer,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between'
        }}>
          <span style={{ fontSize: 16 }}>
            {user && (
              <span>
                <Avatar style={{ backgroundColor: '#1890ff', marginRight: 8 }} icon={<UserOutlined />} />
                {user.username} ({user.role})
              </span>
            )}
          </span>
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <Button type="text" icon={<SettingOutlined />}>设置</Button>
          </Dropdown>
        </Header>
        <Content
          style={{
            margin: '16px',
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
          }}
        >
          <Outlet />
        </Content>
      </Layout>

      {/* 修改密码弹窗 */}
      <Modal
        title="修改密码"
        open={passwordModalVisible}
        onCancel={() => {
          setPasswordModalVisible(false)
          passwordForm.resetFields()
        }}
        footer={null}
        destroyOnClose
      >
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleChangePassword}
        >
          <Form.Item
            name="old_password"
            label="当前密码"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password placeholder="请输入当前密码" />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码长度至少 6 位' }
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认新密码"
            dependencies={['new_password']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'))
                },
              }),
            ]}
          >
            <Input.Password placeholder="请再次输入新密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={passwordLoading} block>
              确认修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  )
}

export default App
