import { useCallback, useContext } from 'react';
import { NotificationContext } from '../contexts/NotificationContext';

export const useNotifications = () => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider');
  }
  return context;
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
      onClick: () => {
        // Navigate to mood page
        window.location.hash = '#/moods';
      }
    });
  }, [showNotification]);

  const showGoalProgress = useCallback((goalTitle: string, progress: number) => {
    const emoji = progress === 100 ? '🎉' : progress >= 75 ? '🚀' : progress >= 50 ? '💪' : '📈';
    
    showNotification({
      title: `${emoji} Goal Progress Update`,
      message: `${goalTitle} is now ${progress}% complete!`,
      type: progress === 100 ? 'success' : 'info',
      duration: 8,
      placement: 'topRight',
      onClick: () => {
        // Navigate to goals page
        window.location.hash = '#/goals';
      }
    });
  }, [showNotification]);

  const showEventReminder = useCallback((eventTitle: string, timeUntil: string) => {
    showNotification({
      title: '⏰ Upcoming Event',
      message: `${eventTitle} starts in ${timeUntil}`,
      type: 'warning',
      duration: 15,
      placement: 'topRight',
      onClick: () => {
        // Navigate to calendar page
        window.location.hash = '#/calendar';
      }
    });
  }, [showNotification]);

  const showSyncComplete = useCallback((source: string, eventsCount: number) => {
    showNotification({
      title: '✅ Sync Complete',
      message: `Successfully synced ${eventsCount} events from ${source}`,
      type: 'success',
      duration: 5,
      placement: 'bottomRight'
    });
  }, [showNotification]);

  const showOfflineMode = useCallback(() => {
    showNotification({
      title: '📱 Offline Mode',
      message: 'You\'re working offline. Changes will sync when connection is restored.',
      type: 'info',
      duration: 8,
      placement: 'bottomLeft'
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
    showGoalProgress,
    showEventReminder,
    showSyncComplete,
    showOfflineMode,
    showWelcomeMessage
  };
};