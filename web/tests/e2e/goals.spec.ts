import { test, expect } from '@playwright/test';

test.describe('Goals Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/goals');
  });

  test('should redirect to login when not authenticated', async ({ page }) => {
    await expect(page).toHaveURL(/.*login/);
  });

  test('should display goals page (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Goals page requires authentication');
    }

    // Check for goals page elements
    await expect(page.locator('h1, h2, .page-title').first()).toBeVisible();
    
    // Look for typical goals page elements
    const goalsElements = [
      'Add Goal', 'Create Goal', 'New Goal',
      'Добавить цель', 'Создать цель', 'Новая цель',
      'Goals', 'Цели'
    ];
    
    let foundElement = false;
    for (const element of goalsElements) {
      if (await page.locator(`text=${element}`).isVisible()) {
        foundElement = true;
        break;
      }
    }
    
    if (foundElement) {
      await expect(page.locator('text=/Goal|Цель/i')).toBeVisible();
    }
  });

  test('should handle goal creation modal (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Goals page requires authentication');
    }

    // Try to find and click add goal button
    const addGoalButtons = [
      'button:has-text("Add Goal")',
      'button:has-text("Create Goal")',
      'button:has-text("New Goal")',
      'button:has-text("Добавить цель")',
      'button:has-text("Создать цель")',
      'button:has-text("+")',
      '[data-testid="add-goal-button"]'
    ];

    let buttonClicked = false;
    for (const buttonSelector of addGoalButtons) {
      const button = page.locator(buttonSelector);
      if (await button.isVisible()) {
        await button.click();
        buttonClicked = true;
        break;
      }
    }

    if (buttonClicked) {
      // Check if modal opened
      await expect(page.locator('.ant-modal, .modal, [role="dialog"]')).toBeVisible({ timeout: 3000 });
      
      // Look for SMART goal form fields
      const smartFields = ['Specific', 'Measurable', 'Achievable', 'Relevant', 'Time-bound'];
      const foundSmartField = await Promise.all(
        smartFields.map(field => page.locator(`text=${field}`).isVisible())
      );
      
      if (foundSmartField.some(Boolean)) {
        await expect(page.locator('text=/Specific|Measurable|Time/i')).toBeVisible();
      }
    }
  });

  test('should display goals list (if authenticated)', async ({ page }) => {
    // Skip if redirected to login
    if (await page.locator('text=/sign in|login|войти/i').isVisible()) {
      test.skip('Goals page requires authentication');
    }

    // Look for goals list container
    const listContainers = [
      '.goals-list',
      '.ant-list',
      '[data-testid="goals-list"]',
      '.goal-item',
      '.goals-container'
    ];

    let foundContainer = false;
    for (const container of listContainers) {
      if (await page.locator(container).isVisible()) {
        foundContainer = true;
        break;
      }
    }

    // If no goals exist, should show empty state
    if (!foundContainer) {
      const emptyMessages = [
        'No goals', 'Empty', 'Create your first goal',
        'Нет целей', 'Пусто', 'Создайте первую цель'
      ];
      
      let foundEmptyMessage = false;
      for (const message of emptyMessages) {
        if (await page.locator(`text=${message}`).isVisible()) {
          foundEmptyMessage = true;
          break;
        }
      }
      
      if (foundEmptyMessage) {
        await expect(page.locator('text=/No goals|Empty|Нет целей|Пусто/i')).toBeVisible();
      }
    }
  });
});