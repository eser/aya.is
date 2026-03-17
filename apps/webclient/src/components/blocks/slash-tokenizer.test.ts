// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertEquals } from "@std/assert";
import { isSlashCommandContext } from "./slash-command-tokenizer.ts";

Deno.test("normal text is safe context", () => {
  assertEquals(isSlashCommandContext("Hello world", 11), true);
});

Deno.test("empty document is safe", () => {
  assertEquals(isSlashCommandContext("", 0), true);
});

Deno.test("cursor at document start is safe", () => {
  assertEquals(isSlashCommandContext("text", 0), true);
});

Deno.test("inside fenced code block", () => {
  const text = "```\ncode here\n```";
  assertEquals(isSlashCommandContext(text, 8), false);
});

Deno.test("after closed fenced code block", () => {
  const text = "```\ncode\n```\nnormal";
  assertEquals(isSlashCommandContext(text, text.length), true);
});

Deno.test("inside inline code", () => {
  const text = "some `inline code` here";
  assertEquals(isSlashCommandContext(text, 10), false);
});

Deno.test("after inline code", () => {
  const text = "some `code` here";
  assertEquals(isSlashCommandContext(text, 16), true);
});

Deno.test("inside JSX opening tag", () => {
  const text = '<Component prop="val"';
  assertEquals(isSlashCommandContext(text, 15), false);
});

Deno.test("after self-closing JSX", () => {
  const text = "<Component />\nnormal";
  assertEquals(isSlashCommandContext(text, text.length), true);
});

Deno.test("after regular JSX close", () => {
  const text = "<div>\ncontent\n</div>\nnormal";
  assertEquals(isSlashCommandContext(text, text.length), true);
});

Deno.test("inside HTML comment", () => {
  const text = "<!-- this is a comment -->normal";
  assertEquals(isSlashCommandContext(text, 10), false);
});

Deno.test("after HTML comment", () => {
  const text = "<!-- comment -->\nnormal";
  assertEquals(isSlashCommandContext(text, text.length), true);
});

Deno.test("unclosed fence is code context", () => {
  const text = "```\ncode without closing";
  assertEquals(isSlashCommandContext(text, text.length), false);
});

Deno.test("multiple code fences - alternating", () => {
  const text = "```\ncode1\n```\nnormal\n```\ncode2\n```";
  // cursor in "normal" section — find the offset of "normal"
  const normalStart = text.indexOf("normal");
  assertEquals(isSlashCommandContext(text, normalStart + 3), true);
});

Deno.test("nested JSX is inside tag context", () => {
  const text = "<Outer><Inner />";
  // position 10 is inside "<Inner />" — after <Outer> closes at position 7, we're in normal.
  // Then < at position 7 starts "<Inner", which enters jsx_tag state.
  // position 10 is inside the Inner tag (between <Inner and />)
  assertEquals(isSlashCommandContext(text, 10), false);
});
