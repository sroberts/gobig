let ASPECT_RATIO = window.BIG_ASPECT_RATIO === undefined ? 1.6 : window.BIG_ASPECT_RATIO;

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
      min-height: 300px;
    }
    .slide-preview iframe {
      border: none;
      transform-origin: top left;
      position: absolute;
      top: 0;
      left: 0;
    }
    .notes-section {
      grid-column: 1 / -1;
      background: #2a2a2a;
      border-radius: 8px;
      padding: 20px;
      max-height: 200px;
      overflow-y: auto;
      border: 2px solid #444;
    }
    .notes-header {
      font-weight: 600;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 1px;
      margin-bottom: 12px;
      color: #999;
    }
    .notes-content {
      font-size: 16px;
      line-height: 1.6;
      white-space: pre-wrap;
      color: #fff;
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
        <div class="notes-header">Speaker Notes</div>
        <div class="notes-content" id="notes"></div>
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
      // Don't navigate if clicking on links or interactive elements
      if (e.target.tagName === "A" || e.target.tagName === "IFRAME") return;
      forward();
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
    
    // Update current slide preview
    const currentSlideEl = doc.getElementById("current-slide");
    if (currentSlideEl) {
      currentSlideEl.innerHTML = "";
      if (slideDivs[currentIdx]) {
        const clone = slideDivs[currentIdx].cloneNode(true);
        const iframe = doc.createElement("iframe");

        // Get actual viewport dimensions
        const viewportWidth = document.documentElement.clientWidth;
        const viewportHeight = document.documentElement.clientHeight;

        // Calculate scale to fit preview (preview container is roughly 45% of presenter window width)
        const previewWidth = 450; // Target preview width
        const previewHeight = 300; // Target preview height
        const scale = Math.min(previewWidth / viewportWidth, previewHeight / viewportHeight);

        iframe.style.width = `${viewportWidth}px`;
        iframe.style.height = `${viewportHeight}px`;
        iframe.style.transform = `scale(${scale})`;
        iframe.style.transformOrigin = "top left";
        iframe.setAttribute("sandbox", "allow-same-origin");
        currentSlideEl.appendChild(iframe);
        const styleEl = document.querySelector("style");
        const styles = styleEl ? styleEl.textContent : "";
        iframe.contentDocument.write(`
          <html>
          <head>
            <style>
              ${styles}
              body { margin: 0; padding: 0; }
            </style>
          </head>
          <body class="${document.body.className}">
            ${clone.innerHTML}
          </body>
          </html>
        `);
        iframe.contentDocument.close();
      }
    }
    
    // Update next slide preview
    const nextSlideEl = doc.getElementById("next-slide");
    if (nextSlideEl) {
      nextSlideEl.innerHTML = "";
      if (nextIdx < big.length && slideDivs[nextIdx]) {
        const clone = slideDivs[nextIdx].cloneNode(true);
        const iframe = doc.createElement("iframe");

        // Get actual viewport dimensions
        const viewportWidth = document.documentElement.clientWidth;
        const viewportHeight = document.documentElement.clientHeight;

        // Calculate scale to fit preview
        const previewWidth = 450;
        const previewHeight = 300;
        const scale = Math.min(previewWidth / viewportWidth, previewHeight / viewportHeight);

        iframe.style.width = `${viewportWidth}px`;
        iframe.style.height = `${viewportHeight}px`;
        iframe.style.transform = `scale(${scale})`;
        iframe.style.transformOrigin = "top left";
        iframe.setAttribute("sandbox", "allow-same-origin");
        nextSlideEl.appendChild(iframe);
        const styleEl = document.querySelector("style");
        const styles = styleEl ? styleEl.textContent : "";
        iframe.contentDocument.write(`
          <html>
          <head>
            <style>
              ${styles}
              body { margin: 0; padding: 0; }
            </style>
          </head>
          <body class="${document.body.className}">
            ${clone.innerHTML}
          </body>
          </html>
        `);
        iframe.contentDocument.close();
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
