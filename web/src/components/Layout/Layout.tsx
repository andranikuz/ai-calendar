import React, { useState } from 'react';
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
} from 'antd';
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
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { logout } from '../../store/slices/authSlice';
import { getIntegration } from '../../store/slices/googleSlice';
import OfflineIndicator from '../Common/OfflineIndicator';

const { Header, Sider, Content } = AntLayout;
const { Text } = Typography;

const Layout: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  
  const { user } = useAppSelector(state => state.auth);
  const { isConnected } = useAppSelector(state => state.google);

  React.useEffect(() => {
    // Check Google integration status on load
    dispatch(getIntegration());
  }, [dispatch]);

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
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
          onClick={({ key }) => navigate(key)}
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
            <OfflineIndicator showDetails />
            
            {isConnected && (
              <Tooltip title="Sync with Google Calendar">
                <Button
                  type="text"
                  icon={<SyncOutlined />}
                  onClick={() => {
                    // TODO: Trigger manual sync
                  }}
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
          <Outlet />
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;