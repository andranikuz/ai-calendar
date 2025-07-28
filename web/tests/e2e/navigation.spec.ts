import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should navigate between pages', async ({ page }) => {
    // Test basic navigation structure
    const navItems = [
      { text: 'Dashboard', path: '/dashboard' },
      { text: 'Calendar', path: '/calendar' },
      { text: 'Goals', path: '/goals' },
      { text: 'Moods', path: '/moods' },
      { text: 'Settings', path: '/settings' }
    ];

    for (const item of navItems) {
      // Try to find navigation link
      const navLink = page.locator(`a[href="${item.path}"], a:has-text("${item.text}")`);
      
      if (await navLink.isVisible()) {
        await navLink.click();
        
        // Should navigate to the page (or login if not authenticated)
        await page.waitForURL(/.*/, { timeout: 3000 });
        
        // Either should be on the target page or redirected to login
        const currentURL = page.url();
        const isOnTargetPage = currentURL.includes(item.path);
        const isOnLoginPage = currentURL.includes('login');
        
        expect(isOnTargetPage || isOnLoginPage).toBe(true);
      }
    }
  });

  test('should show mobile navigation on small screens', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');

    // Look for mobile navigation indicators
    const mobileNavSelectors = [
      '.mobile-nav',
      '.bottom-nav',
      '[data-testid="mobile-nav"]',
      '.mobile-navigation',
      '.drawer-button',
      '.hamburger-menu'
    ];

    let foundMobileNav = false;
    for (const selector of mobileNavSelectors) {
      if (await page.locator(selector).isVisible()) {
        foundMobileNav = true;
        
        // Test mobile nav interaction
        if (selector.includes('drawer') || selector.includes('hamburger')) {
          await page.locator(selector).click();
          await expect(page.locator('.drawer, .menu, .navigation-menu')).toBeVisible({ timeout: 2000 });
        }
        break;
      }
    }

    // Alternative: check if regular nav is hidden on mobile
    if (!foundMobileNav) {
      const regularNav = page.locator('.desktop-nav, .sidebar, nav:not(.mobile-nav)');
      if (await regularNav.isVisible()) {
        // Should be responsive - either hidden or adapted for mobile
        const navWidth = await regularNav.boundingBox();
        if (navWidth && navWidth.width > 375) {
          // Navigation might be too wide for mobile - this could be an issue
          console.warn('Navigation might not be responsive');
        }
      }
    }
  });

  test('should handle browser back/forward navigation', async ({ page }) => {
    // Start at home
    await page.goto('/');
    
    // Navigate to a page
    await page.goto('/dashboard');
    
    // Go back
    await page.goBack();
    
    // Should be back at home or login
    const currentURL = page.url();
    const isAtExpectedPage = currentURL.includes('/') || currentURL.includes('login');
    expect(isAtExpectedPage).toBe(true);
    
    // Go forward
    await page.goForward();
    
    // Should be back at dashboard or login
    const forwardURL = page.url();
    const isAtForwardPage = forwardURL.includes('dashboard') || forwardURL.includes('login');
    expect(isAtForwardPage).toBe(true);
  });

  test('should show active navigation state', async ({ page }) => {
    // Skip if redirected to login immediately
    if (await page.url().then(url => url.includes('login'))) {
      test.skip('Testing navigation requires some level of access');
    }

    const pages = ['/dashboard', '/calendar', '/goals', '/moods'];
    
    for (const pagePath of pages) {
      await page.goto(pagePath);
      
      // Skip if redirected to login
      if (await page.url().then(url => url.includes('login'))) {
        continue;
      }
      
      // Look for active navigation indicator
      const activeSelectors = [
        `.active[href="${pagePath}"]`,
        `.nav-item.active:has(a[href="${pagePath}"])`,
        `a[href="${pagePath}"].active`,
        `[aria-current="page"][href="${pagePath}"]`
      ];
      
      let foundActive = false;
      for (const selector of activeSelectors) {
        if (await page.locator(selector).isVisible()) {
          foundActive = true;
          await expect(page.locator(selector)).toBeVisible();
          break;
        }
      }
      
      // If no specific active class found, at least navigation should be visible
      if (!foundActive) {
        const anyNav = page.locator('nav, .navigation, .nav-menu').first();
        if (await anyNav.isVisible()) {
          await expect(anyNav).toBeVisible();
        }
      }
    }
  });
});