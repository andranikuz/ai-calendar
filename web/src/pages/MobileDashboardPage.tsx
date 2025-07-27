import React, { useEffect } from 'react';
import { 
  Row, 
  Col, 
  Card, 
  Typography, 
  Space, 
  Button, 
  List, 
  Grid,
  Empty 
} from '../utils/antd';
import { 
  CalendarOutlined, 
  FlagOutlined, 
  SmileOutlined, 
  PlusOutlined,
  TrophyOutlined,
  ClockCircleOutlined,
  RightOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { fetchGoals } from '../store/slices/goalsSlice';
import { fetchEvents } from '../store/slices/eventsSlice';
import { getTodayMood, getMoodStats } from '../store/slices/moodsSlice';
import { useNotifications } from '../components/Common/NotificationProvider';
import MobileCard from '../components/Common/MobileCard';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { useBreakpoint } = Grid;

const MobileDashboardPage: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const screens = useBreakpoint();
  const { showInfo } = useNotifications();
  
  const { goals } = useAppSelector(state => state.goals);
  const { events } = useAppSelector(state => state.events);
  const { todayMood } = useAppSelector(state => state.moods);
  const { user } = useAppSelector(state => state.auth);

  const isMobile = !screens.md;

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
    .slice(0, isMobile ? 3 : 5);

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

  const handleQuickAction = (action: string) => {
    switch (action) {
      case 'goal':
        navigate('/goals');
        showInfo('Quick Tip', 'Click the + button to create a new goal!');
        break;
      case 'event':
        navigate('/calendar');
        showInfo('Quick Tip', 'Click on any date to create a new event!');
        break;
      case 'mood':
        navigate('/moods');
        showInfo('Quick Tip', 'Recording your daily mood helps track patterns!');
        break;
      default:
        break;
    }
  };

  const quickActions = [
    {
      key: 'goal',
      title: 'New Goal',
      icon: '🎯',
      description: 'Set a new goal',
      color: '#1677ff'
    },
    {
      key: 'event',
      title: 'New Event',
      icon: '📅',
      description: 'Schedule an event',
      color: '#52c41a'
    },
    {
      key: 'mood',
      title: 'Record Mood',
      icon: '😊',
      description: 'How are you feeling?',
      color: '#fa8c16'
    }
  ];

  const StatCard: React.FC<{
    title: string;
    value: number | string;
    icon: React.ReactNode;
    color: string;
    onClick?: () => void;
  }> = ({ title, value, icon, color, onClick }) => (
    <Card 
      size="small" 
      onClick={onClick}
      style={{ 
        cursor: onClick ? 'pointer' : 'default',
        borderRadius: 12,
        border: `1px solid ${color}20`,
        background: `${color}05`
      }}
      bodyStyle={{ padding: isMobile ? 12 : 16 }}
    >
      <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
        <div style={{
          width: isMobile ? 40 : 48,
          height: isMobile ? 40 : 48,
          borderRadius: '50%',
          background: `${color}20`,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: color,
          fontSize: isMobile ? 16 : 18
        }}>
          {icon}
        </div>
        <div>
          <Text type="secondary" style={{ fontSize: isMobile ? 11 : 12, display: 'block' }}>
            {title}
          </Text>
          <Text strong style={{ fontSize: isMobile ? 18 : 20, color }}>
            {value}
          </Text>
        </div>
      </div>
    </Card>
  );

  return (
    <div style={{ padding: isMobile ? 0 : undefined }}>
      {/* Welcome Section */}
      <div style={{ marginBottom: 24 }}>
        <Space direction="vertical" size="small" style={{ width: '100%' }}>
          <Title level={isMobile ? 3 : 2} style={{ marginBottom: 4 }}>
            {getGreeting()}, {user?.name}! 👋
          </Title>
          <Text type="secondary" style={{ fontSize: isMobile ? 14 : 16 }}>
            Here's your productivity overview for today
          </Text>
        </Space>
      </div>

      {/* Today's Mood Quick Check */}
      {isMobile && (
        <Card 
          size="small" 
          style={{ 
            marginBottom: 16, 
            borderRadius: 12,
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            color: 'white',
            border: 'none'
          }}
          bodyStyle={{ padding: 16 }}
          onClick={() => navigate('/moods')}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div>
              <Text strong style={{ color: 'white', fontSize: 16 }}>
                Today's Mood
              </Text>
              <div style={{ marginTop: 4 }}>
                <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: 12 }}>
                  {todayMood ? 'Recorded' : 'Not recorded yet'}
                </Text>
              </div>
            </div>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 32, marginBottom: 4 }}>
                {getMoodEmoji(todayMood?.level)}
              </div>
              <RightOutlined style={{ color: 'rgba(255,255,255,0.6)' }} />
            </div>
          </div>
        </Card>
      )}

      {/* Stats Cards */}
      <Row gutter={[isMobile ? 8 : 16, isMobile ? 8 : 16]} style={{ marginBottom: 24 }}>
        <Col xs={12} sm={6}>
          <StatCard
            title="Active Goals"
            value={activeGoals.length}
            icon={<FlagOutlined />}
            color="#1677ff"
            onClick={() => navigate('/goals')}
          />
        </Col>
        <Col xs={12} sm={6}>
          <StatCard
            title="Completed"
            value={completedGoals.length}
            icon={<TrophyOutlined />}
            color="#52c41a"
            onClick={() => navigate('/goals')}
          />
        </Col>
        <Col xs={12} sm={6}>
          <StatCard
            title="Today's Events"
            value={todayEvents.length}
            icon={<CalendarOutlined />}
            color="#fa8c16"
            onClick={() => navigate('/calendar')}
          />
        </Col>
        <Col xs={12} sm={6}>
          <StatCard
            title="Mood"
            value={todayMood ? getMoodEmoji(todayMood.level) : '😐'}
            icon={<SmileOutlined />}
            color="#eb2f96"
            onClick={() => navigate('/moods')}
          />
        </Col>
      </Row>

      {/* Quick Actions (Mobile) */}
      {isMobile && (
        <Card 
          title="Quick Actions" 
          size="small"
          style={{ marginBottom: 16, borderRadius: 12 }}
          bodyStyle={{ padding: 12 }}
        >
          <Row gutter={8}>
            {quickActions.map((action) => (
              <Col span={8} key={action.key}>
                <div
                  role="button"
                  tabIndex={0}
                  onClick={() => handleQuickAction(action.key)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                      e.preventDefault();
                      handleQuickAction(action.key);
                    }
                  }}
                  style={{
                    textAlign: 'center',
                    padding: 12,
                    borderRadius: 8,
                    background: `${action.color}10`,
                    border: `1px solid ${action.color}20`,
                    cursor: 'pointer',
                    transition: 'all 0.3s ease'
                  }}
                  aria-label={`${action.title}: ${action.description}`}
                >
                  <div style={{ fontSize: 24, marginBottom: 4 }}>
                    {action.icon}
                  </div>
                  <Text strong style={{ fontSize: 12, color: action.color, display: 'block' }}>
                    {action.title}
                  </Text>
                  <Text type="secondary" style={{ fontSize: 10 }}>
                    {action.description}
                  </Text>
                </div>
              </Col>
            ))}
          </Row>
        </Card>
      )}

      <Row gutter={[isMobile ? 0 : 16, 16]}>
        {/* Active Goals */}
        <Col xs={24} lg={12}>
          <Card 
            title="Active Goals" 
            size={isMobile ? 'small' : 'default'}
            style={{ borderRadius: isMobile ? 12 : 8 }}
            extra={
              !isMobile ? (
                <Button 
                  type="primary" 
                  icon={<PlusOutlined />}
                  onClick={() => navigate('/goals')}
                >
                  New Goal
                </Button>
              ) : (
                <Button 
                  type="text" 
                  size="small"
                  onClick={() => navigate('/goals')}
                >
                  View All
                </Button>
              )
            }
          >
            {activeGoals.length === 0 ? (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description="No active goals yet"
                style={{ margin: isMobile ? '20px 0' : '40px 0' }}
              >
                {!isMobile && (
                  <Button 
                    type="primary" 
                    onClick={() => navigate('/goals')}
                  >
                    Create your first goal
                  </Button>
                )}
              </Empty>
            ) : (
              <List
                dataSource={activeGoals.slice(0, isMobile ? 3 : 5)}
                renderItem={(goal) => (
                  <List.Item style={{ padding: isMobile ? '8px 0' : '12px 0' }}>
                    <MobileCard
                      title={goal.title}
                      description={goal.description}
                      status={goal.status}
                      statusColor="blue"
                      priority={goal.priority}
                      priorityColor={
                        goal.priority === 'critical' ? 'red' :
                        goal.priority === 'high' ? 'orange' :
                        goal.priority === 'medium' ? 'blue' : 'green'
                      }
                      progress={goal.progress}
                      onClick={() => {
                        navigate('/goals');
                        // TODO: Select this goal
                      }}
                      style={{ marginBottom: 0 }}
                    />
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
            size={isMobile ? 'small' : 'default'}
            style={{ borderRadius: isMobile ? 12 : 8 }}
            extra={
              !isMobile ? (
                <Button 
                  type="primary" 
                  icon={<PlusOutlined />}
                  onClick={() => navigate('/calendar')}
                >
                  New Event
                </Button>
              ) : (
                <Button 
                  type="text" 
                  size="small"
                  onClick={() => navigate('/calendar')}
                >
                  View All
                </Button>
              )
            }
          >
            {upcomingEvents.length === 0 ? (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description="No upcoming events"
                style={{ margin: isMobile ? '20px 0' : '40px 0' }}
              >
                {!isMobile && (
                  <Button 
                    type="primary" 
                    onClick={() => navigate('/calendar')}
                  >
                    Add an event
                  </Button>
                )}
              </Empty>
            ) : (
              <List
                dataSource={upcomingEvents}
                renderItem={(event) => (
                  <List.Item style={{ padding: isMobile ? '8px 0' : '12px 0' }}>
                    <List.Item.Meta
                      avatar={<ClockCircleOutlined style={{ color: '#1677ff', fontSize: isMobile ? 16 : 18 }} />}
                      title={
                        <Text strong style={{ fontSize: isMobile ? 14 : 16 }}>
                          {event.title}
                        </Text>
                      }
                      description={
                        <Space direction="vertical" size="small">
                          <Text type="secondary" style={{ fontSize: isMobile ? 12 : 14 }}>
                            {dayjs(event.start_time).format('MMM DD, YYYY HH:mm')}
                          </Text>
                          {event.location && (
                            <Text type="secondary" style={{ fontSize: isMobile ? 11 : 12 }}>
                              📍 {event.location}
                            </Text>
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
      </Row>
    </div>
  );
};

export default MobileDashboardPage;