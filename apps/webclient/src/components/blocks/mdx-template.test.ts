/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import {
  escapeForMdxProp,
  generateContainerMdx,
  generateSelfClosingMdx,
} from "./mdx-template.ts";

// escapeForMdxProp tests

Deno.test("escapeForMdxProp - normal string", async (t) => {
  await assertSnapshot(t, escapeForMdxProp("hello world"));
});

Deno.test("escapeForMdxProp - string with quotes", async (t) => {
  await assertSnapshot(t, escapeForMdxProp('say "hello"'));
});

Deno.test("escapeForMdxProp - string with angle brackets", async (t) => {
  await assertSnapshot(t, escapeForMdxProp("<script>alert(1)</script>"));
});

Deno.test("escapeForMdxProp - string with braces", async (t) => {
  await assertSnapshot(t, escapeForMdxProp("{expression}"));
});

Deno.test("escapeForMdxProp - null returns empty", async (t) => {
  await assertSnapshot(t, escapeForMdxProp(null));
});

Deno.test("escapeForMdxProp - undefined returns empty", async (t) => {
  await assertSnapshot(t, escapeForMdxProp(undefined));
});

Deno.test("escapeForMdxProp - empty string", async (t) => {
  await assertSnapshot(t, escapeForMdxProp(""));
});

// generateSelfClosingMdx tests

Deno.test("generateSelfClosingMdx - no props", async (t) => {
  await assertSnapshot(t, generateSelfClosingMdx("Divider", {}));
});

Deno.test("generateSelfClosingMdx - string props", async (t) => {
  await assertSnapshot(t, generateSelfClosingMdx("Spacer", { size: "md" }));
});

Deno.test("generateSelfClosingMdx - number props", async (t) => {
  await assertSnapshot(t, generateSelfClosingMdx("Component", { cols: 2 }));
});

Deno.test("generateSelfClosingMdx - boolean props", async (t) => {
  await assertSnapshot(t, generateSelfClosingMdx("Component", { open: true }));
});

Deno.test("generateSelfClosingMdx - mixed props", async (t) => {
  await assertSnapshot(t, generateSelfClosingMdx("Card", { title: "Hello", count: 3 }));
});

Deno.test("generateSelfClosingMdx - filters empty strings", async (t) => {
  await assertSnapshot(t, generateSelfClosingMdx("Test", { name: "foo", empty: "" }));
});

// generateContainerMdx tests

Deno.test("generateContainerMdx - simple", async (t) => {
  await assertSnapshot(t, generateContainerMdx("Callout", { variant: "info" }, "Your content here."));
});

Deno.test("generateContainerMdx - no props", async (t) => {
  await assertSnapshot(t, generateContainerMdx("Column", {}, "Column content."));
});

Deno.test("generateContainerMdx - multiple props", async (t) => {
  await assertSnapshot(t, generateContainerMdx("Callout", { variant: "warning", title: "Warning" }, "Be careful!"));
});
