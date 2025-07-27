import React, { createContext, useContext, useCallback, ReactNode } from 'react';
import { notification } from 'antd';
import {
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  InfoCircleOutlined,
  CloseCircleOutlined,
  BellOutlined
} from '@ant-design/icons';

export interface NotificationOptions {
  title: string;
  message: string;
  type?: 'success' | 'info' | 'warning' | 'error';
  duration?: number;
  placement?: 'topLeft' | 'topRight' | 'bottomLeft' | 'bottomRight';
  showProgress?: boolean;
  actions?: Array<{
    label: string;
    onClick: () => void;
    type?: 'primary' | 'default' | 'link';
  }>;
}

interface NotificationContextType {
  showNotification: (options: NotificationOptions) => void;
  showSuccess: (title: string, message: string, duration?: number) => void;
  showError: (title: string, message: string, duration?: number) => void;
  showInfo: (title: string, message: string, duration?: number) => void;
  showWarning: (title: string, message: string, duration?: number) => void;
  clearAll: () => void;
}

const NotificationContext = createContext<NotificationContextType | null>(null);

export const useNotifications = () => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider');
  }
  return context;
};

interface NotificationProviderProps {
  children: ReactNode;
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({ children }) => {
  const [api, contextHolder] = notification.useNotification();

  const getIcon = (type: string) => {
    switch (type) {
      case 'success': return <CheckCircleOutlined style={{ color: '#52c41a' }} />;
      case 'error': return <CloseCircleOutlined style={{ color: '#ff4d4f' }} />;
      case 'warning': return <ExclamationCircleOutlined style={{ color: '#faad14' }} />;
      case 'info': return <InfoCircleOutlined style={{ color: '#1677ff' }} />;
      default: return <BellOutlined style={{ color: '#1677ff' }} />;
    }
  };

  const showNotification = useCallback((options: NotificationOptions) => {
    const {
      title,
      message,
      type = 'info',
      duration = 4.5,
      placement = 'topRight',
      showProgress = false,
      actions = []
    } = options;

    api.open({
      message: title,
      description: message,
      icon: getIcon(type),
      placement,
      duration,
      showProgress,
      btn: actions.length > 0 ? (
        <div style={{ display: 'flex', gap: 8 }}>
          {actions.map((action, index) => (
            <button
              key={index}
              onClick={action.onClick}
              style={{
                background: action.type === 'primary' ? '#1677ff' : 'transparent',
                color: action.type === 'primary' ? 'white' : '#1677ff',
                border: action.type === 'link' ? 'none' : '1px solid #1677ff',
                borderRadius: 4,
                padding: '4px 12px',
                cursor: 'pointer',
                fontSize: 12
              }}
            >
              {action.label}
            </button>
          ))}
        </div>
      ) : undefined,
      style: {
        borderRadius: 8,
        boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
      },
    });
  }, [api]);

  const showSuccess = useCallback((title: string, message: string, duration = 4.5) => {
    showNotification({ title, message, type: 'success', duration });
  }, [showNotification]);

  const showError = useCallback((title: string, message: string, duration = 6) => {
    showNotification({ title, message, type: 'error', duration });
  }, [showNotification]);

  const showInfo = useCallback((title: string, message: string, duration = 4.5) => {
    showNotification({ title, message, type: 'info', duration });
  }, [showNotification]);

  const showWarning = useCallback((title: string, message: string, duration = 5) => {
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

// Hook for system notifications (like reminders)
export const useSystemNotifications = () => {
  const { showNotification } = useNotifications();

  const showMoodReminder = useCallback(() => {
    showNotification({
      title: '🌟 Daily Mood Check',
      message: 'How are you feeling today? Take a moment to record your mood.',
      type: 'info',
      duration: 10,
      placement: 'topRight',
      actions: [
        {
          label: 'Record Mood',
          onClick: () => {
            // Navigate to mood page
            window.location.href = '/moods';
          },
          type: 'primary'
        },
        {
          label: 'Later',
          onClick: () => {},
          type: 'link'
        }
      ]
    });
  }, [showNotification]);

  const showGoalProgressReminder = useCallback((goalTitle: string, progress: number) => {
    showNotification({
      title: '🎯 Goal Progress Update',
      message: `"${goalTitle}" is ${progress}% complete. Keep up the great work!`,
      type: 'success',
      duration: 8,
      placement: 'topRight',
      actions: [
        {
          label: 'View Goals',
          onClick: () => {
            window.location.href = '/goals';
          },
          type: 'primary'
        }
      ]
    });
  }, [showNotification]);

  const showEventReminder = useCallback((eventTitle: string, timeUntilEvent: string) => {
    showNotification({
      title: '⏰ Event Reminder',
      message: `"${eventTitle}" starts ${timeUntilEvent}`,
      type: 'warning',
      duration: 15,
      placement: 'topRight',
      actions: [
        {
          label: 'View Calendar',
          onClick: () => {
            window.location.href = '/calendar';
          },
          type: 'primary'
        },
        {
          label: 'Snooze',
          onClick: () => {
            // Implement snooze logic
          },
          type: 'link'
        }
      ]
    });
  }, [showNotification]);

  const showSyncSuccess = useCallback((calendarName: string, eventsCount: number) => {
    showNotification({
      title: '✅ Sync Complete',
      message: `Successfully synced ${eventsCount} events with ${calendarName}`,
      type: 'success',
      duration: 5,
      placement: 'bottomRight'
    });
  }, [showNotification]);

  const showSyncError = useCallback((calendarName: string, error: string) => {
    showNotification({
      title: '❌ Sync Failed',
      message: `Failed to sync with ${calendarName}: ${error}`,
      type: 'error',
      duration: 8,
      placement: 'bottomRight',
      actions: [
        {
          label: 'Retry',
          onClick: () => {
            // Implement retry logic
          },
          type: 'primary'
        },
        {
          label: 'Settings',
          onClick: () => {
            window.location.href = '/settings';
          },
          type: 'link'
        }
      ]
    });
  }, [showNotification]);

  const showWelcomeMessage = useCallback((userName: string) => {
    showNotification({
      title: `👋 Welcome back, ${userName}!`,
      message: 'Ready to tackle your goals and manage your schedule?',
      type: 'info',
      duration: 6,
      placement: 'topRight'
    });
  }, [showNotification]);

  return {
    showMoodReminder,
    showGoalProgressReminder,
    showEventReminder,
    showSyncSuccess,
    showSyncError,
    showWelcomeMessage,
  };
};