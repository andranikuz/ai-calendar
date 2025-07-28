import { test, expect } from '@playwright/test';

test.describe('Calendar', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/calendar');
  });

  test('should redirect to login when not authenticated', async ({ page }) => {
    await expect(page).toHaveURL(/.*login/);
  });

  test('should display calendar interface (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Calendar requires authentication');
    }

    // Test calendar components
    await expect(page.locator('.fc-toolbar, .calendar-header, h1').first()).toBeVisible({ timeout: 5000 });
    
    // Check for calendar view controls
    const viewControls = [
      'Month', 'Week', 'Day',
      'Месяц', 'Неделя', 'День',
      'month', 'week', 'day'
    ];
    
    let foundViewControl = false;
    for (const control of viewControls) {
      if (await page.locator(`text=${control}`).isVisible()) {
        foundViewControl = true;
        break;
      }
    }
    
    if (foundViewControl) {
      await expect(page.locator('text=/Month|Week|Day|Месяц|Неделя|День/i')).toBeVisible();
    }
  });

  test('should handle calendar navigation', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Calendar requires authentication');
    }

    // Test navigation buttons
    const prevButton = page.locator('button:has-text("❮"), button:has-text("Previous"), button:has-text("Prev"), .fc-prev-button');
    const nextButton = page.locator('button:has-text("❯"), button:has-text("Next"), .fc-next-button');
    
    if (await prevButton.isVisible()) {
      await prevButton.click();
      // Wait for navigation to complete
      await page.waitForTimeout(500);
    }
    
    if (await nextButton.isVisible()) {
      await nextButton.click();
      await page.waitForTimeout(500);
    }
  });

  test('should be responsive on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войti/i').isVisible()) {
      test.skip('Calendar requires authentication');
    }

    // Calendar should adapt to mobile viewport
    await expect(page.locator('body')).toHaveCSS('width', '375px');
  });
});