autoscale: true
slidenumbers: true
footer: gobig + DeckSet = ❤️
theme: dark

---

# [fit] Welcome to gobig
## With DeckSet Format Support

^ Welcome everyone to this presentation
^ We're demonstrating gobig's DeckSet format compatibility

---

## What is DeckSet?

DeckSet is a popular **macOS presentation tool** that uses Markdown.

- Clean, minimal syntax
- Focus on content, not design
- Markdown-based workflow
- Professional results

^ DeckSet has been popular among developers and presenters
^ Now you can use the same format with gobig

---

[.autoscale: true]

## Key Features

1. **Global Configuration** - Set defaults at file top
2. **Per-Slide Directives** - Fine-tune individual slides
3. **Speaker Notes** - Use `^` prefix for notes
4. **Image Modifiers** - Position images with keywords
5. **Fit Headers** - Scale text to fill slides
6. **Backward Compatible** - Works with gobig syntax too

^ This slide has autoscale enabled via directive
^ Notice how all the content fits nicely

---

# [fit] Big Ideas
# [fit] Deserve
# [fit] Big Text

^ The [fit] modifier scales headers to fill the slide
^ Perfect for title slides or emphasis

---

![left](https://pixnio.com/free-images/fauna-animals/dogs/puppy-in-leaves.jpg)

## Image Positioning

DeckSet's image modifiers make layout easy:

- `![left](img.jpg)` - Image on left
- `![right](img.jpg)` - Image on right
- `![fit](img.jpg)` - Scale to fit
- `![inline](img.jpg)` - Inline with text

^ Image modifiers automatically create appropriate layouts
^ No need to manually configure grid layouts

---

![right](https://pixnio.com/free-images/fauna-animals/dogs/puppy-in-leaves.jpg)

## Flexible Layouts

Switch sides easily by changing the modifier.

Content flows naturally around images.

Great for **visual presentations**!

---

## Code Blocks

DeckSet and gobig both support syntax highlighting:

```javascript
function greet(name) {
  return `Hello, ${name}!`;
}

console.log(greet("World"));
```

^ Code blocks work the same in both formats
^ Syntax highlighting is automatic

---

## Tables

| Feature | gobig | DeckSet |
|---------|-------|---------|
| Markdown | ✅ | ✅ |
| Themes | ✅ | ✅ |
| Layouts | ✅ | ✅ |
| Free | ✅ | ❌ |
| Open Source | ✅ | ❌ |

^ gobig brings DeckSet functionality to everyone
^ Plus it's open source and free

---

[.background-color: #2c3e50]

## Custom Styling

You can customize individual slides with directives.

This slide has a custom background color!

Try different colors to match your brand.

^ Per-slide directives give you fine control
^ Override global settings when needed

---

## Why Use DeckSet Format?

✨ **Familiarity** - If you know DeckSet, you know gobig

🔄 **Portability** - Use existing DeckSet presentations

🎯 **Simplicity** - Clean, readable markdown

💪 **Power** - All of big.js's features

^ The format is intuitive and powerful
^ No vendor lock-in

---

## Getting Started

Install gobig:

```bash
# Download from releases
# Or build from source
make install
```

Create your markdown:

```bash
echo "# [fit] Hello World" > slides.md
gobig -o index.html slides.md
```

^ It's that simple to get started
^ No complex setup required

---

## Both Formats Work

gobig supports **both** syntaxes:

**DeckSet style:**
```markdown
^ Speaker note with caret
[.autoscale: true]
```

**gobig style:**
```markdown
<!-- Speaker note -->
<!-- slide
autoscale: true
-->
```

^ Choose whichever format you prefer
^ Or mix and match!

---

## Resources

- 📖 [gobig Documentation](https://github.com/sroberts/gobig)
- 🎨 [DeckSet Syntax Guide](https://docs.deckset.com)
- 🚀 [big.js Framework](https://github.com/tmcw/big)
- 💡 [Examples](https://github.com/sroberts/gobig/tree/main/examples)

^ Check out these resources to learn more
^ Examples folder has many templates

---

# [fit] Thank You!

Questions?

^ Thanks for watching this demonstration
^ Try gobig with DeckSet format today!
