// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
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
test("generateMetaTags - basic website", () => {
  expect(
    generateMetaTags(
      {
        title: "Home",
        description: "Welcome to AYA",
      },
      TEST_DEFAULTS,
    ),
  ).toMatchSnapshot();
});

test("generateMetaTags - article with dates and author", () => {
  expect(
    generateMetaTags(
      {
        title: "My Article",
        description: "An article about testing",
        type: "article",
        locale: "tr",
        url: "https://aya.is/tr/stories/my-article",
        publishedTime: "2026-01-28T14:30:00Z",
        modifiedTime: "2026-02-01T10:00:00Z",
        author: "Eser",
      },
      TEST_DEFAULTS,
    ),
  ).toMatchSnapshot();
});

test("generateMetaTags - noIndex", () => {
  expect(
    generateMetaTags(
      {
        title: "Draft",
        description: "Not for indexing",
        noIndex: true,
      },
      TEST_DEFAULTS,
    ),
  ).toMatchSnapshot();
});

test("generateMetaTags - custom image", () => {
  expect(
    generateMetaTags(
      {
        title: "With Image",
        description: "Has custom image",
        image: "https://aya.is/custom.png",
      },
      TEST_DEFAULTS,
    ),
  ).toMatchSnapshot();
});

test("generateMetaTags - null image", () => {
  expect(
    generateMetaTags(
      {
        title: "No Image",
        description: "Image is null",
        image: null,
      },
      TEST_DEFAULTS,
    ),
  ).toMatchSnapshot();
});

test("generateMetaTags - title matches siteName (no duplication)", () => {
  expect(
    generateMetaTags(
      {
        title: "AYA",
        description: "Main page",
        siteName: "AYA",
      },
      TEST_DEFAULTS,
    ),
  ).toMatchSnapshot();
});

// buildUrl
test("buildUrl - various combinations", () => {
  expect(buildUrl("https://aya.is", "en")).toMatchSnapshot();
  expect(buildUrl("https://aya.is", "tr", "stories")).toMatchSnapshot();
  expect(buildUrl("https://aya.is", "en", "eser", "stories", "my-story")).toMatchSnapshot();
  expect(buildUrl("https://aya.is", "fr", "", "about")).toMatchSnapshot(); // empty segment filtered
});

// generateCanonicalLink
test("generateCanonicalLink", () => {
  expect(generateCanonicalLink("https://aya.is/en/about")).toMatchSnapshot();
});

// computeStoryCanonicalUrl
test("computeStoryCanonicalUrl - managed story with source_url", () => {
  const story = {
    is_managed: true,
    properties: { source_url: "https://medium.com/original-article" },
    locale_code: "en",
    slug: "my-story",
  };
  expect(computeStoryCanonicalUrl("https://aya.is", story, "en", "eser", "stories")).toMatchSnapshot();
});

test("computeStoryCanonicalUrl - regular story", () => {
  const story = {
    is_managed: false,
    properties: null,
    locale_code: "tr",
    slug: "20260128-test-story",
  };
  expect(computeStoryCanonicalUrl("https://aya.is", story, "en", "eser", "stories")).toMatchSnapshot();
});

test("computeStoryCanonicalUrl - fallback to viewer locale", () => {
  const story = {
    is_managed: false,
    properties: null,
    locale_code: "  ",
    slug: "my-story",
  };
  expect(computeStoryCanonicalUrl("https://aya.is", story, "en", "stories")).toMatchSnapshot();
});

// computeContentLanguage
test("computeContentLanguage - same locale", () => {
  expect(computeContentLanguage("en", "en")).toMatchSnapshot();
});

test("computeContentLanguage - different locales", () => {
  expect(computeContentLanguage("en", "tr")).toMatchSnapshot();
});

test("computeContentLanguage - undefined content locale", () => {
  expect(computeContentLanguage("en", undefined)).toMatchSnapshot();
});

// truncateDescription
test("truncateDescription - various cases", () => {
  const fallback = "Default description";
  expect(truncateDescription(null, fallback)).toMatchSnapshot();
  expect(truncateDescription(undefined, fallback)).toMatchSnapshot();
  expect(truncateDescription("", fallback)).toMatchSnapshot();
  expect(truncateDescription("Short text", fallback)).toMatchSnapshot();
  expect(truncateDescription("A".repeat(160), fallback)).toMatchSnapshot(); // exactly 160
  expect(truncateDescription("B".repeat(200), fallback)).toMatchSnapshot(); // truncated
});
