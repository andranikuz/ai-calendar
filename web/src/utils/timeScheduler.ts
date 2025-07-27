import dayjs, { Dayjs } from 'dayjs';
import { Event, Goal, Task } from '../types/api';

export interface TimeSlot {
  start: Dayjs;
  end: Dayjs;
  duration: number; // in minutes
}

export interface SchedulingSuggestion {
  goal: Goal;
  task?: Task;
  suggestedSlots: TimeSlot[];
  totalTimeNeeded: number; // in minutes
  priority: 'critical' | 'high' | 'medium' | 'low';
  deadline?: Dayjs;
  reason: string;
}

export interface WorkingHours {
  start: string; // "09:00"
  end: string;   // "17:00"
  timezone: string;
}

export interface SchedulingPreferences {
  workingHours: WorkingHours;
  workingDays: number[]; // 0-6, where 0 is Sunday
  minSessionDuration: number; // minimum minutes for a work session
  maxSessionDuration: number; // maximum minutes for a work session
  breakBetweenSessions: number; // minutes break between sessions
  avoidTimeSlots: TimeSlot[]; // times to avoid (lunch, commute, etc.)
}

const DEFAULT_PREFERENCES: SchedulingPreferences = {
  workingHours: {
    start: '09:00',
    end: '17:00',
    timezone: 'America/New_York'
  },
  workingDays: [1, 2, 3, 4, 5], // Monday to Friday
  minSessionDuration: 30,
  maxSessionDuration: 120,
  breakBetweenSessions: 15,
  avoidTimeSlots: [
    // Lunch break
    {
      start: dayjs().hour(12).minute(0),
      end: dayjs().hour(13).minute(0),
      duration: 60
    }
  ]
};

/**
 * Find available time slots in the calendar
 */
export function findAvailableTimeSlots(
  events: Event[],
  startDate: Dayjs,
  endDate: Dayjs,
  preferences: SchedulingPreferences = DEFAULT_PREFERENCES
): TimeSlot[] {
  const availableSlots: TimeSlot[] = [];
  const { workingHours, workingDays, minSessionDuration } = preferences;

  // Iterate through each day in the range
  let currentDate = startDate.startOf('day');
  
  while (currentDate.isBefore(endDate) || currentDate.isSame(endDate, 'day')) {
    // Skip non-working days
    if (!workingDays.includes(currentDate.day())) {
      currentDate = currentDate.add(1, 'day');
      continue;
    }

    // Create working hours for this day
    const dayStart = currentDate
      .hour(parseInt(workingHours.start.split(':')[0]))
      .minute(parseInt(workingHours.start.split(':')[1]));
    
    const dayEnd = currentDate
      .hour(parseInt(workingHours.end.split(':')[0]))
      .minute(parseInt(workingHours.end.split(':')[1]));

    // Get events for this day
    const dayEvents = events.filter(event => 
      dayjs(event.start_time).isSame(currentDate, 'day')
    ).sort((a, b) => dayjs(a.start_time).diff(dayjs(b.start_time)));

    // Find gaps between events
    let slotStart = dayStart;
    
    for (const event of dayEvents) {
      const eventStart = dayjs(event.start_time);
      const eventEnd = dayjs(event.end_time);
      
      // Check if there's a gap before this event
      if (eventStart.diff(slotStart, 'minutes') >= minSessionDuration) {
        availableSlots.push({
          start: slotStart,
          end: eventStart,
          duration: eventStart.diff(slotStart, 'minutes')
        });
      }
      
      // Move slot start to after this event
      slotStart = eventEnd.add(preferences.breakBetweenSessions, 'minutes');
    }
    
    // Check for time after the last event
    if (dayEnd.diff(slotStart, 'minutes') >= minSessionDuration) {
      availableSlots.push({
        start: slotStart,
        end: dayEnd,
        duration: dayEnd.diff(slotStart, 'minutes')
      });
    }

    currentDate = currentDate.add(1, 'day');
  }

  // Filter out avoided time slots
  return filterAvoidedTimeSlots(availableSlots, preferences.avoidTimeSlots);
}

/**
 * Filter out avoided time slots
 */
