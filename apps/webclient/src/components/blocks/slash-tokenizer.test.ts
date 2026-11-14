// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { isSlashCommandContext } from "./slash-command-tokenizer.ts";

test("normal text is safe context", () => {
  expect(isSlashCommandContext("Hello world", 11)).toEqual(true);
});

test("empty document is safe", () => {
  expect(isSlashCommandContext("", 0)).toEqual(true);
});

test("cursor at document start is safe", () => {
  expect(isSlashCommandContext("text", 0)).toEqual(true);
});

test("inside fenced code block", () => {
  const text = "```\ncode here\n```";
  expect(isSlashCommandContext(text, 8)).toEqual(false);
});

test("after closed fenced code block", () => {
  const text = "```\ncode\n```\nnormal";
  expect(isSlashCommandContext(text, text.length)).toEqual(true);
});

test("inside inline code", () => {
  const text = "some `inline code` here";
  expect(isSlashCommandContext(text, 10)).toEqual(false);
});

test("after inline code", () => {
  const text = "some `code` here";
  expect(isSlashCommandContext(text, 16)).toEqual(true);
});

test("inside JSX opening tag", () => {
  const text = '<Component prop="val"';
  expect(isSlashCommandContext(text, 15)).toEqual(false);
});

test("after self-closing JSX", () => {
  const text = "<Component />\nnormal";
  expect(isSlashCommandContext(text, text.length)).toEqual(true);
});

test("after regular JSX close", () => {
  const text = "<div>\ncontent\n</div>\nnormal";
  expect(isSlashCommandContext(text, text.length)).toEqual(true);
});

test("inside HTML comment", () => {
  const text = "<!-- this is a comment -->normal";
  expect(isSlashCommandContext(text, 10)).toEqual(false);
});

test("after HTML comment", () => {
  const text = "<!-- comment -->\nnormal";
  expect(isSlashCommandContext(text, text.length)).toEqual(true);
});

test("unclosed fence is code context", () => {
  const text = "```\ncode without closing";
  expect(isSlashCommandContext(text, text.length)).toEqual(false);
});

test("multiple code fences - alternating", () => {
  const text = "```\ncode1\n```\nnormal\n```\ncode2\n```";
  // cursor in "normal" section — find the offset of "normal"
  const normalStart = text.indexOf("normal");
  expect(isSlashCommandContext(text, normalStart + 3)).toEqual(true);
});

test("nested JSX is inside tag context", () => {
  const text = "<Outer><Inner />";
  // position 10 is inside "<Inner />" — after <Outer> closes at position 7, we're in normal.
  // Then < at position 7 starts "<Inner", which enters jsx_tag state.
  // position 10 is inside the Inner tag (between <Inner and />)
  expect(isSlashCommandContext(text, 10)).toEqual(false);
});
