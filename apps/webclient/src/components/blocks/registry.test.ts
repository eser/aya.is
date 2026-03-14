/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals, assertNotEquals } from "@std/assert";
import {
  getAllBlocks,
  getAllPatterns,
  getBlocksByCategory,
  searchBlocks,
  searchPatterns,
} from "./registry.ts";

Deno.test("getAllBlocks returns all blocks", async (t) => {
  const blocks = getAllBlocks();
  await assertSnapshot(t, {
    count: blocks.length,
    ids: blocks.map((b) => b.id),
  });
});

Deno.test("all blocks have unique IDs", () => {
  const blocks = getAllBlocks();
  const ids = blocks.map((b) => b.id);
  const uniqueIds = new Set(ids);
  assertEquals(uniqueIds.size, ids.length);
});

Deno.test("all blocks have valid categories", () => {
  const validCategories = ["text", "media", "layout", "data", "embed", "interactive"];
  const blocks = getAllBlocks();
  for (const block of blocks) {
    assertEquals(
      validCategories.includes(block.category),
      true,
      `Block "${block.id}" has invalid category "${block.category}"`,
    );
  }
});

Deno.test("getBlocksByCategory text returns text blocks", async (t) => {
  const textBlocks = getBlocksByCategory("text");
  await assertSnapshot(t, textBlocks.map((b) => b.id));
});

Deno.test("getBlocksByCategory media returns media blocks", async (t) => {
  const mediaBlocks = getBlocksByCategory("media");
  await assertSnapshot(t, mediaBlocks.map((b) => b.id));
});

Deno.test("getBlocksByCategory embed returns embed blocks", async (t) => {
  const embedBlocks = getBlocksByCategory("embed");
  await assertSnapshot(t, embedBlocks.map((b) => b.id));
});

Deno.test("searchBlocks video finds youtube", async (t) => {
  const results = searchBlocks("video");
  await assertSnapshot(t, results.map((b) => b.id));
});

Deno.test("searchBlocks callout finds callout", async (t) => {
  const results = searchBlocks("callout");
  await assertSnapshot(t, results.map((b) => b.id));
});

Deno.test("searchBlocks nonexistent returns empty", () => {
  const results = searchBlocks("xyznonexistent");
  assertEquals(results.length, 0);
});

Deno.test("getAllPatterns returns patterns", async (t) => {
  const patterns = getAllPatterns();
  await assertSnapshot(t, {
    count: patterns.length,
    ids: patterns.map((p) => p.id),
  });
});

Deno.test("searchPatterns hero finds hero pattern", async (t) => {
  const results = searchPatterns("hero");
  await assertSnapshot(t, results.map((p) => p.id));
});
