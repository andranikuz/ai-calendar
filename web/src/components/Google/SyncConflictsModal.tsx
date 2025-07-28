import React, { useEffect, useState } from 'react';
import {
  Modal,
  List,
  Card,
  Button,
  Space,
  Typography,
  Tag,
  Alert,
  Divider,
  Checkbox,
  Tooltip,
  Badge,
  Row,
  Col,
  Spin,
  Empty,
  message
} from '../../utils/antd';
import {
  ExclamationCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  InfoCircleOutlined,
  CalendarOutlined,
  ClockCircleOutlined,
  MergeOutlined,
  GoogleOutlined,
  DesktopOutlined
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import {
  fetchPendingConflicts,
  resolveConflict,
  bulkResolveConflicts,
  clearError,
  SyncConflict,
  ConflictResolutionAction
} from '../../store/slices/syncConflictsSlice';
import dayjs from 'dayjs';

const { Text, Title } = Typography;

interface SyncConflictsModalProps {
  visible: boolean;
  onCancel: () => void;
}

const SyncConflictsModal: React.FC<SyncConflictsModalProps> = ({
  visible,
  onCancel
}) => {
  const dispatch = useAppDispatch();
  const { conflicts, isLoading, error, resolving } = useAppSelector(state => state.syncConflicts);
  
  const [selectedConflicts, setSelectedConflicts] = useState<string[]>([]);
  const [showResolved, setShowResolved] = useState(false);

  useEffect(() => {
    if (visible) {
      dispatch(fetchPendingConflicts());
    }
  }, [visible, dispatch]);

  useEffect(() => {
    if (error) {
      message.error(error);
      dispatch(clearError());
    }
  }, [error, dispatch]);

  const handleResolveConflict = async (conflictId: string, action: ConflictResolutionAction) => {
    try {
      await dispatch(resolveConflict({ conflictId, action })).unwrap();
      message.success('Конфликт успешно разрешен');
    } catch (error) {
      message.error('Не удалось разрешить конфликт');
    }
  };

  const handleBulkResolve = async (action: 'use_local' | 'use_google' | 'ignore') => {
    if (selectedConflicts.length === 0) {
      message.warning('Выберите конфликты для разрешения');
      return;
    }

    const actionText = {
      use_local: 'локальную версию',
      use_google: 'версию Google Calendar',
      ignore: 'игнорировать конфликты'
    };

    try {
      await dispatch(bulkResolveConflicts({
        conflictIds: selectedConflicts,
        action,
        resolution: `Массовое разрешение: выбрана ${actionText[action]}`
      })).unwrap();
      
      message.success(`Разрешено ${selectedConflicts.length} конфликтов`);
      setSelectedConflicts([]);
    } catch (error) {
      message.error('Не удалось разрешить конфликты');
    }
  };

  const getConflictIcon = (type: SyncConflict['conflict_type']) => {
    switch (type) {
      case 'time_overlap':
        return <ClockCircleOutlined style={{ color: '#ff4d4f' }} />;
      case 'content_diff':
        return <InfoCircleOutlined style={{ color: '#fa8c16' }} />;
      case 'duplicate_event':
        return <ExclamationCircleOutlined style={{ color: '#faad14' }} />;
      case 'deleted_event':
        return <CloseCircleOutlined style={{ color: '#f5222d' }} />;
      default:
        return <ExclamationCircleOutlined />;
    }
  };

  const getConflictTypeText = (type: SyncConflict['conflict_type']) => {
    switch (type) {
      case 'time_overlap':
        return 'Пересечение времени';
      case 'content_diff':
        return 'Различия в содержании';
      case 'duplicate_event':
        return 'Дублирование события';
      case 'deleted_event':
        return 'Удаленное событие';
      default:
        return 'Неизвестный тип';
    }
  };

  const getConflictTypeColor = (type: SyncConflict['conflict_type']) => {
    switch (type) {
      case 'time_overlap':
        return 'error';
      case 'content_diff':
        return 'warning';
      case 'duplicate_event':
        return 'gold';
      case 'deleted_event':
        return 'volcano';
      default:
        return 'default';
    }
  };

  const renderEventComparison = (localEvent?: SyncConflict['local_event'], googleEvent?: SyncConflict['google_event']) => {
    if (!localEvent && !googleEvent) return null;

    return (
      <Row gutter={16}>
        {localEvent && (
          <Col span={12}>
            <Card 
              size="small" 
              title={
                <Space>
                  <DesktopOutlined />
                  <Text strong>Локальная версия</Text>
                </Space>
              }
              style={{ marginBottom: 8 }}
            >
              <Space direction="vertical" size="small" style={{ width: '100%' }}>
                <div><Text strong>Название:</Text> {localEvent.title}</div>
                <div><Text strong>Время:</Text> {dayjs(localEvent.start_time).format('DD.MM.YYYY HH:mm')} - {dayjs(localEvent.end_time).format('HH:mm')}</div>
                {localEvent.description && <div><Text strong>Описание:</Text> {localEvent.description}</div>}
                {localEvent.location && <div><Text strong>Место:</Text> {localEvent.location}</div>}
              </Space>
            </Card>
          </Col>
        )}
        
        {googleEvent && (
          <Col span={12}>
            <Card 
              size="small" 
              title={
                <Space>
                  <GoogleOutlined />
                  <Text strong>Google Calendar</Text>
                </Space>
              }
              style={{ marginBottom: 8 }}
            >
              <Space direction="vertical" size="small" style={{ width: '100%' }}>
                <div><Text strong>Название:</Text> {googleEvent.title}</div>
                <div><Text strong>Время:</Text> {dayjs(googleEvent.start_time).format('DD.MM.YYYY HH:mm')} - {dayjs(googleEvent.end_time).format('HH:mm')}</div>
                {googleEvent.description && <div><Text strong>Описание:</Text> {googleEvent.description}</div>}
                {googleEvent.location && <div><Text strong>Место:</Text> {googleEvent.location}</div>}
              </Space>
            </Card>
          </Col>
        )}
      </Row>
    );
  };

  const renderConflictActions = (conflict: SyncConflict) => {
    const isResolving = resolving[conflict.id];
    
    return (
      <Space>
        {conflict.local_event && (
          <Button
            size="small"
            loading={isResolving}
            onClick={() => handleResolveConflict(conflict.id, {
              action: 'use_local',
              resolution: 'Выбрана локальная версия'
            })}
          >
            <DesktopOutlined /> Локальная
          </Button>
        )}
        
        {conflict.google_event && (
          <Button
            size="small"
            type="primary"
            loading={isResolving}
            onClick={() => handleResolveConflict(conflict.id, {
              action: 'use_google',
              resolution: 'Выбрана версия Google Calendar'
            })}
          >
            <GoogleOutlined /> Google
          </Button>
        )}
        
        <Button
          size="small"
          loading={isResolving}
          onClick={() => handleResolveConflict(conflict.id, {
            action: 'ignore',
            resolution: 'Конфликт проигнорирован'
          })}
        >
          Игнорировать
        </Button>
      </Space>
    );
  };

  const filteredConflicts = conflicts.filter(conflict => 
    showResolved || conflict.status === 'pending'
  );

  return (
    <Modal
      title={
        <Space>
          <ExclamationCircleOutlined style={{ color: '#fa8c16' }} />
          Конфликты синхронизации
          {conflicts.length > 0 && (
            <Badge count={conflicts.length} style={{ backgroundColor: '#fa8c16' }} />
          )}
        </Space>
      }
      open={visible}
      onCancel={onCancel}
      width={900}
      footer={[
        <Button key="close" onClick={onCancel}>
          Закрыть
        </Button>
      ]}
    >
      {isLoading ? (
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <Spin size="large" />
          <div style={{ marginTop: 16 }}>Загрузка конфликтов...</div>
        </div>
      ) : filteredConflicts.length === 0 ? (
        <Empty
          description="Конфликтов синхронизации не найдено"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      ) : (
        <>
          {/* Bulk actions */}
          {selectedConflicts.length > 0 && (
            <Alert
              message={
                <Space>
                  <Text>Выбрано конфликтов: {selectedConflicts.length}</Text>
                  <Button
                    size="small"
                    onClick={() => handleBulkResolve('use_local')}
                  >
                    <DesktopOutlined /> Локальные версии
                  </Button>
                  <Button
                    size="small"
                    type="primary"
                    onClick={() => handleBulkResolve('use_google')}
                  >
                    <GoogleOutlined /> Google версии
                  </Button>
                  <Button
                    size="small"
                    onClick={() => handleBulkResolve('ignore')}
                  >
                    Игнорировать все
                  </Button>
                  <Button
                    size="small"
                    onClick={() => setSelectedConflicts([])}
                  >
                    Отменить выбор
                  </Button>
                </Space>
              }
              type="info"
              style={{ marginBottom: 16 }}
            />
          )}

          {/* Conflicts list */}
          <List
            dataSource={filteredConflicts}
            renderItem={(conflict) => (
              <List.Item style={{ padding: '16px 0' }}>
                <Card
                  size="small"
                  style={{ width: '100%' }}
                  title={
                    <Space>
                      <Checkbox
                        checked={selectedConflicts.includes(conflict.id)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedConflicts([...selectedConflicts, conflict.id]);
                          } else {
                            setSelectedConflicts(selectedConflicts.filter(id => id !== conflict.id));
                          }
                        }}
                      />
                      {getConflictIcon(conflict.conflict_type)}
                      <Text strong>{getConflictTypeText(conflict.conflict_type)}</Text>
                      <Tag color={getConflictTypeColor(conflict.conflict_type)}>
                        {conflict.status}
                      </Tag>
                    </Space>
                  }
                  extra={renderConflictActions(conflict)}
                >
                  <Space direction="vertical" style={{ width: '100%' }}>
                    <Text type="secondary">{conflict.description}</Text>
                    
                    <Divider style={{ margin: '8px 0' }} />
                    
                    {renderEventComparison(conflict.local_event, conflict.google_event)}
                    
                    <div style={{ marginTop: 8 }}>
                      <Text type="secondary" style={{ fontSize: '12px' }}>
                        Создан: {dayjs(conflict.created_at).format('DD.MM.YYYY HH:mm')}
                      </Text>
                    </div>
                  </Space>
                </Card>
              </List.Item>
            )}
          />

          {conflicts.length > 5 && (
            <div style={{ textAlign: 'center', marginTop: 16 }}>
              <Text type="secondary">
                Показано {Math.min(filteredConflicts.length, 10)} из {conflicts.length} конфликтов
              </Text>
            </div>
          )}
        </>
      )}
    </Modal>
  );
};

export default SyncConflictsModal;