function filterAvoidedTimeSlots(slots: TimeSlot[], avoidSlots: TimeSlot[]): TimeSlot[] {
  const filteredSlots: TimeSlot[] = [];

  for (const slot of slots) {
    let currentSlot = slot;
    let hasConflict = false;

    for (const avoidSlot of avoidSlots) {
      // Check if avoid slot overlaps with current slot
      const avoidStart = avoidSlot.start.hour(currentSlot.start.hour()).minute(currentSlot.start.minute());
      const avoidEnd = avoidSlot.end.hour(currentSlot.start.hour()).minute(currentSlot.start.minute());

      if (avoidStart.isBefore(currentSlot.end) && avoidEnd.isAfter(currentSlot.start)) {
        // Split the slot around the avoided time
        if (avoidStart.isAfter(currentSlot.start)) {
          // Add time before avoided slot
          const beforeDuration = avoidStart.diff(currentSlot.start, 'minutes');
          if (beforeDuration >= 30) { // Minimum 30 minutes
            filteredSlots.push({
              start: currentSlot.start,
              end: avoidStart,
              duration: beforeDuration
            });
          }
        }

        if (avoidEnd.isBefore(currentSlot.end)) {
          // Update current slot to time after avoided slot
          currentSlot = {
            start: avoidEnd,
            end: currentSlot.end,
            duration: currentSlot.end.diff(avoidEnd, 'minutes')
          };
        } else {
          hasConflict = true;
          break;
        }
      }
    }

    if (!hasConflict && currentSlot.duration >= 30) {
      filteredSlots.push(currentSlot);
    }
  }

  return filteredSlots;
}

/**
 * Calculate time needed for a goal based on its tasks
 */
export function calculateTimeNeeded(goal: Goal, tasks: Task[]): number {
  const goalTasks = tasks.filter(task => task.goal_id === goal.id);
  
  // If tasks have estimated durations, use them
  const estimatedTime = goalTasks.reduce((total, task) => {
    return total + (task.estimated_duration || 0);
  }, 0);

  if (estimatedTime > 0) {
    return estimatedTime;
  }

  // Fallback: estimate based on goal complexity and deadline
  const daysUntilDeadline = goal.deadline ? 
    dayjs(goal.deadline).diff(dayjs(), 'days') : 30;
  
  // Base time estimates by category (in minutes)
  const categoryTimeEstimates: Record<string, number> = {
    health: 240,      // 4 hours per week
    career: 480,      // 8 hours per week  
    education: 360,   // 6 hours per week
    personal: 180,    // 3 hours per week
    financial: 120,   // 2 hours per week
    relationship: 150 // 2.5 hours per week
  };

  const baseTime = categoryTimeEstimates[goal.category] || 240;
  const weeks = Math.max(1, Math.ceil(daysUntilDeadline / 7));
  
  // Adjust based on priority
  const priorityMultipliers = {
    critical: 1.5,
    high: 1.2,
    medium: 1.0,
    low: 0.8
  };

  return Math.round(baseTime * weeks * (priorityMultipliers[goal.priority] || 1));
}

/**
 * Generate scheduling suggestions for goals
 */
export function generateSchedulingSuggestions(
  goals: Goal[],
  tasks: Task[],
  events: Event[],
  preferences: SchedulingPreferences = DEFAULT_PREFERENCES,
  lookaheadDays: number = 14
): SchedulingSuggestion[] {
  const suggestions: SchedulingSuggestion[] = [];
  const startDate = dayjs();
  const endDate = startDate.add(lookaheadDays, 'days');

  // Get available time slots
  const availableSlots = findAvailableTimeSlots(events, startDate, endDate, preferences);
  
  // Filter active goals that need scheduling
  const activeGoals = goals.filter(goal => 
    goal.status === 'active' && 
    goal.progress < 100
  );

  // Sort goals by priority and deadline
  const sortedGoals = activeGoals.sort((a, b) => {
    const priorityOrder = { critical: 4, high: 3, medium: 2, low: 1 };
    const aPriority = priorityOrder[a.priority] || 1;
    const bPriority = priorityOrder[b.priority] || 1;
    
    if (aPriority !== bPriority) {
      return bPriority - aPriority; // Higher priority first
    }
    
    // If same priority, sort by deadline (earlier first)
    if (a.deadline && b.deadline) {
      return dayjs(a.deadline).diff(dayjs(b.deadline));
    }
    if (a.deadline) return -1;
    if (b.deadline) return 1;
    
    return 0;
  });

  let remainingSlots = [...availableSlots];

  for (const goal of sortedGoals) {
    const timeNeeded = calculateTimeNeeded(goal, tasks);
    const goalTasks = tasks.filter(task => 
      task.goal_id === goal.id && 
      task.status !== 'completed'
    );

    // Find suitable time slots for this goal
    const suggestedSlots = allocateTimeSlots(
      remainingSlots,
      timeNeeded,
      preferences,
      goal.deadline ? dayjs(goal.deadline) : undefined
    );

    if (suggestedSlots.length > 0) {
      // Generate reason for scheduling
      const reason = generateSchedulingReason(goal, timeNeeded, goalTasks.length);
      
      suggestions.push({
        goal,
        task: goalTasks[0], // Next task to work on
        suggestedSlots,
        totalTimeNeeded: timeNeeded,
        priority: goal.priority,
        deadline: goal.deadline ? dayjs(goal.deadline) : undefined,
        reason
      });

      // Remove allocated slots from remaining slots
      remainingSlots = removeAllocatedSlots(remainingSlots, suggestedSlots);
    }
  }

  return suggestions;
}

