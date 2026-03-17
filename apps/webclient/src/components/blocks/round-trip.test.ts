// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertEquals } from "@std/assert";
import { parseBlockProps } from "./parse-block-props.ts";
import { generateBlockTag } from "./generate-block-tag.ts";

Deno.test("round-trip Callout", () => {
  const input = '<Callout variant="info" title="Note">';
  const parsed = parseBlockProps(input, 0);
  assertEquals(parsed !== null, true);

  const generated = generateBlockTag(
    parsed!.componentName,
    parsed!.props,
    parsed!.definition,
    parsed!.selfClosing,
  );

  // Parse again
  const reparsed = parseBlockProps(generated, 0);
  assertEquals(reparsed !== null, true);
  assertEquals(reparsed!.props.variant, "info");
  assertEquals(reparsed!.props.title, "Note");
});

Deno.test("round-trip Spacer self-closing", () => {
  const input = '<Spacer size="md" />';
  const parsed = parseBlockProps(input, 0);
  assertEquals(parsed !== null, true);
  assertEquals(parsed!.selfClosing, true);

  // Edit size to "lg"
  const editedProps = { ...parsed!.props, size: "lg" };
  const generated = generateBlockTag(
    parsed!.componentName,
    editedProps,
    parsed!.definition,
    parsed!.selfClosing,
  );

  // Parse again
  const reparsed = parseBlockProps(generated, 0);
  assertEquals(reparsed !== null, true);
  assertEquals(reparsed!.props.size, "lg");
  assertEquals(reparsed!.selfClosing, true);
});

Deno.test("round-trip Columns with number", () => {
  const input = "<Columns cols={2}>";
  const parsed = parseBlockProps(input, 0);
  assertEquals(parsed !== null, true);
  assertEquals(parsed!.props.cols, "2");

  const generated = generateBlockTag(
    parsed!.componentName,
    parsed!.props,
    parsed!.definition,
    parsed!.selfClosing,
  );

  const reparsed = parseBlockProps(generated, 0);
  assertEquals(reparsed !== null, true);
  assertEquals(reparsed!.props.cols, "2");
});

Deno.test("round-trip multi-line Cover", () => {
  const input = '<Cover\n  src="https://example.com/bg.jpg"\n  overlay="dark">';
  const parsed = parseBlockProps(input, 0);
  assertEquals(parsed !== null, true);
  assertEquals(parsed!.props.src, "https://example.com/bg.jpg");
  assertEquals(parsed!.props.overlay, "dark");

  // Generate will produce single-line
  const generated = generateBlockTag(
    parsed!.componentName,
    parsed!.props,
    parsed!.definition,
    parsed!.selfClosing,
  );

  // Parse the single-line output
  const reparsed = parseBlockProps(generated, 0);
  assertEquals(reparsed !== null, true);
  assertEquals(reparsed!.props.src, "https://example.com/bg.jpg");
  assertEquals(reparsed!.props.overlay, "dark");
});
