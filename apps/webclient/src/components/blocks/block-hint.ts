// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getAllBlocks } from "./registry";
import type { BlockDefinition } from "./types";
import { MDX_NAMES, parseBlockProps, type ParsedBlockProps } from "./parse-block-props";

type BlockHint = {
  blockId: string;
  name: string;
  icon: React.ElementType;
  propsHint: string;
  definition: BlockDefinition;
  currentValues: Record<string, string>;
  parsedBlock: ParsedBlockProps;
};

function formatPropsHint(block: BlockDefinition): string {
  const propsHint = block.props.map((p) => {
    if (p.type === "select" && p.options !== undefined) {
      return `${p.name}: ${p.options.map((o) => o.value).join("|")}`;
    }
    return `${p.name}: ${p.type}`;
  }).join(", ");

  return propsHint === "" ? "no props" : propsHint;
}

/**
 * Get block hint for the current cursor position.
 * Uses parseBlockProps as the source of truth for tag detection and prop parsing.
 * Falls back to line-only detection for read-only hints when multi-line parse fails
 * (e.g. when the user is still typing and tag is not yet closed).
 */
function getBlockHintForLine(content: string, cursorLine: number): BlockHint | null {
  const parsed = parseBlockProps(content, cursorLine);
  if (parsed !== null) {
    return {
      blockId: parsed.blockId,
      name: parsed.definition.name,
      icon: parsed.definition.icon,
      propsHint: formatPropsHint(parsed.definition),
      definition: parsed.definition,
      currentValues: parsed.props,
      parsedBlock: parsed,
    };
  }

  // Fallback: detect the component name from the line for display-only hints
  // (tag not yet closed, user is typing)
  const lines = content.split("\n");
  if (cursorLine < 0 || cursorLine >= lines.length) return null;

  const line = lines[cursorLine];
  const trimmed = line.trim();
  const match = trimmed.match(/^<([A-Z][a-zA-Z]*)/);
  if (match === null) return null;

  const componentName = match[1];
  const blockId = MDX_NAMES[componentName];
  if (blockId === undefined) return null;

  const allBlocks = getAllBlocks();
  const block = allBlocks.find((b) => b.id === blockId);
  if (block === undefined) return null;

  // Compute a minimal startOffset for the fallback
  let startOffset = 0;
  for (let i = 0; i < cursorLine; i++) {
    startOffset += lines[i].length + 1;
  }
  const ltIndex = line.indexOf("<");
  if (ltIndex !== -1) {
    startOffset += ltIndex;
  }

  // Create a synthetic ParsedBlockProps for the fallback case
  const syntheticParsed: ParsedBlockProps = {
    componentName,
    blockId,
    definition: block,
    props: {},
    selfClosing: false,
    startOffset,
    endOffset: startOffset + line.trimEnd().length - (ltIndex === -1 ? 0 : ltIndex),
  };

  return {
    blockId: block.id,
    name: block.name,
    icon: block.icon,
    propsHint: formatPropsHint(block),
    definition: block,
    currentValues: {},
    parsedBlock: syntheticParsed,
  };
}

export { getBlockHintForLine };
export type { BlockHint };
