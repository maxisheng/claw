import { Outlet, Navigate, Link, useLocation, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { Layout, Menu, Button, theme } from 'antd'
import {
  DashboardOutlined,
  FileTextOutlined,
  FolderOutlined,
  LogoutOutlined,
  UserOutlined
} from '@ant-design/icons'

const { Header, Sider, Content } = Layout

function App() {
  const [user, setUser] = useState(null)
  const [collapsed, setCollapsed] = useState(false)
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

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
    navigate('/login')
  }

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">Dashboard</Link>,
    },
    {
      key: '/articles',
      icon: <FileTextOutlined />,
      label: <Link to="/articles">Articles</Link>,
    },
    {
      key: '/categories',
      icon: <FolderOutlined />,
      label: <Link to="/categories">Categories</Link>,
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
          {collapsed ? 'CMS' : 'CMS Admin'}
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
            {user && <><UserOutlined /> {user.username}</>}
          </span>
          <Button type="primary" danger icon={<LogoutOutlined />} onClick={handleLogout}>
            Logout
          </Button>
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
    </Layout>
  )
}

export default App
