// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { parseBlockProps } from "./parse-block-props.ts";

test("single-line tag with string props", () => {
  const result = parseBlockProps('<Callout variant="info" title="Note">', 0);
  expect(result !== null).toEqual(true);
  expect({
    componentName: result!.componentName,
    blockId: result!.blockId,
    props: result!.props,
    selfClosing: result!.selfClosing,
    startOffset: result!.startOffset,
    endOffset: result!.endOffset,
  }).toMatchSnapshot();
});

test("single-line self-closing tag", () => {
  const result = parseBlockProps('<Spacer size="md" />', 0);
  expect(result !== null).toEqual(true);
  expect(result!.selfClosing).toEqual(true);
  expect({
    componentName: result!.componentName,
    blockId: result!.blockId,
    props: result!.props,
    selfClosing: result!.selfClosing,
    startOffset: result!.startOffset,
    endOffset: result!.endOffset,
  }).toMatchSnapshot();
});

test("JSX expression prop", () => {
  const result = parseBlockProps("<Columns cols={2}>", 0);
  expect(result !== null).toEqual(true);
  expect(result!.props.cols).toEqual("2");
});

test("bare boolean prop", () => {
  const result = parseBlockProps("<Details open>", 0);
  expect(result !== null).toEqual(true);
  expect(result!.props.open).toEqual("true");
});

test("mixed props", () => {
  const result = parseBlockProps('<Cover src="url" overlay="dark">', 0);
  expect(result !== null).toEqual(true);
  expect({
    componentName: result!.componentName,
    blockId: result!.blockId,
    props: result!.props,
    selfClosing: result!.selfClosing,
    startOffset: result!.startOffset,
    endOffset: result!.endOffset,
  }).toMatchSnapshot();
});

test("no tag", () => {
  const result = parseBlockProps("Hello world", 0);
  expect(result).toEqual(null);
});

test("closing tag", () => {
  const result = parseBlockProps("</Callout>", 0);
  expect(result).toEqual(null);
});

test("multi-line tag", () => {
  const result = parseBlockProps('<Cover\n  src="url"\n  overlay="dark">', 0);
  expect(result !== null).toEqual(true);
  expect(result!.props.src).toEqual("url");
  expect(result!.props.overlay).toEqual("dark");
});

test("empty props", () => {
  const result = parseBlockProps("<Verse>", 0);
  expect(result !== null).toEqual(true);
  expect(Object.keys(result!.props).length).toEqual(0);
});

test("tag with content before it", () => {
  const result = parseBlockProps('some text\n<Callout variant="info">', 1);
  expect(result !== null).toEqual(true);
  expect(result!.blockId).toEqual("callout");
  expect(result!.props.variant).toEqual("info");
});

test("unknown component", () => {
  const result = parseBlockProps('<UnknownThing prop="val">', 0);
  expect(result).toEqual(null);
});
