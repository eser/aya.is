// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import {
  buildUrl,
  computeContentLanguage,
  computeStoryCanonicalUrl,
  generateCanonicalLink,
  generateMetaTags,
  truncateDescription,
} from "./seo-utils.ts";

const TEST_DEFAULTS = {
  host: "https://aya.is",
  name: "AYA",
  defaultImage: "https://aya.is/og-image.png",
};

// generateMetaTags
Deno.test("generateMetaTags - basic website", async (t) => {
  await assertSnapshot(
    t,
    generateMetaTags({
      title: "Home",
      description: "Welcome to AYA",
    }, TEST_DEFAULTS),
  );
});

Deno.test("generateMetaTags - article with dates and author", async (t) => {
  await assertSnapshot(
    t,
    generateMetaTags({
      title: "My Article",
      description: "An article about testing",
      type: "article",
      locale: "tr",
      url: "https://aya.is/tr/stories/my-article",
      publishedTime: "2026-01-28T14:30:00Z",
      modifiedTime: "2026-02-01T10:00:00Z",
      author: "Eser",
    }, TEST_DEFAULTS),
  );
});

Deno.test("generateMetaTags - noIndex", async (t) => {
  await assertSnapshot(
    t,
    generateMetaTags({
      title: "Draft",
      description: "Not for indexing",
      noIndex: true,
    }, TEST_DEFAULTS),
  );
});

Deno.test("generateMetaTags - custom image", async (t) => {
  await assertSnapshot(
    t,
    generateMetaTags({
      title: "With Image",
      description: "Has custom image",
      image: "https://aya.is/custom.png",
    }, TEST_DEFAULTS),
  );
});

Deno.test("generateMetaTags - null image", async (t) => {
  await assertSnapshot(
    t,
    generateMetaTags({
      title: "No Image",
      description: "Image is null",
      image: null,
    }, TEST_DEFAULTS),
  );
});

Deno.test("generateMetaTags - title matches siteName (no duplication)", async (t) => {
  await assertSnapshot(
    t,
    generateMetaTags({
      title: "AYA",
      description: "Main page",
      siteName: "AYA",
    }, TEST_DEFAULTS),
  );
});

// buildUrl
Deno.test("buildUrl - various combinations", async (t) => {
  await assertSnapshot(t, buildUrl("https://aya.is", "en"));
  await assertSnapshot(t, buildUrl("https://aya.is", "tr", "stories"));
  await assertSnapshot(t, buildUrl("https://aya.is", "en", "eser", "stories", "my-story"));
  await assertSnapshot(t, buildUrl("https://aya.is", "fr", "", "about")); // empty segment filtered
});

// generateCanonicalLink
Deno.test("generateCanonicalLink", async (t) => {
  await assertSnapshot(t, generateCanonicalLink("https://aya.is/en/about"));
});

// computeStoryCanonicalUrl
Deno.test("computeStoryCanonicalUrl - managed story with source_url", async (t) => {
  const story = {
    is_managed: true,
    properties: { source_url: "https://medium.com/original-article" },
    locale_code: "en",
    slug: "my-story",
  };
  await assertSnapshot(t, computeStoryCanonicalUrl("https://aya.is", story, "en", "eser", "stories"));
});

Deno.test("computeStoryCanonicalUrl - regular story", async (t) => {
  const story = {
    is_managed: false,
    properties: null,
    locale_code: "tr",
    slug: "20260128-test-story",
  };
  await assertSnapshot(t, computeStoryCanonicalUrl("https://aya.is", story, "en", "eser", "stories"));
});

Deno.test("computeStoryCanonicalUrl - fallback to viewer locale", async (t) => {
  const story = {
    is_managed: false,
    properties: null,
    locale_code: "  ",
    slug: "my-story",
  };
  await assertSnapshot(t, computeStoryCanonicalUrl("https://aya.is", story, "en", "stories"));
});

// computeContentLanguage
Deno.test("computeContentLanguage - same locale", async (t) => {
  await assertSnapshot(t, computeContentLanguage("en", "en"));
});

Deno.test("computeContentLanguage - different locales", async (t) => {
  await assertSnapshot(t, computeContentLanguage("en", "tr"));
});

Deno.test("computeContentLanguage - undefined content locale", async (t) => {
  await assertSnapshot(t, computeContentLanguage("en", undefined));
});

// truncateDescription
Deno.test("truncateDescription - various cases", async (t) => {
  const fallback = "Default description";
  await assertSnapshot(t, truncateDescription(null, fallback));
  await assertSnapshot(t, truncateDescription(undefined, fallback));
  await assertSnapshot(t, truncateDescription("", fallback));
  await assertSnapshot(t, truncateDescription("Short text", fallback));
  await assertSnapshot(t, truncateDescription("A".repeat(160), fallback)); // exactly 160
  await assertSnapshot(t, truncateDescription("B".repeat(200), fallback)); // truncated
});
