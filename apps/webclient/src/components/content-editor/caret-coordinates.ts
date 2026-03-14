/**
 * Get screen coordinates of the cursor position in a textarea
 * using the "mirror div" technique.
 */
export function getCaretCoordinates(
  textarea: HTMLTextAreaElement,
): { x: number; y: number } {
  const mirror = document.createElement("div");

  const computed = globalThis.getComputedStyle(textarea);

  const stylesToCopy = [
    "fontFamily",
    "fontSize",
    "fontWeight",
    "fontStyle",
    "lineHeight",
    "letterSpacing",
    "wordSpacing",
    "textIndent",
    "textTransform",
    "whiteSpace",
    "wordWrap",
    "overflowWrap",
    "paddingTop",
    "paddingRight",
    "paddingBottom",
    "paddingLeft",
    "borderTopWidth",
    "borderRightWidth",
    "borderBottomWidth",
    "borderLeftWidth",
    "borderTopStyle",
    "borderRightStyle",
    "borderBottomStyle",
    "borderLeftStyle",
    "boxSizing",
  ] as const;

  for (const prop of stylesToCopy) {
    mirror.style[prop] = computed[prop];
  }

  mirror.style.position = "absolute";
  mirror.style.visibility = "hidden";
  mirror.style.left = "-9999px";
  mirror.style.top = "0";
  mirror.style.overflow = "hidden";
  mirror.style.width = `${textarea.offsetWidth}px`;
  mirror.style.height = "auto";
  mirror.style.whiteSpace = "pre-wrap";
  mirror.style.wordWrap = "break-word";

  const textBeforeCursor = textarea.value.substring(
    0,
    textarea.selectionStart,
  );

  const textNode = document.createTextNode(textBeforeCursor);
  mirror.appendChild(textNode);

  const caretSpan = document.createElement("span");
  caretSpan.textContent = "\u200B"; // zero-width space
  mirror.appendChild(caretSpan);

  document.body.appendChild(mirror);

  const rect = textarea.getBoundingClientRect();
  const x =
    rect.left + caretSpan.offsetLeft - textarea.scrollLeft;
  const y =
    rect.top + caretSpan.offsetTop - textarea.scrollTop;

  document.body.removeChild(mirror);

  return { x, y };
}
