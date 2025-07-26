import { apiService } from './api';
import { LoginRequest, RegisterRequest, AuthTokens, User } from '../types/api';

export class AuthService {
  async login(credentials: LoginRequest): Promise<{ tokens: AuthTokens; user: User }> {
    const response = await apiService.post<{ tokens: AuthTokens; user: User }>('/auth/login', credentials);
    return response;
  }

  async register(userData: RegisterRequest): Promise<{ tokens: AuthTokens; user: User }> {
    const response = await apiService.post<{ tokens: AuthTokens; user: User }>('/auth/register', userData);
    return response;
  }

  async getCurrentUser(): Promise<User> {
    return await apiService.get('/users/me');
  }

  async updateUser(userData: Partial<User>): Promise<User> {
    return await apiService.put('/users/me', userData);
  }

  async refreshToken(refreshToken: string): Promise<AuthTokens> {
    const response = await apiService.post<AuthTokens>('/auth/refresh', {
      refresh_token: refreshToken,
    });
    return response;
  }

  setTokens(tokens: AuthTokens) {
    apiService.setTokens(tokens);
  }

  logout() {
    apiService.clearTokens();
  }

  getStoredTokens(): AuthTokens | null {
    try {
      const stored = localStorage.getItem('auth_tokens');
      return stored ? JSON.parse(stored) : null;
    } catch {
      return null;
    }
  }

  isAuthenticated(): boolean {
    const tokens = this.getStoredTokens();
    if (!tokens) return false;

    // Check if token is expired
    const now = Date.now() / 1000;
    const expiresAt = tokens.expires_in + (Date.now() / 1000 - tokens.expires_in);
    
    return now < expiresAt;
  }
}

export const authService = new AuthService();