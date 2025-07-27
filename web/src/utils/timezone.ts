import dayjs, { Dayjs } from 'dayjs';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

// Extend dayjs with timezone support
dayjs.extend(utc);
dayjs.extend(timezone);

export interface TimezoneInfo {
  name: string;
  label: string;
  offset: string;
  abbreviation: string;
}

/**
 * Get user's current timezone
 */
export function getUserTimezone(): string {
  return Intl.DateTimeFormat().resolvedOptions().timeZone;
}

/**
 * Get common timezones for selection
 */
export function getCommonTimezones(): TimezoneInfo[] {
  const now = dayjs();
  const timezones = [
    'America/New_York',
    'America/Chicago',
    'America/Denver', 
    'America/Los_Angeles',
    'America/Toronto',
    'America/Vancouver',
    'Europe/London',
    'Europe/Paris',
    'Europe/Berlin',
    'Europe/Rome',
    'Europe/Amsterdam',
    'Europe/Stockholm',
    'Europe/Moscow',
    'Asia/Tokyo',
    'Asia/Shanghai',
    'Asia/Hong_Kong',
    'Asia/Singapore',
    'Asia/Seoul',
    'Asia/Dubai',
    'Asia/Kolkata',
    'Australia/Sydney',
    'Australia/Melbourne',
    'Pacific/Auckland',
    'UTC'
  ];

  return timezones.map(tz => {
    const inTimezone = now.tz(tz);
    const offset = inTimezone.format('Z');
    const abbreviation = inTimezone.format('z');
    
    return {
      name: tz,
      label: `${tz.replace(/_/g, ' ')} (${offset})`,
      offset,
      abbreviation
    };
  }).sort((a, b) => a.label.localeCompare(b.label));
}

/**
 * Convert a date from one timezone to another
 */
export function convertTimezone(
  date: string | Date | Dayjs, 
  fromTimezone: string, 
  toTimezone: string
): Dayjs {
  return dayjs.tz(date, fromTimezone).tz(toTimezone);
}

/**
 * Convert UTC date to user's local timezone
 */
export function utcToLocal(date: string | Date | Dayjs): Dayjs {
  return dayjs.utc(date).local();
}

/**
 * Convert local date to UTC
 */
export function localToUtc(date: string | Date | Dayjs): Dayjs {
  return dayjs(date).utc();
}

/**
 * Format date in a specific timezone
 */
export function formatInTimezone(
  date: string | Date | Dayjs,
  timezone: string,
  format: string = 'YYYY-MM-DD HH:mm:ss'
): string {
  return dayjs.tz(date, timezone).format(format);
}

/**
 * Get timezone offset in minutes for a specific date and timezone
 */
export function getTimezoneOffset(
  date: string | Date | Dayjs, 
  timezone: string
): number {
  return dayjs.tz(date, timezone).utcOffset();
}

/**
 * Check if daylight saving time is active for a timezone on a specific date
 */
export function isDSTActive(
  date: string | Date | Dayjs,
  timezone: string
): boolean {
  const summer = dayjs.tz('2023-07-01', timezone).utcOffset();
  const winter = dayjs.tz('2023-01-01', timezone).utcOffset();
  const current = dayjs.tz(date, timezone).utcOffset();
  
  return current !== winter && summer !== winter;
}

/**
 * Get all available timezone names
 */
export function getAllTimezones(): string[] {
  // This is a simplified list. In a real application, you might want to use
  // a more comprehensive list from a library like `countries-and-timezones`
  return [
    'UTC',
    'America/New_York',
    'America/Chicago',
    'America/Denver',
    'America/Los_Angeles',
    'America/Toronto',
    'America/Vancouver',
    'America/Mexico_City',
    'America/Sao_Paulo',
    'America/Buenos_Aires',
    'Europe/London',
    'Europe/Paris',
    'Europe/Berlin',
    'Europe/Rome',
    'Europe/Madrid',
    'Europe/Amsterdam',
    'Europe/Stockholm',
    'Europe/Warsaw',
    'Europe/Moscow',
    'Europe/Istanbul',
    'Asia/Tokyo',
    'Asia/Shanghai',
    'Asia/Hong_Kong',
    'Asia/Singapore',
    'Asia/Seoul',
    'Asia/Dubai',
    'Asia/Kolkata',
    'Asia/Bangkok',
    'Asia/Manila',
    'Australia/Sydney',
    'Australia/Melbourne',
    'Australia/Perth',
    'Pacific/Auckland',
    'Pacific/Honolulu',
    'Pacific/Fiji'
  ];
}

/**
 * Validate if a timezone name is valid
 */
export function isValidTimezone(timezone: string): boolean {
  try {
    dayjs.tz('2023-01-01', timezone);
    return true;
  } catch {
    return false;
  }
}

/**
 * Get timezone display name for user interface
 */
export function getTimezoneDisplayName(timezone: string): string {
  const now = dayjs().tz(timezone);
  const offset = now.format('Z');
  const name = timezone.replace(/_/g, ' ').replace('/', ' / ');
  
  return `${name} (UTC${offset})`;
}

/**
 * Convert event times to display timezone for calendar
 */
export function convertEventTimezone(
  event: { start_time: string; end_time: string; timezone?: string },
  displayTimezone: string
): { start_time: string; end_time: string } {
  const eventTimezone = event.timezone || getUserTimezone();
  
  const startTime = dayjs.tz(event.start_time, eventTimezone).tz(displayTimezone);
  const endTime = dayjs.tz(event.end_time, eventTimezone).tz(displayTimezone);
  
  return {
    start_time: startTime.toISOString(),
    end_time: endTime.toISOString()
  };
}

/**
 * Get a human-readable timezone difference
 */
export function getTimezoneDifference(
  fromTimezone: string, 
  toTimezone: string,
  date: string | Date | Dayjs = dayjs()
): string {
  const fromOffset = dayjs.tz(date, fromTimezone).utcOffset();
  const toOffset = dayjs.tz(date, toTimezone).utcOffset();
  const diffMinutes = toOffset - fromOffset;
  const diffHours = Math.abs(diffMinutes) / 60;
  
  if (diffMinutes === 0) {
    return 'Same time';
  }
  
  const direction = diffMinutes > 0 ? 'ahead' : 'behind';
  const hours = Math.floor(diffHours);
  const minutes = Math.abs(diffMinutes) % 60;
  
  let result = `${hours}h`;
  if (minutes > 0) {
    result += ` ${minutes}m`;
  }
  result += ` ${direction}`;
  
  return result;
}

export default {
  getUserTimezone,
  getCommonTimezones,
  convertTimezone,
  utcToLocal,
  localToUtc,
  formatInTimezone,
  getTimezoneOffset,
  isDSTActive,
  getAllTimezones,
  isValidTimezone,
  getTimezoneDisplayName,
  convertEventTimezone,
  getTimezoneDifference
};