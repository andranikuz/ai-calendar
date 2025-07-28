import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  Layout as AntLayout,
  Menu,
  Button,
  Avatar,
  Dropdown,
  Typography,
  Space,
  Badge,
  Tooltip,
  Drawer,
  Grid,
} from '../../utils/antd';
import {
  DashboardOutlined,
  CalendarOutlined,
  FlagOutlined,
  SmileOutlined,
  SettingOutlined,
  LogoutOutlined,
  UserOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SyncOutlined,
  MenuOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { logout } from '../../store/slices/authSlice';
import { getIntegration } from '../../store/slices/googleSlice';
import { NotificationProvider } from '../Common/NotificationProvider';
import { useNotifications, useSystemNotifications } from '../../hooks/useNotifications';

const { Header, Sider, Content } = AntLayout;
const { Text } = Typography;
const { useBreakpoint } = Grid;

interface LayoutContentProps {
  children?: React.ReactNode;
}

const LayoutContent: React.FC<LayoutContentProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);
  const [mobileMenuVisible, setMobileMenuVisible] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  const screens = useBreakpoint();
  
  const { user } = useAppSelector(state => state.auth);
  const { isConnected } = useAppSelector(state => state.google);
  
  const { showSuccess, showInfo } = useNotifications();
  const { showWelcomeMessage } = useSystemNotifications();

  // Determine if we're on mobile
  const isMobile = !screens.md;

  useEffect(() => {
    // Check Google integration status on load
    dispatch(getIntegration());
  }, [dispatch]);

  useEffect(() => {
    // Show welcome message on first load
    if (user && location.pathname === '/dashboard') {
      const hasShownWelcome = sessionStorage.getItem('hasShownWelcome');
      if (!hasShownWelcome) {
        setTimeout(() => {
          showWelcomeMessage(user.name);
          sessionStorage.setItem('hasShownWelcome', 'true');
        }, 1000);
      }
    }
  }, [user, location.pathname, showWelcomeMessage]);

  // Auto-collapse sidebar on mobile
  useEffect(() => {
    if (isMobile) {
      setCollapsed(true);
    }
  }, [isMobile]);

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
    showInfo('Goodbye!', 'You have been logged out successfully.');
  };

  const handleSync = async () => {
    try {
      // TODO: Implement actual sync logic
      showSuccess('Sync Complete', 'Your calendar has been synchronized with Google Calendar.');
    } catch (error) {
      console.error('Sync failed:', error);
    }
  };

  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: 'Dashboard',
    },
    {
      key: '/calendar',
      icon: <CalendarOutlined />,
      label: 'Calendar',
    },
    {
      key: '/goals',
      icon: <FlagOutlined />,
      label: 'Goals',
    },
    {
      key: '/moods',
      icon: <SmileOutlined />,
      label: 'Moods',
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: 'Settings',
    },
  ];

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: 'Profile',
      onClick: () => navigate('/settings'),
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'Logout',
      onClick: handleLogout,
    },
  ];

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
    if (isMobile) {
      setMobileMenuVisible(false);
    }
  };


  // Mobile Layout
  if (isMobile) {
    return (
      <AntLayout style={{ minHeight: '100vh' }}>
        {/* Mobile Header */}
        <Header style={{ 
          padding: '0 16px', 
          background: '#fff', 
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          position: 'sticky',
          top: 0,
          zIndex: 100,
        }}>
          <Space>
            <Button
              type="text"
              icon={<MenuOutlined />}
              onClick={() => setMobileMenuVisible(true)}
            />
            <Typography.Title level={5} style={{ margin: 0, color: '#1677ff' }}>
              Smart Goal Calendar
            </Typography.Title>
          </Space>
          
          <Space>
            {isConnected && (
              <Tooltip title="Sync with Google Calendar">
                <Button
                  type="text"
                  icon={<SyncOutlined />}
                  onClick={handleSync}
                />
              </Tooltip>
            )}
            
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Avatar icon={<UserOutlined />} style={{ cursor: 'pointer' }} />
            </Dropdown>
          </Space>
        </Header>

        {/* Mobile Drawer Menu */}
        <Drawer
          title="Navigation"
          placement="left"
          onClose={() => setMobileMenuVisible(false)}
          open={mobileMenuVisible}
          width={280}
          styles={{
            body: { padding: 0 }
          }}
        >
          <div style={{ 
            padding: '16px', 
            borderBottom: '1px solid #f0f0f0',
            textAlign: 'center'
          }}>
            <Space direction="vertical" align="center">
              <Avatar size={64} icon={<UserOutlined />} />
              <div>
                <Text strong>{user?.name}</Text>
                <br />
                <Text type="secondary">{user?.email}</Text>
              </div>
            </Space>
          </div>
          
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={handleMenuClick}
            style={{ border: 'none' }}
          />
          
          {/* Google Integration Status */}
          <div style={{ 
            position: 'absolute', 
            bottom: 16, 
            left: 16, 
            right: 16,
            borderTop: '1px solid #f0f0f0',
            paddingTop: 16
          }}>
            <Badge 
              status={isConnected ? 'success' : 'default'} 
              text={isConnected ? 'Google Calendar Connected' : 'Connect Google Calendar'}
            />
          </div>
        </Drawer>
        
        <Content style={{ 
          padding: '16px',
          background: '#f0f2f5',
          minHeight: 'calc(100vh - 64px - 64px)', // Header height and bottom nav
          paddingBottom: '80px' // Space for mobile navigation
        }}>
          <div style={{
            background: '#fff',
            borderRadius: '8px',
            padding: '16px',
            minHeight: '100%'
          }}>
            {children || <Outlet />}
          </div>
        </Content>

        {/* Mobile Bottom Navigation will be added here */}
        <div style={{
          position: 'fixed',
          bottom: 0,
          left: 0,
          right: 0,
          height: '64px',
          background: '#fff',
          borderTop: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-around',
          zIndex: 1000,
          boxShadow: '0 -2px 8px rgba(0, 0, 0, 0.1)'
        }}>
          {menuItems.map((item) => (
            <Button
              key={item.key}
              type="text"
              icon={item.icon}
              onClick={() => navigate(item.key)}
              style={{
                display: 'flex',
                flexDirection: 'column',
                height: '48px',
                width: '64px',
                fontSize: '12px',
                color: location.pathname === item.key ? '#1677ff' : '#666',
                background: location.pathname === item.key ? '#f6ffed' : 'transparent'
              }}
            >
              <div style={{ fontSize: '10px', marginTop: '2px' }}>
                {item.label}
              </div>
            </Button>
          ))}
        </div>
      </AntLayout>
    );
  }

  // Desktop Layout
  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider 
        trigger={null} 
        collapsible 
        collapsed={collapsed}
        width={240}
        style={{
          background: '#fff',
          borderRight: '1px solid #f0f0f0',
        }}
        breakpoint="lg"
        collapsedWidth={isMobile ? 0 : 80}
      >
        <div style={{ 
          padding: '16px', 
          borderBottom: '1px solid #f0f0f0',
          textAlign: 'center'
        }}>
          <Typography.Title level={4} style={{ margin: 0, color: '#1677ff' }}>
            {collapsed ? 'SGC' : 'Smart Goal Calendar'}
          </Typography.Title>
        </div>
        
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{ border: 'none' }}
        />
        
        {/* Google Integration Status */}
        <div style={{ 
          position: 'absolute', 
          bottom: 16, 
          left: 16, 
          right: 16,
          borderTop: '1px solid #f0f0f0',
          paddingTop: 16
        }}>
          <Tooltip title={isConnected ? 'Google Calendar Connected' : 'Connect Google Calendar'}>
            <Badge 
              status={isConnected ? 'success' : 'default'} 
              text={collapsed ? '' : (isConnected ? 'Google Connected' : 'Connect Google')}
            />
          </Tooltip>
        </div>
      </Sider>
      
      <AntLayout>
        <Header style={{ 
          padding: '0 24px', 
          background: '#fff', 
          borderBottom: '1px solid #f0f0f0',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between'
        }}>
          <Space>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
            />
          </Space>
          
          <Space>
            {isConnected && (
              <Tooltip title="Sync with Google Calendar">
                <Button
                  type="text"
                  icon={<SyncOutlined />}
                  onClick={handleSync}
                />
              </Tooltip>
            )}
            
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer' }}>
                <Avatar icon={<UserOutlined />} />
                <Text>{user?.name}</Text>
              </Space>
            </Dropdown>
          </Space>
        </Header>
        
        <Content style={{ 
          margin: '24px',
          padding: '24px',
          background: '#fff',
          borderRadius: '8px',
          minHeight: 'calc(100vh - 112px)'
        }}>
          {children || <Outlet />}
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

const ResponsiveLayout: React.FC = () => {
  return (
    <NotificationProvider>
      <LayoutContent />
    </NotificationProvider>
  );
};

export default ResponsiveLayout;