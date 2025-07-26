import { apiService } from './api';
import { Event, CreateEventRequest } from '../types/api';

export class EventsService {
  async getEvents(params?: { 
    start?: string; 
    end?: string; 
    offset?: number; 
    limit?: number 
  }): Promise<{ events: Event[] }> {
    const queryParams = new URLSearchParams();
    if (params?.start) queryParams.append('start', params.start);
    if (params?.end) queryParams.append('end', params.end);
    if (params?.offset) queryParams.append('offset', params.offset.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    
    const url = `/events${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return await apiService.get(url);
  }

  async getEvent(id: string): Promise<Event> {
    return await apiService.get(`/events/${id}`);
  }

  async createEvent(eventData: CreateEventRequest): Promise<Event> {
    return await apiService.post('/events', eventData);
  }

  async updateEvent(id: string, data: Partial<Event>): Promise<Event> {
    return await apiService.put(`/events/${id}`, data);
  }

  async deleteEvent(id: string): Promise<void> {
    await apiService.delete(`/events/${id}`);
  }

  async moveEvent(id: string, startTime: string, endTime: string): Promise<Event> {
    return await apiService.post(`/events/${id}/move`, {
      start_time: startTime,
      end_time: endTime,
    });
  }

  async duplicateEvent(id: string): Promise<Event> {
    return await apiService.post(`/events/${id}/duplicate`);
  }

  async linkToGoal(eventId: string, goalId: string): Promise<Event> {
    return await apiService.post(`/events/${eventId}/link-goal`, { goal_id: goalId });
  }

  async unlinkFromGoal(eventId: string): Promise<Event> {
    return await apiService.post(`/events/${eventId}/unlink-goal`);
  }

  async getUpcomingEvents(limit: number = 10): Promise<{ events: Event[] }> {
    return await apiService.get(`/events/upcoming?limit=${limit}`);
  }

  async getTodayEvents(): Promise<{ events: Event[] }> {
    return await apiService.get('/events/today');
  }

  async searchEvents(query: string, limit: number = 20): Promise<{ events: Event[] }> {
    return await apiService.get(`/events/search?q=${encodeURIComponent(query)}&limit=${limit}`);
  }

  async checkConflicts(startTime: string, endTime: string, excludeEventId?: string): Promise<{ 
    has_conflict: boolean; 
    conflicting_events: Event[] 
  }> {
    const params = new URLSearchParams({
      start: startTime,
      end: endTime,
    });
    
    if (excludeEventId) {
      params.append('exclude_event_id', excludeEventId);
    }
    
    return await apiService.get(`/events/conflict-check?${params.toString()}`);
  }

  async getEventsByTimeRange(start: string, end: string): Promise<{ events: Event[] }> {
    return await apiService.get(`/events/time-range?start=${start}&end=${end}`);
  }
}

export const eventsService = new EventsService();