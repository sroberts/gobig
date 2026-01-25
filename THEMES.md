# GoBig Themes

GoBig supports customizable themes defined in YAML format. Themes control colors, typography, and other visual aspects of your presentations.

## Built-in Themes

GoBig comes with several built-in themes:

- **dark** (default) - High-contrast dark theme with white text
- **light** - Light gray background with dark text
- **white** - Pure white background with bold black text
- **solarized-dark** - Based on the Solarized color scheme
- **oceanic** - Cool blue tones inspired by the ocean

## Using Themes

Specify a theme with the `-theme` flag:

```bash
gobig -theme light -o presentation.html slides.md
gobig -serve -theme solarized-dark presentation.md
```

## Creating Custom Themes

Themes are defined in YAML files located in `internal/assets/embed/themes/`. You can create your own theme by creating a new YAML file.

### Theme File Structure

Here's a complete example of a theme file:

```yaml
# Theme metadata
name: my-theme
description: My custom presentation theme

# Colors
colors:
  background: "#1a1a1a"           # Slide background color
  foreground: "#ffffff"           # Default text color
  emphasis: "#ffcc00"             # Color for emphasized text (em tags)
  link: "#4a9eff"                 # Link color
  link_hover_background: "#333333" # Background color when hovering over links

# Typography
typography:
  # Font family - supports CSS font stacks
  font_family: "-apple-system, BlinkMacSystemFont, avenir next, avenir, helvetica neue, helvetica, Ubuntu, roboto, noto, segoe ui, arial, sans-serif"
  # Font weight for body text (100-900)
  font_weight: 700

# Link styling
links:
  text_decoration_skip: "ink"     # How text decoration skips descenders

# Advanced: Custom CSS
# Add any additional CSS rules here
custom_css: |
  h1 {
    color: #ffcc00;
    text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
  }

  code {
    background: #2a2a2a;
    padding: 2px 6px;
    border-radius: 3px;
  }
```

### Theme Properties

#### Colors

All color properties are optional. If omitted, they won't generate corresponding CSS.

- `background` - Slide background color
- `foreground` - Default text color
- `emphasis` - Color for `<em>` tags (italicized text in Markdown)
- `link` - Color for hyperlinks
- `link_hover_background` - Background color when hovering over links

#### Typography

- `font_family` - CSS font stack for the presentation
- `font_weight` - Font weight for body text (100-900)

#### Links

- `text_decoration_skip` - How underlines skip descenders (typically "ink")

#### Custom CSS

The `custom_css` field allows you to add arbitrary CSS rules that will be appended to the generated theme CSS. This is useful for:

- Styling specific elements (h1, h2, code blocks, etc.)
- Adding text shadows or other effects
- Overriding default styles
- Adding responsive design rules

### Installing a Custom Theme

1. Create a new YAML file in `internal/assets/embed/themes/`
2. Name it `mytheme.yaml` (where "mytheme" is your theme name)
3. Define your theme properties following the structure above
4. Rebuild GoBig: `go build ./cmd/gobig`
5. Use your theme: `gobig -theme mytheme presentation.md`

### Example: Creating a "Midnight" Theme

Create `internal/assets/embed/themes/midnight.yaml`:

```yaml
name: midnight
description: Deep blue midnight theme

colors:
  background: "#0a0e27"
  foreground: "#e0e6f0"
  emphasis: "#ffd700"
  link: "#66b3ff"
  link_hover_background: "#1a2847"

typography:
  font_family: "-apple-system, BlinkMacSystemFont, avenir next, avenir, helvetica neue, helvetica, Ubuntu, roboto, noto, segoe ui, arial, sans-serif"
  font_weight: 600

links:
  text_decoration_skip: "ink"

custom_css: |
  h1, h2, h3 {
    color: #66d9ff;
  }
```

Rebuild and use:

```bash
go build ./cmd/gobig
gobig -theme midnight -o presentation.html slides.md
```

## Color Scheme Resources

When creating custom themes, consider using these resources:

- [Coolors](https://coolors.co/) - Color palette generator
- [Adobe Color](https://color.adobe.com/) - Color wheel and schemes
- [Solarized](https://ethanschoonover.com/solarized/) - Popular color scheme
- [Nord](https://www.nordtheme.com/) - Arctic-inspired color palette
- [Dracula](https://draculatheme.com/) - Dark theme

## Accessibility

When creating themes, consider accessibility:

- Ensure sufficient contrast between background and foreground colors
- Test your theme with color blindness simulators
- Use WCAG 2.1 Level AA contrast ratios (4.5:1 for normal text)
- Test readability in different lighting conditions

## Tips

1. **Start with an existing theme** - Copy a built-in theme and modify it
2. **Test thoroughly** - View your presentations on different displays
3. **Keep it simple** - Less is often more in presentation design
4. **Consider your audience** - Dark themes work well in dark rooms, light themes for bright environments
5. **Use semantic colors** - Make emphasis and links distinct from body text
6. **Font weights matter** - Heavier weights work better for projectors
