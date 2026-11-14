// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { getBlockHintForLine } from "./block-hint.ts";

test("getBlockHintForLine returns hint for <Callout", () => {
  const result = getBlockHintForLine("<Callout", 0);
  expect(result !== null).toEqual(true);
  expect(result!.blockId).toEqual("callout");
});

test("getBlockHintForLine returns null for <UnknownComponent", () => {
  const result = getBlockHintForLine("<UnknownComponent", 0);
  expect(result).toEqual(null);
});

test("getBlockHintForLine returns null for regular text", () => {
  const result = getBlockHintForLine("Hello world", 0);
  expect(result).toEqual(null);
});

test("getBlockHintForLine returns null for closing tag </Callout>", () => {
  const result = getBlockHintForLine("</Callout>", 0);
  expect(result).toEqual(null);
});

test('getBlockHintForLine returns hint for <Button href="..."', () => {
  const result = getBlockHintForLine('<Button href="...">', 0);
  expect(result !== null).toEqual(true);
  expect(result!.blockId).toEqual("button");
});

test("getBlockHintForLine handles leading whitespace", () => {
  const result = getBlockHintForLine('  <Row gap="md">', 0);
  expect(result !== null).toEqual(true);
  expect(result!.blockId).toEqual("row");
});
