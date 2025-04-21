// End-to-End Test Script for Hustler Trading Bot
// This script tests the core functionality of the application

const puppeteer = require('puppeteer');

(async () => {
  console.log('Starting End-to-End Tests for Hustler Trading Bot...');
  
  // Launch browser
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });
  
  const page = await browser.newPage();
  console.log('Browser launched successfully');
  
  try {
    // Test 1: Load the application
    console.log('\nTest 1: Loading application...');
    await page.goto('http://localhost:8000', { waitUntil: 'networkidle2' });
    console.log('✅ Application loaded successfully');
    
    // Test 2: Verify navigation components
    console.log('\nTest 2: Verifying navigation components...');
    const navItems = await page.$$eval('nav a', items => items.length);
    console.log(`Found ${navItems} navigation items`);
    if (navItems >= 5) {
      console.log('✅ Navigation components verified');
    } else {
      throw new Error('Navigation components not found or incomplete');
    }
    
    // Test 3: Test dark mode toggle
    console.log('\nTest 3: Testing dark mode toggle...');
    const initialTheme = await page.evaluate(() => {
      return document.documentElement.classList.contains('dark');
    });
    
    await page.click('button[aria-label="Toggle dark mode"]');
    
    const newTheme = await page.evaluate(() => {
      return document.documentElement.classList.contains('dark');
    });
    
    if (initialTheme !== newTheme) {
      console.log('✅ Dark mode toggle works correctly');
    } else {
      throw new Error('Dark mode toggle failed');
    }
    
    // Test 4: Navigate to different pages
    console.log('\nTest 4: Testing navigation between pages...');
    
    // Dashboard should be visible by default
    let isDashboardVisible = await page.evaluate(() => {
      return document.querySelector('[x-show="currentPage === \'dashboard\'"]').style.display !== 'none';
    });
    
    if (isDashboardVisible) {
      console.log('✅ Dashboard page loaded by default');
    } else {
      throw new Error('Dashboard not visible by default');
    }
    
    // Navigate to Trading Groups
    await page.click('nav a[href="#"]:nth-child(2)');
    await page.waitForTimeout(500); // Wait for Alpine.js to update the DOM
    
    let isGroupsVisible = await page.evaluate(() => {
      return document.querySelector('[x-show="currentPage === \'groups\'"]').style.display !== 'none';
    });
    
    if (isGroupsVisible) {
      console.log('✅ Successfully navigated to Trading Groups page');
    } else {
      throw new Error('Trading Groups page not visible after navigation');
    }
    
    // Navigate to Stock Setup
    await page.click('nav a[href="#"]:nth-child(3)');
    await page.waitForTimeout(500);
    
    let isStocksVisible = await page.evaluate(() => {
      return document.querySelector('[x-show="currentPage === \'stocks\'"]').style.display !== 'none';
    });
    
    if (isStocksVisible) {
      console.log('✅ Successfully navigated to Stock Setup page');
    } else {
      throw new Error('Stock Setup page not visible after navigation');
    }
    
    // Test 5: Verify charts are rendered
    console.log('\nTest 5: Verifying chart rendering...');
    
    // Navigate back to dashboard to check the P&L chart
    await page.click('nav a[href="#"]:nth-child(1)');
    await page.waitForTimeout(500);
    
    const chartElements = await page.$$eval('canvas', canvases => canvases.length);
    
    if (chartElements >= 1) {
      console.log(`Found ${chartElements} chart elements`);
      console.log('✅ Charts rendered successfully');
    } else {
      throw new Error('Charts not rendered properly');
    }
    
    // Test 6: Test mobile responsiveness by resizing viewport
    console.log('\nTest 6: Testing mobile responsiveness...');
    
    // Resize to mobile dimensions
    await page.setViewport({ width: 375, height: 667 });
    await page.waitForTimeout(500);
    
    // Check if mobile menu button is visible
    const mobileMenuVisible = await page.evaluate(() => {
      return window.getComputedStyle(document.querySelector('.sm\\:hidden button')).display !== 'none';
    });
    
    if (mobileMenuVisible) {
      console.log('✅ Mobile menu button visible on small viewport');
    } else {
      throw new Error('Mobile menu button not visible on small viewport');
    }
    
    // Click mobile menu button
    await page.click('.sm\\:hidden button');
    await page.waitForTimeout(500);
    
    // Check if mobile menu is expanded
    const mobileMenuExpanded = await page.evaluate(() => {
      return document.querySelector('[x-show="showMobileMenu"]').style.display !== 'none';
    });
    
    if (mobileMenuExpanded) {
      console.log('✅ Mobile menu expands correctly');
    } else {
      throw new Error('Mobile menu does not expand');
    }
    
    // Reset viewport to desktop size
    await page.setViewport({ width: 1280, height: 800 });
    
    console.log('\nAll tests completed successfully! ✅');
    
  } catch (error) {
    console.error('❌ Test failed:', error);
  } finally {
    await browser.close();
    console.log('\nEnd-to-End tests completed. Browser closed.');
  }
})();
