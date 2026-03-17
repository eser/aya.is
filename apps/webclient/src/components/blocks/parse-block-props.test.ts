// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals } from "@std/assert";
import { parseBlockProps } from "./parse-block-props.ts";

Deno.test("single-line tag with string props", async (t) => {
  const result = parseBlockProps('<Callout variant="info" title="Note">', 0);
  assertEquals(result !== null, true);
  await assertSnapshot(t, {
    componentName: result!.componentName,
    blockId: result!.blockId,
    props: result!.props,
    selfClosing: result!.selfClosing,
    startOffset: result!.startOffset,
    endOffset: result!.endOffset,
  });
});

Deno.test("single-line self-closing tag", async (t) => {
  const result = parseBlockProps('<Spacer size="md" />', 0);
  assertEquals(result !== null, true);
  assertEquals(result!.selfClosing, true);
  await assertSnapshot(t, {
    componentName: result!.componentName,
    blockId: result!.blockId,
    props: result!.props,
    selfClosing: result!.selfClosing,
    startOffset: result!.startOffset,
    endOffset: result!.endOffset,
  });
});

Deno.test("JSX expression prop", () => {
  const result = parseBlockProps("<Columns cols={2}>", 0);
  assertEquals(result !== null, true);
  assertEquals(result!.props.cols, "2");
});

Deno.test("bare boolean prop", () => {
  const result = parseBlockProps("<Details open>", 0);
  assertEquals(result !== null, true);
  assertEquals(result!.props.open, "true");
});

Deno.test("mixed props", async (t) => {
  const result = parseBlockProps('<Cover src="url" overlay="dark">', 0);
  assertEquals(result !== null, true);
  await assertSnapshot(t, {
    componentName: result!.componentName,
    blockId: result!.blockId,
    props: result!.props,
    selfClosing: result!.selfClosing,
    startOffset: result!.startOffset,
    endOffset: result!.endOffset,
  });
});

Deno.test("no tag", () => {
  const result = parseBlockProps("Hello world", 0);
  assertEquals(result, null);
});

Deno.test("closing tag", () => {
  const result = parseBlockProps("</Callout>", 0);
  assertEquals(result, null);
});

Deno.test("multi-line tag", () => {
  const result = parseBlockProps('<Cover\n  src="url"\n  overlay="dark">', 0);
  assertEquals(result !== null, true);
  assertEquals(result!.props.src, "url");
  assertEquals(result!.props.overlay, "dark");
});

Deno.test("empty props", () => {
  const result = parseBlockProps("<Verse>", 0);
  assertEquals(result !== null, true);
  assertEquals(Object.keys(result!.props).length, 0);
});

Deno.test("tag with content before it", () => {
  const result = parseBlockProps('some text\n<Callout variant="info">', 1);
  assertEquals(result !== null, true);
  assertEquals(result!.blockId, "callout");
  assertEquals(result!.props.variant, "info");
});

Deno.test("unknown component", () => {
  const result = parseBlockProps('<UnknownThing prop="val">', 0);
  assertEquals(result, null);
});
