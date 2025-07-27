import { RRule, Frequency } from 'rrule';
import { Recurrence, Event } from '../types/api';
import dayjs from 'dayjs';

/**
 * Convert our Recurrence type to RRule options
 */
export function recurrenceToRRule(recurrence: Recurrence, startDate: Date): RRule {
  const freq = getFrequency(recurrence.freq);
  
  const options: Record<string, unknown> = {
    freq,
    dtstart: startDate,
    interval: recurrence.interval || 1,
  };

  if (recurrence.count) {
    options.count = recurrence.count;
  }

  if (recurrence.until) {
    options.until = new Date(recurrence.until);
  }

  if (recurrence.byweekday && recurrence.byweekday.length > 0) {
    options.byweekday = recurrence.byweekday;
  }

  if (recurrence.bymonthday && recurrence.bymonthday.length > 0) {
    options.bymonthday = recurrence.bymonthday;
  }

  if (recurrence.bymonth && recurrence.bymonth.length > 0) {
    options.bymonth = recurrence.bymonth;
  }

  if (recurrence.byhour && recurrence.byhour.length > 0) {
    options.byhour = recurrence.byhour;
  }

  if (recurrence.byminute && recurrence.byminute.length > 0) {
    options.byminute = recurrence.byminute;
  }

  return new RRule(options);
}

/**
 * Convert RRule string to our Recurrence type
 */
export function rruleToRecurrence(rruleString: string): Recurrence | null {
  try {
    const rule = RRule.fromString(rruleString);
    const options = rule.options;

    const recurrence: Recurrence = {
      freq: getFrequencyString(options.freq),
      interval: options.interval || 1,
    };

    if (options.count) {
      recurrence.count = options.count;
    }

    if (options.until) {
      recurrence.until = options.until.toISOString();
    }

    if (options.byweekday && options.byweekday.length > 0) {
      recurrence.byweekday = options.byweekday.map((day: any) => 
        typeof day === 'number' ? day : day.weekday
      );
    }

    if (options.bymonthday && options.bymonthday.length > 0) {
      recurrence.bymonthday = options.bymonthday;
    }

    if (options.bymonth && options.bymonth.length > 0) {
      recurrence.bymonth = options.bymonth;
    }

    if (options.byhour && options.byhour.length > 0) {
      recurrence.byhour = options.byhour;
    }

    if (options.byminute && options.byminute.length > 0) {
      recurrence.byminute = options.byminute;
    }

    return recurrence;
  } catch (error) {
    console.error('Failed to parse RRule:', error);
    return null;
  }
}

/**
 * Generate event instances from a recurring event
 */
export function generateEventInstances(
  event: Event, 
  startDate: Date, 
  endDate: Date
): Array<Omit<Event, 'id'> & { instanceDate: string; isRecurring: boolean }> {
  if (!event.recurrence) {
    return [];
  }

  try {
    const rule = recurrenceToRRule(event.recurrence, new Date(event.start_time));
    const instances = rule.between(startDate, endDate, true);
    
    const eventDuration = dayjs(event.end_time).diff(dayjs(event.start_time));
    
    return instances.map(instanceStart => {
      const instanceEnd = dayjs(instanceStart).add(eventDuration, 'milliseconds');
      
      return {
        ...event,
        start_time: instanceStart.toISOString(),
        end_time: instanceEnd.toISOString(),
        instanceDate: instanceStart.toISOString(),
        isRecurring: true,
      };
    });
  } catch (error) {
    console.error('Failed to generate event instances:', error);
    return [];
  }
}

/**
 * Check if a date is an exception (deleted occurrence) for a recurring event
 */
export function isExceptionDate(event: Event, date: Date): boolean {
  // TODO: Implement exception handling when we add support for it
  return false;
}

/**
 * Get the next occurrence of a recurring event after a given date
 */
export function getNextOccurrence(event: Event, after: Date): Date | null {
  if (!event.recurrence) {
    return null;
  }

  try {
    const rule = recurrenceToRRule(event.recurrence, new Date(event.start_time));
    const nextOccurrence = rule.after(after);
    return nextOccurrence;
  } catch (error) {
    console.error('Failed to get next occurrence:', error);
    return null;
  }
}

/**
 * Get a human-readable description of the recurrence
 */
export function getRecurrenceDescription(recurrence: Recurrence): string {
  try {
    const rule = recurrenceToRRule(recurrence, new Date());
    return rule.toText();
  } catch (error) {
    console.error('Failed to get recurrence description:', error);
    return 'Invalid recurrence pattern';
  }
}

/**
 * Validate a recurrence pattern
 */
export function validateRecurrence(recurrence: Recurrence): { valid: boolean; error?: string } {
  try {
    recurrenceToRRule(recurrence, new Date());
    return { valid: true };
  } catch (error) {
    return { 
      valid: false, 
      error: error instanceof Error ? error.message : 'Invalid recurrence pattern' 
    };
  }
}

// Helper functions
function getFrequency(freq: Recurrence['freq']): Frequency {
  switch (freq) {
    case 'DAILY': return RRule.DAILY;
    case 'WEEKLY': return RRule.WEEKLY;
    case 'MONTHLY': return RRule.MONTHLY;
    case 'YEARLY': return RRule.YEARLY;
    default: return RRule.WEEKLY;
  }
}

function getFrequencyString(freq: Frequency): Recurrence['freq'] {
  switch (freq) {
    case RRule.DAILY: return 'DAILY';
    case RRule.WEEKLY: return 'WEEKLY';
    case RRule.MONTHLY: return 'MONTHLY';
    case RRule.YEARLY: return 'YEARLY';
    default: return 'WEEKLY';
  }
}

/**
 * Convert recurrence to RFC 5545 RRULE string for API
 */
export function recurrenceToRRuleString(recurrence: Recurrence, startDate: Date): string {
  const rule = recurrenceToRRule(recurrence, startDate);
  return rule.toString();
}

/**
 * Parse RFC 5545 RRULE string from API
 */
export function parseRRuleString(rruleString: string): Recurrence | null {
  return rruleToRecurrence(rruleString);
}