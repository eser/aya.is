import { escapeForMdxProp } from "./mdx-template";
import type { BlockDefinition } from "./types";

/**
 * Regenerates an MDX opening tag from component name + prop values.
 * Uses BlockDefinition.props metadata to determine formatting:
 * - "number" and "boolean" types -> {value} (JSX expression)
 * - "string", "select", etc -> "value" (quoted)
 */
export function generateBlockTag(
  componentName: string,
  props: Record<string, string>,
  definition: BlockDefinition,
  selfClosing: boolean,
): string {
  // Build a map of prop name -> type from definition
  const propTypes: Record<string, string> = {};
  for (const p of definition.props) {
    propTypes[p.name] = p.type;
  }

  const parts: string[] = [];

  for (const [name, value] of Object.entries(props)) {
    if (value === "" || value === undefined) continue;

    const propType = propTypes[name];

    if (propType === "number") {
      parts.push(`${name}={${value}}`);
    } else if (propType === "boolean") {
      // For boolean true, can use bare prop or {true}
      if (value === "true") {
        parts.push(`${name}={true}`);
      } else {
        parts.push(`${name}={false}`);
      }
    } else {
      // String, select, color, date, rich-text — use quoted string
      parts.push(`${name}="${escapeForMdxProp(value)}"`);
    }
  }

  const propsStr = parts.length > 0 ? ` ${parts.join(" ")}` : "";

  if (selfClosing) {
    return `<${componentName}${propsStr} />`;
  }
  return `<${componentName}${propsStr}>`;
}
