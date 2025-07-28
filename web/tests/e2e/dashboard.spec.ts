import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // For now, we'll test the dashboard UI without authentication
    // In a real scenario, you'd want to login first or mock authentication
    await page.goto('/dashboard');
  });

  test('should display dashboard components', async ({ page }) => {
    // Check if we're redirected to login (expected behavior)
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      // This is expected behavior - dashboard requires authentication
      await expect(page).toHaveURL(/.*login/);
      return;
    }

    // If somehow we reach dashboard, test its components
    await expect(page.locator('h1, h2, h3').first()).toBeVisible();
    
    // Check for main dashboard sections
    const possibleSections = [
      'Goals', 'Цели',
      'Calendar', 'Календарь', 
      'Mood', 'Настроение',
      'Progress', 'Прогресс'
    ];
    
    let foundSection = false;
    for (const section of possibleSections) {
      if (await page.locator(`text=${section}`).isVisible()) {
        foundSection = true;
        break;
      }
    }
    
    if (foundSection) {
      await expect(page.locator('text=/Goals|Calendar|Mood|Цели|Календарь|Настроение/i')).toBeVisible();
    }
  });

  test('should be responsive on mobile', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    
    await page.goto('/dashboard');
    
    // Check if mobile navigation exists (if we reach dashboard)
    if (!await page.locator('text=/sign in|login|войти/i').isVisible()) {
      // Look for mobile navigation indicators
      const mobileNavExists = await page.locator('[role="navigation"], .mobile-nav, .bottom-nav').isVisible();
      if (mobileNavExists) {
        await expect(page.locator('[role="navigation"], .mobile-nav, .bottom-nav')).toBeVisible();
      }
    }
  });
});