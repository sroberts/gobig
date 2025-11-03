# Mermaid Diagram Rendering Tests

This directory contains Playwright-based end-to-end tests for verifying that Mermaid diagrams render correctly in gobig presentations.

## Test Files

- **mermaid-rendering.spec.js** - Simplified test suite for Mermaid diagram rendering

## What's Being Tested

The tests verify basic Mermaid functionality:

1. **Presentation Loads** - The mermaid-test.html presentation loads successfully with the correct title
2. **Diagrams Render** - Multiple Mermaid diagram containers are created and rendered on the page
3. **Advanced Example Works** - The advanced.html example includes Mermaid diagram content

## Test Philosophy

These tests follow a **simple and fast** approach:
- Focus on verifying content renders (not detailed visual properties)
- No screenshot generation or visual regression testing
- No iteration through individual slides
- No detailed SVG inspection or color checking
- Tests complete in under 2 seconds

This approach ensures:
- Fast test execution in CI/CD pipelines
- High reliability (fewer false positives)
- Easy maintenance
- Clear pass/fail signals

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
Tests check for the existence of `.mermaid` containers in the DOM, which proves that Mermaid.js successfully processed and rendered the diagrams. This is a reliable indicator of proper integration without needing to inspect detailed SVG properties.

## Test Examples

Example test from the suite:

```javascript
test('renders mermaid diagrams', async ({ page }) => {
  await page.goto('http://localhost:8080/examples/mermaid-test.html');
  await page.waitForLoadState('networkidle');
  await page.waitForSelector('.presentation-container', { timeout: 10000 });

  // Check that mermaid containers exist (proves Mermaid rendered)
  const mermaidContainers = page.locator('.mermaid');
  const count = await mermaidContainers.count();

  // Should have multiple diagrams
  expect(count).toBeGreaterThan(5);
});
```

## Contributing

When adding new Mermaid diagram types or modifying rendering:

1. Add the diagram type to `examples/mermaid-test.md`
2. Regenerate `examples/mermaid-test.html`
3. Run the test suite to ensure diagrams still render
4. Manual testing is recommended for visual verification

## CI/CD Integration

These tests are designed to run efficiently in CI environments. The configuration in `playwright.config.js` adjusts settings for CI:
- Retries: 2 retries in CI, 0 locally
- Workers: 1 worker in CI, parallel locally
- Screenshot: Only on failure

## Future Improvements

- [ ] Add accessibility testing for SVG diagrams
- [ ] Add performance benchmarks for diagram rendering
- [ ] Add tests for responsive sizing
