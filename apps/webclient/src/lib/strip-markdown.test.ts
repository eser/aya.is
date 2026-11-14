// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { stripMarkdown } from "./strip-markdown.ts";

test("stripMarkdown - bold removal", () => {
  expect(stripMarkdown("This is **bold** text")).toMatchSnapshot();
  expect(stripMarkdown("This is __bold__ text")).toMatchSnapshot();
});

test("stripMarkdown - italic removal", () => {
  expect(stripMarkdown("This is *italic* text")).toMatchSnapshot();
  expect(stripMarkdown("This is _italic_ text")).toMatchSnapshot();
});

test("stripMarkdown - strikethrough removal", () => {
  expect(stripMarkdown("This is ~~deleted~~ text")).toMatchSnapshot();
});

test("stripMarkdown - link removal keeps text", () => {
  expect(stripMarkdown("Visit [our site](https://aya.is) today")).toMatchSnapshot();
});

test("stripMarkdown - image removal keeps alt text", () => {
  expect(stripMarkdown("Look at ![a cat](https://example.com/cat.png) here")).toMatchSnapshot();
});

test("stripMarkdown - code block removal", () => {
  expect(stripMarkdown("Before ```const x = 1;``` after")).toMatchSnapshot();
});

test("stripMarkdown - inline code removal", () => {
  expect(stripMarkdown("Use `console.log()` for debugging")).toMatchSnapshot();
});

test("stripMarkdown - heading markers removal", () => {
  expect(stripMarkdown("# Heading 1")).toMatchSnapshot();
  expect(stripMarkdown("## Heading 2")).toMatchSnapshot();
  expect(stripMarkdown("### Heading 3")).toMatchSnapshot();
});

test("stripMarkdown - blockquote removal", () => {
  expect(stripMarkdown("> This is a quote")).toMatchSnapshot();
});

test("stripMarkdown - list markers removal", () => {
  expect(stripMarkdown("- Item one")).toMatchSnapshot();
  expect(stripMarkdown("* Item two")).toMatchSnapshot();
  expect(stripMarkdown("+ Item three")).toMatchSnapshot();
  expect(stripMarkdown("1. Ordered item")).toMatchSnapshot();
});

test("stripMarkdown - combined formatting", () => {
  const complex = "## **Bold heading** with [a link](https://x.com) and `code`";
  expect(stripMarkdown(complex)).toMatchSnapshot();
});

test("stripMarkdown - edge cases", () => {
  expect(stripMarkdown("")).toMatchSnapshot();
  expect(stripMarkdown("Plain text with no markdown")).toMatchSnapshot();
  expect(stripMarkdown("  leading and trailing whitespace  ")).toMatchSnapshot();
});
