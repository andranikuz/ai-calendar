import React, { useEffect } from 'react';
import {
  Modal,
  Form,
  Input,
  Select,
  DatePicker,
  Slider,
  Button,
  Space,
  message,
  Typography,
  Divider,
  Alert
} from 'antd';
import { 
  FlagOutlined, 
  CalendarOutlined, 
  AimOutlined,
  BulbOutlined 
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { createGoal, updateGoal } from '../../store/slices/goalsSlice';
import { Goal } from '../../types/api';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { Text } = Typography;

interface GoalModalProps {
  visible: boolean;
  goal: Goal | null;
  onCancel: () => void;
  onSuccess: () => void;
}

const GoalModal: React.FC<GoalModalProps> = ({
  visible,
  goal,
  onCancel,
  onSuccess
}) => {
  const [form] = Form.useForm();
  const dispatch = useAppDispatch();
  const { isLoading } = useAppSelector(state => state.goals);

  const isEditing = !!goal;

  useEffect(() => {
    if (visible) {
      if (goal) {
        // Editing existing goal
        form.setFieldsValue({
          title: goal.title,
          description: goal.description,
          category: goal.category,
          priority: goal.priority,
          deadline: goal.deadline ? dayjs(goal.deadline) : null,
          progress: goal.progress,
        });
      } else {
        // Creating new goal
        form.setFieldsValue({
          title: '',
          description: '',
          category: 'personal',
          priority: 'medium',
          deadline: null,
          progress: 0,
        });
      }
    } else {
      form.resetFields();
    }
  }, [visible, goal, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      const goalData = {
        title: values.title,
        description: values.description,
        category: values.category,
        priority: values.priority,
        deadline: values.deadline ? values.deadline.toISOString() : undefined,
        progress: values.progress || 0,
      };

      if (isEditing && goal) {
        await dispatch(updateGoal({
          id: goal.id,
          data: goalData
        })).unwrap();
        message.success('Goal updated successfully');
      } else {
        await dispatch(createGoal(goalData)).unwrap();
        message.success('Goal created successfully');
      }

      onSuccess();
    } catch (error) {
      console.error('Failed to save goal:', error);
      message.error('Failed to save goal');
    }
  };

  const categoryOptions = [
    { value: 'health', label: '🏥 Health', description: 'Physical and mental wellness goals' },
    { value: 'career', label: '💼 Career', description: 'Professional development and work goals' },
    { value: 'education', label: '📚 Education', description: 'Learning and skill development goals' },
    { value: 'personal', label: '👤 Personal', description: 'Personal growth and self-improvement' },
    { value: 'financial', label: '💰 Financial', description: 'Money, savings, and investment goals' },
    { value: 'relationship', label: '❤️ Relationship', description: 'Family, friends, and social goals' },
  ];

  const priorityOptions = [
    { value: 'low', label: 'Low', color: '#52c41a' },
    { value: 'medium', label: 'Medium', color: '#1677ff' },
    { value: 'high', label: 'High', color: '#fa8c16' },
    { value: 'critical', label: 'Critical', color: '#ff4d4f' },
  ];

  const smartGoalTips = [
    "📌 Specific: Be clear about what you want to achieve",
    "📏 Measurable: Include metrics to track progress",
    "🎯 Achievable: Set realistic and attainable goals",
    "💡 Relevant: Align with your values and long-term objectives",
    "⏰ Time-bound: Set a clear deadline"
  ];

  return (
    <Modal
      title={
        <Space>
          <FlagOutlined />
          {isEditing ? 'Edit Goal' : 'Create New Goal'}
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
          {isEditing ? 'Update Goal' : 'Create Goal'}
        </Button>
      ]}
      width={700}
      destroyOnClose
    >
      <Alert
        message="SMART Goal Framework"
        description="Create goals that are Specific, Measurable, Achievable, Relevant, and Time-bound"
        type="info"
        icon={<BulbOutlined />}
        style={{ marginBottom: 24 }}
        showIcon
      />

      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        <Form.Item
          name="title"
          label="Goal Title"
          rules={[
            { required: true, message: 'Please enter a goal title' },
            { max: 255, message: 'Title must be less than 255 characters' }
          ]}
        >
          <Input 
            placeholder="e.g., Run a marathon, Learn Spanish, Save $10,000"
            size="large"
          />
        </Form.Item>

        <Form.Item
          name="description"
          label="Description"
          rules={[
            { required: true, message: 'Please enter a goal description' },
            { max: 1000, message: 'Description must be less than 1000 characters' }
          ]}
        >
          <TextArea
            rows={4}
            placeholder="Describe your goal in detail. What exactly do you want to achieve? Why is this important to you?"
            showCount
            maxLength={1000}
          />
        </Form.Item>

        <Form.Item
          name="category"
          label="Category"
          rules={[{ required: true, message: 'Please select a category' }]}
        >
          <Select size="large" placeholder="Select goal category">
            {categoryOptions.map(option => (
              <Select.Option key={option.value} value={option.value}>
                <div>
                  <div>{option.label}</div>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    {option.description}
                  </Text>
                </div>
              </Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="priority"
          label="Priority"
          rules={[{ required: true, message: 'Please select a priority' }]}
        >
          <Select size="large" placeholder="Select priority level">
            {priorityOptions.map(option => (
              <Select.Option key={option.value} value={option.value}>
                <Space>
                  <div
                    style={{
                      width: 12,
                      height: 12,
                      borderRadius: '50%',
                      backgroundColor: option.color
                    }}
                  />
                  {option.label}
                </Space>
              </Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="deadline"
          label={
            <Space>
              <CalendarOutlined />
              Target Deadline (Optional)
            </Space>
          }
        >
          <DatePicker
            size="large"
            style={{ width: '100%' }}
            placeholder="Select target completion date"
            disabledDate={(current) => current && current.isBefore(dayjs(), 'day')}
          />
        </Form.Item>

        {isEditing && (
          <Form.Item
            name="progress"
            label={
              <Space>
                <AimOutlined />
                Current Progress
              </Space>
            }
          >
            <div>
              <Slider
                min={0}
                max={100}
                marks={{
                  0: '0%',
                  25: '25%',
                  50: '50%',
                  75: '75%',
                  100: '100%'
                }}
                tooltip={{ formatter: (value) => `${value}%` }}
              />
              <Text type="secondary" style={{ fontSize: 12 }}>
                Update your progress as you work towards your goal
              </Text>
            </div>
          </Form.Item>
        )}

        <Divider />

        <div style={{ background: '#f8f9fa', padding: 16, borderRadius: 8 }}>
          <Text strong style={{ color: '#1677ff' }}>💡 SMART Goal Tips:</Text>
          <ul style={{ marginTop: 8, marginBottom: 0, paddingLeft: 20 }}>
            {smartGoalTips.map((tip, index) => (
              <li key={index}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {tip}
                </Text>
              </li>
            ))}
          </ul>
        </div>
      </Form>
    </Modal>
  );
};

export default GoalModal;