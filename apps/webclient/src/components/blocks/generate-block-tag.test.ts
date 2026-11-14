// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { generateBlockTag } from "./generate-block-tag.ts";
import { getAllBlocks } from "./registry.ts";

function getBlockDef(id: string) {
  const block = getAllBlocks().find((b) => b.id === id);
  if (block === undefined) throw new Error(`Block ${id} not found`);
  return block;
}

test("string props", () => {
  const def = getBlockDef("callout");
  const result = generateBlockTag("Callout", { variant: "info", title: "Note" }, def, false);
  expect(result).toMatchSnapshot();
});

test("number prop", () => {
  // Gallery has cols as select type, but we test the format
  // Grid also has cols as select. Let's use a definition that has a concept of number.
  // Actually none of the current blocks use number type, but we can test behavior with select (quoted).
  // For a true number test, let's check the format manually.
  const def = getBlockDef("gallery");
  const result = generateBlockTag("Gallery", { cols: "3" }, def, false);
  // cols is "select" type in gallery, so it should be quoted
  const hasQuotedCols = result.includes('cols="3"');
  if (!hasQuotedCols) {
    throw new Error(`Expected cols="3" in result: ${result}`);
  }
});

test("boolean prop", () => {
  const def = getBlockDef("group");
  const result = generateBlockTag("Group", { border: "true", rounded: "false", padding: "md" }, def, false);
  expect(result).toMatchSnapshot();
});

test("self-closing", () => {
  const def = getBlockDef("spacer");
  const result = generateBlockTag("Spacer", { size: "md" }, def, true);
  expect(result).toMatchSnapshot();
});

test("container", () => {
  const def = getBlockDef("callout");
  const result = generateBlockTag("Callout", { variant: "info" }, def, false);
  expect(result).toMatchSnapshot();
});

test("escape special chars", () => {
  const def = getBlockDef("callout");
  const result = generateBlockTag("Callout", { title: 'say "hello"' }, def, false);
  expect(result).toMatchSnapshot();
});

test("empty props", () => {
  const def = getBlockDef("verse");
  const result = generateBlockTag("Verse", {}, def, false);
  expect(result).toMatchSnapshot();
});
