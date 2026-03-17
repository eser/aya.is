// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/**
 * MDX-aware tokenizer that determines whether the cursor is in a "safe"
 * context for the slash command inserter to open.
 *
 * The slash command should NOT trigger when the cursor is inside:
 * - Fenced code blocks (``` ... ```)
 * - Inline code (` ... `)
 * - JSX tags (< ... >)
 * - HTML comments (<!-- ... -->)
 *
 * State machine:
 * ┌─────────┐  ```   ┌───────────┐
 * │ NORMAL  │───────▶│ CODE_FENCE │
 * │         │◀───────│           │
 * └────┬────┘  ```   └───────────┘
 *      │  `
 *      ▼
 * ┌───────────┐
 * │INLINE_CODE│  (until matching `)
 * └───────────┘
 *      │  <
 *      ▼
 * ┌─────────┐    ┌─────────────┐
 * │ JSX_TAG │    │HTML_COMMENT  │
 * │(< ... >)│    │(<!-- ... -->)│
 * └─────────┘    └─────────────┘
 */

type TokenizerState = "normal" | "code_fence" | "inline_code" | "jsx_tag" | "html_comment";

const TIMEOUT_MS = 50;

/**
 * Determines if the cursor position is in a safe context for the slash command.
 * Returns true if the slash command inserter should open, false if suppressed.
 *
 * Pure function with a 50ms timeout guard for very long documents.
 * On timeout, falls back to "allow" (safe default — worst case is the
 * inserter opens inside a code block, and user just presses Escape).
 */
export function isSlashCommandContext(text: string, cursorPos: number): boolean {
  const startTime = performance.now();
  const end = Math.min(cursorPos, text.length);

  let state: TokenizerState = "normal";
  let i = 0;

  while (i < end) {
    // Timeout guard: fall back to "allow" for very long documents
    if ((i & 0x3FF) === 0 && performance.now() - startTime > TIMEOUT_MS) {
      return true;
    }

    const char = text[i];
    const next = i + 1 < end ? text[i + 1] : "";
    const next2 = i + 2 < end ? text[i + 2] : "";

    switch (state) {
      case "normal": {
        // Check for fenced code block opening: ``` at line start
        if (char === "`" && next === "`" && next2 === "`") {
          // Verify it's at the start of a line (or document start)
          if (i === 0 || text[i - 1] === "\n") {
            state = "code_fence";
            // Skip past the opening ``` and any language identifier
            i += 3;
            while (i < end && text[i] !== "\n") {
              i++;
            }
            continue;
          }
        }

        // Check for inline code
        if (char === "`") {
          state = "inline_code";
          i++;
          continue;
        }

        // Check for HTML comment opening: <!--
        if (char === "<" && next === "!" && next2 === "-" && i + 3 < end && text[i + 3] === "-") {
          state = "html_comment";
          i += 4;
          continue;
        }

        // Check for JSX tag opening: < followed by a letter or /
        if (char === "<" && (/[a-zA-Z/]/).test(next)) {
          state = "jsx_tag";
          i++;
          continue;
        }

        i++;
        break;
      }

      case "code_fence": {
        // Look for closing ``` at line start
        if (char === "`" && next === "`" && next2 === "`") {
          if (i === 0 || text[i - 1] === "\n") {
            state = "normal";
            i += 3;
            continue;
          }
        }
        i++;
        break;
      }

      case "inline_code": {
        // Closing backtick ends inline code
        if (char === "`") {
          state = "normal";
        }
        // Newline also ends inline code (inline code can't span lines in standard MD)
        if (char === "\n") {
          state = "normal";
        }
        i++;
        break;
      }

      case "jsx_tag": {
        // Self-closing tag end: />
        if (char === "/" && next === ">") {
          state = "normal";
          i += 2;
          continue;
        }
        // Regular closing: >
        if (char === ">") {
          state = "normal";
          i++;
          continue;
        }
        // Newline inside a JSX tag is OK (multi-line props) — stay in jsx_tag
        i++;
        break;
      }

      case "html_comment": {
        // Closing: -->
        if (char === "-" && next === "-" && next2 === ">") {
          state = "normal";
          i += 3;
          continue;
        }
        i++;
        break;
      }
    }
  }

  // Safe context only when in normal state
  return state === "normal";
}
