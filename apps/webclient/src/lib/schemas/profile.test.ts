// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals } from "@std/assert";
import { checkSlugSchema, createProfileSchema, updateProfileSchema } from "./profile.ts";

// Snapshot schema shapes to detect validation rule changes
Deno.test("createProfileSchema - shape snapshot", async (t) => {
  await assertSnapshot(t, Object.keys(createProfileSchema.shape));
});

Deno.test("updateProfileSchema - shape snapshot", async (t) => {
  await assertSnapshot(t, Object.keys(updateProfileSchema.shape));
});

Deno.test("checkSlugSchema - shape snapshot", async (t) => {
  await assertSnapshot(t, Object.keys(checkSlugSchema.shape));
});

// createProfileSchema - valid input
Deno.test("createProfileSchema - valid input passes", () => {
  const result = createProfileSchema.safeParse({
    slug: "my-profile",
    title: "My Profile",
    description: "A test profile",
    kind: "individual",
  });
  assertEquals(result.success, true);
});

// createProfileSchema - invalid cases
Deno.test("createProfileSchema - slug too short", async (t) => {
  const result = createProfileSchema.safeParse({
    slug: "a",
    title: "Title",
    kind: "individual",
  });
  assertEquals(result.success, false);
  if (!result.success) {
    await assertSnapshot(t, result.error.issues.map((i) => ({ path: i.path, message: i.message })));
  }
});

Deno.test("createProfileSchema - slug with special chars", async (t) => {
  const result = createProfileSchema.safeParse({
    slug: "My Profile!",
    title: "Title",
    kind: "individual",
  });
  assertEquals(result.success, false);
  if (!result.success) {
    await assertSnapshot(t, result.error.issues.map((i) => ({ path: i.path, message: i.message })));
  }
});

Deno.test("createProfileSchema - missing title", async (t) => {
  const result = createProfileSchema.safeParse({
    slug: "valid-slug",
    title: "",
    kind: "individual",
  });
  assertEquals(result.success, false);
  if (!result.success) {
    await assertSnapshot(t, result.error.issues.map((i) => ({ path: i.path, message: i.message })));
  }
});

Deno.test("createProfileSchema - invalid kind", async (t) => {
  const result = createProfileSchema.safeParse({
    slug: "valid-slug",
    title: "Title",
    kind: "invalid",
  });
  assertEquals(result.success, false);
  if (!result.success) {
    await assertSnapshot(t, result.error.issues.map((i) => ({ path: i.path, message: i.message })));
  }
});

Deno.test("createProfileSchema - valid kinds", () => {
  for (const kind of ["individual", "organization", "product"]) {
    const result = createProfileSchema.safeParse({
      slug: "valid-slug",
      title: "Title",
      kind,
    });
    assertEquals(result.success, true, `Expected kind "${kind}" to be valid`);
  }
});

// updateProfileSchema - all optional
Deno.test("updateProfileSchema - empty input passes", () => {
  const result = updateProfileSchema.safeParse({});
  assertEquals(result.success, true);
});

Deno.test("updateProfileSchema - title too long", async (t) => {
  const result = updateProfileSchema.safeParse({
    title: "T".repeat(101),
  });
  assertEquals(result.success, false);
  if (!result.success) {
    await assertSnapshot(t, result.error.issues.map((i) => ({ path: i.path, message: i.message })));
  }
});
