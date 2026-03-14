/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals } from "@std/assert";
import { SPACER_SIZES } from "./spacer/spacer-block.tsx";

Deno.test("SPACER_SIZES mapping", async (t) => {
  await assertSnapshot(t, SPACER_SIZES);
});

Deno.test("SPACER_SIZES has all sizes", () => {
  assertEquals("sm" in SPACER_SIZES, true);
  assertEquals("md" in SPACER_SIZES, true);
  assertEquals("lg" in SPACER_SIZES, true);
  assertEquals("xl" in SPACER_SIZES, true);
});
