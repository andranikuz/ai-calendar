import React, { useEffect, useState } from 'react';
import {
  Row,
  Col,
  Card,
  Button,
  Typography,
  Space,
  Progress,
  Tag,
  List,
  Avatar,
  Dropdown,
  Input,
  Select,
  Empty,
  Statistic,
  Divider
} from 'antd';
import {
  PlusOutlined,
  FlagOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  MoreOutlined,
  SearchOutlined,
  TrophyOutlined,
  AimOutlined,
  RocketOutlined,
  EditOutlined,
  DeleteOutlined
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { 
  fetchGoals, 
  setCurrentGoal, 
  updateGoalProgress,
  deleteGoal 
} from '../store/slices/goalsSlice';
import GoalModal from '../components/Goals/GoalModal';
import GoalDetailPanel from '../components/Goals/GoalDetailPanel';
import { Goal } from '../types/api';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { Search } = Input;

const GoalsPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { goals, currentGoal, isLoading } = useAppSelector(state => state.goals);
  
  const [goalModalVisible, setGoalModalVisible] = useState(false);
  const [editingGoal, setEditingGoal] = useState<Goal | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [filterCategory, setFilterCategory] = useState<string>('all');

  useEffect(() => {
    dispatch(fetchGoals());
  }, [dispatch]);

  // Filter goals based on search and filters
  const filteredGoals = goals.filter(goal => {
    const matchesSearch = goal.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         goal.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = filterStatus === 'all' || goal.status === filterStatus;
    const matchesCategory = filterCategory === 'all' || goal.category === filterCategory;
    
    return matchesSearch && matchesStatus && matchesCategory;
  });

  // Calculate statistics
  const stats = {
    total: goals.length,
    active: goals.filter(g => g.status === 'active').length,
    completed: goals.filter(g => g.status === 'completed').length,
    averageProgress: goals.length > 0 
      ? Math.round(goals.reduce((sum, g) => sum + g.progress, 0) / goals.length)
      : 0
  };

  const handleCreateGoal = () => {
    setEditingGoal(null);
    setGoalModalVisible(true);
  };

  const handleEditGoal = (goal: Goal) => {
    setEditingGoal(goal);
    setGoalModalVisible(true);
  };

  const handleDeleteGoal = async (goalId: string) => {
    try {
      await dispatch(deleteGoal(goalId)).unwrap();
    } catch (error) {
      console.error('Failed to delete goal:', error);
    }
  };

  const handleViewGoal = (goal: Goal) => {
    dispatch(setCurrentGoal(goal));
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'blue';
      case 'completed': return 'green';
      case 'paused': return 'orange';
      case 'cancelled': return 'red';
      default: return 'default';
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'critical': return 'red';
      case 'high': return 'orange';
      case 'medium': return 'blue';
      case 'low': return 'green';
      default: return 'default';
    }
  };

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'health': return '🏥';
      case 'career': return '💼';
      case 'education': return '📚';
      case 'personal': return '👤';
      case 'financial': return '💰';
      case 'relationship': return '❤️';
      default: return '🎯';
    }
  };

  const getGoalActions = (goal: Goal) => [
    {
      key: 'view',
      label: 'View Details',
      icon: <AimOutlined />,
      onClick: () => handleViewGoal(goal),
    },
    {
      key: 'edit',
      label: 'Edit Goal',
      icon: <EditOutlined />,
      onClick: () => handleEditGoal(goal),
    },
    {
      key: 'delete',
      label: 'Delete Goal',
      icon: <DeleteOutlined />,
      danger: true,
      onClick: () => handleDeleteGoal(goal.id),
    },
  ];

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
            Goals Management
          </Title>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            size="large"
            onClick={handleCreateGoal}
          >
            New Goal
          </Button>
        </div>

        {/* Statistics Cards */}
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="Total Goals"
                value={stats.total}
                prefix={<FlagOutlined />}
                valueStyle={{ color: '#1677ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="Active Goals"
                value={stats.active}
                prefix={<RocketOutlined />}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="Completed"
                value={stats.completed}
                prefix={<TrophyOutlined />}
                valueStyle={{ color: '#faad14' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={6}>
            <Card>
              <Statistic
                title="Avg Progress"
                value={stats.averageProgress}
                suffix="%"
                prefix={<AimOutlined />}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
        </Row>

        {/* Filters */}
        <Card>
          <Row gutter={[16, 16]} align="middle">
            <Col xs={24} md={8}>
              <Search
                placeholder="Search goals..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                prefix={<SearchOutlined />}
                allowClear
              />
            </Col>
            <Col xs={12} md={4}>
              <Select
                placeholder="Status"
                value={filterStatus}
                onChange={setFilterStatus}
                style={{ width: '100%' }}
              >
                <Select.Option value="all">All Status</Select.Option>
                <Select.Option value="active">Active</Select.Option>
                <Select.Option value="completed">Completed</Select.Option>
                <Select.Option value="paused">Paused</Select.Option>
                <Select.Option value="cancelled">Cancelled</Select.Option>
              </Select>
            </Col>
            <Col xs={12} md={4}>
              <Select
                placeholder="Category"
                value={filterCategory}
                onChange={setFilterCategory}
                style={{ width: '100%' }}
              >
                <Select.Option value="all">All Categories</Select.Option>
                <Select.Option value="health">Health</Select.Option>
                <Select.Option value="career">Career</Select.Option>
                <Select.Option value="education">Education</Select.Option>
                <Select.Option value="personal">Personal</Select.Option>
                <Select.Option value="financial">Financial</Select.Option>
                <Select.Option value="relationship">Relationship</Select.Option>
              </Select>
            </Col>
          </Row>
        </Card>
      </div>

      <Row gutter={[16, 16]}>
        {/* Goals List */}
        <Col xs={24} lg={currentGoal ? 16 : 24}>
          <Card title="Your Goals" loading={isLoading}>
            {filteredGoals.length === 0 ? (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description={
                  goals.length === 0 
                    ? "No goals yet. Create your first goal to get started!"
                    : "No goals match your search criteria"
                }
              >
                {goals.length === 0 && (
                  <Button type="primary" onClick={handleCreateGoal}>
                    Create Goal
                  </Button>
                )}
              </Empty>
            ) : (
              <List
                itemLayout="vertical"
                dataSource={filteredGoals}
                renderItem={(goal) => (
                  <List.Item
                    key={goal.id}
                    actions={[
                      <Button
                        key="view"
                        type="link"
                        onClick={() => handleViewGoal(goal)}
                      >
                        View Details
                      </Button>,
                      <Dropdown
                        key="more"
                        menu={{ items: getGoalActions(goal) }}
                        trigger={['click']}
                      >
                        <Button type="text" icon={<MoreOutlined />} />
                      </Dropdown>
                    ]}
                    style={{
                      border: '1px solid #f0f0f0',
                      borderRadius: 8,
                      padding: 16,
                      marginBottom: 16,
                      background: goal.id === currentGoal?.id ? '#f6ffed' : 'white'
                    }}
                  >
                    <List.Item.Meta
                      avatar={
                        <Avatar
                          size={48}
                          style={{ 
                            backgroundColor: getPriorityColor(goal.priority),
                            fontSize: 20
                          }}
                        >
                          {getCategoryIcon(goal.category)}
                        </Avatar>
                      }
                      title={
                        <Space>
                          <Text strong style={{ fontSize: 16 }}>
                            {goal.title}
                          </Text>
                          <Tag color={getStatusColor(goal.status)}>
                            {goal.status.toUpperCase()}
                          </Tag>
                          <Tag color={getPriorityColor(goal.priority)}>
                            {goal.priority.toUpperCase()}
                          </Tag>
                        </Space>
                      }
                      description={
                        <Space direction="vertical" style={{ width: '100%' }}>
                          <Text type="secondary">{goal.description}</Text>
                          <div>
                            <Text strong>Progress: </Text>
                            <Progress 
                              percent={goal.progress} 
                              size="small" 
                              style={{ width: 200 }}
                              status={goal.progress === 100 ? 'success' : 'active'}
                            />
                          </div>
                          {goal.deadline && (
                            <Space>
                              <CalendarOutlined />
                              <Text type="secondary">
                                Due: {dayjs(goal.deadline).format('MMM DD, YYYY')}
                              </Text>
                              {dayjs(goal.deadline).isBefore(dayjs()) && (
                                <Tag color="red">Overdue</Tag>
                              )}
                            </Space>
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

        {/* Goal Detail Panel */}
        {currentGoal && (
          <Col xs={24} lg={8}>
            <GoalDetailPanel
              goal={currentGoal}
              onClose={() => dispatch(setCurrentGoal(null))}
              onEdit={() => handleEditGoal(currentGoal)}
            />
          </Col>
        )}
      </Row>

      {/* Goal Modal */}
      <GoalModal
        visible={goalModalVisible}
        goal={editingGoal}
        onCancel={() => {
          setGoalModalVisible(false);
          setEditingGoal(null);
        }}
        onSuccess={() => {
          setGoalModalVisible(false);
          setEditingGoal(null);
          dispatch(fetchGoals());
        }}
      />
    </div>
  );
};

export default GoalsPage;