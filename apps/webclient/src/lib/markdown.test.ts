/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { formatFrontmatter, formatStoryListItem } from "./markdown.ts";

// formatFrontmatter
Deno.test("formatFrontmatter - simple values", async (t) => {
  await assertSnapshot(
    t,
    formatFrontmatter({
      title: "Hello World",
      date: "2026-01-28",
      draft: false,
      count: 42,
    }),
  );
});

Deno.test("formatFrontmatter - strings needing quotes", async (t) => {
  await assertSnapshot(
    t,
    formatFrontmatter({
      title: 'Title with "quotes"',
      subtitle: "Has: colon in value",
      multiline: "First line\nSecond line",
      padded: " leading space",
    }),
  );
});

Deno.test("formatFrontmatter - null/undefined values skipped", async (t) => {
  await assertSnapshot(
    t,
    formatFrontmatter({
      title: "Present",
      missing: null,
      absent: undefined,
    }),
  );
});

// formatStoryListItem
Deno.test("formatStoryListItem - full story with all fields", async (t) => {
  const story = {
    title: "Microsoft's Vision",
    slug: "20160610-microsoftun-yari-vizyonu",
    summary: "A look at Microsoft's development strategy",
    locale_code: "tr",
    author_profile: { title: "Eser Ozvataf" },
  };
  await assertSnapshot(t, formatStoryListItem(story, "en", "eser/stories"));
});

Deno.test("formatStoryListItem - minimal story", async (t) => {
  const story = {
    title: null,
    slug: null,
    summary: null,
  };
  await assertSnapshot(t, formatStoryListItem(story, "en", "stories"));
});

Deno.test("formatStoryListItem - same locale (no badge)", async (t) => {
  const story = {
    title: "Turkish Article",
    slug: "20260128-turkish-article",
    summary: "Written in Turkish",
    locale_code: "tr",
    author_profile: { title: "Author" },
  };
  await assertSnapshot(t, formatStoryListItem(story, "tr", "stories"));
});

Deno.test("formatStoryListItem - no date in slug", async (t) => {
  const story = {
    title: "About Page",
    slug: "about-us",
    summary: null,
    locale_code: "en",
  };
  await assertSnapshot(t, formatStoryListItem(story, "en", "pages"));
});
