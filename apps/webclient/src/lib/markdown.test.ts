// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { formatFrontmatter, formatStoryListItem } from "./markdown.ts";

// formatFrontmatter
test("formatFrontmatter - simple values", () => {
  expect(
    formatFrontmatter({
      title: "Hello World",
      date: "2026-01-28",
      draft: false,
      count: 42,
    }),
  ).toMatchSnapshot();
});

test("formatFrontmatter - strings needing quotes", () => {
  expect(
    formatFrontmatter({
      title: 'Title with "quotes"',
      subtitle: "Has: colon in value",
      multiline: "First line\nSecond line",
      padded: " leading space",
    }),
  ).toMatchSnapshot();
});

test("formatFrontmatter - null/undefined values skipped", () => {
  expect(
    formatFrontmatter({
      title: "Present",
      missing: null,
      absent: undefined,
    }),
  ).toMatchSnapshot();
});

// formatStoryListItem
test("formatStoryListItem - full story with all fields", () => {
  const story = {
    title: "Microsoft's Vision",
    slug: "20160610-microsoftun-yari-vizyonu",
    summary: "A look at Microsoft's development strategy",
    locale_code: "tr",
    author_profile: { title: "Eser Ozvataf" },
  };
  expect(formatStoryListItem(story, "en", "eser/stories")).toMatchSnapshot();
});

test("formatStoryListItem - minimal story", () => {
  const story = {
    title: null,
    slug: null,
    summary: null,
  };
  expect(formatStoryListItem(story, "en", "stories")).toMatchSnapshot();
});

test("formatStoryListItem - same locale (no badge)", () => {
  const story = {
    title: "Turkish Article",
    slug: "20260128-turkish-article",
    summary: "Written in Turkish",
    locale_code: "tr",
    author_profile: { title: "Author" },
  };
  expect(formatStoryListItem(story, "tr", "stories")).toMatchSnapshot();
});

test("formatStoryListItem - no date in slug", () => {
  const story = {
    title: "About Page",
    slug: "about-us",
    summary: null,
    locale_code: "en",
  };
  expect(formatStoryListItem(story, "en", "pages")).toMatchSnapshot();
});
