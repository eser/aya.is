// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { checkSlugSchema, createProfileSchema, updateProfileSchema } from "./profile.ts";

// Snapshot schema shapes to detect validation rule changes
test("createProfileSchema - shape snapshot", () => {
  expect(Object.keys(createProfileSchema.shape)).toMatchSnapshot();
});

test("updateProfileSchema - shape snapshot", () => {
  expect(Object.keys(updateProfileSchema.shape)).toMatchSnapshot();
});

test("checkSlugSchema - shape snapshot", () => {
  expect(Object.keys(checkSlugSchema.shape)).toMatchSnapshot();
});

// createProfileSchema - valid input
test("createProfileSchema - valid input passes", () => {
  const result = createProfileSchema.safeParse({
    slug: "my-profile",
    title: "My Profile",
    description: "A test profile",
    kind: "individual",
  });
  expect(result.success).toEqual(true);
});

// createProfileSchema - invalid cases
test("createProfileSchema - slug too short", () => {
  const result = createProfileSchema.safeParse({
    slug: "a",
    title: "Title",
    kind: "individual",
  });
  expect(result.success).toEqual(false);
  if (!result.success) {
    expect(result.error.issues.map((i) => ({ path: i.path, message: i.message }))).toMatchSnapshot();
  }
});

test("createProfileSchema - slug with special chars", () => {
  const result = createProfileSchema.safeParse({
    slug: "My Profile!",
    title: "Title",
    kind: "individual",
  });
  expect(result.success).toEqual(false);
  if (!result.success) {
    expect(result.error.issues.map((i) => ({ path: i.path, message: i.message }))).toMatchSnapshot();
  }
});

test("createProfileSchema - missing title", () => {
  const result = createProfileSchema.safeParse({
    slug: "valid-slug",
    title: "",
    kind: "individual",
  });
  expect(result.success).toEqual(false);
  if (!result.success) {
    expect(result.error.issues.map((i) => ({ path: i.path, message: i.message }))).toMatchSnapshot();
  }
});

test("createProfileSchema - invalid kind", () => {
  const result = createProfileSchema.safeParse({
    slug: "valid-slug",
    title: "Title",
    kind: "invalid",
  });
  expect(result.success).toEqual(false);
  if (!result.success) {
    expect(result.error.issues.map((i) => ({ path: i.path, message: i.message }))).toMatchSnapshot();
  }
});

test("createProfileSchema - valid kinds", () => {
  for (const kind of ["individual", "organization", "product"]) {
    const result = createProfileSchema.safeParse({
      slug: "valid-slug",
      title: "Title",
      kind,
    });
    expect(result.success, `Expected kind "${kind}" to be valid`).toEqual(true);
  }
});

// updateProfileSchema - all optional
test("updateProfileSchema - empty input passes", () => {
  const result = updateProfileSchema.safeParse({});
  expect(result.success).toEqual(true);
});

test("updateProfileSchema - title too long", () => {
  const result = updateProfileSchema.safeParse({
    title: "T".repeat(101),
  });
  expect(result.success).toEqual(false);
  if (!result.success) {
    expect(result.error.issues.map((i) => ({ path: i.path, message: i.message }))).toMatchSnapshot();
  }
});
