import React, { useState } from 'react';
import {
  Tree,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  DatePicker,
  message,
  Typography,
  Progress,
  Tag,
  Dropdown,
  Card
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  MoreOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  PlayCircleOutlined,
  FlagOutlined,
  BranchesOutlined
} from '@ant-design/icons';
import { Task } from '../../types/api';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { Text } = Typography;

interface TaskTreeViewProps {
  tasks: Task[];
  onTaskCreate: (task: Omit<Task, 'id' | 'created_at' | 'updated_at'>) => Promise<void>;
  onTaskUpdate: (taskId: string, updates: Partial<Task>) => Promise<void>;
  onTaskDelete: (taskId: string) => Promise<void>;
  goalId: string;
}


const TaskTreeView: React.FC<TaskTreeViewProps> = ({
  tasks,
  onTaskCreate,
  onTaskUpdate,
  onTaskDelete,
  goalId
}) => {
  const [taskModalVisible, setTaskModalVisible] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [parentTaskId, setParentTaskId] = useState<string | null>(null);
  const [form] = Form.useForm();

  // Build tree structure from flat task list
  const buildTaskTree = (tasks: Task[]): any[] => {
    const taskMap = new Map<string, any>();
    const rootTasks: any[] = [];

    // Create task map
    tasks.forEach(task => {
      taskMap.set(task.id, { ...task, children: [] });
    });

    // Build tree structure
    tasks.forEach(task => {
      const taskNode = taskMap.get(task.id)!;
      const treeNode = {
        key: task.id,
        title: renderTaskNode(taskNode),
        children: [],
        task: taskNode
      };

      if (task.parent_task_id && taskMap.has(task.parent_task_id)) {
        const parent = taskMap.get(task.parent_task_id)!;
        if (!parent.children) parent.children = [];
        parent.children.push(treeNode);
      } else {
        rootTasks.push(treeNode);
      }
    });

    return rootTasks;
  };

  const renderTaskNode = (task: Task) => {
    const getStatusIcon = (status: Task['status']) => {
      switch (status) {
        case 'completed':
          return <CheckCircleOutlined style={{ color: '#52c41a' }} />;
        case 'in_progress':
          return <PlayCircleOutlined style={{ color: '#1677ff' }} />;
        case 'pending':
          return <ClockCircleOutlined style={{ color: '#faad14' }} />;
        default:
          return <ClockCircleOutlined style={{ color: '#d9d9d9' }} />;
      }
    };

    const getPriorityColor = (priority: Task['priority']) => {
      switch (priority) {
        case 'critical': return '#ff4d4f';
        case 'high': return '#fa8c16';
        case 'medium': return '#1677ff';
        case 'low': return '#52c41a';
        default: return '#d9d9d9';
      }
    };

    const getProgressFromSubtasks = (task: Task): number => {
      const subtasks = tasks.filter(t => t.parent_task_id === task.id);
      if (subtasks.length === 0) {
        return task.status === 'completed' ? 100 : 0;
      }
      const completedSubtasks = subtasks.filter(t => t.status === 'completed').length;
      return Math.round((completedSubtasks / subtasks.length) * 100);
    };

    const progress = getProgressFromSubtasks(task);
    const subtaskCount = tasks.filter(t => t.parent_task_id === task.id).length;

    return (
      <div style={{ 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'space-between',
        width: '100%',
        minHeight: '32px'
      }}>
        <Space size="small">
          {getStatusIcon(task.status)}
          <Text 
            style={{ 
              textDecoration: task.status === 'completed' ? 'line-through' : 'none',
              opacity: task.status === 'completed' ? 0.6 : 1
            }}
          >
            {task.title}
          </Text>
          
          <Tag color={getPriorityColor(task.priority)}>
            {task.priority}
          </Tag>

          {subtaskCount > 0 && (
            <Tag icon={<BranchesOutlined />} color="blue">
              {subtaskCount} subtasks
            </Tag>
          )}

          {task.deadline && (
            <Tag color={dayjs(task.deadline).isBefore(dayjs()) ? 'red' : 'green'}>
              {dayjs(task.deadline).format('MMM DD')}
            </Tag>
          )}
        </Space>

        <Space size="small">
          {progress > 0 && (
            <Progress 
              percent={progress} 
              size="small" 
              style={{ width: '60px' }}
              showInfo={false}
            />
          )}
          
          <Dropdown
            menu={{
              items: [
                {
                  key: 'add-subtask',
                  label: 'Add Subtask',
                  icon: <PlusOutlined />,
                  onClick: () => handleAddSubtask(task.id)
                },
                {
                  key: 'edit',
                  label: 'Edit Task',
                  icon: <EditOutlined />,
                  onClick: () => handleEditTask(task)
                },
                {
                  key: 'toggle-status',
                  label: task.status === 'completed' ? 'Mark Incomplete' : 'Mark Complete',
                  icon: <CheckCircleOutlined />,
                  onClick: () => handleToggleStatus(task)
                },
                { type: 'divider' },
                {
                  key: 'delete',
                  label: 'Delete Task',
                  icon: <DeleteOutlined />,
                  danger: true,
                  onClick: () => handleDeleteTask(task.id)
                }
              ]
            }}
            trigger={['click']}
          >
            <Button type="text" size="small" icon={<MoreOutlined />} />
          </Dropdown>
        </Space>
      </div>
    );
  };

  const handleAddTask = () => {
    setEditingTask(null);
    setParentTaskId(null);
    setTaskModalVisible(true);
  };

  const handleAddSubtask = (parentId: string) => {
    setEditingTask(null);
    setParentTaskId(parentId);
    setTaskModalVisible(true);
  };

  const handleEditTask = (task: Task) => {
    setEditingTask(task);
    setParentTaskId(task.parent_task_id || null);
    form.setFieldsValue({
      title: task.title,
      description: task.description,
      priority: task.priority,
      status: task.status,
      deadline: task.deadline ? dayjs(task.deadline) : null,
      estimated_duration: task.estimated_duration
    });
    setTaskModalVisible(true);
  };

  const handleToggleStatus = async (task: Task) => {
    const newStatus = task.status === 'completed' ? 'pending' : 'completed';
    try {
      await onTaskUpdate(task.id, {
        status: newStatus,
        completed_at: newStatus === 'completed' ? new Date().toISOString() : undefined
      });
      message.success(`Task marked as ${newStatus}`);
    } catch {
      message.error('Failed to update task status');
    }
  };

  const handleDeleteTask = async (taskId: string) => {
    Modal.confirm({
      title: 'Delete Task',
      content: 'Are you sure you want to delete this task and all its subtasks?',
      okText: 'Delete',
      okType: 'danger',
      onOk: async () => {
        try {
          await onTaskDelete(taskId);
          message.success('Task deleted successfully');
        } catch {
          message.error('Failed to delete task');
        }
      }
    });
  };

  const handleTaskSubmit = async () => {
    try {
      const values = await form.validateFields();
      const taskData: Omit<Task, 'id' | 'created_at' | 'updated_at'> = {
        goal_id: goalId,
        parent_task_id: parentTaskId || undefined,
        title: values.title,
        description: values.description || '',
        priority: values.priority,
        status: values.status,
        deadline: values.deadline?.toISOString(),
        estimated_duration: values.estimated_duration,
        order_index: tasks.length + 1
      };

      if (editingTask) {
        await onTaskUpdate(editingTask.id, taskData);
        message.success('Task updated successfully');
      } else {
        await onTaskCreate(taskData);
        message.success('Task created successfully');
      }

      setTaskModalVisible(false);
      form.resetFields();
    } catch {
      message.error('Failed to save task');
    }
  };

  const treeData = buildTaskTree(tasks);

  return (
    <Card 
      title={
        <Space>
          <BranchesOutlined />
          Task Breakdown
        </Space>
      }
      extra={
        <Button 
          type="primary" 
          icon={<PlusOutlined />} 
          size="small"
          onClick={handleAddTask}
        >
          Add Task
        </Button>
      }
    >
      {treeData.length > 0 ? (
        <Tree
          treeData={treeData}
          defaultExpandAll
          showLine={{ showLeafIcon: false }}
          selectable={false}
        />
      ) : (
        <div style={{ textAlign: 'center', padding: '20px' }}>
          <Text type="secondary">No tasks yet. Add your first task to break down this goal.</Text>
        </div>
      )}

      {/* Task Modal */}
      <Modal
        title={
          <Space>
            <FlagOutlined />
            {editingTask ? 'Edit Task' : parentTaskId ? 'Add Subtask' : 'Add Task'}
          </Space>
        }
        open={taskModalVisible}
        onCancel={() => {
          setTaskModalVisible(false);
          form.resetFields();
        }}
        onOk={handleTaskSubmit}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            priority: 'medium',
            status: 'pending'
          }}
        >
          <Form.Item
            name="title"
            label="Task Title"
            rules={[
              { required: true, message: 'Please enter task title' },
              { max: 100, message: 'Title should be less than 100 characters' }
            ]}
          >
            <Input placeholder="What needs to be done?" maxLength={100} />
          </Form.Item>

          <Form.Item
            name="description"
            label="Description"
          >
            <TextArea
              rows={3}
              placeholder="Describe the task in detail..."
              maxLength={500}
              showCount
            />
          </Form.Item>

          <div style={{ display: 'flex', gap: '16px' }}>
            <Form.Item
              name="priority"
              label="Priority"
              style={{ flex: 1 }}
              rules={[{ required: true, message: 'Please select priority' }]}
            >
              <Select>
                <Select.Option value="low">Low</Select.Option>
                <Select.Option value="medium">Medium</Select.Option>
                <Select.Option value="high">High</Select.Option>
                <Select.Option value="critical">Critical</Select.Option>
              </Select>
            </Form.Item>

            <Form.Item
              name="status"
              label="Status"
              style={{ flex: 1 }}
              rules={[{ required: true, message: 'Please select status' }]}
            >
              <Select>
                <Select.Option value="pending">Pending</Select.Option>
                <Select.Option value="in_progress">In Progress</Select.Option>
                <Select.Option value="completed">Completed</Select.Option>
                <Select.Option value="cancelled">Cancelled</Select.Option>
              </Select>
            </Form.Item>
          </div>

          <div style={{ display: 'flex', gap: '16px' }}>
            <Form.Item
              name="deadline"
              label="Deadline"
              style={{ flex: 1 }}
            >
              <DatePicker 
                style={{ width: '100%' }}
                showTime={{ format: 'HH:mm' }}
                format="YYYY-MM-DD HH:mm"
              />
            </Form.Item>

            <Form.Item
              name="estimated_duration"
              label="Est. Duration (hours)"
              style={{ flex: 1 }}
            >
              <Input type="number" min={0} max={999} placeholder="Hours" />
            </Form.Item>
          </div>
        </Form>
      </Modal>
    </Card>
  );
};

export default TaskTreeView;