/**
 * Allocate time slots for a specific time requirement
 */
function allocateTimeSlots(
  availableSlots: TimeSlot[],
  timeNeeded: number,
  preferences: SchedulingPreferences,
  deadline?: Dayjs
): TimeSlot[] {
  const allocatedSlots: TimeSlot[] = [];
  let remainingTime = timeNeeded;
  const { maxSessionDuration, minSessionDuration } = preferences;

  // Sort slots by preference (earlier times first, closer to deadline if applicable)
  const sortedSlots = [...availableSlots].sort((a, b) => {
    if (deadline) {
      // Prefer slots closer to deadline
      const aDistance = Math.abs(deadline.diff(a.start, 'days'));
      const bDistance = Math.abs(deadline.diff(b.start, 'days'));
      if (aDistance !== bDistance) {
        return aDistance - bDistance;
      }
    }
    
    // Prefer earlier times in the day
    return a.start.hour() - b.start.hour();
  });

  for (const slot of sortedSlots) {
    if (remainingTime <= 0) break;

    const sessionDuration = Math.min(
      remainingTime,
      slot.duration,
      maxSessionDuration
    );

    if (sessionDuration >= minSessionDuration) {
      allocatedSlots.push({
        start: slot.start,
        end: slot.start.add(sessionDuration, 'minutes'),
        duration: sessionDuration
      });

      remainingTime -= sessionDuration;
    }
  }

  return allocatedSlots;
}

/**
 * Remove allocated slots from available slots
 */
function removeAllocatedSlots(availableSlots: TimeSlot[], allocatedSlots: TimeSlot[]): TimeSlot[] {
  let remaining = [...availableSlots];

  for (const allocated of allocatedSlots) {
    remaining = remaining.reduce<TimeSlot[]>((acc, slot) => {
      // Check if allocated slot overlaps with available slot
      if (allocated.start.isBefore(slot.end) && allocated.end.isAfter(slot.start)) {
        // Split the slot
        if (allocated.start.isAfter(slot.start)) {
          // Add time before allocated slot
          acc.push({
            start: slot.start,
            end: allocated.start,
            duration: allocated.start.diff(slot.start, 'minutes')
          });
        }
        
        if (allocated.end.isBefore(slot.end)) {
          // Add time after allocated slot
          acc.push({
            start: allocated.end,
            end: slot.end,
            duration: slot.end.diff(allocated.end, 'minutes')
          });
        }
      } else {
        // No overlap, keep the slot
        acc.push(slot);
      }
      
      return acc;
    }, []);
  }

  return remaining.filter(slot => slot.duration >= 30); // Filter out too small slots
}

/**
 * Generate human-readable reason for scheduling
 */
function generateSchedulingReason(goal: Goal, timeNeeded: number, pendingTasks: number): string {
  const hours = Math.round(timeNeeded / 60 * 10) / 10;
  const days = goal.deadline ? dayjs(goal.deadline).diff(dayjs(), 'days') : null;
  
  let reason = `Allocate ${hours}h total`;
  
  if (pendingTasks > 0) {
    reason += ` across ${pendingTasks} tasks`;
  }
  
  if (days !== null) {
    if (days <= 7) {
      reason += ` (deadline in ${days} days - urgent!)`;
    } else if (days <= 30) {
      reason += ` (deadline in ${days} days)`;
    }
  }
  
  switch (goal.priority) {
    case 'critical':
      reason += ' - Critical priority';
      break;
    case 'high':
      reason += ' - High priority';
      break;
  }
  
  return reason;
}

/**
 * Create calendar events from scheduling suggestions
 */
export function createScheduledEvents(
  suggestions: SchedulingSuggestion[],
  userId: string
): Omit<Event, 'id' | 'created_at' | 'updated_at'>[] {
  const events: Omit<Event, 'id' | 'created_at' | 'updated_at'>[] = [];

  for (const suggestion of suggestions) {
    for (const slot of suggestion.suggestedSlots) {
      const taskTitle = suggestion.task ? 
        `${suggestion.goal.title}: ${suggestion.task.title}` :
        `Work on: ${suggestion.goal.title}`;

      events.push({
        user_id: userId,
        goal_id: suggestion.goal.id,
        title: taskTitle,
        description: `Scheduled work session: ${suggestion.reason}`,
        start_time: slot.start.toISOString(),
        end_time: slot.end.toISOString(),
        timezone: dayjs.tz.guess(),
        location: '',
        attendees: [],
        status: 'tentative', // Mark as tentative until user confirms
      });
    }
  }

  return events;
}

export default {
  findAvailableTimeSlots,
  calculateTimeNeeded,
  generateSchedulingSuggestions,
  createScheduledEvents,
  DEFAULT_PREFERENCES
};