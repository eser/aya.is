// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { validateBlockProps } from "./validate-props.ts";
import type { BlockDefinition } from "./types.ts";

// Helper to create a minimal BlockDefinition for testing
function mockDefinition(props: BlockDefinition["props"]): BlockDefinition {
  return {
    id: "test-block",
    name: "Test Block",
    description: "A test block",
    icon: () => null,
    category: "text",
    keywords: [],
    props,
    generateMdx: () => "",
  } as unknown as BlockDefinition;
}

test("valid - all required fields present", () => {
  const def = mockDefinition([
    { name: "title", type: "string", label: "Title", required: true, defaultValue: "" },
    { name: "description", type: "string", label: "Description", required: true, defaultValue: "" },
  ]);
  const result = validateBlockProps(def, { title: "Hello", description: "World" });
  expect(result).toMatchSnapshot();
});

test("invalid - required field missing", () => {
  const def = mockDefinition([{ name: "title", type: "string", label: "Title", required: true, defaultValue: "" }]);
  const result = validateBlockProps(def, {});
  expect(result).toMatchSnapshot();
});

test("invalid - required field empty string", () => {
  const def = mockDefinition([{ name: "title", type: "string", label: "Title", required: true, defaultValue: "" }]);
  const result = validateBlockProps(def, { title: "" });
  expect(result).toMatchSnapshot();
});

test("valid - optional field missing", () => {
  const def = mockDefinition([
    { name: "subtitle", type: "string", label: "Subtitle", required: false, defaultValue: "" },
  ]);
  const result = validateBlockProps(def, {});
  expect(result).toMatchSnapshot();
});

test("invalid - number field with non-numeric", () => {
  const def = mockDefinition([{ name: "count", type: "number", label: "Count", required: false, defaultValue: 0 }]);
  const result = validateBlockProps(def, { count: "abc" });
  expect(result).toMatchSnapshot();
});

test("valid - number field with numeric string", () => {
  const def = mockDefinition([{ name: "count", type: "number", label: "Count", required: false, defaultValue: 0 }]);
  const result = validateBlockProps(def, { count: "42" });
  expect(result).toMatchSnapshot();
});

test("invalid - boolean field with invalid value", () => {
  const def = mockDefinition([{ name: "open", type: "boolean", label: "Open", required: false, defaultValue: false }]);
  const result = validateBlockProps(def, { open: "maybe" });
  expect(result).toMatchSnapshot();
});

test("valid - boolean field with true/false string", () => {
  const def = mockDefinition([{ name: "open", type: "boolean", label: "Open", required: false, defaultValue: false }]);
  const result = validateBlockProps(def, { open: "true" });
  expect(result).toMatchSnapshot();
});

test("invalid - select field with invalid option", () => {
  const def = mockDefinition([
    {
      name: "variant",
      type: "select",
      label: "Variant",
      required: false,
      defaultValue: "info",
      options: [
        { value: "info", label: "Info" },
        { value: "warning", label: "Warning" },
      ],
    },
  ]);
  const result = validateBlockProps(def, { variant: "danger" });
  expect(result).toMatchSnapshot();
});

test("valid - select field with valid option", () => {
  const def = mockDefinition([
    {
      name: "variant",
      type: "select",
      label: "Variant",
      required: false,
      defaultValue: "info",
      options: [
        { value: "info", label: "Info" },
        { value: "warning", label: "Warning" },
      ],
    },
  ]);
  const result = validateBlockProps(def, { variant: "info" });
  expect(result).toMatchSnapshot();
});

test("multiple errors returned", () => {
  const def = mockDefinition([
    { name: "title", type: "string", label: "Title", required: true, defaultValue: "" },
    { name: "count", type: "number", label: "Count", required: false, defaultValue: 0 },
    { name: "open", type: "boolean", label: "Open", required: false, defaultValue: false },
  ]);
  const result = validateBlockProps(def, { count: "abc", open: "maybe" });
  expect(result).toMatchSnapshot();
});

test("all optional empty is valid", () => {
  const def = mockDefinition([
    { name: "a", type: "string", label: "A", required: false, defaultValue: "" },
    { name: "b", type: "number", label: "B", required: false, defaultValue: 0 },
    { name: "c", type: "boolean", label: "C", required: false, defaultValue: false },
  ]);
  const result = validateBlockProps(def, {});
  expect(result).toMatchSnapshot();
});
