import React, { useEffect } from 'react';
import { Row, Col, Card, Typography, Space, Button, Progress, List, Badge } from '../utils/antd';
import { 
  CalendarOutlined, 
  FlagOutlined, 
  SmileOutlined, 
  PlusOutlined,
  TrophyOutlined,
  ClockCircleOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { fetchGoals } from '../store/slices/goalsSlice';
import { fetchEvents } from '../store/slices/eventsSlice';
import { getTodayMood, getMoodStats } from '../store/slices/moodsSlice';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

const DashboardPage: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  
  const { goals } = useAppSelector(state => state.goals);
  const { events } = useAppSelector(state => state.events);
  const { todayMood, stats } = useAppSelector(state => state.moods);
  const { user } = useAppSelector(state => state.auth);

  useEffect(() => {
    // Fetch dashboard data
    dispatch(fetchGoals());
    dispatch(fetchEvents({}));
    dispatch(getTodayMood());
    dispatch(getMoodStats({ days: 7 }));
  }, [dispatch]);

  // Calculate stats
  const activeGoals = goals.filter(g => g.status === 'active');
  const completedGoals = goals.filter(g => g.status === 'completed');
  const todayEvents = events.filter(event => 
    dayjs(event.start_time).isSame(dayjs(), 'day')
  );
  const upcomingEvents = events
    .filter(event => dayjs(event.start_time).isAfter(dayjs()))
    .slice(0, 5);

  const getGreeting = () => {
    const hour = dayjs().hour();
    if (hour < 12) return 'Good morning';
    if (hour < 17) return 'Good afternoon';
    return 'Good evening';
  };

  const getMoodEmoji = (level?: number) => {
    if (!level) return '😐';
    const emojis = ['😢', '😔', '😐', '😊', '😄'];
    return emojis[level - 1] || '😐';
  };

  return (
    <div>
      {/* Welcome Section */}
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ marginBottom: 8 }}>
          {getGreeting()}, {user?.name}! 👋
        </Title>
        <Text type="secondary">
          Here's your productivity overview for today
        </Text>
      </div>

      {/* Stats Cards */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Space>
              <div style={{ 
                background: '#1677ff', 
                borderRadius: '50%', 
                width: 40, 
                height: 40, 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'center' 
              }}>
                <FlagOutlined style={{ color: 'white', fontSize: 18 }} />
              </div>
              <div>
                <Text type="secondary">Active Goals</Text>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>
                  {activeGoals.length}
                </div>
              </div>
            </Space>
          </Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Space>
              <div style={{ 
                background: '#52c41a', 
                borderRadius: '50%', 
                width: 40, 
                height: 40, 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'center' 
              }}>
                <TrophyOutlined style={{ color: 'white', fontSize: 18 }} />
              </div>
              <div>
                <Text type="secondary">Completed Goals</Text>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>
                  {completedGoals.length}
                </div>
              </div>
            </Space>
          </Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Space>
              <div style={{ 
                background: '#fa8c16', 
                borderRadius: '50%', 
                width: 40, 
                height: 40, 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'center' 
              }}>
                <CalendarOutlined style={{ color: 'white', fontSize: 18 }} />
              </div>
              <div>
                <Text type="secondary">Today's Events</Text>
                <div style={{ fontSize: 24, fontWeight: 'bold' }}>
                  {todayEvents.length}
                </div>
              </div>
            </Space>
          </Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Space>
              <div style={{ 
                background: '#eb2f96', 
                borderRadius: '50%', 
                width: 40, 
                height: 40, 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'center' 
              }}>
                <SmileOutlined style={{ color: 'white', fontSize: 18 }} />
              </div>
              <div>
                <Text type="secondary">Today's Mood</Text>
                <div style={{ fontSize: 24 }}>
                  {getMoodEmoji(todayMood?.level)}
                </div>
              </div>
            </Space>
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        {/* Active Goals */}
        <Col xs={24} lg={12}>
          <Card 
            title="Active Goals" 
            extra={
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => navigate('/goals')}
              >
                New Goal
              </Button>
            }
          >
            {activeGoals.length === 0 ? (
              <div style={{ textAlign: 'center', padding: '40px 0' }}>
                <Text type="secondary">No active goals yet</Text>
                <br />
                <Button 
                  type="link" 
                  onClick={() => navigate('/goals')}
                >
                  Create your first goal
                </Button>
              </div>
            ) : (
              <List
                dataSource={activeGoals.slice(0, 5)}
                renderItem={(goal) => (
                  <List.Item>
                    <div style={{ width: '100%' }}>
                      <div style={{ 
                        display: 'flex', 
                        justifyContent: 'space-between', 
                        marginBottom: 8 
                      }}>
                        <Text strong>{goal.title}</Text>
                        <Badge 
                          color={
                            goal.priority === 'critical' ? 'red' :
                            goal.priority === 'high' ? 'orange' :
                            goal.priority === 'medium' ? 'blue' : 'green'
                          } 
                          text={goal.priority}
                        />
                      </div>
                      <Progress 
                        percent={goal.progress} 
                        size="small"
                        status={goal.progress === 100 ? 'success' : 'active'}
                      />
                      {goal.deadline && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          Due: {dayjs(goal.deadline).format('MMM DD, YYYY')}
                        </Text>
                      )}
                    </div>
                  </List.Item>
                )}
              />
            )}
          </Card>
        </Col>

        {/* Upcoming Events */}
        <Col xs={24} lg={12}>
          <Card 
            title="Upcoming Events" 
            extra={
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => navigate('/calendar')}
              >
                New Event
              </Button>
            }
          >
            {upcomingEvents.length === 0 ? (
              <div style={{ textAlign: 'center', padding: '40px 0' }}>
                <Text type="secondary">No upcoming events</Text>
                <br />
                <Button 
                  type="link" 
                  onClick={() => navigate('/calendar')}
                >
                  Add an event
                </Button>
              </div>
            ) : (
              <List
                dataSource={upcomingEvents}
                renderItem={(event) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={<ClockCircleOutlined style={{ color: '#1677ff' }} />}
                      title={event.title}
                      description={
                        <Space direction="vertical" size="small">
                          <Text type="secondary">
                            {dayjs(event.start_time).format('MMM DD, YYYY HH:mm')}
                          </Text>
                          {event.location && (
                            <Text type="secondary">{event.location}</Text>
                          )}
                        </Space>
                      }
                    />
                  </List.Item>
                )}
              />
            )}
          </Card>
        </Col>

        {/* Mood Overview */}
        <Col xs={24}>
          <Card title="Weekly Mood Trends">
            {stats ? (
              <Row gutter={16}>
                <Col span={8}>
                  <div style={{ textAlign: 'center' }}>
                    <div style={{ fontSize: 32, marginBottom: 8 }}>
                      {getMoodEmoji(Math.round(stats.average))}
                    </div>
                    <Text type="secondary">Average Mood</Text>
                    <div style={{ fontSize: 18, fontWeight: 'bold' }}>
                      {stats.average.toFixed(1)}/5
                    </div>
                  </div>
                </Col>
                <Col span={8}>
                  <div style={{ textAlign: 'center' }}>
                    <div style={{ fontSize: 24, marginBottom: 8, color: '#52c41a' }}>
                      {stats.total}
                    </div>
                    <Text type="secondary">Days Tracked</Text>
                  </div>
                </Col>
                <Col span={8}>
                  <div style={{ textAlign: 'center' }}>
                    <div style={{ 
                      fontSize: 24, 
                      marginBottom: 8,
                      color: stats.trend === 'up' ? '#52c41a' : 
                            stats.trend === 'down' ? '#ff4d4f' : '#1677ff'
                    }}>
                      {stats.trend === 'up' ? '📈' : 
                       stats.trend === 'down' ? '📉' : '➡️'}
                    </div>
                    <Text type="secondary">Trend</Text>
                  </div>
                </Col>
              </Row>
            ) : (
              <div style={{ textAlign: 'center', padding: '40px 0' }}>
                <Text type="secondary">Track your mood to see trends</Text>
                <br />
                <Button 
                  type="link" 
                  onClick={() => navigate('/moods')}
                >
                  Record your mood
                </Button>
              </div>
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default DashboardPage;