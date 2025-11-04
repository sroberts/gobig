autoscale: true
slidenumbers: true
footer: DeckSet Format Demo

---

# [fit] DeckSet Format
## Testing gobig with DeckSet syntax

^ This is a speaker note using DeckSet syntax
^ It uses the caret prefix

---

## Features

- Simple markdown syntax
- Horizontal rules (`---`) separate slides
- DeckSet-style speaker notes
- Multiple themes

^ Remember to mention the DeckSet compatibility

---

[.autoscale: true]
[.build-lists: true]

## Per-Slide Directives

- This slide has autoscale enabled
- Build lists are also enabled
- These are DeckSet directives

---

![left](https://pixnio.com/free-images/fauna-animals/dogs/puppy-in-leaves.jpg)

## Side by Side

Text appears next to the image with DeckSet's **left** modifier.

Perfect for *image + description* slides!

^ This demonstrates DeckSet's image positioning

---

![right](https://pixnio.com/free-images/fauna-animals/dogs/puppy-in-leaves.jpg)

## Right Side Image

Image on the right, content on the left.

Great for alternating layouts.

---

## Code Example

```javascript
function demo() {
  console.log("Keep code snippets short!");
  return "8 lines or fewer works best";
}
```

^ Remind the audience that code should be minimal and focused

---

[.background-color: #1a1a2e]

## Custom Background

This slide has a custom background color using DeckSet directive.

---

# Thank You!

**gobig** now supports DeckSet format!

Built with ❤️ using Go
