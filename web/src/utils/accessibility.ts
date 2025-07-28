/**
 * Utility function to announce messages to screen readers
 */
export const announceToScreenReader = (message: string, priority: 'polite' | 'assertive' = 'polite') => {
  const announcement = document.createElement('div');
  announcement.setAttribute('aria-live', priority);
  announcement.setAttribute('aria-atomic', 'true');
  announcement.className = 'sr-only';
  
  // Hide visually but keep accessible to screen readers
  Object.assign(announcement.style, {
    position: 'absolute',
    left: '-10000px',
    width: '1px',
    height: '1px',
    overflow: 'hidden',
    clip: 'rect(1px, 1px, 1px, 1px)',
    clipPath: 'inset(50%)',
    border: '0',
    padding: '0',
    margin: '0'
  });
  
  document.body.appendChild(announcement);
  announcement.textContent = message;
  
  // Remove the announcement after a short delay
  setTimeout(() => {
    if (announcement.parentNode) {
      announcement.parentNode.removeChild(announcement);
    }
  }, 1000);
};

/**
 * Custom hook for screen reader announcements
 */
export const useScreenReader = () => {
  const announce = (message: string, priority: 'polite' | 'assertive' = 'polite') => {
    announceToScreenReader(message, priority);
  };

  const announcePageChange = (pageName: string) => {
    announce(`Navigated to ${pageName} page`, 'polite');
  };

  const announceActionComplete = (action: string) => {
    announce(`${action} completed successfully`, 'polite');
  };

  const announceError = (error: string) => {
    announce(`Error: ${error}`, 'assertive');
  };

  const announceSuccess = (success: string) => {
    announce(`Success: ${success}`, 'polite');
  };

  return {
    announce,
    announcePageChange,
    announceActionComplete,
    announceError,
    announceSuccess
  };
};