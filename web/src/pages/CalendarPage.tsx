import React, { useEffect, useRef, useState } from 'react';
import { Card, Button, Space, Select, Tooltip, Modal, message } from '../utils/antd';
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
import { EventClickArg, EventDropArg, DateSelectArg } from '@fullcalendar/core';
import type { EventResizeDoneArg } from '@fullcalendar/interaction';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { 
  fetchEvents,
  updateEvent,
  moveEvent,
  setCalendarView,
  navigateDate
} from '../store/slices/eventsSlice';
import { getIntegration, getCalendarSyncs, triggerSync } from '../store/slices/googleSlice';
import EventModal from '../components/Calendar/EventModal';
import { generateEventInstances } from '../utils/rrule';
import { Event } from '../types/api';
import dayjs from 'dayjs';

// Type alias for calendar usage
type CalendarEvent = Event;

// Extended type for FullCalendar events
interface FullCalendarEvent {
  id: string;
  title: string;
  start: string;
  end: string;
  color: string;
  borderColor?: string;
  opacity?: number;
  extendedProps: Record<string, unknown>;
}

const CalendarPage: React.FC = () => {
  const calendarRef = useRef<FullCalendar>(null);
  const dispatch = useAppDispatch();
  
  const { events, calendarView, currentDate, isLoading } = useAppSelector(state => state.events);
  const { calendarSyncs, isConnected } = useAppSelector(state => state.google);
  
  const [eventModalVisible, setEventModalVisible] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState<CalendarEvent | null>(null);
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

  const handleDateSelect = (selectInfo: DateSelectArg) => {
    setSelectedDate(selectInfo.startStr);
    setSelectedEvent(null);
    setEventModalVisible(true);
  };

  const handleEventClick = (clickInfo: EventClickArg) => {
    const { extendedProps } = clickInfo.event;
    
    if (extendedProps.isInstance) {
      // For recurring instances, edit the original event
      const originalEvent = extendedProps.originalEvent as CalendarEvent;
      setSelectedEvent(originalEvent);
    } else {
      // For regular events or the main recurring event
      const event = events.find(e => e.id === clickInfo.event.id);
      setSelectedEvent(event || null);
    }
    
    setSelectedDate(null);
    setEventModalVisible(true);
  };

  const handleEventDrop = async (dropInfo: EventDropArg) => {
    const { event } = dropInfo;
    const { extendedProps } = event;
    
    try {
      if (extendedProps.isInstance) {
        // For recurring instances, show a dialog asking how to handle the change
        Modal.confirm({
          title: 'Move Recurring Event',
          content: 'This is a recurring event. How would you like to move it?',
          okText: 'Move all occurrences',
          cancelText: 'Move only this occurrence',
          onOk: async () => {
            // Move the entire series by updating the original event
            const originalEvent = extendedProps.originalEvent as CalendarEvent;
            const timeDiff = dayjs(event.start).diff(dayjs(originalEvent.start_time));
            const newStart = dayjs(originalEvent.start_time).add(timeDiff, 'milliseconds');
            const newEnd = dayjs(originalEvent.end_time).add(timeDiff, 'milliseconds');
            
            await dispatch(updateEvent({
              id: originalEvent.id,
              data: {
                start_time: newStart.toISOString(),
                end_time: newEnd.toISOString()
              }
            })).unwrap();
            message.success('All occurrences moved successfully');
          },
          onCancel: () => {
            // TODO: In a full implementation, this would create an exception
            // For now, we'll just revert the change
            message.info('Moving single occurrences is not yet implemented');
            dropInfo.revert();
          }
        });
      } else {
        // Handle regular events
        await dispatch(moveEvent({
          id: event.id,
          start: event.start!.toISOString(),
          end: event.end?.toISOString() || event.start!.toISOString()
        })).unwrap();
        message.success('Event moved successfully');
      }
    } catch {
      message.error('Failed to move event');
      dropInfo.revert();
    }
  };

  const handleEventResize = async (resizeInfo: EventResizeDoneArg) => {
    const { event } = resizeInfo;
    const { extendedProps } = event;
    
    try {
      if (extendedProps.isInstance) {
        // For recurring instances, show a dialog
        Modal.confirm({
          title: 'Resize Recurring Event',
          content: 'This is a recurring event. How would you like to resize it?',
          okText: 'Resize all occurrences',
          cancelText: 'Resize only this occurrence',
          onOk: async () => {
            // Resize the entire series by updating the original event duration
            const originalEvent = extendedProps.originalEvent as CalendarEvent;
            const newDuration = dayjs(event.end).diff(dayjs(event.start));
            const originalStart = dayjs(originalEvent.start_time);
            const newEnd = originalStart.add(newDuration, 'milliseconds');
            
            await dispatch(updateEvent({
              id: originalEvent.id,
              data: {
                end_time: newEnd.toISOString()
              }
            })).unwrap();
            message.success('All occurrences resized successfully');
          },
          onCancel: () => {
            message.info('Resizing single occurrences is not yet implemented');
            resizeInfo.revert();
          }
        });
      } else {
        // Handle regular events
        await dispatch(updateEvent({
          id: event.id,
          data: {
            start_time: event.start!.toISOString(),
            end_time: event.end?.toISOString() || event.start!.toISOString()
          }
        })).unwrap();
        message.success('Event updated successfully');
      }
    } catch {
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
    } catch {
      message.error('Failed to sync with Google Calendar');
    }
  };

  const formatEventsForCalendar = (): FullCalendarEvent[] => {
    const calendarApi = calendarRef.current?.getApi();
    const currentView = calendarApi?.view;
    
    // Get the current view date range
    const viewStart = currentView?.currentStart || dayjs().subtract(1, 'month').toDate();
    const viewEnd = currentView?.currentEnd || dayjs().add(1, 'month').toDate();
    
    const allEvents: FullCalendarEvent[] = [];
    
    events.forEach(event => {
      // Add the main event
      allEvents.push({
        id: event.id,
        title: event.title,
        start: event.start_time,
        end: event.end_time,
        color: event.goal_id ? '#52c41a' : '#1677ff', // Green for goal-linked events
        borderColor: event.recurrence ? '#ff4d4f' : undefined, // Red border for recurring events
        extendedProps: {
          description: event.description,
          location: event.location,
          goalId: event.goal_id,
          externalSource: event.external_source,
          isRecurring: !!event.recurrence,
          originalEvent: event
        }
      });
      
      // Generate recurring instances if this is a recurring event
      if (event.recurrence) {
        try {
          const instances = generateEventInstances(event, viewStart, viewEnd);
          instances.forEach((instance, index) => {
            // Skip the first instance if it's the same as the original event
            const instanceStart = dayjs(instance.start_time);
            const originalStart = dayjs(event.start_time);
            
            if (index === 0 && instanceStart.isSame(originalStart, 'minute')) {
              return; // Skip duplicate of original event
            }
            
            allEvents.push({
              id: `${event.id}_${instance.instanceDate}`,
              title: `${event.title} (recurring)`,
              start: instance.start_time,
              end: instance.end_time,
              color: event.goal_id ? '#52c41a' : '#1677ff',
              borderColor: '#ff4d4f', // Red border for recurring instances
              opacity: 0.8, // Slightly transparent for recurring instances
              extendedProps: {
                description: instance.description,
                location: instance.location,
                goalId: instance.goal_id,
                externalSource: instance.external_source,
                isRecurring: true,
                isInstance: true,
                originalEvent: event,
                instanceDate: instance.instanceDate
              }
            });
          });
        } catch (error) {
          console.error('Error generating recurring instances for event:', event.id, error);
        }
      }
    });
    
    return allEvents;
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