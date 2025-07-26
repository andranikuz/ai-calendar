import React, { useEffect } from 'react';
import {
  Modal,
  Form,
  Input,
  DatePicker,
  Select,
  Button,
  Space,
  message,
  Popconfirm
} from 'antd';
import { DeleteOutlined } from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { createEvent, updateEvent, deleteEvent } from '../../store/slices/eventsSlice';
import { fetchGoals } from '../../store/slices/goalsSlice';
import { Event } from '../../types/api';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { RangePicker } = DatePicker;

interface EventModalProps {
  visible: boolean;
  event: Event | null;
  selectedDate: string | null;
  onCancel: () => void;
  onSuccess: () => void;
}

const EventModal: React.FC<EventModalProps> = ({
  visible,
  event,
  selectedDate,
  onCancel,
  onSuccess
}) => {
  const [form] = Form.useForm();
  const dispatch = useAppDispatch();
  
  const { goals } = useAppSelector(state => state.goals);
  const { isLoading } = useAppSelector(state => state.events);

  const isEditing = !!event;

  useEffect(() => {
    // Fetch goals for linking
    dispatch(fetchGoals());
  }, [dispatch]);

  useEffect(() => {
    if (visible) {
      if (event) {
        // Editing existing event
        form.setFieldsValue({
          title: event.title,
          description: event.description,
          location: event.location,
          dateTime: [
            dayjs(event.start_time),
            dayjs(event.end_time)
          ],
          goal_id: event.goal_id || undefined,
        });
      } else {
        // Creating new event
        const startDate = selectedDate ? dayjs(selectedDate) : dayjs();
        const endDate = startDate.add(1, 'hour');
        
        form.setFieldsValue({
          title: '',
          description: '',
          location: '',
          dateTime: [startDate, endDate],
          goal_id: undefined,
        });
      }
    } else {
      form.resetFields();
    }
  }, [visible, event, selectedDate, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const [startTime, endTime] = values.dateTime;

      const eventData = {
        title: values.title,
        description: values.description || '',
        location: values.location || '',
        start_time: startTime.toISOString(),
        end_time: endTime.toISOString(),
        goal_id: values.goal_id || undefined,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      };

      if (isEditing && event) {
        await dispatch(updateEvent({
          id: event.id,
          data: eventData
        })).unwrap();
        message.success('Event updated successfully');
      } else {
        await dispatch(createEvent(eventData)).unwrap();
        message.success('Event created successfully');
      }

      onSuccess();
    } catch (error) {
      console.error('Failed to save event:', error);
      message.error('Failed to save event');
    }
  };

  const handleDelete = async () => {
    if (!event) return;

    try {
      await dispatch(deleteEvent(event.id)).unwrap();
      message.success('Event deleted successfully');
      onSuccess();
    } catch (error) {
      console.error('Failed to delete event:', error);
      message.error('Failed to delete event');
    }
  };

  const validateTimeRange = (_rule: any, value: any) => {
    if (!value || !Array.isArray(value) || value.length !== 2) {
      return Promise.reject('Please select start and end time');
    }
    
    const [start, end] = value;
    if (!start || !end) {
      return Promise.reject('Please select start and end time');
    }
    
    if (end.isBefore(start)) {
      return Promise.reject('End time must be after start time');
    }
    
    return Promise.resolve();
  };

  return (
    <Modal
      title={isEditing ? 'Edit Event' : 'Create New Event'}
      open={visible}
      onCancel={onCancel}
      footer={[
        <Space key="actions" style={{ width: '100%', justifyContent: 'space-between' }}>
          <div>
            {isEditing && (
              <Popconfirm
                title="Are you sure you want to delete this event?"
                onConfirm={handleDelete}
                okText="Yes"
                cancelText="No"
              >
                <Button danger icon={<DeleteOutlined />}>
                  Delete
                </Button>
              </Popconfirm>
            )}
          </div>
          <Space>
            <Button onClick={onCancel}>
              Cancel
            </Button>
            <Button
              type="primary"
              loading={isLoading}
              onClick={handleSubmit}
            >
              {isEditing ? 'Update' : 'Create'}
            </Button>
          </Space>
        </Space>
      ]}
      width={600}
    >
      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        <Form.Item
          name="title"
          label="Event Title"
          rules={[
            { required: true, message: 'Please enter event title' },
            { max: 255, message: 'Title must be less than 255 characters' }
          ]}
        >
          <Input placeholder="Enter event title" />
        </Form.Item>

        <Form.Item
          name="dateTime"
          label="Date & Time"
          rules={[
            { required: true, message: 'Please select date and time' },
            { validator: validateTimeRange }
          ]}
        >
          <RangePicker
            showTime={{ format: 'HH:mm' }}
            format="YYYY-MM-DD HH:mm"
            style={{ width: '100%' }}
            placeholder={['Start time', 'End time']}
          />
        </Form.Item>

        <Form.Item
          name="description"
          label="Description"
        >
          <TextArea
            rows={3}
            placeholder="Enter event description (optional)"
            maxLength={1000}
            showCount
          />
        </Form.Item>

        <Form.Item
          name="location"
          label="Location"
        >
          <Input placeholder="Enter location (optional)" />
        </Form.Item>

        <Form.Item
          name="goal_id"
          label="Link to Goal"
          help="Link this event to a goal to track progress"
        >
          <Select
            placeholder="Select a goal (optional)"
            allowClear
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) =>
              String(option?.children || '').toLowerCase().includes(input.toLowerCase())
            }
          >
            {goals.map(goal => (
              <Select.Option key={goal.id} value={goal.id}>
                {goal.title} ({goal.category})
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default EventModal;