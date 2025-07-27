import React, { useState, useEffect } from 'react';
import {
  Modal,
  Button,
  Card,
  List,
  Space,
  Typography,
  Tag,
  Progress,
  Alert,
  Divider,
  Form,
  InputNumber,
  TimePicker,
  Checkbox,
  Row,
  Col,
  Badge,
  Tooltip,
  message
} from 'antd';
import {
  ClockCircleOutlined,
  CalendarOutlined,
  ThunderboltOutlined,
  CheckCircleOutlined,
  WarningOutlined,
  SettingOutlined,
  PlusOutlined
} from '@ant-design/icons';
import { useAppSelector, useAppDispatch } from '../../hooks/redux';
import { createEvent } from '../../store/slices/eventsSlice';
import timeScheduler, {
  generateSchedulingSuggestions,
  createScheduledEvents,
  SchedulingSuggestion,
  SchedulingPreferences
} from '../../utils/timeScheduler';
import { Goal, Task, Event } from '../../types/api';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

interface TimeSchedulerModalProps {
  visible: boolean;
  goal: Goal;
  tasks: Task[];
  onCancel: () => void;
  onSuccess: () => void;
}

const TimeSchedulerModal: React.FC<TimeSchedulerModalProps> = ({
  visible,
  goal,
  tasks,
  onCancel,
  onSuccess
}) => {
  const dispatch = useAppDispatch();
  const { events } = useAppSelector(state => state.events);
  const { user } = useAppSelector(state => state.auth);
  
  const [suggestions, setSuggestions] = useState<SchedulingSuggestion[]>([]);
  const [preferences, setPreferences] = useState<SchedulingPreferences>(timeScheduler.DEFAULT_PREFERENCES);
  const [isGenerating, setIsGenerating] = useState(false);
  const [showPreferences, setShowPreferences] = useState(false);
  const [selectedSuggestions, setSelectedSuggestions] = useState<Set<string>>(new Set());
  const [isScheduling, setIsScheduling] = useState(false);

  useEffect(() => {
    if (visible && goal) {
      generateSuggestions();
    }
  }, [visible, goal]);

  const generateSuggestions = async () => {
    setIsGenerating(true);
    try {
      const allGoals = [goal]; // Focus on current goal
      const goalTasks = tasks.filter(task => task.goal_id === goal.id);
      
      const suggestions = generateSchedulingSuggestions(
        allGoals,
        goalTasks,
        events,
        preferences,
        14 // 2 weeks lookahead
      );
      
      setSuggestions(suggestions);
      
      if (suggestions.length === 0) {
        message.warning('No available time slots found for the next 2 weeks. Try adjusting your preferences.');
      }
    } catch (error) {
      message.error('Failed to generate scheduling suggestions');
    }
    setIsGenerating(false);
  };

  const handleScheduleSelected = async () => {
    if (selectedSuggestions.size === 0) {
      message.warning('Please select at least one time slot to schedule');
      return;
    }

    if (!user?.id) {
      message.error('User not authenticated');
      return;
    }

    setIsScheduling(true);
    try {
      const selectedSuggestionsList = suggestions.filter(s => 
        selectedSuggestions.has(s.goal.id)
      );
      
      const newEvents = createScheduledEvents(selectedSuggestionsList, user.id);
      
      // Create events one by one
      for (const eventData of newEvents) {
        await dispatch(createEvent(eventData)).unwrap();
      }
      
      message.success(`Successfully scheduled ${newEvents.length} time slots for your goal`);
      onSuccess();
    } catch (error) {
      message.error('Failed to schedule time slots');
    }
    setIsScheduling(false);
  };

  const toggleSuggestionSelection = (goalId: string) => {
    const newSelected = new Set(selectedSuggestions);
    if (newSelected.has(goalId)) {
      newSelected.delete(goalId);
    } else {
      newSelected.add(goalId);
    }
    setSelectedSuggestions(newSelected);
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'critical': return '#ff4d4f';
      case 'high': return '#fa8c16';
      case 'medium': return '#1677ff';
      case 'low': return '#52c41a';
      default: return '#d9d9d9';
    }
  };

  const formatTimeSlot = (slot: { start: { format: (format: string) => string }; end: { format: (format: string) => string }; duration: number }) => {
    return `${slot.start.format('MMM DD, HH:mm')} - ${slot.end.format('HH:mm')} (${slot.duration}min)`;
  };

  const PreferencesForm = () => (
    <Card size="small" title={
      <Space>
        <SettingOutlined />
        Scheduling Preferences
      </Space>
    }>
      <Form layout="vertical">
        <Row gutter={16}>
          <Col span={12}>
            <Form.Item label="Working Hours">
              <Space>
                <TimePicker
                  value={dayjs(preferences.workingHours.start, 'HH:mm')}
                  format="HH:mm"
                  onChange={(time) => {
                    if (time) {
                      setPreferences(prev => ({
                        ...prev,
                        workingHours: {
                          ...prev.workingHours,
                          start: time.format('HH:mm')
                        }
                      }));
                    }
                  }}
                />
                <Text>to</Text>
                <TimePicker
                  value={dayjs(preferences.workingHours.end, 'HH:mm')}
                  format="HH:mm"
                  onChange={(time) => {
                    if (time) {
                      setPreferences(prev => ({
                        ...prev,
                        workingHours: {
                          ...prev.workingHours,
                          end: time.format('HH:mm')
                        }
                      }));
                    }
                  }}
                />
              </Space>
            </Form.Item>
          </Col>
          
          <Col span={12}>
            <Form.Item label="Session Duration (minutes)">
              <Space>
                <InputNumber
                  min={15}
                  max={180}
                  value={preferences.minSessionDuration}
                  onChange={(value) => {
                    if (value) {
                      setPreferences(prev => ({
                        ...prev,
                        minSessionDuration: value
                      }));
                    }
                  }}
                  addonBefore="Min"
                />
                <InputNumber
                  min={30}
                  max={480}
                  value={preferences.maxSessionDuration}
                  onChange={(value) => {
                    if (value) {
                      setPreferences(prev => ({
                        ...prev,
                        maxSessionDuration: value
                      }));
                    }
                  }}
                  addonBefore="Max"
                />
              </Space>
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="Working Days">
          <Checkbox.Group
            value={preferences.workingDays}
            onChange={(checkedValues) => {
              setPreferences(prev => ({
                ...prev,
                workingDays: checkedValues as number[]
              }));
            }}
          >
            <Row>
              <Col span={8}><Checkbox value={1}>Monday</Checkbox></Col>
              <Col span={8}><Checkbox value={2}>Tuesday</Checkbox></Col>
              <Col span={8}><Checkbox value={3}>Wednesday</Checkbox></Col>
              <Col span={8}><Checkbox value={4}>Thursday</Checkbox></Col>
              <Col span={8}><Checkbox value={5}>Friday</Checkbox></Col>
              <Col span={8}><Checkbox value={6}>Saturday</Checkbox></Col>
              <Col span={8}><Checkbox value={0}>Sunday</Checkbox></Col>
            </Row>
          </Checkbox.Group>
        </Form.Item>

        <Button 
          type="primary" 
          icon={<ThunderboltOutlined />}
          onClick={generateSuggestions}
          loading={isGenerating}
        >
          Regenerate Suggestions
        </Button>
      </Form>
    </Card>
  );

  return (
    <Modal
      title={
        <Space>
          <ClockCircleOutlined />
          Schedule Time for: {goal.title}
        </Space>
      }
      open={visible}
      onCancel={onCancel}
      width={900}
      footer={[
        <Button 
          key="preferences" 
          onClick={() => setShowPreferences(!showPreferences)}
        >
          <SettingOutlined /> {showPreferences ? 'Hide' : 'Show'} Preferences
        </Button>,
        <Button key="cancel" onClick={onCancel}>
          Cancel
        </Button>,
        <Button
          key="schedule"
          type="primary"
          loading={isScheduling}
          disabled={selectedSuggestions.size === 0}
          onClick={handleScheduleSelected}
        >
          <PlusOutlined /> Schedule {selectedSuggestions.size} Time Slot{selectedSuggestions.size !== 1 ? 's' : ''}
        </Button>
      ]}
    >
      {/* Goal Summary */}
      <Card size="small" style={{ marginBottom: 16 }}>
        <Row gutter={16}>
          <Col span={12}>
            <Space direction="vertical" size="small">
              <div>
                <Text strong>Goal: </Text>
                <Text>{goal.title}</Text>
              </div>
              <div>
                <Text strong>Priority: </Text>
                <Tag color={getPriorityColor(goal.priority)}>{goal.priority}</Tag>
              </div>
              {goal.deadline && (
                <div>
                  <Text strong>Deadline: </Text>
                  <Text>{dayjs(goal.deadline).format('MMM DD, YYYY')}</Text>
                  {dayjs(goal.deadline).diff(dayjs(), 'days') <= 7 && (
                    <Tag color="red" style={{ marginLeft: 8 }}>Urgent!</Tag>
                  )}
                </div>
              )}
            </Space>
          </Col>
          <Col span={12}>
            <Space direction="vertical" size="small">
              <div>
                <Text strong>Progress: </Text>
                <Progress 
                  percent={goal.progress} 
                  size="small" 
                  style={{ width: '100px' }}
                />
              </div>
              <div>
                <Text strong>Tasks: </Text>
                <Text>{tasks.filter(t => t.status !== 'completed').length} pending</Text>
              </div>
            </Space>
          </Col>
        </Row>
      </Card>

      {/* Preferences Panel */}
      {showPreferences && (
        <>
          {PreferencesForm()}
          <Divider />
        </>
      )}

      {/* Suggestions List */}
      {isGenerating ? (
        <Card style={{ textAlign: 'center', padding: '40px' }}>
          <ThunderboltOutlined style={{ fontSize: '24px', marginBottom: '16px' }} />
          <div>Analyzing your calendar and generating optimal time slots...</div>
        </Card>
      ) : suggestions.length > 0 ? (
        <List
          dataSource={suggestions}
          renderItem={(suggestion) => (
            <List.Item
              style={{
                cursor: 'pointer',
                backgroundColor: selectedSuggestions.has(suggestion.goal.id) ? '#e6f7ff' : 'transparent',
                padding: '16px',
                border: selectedSuggestions.has(suggestion.goal.id) ? '2px solid #1677ff' : '1px solid #f0f0f0',
                borderRadius: '8px',
                marginBottom: '8px'
              }}
              onClick={() => toggleSuggestionSelection(suggestion.goal.id)}
            >
              <List.Item.Meta
                avatar={
                  <Badge 
                    status={selectedSuggestions.has(suggestion.goal.id) ? 'processing' : 'default'}
                    style={{ marginTop: '8px' }}
                  />
                }
                title={
                  <Space>
                    <Text strong>{suggestion.reason}</Text>
                    <Tag color={getPriorityColor(suggestion.priority)}>
                      {suggestion.priority}
                    </Tag>
                    {suggestion.deadline && dayjs(suggestion.deadline).diff(dayjs(), 'days') <= 7 && (
                      <Tag color="red" icon={<WarningOutlined />}>
                        Urgent
                      </Tag>
                    )}
                  </Space>
                }
                description={
                  <div>
                    <Text type="secondary">
                      {Math.round(suggestion.totalTimeNeeded / 60 * 10) / 10}h total across {suggestion.suggestedSlots.length} sessions
                    </Text>
                    {suggestion.task && (
                      <div style={{ marginTop: '4px' }}>
                        <Text type="secondary">Next task: {suggestion.task.title}</Text>
                      </div>
                    )}
                  </div>
                }
              />
              <div style={{ minWidth: '300px' }}>
                {suggestion.suggestedSlots.map((slot, index) => (
                  <div key={index} style={{ marginBottom: '4px' }}>
                    <CalendarOutlined style={{ marginRight: '8px', color: '#1677ff' }} />
                    <Text style={{ fontSize: '12px' }}>
                      {formatTimeSlot(slot)}
                    </Text>
                  </div>
                ))}
              </div>
            </List.Item>
          )}
        />
      ) : (
        <Alert
          message="No Available Time Slots"
          description="No suitable time slots found for the next 2 weeks. Try adjusting your working hours or preferences."
          type="warning"
          showIcon
          style={{ textAlign: 'center' }}
          action={
            <Button size="small" onClick={() => setShowPreferences(true)}>
              Adjust Preferences
            </Button>
          }
        />
      )}

      {selectedSuggestions.size > 0 && (
        <Alert
          message={`${selectedSuggestions.size} time slot${selectedSuggestions.size !== 1 ? 's' : ''} selected`}
          description="Selected time slots will be added to your calendar as tentative events."
          type="info"
          showIcon
          style={{ marginTop: '16px' }}
        />
      )}
    </Modal>
  );
};

export default TimeSchedulerModal;