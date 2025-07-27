import React, { useEffect, useRef } from 'react';

interface FocusManagerProps {
  children: React.ReactNode;
  announceOnMount?: string;
  autoFocus?: boolean;
  restoreFocus?: boolean;
}

/**
 * FocusManager component for accessibility
 * Manages focus state and announces changes to screen readers
 */
const FocusManager: React.FC<FocusManagerProps> = ({
  children,
  announceOnMount,
  autoFocus = false,
  restoreFocus = false
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const previousActiveElement = useRef<HTMLElement | null>(null);

  useEffect(() => {
    // Store the currently focused element
    if (restoreFocus) {
      previousActiveElement.current = document.activeElement as HTMLElement;
    }

    // Auto focus the first focusable element in the container
    if (autoFocus && containerRef.current) {
      const focusableElement = containerRef.current.querySelector(
        'a[href], button, textarea, input[type="text"], input[type="radio"], input[type="checkbox"], select, [tabindex]:not([tabindex="-1"])'
      ) as HTMLElement;
      
      if (focusableElement) {
        // Small delay to ensure the element is rendered
        setTimeout(() => focusableElement.focus(), 100);
      }
    }

    // Announce to screen readers
    if (announceOnMount) {
      announceToScreenReader(announceOnMount);
    }

    // Cleanup function to restore focus
    return () => {
      if (restoreFocus && previousActiveElement.current) {
        previousActiveElement.current.focus();
      }
    };
  }, [autoFocus, restoreFocus, announceOnMount]);

  return (
    <div ref={containerRef}>
      {children}
    </div>
  );
};

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
    clip: 'rect(0, 0, 0, 0)',
    whiteSpace: 'nowrap',
    border: '0',
    margin: '0',
    padding: '0'
  });

  announcement.textContent = message;
  document.body.appendChild(announcement);

  // Remove after screen readers have announced it
  setTimeout(() => {
    if (document.body.contains(announcement)) {
      document.body.removeChild(announcement);
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

  return {
    announce,
    announcePageChange,
    announceActionComplete,
    announceError
  };
};

export default FocusManager;