# gobig

A markdown-to-big.js presentation tool

<!-- This is a speaker note. It won't appear on the slide but will be visible in the browser console -->

---

## Features

- Simple markdown syntax
- Horizontal rules (`---`) separate slides
- Speaker notes in HTML comments
- Multiple themes
- Grid-based layouts

---

## Getting Started

Install and run:

```bash
go install gobig
gobig presentation.md -o index.html
```

<!-- Remember to mention the installation options -->

---

## Why big.js?

- **Minimal**: Only ~16KB total
- **Fast**: No dependencies
- **Focused**: Encourages simple, clear slides
- **Flexible**: Supports layouts, themes, and customization

---

## Themes

Three built-in themes:

1. **dark** - Near-black background (default)
2. **light** - Light background with dark text
3. **white** - Stark black and white

```bash
gobig slides.md -theme light -o output.html
```

---

## Thank You!

Questions?
