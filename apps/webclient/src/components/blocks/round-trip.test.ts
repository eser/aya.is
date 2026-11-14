// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { parseBlockProps } from "./parse-block-props.ts";
import { generateBlockTag } from "./generate-block-tag.ts";

test("round-trip Callout", () => {
  const input = '<Callout variant="info" title="Note">';
  const parsed = parseBlockProps(input, 0);
  expect(parsed !== null).toEqual(true);

  const generated = generateBlockTag(parsed!.componentName, parsed!.props, parsed!.definition, parsed!.selfClosing);

  // Parse again
  const reparsed = parseBlockProps(generated, 0);
  expect(reparsed !== null).toEqual(true);
  expect(reparsed!.props.variant).toEqual("info");
  expect(reparsed!.props.title).toEqual("Note");
});

test("round-trip Spacer self-closing", () => {
  const input = '<Spacer size="md" />';
  const parsed = parseBlockProps(input, 0);
  expect(parsed !== null).toEqual(true);
  expect(parsed!.selfClosing).toEqual(true);

  // Edit size to "lg"
  const editedProps = { ...parsed!.props, size: "lg" };
  const generated = generateBlockTag(parsed!.componentName, editedProps, parsed!.definition, parsed!.selfClosing);

  // Parse again
  const reparsed = parseBlockProps(generated, 0);
  expect(reparsed !== null).toEqual(true);
  expect(reparsed!.props.size).toEqual("lg");
  expect(reparsed!.selfClosing).toEqual(true);
});

test("round-trip Columns with number", () => {
  const input = "<Columns cols={2}>";
  const parsed = parseBlockProps(input, 0);
  expect(parsed !== null).toEqual(true);
  expect(parsed!.props.cols).toEqual("2");

  const generated = generateBlockTag(parsed!.componentName, parsed!.props, parsed!.definition, parsed!.selfClosing);

  const reparsed = parseBlockProps(generated, 0);
  expect(reparsed !== null).toEqual(true);
  expect(reparsed!.props.cols).toEqual("2");
});

test("round-trip multi-line Cover", () => {
  const input = '<Cover\n  src="https://example.com/bg.jpg"\n  overlay="dark">';
  const parsed = parseBlockProps(input, 0);
  expect(parsed !== null).toEqual(true);
  expect(parsed!.props.src).toEqual("https://example.com/bg.jpg");
  expect(parsed!.props.overlay).toEqual("dark");

  // Generate will produce single-line
  const generated = generateBlockTag(parsed!.componentName, parsed!.props, parsed!.definition, parsed!.selfClosing);

  // Parse the single-line output
  const reparsed = parseBlockProps(generated, 0);
  expect(reparsed !== null).toEqual(true);
  expect(reparsed!.props.src).toEqual("https://example.com/bg.jpg");
  expect(reparsed!.props.overlay).toEqual("dark");
});
