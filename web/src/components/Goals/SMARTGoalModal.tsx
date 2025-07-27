import React, { useEffect, useState } from 'react';
import {
  Modal,
  Form,
  Input,
  Select,
  DatePicker,
  Button,
  Space,
  message,
  Typography,
  Divider,
  Alert,
  Progress,
  Card,
  Badge,
  Tooltip,
  Row,
  Col
} from '../../utils/antd';
import { 
  FlagOutlined, 
  CalendarOutlined, 
  AimOutlined,
  BulbOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  InfoCircleOutlined
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { createGoal, updateGoal } from '../../store/slices/goalsSlice';
import { Goal } from '../../types/api';
import { 
  validateSMARTGoal, 
  getSMARTSuggestions,
  SMARTValidationResult 
} from '../../utils/smartGoals';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { Text, Title } = Typography;

interface SMARTGoalModalProps {
  visible: boolean;
  goal: Goal | null;
  onCancel: () => void;
  onSuccess: () => void;
}

const SMARTGoalModal: React.FC<SMARTGoalModalProps> = ({
  visible,
  goal,
  onCancel,
  onSuccess
}) => {
  const [form] = Form.useForm();
  const dispatch = useAppDispatch();
  const { isLoading } = useAppSelector(state => state.goals);

  const [smartValidation, setSmartValidation] = useState<SMARTValidationResult | null>(null);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<string>('personal');

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
        });
        setSelectedCategory(goal.category);
      } else {
        // Creating new goal
        form.setFieldsValue({
          title: '',
          description: '',
          category: 'personal',
          priority: 'medium',
          deadline: null,
        });
        setSelectedCategory('personal');
      }
    } else {
      form.resetFields();
      setSmartValidation(null);
      setShowSuggestions(false);
    }
  }, [visible, goal, form]);

  const handleFormChange = () => {
    const values = form.getFieldsValue();
    if (values.title || values.description) {
      const validation = validateSMARTGoal({
        title: values.title || '',
        description: values.description || '',
        category: values.category || 'personal',
        deadline: values.deadline?.toISOString(),
        priority: values.priority || 'medium'
      });
      setSmartValidation(validation);
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      
      // Validate SMART criteria
      const validation = validateSMARTGoal({
        title: values.title,
        description: values.description,
        category: values.category,
        deadline: values.deadline?.toISOString(),
        priority: values.priority
      });

      if (validation.score < 60) {
        message.warning('Consider improving your goal to meet more SMART criteria for better success');
      }

      const goalData = {
        title: values.title,
        description: values.description,
        category: values.category,
        priority: values.priority,
        deadline: values.deadline ? values.deadline.toISOString() : undefined,
      };

      if (isEditing && goal) {
        await dispatch(updateGoal({
          id: goal.id,
          data: goalData
        })).unwrap();
        message.success('SMART goal updated successfully');
      } else {
        await dispatch(createGoal(goalData)).unwrap();
        message.success('SMART goal created successfully');
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

  const getCriteriaIcon = (met: boolean) => {
    return met ? 
      <CheckCircleOutlined style={{ color: '#52c41a' }} /> : 
      <CloseCircleOutlined style={{ color: '#ff4d4f' }} />;
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return '#52c41a';
    if (score >= 60) return '#fa8c16';
    return '#ff4d4f';
  };

  const suggestions = getSMARTSuggestions(selectedCategory);

  return (
    <Modal
      title={
        <Space>
          <FlagOutlined />
          {isEditing ? 'Edit SMART Goal' : 'Create SMART Goal'}
        </Space>
      }
      open={visible}
      onCancel={onCancel}
      width={800}
      footer={[
        <Button key="suggestions" onClick={() => setShowSuggestions(!showSuggestions)}>
          <BulbOutlined /> {showSuggestions ? 'Hide' : 'Show'} Examples
        </Button>,
        <Button key="cancel" onClick={onCancel}>
          Cancel
        </Button>,
        <Button
          key="submit"
          type="primary"
          loading={isLoading}
          onClick={handleSubmit}
          disabled={smartValidation ? smartValidation.score < 40 : false}
        >
          {isEditing ? 'Update Goal' : 'Create Goal'}
        </Button>
      ]}
    >
      <Row gutter={[16, 16]}>
        <Col span={16}>
          <Form
            form={form}
            layout="vertical"
            onValuesChange={handleFormChange}
            preserve={false}
          >
            <Form.Item
              name="title"
              label="Goal Title"
              rules={[
                { required: true, message: 'Please enter goal title' },
                { min: 5, message: 'Title should be at least 5 characters' },
                { max: 100, message: 'Title should be less than 100 characters' }
              ]}
            >
              <Input 
                placeholder="What specific outcome do you want to achieve?"
                showCount
                maxLength={100}
              />
            </Form.Item>

            <Form.Item
              name="description"
              label="Detailed Description"
              rules={[
                { required: true, message: 'Please enter goal description' },
                { min: 20, message: 'Description should be at least 20 characters for clarity' },
                { max: 500, message: 'Description should be less than 500 characters' }
              ]}
            >
              <TextArea
                rows={4}
                placeholder="Describe your goal in detail. Include what, when, where, how, and why..."
                showCount
                maxLength={500}
              />
            </Form.Item>

            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="category"
                  label="Category"
                  rules={[{ required: true, message: 'Please select a category' }]}
                >
                  <Select
                    placeholder="Select goal category"
                    onChange={setSelectedCategory}
                  >
                    {categoryOptions.map(option => (
                      <Select.Option key={option.value} value={option.value}>
                        <div>
                          <div>{option.label}</div>
                          <Text type="secondary" style={{ fontSize: '12px' }}>
                            {option.description}
                          </Text>
                        </div>
                      </Select.Option>
                    ))}
                  </Select>
                </Form.Item>
              </Col>

              <Col span={12}>
                <Form.Item
                  name="priority"
                  label="Priority"
                  rules={[{ required: true, message: 'Please select priority' }]}
                >
                  <Select placeholder="Select priority level">
                    {priorityOptions.map(option => (
                      <Select.Option key={option.value} value={option.value}>
                        <Badge color={option.color} text={option.label} />
                      </Select.Option>
                    ))}
                  </Select>
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="deadline"
              label={
                <Space>
                  <CalendarOutlined />
                  Deadline (Time-bound)
                </Space>
              }
              rules={[
                { required: true, message: 'Setting a deadline is crucial for SMART goals' },
                {
                  validator: (_, value) => {
                    if (!value) return Promise.reject('Please set a deadline');
                    if (value.isBefore(dayjs())) {
                      return Promise.reject('Deadline must be in the future');
                    }
                    if (value.isAfter(dayjs().add(3, 'years'))) {
                      return Promise.reject('Deadline should be within 3 years for better focus');
                    }
                    return Promise.resolve();
                  }
                }
              ]}
            >
              <DatePicker
                style={{ width: '100%' }}
                disabledDate={(current) => current && current < dayjs().startOf('day')}
                showTime={{ format: 'HH:mm' }}
                format="YYYY-MM-DD HH:mm"
                placeholder="When do you want to achieve this goal?"
              />
            </Form.Item>
          </Form>
        </Col>

        <Col span={8}>
          {/* SMART Validation Panel */}
          {smartValidation && (
            <Card size="small" title={
              <Space>
                <AimOutlined />
                SMART Score
              </Space>
            }>
              <div style={{ textAlign: 'center', marginBottom: 16 }}>
                <Progress
                  type="circle"
                  size={80}
                  percent={smartValidation.score}
                  strokeColor={getScoreColor(smartValidation.score)}
                  format={(percent) => `${percent}%`}
                />
                <div style={{ marginTop: 8 }}>
                  <Text strong style={{ color: getScoreColor(smartValidation.score) }}>
                    {smartValidation.score >= 80 ? 'Excellent!' : 
                     smartValidation.score >= 60 ? 'Good' : 'Needs Work'}
                  </Text>
                </div>
              </div>

              <Divider style={{ margin: '12px 0' }} />

              <div>
                <Title level={5} style={{ margin: '8px 0' }}>SMART Criteria:</Title>
                
                <div style={{ marginBottom: 8 }}>
                  <Space>
                    {getCriteriaIcon(smartValidation.criteria.specific)}
                    <Tooltip title="Clear and well-defined outcome">
                      <Text strong>Specific</Text>
                    </Tooltip>
                  </Space>
                </div>

                <div style={{ marginBottom: 8 }}>
                  <Space>
                    {getCriteriaIcon(smartValidation.criteria.measurable)}
                    <Tooltip title="Includes metrics and measurable outcomes">
                      <Text strong>Measurable</Text>
                    </Tooltip>
                  </Space>
                </div>

                <div style={{ marginBottom: 8 }}>
                  <Space>
                    {getCriteriaIcon(smartValidation.criteria.achievable)}
                    <Tooltip title="Realistic and attainable">
                      <Text strong>Achievable</Text>
                    </Tooltip>
                  </Space>
                </div>

                <div style={{ marginBottom: 8 }}>
                  <Space>
                    {getCriteriaIcon(smartValidation.criteria.relevant)}
                    <Tooltip title="Aligns with values and objectives">
                      <Text strong>Relevant</Text>
                    </Tooltip>
                  </Space>
                </div>

                <div style={{ marginBottom: 8 }}>
                  <Space>
                    {getCriteriaIcon(smartValidation.criteria.timeBound)}
                    <Tooltip title="Has a clear deadline">
                      <Text strong>Time-bound</Text>
                    </Tooltip>
                  </Space>
                </div>
              </div>

              {smartValidation.suggestions.length > 0 && (
                <>
                  <Divider style={{ margin: '12px 0' }} />
                  <div>
                    <Title level={5} style={{ margin: '8px 0' }}>Suggestions:</Title>
                    {smartValidation.suggestions.map((suggestion, index) => (
                      <div key={index} style={{ marginBottom: 4, fontSize: '12px' }}>
                        <InfoCircleOutlined style={{ color: '#1677ff', marginRight: 4 }} />
                        <Text type="secondary">{suggestion}</Text>
                      </div>
                    ))}
                  </div>
                </>
              )}

              {smartValidation.warnings.length > 0 && (
                <>
                  <Divider style={{ margin: '12px 0' }} />
                  <Alert
                    message="Warnings"
                    description={smartValidation.warnings.join('. ')}
                    type="warning"
                    showIcon
                  />
                </>
              )}
            </Card>
          )}
        </Col>
      </Row>

      {/* Examples Section */}
      {showSuggestions && (
        <>
          <Divider />
          <Card size="small" title={
            <Space>
              <BulbOutlined />
              SMART Goal Examples for {categoryOptions.find(c => c.value === selectedCategory)?.label}
            </Space>
          }>
            {suggestions.map((suggestion, index) => (
              <Alert
                key={index}
                message={suggestion}
                type="info"
                style={{ marginBottom: 8 }}
                action={
                  <Button 
                    size="small" 
                    type="text"
                    onClick={() => {
                      // Parse suggestion and populate form
                      const parts = suggestion.split(' by ');
                      if (parts.length >= 2) {
                        form.setFieldsValue({
                          title: parts[0],
                          description: suggestion
                        });
                        handleFormChange();
                      }
                    }}
                  >
                    Use as Template
                  </Button>
                }
              />
            ))}
          </Card>
        </>
      )}
    </Modal>
  );
};

export default SMARTGoalModal;