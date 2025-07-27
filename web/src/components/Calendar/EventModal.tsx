import React, { useEffect, useState } from 'react';
import {
  Modal,
  Form,
  Input,
  DatePicker,
  Select,
  Button,
  Space,
  message,
  Popconfirm,
  Divider,
  Typography
} from 'antd';
import { DeleteOutlined, ReloadOutlined } from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../hooks/redux';
import { createEvent, updateEvent, deleteEvent } from '../../store/slices/eventsSlice';
import { fetchGoals } from '../../store/slices/goalsSlice';
import { Event, Recurrence } from '../../types/api';
import RecurrenceModal from './RecurrenceModal';
import { getUserTimezone, getCommonTimezones } from '../../utils/timezone';
import dayjs from 'dayjs';

const { TextArea } = Input;
const { RangePicker } = DatePicker;
const { Text } = Typography;

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

  const [recurrenceModalVisible, setRecurrenceModalVisible] = useState(false);
  const [currentRecurrence, setCurrentRecurrence] = useState<Recurrence | null>(null);
  const [selectedTimezone, setSelectedTimezone] = useState<string>(getUserTimezone());

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
        setCurrentRecurrence(event.recurrence || null);
        setSelectedTimezone(event.timezone || getUserTimezone());
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
        setCurrentRecurrence(null);
        setSelectedTimezone(getUserTimezone());
      }
    } else {
      form.resetFields();
      setCurrentRecurrence(null);
      setSelectedTimezone(getUserTimezone());
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
        timezone: selectedTimezone,
        recurrence: currentRecurrence || undefined,
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

  const handleRecurrenceChange = (recurrence: Recurrence | null) => {
    setCurrentRecurrence(recurrence);
    setRecurrenceModalVisible(false);
  };

  const formatRecurrenceText = (recurrence: Recurrence | null): string => {
    if (!recurrence) return 'Does not repeat';
    
    const { freq, interval = 1, count, until, byweekday, bymonthday, bymonth } = recurrence;
    
    let text = '';
    
    // Frequency and interval
    if (interval === 1) {
      text = freq.toLowerCase();
    } else {
      text = `every ${interval} ${freq.toLowerCase().slice(0, -2)}s`;
    }
    
    // Additional specifics
    if (freq === 'WEEKLY' && byweekday?.length) {
      const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
      const selectedDays = byweekday.map(day => days[day]).join(', ');
      text += ` on ${selectedDays}`;
    } else if (freq === 'MONTHLY' && bymonthday?.length) {
      text += ` on day ${bymonthday.join(', ')}`;
    } else if (freq === 'YEARLY' && bymonth?.length) {
      const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
      const selectedMonths = bymonth.map(month => months[month - 1]).join(', ');
      text += ` in ${selectedMonths}`;
    }
    
    // End condition
    if (count) {
      text += `, ${count} times`;
    } else if (until) {
      text += `, until ${dayjs(until).format('MMM D, YYYY')}`;
    }
    
    return text.charAt(0).toUpperCase() + text.slice(1);
  };

  const getCurrentStartDate = () => {
    const dateTime = form.getFieldValue('dateTime');
    return dateTime?.[0] || dayjs();
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
          label="Timezone"
          help="Select the timezone for this event"
        >
          <Select
            value={selectedTimezone}
            onChange={setSelectedTimezone}
            placeholder="Select timezone"
            showSearch
            optionFilterProp="children"
            filterOption={(input, option) =>
              String(option?.children || '').toLowerCase().includes(input.toLowerCase())
            }
          >
            {getCommonTimezones().map(tz => (
              <Select.Option key={tz.name} value={tz.name}>
                {tz.label}
              </Select.Option>
            ))}
          </Select>
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

        <Divider />

        <Form.Item label="Repeat">
          <Space direction="vertical" style={{ width: '100%' }}>
            <div style={{ 
              border: '1px solid #d9d9d9', 
              borderRadius: 6, 
              padding: 12,
              backgroundColor: '#fafafa',
              cursor: 'pointer',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}
            onClick={() => setRecurrenceModalVisible(true)}
            >
              <Text>{formatRecurrenceText(currentRecurrence)}</Text>
              <Button 
                type="text" 
                icon={<ReloadOutlined />}
                size="small"
              >
                Edit
              </Button>
            </div>
          </Space>
        </Form.Item>
      </Form>

      {/* Recurrence Modal */}
      <RecurrenceModal
        visible={recurrenceModalVisible}
        recurrence={currentRecurrence}
        startDate={getCurrentStartDate()}
        onOk={handleRecurrenceChange}
        onCancel={() => setRecurrenceModalVisible(false)}
      />
    </Modal>
  );
};

export default EventModal;