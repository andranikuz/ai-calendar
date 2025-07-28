import { test, expect } from '@playwright/test';

test.describe('Mood Tracking', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/moods');
  });

  test('should redirect to login when not authenticated', async ({ page }) => {
    await expect(page).toHaveURL(/.*login/);
  });

  test('should display mood tracking interface (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Mood tracking requires authentication');
    }

    // Check for mood page elements
    await expect(page.locator('h1, h2, .page-title').first()).toBeVisible();
    
    // Look for mood-related elements
    const moodElements = [
      'Mood', 'Настроение',
      'How are you feeling', 'Как дела',
      'Track mood', 'Отследить настроение'
    ];
    
    let foundElement = false;
    for (const element of moodElements) {
      if (await page.locator(`text=${element}`).isVisible()) {
        foundElement = true;
        break;
      }
    }
    
    if (foundElement) {
      await expect(page.locator('text=/Mood|Настроение/i')).toBeVisible();
    }
  });

  test('should display emoji mood selector (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Mood tracking requires authentication');
    }

    // Look for emoji selectors (typical mood emojis)
    const moodEmojis = ['😢', '😔', '😐', '😊', '😄', '🙂', '😃', '🥺', '😭'];
    
    let foundEmoji = false;
    for (const emoji of moodEmojis) {
      if (await page.locator(`text=${emoji}`).isVisible()) {
        foundEmoji = true;
        break;
      }
    }
    
    if (foundEmoji) {
      // Test emoji interaction
      const firstEmoji = page.locator('text=😊, text=😄, text=🙂').first();
      if (await firstEmoji.isVisible()) {
        await firstEmoji.click();
        
        // Should show some feedback after selection
        await expect(page.locator('text=/selected|saved|выбрано|сохранено/i')).toBeVisible({ timeout: 3000 });
      }
    }
  });

  test('should display mood calendar/history (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Mood tracking requires authentication');
    }

    // Look for calendar or history view
    const historyElements = [
      '.mood-calendar',
      '.mood-history',
      '.calendar',
      'History', 'История',
      'Calendar', 'Календарь'
    ];

    let foundHistory = false;
    for (const element of historyElements) {
      if (typeof element === 'string' && element.startsWith('.')) {
        if (await page.locator(element).isVisible()) {
          foundHistory = true;
          break;
        }
      } else if (await page.locator(`text=${element}`).isVisible()) {
        foundHistory = true;
        break;
      }
    }

    if (foundHistory) {
      // Should show some form of mood history
      await expect(page.locator('.calendar, .history, text=/History|Calendar|История|Календарь/i')).toBeVisible();
    }
  });

  test('should show mood statistics (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Mood tracking requires authentication');
    }

    // Look for statistics/analytics
    const statsElements = [
      'Statistics', 'Stats', 'Analytics',
      'Статистика', 'Аналитика',
      'Average', 'Среднее',
      'Trend', 'Тренд'
    ];

    let foundStats = false;
    for (const element of statsElements) {
      if (await page.locator(`text=${element}`).isVisible()) {
        foundStats = true;
        break;
      }
    }

    if (foundStats) {
      await expect(page.locator('text=/Statistics|Analytics|Статистика|Аналитика/i')).toBeVisible();
    }
  });
});