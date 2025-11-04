let ASPECT_RATIO = window.BIG_ASPECT_RATIO === undefined ? 1.6 : window.BIG_ASPECT_RATIO;
// Delay in milliseconds to allow DOM layout to complete before scaling notes
const NOTES_SCALING_DELAY = 50;

function parseHash() {
  return parseInt(window.location.hash.substring(1), 10);
}

function emptyNode(node) {
  while (node.hasChildNodes()) node.removeChild(node.lastChild);
}

function ce(type, className = "") {
  return Object.assign(document.createElement(type), { className });
}

addEventListener("load", () => {
  let slideDivs = Array.from(document.querySelectorAll("body > div"));
  let pc = document.body.appendChild(ce("div", "presentation-container"));
  slideDivs = slideDivs.map((slide, _i) => {
    slide.setAttribute("tabindex", 0);
    slide.classList.add("slide");
    let sc = pc.appendChild(ce("div", "slide-container"));
    sc.appendChild(slide);
    return Object.assign(sc, {
      _notes: Array.from(slide.querySelectorAll("notes"), noteElement => {
        noteElement.parentNode.removeChild(noteElement);
        return noteElement.innerHTML.trim();
      }),
      _i
    });
  });
  let timeoutInterval,
    presenterWindow = null,
    presenterStartTime = null,
    presenterTimerInterval = null,
    presenterCheckInterval = null,
    { body } = document,
    {
      className: initialBodyClass,
      style: { cssText: initialBodyStyle }
    } = body,
    big = (window.big = {
      current: -1,
      mode: "talk",
      length: slideDivs.length,
      forward,
      reverse,
      go
    });

  function forward() {
    go(big.current + 1);
  }

  function reverse() {
    go(big.current - 1);
  }

  function go(n, force) {
    n = Math.max(0, Math.min(big.length - 1, n));
    if (!force && big.current === n) return;
    big.current = n;
    let sc = slideDivs[n],
      slideDiv = sc.firstChild;
    if (sc._notes.length) {
      console.group(n);
      for (let note of sc._notes) console.log("%c%s", "padding:5px;font-family:serif;font-size:18px;line-height:150%;", note);
      console.groupEnd();
    }
    for (let slide of slideDivs) slide.style.display = slide._i === n ? "" : "none";
    body.className = `talk-mode ${slideDiv.dataset.bodyClass || ""} ${initialBodyClass}`;
    body.style.cssText = `${initialBodyStyle} ${slideDiv.dataset.bodyStyle || ""}`;
    window.clearInterval(timeoutInterval);
    if (slideDiv.dataset.timeToNext) timeoutInterval = window.setTimeout(forward, parseFloat(slideDiv.dataset.timeToNext) * 1000);
    onResize();
    if (window.location.hash !== n) window.location.hash = n;
    document.title = slideDiv.textContent;
    updatePresenterView();
  }

  function resizeTo(sc, width, height) {
    let slideDiv = sc.firstChild,
      padding = Math.min(width * 0.04),
      fontSize = height;
    sc.style.width = `${width}px`;
    sc.style.height = `${height}px`;
    slideDiv.style.padding = `${padding}px`;
    if (getComputedStyle(slideDiv).display === "grid") slideDiv.style.height = `${height - padding * 2}px`;
    for (let step of [100, 50, 10, 2]) {
      for (; fontSize > 0; fontSize -= step) {
        slideDiv.style.fontSize = `${fontSize}px`;
        if (
          slideDiv.scrollWidth <= width &&
          slideDiv.offsetHeight <= height &&
          Array.from(slideDiv.querySelectorAll("div")).every(elem => elem.scrollWidth <= elem.clientWidth && elem.scrollHeight <= elem.clientHeight)
        ) {
          break;
        }
      }
      fontSize += step;
    }
  }

  function openPresenterView() {
    if (presenterWindow && !presenterWindow.closed) {
      presenterWindow.focus();
      return;
    }
    
    // Start timer when presenter view opens
    presenterStartTime = Date.now();
    
    presenterWindow = window.open("", "big-presenter", "width=1000,height=700,menubar=no,toolbar=no,location=no,status=no");
    
    if (!presenterWindow) {
      alert("Could not open presenter window. Please allow pop-ups for this site.");
      return;
    }
    
    presenterWindow.document.write(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Presenter View</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
      background: #1a1a1a;
      color: #fff;
      overflow: hidden;
      height: 100vh;
    }
    .presenter-container {
      display: grid;
      grid-template-rows: auto 1fr;
      height: 100vh;
    }
    .timer-bar {
      background: #2a2a2a;
      padding: 12px 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      border-bottom: 2px solid #444;
      font-size: 18px;
      font-weight: 600;
    }
    .timer-section {
      display: flex;
      gap: 30px;
    }
    .timer-item {
      display: flex;
      flex-direction: column;
      align-items: center;
    }
    .timer-label {
      font-size: 12px;
      color: #999;
      text-transform: uppercase;
      letter-spacing: 1px;
      margin-bottom: 4px;
    }
    .timer-value {
      font-size: 24px;
      font-weight: 700;
      color: #4CAF50;
    }
    .slide-number {
      color: #fff;
      font-size: 20px;
    }
    .slides-section {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 20px;
      padding: 20px;
      overflow: hidden;
    }
    .slide-preview {
      background: #2a2a2a;
      border-radius: 8px;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      border: 2px solid #444;
    }
    .slide-preview-header {
      background: #333;
      padding: 10px 15px;
      font-weight: 600;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 1px;
      border-bottom: 2px solid #444;
    }
    .slide-preview-content {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;
      position: relative;
      background: #fff;
      aspect-ratio: ${ASPECT_RATIO};
    }
    .slide-preview iframe {
      border: none;
    }
    .notes-section {
      grid-column: 1 / -1;
      background: #2a2a2a;
      border-radius: 8px;
      padding: 20px;
      border: 2px solid #444;
      display: flex;
      flex-direction: column;
      min-height: 150px;
      max-height: 200px;
    }
    .notes-header-container {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
    }
    .notes-header {
      font-weight: 600;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 1px;
      color: #999;
    }
    .notes-toggle-btn {
      background: #444;
      border: 1px solid #666;
      color: #fff;
      padding: 6px 12px;
      font-size: 12px;
      border-radius: 4px;
      cursor: pointer;
      text-transform: uppercase;
      letter-spacing: 0.5px;
      transition: background 0.2s;
    }
    .notes-toggle-btn:hover {
      background: #555;
    }
    .notes-content-wrapper {
      flex: 1;
      overflow: hidden;
      display: flex;
      align-items: center;
    }
    .notes-content-wrapper.scrollable {
      overflow-y: auto;
      align-items: flex-start;
    }
    .notes-content {
      font-size: 16px;
      line-height: 1.6;
      white-space: pre-wrap;
      color: #fff;
      width: 100%;
    }
    .notes-content-wrapper.scale-to-fit .notes-content {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
    }
    .no-notes {
      color: #666;
      font-style: italic;
    }
  </style>
</head>
<body>
  <div class="presenter-container">
    <div class="timer-bar">
      <div class="timer-section">
        <div class="timer-item">
          <div class="timer-label">Elapsed</div>
          <div class="timer-value" id="elapsed">00:00</div>
        </div>
        <div class="timer-item">
          <div class="timer-label">Current Time</div>
          <div class="timer-value" id="current-time">00:00</div>
        </div>
      </div>
      <div class="slide-number" id="slide-number">Slide 1 / 1</div>
    </div>
    <div class="slides-section">
      <div class="slide-preview">
        <div class="slide-preview-header">Current Slide</div>
        <div class="slide-preview-content" id="current-slide"></div>
      </div>
      <div class="slide-preview">
        <div class="slide-preview-header">Next Slide</div>
        <div class="slide-preview-content" id="next-slide"></div>
      </div>
      <div class="notes-section">
        <div class="notes-header-container">
          <div class="notes-header">Speaker Notes</div>
          <button class="notes-toggle-btn" id="notes-toggle">Scrollable</button>
        </div>
        <div class="notes-content-wrapper scale-to-fit" id="notes-wrapper">
          <div class="notes-content" id="notes"></div>
        </div>
      </div>
    </div>
  </div>
</body>
</html>
    `);
    presenterWindow.document.close();

    // Add keyboard navigation to presenter view
    presenterWindow.document.addEventListener("keydown", e => {
      switch (e.key) {
        case "ArrowLeft":
        case "ArrowUp":
        case "PageUp":
          return reverse();
        case "ArrowRight":
        case "ArrowDown":
        case "PageDown":
        case " ":
          return forward();
      }
    });

    // Add click navigation to presenter view
    presenterWindow.document.addEventListener("click", e => {
      // Don't navigate if clicking on links, interactive elements, or the notes toggle button
      if (e.target.tagName === "A" || e.target.tagName === "IFRAME" || e.target.tagName === "BUTTON") return;
      forward();
    });

    // Add notes toggle functionality
    const notesToggleBtn = presenterWindow.document.getElementById("notes-toggle");
    const notesWrapper = presenterWindow.document.getElementById("notes-wrapper");
    let notesMode = "scale-to-fit"; // Start in scale-to-fit mode

    notesToggleBtn.addEventListener("click", () => {
      if (notesMode === "scale-to-fit") {
        notesMode = "scrollable";
        notesWrapper.classList.remove("scale-to-fit");
        notesWrapper.classList.add("scrollable");
        notesToggleBtn.textContent = "Scale to Fit";
      } else {
        notesMode = "scale-to-fit";
        notesWrapper.classList.remove("scrollable");
        notesWrapper.classList.add("scale-to-fit");
        notesToggleBtn.textContent = "Scrollable";
        // Re-scale the notes
        scaleNotesToFit();
      }
    });

    // Clear any existing intervals
    if (presenterTimerInterval) clearInterval(presenterTimerInterval);

    // Combined interval for timer updates and window check
    presenterTimerInterval = setInterval(() => {
      if (!presenterWindow || presenterWindow.closed) {
        presenterWindow = null;
        presenterStartTime = null;
        clearInterval(presenterTimerInterval);
        presenterTimerInterval = null;
      } else {
        updatePresenterTimers();
      }
    }, 1000);

    // Update presenter view with current slide
    setTimeout(() => updatePresenterView(), 100);
  }

  function updatePresenterTimers() {
    if (!presenterWindow || presenterWindow.closed) return;

    const doc = presenterWindow.document;

    // Update elapsed time
    if (presenterStartTime) {
      const elapsed = Math.floor((Date.now() - presenterStartTime) / 1000);
      const minutes = Math.floor(elapsed / 60);
      const seconds = elapsed % 60;
      const elapsedEl = doc.getElementById("elapsed");
      if (elapsedEl) {
        elapsedEl.textContent = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
      }
    }

    // Update current time
    const now = new Date();
    const hours = now.getHours();
    const minutes = now.getMinutes();
    const timeEl = doc.getElementById("current-time");
    if (timeEl) {
      timeEl.textContent = `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}`;
    }
  }

  function scaleNotesToFit() {
    if (!presenterWindow || presenterWindow.closed) return;

    const doc = presenterWindow.document;
    const notesWrapper = doc.getElementById("notes-wrapper");
    const notesContent = doc.getElementById("notes");

    // Only scale if in scale-to-fit mode
    if (!notesWrapper || !notesWrapper.classList.contains("scale-to-fit")) return;
    if (!notesContent || notesContent.classList.contains("no-notes")) return;

    // Reset font size to default
    notesContent.style.fontSize = "16px";

    // Get container dimensions
    const wrapperHeight = notesWrapper.clientHeight;
    const wrapperWidth = notesWrapper.clientWidth;

    // Calculate if content fits
    let fontSize = 16;
    const minFontSize = 8;
    const maxFontSize = 20;

    // Try to find optimal font size
    for (let step of [2, 1, 0.5]) {
      let hitMinFontSize = true;
      while (fontSize >= minFontSize) {
        notesContent.style.fontSize = `${fontSize}px`;

        if (notesContent.scrollHeight <= wrapperHeight &&
            notesContent.scrollWidth <= wrapperWidth) {
          hitMinFontSize = false;
          break;
        }
        fontSize -= step;
      }
      // Always add back the last step after the loop, since fontSize was reduced one extra time
      fontSize += step;
    }

    // Clamp to min/max
    fontSize = Math.max(minFontSize, Math.min(maxFontSize, fontSize));
    notesContent.style.fontSize = `${fontSize}px`;
  }

  function createSlidePreview(slideContainer, targetElement, viewportWidth, viewportHeight, scale) {
    const clone = slideContainer.cloneNode(true);

    // Create a wrapper div to handle the scaled iframe
    const wrapper = targetElement.ownerDocument.createElement("div");
    wrapper.style.width = `${Math.floor(viewportWidth * scale)}px`;
    wrapper.style.height = `${Math.floor(viewportHeight * scale)}px`;
    wrapper.style.overflow = "hidden";
    wrapper.style.position = "relative";

    const iframe = targetElement.ownerDocument.createElement("iframe");

    // Set iframe dimensions and scaling
    iframe.style.width = `${viewportWidth}px`;
    iframe.style.height = `${viewportHeight}px`;
    iframe.style.transform = `scale(${scale})`;
    iframe.style.transformOrigin = "top left";
    iframe.style.border = "none";
    iframe.style.position = "absolute";
    iframe.style.top = "0";
    iframe.style.left = "0";

    iframe.setAttribute("sandbox", "allow-same-origin allow-scripts");
    wrapper.appendChild(iframe);
    targetElement.appendChild(wrapper);

    const styleEls = document.querySelectorAll("style");
    const styles = Array.from(styleEls).map(el => el.textContent).join("\n");
    iframe.contentDocument.write(`
      <html>
      <head>
        <style>
          ${styles}
          html, body {
            margin: 0;
            padding: 0;
            width: 100%;
            height: 100%;
            overflow: hidden;
          }
        </style>
      </head>
      <body class="${document.body.className}">
        ${clone.outerHTML}
      </body>
      </html>
    `);
    iframe.contentDocument.close();

    // Apply proper sizing using the same resizeTo logic as main presentation
    const sc = iframe.contentDocument.querySelector('.slide-container');
    if (sc) {
      // Ensure slide-container is visible (it may have display:none from main view)
      sc.style.display = 'flex';

      const slideDiv = sc.firstChild;
      if (slideDiv) {
        let padding = Math.min(viewportWidth * 0.04);
        let fontSize = viewportHeight;
        sc.style.width = `${viewportWidth}px`;
        sc.style.height = `${viewportHeight}px`;
        slideDiv.style.padding = `${padding}px`;
        if (iframe.contentWindow.getComputedStyle(slideDiv).display === "grid") {
          slideDiv.style.height = `${viewportHeight - padding * 2}px`;
        }
        // Calculate optimal font size
        for (let step of [100, 50, 10, 2]) {
          for (; fontSize > 0; fontSize -= step) {
            slideDiv.style.fontSize = `${fontSize}px`;
            if (
              slideDiv.scrollWidth <= viewportWidth &&
              slideDiv.offsetHeight <= viewportHeight &&
              Array.from(slideDiv.querySelectorAll("div")).every(elem =>
                elem.scrollWidth <= elem.clientWidth && elem.scrollHeight <= elem.clientHeight
              )
            ) {
              break;
            }
          }
          fontSize += step;
        }
      }
    }
  }

  function updatePresenterView() {
    if (!presenterWindow || presenterWindow.closed) return;

    const doc = presenterWindow.document;
    const currentIdx = big.current;
    const nextIdx = currentIdx + 1;

    // Update slide number
    const slideNumberEl = doc.getElementById("slide-number");
    if (slideNumberEl) {
      slideNumberEl.textContent = `Slide ${currentIdx + 1} / ${big.length}`;
    }

    // Use main window's viewport dimensions for proper scaling
    const mainViewportWidth = document.documentElement.clientWidth;
    const mainViewportHeight = document.documentElement.clientHeight;

    // Update current slide preview
    const currentSlideEl = doc.getElementById("current-slide");
    if (currentSlideEl) {
      currentSlideEl.innerHTML = "";
      if (slideDivs[currentIdx]) {
        // Get actual container dimensions for scaling
        const containerWidth = currentSlideEl.clientWidth;
        const containerHeight = currentSlideEl.clientHeight;

        // Calculate scale to fit main viewport in container
        const scale = Math.min(
          containerWidth / mainViewportWidth,
          containerHeight / mainViewportHeight
        );

        createSlidePreview(slideDivs[currentIdx], currentSlideEl, mainViewportWidth, mainViewportHeight, scale);
      }
    }

    // Update next slide preview
    const nextSlideEl = doc.getElementById("next-slide");
    if (nextSlideEl) {
      nextSlideEl.innerHTML = "";
      if (nextIdx < big.length && slideDivs[nextIdx]) {
        // Get actual container dimensions for scaling
        const containerWidth = nextSlideEl.clientWidth;
        const containerHeight = nextSlideEl.clientHeight;

        // Calculate scale to fit main viewport in container
        const scale = Math.min(
          containerWidth / mainViewportWidth,
          containerHeight / mainViewportHeight
        );

        createSlidePreview(slideDivs[nextIdx], nextSlideEl, mainViewportWidth, mainViewportHeight, scale);
      } else {
        nextSlideEl.innerHTML = '<div style="color:#666;padding:20px;text-align:center;font-size:18px;">End of presentation</div>';
      }
    }
    
    // Update notes
    const notesEl = doc.getElementById("notes");
    if (notesEl) {
      const currentNotes = slideDivs[currentIdx]._notes;
      if (currentNotes && currentNotes.length > 0) {
        // Clear existing content
        notesEl.innerHTML = "";
        notesEl.className = "notes-content";
        // Safely add notes as text nodes to prevent XSS
        currentNotes.forEach((note, i) => {
          const p = doc.createElement("p");
          p.textContent = note;
          notesEl.appendChild(p);
          if (i < currentNotes.length - 1) {
            notesEl.appendChild(doc.createElement("br"));
          }
        });
      } else {
        notesEl.textContent = "No speaker notes for this slide";
        notesEl.className = "notes-content no-notes";
      }
    }

    updatePresenterTimers();

    // Scale notes to fit if in scale-to-fit mode
    setTimeout(() => scaleNotesToFit(), NOTES_SCALING_DELAY);
  }

  function onPrint() {
    if (big.mode === "print") return;
    body.className = `print-mode ${initialBodyClass}`;
    body.style.cssText = initialBodyStyle;
    emptyNode(pc);
    for (let sc of slideDivs) {
      let subContainer = pc.appendChild(ce("div", "sub-container")),
        sbc = subContainer.appendChild(ce("div", sc.firstChild.dataset.bodyClass || ""));
      sbc.appendChild(sc);
      sbc.style.cssText = sc.dataset.bodyStyle || "";
      sc.style.display = "flex";
      resizeTo(sc, 512, 320);
      if (sc._notes.length) continue;
      let notesUl = subContainer.appendChild(ce("ul", "notes-list"));
      for (let note of sc._notes) {
        let li = notesUl.appendChild(ce("li"));
        li.innerText = note;
      }
    }
    big.mode = "print";
  }

  function onTalk(i) {
    if (big.mode === "talk") return;
    big.mode = "talk";
    body.className = `talk-mode ${initialBodyClass}`;
    emptyNode(pc);
    for (let sc of slideDivs) pc.appendChild(sc);
    go(i, true);
  }

  function onJump() {
    if (big.mode === "jump") return;
    big.mode = "jump";
    body.className = "jump-mode " + initialBodyClass;
    body.style.cssText = initialBodyStyle;
    emptyNode(pc);
    slideDivs.forEach(sc => {
      let subContainer = pc.appendChild(ce("div", "sub-container"));
      subContainer.addEventListener("keypress", e => {
        if (e.key !== "Enter") return;
        subContainer.removeEventListener("click", onClickSlide);
        e.stopPropagation();
        e.preventDefault();
        onTalk(sc._i);
      });
      let sbc = subContainer.appendChild(ce("div", sc.firstChild.dataset.bodyClass || ""));
      sbc.appendChild(sc);
      sc.style.display = "flex";
      sbc.style.cssText = sc.dataset.bodyStyle || "";
      resizeTo(sc, 192, 120);
      function onClickSlide(e) {
        subContainer.removeEventListener("click", onClickSlide);
        e.stopPropagation();
        e.preventDefault();
        onTalk(sc._i);
      }
      subContainer.addEventListener("click", onClickSlide);
    });
  }

  function onClick(e) {
    if (big.mode !== "talk") return;
    if (e.target.tagName !== "A") go((big.current + 1) % big.length);
  }

  function onKeyDown(e) {
    if (big.mode === "talk") {
      switch (e.key) {
        case "ArrowLeft":
        case "ArrowUp":
        case "PageUp":
          return reverse();
        case "ArrowRight":
        case "ArrowDown":
        case "PageDown":
          return forward();
      }
    }
    let m = { p: onPrint, t: onTalk, j: onJump, r: openPresenterView }[e.key];
    if (m) m(big.current);
  }

  function onResize() {
    if (big.mode !== "talk") return;
    let { clientWidth: width, clientHeight: height } = document.documentElement;
    if (ASPECT_RATIO !== false) {
      if (width / height > ASPECT_RATIO) width = Math.ceil(height * ASPECT_RATIO);
      else height = Math.ceil(width / ASPECT_RATIO);
    }
    resizeTo(slideDivs[big.current], width, height);
  }

  window.matchMedia("print").addListener(onPrint);
  document.addEventListener("click", onClick);
  document.addEventListener("keydown", onKeyDown);
  document.addEventListener("touchstart", e => {
    if (big.mode !== "talk") return;
    let { pageX: startingPageX } = e.changedTouches[0];
    document.addEventListener(
      "touchend",
      e2 => {
        let distanceTraveled = e2.changedTouches[0].pageX - startingPageX;
        // Don't navigate if the person didn't swipe by fewer than 4 pixels
        if (Math.abs(distanceTraveled) < 4) return;
        if (distanceTraveled < 0) forward();
        else reverse();
      },
      { once: true }
    );
  });
  addEventListener("hashchange", () => {
    if (big.mode === "talk") go(parseHash());
  });
  addEventListener("resize", onResize);
  console.log("This is a big presentation. You can: \n\n* press j to jump to a slide\n" + "* press p to see the print view\n* press t to go back to the talk view\n* press r to open the presenter view");
  body.className = `talk-mode ${initialBodyClass}`;
  go(parseHash() || big.current);
});
