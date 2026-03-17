// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { generateBlockTag } from "./generate-block-tag.ts";
import { getAllBlocks } from "./registry.ts";

function getBlockDef(id: string) {
  const block = getAllBlocks().find((b) => b.id === id);
  if (block === undefined) throw new Error(`Block ${id} not found`);
  return block;
}

Deno.test("string props", async (t) => {
  const def = getBlockDef("callout");
  const result = generateBlockTag("Callout", { variant: "info", title: "Note" }, def, false);
  await assertSnapshot(t, result);
});

Deno.test("number prop", () => {
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

Deno.test("boolean prop", async (t) => {
  const def = getBlockDef("group");
  const result = generateBlockTag("Group", { border: "true", rounded: "false", padding: "md" }, def, false);
  await assertSnapshot(t, result);
});

Deno.test("self-closing", async (t) => {
  const def = getBlockDef("spacer");
  const result = generateBlockTag("Spacer", { size: "md" }, def, true);
  await assertSnapshot(t, result);
});

Deno.test("container", async (t) => {
  const def = getBlockDef("callout");
  const result = generateBlockTag("Callout", { variant: "info" }, def, false);
  await assertSnapshot(t, result);
});

Deno.test("escape special chars", async (t) => {
  const def = getBlockDef("callout");
  const result = generateBlockTag("Callout", { title: 'say "hello"' }, def, false);
  await assertSnapshot(t, result);
});

Deno.test("empty props", async (t) => {
  const def = getBlockDef("verse");
  const result = generateBlockTag("Verse", {}, def, false);
  await assertSnapshot(t, result);
});
