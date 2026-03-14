/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
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

Deno.test("valid - all required fields present", async (t) => {
  const def = mockDefinition([
    { name: "title", type: "string", label: "Title", required: true, defaultValue: "" },
    { name: "description", type: "string", label: "Description", required: true, defaultValue: "" },
  ]);
  const result = validateBlockProps(def, { title: "Hello", description: "World" });
  await assertSnapshot(t, result);
});

Deno.test("invalid - required field missing", async (t) => {
  const def = mockDefinition([
    { name: "title", type: "string", label: "Title", required: true, defaultValue: "" },
  ]);
  const result = validateBlockProps(def, {});
  await assertSnapshot(t, result);
});

Deno.test("invalid - required field empty string", async (t) => {
  const def = mockDefinition([
    { name: "title", type: "string", label: "Title", required: true, defaultValue: "" },
  ]);
  const result = validateBlockProps(def, { title: "" });
  await assertSnapshot(t, result);
});

Deno.test("valid - optional field missing", async (t) => {
  const def = mockDefinition([
    { name: "subtitle", type: "string", label: "Subtitle", required: false, defaultValue: "" },
  ]);
  const result = validateBlockProps(def, {});
  await assertSnapshot(t, result);
});

Deno.test("invalid - number field with non-numeric", async (t) => {
  const def = mockDefinition([
    { name: "count", type: "number", label: "Count", required: false, defaultValue: 0 },
  ]);
  const result = validateBlockProps(def, { count: "abc" });
  await assertSnapshot(t, result);
});

Deno.test("valid - number field with numeric string", async (t) => {
  const def = mockDefinition([
    { name: "count", type: "number", label: "Count", required: false, defaultValue: 0 },
  ]);
  const result = validateBlockProps(def, { count: "42" });
  await assertSnapshot(t, result);
});

Deno.test("invalid - boolean field with invalid value", async (t) => {
  const def = mockDefinition([
    { name: "open", type: "boolean", label: "Open", required: false, defaultValue: false },
  ]);
  const result = validateBlockProps(def, { open: "maybe" });
  await assertSnapshot(t, result);
});

Deno.test("valid - boolean field with true/false string", async (t) => {
  const def = mockDefinition([
    { name: "open", type: "boolean", label: "Open", required: false, defaultValue: false },
  ]);
  const result = validateBlockProps(def, { open: "true" });
  await assertSnapshot(t, result);
});

Deno.test("invalid - select field with invalid option", async (t) => {
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
  await assertSnapshot(t, result);
});

Deno.test("valid - select field with valid option", async (t) => {
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
  await assertSnapshot(t, result);
});

Deno.test("multiple errors returned", async (t) => {
  const def = mockDefinition([
    { name: "title", type: "string", label: "Title", required: true, defaultValue: "" },
    { name: "count", type: "number", label: "Count", required: false, defaultValue: 0 },
    { name: "open", type: "boolean", label: "Open", required: false, defaultValue: false },
  ]);
  const result = validateBlockProps(def, { count: "abc", open: "maybe" });
  await assertSnapshot(t, result);
});

Deno.test("all optional empty is valid", async (t) => {
  const def = mockDefinition([
    { name: "a", type: "string", label: "A", required: false, defaultValue: "" },
    { name: "b", type: "number", label: "B", required: false, defaultValue: 0 },
    { name: "c", type: "boolean", label: "C", required: false, defaultValue: false },
  ]);
  const result = validateBlockProps(def, {});
  await assertSnapshot(t, result);
});
