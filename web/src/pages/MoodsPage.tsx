import React, { useEffect, useState } from 'react';
import {
  Row,
  Col,
  Card,
  Typography,
  Space,
  Button,
  Calendar,
  Statistic,
  Select,
  DatePicker,
  Input,
  Form,
  Modal,
  Tag,
  Progress,
  Empty,
  Divider,
  Tooltip,
  message
} from '../utils/antd';
import {
  SmileOutlined,
  PlusOutlined,
  CalendarOutlined,
  RiseOutlined,
  BarChartOutlined,
  EditOutlined,
  EyeOutlined
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import {
  fetchMoods,
  createMood,
  updateMood,
  getTodayMood,
  getMoodStats
} from '../store/slices/moodsSlice';
import { Mood } from '../types/api';
import dayjs, { Dayjs } from 'dayjs';

const { Title, Text } = Typography;
const { TextArea } = Input;

interface MoodModalProps {
  visible: boolean;
  mood: Mood | null;
  selectedDate: string | null;
  onCancel: () => void;
  onSuccess: () => void;
}

const MoodModal: React.FC<MoodModalProps> = ({ visible, mood, selectedDate, onCancel, onSuccess }) => {
  const [form] = Form.useForm();
  const dispatch = useAppDispatch();
  const { isLoading } = useAppSelector(state => state.moods);

  const isEditing = !!mood;
  const targetDate = selectedDate || dayjs().format('YYYY-MM-DD');

  const moodLevels = [
    { value: 1, emoji: '😢', label: 'Very Sad', color: '#ff4d4f', description: 'Feeling very down or upset' },
    { value: 2, emoji: '😔', label: 'Sad', color: '#fa8c16', description: 'Feeling somewhat down' },
    { value: 3, emoji: '😐', label: 'Neutral', color: '#1677ff', description: 'Feeling okay, neither good nor bad' },
    { value: 4, emoji: '😊', label: 'Happy', color: '#52c41a', description: 'Feeling good and positive' },
    { value: 5, emoji: '😄', label: 'Very Happy', color: '#389e0d', description: 'Feeling great and energetic' }
  ];

  const moodTags = [
    'productive', 'stressed', 'relaxed', 'energetic', 'tired', 'focused',
    'anxious', 'confident', 'grateful', 'motivated', 'creative', 'social',
    'lonely', 'excited', 'calm', 'overwhelmed', 'optimistic', 'peaceful'
  ];

  useEffect(() => {
    if (visible) {
      if (mood) {
        form.setFieldsValue({
          level: mood.level,
          notes: mood.notes,
          tags: mood.tags
        });
      } else {
        form.setFieldsValue({
          level: 3,
          notes: '',
          tags: []
        });
      }
    } else {
      form.resetFields();
    }
  }, [visible, mood, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      const moodData = {
        date: targetDate,
        level: values.level,
        notes: values.notes || '',
        tags: values.tags || []
      };

      if (isEditing && mood) {
        await dispatch(updateMood({
          id: mood.id,
          data: moodData
        })).unwrap();
        message.success('Mood updated successfully');
      } else {
        await dispatch(createMood(moodData)).unwrap();
        message.success('Mood recorded successfully');
      }

      onSuccess();
    } catch (error) {
      console.error('Failed to save mood:', error);
      message.error('Failed to save mood');
    }
  };

  const selectedLevel = Form.useWatch('level', form);
  const selectedMood = moodLevels.find(m => m.value === selectedLevel);

  return (
    <Modal
      title={
        <Space>
          <SmileOutlined />
          {isEditing ? 'Edit Mood' : 'Record Your Mood'}
          <Text type="secondary">- {dayjs(targetDate).format('MMM DD, YYYY')}</Text>
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
          {isEditing ? 'Update Mood' : 'Record Mood'}
        </Button>
      ]}
      width={600}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="level"
          label="How are you feeling?"
          rules={[{ required: true, message: 'Please select your mood level' }]}
        >
          <div>
            <div style={{ 
              display: 'flex', 
              justifyContent: 'space-between', 
              marginBottom: 16,
              gap: 8
            }}>
              {moodLevels.map(mood => (
                <Tooltip key={mood.value} title={mood.description}>
                  <Button
                    size="large"
                    style={{
                      height: 80,
                      display: 'flex',
                      flexDirection: 'column',
                      alignItems: 'center',
                      justifyContent: 'center',
                      border: selectedLevel === mood.value ? `2px solid ${mood.color}` : '1px solid #d9d9d9',
                      backgroundColor: selectedLevel === mood.value ? `${mood.color}10` : 'white'
                    }}
                    onClick={() => form.setFieldsValue({ level: mood.value })}
                  >
                    <div style={{ fontSize: 24, marginBottom: 4 }}>{mood.emoji}</div>
                    <Text style={{ fontSize: 12, textAlign: 'center' }}>{mood.label}</Text>
                  </Button>
                </Tooltip>
              ))}
            </div>
            {selectedMood && (
              <div style={{ 
                textAlign: 'center', 
                padding: 12, 
                background: `${selectedMood.color}10`,
                borderRadius: 6,
                border: `1px solid ${selectedMood.color}30`
              }}>
                <Text style={{ color: selectedMood.color, fontWeight: 'bold' }}>
                  {selectedMood.emoji} {selectedMood.label}
                </Text>
                <br />
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {selectedMood.description}
                </Text>
              </div>
            )}
          </div>
        </Form.Item>

        <Form.Item
          name="tags"
          label="What's influencing your mood? (Optional)"
          help="Select tags that describe your current state"
        >
          <Select
            mode="multiple"
            placeholder="Select mood influences"
            style={{ width: '100%' }}
            maxTagCount={6}
            options={moodTags.map(tag => ({ label: tag, value: tag }))}
          />
        </Form.Item>

        <Form.Item
          name="notes"
          label="Additional Notes (Optional)"
          help="Reflect on your day or add any additional context"
        >
          <TextArea
            rows={3}
            placeholder="What happened today? How are you feeling about your progress towards your goals?"
            maxLength={500}
            showCount
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

const MoodsPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { moods, todayMood, stats, isLoading } = useAppSelector(state => state.moods);
  
  const [moodModalVisible, setMoodModalVisible] = useState(false);
  const [editingMood, setEditingMood] = useState<Mood | null>(null);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [calendarValue] = useState<Dayjs>(dayjs());
  const [viewMode, setViewMode] = useState<'calendar' | 'list'>('calendar');

  useEffect(() => {
    dispatch(getTodayMood());
    dispatch(getMoodStats({ days: 30 }));
    dispatch(fetchMoods({}));
  }, [dispatch]);

  const getMoodEmoji = (level: number) => {
    const emojis = ['😢', '😔', '😐', '😊', '😄'];
    return emojis[level - 1] || '😐';
  };

  const getMoodColor = (level: number) => {
    const colors = ['#ff4d4f', '#fa8c16', '#1677ff', '#52c41a', '#389e0d'];
    return colors[level - 1] || '#1677ff';
  };

  const getMoodLabel = (level: number) => {
    const labels = ['Very Sad', 'Sad', 'Neutral', 'Happy', 'Very Happy'];
    return labels[level - 1] || 'Neutral';
  };

  const handleRecordMood = (date?: string) => {
    setSelectedDate(date || dayjs().format('YYYY-MM-DD'));
    setEditingMood(null);
    setMoodModalVisible(true);
  };

  const handleEditMood = (mood: Mood) => {
    setEditingMood(mood);
    setSelectedDate(mood.date);
    setMoodModalVisible(true);
  };

  const handleViewMood = (mood: Mood) => {
    // Could implement a read-only view modal
    console.log('View mood:', mood);
  };

  const onCalendarSelect = (date: Dayjs) => {
    const dateStr = date.format('YYYY-MM-DD');
    const existingMood = moods.find(m => m.date === dateStr);
    
    if (existingMood) {
      handleEditMood(existingMood);
    } else {
      handleRecordMood(dateStr);
    }
  };

  const dateCellRender = (date: Dayjs) => {
    const dateStr = date.format('YYYY-MM-DD');
    const mood = moods.find(m => m.date === dateStr);
    
    if (mood) {
      return (
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: 20 }}>{getMoodEmoji(mood.level)}</div>
        </div>
      );
    }
    return null;
  };

  const monthCellRender = (date: Dayjs) => {
    const monthMoods = moods.filter(m => dayjs(m.date).month() === date.month());
    const avgMood = monthMoods.length > 0 
      ? monthMoods.reduce((sum, m) => sum + m.level, 0) / monthMoods.length 
      : 0;
    
    if (avgMood > 0) {
      return (
        <div style={{ textAlign: 'center', fontSize: 12 }}>
          <div>{getMoodEmoji(Math.round(avgMood))}</div>
          <div>{avgMood.toFixed(1)}</div>
        </div>
      );
    }
    return null;
  };

  const recentMoods = moods
    .sort((a, b) => dayjs(b.date).valueOf() - dayjs(a.date).valueOf())
    .slice(0, 10);

  return (
    <div>
      {/* Header */}
      <div style={{ marginBottom: 24 }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center', 
          marginBottom: 16 
        }}>
          <Title level={2} style={{ margin: 0 }}>
            Mood Tracking
          </Title>
          <Space>
            <Select
              value={viewMode}
              onChange={setViewMode}
              style={{ width: 120 }}
            >
              <Select.Option value="calendar">Calendar</Select.Option>
              <Select.Option value="list">List View</Select.Option>
            </Select>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              size="large"
              onClick={() => handleRecordMood()}
            >
              Record Mood
            </Button>
          </Space>
        </div>

        {/* Today's Mood Quick Action */}
        <Card size="small" style={{ marginBottom: 16 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <div style={{ fontSize: 32 }}>
                {todayMood ? getMoodEmoji(todayMood.level) : '😐'}
              </div>
              <div>
                <Text strong>Today's Mood</Text>
                <br />
                <Text type="secondary">
                  {todayMood 
                    ? `${getMoodLabel(todayMood.level)} (${todayMood.level}/5)`
                    : 'Not recorded yet'
                  }
                </Text>
              </div>
            </Space>
            <Button
              type={todayMood ? 'default' : 'primary'}
              onClick={() => handleRecordMood()}
            >
              {todayMood ? 'Update Today' : 'Record Today'}
            </Button>
          </div>
        </Card>

        {/* Statistics */}
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="Days Tracked"
                value={stats?.total || 0}
                prefix={<CalendarOutlined />}
                valueStyle={{ color: '#1677ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="Average Mood"
                value={stats?.average || 0}
                precision={1}
                suffix="/5"
                prefix={<SmileOutlined />}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="Trend"
                value={stats?.trend === 'up' ? '↗️' : stats?.trend === 'down' ? '↘️' : '➡️'}
                prefix={<RiseOutlined />}
                valueStyle={{ 
                  color: stats?.trend === 'up' ? '#52c41a' : 
                         stats?.trend === 'down' ? '#ff4d4f' : '#1677ff'
                }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="This Week"
                value={stats?.weekAverage || 0}
                precision={1}
                suffix="/5"
                prefix={<BarChartOutlined />}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
        </Row>
      </div>

      <Row gutter={[16, 16]}>
        {/* Main Content */}
        <Col xs={24} lg={16}>
          <Card title={viewMode === 'calendar' ? 'Mood Calendar' : 'Mood History'} loading={isLoading}>
            {viewMode === 'calendar' ? (
              <Calendar
                value={calendarValue}
                onSelect={onCalendarSelect}
                dateCellRender={dateCellRender}
                monthCellRender={monthCellRender}
                fullscreen={false}
              />
            ) : (
              <div>
                {recentMoods.length === 0 ? (
                  <Empty
                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                    description="No mood records yet"
                  >
                    <Button type="primary" onClick={() => handleRecordMood()}>
                      Record Your First Mood
                    </Button>
                  </Empty>
                ) : (
                  <div style={{ maxHeight: 600, overflowY: 'auto' }}>
                    {recentMoods.map((mood, index) => (
                      <div key={mood.id}>
                        <div style={{ 
                          display: 'flex', 
                          justifyContent: 'space-between', 
                          alignItems: 'flex-start',
                          padding: '12px 0'
                        }}>
                          <Space align="start">
                            <div style={{ fontSize: 32 }}>{getMoodEmoji(mood.level)}</div>
                            <div>
                              <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
                                <Text strong>{dayjs(mood.date).format('MMM DD, YYYY')}</Text>
                                <Tag color={getMoodColor(mood.level)}>
                                  {getMoodLabel(mood.level)} ({mood.level}/5)
                                </Tag>
                              </div>
                              {mood.notes && (
                                <Text type="secondary" style={{ display: 'block', marginBottom: 4 }}>
                                  {mood.notes}
                                </Text>
                              )}
                              {mood.tags.length > 0 && (
                                <Space wrap>
                                  {mood.tags.map(tag => (
                                    <Tag key={tag}>{tag}</Tag>
                                  ))}
                                </Space>
                              )}
                            </div>
                          </Space>
                          <Space>
                            <Button
                              type="text"
                              size="small"
                              icon={<EyeOutlined />}
                              onClick={() => handleViewMood(mood)}
                            />
                            <Button
                              type="text"
                              size="small"
                              icon={<EditOutlined />}
                              onClick={() => handleEditMood(mood)}
                            />
                          </Space>
                        </div>
                        {index < recentMoods.length - 1 && <Divider />}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </Card>
        </Col>

        {/* Sidebar */}
        <Col xs={24} lg={8}>
          {/* Weekly Overview */}
          <Card title="This Week" style={{ marginBottom: 16 }}>
            {stats && stats.weeklyData ? (
              <div>
                {['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((day, index) => {
                  const mood = stats.weeklyData?.[index];
                  return (
                    <div key={day} style={{ 
                      display: 'flex', 
                      justifyContent: 'space-between', 
                      alignItems: 'center',
                      marginBottom: 8
                    }}>
                      <Text>{day}</Text>
                      <div>
                        {mood ? (
                          <Space>
                            <span style={{ fontSize: 18 }}>{getMoodEmoji(mood.level)}</span>
                            <Text type="secondary">{mood.level}/5</Text>
                          </Space>
                        ) : (
                          <Text type="secondary">-</Text>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <Text type="secondary">Track your mood daily to see weekly patterns</Text>
            )}
          </Card>

          {/* Mood Distribution */}
          <Card title="Mood Distribution">
            {stats && stats.distribution ? (
              <div>
                {[5, 4, 3, 2, 1].map(level => {
                  const count = stats.distribution?.[level] || 0;
                  const percentage = stats.total > 0 ? (count / stats.total) * 100 : 0;
                  return (
                    <div key={level} style={{ marginBottom: 12 }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                        <Space>
                          <span>{getMoodEmoji(level)}</span>
                          <Text>{getMoodLabel(level)}</Text>
                        </Space>
                        <Text type="secondary">{count}</Text>
                      </div>
                      <Progress
                        percent={percentage}
                        size="small"
                        strokeColor={getMoodColor(level)}
                        showInfo={false}
                      />
                    </div>
                  );
                })}
              </div>
            ) : (
              <Text type="secondary">Record more moods to see distribution</Text>
            )}
          </Card>
        </Col>
      </Row>

      {/* Mood Modal */}
      <MoodModal
        visible={moodModalVisible}
        mood={editingMood}
        selectedDate={selectedDate}
        onCancel={() => {
          setMoodModalVisible(false);
          setEditingMood(null);
          setSelectedDate(null);
        }}
        onSuccess={() => {
          setMoodModalVisible(false);
          setEditingMood(null);
          setSelectedDate(null);
          dispatch(fetchMoods({}));
          dispatch(getTodayMood());
          dispatch(getMoodStats({ days: 30 }));
        }}
      />
    </div>
  );
};

export default MoodsPage;