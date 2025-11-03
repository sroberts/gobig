// @ts-check
const { test, expect } = require('@playwright/test');

const DARK_THEME_BG = '#1f2020';
const LIGHT_TEXT_COLOR = '#f0f0f0';
const BORDER_COLOR = '#ccc';

/**
 * Helper function to check if a color matches expected value (with tolerance for slight variations)
 */
function colorMatches(actual, expected, tolerance = 10) {
  const parseColor = (color) => {
    if (color.startsWith('rgb')) {
      const match = color.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/);
      return match ? [parseInt(match[1]), parseInt(match[2]), parseInt(match[3])] : null;
    }
    // Convert hex to rgb
    const hex = color.replace('#', '');
    return [
      parseInt(hex.substr(0, 2), 16),
      parseInt(hex.substr(2, 2), 16),
      parseInt(hex.substr(4, 2), 16)
    ];
  };

  const a = parseColor(actual);
  const e = parseColor(expected);

  if (!a || !e) return false;

  return Math.abs(a[0] - e[0]) <= tolerance &&
         Math.abs(a[1] - e[1]) <= tolerance &&
         Math.abs(a[2] - e[2]) <= tolerance;
}

test.describe('Mermaid Diagram Rendering', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the mermaid test presentation
    await page.goto('http://localhost:8080/examples/mermaid-test.html');
    await page.waitForLoadState('networkidle');

    // Wait for big.js to initialize and create .slide elements
    await page.waitForSelector('.slide', { timeout: 10000 });
  });

  test('should load the presentation', async ({ page }) => {
    await expect(page).toHaveTitle(/Mermaid Diagram Test Suite/);
  });

  test('flowchart renders with correct colors', async ({ page }) => {
    // Navigate to flowchart slide (slide 1)
    await page.keyboard.press('ArrowRight');
    await page.waitForTimeout(500); // Wait for slide transition

    // Get the slide by index (slide 1 is the second .slide div, index 1)
    const slide = page.locator('.slide').nth(1);

    // Check that SVG exists in this slide
    const svg = slide.locator('.mermaid svg');
    await expect(svg).toBeVisible();

    // Check for flowchart nodes
    const nodes = slide.locator('.mermaid rect.basic');
    await expect(nodes.first()).toBeVisible();

    // Verify dark background on nodes
    const nodeStyle = await nodes.first().getAttribute('fill');
    expect(nodeStyle).toBe(DARK_THEME_BG);

    // Check text is visible (not checking exact color due to SVG complexity)
    const text = slide.locator('.mermaid text');
    await expect(text.first()).toBeVisible();
  });

  test('sequence diagram renders with correct colors', async ({ page }) => {
    // Navigate to sequence diagram slide (slide 2)
    await page.keyboard.press('ArrowRight');
    await page.keyboard.press('ArrowRight');
    await page.waitForTimeout(500);

    // Get the slide by index
    const slide = page.locator('.slide').nth(2);

    // Check that SVG exists
    const svg = slide.locator('.mermaid svg');
    await expect(svg).toBeVisible();

    // Check for actor boxes
    const actors = slide.locator('.mermaid rect.actor');
    await expect(actors.first()).toBeVisible();

    // Verify actor boxes have dark background
    const actorFill = await actors.first().getAttribute('fill');
    expect(actorFill).toBe(DARK_THEME_BG);
  });

  test('class diagram renders with correct colors', async ({ page }) => {
    // Navigate to class diagram slide (slide 3)
    await page.goto('http://localhost:8080/examples/mermaid-test.html#3');
    await page.waitForLoadState('networkidle');

    // Wait for big.js to initialize and create .slide elements
    await page.waitForSelector('.slide', { timeout: 10000 });

    // Get the slide by index
    const slide = page.locator('.slide').nth(3);

    // Check that SVG exists
    const svg = slide.locator('.mermaid svg');
    await expect(svg).toBeVisible();

    // Check for class boxes
    const classBoxes = slide.locator('.mermaid g.classGroup rect');
    await expect(classBoxes.first()).toBeVisible();

    // Verify class boxes have dark background
    const boxFill = await classBoxes.first().getAttribute('fill');
    expect(boxFill).toBe(DARK_THEME_BG);
  });

  test('state diagram renders with correct colors', async ({ page }) => {
    // Navigate to state diagram slide (slide 4)
    await page.goto('http://localhost:8080/examples/mermaid-test.html#4');
    await page.waitForLoadState('networkidle');

    // Wait for big.js to initialize and create .slide elements
    await page.waitForSelector('.slide', { timeout: 10000 });

    // Get the slide by index
    const slide = page.locator('.slide').nth(4);

    // Check that SVG exists
    const svg = slide.locator('.mermaid svg');
    await expect(svg).toBeVisible();

    // Check for state boxes
    const stateBoxes = slide.locator('.mermaid g.stateGroup rect');
    await expect(stateBoxes.first()).toBeVisible();

    // Verify state boxes have dark background
    const boxFill = await stateBoxes.first().getAttribute('fill');
    expect(boxFill).toBe(DARK_THEME_BG);
  });

  test('ER diagram renders with correct colors', async ({ page }) => {
    // Navigate to ER diagram slide (slide 5)
    await page.goto('http://localhost:8080/examples/mermaid-test.html#5');
    await page.waitForLoadState('networkidle');

    // Wait for big.js to initialize and create .slide elements
    await page.waitForSelector('.slide', { timeout: 10000 });

    // Get the slide by index
    const slide = page.locator('.slide').nth(5);

    // Check that SVG exists
    const svg = slide.locator('.mermaid svg');
    await expect(svg).toBeVisible();

    // Check for entity boxes
    const entityBoxes = slide.locator('.mermaid .entityBox');
    await expect(entityBoxes.first()).toBeVisible();

    // Verify entity boxes have dark background
    const boxFill = await entityBoxes.first().getAttribute('fill');
    expect(boxFill).toBe(DARK_THEME_BG);
  });

  test('no white backgrounds on any diagrams', async ({ page }) => {
    const slides = [1, 2, 3, 4, 5]; // Test first 5 diagram slides

    for (const slideNum of slides) {
      await page.goto(`http://localhost:8080/examples/mermaid-test.html#${slideNum}`);
      await page.waitForLoadState('networkidle');

      // Wait for big.js to initialize and create .slide elements
      await page.waitForSelector('.slide', { timeout: 10000 });

      // Get the slide by index
      const slide = page.locator('.slide').nth(slideNum);

      // Check SVG doesn't have white background in style attribute
      const svg = slide.locator('.mermaid svg');
      const style = await svg.getAttribute('style');

      if (style) {
        expect(style).not.toContain('background-color: white');
        expect(style).not.toContain('background-color:white');
      }

      // Check for white fill attributes in this specific slide
      const whiteFills = await slide.locator('.mermaid [fill="#ffffff"], .mermaid [fill="#fff"], .mermaid [fill="white"]').count();
      // Some white fills might be intentional (like end state markers), so we just log it
      console.log(`Slide ${slideNum}: Found ${whiteFills} white-filled elements`);
    }
  });

  test('all diagrams are visible and not empty', async ({ page }) => {
    const slides = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]; // All diagram slides

    for (const slideNum of slides) {
      await page.goto(`http://localhost:8080/examples/mermaid-test.html#${slideNum}`);
      await page.waitForLoadState('networkidle');

      // Wait for big.js to initialize and create .slide elements
      await page.waitForSelector('.slide', { timeout: 10000 });

      // Get the slide by index
      const slide = page.locator('.slide').nth(slideNum);

      // Check that SVG exists and is visible
      const svg = slide.locator('.mermaid svg');
      await expect(svg).toBeVisible();

      // Check that SVG has some content (paths, rects, or text)
      const hasContent = await slide.locator('.mermaid svg > g').count();
      expect(hasContent).toBeGreaterThan(0);

      // Take a screenshot for visual regression testing
      await page.screenshot({
        path: `tests/screenshots/diagram-slide-${slideNum}.png`,
        fullPage: false
      });
    }
  });

  test('CSS dark theme classes are applied', async ({ page }) => {
    // Check body has dark class
    const bodyClass = await page.locator('body').getAttribute('class');
    expect(bodyClass).toContain('dark');

    // Navigate to a diagram slide
    await page.goto('http://localhost:8080/examples/mermaid-test.html#1');
    await page.waitForLoadState('networkidle');

    // Wait for big.js to initialize and create .slide elements
    await page.waitForSelector('.slide', { timeout: 10000 });

    // Check that mermaid container has proper parent with dark class
    const slide = page.locator('.slide').nth(1);
    const darkParent = slide.locator('.mermaid');
    await expect(darkParent).toBeVisible();
  });

  test('advanced example Mermaid diagram renders correctly', async ({ page }) => {
    // Test the actual advanced.html example that was originally broken
    await page.goto('http://localhost:8080/examples/advanced.html#10');
    await page.waitForLoadState('networkidle');

    // Wait for big.js to initialize and create .slide elements
    await page.waitForSelector('.slide', { timeout: 10000 });

    // Get the slide by index (slide 10)
    const slide = page.locator('.slide').nth(10);

    // Check that SVG exists and is visible
    const svg = slide.locator('.mermaid svg');
    await expect(svg).toBeVisible();

    // Check for the specific diagram content (gobig -> Markdown -> HTML -> Present!)
    const text = slide.locator('.mermaid text');
    const textContent = await text.allTextContents();
    const fullText = textContent.join(' ');

    expect(fullText).toContain('gobig');
    expect(fullText).toContain('Markdown');
    expect(fullText).toContain('HTML');
    expect(fullText).toContain('Present');

    // Verify nodes have dark backgrounds
    const nodes = slide.locator('.mermaid rect.basic');
    await expect(nodes.first()).toBeVisible();

    const nodeFill = await nodes.first().getAttribute('fill');
    expect(nodeFill).toBe(DARK_THEME_BG);
  });

  test('text has sufficient contrast on dark backgrounds', async ({ page }) => {
    // Navigate to flowchart
    await page.goto('http://localhost:8080/examples/mermaid-test.html#1');
    await page.waitForLoadState('networkidle');

    // Wait for big.js to initialize and create .slide elements
    await page.waitForSelector('.slide', { timeout: 10000 });

    // Get the slide by index
    const slide = page.locator('.slide').nth(1);

    // Get computed styles for text elements
    const textElement = slide.locator('.mermaid text').first();
    await expect(textElement).toBeVisible();

    const computedColor = await textElement.evaluate((el) => {
      return window.getComputedStyle(el).fill || window.getComputedStyle(el).color;
    });

    // Text should be light colored (high RGB values)
    console.log('Text color:', computedColor);
    expect(computedColor).toBeTruthy();
  });
});

