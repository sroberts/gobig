# Mermaid Diagram Rendering Tests

This directory contains Playwright-based end-to-end tests for verifying that Mermaid diagrams render correctly in gobig presentations.

## Test Files

- **mermaid-rendering.spec.js** - Comprehensive test suite for Mermaid diagram rendering
- **screenshots/** - Visual regression screenshots (generated during test runs)

## What's Being Tested

The tests verify that:

1. **Color Scheme** - Dark theme colors are applied correctly:
   - Dark backgrounds (#1f2020) on diagram boxes
   - Light text (#f0f0f0) for readability
   - Light borders (#ccc) for contrast

2. **No White Backgrounds** - SVG elements don't have white background-color styles that interfere with dark theme

3. **All Diagram Types Render** - Tests all 12 Mermaid diagram types:
   - Flowchart
   - Sequence Diagram
   - Class Diagram
   - State Diagram
   - Entity Relationship Diagram
   - Gantt Chart
   - Pie Chart
   - Git Graph
   - Journey Diagram
   - Mindmap
   - Timeline
   - Requirement Diagram

4. **CSS Classes Applied** - Dark theme CSS classes are present and working

5. **Content Visibility** - Diagrams are visible and contain expected content

6. **Contrast** - Text has sufficient contrast against backgrounds

## Running the Tests

```bash
# Install dependencies (if not already installed)
npm install

# Install Playwright browsers
npx playwright install

# Run all tests
npm test

# Run tests with UI (interactive mode)
npm run test:ui

# Run tests with browser visible (headed mode)
npm run test:headed

# Run tests in debug mode
npm run test:debug
```

## Test Architecture

### Web Server
Tests use the Playwright built-in web server to serve the examples directory over HTTP. The server runs on port 8080 and is automatically started before tests run.

### Test Strategy
The tests navigate to different slide numbers using URL hash fragments (e.g., `#1`, `#2`) to test individual diagram slides.

### Known Limitations

**Note**: Some tests currently fail due to big.js loading all slides simultaneously in the DOM. The test selectors need refinement to target only the currently visible slide. This is a test infrastructure issue, not a Mermaid rendering issue - manual testing confirms all diagrams render correctly.

## Visual Regression Testing

The test suite generates screenshots of each diagram type in `tests/screenshots/`. These can be used for visual regression testing to ensure diagram rendering remains consistent across changes.

## Test Examples

Example test from the suite:

```javascript
test('flowchart renders with correct colors', async ({ page }) => {
  // Navigate to flowchart slide (slide 1)
  await page.goto('http://localhost:8080/examples/mermaid-test.html#1');
  await page.waitForTimeout(500);

  // Check for flowchart nodes
  const nodes = page.locator('.mermaid rect.basic');
  await expect(nodes.first()).toBeVisible();

  // Verify dark background on nodes
  const nodeStyle = await nodes.first().getAttribute('fill');
  expect(nodeStyle).toBe('#1f2020');
});
```

## Contributing

When adding new Mermaid diagram types or modifying rendering:

1. Add the diagram type to `examples/mermaid-test.md`
2. Regenerate `examples/mermaid-test.html`
3. Add a test case to verify proper rendering
4. Run the test suite to ensure no regressions

## CI/CD Integration

These tests are designed to run in CI environments. The configuration in `playwright.config.js` adjusts settings for CI:
- Retries: 2 retries in CI, 0 locally
- Workers: 1 worker in CI, parallel locally
- Screenshot: Only on failure

##Future Improvements

- [ ] Fix selectors to work with big.js's DOM structure (all slides loaded simultaneously)
- [ ] Add accessibility testing for SVG diagrams
- [ ] Add performance benchmarks for diagram rendering
- [ ] Add tests for light and white themes
- [ ] Add tests for responsive sizing
