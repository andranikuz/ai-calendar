import React, { useEffect, useState } from 'react';
import {
  Card,
  Typography,
  Space,
  Progress,
  Tag,
  Button,
  List,
  Divider,
  Avatar,
  Tooltip,
  Input,
  Form,
  Modal,
  DatePicker,
  Select,
  message,
  Popconfirm,
  Empty
} from 'antd';
import {
  CloseOutlined,
  EditOutlined,
  PlusOutlined,
  CheckCircleOutlined,
  DeleteOutlined,
  CalendarOutlined,
  FlagOutlined,
  AimOutlined,
  TrophyOutlined,
  ClockCircleOutlined
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { updateGoalProgress } from '../../store/slices/goalsSlice';
import { Goal, Task, Milestone } from '../../types/api';
import TaskTreeView from './TaskTreeView';
import dayjs from 'dayjs';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;

interface GoalDetailPanelProps {
  goal: Goal;
  onClose: () => void;
  onEdit: () => void;
}

interface TaskModalProps {
  visible: boolean;
  task: Task | null;
  goalId: string;
  onCancel: () => void;
  onSuccess: () => void;
}

interface MilestoneModalProps {
  visible: boolean;
  milestone: Milestone | null;
  goalId: string;
  onCancel: () => void;
  onSuccess: () => void;
}

const TaskModal: React.FC<TaskModalProps> = ({ visible, task, goalId, onCancel, onSuccess }) => {
  const [form] = Form.useForm();
  const isEditing = !!task;

  useEffect(() => {
    if (visible) {
      if (task) {
        form.setFieldsValue({
          title: task.title,
          description: task.description,
          deadline: task.deadline ? dayjs(task.deadline) : null,
          priority: task.priority
        });
      } else {
        form.resetFields();
      }
    }
  }, [visible, task, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      // Here you would dispatch to create/update task
      console.log('Task data:', { goalId, ...values });
      message.success(isEditing ? 'Task updated successfully' : 'Task created successfully');
      onSuccess();
    } catch (error) {
      console.error('Failed to save task:', error);
      message.error('Failed to save task');
    }
  };

  return (
    <Modal
      title={isEditing ? 'Edit Task' : 'Create New Task'}
      open={visible}
      onCancel={onCancel}
      footer={[
        <Button key="cancel" onClick={onCancel}>Cancel</Button>,
        <Button key="submit" type="primary" onClick={handleSubmit}>
          {isEditing ? 'Update' : 'Create'}
        </Button>
      ]}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="title"
          label="Task Title"
          rules={[{ required: true, message: 'Please enter task title' }]}
        >
          <Input placeholder="Enter task title" />
        </Form.Item>
        
        <Form.Item name="description" label="Description">
          <TextArea rows={3} placeholder="Enter task description (optional)" />
        </Form.Item>
        
        <Form.Item name="deadline" label="Deadline (Optional)">
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>
        
        <Form.Item name="priority" label="Priority">
          <Select placeholder="Select priority">
            <Select.Option value="low">Low</Select.Option>
            <Select.Option value="medium">Medium</Select.Option>
            <Select.Option value="high">High</Select.Option>
            <Select.Option value="critical">Critical</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

const MilestoneModal: React.FC<MilestoneModalProps> = ({ visible, milestone, goalId, onCancel, onSuccess }) => {
  const [form] = Form.useForm();
  const isEditing = !!milestone;

  useEffect(() => {
    if (visible) {
      if (milestone) {
        form.setFieldsValue({
          title: milestone.title,
          description: milestone.description,
          deadline: milestone.deadline ? dayjs(milestone.deadline) : null,
          target_value: milestone.target_value
        });
      } else {
        form.resetFields();
      }
    }
  }, [visible, milestone, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      // Here you would dispatch to create/update milestone
      console.log('Milestone data:', { goalId, ...values });
      message.success(isEditing ? 'Milestone updated successfully' : 'Milestone created successfully');
      onSuccess();
    } catch (error) {
      console.error('Failed to save milestone:', error);
      message.error('Failed to save milestone');
    }
  };

  return (
    <Modal
      title={isEditing ? 'Edit Milestone' : 'Create New Milestone'}
      open={visible}
      onCancel={onCancel}
      footer={[
        <Button key="cancel" onClick={onCancel}>Cancel</Button>,
        <Button key="submit" type="primary" onClick={handleSubmit}>
          {isEditing ? 'Update' : 'Create'}
        </Button>
      ]}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="title"
          label="Milestone Title"
          rules={[{ required: true, message: 'Please enter milestone title' }]}
        >
          <Input placeholder="Enter milestone title" />
        </Form.Item>
        
        <Form.Item name="description" label="Description">
          <TextArea rows={3} placeholder="Enter milestone description (optional)" />
        </Form.Item>
        
        <Form.Item name="deadline" label="Target Date">
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>
        
        <Form.Item name="target_value" label="Target Value">
          <Input placeholder="e.g., 50% completion, 10kg weight loss" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

const GoalDetailPanel: React.FC<GoalDetailPanelProps> = ({ goal, onClose, onEdit }) => {
  const dispatch = useAppDispatch();
  const { isLoading } = useAppSelector(state => state.goals);
  
  const [taskModalVisible, setTaskModalVisible] = useState(false);
  const [milestoneModalVisible, setMilestoneModalVisible] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [editingMilestone, setEditingMilestone] = useState<Milestone | null>(null);

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

  const handleProgressUpdate = async (newProgress: number) => {
    try {
      await dispatch(updateGoalProgress({
        id: goal.id,
        progress: newProgress
      }));
      message.success('Progress updated successfully');
    } catch (error) {
      console.error('Failed to update progress:', error);
      message.error('Failed to update progress');
    }
  };

  const handleCreateTask = () => {
    setEditingTask(null);
    setTaskModalVisible(true);
  };

  const handleEditTask = (task: Task) => {
    setEditingTask(task);
    setTaskModalVisible(true);
  };

  const handleCreateMilestone = () => {
    setEditingMilestone(null);
    setMilestoneModalVisible(true);
  };

  const handleEditMilestone = (milestone: Milestone) => {
    setEditingMilestone(milestone);
    setMilestoneModalVisible(true);
  };

  const handleCompleteTask = async (taskId: string) => {
    try {
      // Here you would dispatch to complete task
      console.log('Completing task:', taskId);
      message.success('Task completed successfully');
    } catch (error) {
      console.error('Failed to complete task:', error);
      message.error('Failed to complete task');
    }
  };

  const handleTaskCreate = async (taskData: Omit<Task, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      // TODO: Implement task creation API call
      console.log('Create task:', taskData);
      message.success('Task created successfully!');
    } catch (error) {
      console.error('Failed to create task:', error);
      message.error('Failed to create task');
      throw error;
    }
  };

  const handleTaskUpdate = async (taskId: string, updates: Partial<Task>) => {
    try {
      // TODO: Implement task update API call
      console.log('Update task:', taskId, updates);
      message.success('Task updated successfully!');
    } catch (error) {
      console.error('Failed to update task:', error);
      message.error('Failed to update task');
      throw error;
    }
  };

  const handleTaskDelete = async (taskId: string) => {
    try {
      // TODO: Implement task deletion API call
      console.log('Delete task:', taskId);
      message.success('Task deleted successfully!');
    } catch (error) {
      console.error('Failed to delete task:', error);
      message.error('Failed to delete task');
      throw error;
    }
  };

  const handleCompleteMilestone = async (milestoneId: string) => {
    try {
      // Here you would dispatch to complete milestone
      console.log('Completing milestone:', milestoneId);
      message.success('Milestone completed successfully');
    } catch (error) {
      console.error('Failed to complete milestone:', error);
      message.error('Failed to complete milestone');
    }
  };

  // Mock data for tasks and milestones - replace with real data from API
  const mockTasks: Task[] = [
    {
      id: '1',
      title: 'Research marathon training plans',
      description: 'Find a suitable 16-week training program',
      status: 'active',
      priority: 'high',
      deadline: dayjs().add(1, 'week').toISOString(),
      goal_id: goal.id,
      order_index: 1,
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    },
    {
      id: '1-1',
      title: 'Compare beginner programs',
      description: 'Research Couch to 5K, C25K+ programs',
      status: 'completed',
      priority: 'medium',
      goal_id: goal.id,
      parent_task_id: '1',
      order_index: 1,
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    },
    {
      id: '1-2',
      title: 'Choose final training plan',
      description: 'Select the most suitable plan based on research',
      status: 'pending',
      priority: 'high',
      goal_id: goal.id,
      parent_task_id: '1',
      order_index: 2,
      deadline: dayjs().add(3, 'days').toISOString(),
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    },
    {
      id: '2',
      title: 'Buy running gear',
      description: 'Get proper equipment for training',
      status: 'in_progress',
      priority: 'medium',
      goal_id: goal.id,
      order_index: 2,
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    },
    {
      id: '2-1',
      title: 'Buy running shoes',
      description: 'Get proper running shoes for training',
      status: 'completed',
      priority: 'high',
      goal_id: goal.id,
      parent_task_id: '2',
      order_index: 1,
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    },
    {
      id: '2-2',
      title: 'Buy moisture-wicking clothes',
      description: 'Get appropriate running attire',
      status: 'pending',
      priority: 'low',
      goal_id: goal.id,
      parent_task_id: '2',
      order_index: 2,
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    },
    {
      id: '3',
      title: 'Start training routine',
      description: 'Begin following the selected training plan',
      status: 'pending',
      priority: 'critical',
      goal_id: goal.id,
      order_index: 3,
      deadline: dayjs().add(2, 'weeks').toISOString(),
      estimated_duration: 120, // 2 hours per session
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    }
  ];

  const mockMilestones: Milestone[] = [
    {
      id: '1',
      title: 'Run 5K without stopping',
      description: 'First milestone towards marathon goal',
      status: 'completed',
      deadline: dayjs().subtract(1, 'month').toISOString(),
      target_value: '5K continuous run',
      goal_id: goal.id,
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    },
    {
      id: '2',
      title: 'Complete half marathon',
      description: 'Run 21K as preparation for full marathon',
      status: 'active',
      deadline: dayjs().add(2, 'months').toISOString(),
      target_value: '21K race completion',
      goal_id: goal.id,
      created_at: dayjs().toISOString(),
      updated_at: dayjs().toISOString()
    }
  ];

  return (
    <Card
      style={{ height: 'fit-content' }}
      actions={[
        <Button
          key="edit"
          type="text"
          icon={<EditOutlined />}
          onClick={onEdit}
        >
          Edit Goal
        </Button>,
        <Button
          key="close"
          type="text"
          icon={<CloseOutlined />}
          onClick={onClose}
        >
          Close
        </Button>
      ]}
    >
      {/* Goal Header */}
      <div style={{ marginBottom: 24 }}>
        <Space style={{ marginBottom: 16 }}>
          <Avatar
            size={48}
            style={{
              backgroundColor: getPriorityColor(goal.priority),
              fontSize: 20
            }}
          >
            {getCategoryIcon(goal.category)}
          </Avatar>
          <div>
            <Title level={4} style={{ margin: 0 }}>
              {goal.title}
            </Title>
            <Space>
              <Tag color={getStatusColor(goal.status)}>
                {goal.status.toUpperCase()}
              </Tag>
              <Tag color={getPriorityColor(goal.priority)}>
                {goal.priority.toUpperCase()}
              </Tag>
            </Space>
          </div>
        </Space>

        <Paragraph type="secondary">
          {goal.description}
        </Paragraph>

        {/* Progress */}
        <div style={{ marginBottom: 16 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
            <Text strong>Progress</Text>
            <Text>{goal.progress}%</Text>
          </div>
          <Progress
            percent={goal.progress}
            status={goal.progress === 100 ? 'success' : 'active'}
            strokeColor={goal.progress === 100 ? '#52c41a' : '#1677ff'}
          />
          <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
            {[25, 50, 75, 100].map(value => (
              <Button
                key={value}
                size="small"
                type={goal.progress >= value ? 'primary' : 'default'}
                onClick={() => handleProgressUpdate(value)}
                disabled={isLoading}
              >
                {value}%
              </Button>
            ))}
          </div>
        </div>

        {/* Goal Info */}
        <Space direction="vertical" size="small" style={{ width: '100%' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Text type="secondary">Category:</Text>
            <Text>{goal.category}</Text>
          </div>
          {goal.deadline && (
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <Text type="secondary">Deadline:</Text>
              <Space>
                <CalendarOutlined />
                <Text>{dayjs(goal.deadline).format('MMM DD, YYYY')}</Text>
                {dayjs(goal.deadline).isBefore(dayjs()) && (
                  <Tag color="red">Overdue</Tag>
                )}
              </Space>
            </div>
          )}
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Text type="secondary">Created:</Text>
            <Text>{dayjs(goal.created_at).format('MMM DD, YYYY')}</Text>
          </div>
        </Space>
      </div>

      <Divider />

      {/* Milestones Section */}
      <div style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
          <Title level={5} style={{ margin: 0 }}>
            <TrophyOutlined /> Milestones
          </Title>
          <Button
            type="primary"
            size="small"
            icon={<PlusOutlined />}
            onClick={handleCreateMilestone}
          >
            Add Milestone
          </Button>
        </div>

        {mockMilestones.length === 0 ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description="No milestones yet"
          >
            <Button type="primary" onClick={handleCreateMilestone}>
              Create Milestone
            </Button>
          </Empty>
        ) : (
          <List
            size="small"
            dataSource={mockMilestones}
            renderItem={(milestone) => (
              <List.Item
                actions={[
                  milestone.status !== 'completed' && (
                    <Tooltip title="Mark as completed">
                      <Button
                        type="text"
                        size="small"
                        icon={<CheckCircleOutlined />}
                        onClick={() => handleCompleteMilestone(milestone.id)}
                      />
                    </Tooltip>
                  ),
                  <Tooltip title="Edit milestone">
                    <Button
                      type="text"
                      size="small"
                      icon={<EditOutlined />}
                      onClick={() => handleEditMilestone(milestone)}
                    />
                  </Tooltip>,
                  <Popconfirm
                    title="Are you sure you want to delete this milestone?"
                    onConfirm={() => console.log('Delete milestone:', milestone.id)}
                  >
                    <Tooltip title="Delete milestone">
                      <Button
                        type="text"
                        size="small"
                        icon={<DeleteOutlined />}
                        danger
                      />
                    </Tooltip>
                  </Popconfirm>
                ].filter(Boolean)}
              >
                <List.Item.Meta
                  avatar={
                    <Avatar
                      size="small"
                      style={{
                        backgroundColor: milestone.status === 'completed' ? '#52c41a' : '#1677ff'
                      }}
                    >
                      {milestone.status === 'completed' ? <CheckCircleOutlined /> : <FlagOutlined />}
                    </Avatar>
                  }
                  title={
                    <Space>
                      <Text strong={milestone.status !== 'completed'} delete={milestone.status === 'completed'}>
                        {milestone.title}
                      </Text>
                      {milestone.deadline && (
                        <Tag color={milestone.status === 'completed' ? 'green' : 'blue'}>
                          <ClockCircleOutlined /> {dayjs(milestone.deadline).format('MMM DD')}
                        </Tag>
                      )}
                    </Space>
                  }
                  description={milestone.description}
                />
              </List.Item>
            )}
          />
        )}
      </div>

      <Divider />

      {/* Tasks Section with Tree View */}
      <TaskTreeView
        tasks={mockTasks}
        goalId={goal.id}
        onTaskCreate={handleTaskCreate}
        onTaskUpdate={handleTaskUpdate}
        onTaskDelete={handleTaskDelete}
      />

      {/* Modals */}
      <TaskModal
        visible={taskModalVisible}
        task={editingTask}
        goalId={goal.id}
        onCancel={() => {
          setTaskModalVisible(false);
          setEditingTask(null);
        }}
        onSuccess={() => {
          setTaskModalVisible(false);
          setEditingTask(null);
          // Refresh tasks data
        }}
      />

      <MilestoneModal
        visible={milestoneModalVisible}
        milestone={editingMilestone}
        goalId={goal.id}
        onCancel={() => {
          setMilestoneModalVisible(false);
          setEditingMilestone(null);
        }}
        onSuccess={() => {
          setMilestoneModalVisible(false);
          setEditingMilestone(null);
          // Refresh milestones data
        }}
      />
    </Card>
  );
};

export default GoalDetailPanel;