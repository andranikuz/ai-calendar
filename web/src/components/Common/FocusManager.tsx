import React, { useEffect, useRef } from 'react';
import { announceToScreenReader } from '../../utils/accessibility';

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


export default FocusManager;