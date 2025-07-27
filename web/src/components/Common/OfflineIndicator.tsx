import React from 'react';
import { Badge, Tooltip, Button, Space } from '../../utils/antd';
import { WifiOutlined, SyncOutlined, CloudSyncOutlined } from '@ant-design/icons';
import { useOffline } from '../../hooks/useOffline';

interface OfflineIndicatorProps {
  showDetails?: boolean;
  style?: React.CSSProperties;
}

const OfflineIndicator: React.FC<OfflineIndicatorProps> = ({ 
  showDetails = false, 
  style 
}) => {
  const {
    isOnline,
    isOfflineMode,
    pendingActions,
    syncStatus,
    syncPendingActions
  } = useOffline();

  const getStatusColor = () => {
    if (syncStatus === 'syncing') return 'processing';
    if (syncStatus === 'error') return 'error';
    if (!isOnline) return 'default';
    if (pendingActions.length > 0) return 'warning';
    return 'success';
  };

  const getStatusText = () => {
    if (syncStatus === 'syncing') return 'Syncing...';
    if (syncStatus === 'error') return 'Sync Error';
    if (!isOnline) return 'Offline';
    if (pendingActions.length > 0) return `${pendingActions.length} pending`;
    return 'Online';
  };

  const getStatusIcon = () => {
    if (syncStatus === 'syncing') return <SyncOutlined spin />;
    if (!isOnline) return <WifiOutlined style={{ opacity: 0.5 }} />;
    if (pendingActions.length > 0) return <CloudSyncOutlined />;
    return <WifiOutlined />;
  };

  const handleSyncClick = async () => {
    if (isOnline && syncStatus !== 'syncing') {
      await syncPendingActions();
    }
  };

  if (!showDetails && isOnline && pendingActions.length === 0 && syncStatus === 'idle') {
    return null; // Hide when everything is normal
  }

  const indicator = (
    <Badge 
      color={getStatusColor()} 
      text={showDetails ? getStatusText() : undefined}
      style={style}
    >
      {getStatusIcon()}
    </Badge>
  );

  const tooltipTitle = (
    <div>
      <div><strong>Status:</strong> {getStatusText()}</div>
      {isOfflineMode && (
        <div style={{ marginTop: 4 }}>
          Your changes will sync when you're back online
        </div>
      )}
      {pendingActions.length > 0 && (
        <div style={{ marginTop: 4 }}>
          <strong>Pending actions:</strong> {pendingActions.length}
        </div>
      )}
      {isOnline && pendingActions.length > 0 && syncStatus === 'idle' && (
        <div style={{ marginTop: 4, color: '#1677ff' }}>
          Click to sync now
        </div>
      )}
    </div>
  );

  if (showDetails) {
    return (
      <Tooltip title={tooltipTitle}>
        <Space>
          {indicator}
          {isOnline && pendingActions.length > 0 && syncStatus === 'idle' && (
            <Button 
              size="small" 
              type="text" 
              icon={<SyncOutlined />}
              onClick={handleSyncClick}
            >
              Sync
            </Button>
          )}
        </Space>
      </Tooltip>
    );
  }

  return (
    <Tooltip title={tooltipTitle}>
      <div 
        role={isOnline && pendingActions.length > 0 && syncStatus === 'idle' ? 'button' : undefined}
        tabIndex={isOnline && pendingActions.length > 0 && syncStatus === 'idle' ? 0 : undefined}
        onClick={isOnline && pendingActions.length > 0 && syncStatus === 'idle' ? handleSyncClick : undefined}
        onKeyDown={isOnline && pendingActions.length > 0 && syncStatus === 'idle' ? (e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            handleSyncClick();
          }
        } : undefined}
        style={{ 
          cursor: isOnline && pendingActions.length > 0 && syncStatus === 'idle' ? 'pointer' : 'default'
        }}
        aria-label={isOnline && pendingActions.length > 0 && syncStatus === 'idle' ? 'Sync pending changes' : undefined}
      >
        {indicator}
      </div>
    </Tooltip>
  );
};

export default OfflineIndicator;