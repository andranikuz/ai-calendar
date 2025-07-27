import { apiService } from './api';
import { Mood, CreateMoodRequest } from '../types/api';

export class MoodsService {
  async getMoods(params?: { 
    start?: string; 
    end?: string; 
    offset?: number; 
    limit?: number 
  }): Promise<{ moods: Mood[] }> {
    const queryParams = new URLSearchParams();
    if (params?.start) queryParams.append('start', params.start);
    if (params?.end) queryParams.append('end', params.end);
    if (params?.offset) queryParams.append('offset', params.offset.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    
    const url = `/moods${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return await apiService.get(url);
  }

  async getMood(id: string): Promise<Mood> {
    return await apiService.get(`/moods/${id}`);
  }

  async createMood(moodData: CreateMoodRequest): Promise<Mood> {
    return await apiService.post('/moods', moodData);
  }

  async updateMood(id: string, data: Partial<Mood>): Promise<Mood> {
    return await apiService.put(`/moods/${id}`, data);
  }

  async deleteMood(id: string): Promise<void> {
    await apiService.delete(`/moods/${id}`);
  }

  async getTodayMood(): Promise<Mood | null> {
    try {
      return await apiService.get('/moods/latest');
    } catch (error: unknown) {
      if ((error as { response?: { status: number } }).response?.status === 404) {
        return null;
      }
      throw error;
    }
  }

  async getMoodByDate(date: string): Promise<Mood | null> {
    try {
      return await apiService.get(`/moods/by-date?date=${date}`);
    } catch (error: unknown) {
      if ((error as { response?: { status: number } }).response?.status === 404) {
        return null;
      }
      throw error;
    }
  }

  async upsertMoodByDate(date: string, level: number, notes?: string, tags?: string[]): Promise<Mood> {
    return await apiService.post('/moods/upsert-by-date', {
      date,
      level,
      notes,
      tags,
    });
  }

  async getMoodStats(params?: { days?: number }): Promise<{
    average: number;
    total: number;
    trend: 'up' | 'down' | 'stable';
    weekAverage?: number;
    weeklyData?: Array<{ level: number; date: string } | null>;
    distribution?: { [key: number]: number };
    daily_averages: Array<{ date: string; average: number }>;
  }> {
    const queryParams = new URLSearchParams();
    if (params?.days) queryParams.append('days', params.days.toString());
    
    const url = `/moods/stats${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return await apiService.get(url);
  }

  async getMoodTrends(params?: { days?: number }): Promise<{
    trends: Array<{ 
      date: string; 
      level: number; 
      change: number;
      trend: 'up' | 'down' | 'stable';
    }>;
    overall_trend: 'up' | 'down' | 'stable';
    correlation_with_productivity?: number;
  }> {
    const queryParams = new URLSearchParams();
    if (params?.days) queryParams.append('days', params.days.toString());
    
    const url = `/moods/trends${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return await apiService.get(url);
  }

  async getMoodsByDateRange(start: string, end: string): Promise<{ moods: Mood[] }> {
    return await apiService.get(`/moods/date-range?start=${start}&end=${end}`);
  }
}

export const moodsService = new MoodsService();