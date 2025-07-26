import React, { useEffect, useRef, useState } from 'react';
import { Card, Button, Space, Select, Tooltip, Modal, message } from 'antd';
import { 
  LeftOutlined, 
  RightOutlined, 
  PlusOutlined,
  SyncOutlined,
  CalendarOutlined
} from '@ant-design/icons';
import FullCalendar from '@fullcalendar/react';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import interactionPlugin from '@fullcalendar/interaction';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { 
  fetchEvents,
  createEvent,
  updateEvent,
  deleteEvent,
  moveEvent,
  setCalendarView,
  setCurrentDate,
  navigateDate
} from '../store/slices/eventsSlice';
import { getIntegration, getCalendarSyncs, triggerSync } from '../store/slices/googleSlice';
import EventModal from '../components/Calendar/EventModal';
import dayjs from 'dayjs';

const CalendarPage: React.FC = () => {
  const calendarRef = useRef<FullCalendar>(null);
  const dispatch = useAppDispatch();
  
  const { events, calendarView, currentDate, isLoading } = useAppSelector(state => state.events);
  const { integration, calendarSyncs, isConnected } = useAppSelector(state => state.google);
  
  const [eventModalVisible, setEventModalVisible] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState<any>(null);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);

  useEffect(() => {
    // Fetch events and Google integration status
    dispatch(fetchEvents({}));
    dispatch(getIntegration());
    dispatch(getCalendarSyncs());
  }, [dispatch]);

  useEffect(() => {
    // Update FullCalendar view when Redux state changes
    const calendarApi = calendarRef.current?.getApi();
    if (calendarApi) {
      calendarApi.changeView(calendarView === 'month' ? 'dayGridMonth' : 
                           calendarView === 'week' ? 'timeGridWeek' : 'timeGridDay');
      calendarApi.gotoDate(currentDate);
    }
  }, [calendarView, currentDate]);

  const handleViewChange = (view: 'month' | 'week' | 'day') => {
    dispatch(setCalendarView(view));
  };

  const handleNavigation = (direction: 'prev' | 'next' | 'today') => {
    dispatch(navigateDate(direction));
  };

  const handleDateSelect = (selectInfo: any) => {
    setSelectedDate(selectInfo.startStr);
    setSelectedEvent(null);
    setEventModalVisible(true);
  };

  const handleEventClick = (clickInfo: any) => {
    const event = events.find(e => e.id === clickInfo.event.id);
    setSelectedEvent(event);
    setSelectedDate(null);
    setEventModalVisible(true);
  };

  const handleEventDrop = async (dropInfo: any) => {
    const { event } = dropInfo;
    try {
      await dispatch(moveEvent({
        id: event.id,
        start: event.start.toISOString(),
        end: event.end?.toISOString() || event.start.toISOString()
      })).unwrap();
      message.success('Event moved successfully');
    } catch (error) {
      message.error('Failed to move event');
      dropInfo.revert();
    }
  };

  const handleEventResize = async (resizeInfo: any) => {
    const { event } = resizeInfo;
    try {
      await dispatch(updateEvent({
        id: event.id,
        data: {
          start_time: event.start.toISOString(),
          end_time: event.end?.toISOString() || event.start.toISOString()
        }
      })).unwrap();
      message.success('Event updated successfully');
    } catch (error) {
      message.error('Failed to update event');
      resizeInfo.revert();
    }
  };

  const handleSyncWithGoogle = async () => {
    if (!isConnected || calendarSyncs.length === 0) {
      message.warning('Google Calendar not connected or no sync configurations found');
      return;
    }

    try {
      for (const sync of calendarSyncs) {
        await dispatch(triggerSync(sync.id)).unwrap();
      }
      message.success('Calendar synchronized successfully');
      dispatch(fetchEvents({})); // Refresh events
    } catch (error) {
      message.error('Failed to sync with Google Calendar');
    }
  };

  const formatEventsForCalendar = () => {
    return events.map(event => ({
      id: event.id,
      title: event.title,
      start: event.start_time,
      end: event.end_time,
      color: event.goal_id ? '#52c41a' : '#1677ff', // Green for goal-linked events
      extendedProps: {
        description: event.description,
        location: event.location,
        goalId: event.goal_id,
        externalSource: event.external_source
      }
    }));
  };

  return (
    <div>
      {/* Calendar Header */}
      <div style={{ marginBottom: 16 }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center',
          marginBottom: 16 
        }}>
          <Space>
            <Button.Group>
              <Button 
                icon={<LeftOutlined />} 
                onClick={() => handleNavigation('prev')}
              />
              <Button onClick={() => handleNavigation('today')}>
                Today
              </Button>
              <Button 
                icon={<RightOutlined />} 
                onClick={() => handleNavigation('next')}
              />
            </Button.Group>
            
            <h2 style={{ margin: 0 }}>
              {dayjs(currentDate).format('MMMM YYYY')}
            </h2>
          </Space>

          <Space>
            <Select
              value={calendarView}
              onChange={handleViewChange}
              style={{ width: 120 }}
            >
              <Select.Option value="month">Month</Select.Option>
              <Select.Option value="week">Week</Select.Option>
              <Select.Option value="day">Day</Select.Option>
            </Select>

            {isConnected && (
              <Tooltip title="Sync with Google Calendar">
                <Button
                  icon={<SyncOutlined />}
                  onClick={handleSyncWithGoogle}
                  loading={isLoading}
                >
                  Sync
                </Button>
              </Tooltip>
            )}

            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setSelectedEvent(null);
                setSelectedDate(null);
                setEventModalVisible(true);
              }}
            >
              New Event
            </Button>
          </Space>
        </div>

        {/* Calendar Integration Status */}
        {!isConnected && (
          <div style={{ 
            background: '#fff7e6', 
            border: '1px solid #ffd591',
            borderRadius: 6,
            padding: 12,
            marginBottom: 16
          }}>
            <Space>
              <CalendarOutlined />
              <span>Connect your Google Calendar for automatic synchronization</span>
              <Button type="link" onClick={() => {
                // Navigate to settings to connect Google
              }}>
                Connect Now
              </Button>
            </Space>
          </div>
        )}
      </div>

      {/* Calendar */}
      <Card>
        <FullCalendar
          ref={calendarRef}
          plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin]}
          headerToolbar={false} // We handle navigation manually
          initialView="dayGridMonth"
          events={formatEventsForCalendar()}
          selectable={true}
          selectMirror={true}
          editable={true}
          dayMaxEvents={true}
          weekends={true}
          height="auto"
          select={handleDateSelect}
          eventClick={handleEventClick}
          eventDrop={handleEventDrop}
          eventResize={handleEventResize}
          eventDisplay="block"
          dayHeaderFormat={{ weekday: 'short' }}
          slotMinTime="06:00:00"
          slotMaxTime="22:00:00"
          nowIndicator={true}
          eventTimeFormat={{
            hour: 'numeric',
            minute: '2-digit',
            omitZeroMinute: false
          }}
        />
      </Card>

      {/* Event Modal */}
      <EventModal
        visible={eventModalVisible}
        event={selectedEvent}
        selectedDate={selectedDate}
        onCancel={() => {
          setEventModalVisible(false);
          setSelectedEvent(null);
          setSelectedDate(null);
        }}
        onSuccess={() => {
          setEventModalVisible(false);
          setSelectedEvent(null);
          setSelectedDate(null);
          dispatch(fetchEvents({})); // Refresh events
        }}
      />
    </div>
  );
};

export default CalendarPage;