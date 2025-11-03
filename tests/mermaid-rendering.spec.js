// @ts-check
const { test, expect } = require('@playwright/test');

test.describe('Mermaid Diagram Rendering', () => {
  test('loads mermaid test presentation', async ({ page }) => {
    await page.goto('http://localhost:8080/examples/mermaid-test.html');
    await page.waitForLoadState('networkidle');

    // Verify title
    await expect(page).toHaveTitle(/Mermaid Diagram Test Suite/);

    // Wait for big.js to initialize
    await page.waitForSelector('.presentation-container', { timeout: 10000 });
  });

  test('renders mermaid diagrams', async ({ page }) => {
    await page.goto('http://localhost:8080/examples/mermaid-test.html');
    await page.waitForLoadState('networkidle');
    await page.waitForSelector('.presentation-container', { timeout: 10000 });

    // Wait for mermaid diagrams to be rendered (looking for SVG elements inside mermaid containers)
    // Use state: 'attached' instead of default 'visible' since not all slides are visible
    await page.waitForSelector('.mermaid svg', { timeout: 10000, state: 'attached' });

    // Check that mermaid containers exist (proves Mermaid rendered)
    const mermaidContainers = page.locator('.mermaid');
    const count = await mermaidContainers.count();

    // Should have multiple diagrams
    expect(count).toBeGreaterThan(5);
  });

  test('renders mermaid in advanced example', async ({ page }) => {
    await page.goto('http://localhost:8080/examples/advanced.html');
    await page.waitForLoadState('networkidle');
    await page.waitForSelector('.presentation-container', { timeout: 10000 });

    // Wait for mermaid diagrams to be rendered (looking for SVG elements inside mermaid containers)
    // Use state: 'attached' instead of default 'visible' since not all slides are visible
    await page.waitForSelector('.mermaid svg', { timeout: 10000, state: 'attached' });

    // Check that mermaid content exists
    const mermaidContainers = page.locator('.mermaid');
    const count = await mermaidContainers.count();
    expect(count).toBeGreaterThan(0);
  });
});
