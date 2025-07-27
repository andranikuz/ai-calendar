import { apiService } from './api';
import { Goal, Task, Milestone, CreateGoalRequest } from '../types/api';

export class GoalsService {
  async getGoals(params?: { offset?: number; limit?: number }): Promise<{ goals: Goal[] }> {
    const queryParams = new URLSearchParams();
    if (params?.offset) queryParams.append('offset', params.offset.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    
    const url = `/goals${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return await apiService.get(url);
  }

  async getGoal(id: string): Promise<Goal> {
    return await apiService.get(`/goals/${id}`);
  }

  async createGoal(goalData: CreateGoalRequest): Promise<Goal> {
    return await apiService.post('/goals', goalData);
  }

  async updateGoal(id: string, data: Partial<Goal>): Promise<Goal> {
    return await apiService.put(`/goals/${id}`, data);
  }

  async deleteGoal(id: string): Promise<void> {
    await apiService.delete(`/goals/${id}`);
  }

  // Tasks
  async getGoalTasks(goalId: string): Promise<{ tasks: Task[] }> {
    return await apiService.get(`/goals/${goalId}/tasks`);
  }

  async createTask(goalId: string, taskData: Partial<Task>): Promise<Task> {
    return await apiService.post(`/goals/${goalId}/tasks`, taskData);
  }

  async completeTask(taskId: string): Promise<Task> {
    return await apiService.post(`/goals/tasks/${taskId}/complete`);
  }

  // Milestones
  async getGoalMilestones(goalId: string): Promise<{ milestones: Milestone[] }> {
    return await apiService.get(`/goals/${goalId}/milestones`);
  }

  async createMilestone(goalId: string, milestoneData: Partial<Milestone>): Promise<Milestone> {
    return await apiService.post(`/goals/${goalId}/milestones`, milestoneData);
  }

  async completeMilestone(milestoneId: string): Promise<Milestone> {
    return await apiService.post(`/goals/milestones/${milestoneId}/complete`);
  }

  // Goal progress
  async updateProgress(goalId: string, progress: number): Promise<Goal> {
    return await apiService.patch(`/goals/${goalId}`, { progress });
  }
}

export const goalsService = new GoalsService();