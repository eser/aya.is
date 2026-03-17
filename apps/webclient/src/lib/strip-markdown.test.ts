// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { stripMarkdown } from "./strip-markdown.ts";

Deno.test("stripMarkdown - bold removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("This is **bold** text"));
  await assertSnapshot(t, stripMarkdown("This is __bold__ text"));
});

Deno.test("stripMarkdown - italic removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("This is *italic* text"));
  await assertSnapshot(t, stripMarkdown("This is _italic_ text"));
});

Deno.test("stripMarkdown - strikethrough removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("This is ~~deleted~~ text"));
});

Deno.test("stripMarkdown - link removal keeps text", async (t) => {
  await assertSnapshot(t, stripMarkdown("Visit [our site](https://aya.is) today"));
});

Deno.test("stripMarkdown - image removal keeps alt text", async (t) => {
  await assertSnapshot(t, stripMarkdown("Look at ![a cat](https://example.com/cat.png) here"));
});

Deno.test("stripMarkdown - code block removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("Before ```const x = 1;``` after"));
});

Deno.test("stripMarkdown - inline code removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("Use `console.log()` for debugging"));
});

Deno.test("stripMarkdown - heading markers removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("# Heading 1"));
  await assertSnapshot(t, stripMarkdown("## Heading 2"));
  await assertSnapshot(t, stripMarkdown("### Heading 3"));
});

Deno.test("stripMarkdown - blockquote removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("> This is a quote"));
});

Deno.test("stripMarkdown - list markers removal", async (t) => {
  await assertSnapshot(t, stripMarkdown("- Item one"));
  await assertSnapshot(t, stripMarkdown("* Item two"));
  await assertSnapshot(t, stripMarkdown("+ Item three"));
  await assertSnapshot(t, stripMarkdown("1. Ordered item"));
});

Deno.test("stripMarkdown - combined formatting", async (t) => {
  const complex = "## **Bold heading** with [a link](https://x.com) and `code`";
  await assertSnapshot(t, stripMarkdown(complex));
});

Deno.test("stripMarkdown - edge cases", async (t) => {
  await assertSnapshot(t, stripMarkdown(""));
  await assertSnapshot(t, stripMarkdown("Plain text with no markdown"));
  await assertSnapshot(t, stripMarkdown("  leading and trailing whitespace  "));
});
