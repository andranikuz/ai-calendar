import React, { useCallback, ReactNode } from 'react';
import { notification } from '../../utils/antd';
import {
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  InfoCircleOutlined,
  CloseCircleOutlined,
  BellOutlined
} from '@ant-design/icons';
import { NotificationContext, NotificationContextType, NotificationOptions } from '../../contexts/NotificationContext';

interface NotificationProviderProps {
  children: ReactNode;
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({ children }) => {
  const [api, contextHolder] = notification.useNotification();

  const getIcon = (type: string) => {
    switch (type) {
      case 'success':
        return <CheckCircleOutlined style={{ color: '#52c41a' }} />;
      case 'error':
        return <CloseCircleOutlined style={{ color: '#ff4d4f' }} />;
      case 'warning':
        return <ExclamationCircleOutlined style={{ color: '#faad14' }} />;
      case 'info':
        return <InfoCircleOutlined style={{ color: '#1677ff' }} />;
      default:
        return <BellOutlined />;
    }
  };

  const showNotification = useCallback((options: NotificationOptions) => {
    const {
      title,
      message,
      type = 'info',
      duration = 4.5,
      placement = 'topRight',
      actions = []
    } = options;

    api.open({
      message: title,
      description: message,
      icon: getIcon(type),
      duration,
      placement,
      btn: actions.length > 0 ? (
        <div style={{ display: 'flex', gap: '8px' }}>
          {actions.map((action, index) => (
            <button
              key={index}
              onClick={action.onClick}
              style={{
                padding: '4px 12px',
                border: action.type === 'primary' ? 'none' : '1px solid #d9d9d9',
                backgroundColor: action.type === 'primary' ? '#1677ff' : 'transparent',
                color: action.type === 'primary' ? 'white' : action.type === 'link' ? '#1677ff' : '#000',
                cursor: 'pointer',
                borderRadius: '6px',
                fontSize: '14px'
              }}
            >
              {action.label}
            </button>
          ))}
        </div>
      ) : undefined,
    });
  }, [api]);

  const showSuccess = useCallback((title: string, message: string, duration = 4.5) => {
    showNotification({ title, message, type: 'success', duration });
  }, [showNotification]);

  const showError = useCallback((title: string, message: string, duration = 4.5) => {
    showNotification({ title, message, type: 'error', duration });
  }, [showNotification]);

  const showInfo = useCallback((title: string, message: string, duration = 4.5) => {
    showNotification({ title, message, type: 'info', duration });
  }, [showNotification]);

  const showWarning = useCallback((title: string, message: string, duration = 4.5) => {
    showNotification({ title, message, type: 'warning', duration });
  }, [showNotification]);

  const clearAll = useCallback(() => {
    api.destroy();
  }, [api]);

  const contextValue: NotificationContextType = {
    showNotification,
    showSuccess,
    showError,
    showInfo,
    showWarning,
    clearAll,
  };

  return (
    <NotificationContext.Provider value={contextValue}>
      {contextHolder}
      {children}
    </NotificationContext.Provider>
  );
};