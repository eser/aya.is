/**
 * Escapes a string value for safe use inside an MDX JSX prop.
 * Prevents MDX injection by replacing characters that could break
 * out of a quoted attribute or inject JSX/expressions.
 *
 * Returns "" for null/undefined inputs.
 */
export function escapeForMdxProp(value: string | null | undefined): string {
  if (value === null || value === undefined) {
    return "";
  }

  return value
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/\{/g, "&#123;")
    .replace(/\}/g, "&#125;");
}

/**
 * Formats a single prop for MDX JSX output.
 * Strings use double quotes, numbers and booleans use curly braces.
 */
function formatProp(name: string, value: string | number | boolean): string {
  if (typeof value === "string") {
    return `${name}="${escapeForMdxProp(value)}"`;
  }
  return `${name}={${String(value)}}`;
}

/**
 * DRY helper for generating a self-closing MDX tag.
 *
 * Example: generateSelfClosingMdx("Spacer", { size: "md" })
 * → '\n<Spacer size="md" />\n'
 */
export function generateSelfClosingMdx(
  name: string,
  props: Record<string, string | number | boolean>,
): string {
  const entries = Object.entries(props).filter(
    ([_, v]) => v !== "" && v !== undefined,
  );

  if (entries.length === 0) {
    return `\n<${name} />\n`;
  }

  const propsStr = entries.map(([k, v]) => formatProp(k, v)).join(" ");
  return `\n<${name} ${propsStr} />\n`;
}

/**
 * DRY helper for generating a container MDX tag with children placeholder.
 *
 * Example: generateContainerMdx("Callout", { variant: "info" }, "Your content here.")
 * → '\n<Callout variant="info">\n\nYour content here.\n\n</Callout>\n'
 *
 * Blank lines around children are critical for MDX to parse
 * nested content as markdown (not raw text).
 */
export function generateContainerMdx(
  name: string,
  props: Record<string, string | number | boolean>,
  childrenPlaceholder: string,
): string {
  const entries = Object.entries(props).filter(
    ([_, v]) => v !== "" && v !== undefined,
  );

  const openTag = entries.length === 0
    ? `<${name}>`
    : `<${name} ${entries.map(([k, v]) => formatProp(k, v)).join(" ")}>`;

  return `\n${openTag}\n\n${childrenPlaceholder}\n\n</${name}>\n`;
}
