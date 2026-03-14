/// <reference lib="deno.ns" />
import { assertEquals } from "@std/assert";
import { getBlockHintForLine } from "./block-hint.ts";

Deno.test("getBlockHintForLine returns hint for <Callout", () => {
  const result = getBlockHintForLine("<Callout", 0);
  assertEquals(result !== null, true);
  assertEquals(result!.blockId, "callout");
});

Deno.test("getBlockHintForLine returns null for <UnknownComponent", () => {
  const result = getBlockHintForLine("<UnknownComponent", 0);
  assertEquals(result, null);
});

Deno.test("getBlockHintForLine returns null for regular text", () => {
  const result = getBlockHintForLine("Hello world", 0);
  assertEquals(result, null);
});

Deno.test("getBlockHintForLine returns null for closing tag </Callout>", () => {
  const result = getBlockHintForLine("</Callout>", 0);
  assertEquals(result, null);
});

Deno.test('getBlockHintForLine returns hint for <Button href="..."', () => {
  const result = getBlockHintForLine('<Button href="...">', 0);
  assertEquals(result !== null, true);
  assertEquals(result!.blockId, "button");
});

Deno.test("getBlockHintForLine handles leading whitespace", () => {
  const result = getBlockHintForLine('  <Row gap="md">', 0);
  assertEquals(result !== null, true);
  assertEquals(result!.blockId, "row");
});
