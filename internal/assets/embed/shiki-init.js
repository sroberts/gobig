// Shiki initialization script
(function() {
  'use strict';
  
  // Import Shiki from CDN using dynamic import
  async function initShiki() {
    try {
      // Import from esm.sh CDN
      const shiki = await import('https://esm.sh/shiki@1.0.0/bundle/web');
      
      // Get theme based on body class
      const bodyClass = document.body.className;
      let theme = 'github-dark';
      if (bodyClass.includes('light')) {
        theme = 'github-light';
      } else if (bodyClass.includes('white')) {
        theme = 'github-light';
      }
      
      // Find all code blocks and highlight them
      const codeBlocks = document.querySelectorAll('pre code');
      
      for (const codeBlock of codeBlocks) {
        const pre = codeBlock.parentElement;
        const code = codeBlock.textContent;
        
        // Get language from class
        let lang = 'text';
        const classes = Array.from(codeBlock.classList);
        for (const cls of classes) {
          if (cls.startsWith('language-')) {
            lang = cls.substring(9); // Remove 'language-' prefix
            break;
          }
        }
        
        try {
          // Highlight the code
          const html = await shiki.codeToHtml(code, {
            lang: lang,
            theme: theme
          });
          
          // Replace the pre element
          const tempDiv = document.createElement('div');
          tempDiv.innerHTML = html;
          const newPre = tempDiv.firstElementChild;
          
          // Preserve any existing classes on the pre element
          if (pre.className) {
            newPre.className = pre.className + ' ' + newPre.className;
          }
          
          pre.replaceWith(newPre);
        } catch (err) {
          console.warn(`Shiki: Could not highlight ${lang}:`, err.message);
          // Keep original code block if highlighting fails
        }
      }
      
      console.log('Shiki highlighting complete');
    } catch (err) {
      console.warn('Shiki failed to load:', err.message);
      // Fallback: Code blocks will display with default styling
    }
  }
  
  // Run when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initShiki);
  } else {
    initShiki();
  }
})();
