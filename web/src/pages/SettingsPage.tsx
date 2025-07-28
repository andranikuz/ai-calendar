import React, { useEffect, useState } from 'react';
import {
  Row,
  Col,
  Card,
  Typography,
  Form,
  Input,
  Button,
  Switch,
  Space,
  Divider,
  Alert,
  List,
  Tag,
  Modal,
  Select,
  TimePicker,
  message,
  Avatar,
  Popconfirm
} from '../utils/antd';
import {
  GoogleOutlined,
  UserOutlined,
  BellOutlined,
  SecurityScanOutlined,
  LinkOutlined,
  DisconnectOutlined,
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  CalendarOutlined,
  SyncOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { updateUser } from '../store/slices/authSlice';
import { 
  fetchGoogleIntegration, 
  disconnectGoogle, 
  fetchGoogleCalendars,
  createCalendarSync,
  updateCalendarSync,
  deleteCalendarSync,
  manualSync
} from '../store/slices/googleSlice';
import { fetchPendingConflicts } from '../store/slices/syncConflictsSlice';
import { GoogleCalendar, GoogleCalendarSync } from '../types/api';
import SyncConflictsModal from '../components/Google/SyncConflictsModal';
import dayjs from 'dayjs';

const { Title, Text, Paragraph } = Typography;

interface CalendarSyncModalProps {
  visible: boolean;
  calendar: GoogleCalendar | null;
  sync: GoogleCalendarSync | null;
  onCancel: () => void;
  onSuccess: () => void;
}

const CalendarSyncModal: React.FC<CalendarSyncModalProps> = ({ 
  visible, 
  calendar, 
  sync, 
  onCancel, 
  onSuccess 
}) => {
  const [form] = Form.useForm();
  const dispatch = useAppDispatch();
  const { isLoading } = useAppSelector(state => state.google);

  const isEditing = !!sync;

  useEffect(() => {
    if (visible) {
      if (sync) {
        form.setFieldsValue({
          sync_direction: sync.sync_direction,
          auto_sync: sync.settings.auto_sync,
          sync_interval: sync.settings.sync_interval,
          sync_past_events: sync.settings.sync_past_events,
          sync_future_events: sync.settings.sync_future_events,
          conflict_resolution: sync.settings.conflict_resolution
        });
      } else {
        form.setFieldsValue({
          sync_direction: 'bidirectional',
          auto_sync: true,
          sync_interval: 15,
          sync_past_events: false,
          sync_future_events: true,
          conflict_resolution: 'google_wins'
        });
      }
    } else {
      form.resetFields();
    }
  }, [visible, sync, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      const syncData = {
        calendar_id: calendar?.id || sync?.calendar_id || '',
        sync_direction: values.sync_direction,
        settings: {
          auto_sync: values.auto_sync,
          sync_interval: values.sync_interval,
          sync_past_events: values.sync_past_events,
          sync_future_events: values.sync_future_events,
          conflict_resolution: values.conflict_resolution
        }
      };

      if (isEditing && sync) {
        await dispatch(updateCalendarSync({
          id: sync.id,
          data: syncData
        })).unwrap();
        message.success('Calendar sync updated successfully');
      } else {
        await dispatch(createCalendarSync(syncData)).unwrap();
        message.success('Calendar sync created successfully');
      }

      onSuccess();
    } catch (error) {
      console.error('Failed to save calendar sync:', error);
      message.error('Failed to save calendar sync');
    }
  };

  return (
    <Modal
      title={
        <Space>
          <CalendarOutlined />
          {isEditing ? 'Edit Calendar Sync' : 'Setup Calendar Sync'}
          {calendar && <Text type="secondary">- {calendar.summary}</Text>}
        </Space>
      }
      open={visible}
      onCancel={onCancel}
      footer={[
        <Button key="cancel" onClick={onCancel}>
          Cancel
        </Button>,
        <Button
          key="submit"
          type="primary"
          loading={isLoading}
          onClick={handleSubmit}
        >
          {isEditing ? 'Update Sync' : 'Setup Sync'}
        </Button>
      ]}
      width={600}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="sync_direction"
          label="Sync Direction"
          rules={[{ required: true, message: 'Please select sync direction' }]}
        >
          <Select>
            <Select.Option value="bidirectional">
              ↔️ Bidirectional (both ways)
            </Select.Option>
            <Select.Option value="from_google">
              ⬇️ From Google to App
            </Select.Option>
            <Select.Option value="to_google">
              ⬆️ From App to Google
            </Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="auto_sync"
          label="Automatic Sync"
          valuePropName="checked"
        >
          <Switch checkedChildren="On" unCheckedChildren="Off" />
        </Form.Item>

        <Form.Item
          name="sync_interval"
          label="Sync Interval (minutes)"
          help="How often to sync automatically"
        >
          <Select>
            <Select.Option value={5}>Every 5 minutes</Select.Option>
            <Select.Option value={15}>Every 15 minutes</Select.Option>
            <Select.Option value={30}>Every 30 minutes</Select.Option>
            <Select.Option value={60}>Every hour</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="sync_past_events"
          label="Sync Past Events"
          valuePropName="checked"
        >
          <Switch checkedChildren="Yes" unCheckedChildren="No" />
        </Form.Item>

        <Form.Item
          name="sync_future_events"
          label="Sync Future Events"
          valuePropName="checked"
        >
          <Switch checkedChildren="Yes" unCheckedChildren="No" />
        </Form.Item>

        <Form.Item
          name="conflict_resolution"
          label="Conflict Resolution"
          help="How to handle conflicts when the same event is modified in both places"
        >
          <Select>
            <Select.Option value="google_wins">Google Calendar Wins</Select.Option>
            <Select.Option value="local_wins">Local App Wins</Select.Option>
            <Select.Option value="manual">Manual Resolution</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

const SettingsPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { user } = useAppSelector(state => state.auth);
  const { 
    integration, 
    calendars, 
    calendarSyncs, 
    isLoading 
  } = useAppSelector(state => state.google);

  const [userForm] = Form.useForm();
  const [calendarSyncModalVisible, setCalendarSyncModalVisible] = useState(false);
  const [syncConflictsModalVisible, setSyncConflictsModalVisible] = useState(false);
  const [selectedCalendar, setSelectedCalendar] = useState<GoogleCalendar | null>(null);
  const [editingSync, setEditingSync] = useState<GoogleCalendarSync | null>(null);

  useEffect(() => {
    if (user) {
      userForm.setFieldsValue({
        name: user.name,
        email: user.email
      });
    }
  }, [user, userForm]);

  useEffect(() => {
    dispatch(fetchGoogleIntegration());
    if (integration?.enabled) {
      dispatch(fetchGoogleCalendars());
    }
  }, [dispatch, integration?.enabled]);

  const handleUpdateProfile = async () => {
    try {
      const values = await userForm.validateFields();
      await dispatch(updateUser(values)).unwrap();
      message.success('Profile updated successfully');
    } catch (error) {
      console.error('Failed to update profile:', error);
      message.error('Failed to update profile');
    }
  };

  const handleConnectGoogle = () => {
    // This would redirect to Google OAuth flow
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
    window.location.href = `${apiUrl}/api/v1/google/auth-url`;
  };

  const handleDisconnectGoogle = async () => {
    try {
      await dispatch(disconnectGoogle()).unwrap();
      message.success('Google Calendar disconnected successfully');
    } catch (error) {
      console.error('Failed to disconnect Google:', error);
      message.error('Failed to disconnect Google Calendar');
    }
  };

  const handleSetupCalendarSync = (calendar: GoogleCalendar) => {
    setSelectedCalendar(calendar);
    setEditingSync(null);
    setCalendarSyncModalVisible(true);
  };

  const handleEditCalendarSync = (sync: GoogleCalendarSync) => {
    const calendar = calendars.find(c => c.id === sync.calendar_id);
    setSelectedCalendar(calendar || null);
    setEditingSync(sync);
    setCalendarSyncModalVisible(true);
  };

  const handleDeleteCalendarSync = async (syncId: string) => {
    try {
      await dispatch(deleteCalendarSync(syncId)).unwrap();
      message.success('Calendar sync deleted successfully');
    } catch (error) {
      console.error('Failed to delete calendar sync:', error);
      message.error('Failed to delete calendar sync');
    }
  };

  const handleManualSync = async (syncId: string) => {
    try {
      await dispatch(manualSync(syncId)).unwrap();
      message.success('Manual sync completed successfully');
    } catch (error) {
      console.error('Failed to manual sync:', error);
      message.error('Failed to perform manual sync');
    }
  };

  const getSyncStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'green';
      case 'paused': return 'orange';
      case 'error': return 'red';
      case 'disabled': return 'default';
      default: return 'default';
    }
  };

  const getSyncDirectionIcon = (direction: string) => {
    switch (direction) {
      case 'bidirectional': return '↔️';
      case 'from_google': return '⬇️';
      case 'to_google': return '⬆️';
      default: return '↔️';
    }
  };

  return (
    <div>
      <Title level={2} style={{ marginBottom: 24 }}>Settings</Title>

      <Row gutter={[24, 24]}>
        {/* Profile Settings */}
        <Col xs={24} lg={12}>
          <Card
            title={
              <Space>
                <UserOutlined />
                Profile Settings
              </Space>
            }
          >
            <Form
              form={userForm}
              layout="vertical"
              onFinish={handleUpdateProfile}
            >
              <Form.Item
                name="name"
                label="Display Name"
                rules={[{ required: true, message: 'Please enter your name' }]}
              >
                <Input />
              </Form.Item>

              <Form.Item
                name="email"
                label="Email Address"
                rules={[
                  { required: true, message: 'Please enter your email' },
                  { type: 'email', message: 'Please enter a valid email' }
                ]}
              >
                <Input disabled />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={isLoading}>
                  Update Profile
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>

        {/* Notification Settings */}
        <Col xs={24} lg={12}>
          <Card
            title={
              <Space>
                <BellOutlined />
                Notification Settings
              </Space>
            }
          >
            <Space direction="vertical" style={{ width: '100%' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <Text strong>Email Notifications</Text>
                  <br />
                  <Text type="secondary">Receive email reminders for events</Text>
                </div>
                <Switch defaultChecked />
              </div>

              <Divider />

              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <Text strong>Daily Mood Reminders</Text>
                  <br />
                  <Text type="secondary">Get reminded to track your daily mood</Text>
                </div>
                <Switch defaultChecked />
              </div>

              <Divider />

              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <Text strong>Goal Progress Updates</Text>
                  <br />
                  <Text type="secondary">Weekly updates on your goal progress</Text>
                </div>
                <Switch defaultChecked />
              </div>

              <Divider />

              <div>
                <Text strong>Reminder Time</Text>
                <br />
                <Text type="secondary">When to send daily reminders</Text>
                <br />
                <TimePicker 
                  defaultValue={dayjs('09:00', 'HH:mm')} 
                  format="HH:mm" 
                  style={{ marginTop: 8 }}
                />
              </div>
            </Space>
          </Card>
        </Col>

        {/* Google Calendar Integration */}
        <Col xs={24}>
          <Card
            title={
              <Space>
                <GoogleOutlined />
                Google Calendar Integration
              </Space>
            }
            extra={
              integration?.enabled ? (
                <Space>
                  <Button 
                    icon={<ExclamationCircleOutlined />}
                    onClick={() => setSyncConflictsModalVisible(true)}
                  >
                    View Conflicts
                  </Button>
                  <Popconfirm
                    title="Are you sure you want to disconnect Google Calendar?"
                    onConfirm={handleDisconnectGoogle}
                    okText="Yes"
                    cancelText="No"
                  >
                    <Button danger icon={<DisconnectOutlined />}>
                      Disconnect
                    </Button>
                  </Popconfirm>
                </Space>
              ) : (
                <Button 
                  type="primary" 
                  icon={<LinkOutlined />}
                  onClick={handleConnectGoogle}
                >
                  Connect Google Calendar
                </Button>
              )
            }
          >
            {integration?.enabled ? (
              <div>
                <Alert
                  message="Google Calendar Connected"
                  description={`Connected as ${integration.email} (${integration.name})`}
                  type="success"
                  showIcon
                  style={{ marginBottom: 24 }}
                />

                {/* Available Calendars */}
                <div style={{ marginBottom: 24 }}>
                  <Title level={4}>Available Calendars</Title>
                  <List
                    dataSource={calendars}
                    renderItem={(calendar) => {
                      const existingSync = calendarSyncs.find(s => s.calendar_id === calendar.id);
                      return (
                        <List.Item
                          actions={[
                            existingSync ? (
                              <Space key="actions">
                                <Button
                                  type="text"
                                  icon={<SyncOutlined />}
                                  onClick={() => handleManualSync(existingSync.id)}
                                  loading={isLoading}
                                >
                                  Sync Now
                                </Button>
                                <Button
                                  type="text"
                                  icon={<EditOutlined />}
                                  onClick={() => handleEditCalendarSync(existingSync)}
                                >
                                  Edit
                                </Button>
                                <Popconfirm
                                  title="Are you sure you want to delete this sync?"
                                  onConfirm={() => handleDeleteCalendarSync(existingSync.id)}
                                >
                                  <Button
                                    type="text"
                                    danger
                                    icon={<DeleteOutlined />}
                                  />
                                </Popconfirm>
                              </Space>
                            ) : (
                              <Button
                                key="setup"
                                type="primary"
                                icon={<PlusOutlined />}
                                onClick={() => handleSetupCalendarSync(calendar)}
                              >
                                Setup Sync
                              </Button>
                            )
                          ]}
                        >
                          <List.Item.Meta
                            avatar={
                              <Avatar 
                                style={{ 
                                  backgroundColor: calendar.primary ? '#1677ff' : '#52c41a' 
                                }}
                              >
                                {calendar.primary ? '🏠' : '📅'}
                              </Avatar>
                            }
                            title={
                              <Space>
                                {calendar.summary}
                                {calendar.primary && <Tag color="blue">Primary</Tag>}
                                {existingSync && (
                                  <Tag color={getSyncStatusColor(existingSync.sync_status)}>
                                    {getSyncDirectionIcon(existingSync.sync_direction)} {existingSync.sync_status}
                                  </Tag>
                                )}
                              </Space>
                            }
                            description={
                              <div>
                                <Text type="secondary">{calendar.description || 'No description'}</Text>
                                {existingSync && existingSync.last_sync_at && (
                                  <div>
                                    <Text type="secondary" style={{ fontSize: 12 }}>
                                      Last sync: {dayjs(existingSync.last_sync_at).format('MMM DD, YYYY HH:mm')}
                                    </Text>
                                  </div>
                                )}
                                {existingSync?.last_sync_error && (
                                  <div>
                                    <Text type="danger" style={{ fontSize: 12 }}>
                                      Error: {existingSync.last_sync_error}
                                    </Text>
                                  </div>
                                )}
                              </div>
                            }
                          />
                        </List.Item>
                      );
                    }}
                  />
                </div>
              </div>
            ) : (
              <div style={{ textAlign: 'center', padding: '40px 0' }}>
                <GoogleOutlined style={{ fontSize: 48, color: '#ccc', marginBottom: 16 }} />
                <Paragraph type="secondary">
                  Connect your Google Calendar to enable two-way synchronization of events.
                  This allows you to manage your schedule from both platforms seamlessly.
                </Paragraph>
                <Button 
                  type="primary" 
                  size="large"
                  icon={<LinkOutlined />}
                  onClick={handleConnectGoogle}
                >
                  Connect Google Calendar
                </Button>
              </div>
            )}
          </Card>
        </Col>

        {/* Data & Privacy */}
        <Col xs={24} lg={12}>
          <Card
            title={
              <Space>
                <SecurityScanOutlined />
                Data & Privacy
              </Space>
            }
          >
            <Space direction="vertical" style={{ width: '100%' }}>
              <div>
                <Text strong>Data Export</Text>
                <br />
                <Text type="secondary">Download all your data</Text>
                <br />
                <Button style={{ marginTop: 8 }}>Export Data</Button>
              </div>

              <Divider />

              <div>
                <Text strong>Delete Account</Text>
                <br />
                <Text type="secondary">Permanently delete your account and all data</Text>
                <br />
                <Popconfirm
                  title="Are you sure you want to delete your account? This action cannot be undone."
                  okText="Yes, Delete"
                  cancelText="Cancel"
                  okButtonProps={{ danger: true }}
                >
                  <Button danger style={{ marginTop: 8 }}>
                    Delete Account
                  </Button>
                </Popconfirm>
              </div>
            </Space>
          </Card>
        </Col>

        {/* App Information */}
        <Col xs={24} lg={12}>
          <Card
            title="App Information"
          >
            <Space direction="vertical" style={{ width: '100%' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <Text type="secondary">Version:</Text>
                <Text>1.0.0</Text>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <Text type="secondary">Last Updated:</Text>
                <Text>{dayjs().format('MMM DD, YYYY')}</Text>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <Text type="secondary">Support:</Text>
                <Text>help@smartcalendar.app</Text>
              </div>

              <Divider />

              <div>
                <Button type="link" style={{ paddingLeft: 0 }}>
                  Privacy Policy
                </Button>
                <Button type="link">
                  Terms of Service
                </Button>
                <Button type="link">
                  Help Center
                </Button>
              </div>
            </Space>
          </Card>
        </Col>
      </Row>

      {/* Calendar Sync Modal */}
      <CalendarSyncModal
        visible={calendarSyncModalVisible}
        calendar={selectedCalendar}
        sync={editingSync}
        onCancel={() => {
          setCalendarSyncModalVisible(false);
          setSelectedCalendar(null);
          setEditingSync(null);
        }}
        onSuccess={() => {
          setCalendarSyncModalVisible(false);
          setSelectedCalendar(null);
          setEditingSync(null);
          // Refresh calendar syncs
        }}
      />

      {/* Sync Conflicts Modal */}
      <SyncConflictsModal
        visible={syncConflictsModalVisible}
        onCancel={() => setSyncConflictsModalVisible(false)}
      />
    </div>
  );
};

export default SettingsPage;