import React, { useState, useEffect } from 'react';
import {
  Modal,
  Form,
  Select,
  DatePicker,
  Checkbox,
  Space,
  InputNumber,
  Row,
  Col,
  Button,
  Typography,
  Alert
} from 'antd';
import { Recurrence } from '../../types/api';
import dayjs from 'dayjs';

const { Text } = Typography;

interface RecurrenceModalProps {
  visible: boolean;
  recurrence: Recurrence | null;
  startDate: dayjs.Dayjs;
  onOk: (recurrence: Recurrence | null) => void;
  onCancel: () => void;
}

const RecurrenceModal: React.FC<RecurrenceModalProps> = ({
  visible,
  recurrence,
  startDate,
  onOk,
  onCancel
}) => {
  const [form] = Form.useForm();
  const [frequency, setFrequency] = useState<'DAILY' | 'WEEKLY' | 'MONTHLY' | 'YEARLY'>('WEEKLY');
  const [endType, setEndType] = useState<'never' | 'after' | 'on'>('never');

  useEffect(() => {
    if (visible) {
      if (recurrence) {
        form.setFieldsValue({
          freq: recurrence.freq,
          interval: recurrence.interval || 1,
          count: recurrence.count,
          until: recurrence.until ? dayjs(recurrence.until) : null,
          byweekday: recurrence.byweekday || [],
          bymonthday: recurrence.bymonthday || [],
          bymonth: recurrence.bymonth || []
        });
        setFrequency(recurrence.freq);
        if (recurrence.count) {
          setEndType('after');
        } else if (recurrence.until) {
          setEndType('on');
        } else {
          setEndType('never');
        }
      } else {
        form.setFieldsValue({
          freq: 'WEEKLY',
          interval: 1,
          byweekday: [startDate.day()],
        });
        setFrequency('WEEKLY');
        setEndType('never');
      }
    }
  }, [visible, recurrence, startDate, form]);

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      
      const newRecurrence: Recurrence = {
        freq: values.freq,
        interval: values.interval || 1,
      };

      // Add end condition
      if (endType === 'after' && values.count) {
        newRecurrence.count = values.count;
      } else if (endType === 'on' && values.until) {
        newRecurrence.until = values.until.toISOString();
      }

      // Add frequency-specific fields
      if (frequency === 'WEEKLY' && values.byweekday?.length > 0) {
        newRecurrence.byweekday = values.byweekday;
      } else if (frequency === 'MONTHLY' && values.bymonthday?.length > 0) {
        newRecurrence.bymonthday = values.bymonthday;
      } else if (frequency === 'YEARLY' && values.bymonth?.length > 0) {
        newRecurrence.bymonth = values.bymonth;
      }

      onOk(newRecurrence);
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  const handleRemoveRecurrence = () => {
    onOk(null);
  };

  const weekdayOptions = [
    { label: 'Sun', value: 0 },
    { label: 'Mon', value: 1 },
    { label: 'Tue', value: 2 },
    { label: 'Wed', value: 3 },
    { label: 'Thu', value: 4 },
    { label: 'Fri', value: 5 },
    { label: 'Sat', value: 6 },
  ];

  const monthOptions = [
    { label: 'Jan', value: 1 },
    { label: 'Feb', value: 2 },
    { label: 'Mar', value: 3 },
    { label: 'Apr', value: 4 },
    { label: 'May', value: 5 },
    { label: 'Jun', value: 6 },
    { label: 'Jul', value: 7 },
    { label: 'Aug', value: 8 },
    { label: 'Sep', value: 9 },
    { label: 'Oct', value: 10 },
    { label: 'Nov', value: 11 },
    { label: 'Dec', value: 12 },
  ];

  const renderFrequencySpecificFields = () => {
    switch (frequency) {
      case 'WEEKLY':
        return (
          <Form.Item
            name="byweekday"
            label="Repeat on"
            rules={[
              { required: true, message: 'Please select at least one day' },
              { type: 'array', min: 1, message: 'Please select at least one day' }
            ]}
          >
            <Checkbox.Group options={weekdayOptions} />
          </Form.Item>
        );
      
      case 'MONTHLY':
        return (
          <Form.Item
            name="bymonthday"
            label="On day(s) of month"
            help="Select specific days of the month (1-31)"
          >
            <Select
              mode="multiple"
              placeholder="Select days"
              style={{ width: '100%' }}
            >
              {Array.from({ length: 31 }, (_, i) => (
                <Select.Option key={i + 1} value={i + 1}>
                  {i + 1}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        );
      
      case 'YEARLY':
        return (
          <Form.Item
            name="bymonth"
            label="In month(s)"
            help="Select specific months"
          >
            <Checkbox.Group options={monthOptions} />
          </Form.Item>
        );
      
      default:
        return null;
    }
  };

  return (
    <Modal
      title="Repeat Event"
      open={visible}
      onOk={handleOk}
      onCancel={onCancel}
      width={600}
      footer={[
        <Space key="actions" style={{ width: '100%', justifyContent: 'space-between' }}>
          <div>
            {recurrence && (
              <Button type="default" onClick={handleRemoveRecurrence}>
                Remove Recurrence
              </Button>
            )}
          </div>
          <Space>
            <Button onClick={onCancel}>Cancel</Button>
            <Button type="primary" onClick={handleOk}>
              Save
            </Button>
          </Space>
        </Space>
      ]}
    >
      <Alert
        message="Create Recurring Event"
        description="Configure how often this event should repeat. Changes will apply to this and future occurrences."
        type="info"
        showIcon
        style={{ marginBottom: 16 }}
      />

      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        <Row gutter={16}>
          <Col span={12}>
            <Form.Item
              name="freq"
              label="Frequency"
              rules={[{ required: true, message: 'Please select frequency' }]}
            >
              <Select
                value={frequency}
                onChange={(value) => setFrequency(value)}
                placeholder="Select frequency"
              >
                <Select.Option value="DAILY">Daily</Select.Option>
                <Select.Option value="WEEKLY">Weekly</Select.Option>
                <Select.Option value="MONTHLY">Monthly</Select.Option>
                <Select.Option value="YEARLY">Yearly</Select.Option>
              </Select>
            </Form.Item>
          </Col>
          
          <Col span={12}>
            <Form.Item
              name="interval"
              label="Every"
              help={`Repeat every ${frequency === 'DAILY' ? 'day(s)' : 
                     frequency === 'WEEKLY' ? 'week(s)' : 
                     frequency === 'MONTHLY' ? 'month(s)' : 'year(s)'}`}
            >
              <InputNumber
                min={1}
                max={99}
                placeholder="1"
                style={{ width: '100%' }}
              />
            </Form.Item>
          </Col>
        </Row>

        {renderFrequencySpecificFields()}

        <Form.Item label="Ends">
          <Space direction="vertical" style={{ width: '100%' }}>
            <Select
              value={endType}
              onChange={setEndType}
              style={{ width: '100%' }}
            >
              <Select.Option value="never">Never</Select.Option>
              <Select.Option value="after">After a number of occurrences</Select.Option>
              <Select.Option value="on">On a specific date</Select.Option>
            </Select>

            {endType === 'after' && (
              <Form.Item
                name="count"
                rules={[
                  { required: true, message: 'Please enter number of occurrences' },
                  { type: 'number', min: 1, max: 999, message: 'Must be between 1 and 999' }
                ]}
              >
                <Space>
                  <Text>After</Text>
                  <InputNumber
                    min={1}
                    max={999}
                    placeholder="10"
                    style={{ width: 100 }}
                  />
                  <Text>occurrences</Text>
                </Space>
              </Form.Item>
            )}

            {endType === 'on' && (
              <Form.Item
                name="until"
                rules={[
                  { required: true, message: 'Please select end date' },
                  {
                    validator: (_, value) => {
                      if (!value) return Promise.reject('Please select end date');
                      if (value.isBefore(startDate)) {
                        return Promise.reject('End date must be after start date');
                      }
                      return Promise.resolve();
                    }
                  }
                ]}
              >
                <DatePicker
                  placeholder="Select end date"
                  style={{ width: '100%' }}
                  disabledDate={(current) => current && current.isBefore(startDate)}
                />
              </Form.Item>
            )}
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default RecurrenceModal;