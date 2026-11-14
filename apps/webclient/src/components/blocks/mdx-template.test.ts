// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { escapeForMdxProp, generateContainerMdx, generateSelfClosingMdx } from "./mdx-template.ts";

// escapeForMdxProp tests

test("escapeForMdxProp - normal string", () => {
  expect(escapeForMdxProp("hello world")).toMatchSnapshot();
});

test("escapeForMdxProp - string with quotes", () => {
  expect(escapeForMdxProp('say "hello"')).toMatchSnapshot();
});

test("escapeForMdxProp - string with angle brackets", () => {
  expect(escapeForMdxProp("<script>alert(1)</script>")).toMatchSnapshot();
});

test("escapeForMdxProp - string with braces", () => {
  expect(escapeForMdxProp("{expression}")).toMatchSnapshot();
});

test("escapeForMdxProp - null returns empty", () => {
  expect(escapeForMdxProp(null)).toMatchSnapshot();
});

test("escapeForMdxProp - undefined returns empty", () => {
  expect(escapeForMdxProp(undefined)).toMatchSnapshot();
});

test("escapeForMdxProp - empty string", () => {
  expect(escapeForMdxProp("")).toMatchSnapshot();
});

// generateSelfClosingMdx tests

test("generateSelfClosingMdx - no props", () => {
  expect(generateSelfClosingMdx("Divider", {})).toMatchSnapshot();
});

test("generateSelfClosingMdx - string props", () => {
  expect(generateSelfClosingMdx("Spacer", { size: "md" })).toMatchSnapshot();
});

test("generateSelfClosingMdx - number props", () => {
  expect(generateSelfClosingMdx("Component", { cols: 2 })).toMatchSnapshot();
});

test("generateSelfClosingMdx - boolean props", () => {
  expect(generateSelfClosingMdx("Component", { open: true })).toMatchSnapshot();
});

test("generateSelfClosingMdx - mixed props", () => {
  expect(generateSelfClosingMdx("Card", { title: "Hello", count: 3 })).toMatchSnapshot();
});

test("generateSelfClosingMdx - filters empty strings", () => {
  expect(generateSelfClosingMdx("Test", { name: "foo", empty: "" })).toMatchSnapshot();
});

// generateContainerMdx tests

test("generateContainerMdx - simple", () => {
  expect(generateContainerMdx("Callout", { variant: "info" }, "Your content here.")).toMatchSnapshot();
});

test("generateContainerMdx - no props", () => {
  expect(generateContainerMdx("Column", {}, "Column content.")).toMatchSnapshot();
});

test("generateContainerMdx - multiple props", () => {
  expect(generateContainerMdx("Callout", { variant: "warning", title: "Warning" }, "Be careful!")).toMatchSnapshot();
});
