// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getAllBlocks } from "./registry";
import type { BlockDefinition } from "./types";

// MDX component name -> block ID mapping
const MDX_NAMES: Record<string, string> = {
  "Callout": "callout",
  "Details": "details",
  "Columns": "columns-2",
  "Column": "columns-2",
  "Divider": "divider",
  "Spacer": "spacer",
  "Button": "button",
  "Tabs": "tabs",
  "Tab": "tabs",
  "Card": "card",
  "Cards": "cards",
  "YouTubeEmbed": "youtube",
  "TwitterEmbed": "twitter",
  "PDF": "pdf",
  "Embed": "embed",
  "SiteLink": "site-link",
  "Gallery": "gallery",
  "Audio": "audio",
  "Video": "video",
  "Cover": "cover",
  "MediaText": "media-text",
  "Row": "row",
  "Stack": "stack",
  "Grid": "grid",
  "Group": "group",
  "Center": "center",
  "Pullquote": "pullquote",
  "Verse": "verse",
};

export { MDX_NAMES };

type ParsedBlockProps = {
  componentName: string;
  blockId: string;
  definition: BlockDefinition;
  props: Record<string, string>;
  selfClosing: boolean;
  startOffset: number;
  endOffset: number;
};

export type { ParsedBlockProps };

export function parseBlockProps(fullText: string, cursorLine: number): ParsedBlockProps | null {
  const lines = fullText.split("\n");
  if (cursorLine < 0 || cursorLine >= lines.length) return null;

  const line = lines[cursorLine];
  const trimmed = line.trim();

  // Must start with <ComponentName (not closing tag)
  const tagMatch = trimmed.match(/^<([A-Z][a-zA-Z]*)/);
  if (tagMatch === null) return null;

  const componentName = tagMatch[1];
  const blockId = MDX_NAMES[componentName];
  if (blockId === undefined) return null;

  const allBlocks = getAllBlocks();
  const definition = allBlocks.find((b) => b.id === blockId);
  if (definition === undefined) return null;

  // Calculate start offset (absolute) — find where `<` is in the full text
  let startOffset = 0;
  for (let i = 0; i < cursorLine; i++) {
    startOffset += lines[i].length + 1; // +1 for \n
  }
  // Find the `<` within the current line
  const ltIndex = line.indexOf("<");
  if (ltIndex === -1) return null;
  startOffset += ltIndex;

  // Now scan forward to find the closing > or />
  // Accumulate the tag text across lines
  let tagText = "";
  let endOffset = startOffset;
  let found = false;
  let selfClosing = false;

  for (let lineIdx = cursorLine; lineIdx < lines.length; lineIdx++) {
    const currentLine = lineIdx === cursorLine ? line.substring(ltIndex) : lines[lineIdx];

    for (let charIdx = 0; charIdx < currentLine.length; charIdx++) {
      const ch = currentLine[charIdx];
      tagText += ch;

      if (ch === ">") {
        // Check if self-closing
        selfClosing = tagText.endsWith("/>");

        // Compute endOffset properly for multi-line
        if (lineIdx !== cursorLine) {
          endOffset = 0;
          for (let i = 0; i < lineIdx; i++) {
            endOffset += lines[i].length + 1;
          }
          endOffset += charIdx + 1;
        } else {
          endOffset = startOffset + charIdx + 1;
        }

        found = true;
        break;
      }
    }

    if (found) break;
    tagText += "\n"; // line break between lines
  }

  if (!found) return null;

  // Parse props from tagText
  // tagText is like: <ComponentName prop="val" prop={val} open>
  // or <ComponentName prop="val" />
  const props: Record<string, string> = {};

  // Remove the opening <ComponentName and closing > or />
  let propsString = tagText;
  // Remove <ComponentName
  propsString = propsString.replace(/^<[A-Z][a-zA-Z]*/, "");
  // Remove trailing /> or >
  propsString = propsString.replace(/\s*\/?>$/, "");

  // Parse individual props using regex
  // Match: name="value", name='value', name={value}, or bare name (boolean)
  const propPattern = /(\w+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|\{([^}]*)\}))?/g;
  let propMatch;
  while ((propMatch = propPattern.exec(propsString)) !== null) {
    const propName = propMatch[1];
    if (propName === componentName) continue; // skip if it's the component name leftover

    const stringDoubleQuote = propMatch[2];
    const stringSingleQuote = propMatch[3];
    const jsxExpression = propMatch[4];

    if (stringDoubleQuote !== undefined) {
      props[propName] = stringDoubleQuote;
    } else if (stringSingleQuote !== undefined) {
      props[propName] = stringSingleQuote;
    } else if (jsxExpression !== undefined) {
      props[propName] = jsxExpression.trim();
    } else {
      // Bare boolean prop like `open`
      props[propName] = "true";
    }
  }

  return {
    componentName,
    blockId,
    definition,
    props,
    selfClosing,
    startOffset,
    endOffset,
  };
}
