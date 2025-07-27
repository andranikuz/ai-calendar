import React from 'react';

interface SkipLinksProps {
  links?: Array<{
    href: string;
    label: string;
  }>;
}

/**
 * SkipLinks component for accessibility
 * Provides keyboard users quick navigation to main content areas
 */
const SkipLinks: React.FC<SkipLinksProps> = ({ 
  links = [
    { href: '#main-content', label: 'Skip to main content' },
    { href: '#main-navigation', label: 'Skip to navigation' },
    { href: '#user-actions', label: 'Skip to user actions' }
  ] 
}) => {
  const skipLinksStyle: React.CSSProperties = {
    position: 'fixed',
    top: 0,
    left: 0,
    zIndex: 9999,
  };

  const skipLinkStyle: React.CSSProperties = {
    position: 'absolute',
    top: 8,
    left: 8,
    padding: '8px 16px',
    background: '#1677ff',
    color: 'white',
    textDecoration: 'none',
    borderRadius: 4,
    fontWeight: 500,
    transform: 'translateY(-100%)',
    transition: 'transform 0.2s ease-in-out',
    zIndex: 10000,
  };

  const focusedStyle: React.CSSProperties = {
    ...skipLinkStyle,
    transform: 'translateY(0)',
    outline: '2px solid #ffffff',
    outlineOffset: 2,
  };

  return (
    <div style={skipLinksStyle}>
      {links.map((link, index) => (
        <a
          key={index}
          href={link.href}
          style={skipLinkStyle}
          onFocus={(e) => Object.assign(e.currentTarget.style, focusedStyle)}
          onBlur={(e) => Object.assign(e.currentTarget.style, skipLinkStyle)}
          onMouseEnter={(e) => e.currentTarget.style.background = '#0958d9'}
          onMouseLeave={(e) => e.currentTarget.style.background = '#1677ff'}
        >
          {link.label}
        </a>
      ))}
    </div>
  );
};

export default SkipLinks;