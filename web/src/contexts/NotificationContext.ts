import { createContext } from 'react';

export interface NotificationOptions {
  title: string;
  message: string;
  type?: 'success' | 'info' | 'warning' | 'error';
  duration?: number;
  placement?: 'topLeft' | 'topRight' | 'bottomLeft' | 'bottomRight';
  showProgress?: boolean;
  onClick?: () => void;
  actions?: Array<{
    label: string;
    onClick: () => void;
    type?: 'primary' | 'default' | 'link';
  }>;
}

export interface NotificationContextType {
  showNotification: (options: NotificationOptions) => void;
  showSuccess: (title: string, message: string, duration?: number) => void;
  showError: (title: string, message: string, duration?: number) => void;
  showInfo: (title: string, message: string, duration?: number) => void;
  showWarning: (title: string, message: string, duration?: number) => void;
  clearAll: () => void;
}

export const NotificationContext = createContext<NotificationContextType | null>(null);