import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should redirect to login page when not authenticated', async ({ page }) => {
    await expect(page).toHaveURL(/.*login/);
    await expect(page.locator('h1')).toContainText(['Sign In', 'Login', 'Войдите']);
  });

  test('should show registration form', async ({ page }) => {
    // Navigate to registration if not already there
    const signUpButton = page.locator('a[href*="register"], button:has-text("Sign Up"), button:has-text("Регистрация")');
    if (await signUpButton.isVisible()) {
      await signUpButton.click();
    } else {
      await page.goto('/register');
    }

    await expect(page).toHaveURL(/.*register/);
    await expect(page.locator('h1')).toContainText(['Sign Up', 'Register', 'Регистрация']);
    
    // Check form fields exist
    await expect(page.locator('input[type="email"], input[name="email"]')).toBeVisible();
    await expect(page.locator('input[type="password"], input[name="password"]')).toBeVisible();
  });

  test('should validate login form', async ({ page }) => {
    await page.goto('/login');
    
    // Try to submit empty form
    const submitButton = page.locator('button[type="submit"], button:has-text("Sign In"), button:has-text("Войти")').first();
    await submitButton.click();
    
    // Should show validation errors
    await expect(page.locator('text=/email.*required|Email.*обязательно|required/i')).toBeVisible({ timeout: 3000 });
  });

  test('should handle invalid credentials', async ({ page }) => {
    await page.goto('/login');
    
    // Fill invalid credentials
    await page.fill('input[type="email"], input[name="email"]', 'invalid@test.com');
    await page.fill('input[type="password"], input[name="password"]', 'wrongpassword');
    
    const submitButton = page.locator('button[type="submit"], button:has-text("Sign In"), button:has-text("Войти")').first();
    await submitButton.click();
    
    // Should show error message
    await expect(page.locator('text=/invalid.*credentials|неверные.*данные|error/i')).toBeVisible({ timeout: 5000 });
  });
});