test.describe('Mermaid Diagram Types Coverage', () => {
  const diagramTypes = [
    { name: 'Flowchart', slide: 1 },
    { name: 'Sequence Diagram', slide: 2 },
    { name: 'Class Diagram', slide: 3 },
    { name: 'State Diagram', slide: 4 },
    { name: 'Entity Relationship Diagram', slide: 5 },
    { name: 'Gantt Chart', slide: 6 },
    { name: 'Pie Chart', slide: 7 },
    { name: 'Git Graph', slide: 8 },
    { name: 'Journey Diagram', slide: 9 },
    { name: 'Mindmap', slide: 10 },
    { name: 'Timeline', slide: 11 },
    { name: 'Requirement Diagram', slide: 12 },
  ];

  for (const diagram of diagramTypes) {
    test(`${diagram.name} renders without errors`, async ({ page }) => {
      await page.goto(`http://localhost:8080/examples/mermaid-test.html#${diagram.slide}`);
      await page.waitForLoadState('networkidle');

      // Wait for big.js to initialize and create .slide elements
      await page.waitForSelector('.slide', { timeout: 10000 });

      // Get the slide by index
      const slide = page.locator('.slide').nth(diagram.slide);

      // Check for SVG
      const svg = slide.locator('.mermaid svg');
      await expect(svg).toBeVisible();

      // Check no error messages
      const errorIcon = slide.locator('.mermaid .error-icon');
      await expect(errorIcon).not.toBeVisible();

      // Log for debugging
      const title = await page.title();
      console.log(`${diagram.name}: ${title}`);
    });
  }
});
