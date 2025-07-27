import React from 'react';
import { Card, Space, Typography, Button, Tag, Grid } from '../../utils/antd';
import { MoreOutlined } from '@ant-design/icons';

const { Text } = Typography;
const { useBreakpoint } = Grid;

interface MobileCardProps {
  title: string;
  subtitle?: string;
  description?: string;
  status?: string;
  statusColor?: string;
  priority?: string;
  priorityColor?: string;
  progress?: number;
  icon?: React.ReactNode;
  avatar?: React.ReactNode;
  actions?: Array<{
    key: string;
    label: string;
    icon?: React.ReactNode;
    onClick: () => void;
    type?: 'primary' | 'default' | 'text' | 'link';
    danger?: boolean;
  }>;
  extra?: React.ReactNode;
  onClick?: () => void;
  style?: React.CSSProperties;
  bodyStyle?: React.CSSProperties;
  selected?: boolean;
}

const MobileCard: React.FC<MobileCardProps> = ({
  title,
  subtitle,
  description,
  status,
  statusColor,
  priority,
  priorityColor,
  progress,
  icon,
  avatar,
  actions = [],
  extra,
  onClick,
  style,
  bodyStyle,
  selected = false
}) => {
  const screens = useBreakpoint();
  const isMobile = !screens.md;

  const cardStyle: React.CSSProperties = {
    marginBottom: isMobile ? 12 : 16,
    borderRadius: isMobile ? 12 : 8,
    border: selected ? '2px solid #1677ff' : '1px solid #f0f0f0',
    backgroundColor: selected ? '#f6ffed' : '#fff',
    cursor: onClick ? 'pointer' : 'default',
    transition: 'all 0.3s ease',
    boxShadow: isMobile 
      ? '0 2px 8px rgba(0, 0, 0, 0.1)' 
      : '0 1px 2px rgba(0, 0, 0, 0.03)',
    ...style
  };

  const cardBodyStyle: React.CSSProperties = {
    padding: isMobile ? '16px' : '20px',
    ...bodyStyle
  };

  const handleCardClick = (e: React.MouseEvent) => {
    // Don't trigger card click if clicking on actions
    if ((e.target as Element).closest('.mobile-card-actions')) {
      return;
    }
    onClick?.();
  };

  const renderMobileActions = () => {
    if (actions.length === 0) return null;

    if (isMobile && actions.length > 2) {
      // Show first action and "more" button on mobile
      const primaryAction = actions[0];
      return (
        <Space className="mobile-card-actions">
          <Button
            type={primaryAction.type || 'default'}
            size="small"
            icon={primaryAction.icon}
            onClick={primaryAction.onClick}
            danger={primaryAction.danger}
          >
            {primaryAction.label}
          </Button>
          <Button
            type="text"
            size="small"
            icon={<MoreOutlined />}
            // TODO: Show action sheet with all actions
          />
        </Space>
      );
    }

    return (
      <Space className="mobile-card-actions" wrap>
        {actions.slice(0, isMobile ? 2 : actions.length).map((action) => (
          <Button
            key={action.key}
            type={action.type || 'default'}
            size={isMobile ? 'small' : 'middle'}
            icon={action.icon}
            onClick={action.onClick}
            danger={action.danger}
          >
            {isMobile && actions.length > 1 ? action.icon : action.label}
          </Button>
        ))}
      </Space>
    );
  };

  return (
    <Card
      style={cardStyle}
      bodyStyle={cardBodyStyle}
      onClick={handleCardClick}
      hoverable={!!onClick}
    >
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'flex-start',
        marginBottom: description || progress !== undefined ? 12 : 0
      }}>
        <div style={{ flex: 1, minWidth: 0 }}>
          {/* Header with icon/avatar and title */}
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            marginBottom: subtitle ? 4 : 8,
            gap: 12
          }}>
            {avatar || icon}
            <div style={{ flex: 1, minWidth: 0 }}>
              <Text 
                strong 
                style={{ 
                  fontSize: isMobile ? 14 : 16,
                  display: 'block',
                  wordBreak: 'break-word'
                }}
              >
                {title}
              </Text>
              {subtitle && (
                <Text 
                  type="secondary" 
                  style={{ 
                    fontSize: isMobile ? 12 : 14,
                    display: 'block'
                  }}
                >
                  {subtitle}
                </Text>
              )}
            </div>
          </div>

          {/* Tags for status and priority */}
          {(status || priority) && (
            <Space wrap style={{ marginBottom: 8 }}>
              {status && (
                <Tag color={statusColor}>
                  {status.toUpperCase()}
                </Tag>
              )}
              {priority && (
                <Tag color={priorityColor}>
                  {priority.toUpperCase()}
                </Tag>
              )}
            </Space>
          )}

          {/* Description */}
          {description && (
            <Text 
              type="secondary" 
              style={{ 
                fontSize: isMobile ? 12 : 14,
                display: 'block',
                marginBottom: 12,
                lineHeight: 1.4,
                wordBreak: 'break-word'
              }}
            >
              {description.length > (isMobile ? 80 : 120) 
                ? `${description.substring(0, isMobile ? 80 : 120)}...` 
                : description
              }
            </Text>
          )}

          {/* Progress bar */}
          {progress !== undefined && (
            <div style={{ marginBottom: 12 }}>
              <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: 4
              }}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  Progress
                </Text>
                <Text strong style={{ fontSize: 12 }}>
                  {progress}%
                </Text>
              </div>
              <div style={{
                height: 6,
                background: '#f0f0f0',
                borderRadius: 3,
                overflow: 'hidden'
              }}>
                <div style={{
                  height: '100%',
                  width: `${progress}%`,
                  background: progress === 100 ? '#52c41a' : '#1677ff',
                  borderRadius: 3,
                  transition: 'width 0.3s ease'
                }} />
              </div>
            </div>
          )}
        </div>

        {/* Extra content */}
        {extra && (
          <div style={{ marginLeft: 12, flexShrink: 0 }}>
            {extra}
          </div>
        )}
      </div>

      {/* Actions */}
      {actions.length > 0 && (
        <div style={{
          display: 'flex',
          justifyContent: 'flex-end',
          borderTop: isMobile ? '1px solid #f0f0f0' : 'none',
          paddingTop: isMobile ? 12 : 0,
          marginTop: isMobile ? 12 : 0
        }}>
          {renderMobileActions()}
        </div>
      )}
    </Card>
  );
};

export default MobileCard;