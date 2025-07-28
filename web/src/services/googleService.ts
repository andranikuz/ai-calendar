import { apiService } from './api';
import { GoogleIntegration, GoogleCalendarSync, GoogleCalendar } from '../types/api';

export class GoogleService {
  async getAuthUrl(): Promise<{ auth_url: string; state: string }> {
    return await apiService.get('/google/auth-url');
  }

  async handleCallback(code: string, state: string): Promise<GoogleIntegration> {
    return await apiService.post('/google/callback', { code, state });
  }

  async getIntegration(): Promise<GoogleIntegration | null> {
    try {
      return await apiService.get('/google/integration');
    } catch (error: unknown) {
      if ((error as { response?: { status: number } }).response?.status === 404) {
        return null;
      }
      throw error;
    }
  }

  async disconnect(): Promise<void> {
    await apiService.delete('/google/integration');
  }

  async getCalendars(): Promise<{ calendars: GoogleCalendar[] }> {
    return await apiService.get('/google/calendars');
  }

  // Calendar Sync Management
  async getCalendarSyncs(): Promise<{ syncs: GoogleCalendarSync[] }> {
    return await apiService.get('/google/calendar-syncs');
  }

  async createCalendarSync(syncData: Partial<GoogleCalendarSync>): Promise<GoogleCalendarSync> {
    return await apiService.post('/google/calendar-syncs', syncData);
  }

  async updateCalendarSync(id: string, data: Partial<GoogleCalendarSync>): Promise<GoogleCalendarSync> {
    return await apiService.put(`/google/calendar-syncs/${id}`, data);
  }

  async deleteCalendarSync(id: string): Promise<void> {
    await apiService.delete(`/google/calendar-syncs/${id}`);
  }

  async triggerSync(syncId: string): Promise<{
    message: string;
    synced_count: number;
    synced_at: string;
  }> {
    return await apiService.post(`/google/calendar-syncs/${syncId}/sync`);
  }

  async triggerSyncWithConflictDetection(syncId: string): Promise<{
    message: string;
    synced_events: number;
    detected_conflicts: number;
    resolved_conflicts: number;
    sync_status: string;
    conflicts: Array<{
      id: string;
      conflict_type: string;
      description: string;
      local_event?: {
        id: string;
        title: string;
        description?: string;
        location?: string;
        start_time: string;
        end_time: string;
      };
      google_event?: {
        id: string;
        title: string;
        description?: string;
        location?: string;
        start_time: string;
        end_time: string;
      };
    }>;
  }> {
    return await apiService.post(`/google/calendar-syncs/${syncId}/sync-with-conflicts`);
  }

  // Helper methods
  async isConnected(): Promise<boolean> {
    try {
      const integration = await this.getIntegration();
      return !!integration && integration.enabled;
    } catch {
      return false;
    }
  }

  getGoogleAuthUrl(authUrl: string): string {
    // Open in new window for OAuth flow
    return authUrl;
  }

  async waitForCallback(): Promise<GoogleIntegration> {
    return new Promise((resolve, reject) => {
      const checkInterval = setInterval(async () => {
        try {
          const integration = await this.getIntegration();
          if (integration) {
            clearInterval(checkInterval);
            resolve(integration);
          }
        } catch {
          // Continue checking
        }
      }, 1000);

      // Timeout after 5 minutes
      setTimeout(() => {
        clearInterval(checkInterval);
        reject(new Error('OAuth callback timeout'));
      }, 5 * 60 * 1000);
    });
  }
}

export const googleService = new GoogleService();