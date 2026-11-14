// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { getAllBlocks, getAllPatterns, getBlocksByCategory, searchBlocks, searchPatterns } from "./registry.ts";

test("getAllBlocks returns all blocks", () => {
  const blocks = getAllBlocks();
  expect({
    count: blocks.length,
    ids: blocks.map((b) => b.id),
  }).toMatchSnapshot();
});

test("all blocks have unique IDs", () => {
  const blocks = getAllBlocks();
  const ids = blocks.map((b) => b.id);
  const uniqueIds = new Set(ids);
  expect(uniqueIds.size).toEqual(ids.length);
});

test("all blocks have valid categories", () => {
  const validCategories = ["text", "media", "layout", "data", "embed", "interactive"];
  const blocks = getAllBlocks();
  for (const block of blocks) {
    expect(
      validCategories.includes(block.category),
      `Block "${block.id}" has invalid category "${block.category}"`,
    ).toEqual(true);
  }
});

test("getBlocksByCategory text returns text blocks", () => {
  const textBlocks = getBlocksByCategory("text");
  expect(textBlocks.map((b) => b.id)).toMatchSnapshot();
});

test("getBlocksByCategory media returns media blocks", () => {
  const mediaBlocks = getBlocksByCategory("media");
  expect(mediaBlocks.map((b) => b.id)).toMatchSnapshot();
});

test("getBlocksByCategory embed returns embed blocks", () => {
  const embedBlocks = getBlocksByCategory("embed");
  expect(embedBlocks.map((b) => b.id)).toMatchSnapshot();
});

test("searchBlocks video finds youtube", () => {
  const results = searchBlocks("video");
  expect(results.map((b) => b.id)).toMatchSnapshot();
});

test("searchBlocks callout finds callout", () => {
  const results = searchBlocks("callout");
  expect(results.map((b) => b.id)).toMatchSnapshot();
});

test("searchBlocks nonexistent returns empty", () => {
  const results = searchBlocks("xyznonexistent");
  expect(results.length).toEqual(0);
});

test("getAllPatterns returns patterns", () => {
  const patterns = getAllPatterns();
  expect({
    count: patterns.length,
    ids: patterns.map((p) => p.id),
  }).toMatchSnapshot();
});

test("searchPatterns hero finds hero pattern", () => {
  const results = searchPatterns("hero");
  expect(results.map((p) => p.id)).toMatchSnapshot();
});
