/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals } from "@std/assert";
import {
  DEFAULT_LOCALE,
  FALLBACK_LOCALE,
  getLocaleData,
  isAllowedURI,
  isValidLocale,
  SUPPORTED_LOCALES,
  supportedLocales,
} from "./locale-utils.ts";

// Snapshot the supported locales array — catches accidental additions/removals
Deno.test("SUPPORTED_LOCALES - complete list", async (t) => {
  await assertSnapshot(t, SUPPORTED_LOCALES);
});

// Snapshot the full locale record structure — catches property changes
Deno.test("supportedLocales - all entries", async (t) => {
  await assertSnapshot(t, supportedLocales);
});

Deno.test("DEFAULT_LOCALE and FALLBACK_LOCALE", () => {
  assertEquals(DEFAULT_LOCALE, "en");
  assertEquals(FALLBACK_LOCALE, "en");
});

// isValidLocale
Deno.test("isValidLocale - all valid locales return true", () => {
  for (const locale of SUPPORTED_LOCALES) {
    assertEquals(isValidLocale(locale), true, `Expected ${locale} to be valid`);
  }
});

Deno.test("isValidLocale - invalid locales return false", () => {
  assertEquals(isValidLocale("xx"), false);
  assertEquals(isValidLocale(""), false);
  assertEquals(isValidLocale("EN"), false); // case-sensitive
  assertEquals(isValidLocale("pt"), false); // only pt-PT is valid
  assertEquals(isValidLocale("zh"), false); // only zh-CN is valid
});

// getLocaleData
Deno.test("getLocaleData - returns data for valid locale", async (t) => {
  await assertSnapshot(t, getLocaleData("en"));
  await assertSnapshot(t, getLocaleData("ar")); // RTL locale
  await assertSnapshot(t, getLocaleData("pt-PT")); // hyphenated locale
});

Deno.test("getLocaleData - returns undefined for invalid locale", () => {
  assertEquals(getLocaleData("invalid"), undefined);
  assertEquals(getLocaleData(""), undefined);
});

// isAllowedURI
Deno.test("isAllowedURI - null/empty/undefined are allowed", () => {
  assertEquals(isAllowedURI(null, ["https://example.com/"]), true);
  assertEquals(isAllowedURI(undefined, ["https://example.com/"]), true);
  assertEquals(isAllowedURI("", ["https://example.com/"]), true);
});

Deno.test("isAllowedURI - matching prefix", () => {
  const prefixes = ["https://objects.aya.is/", "https://avatars.githubusercontent.com/"];
  assertEquals(isAllowedURI("https://objects.aya.is/image.png", prefixes), true);
  assertEquals(isAllowedURI("https://avatars.githubusercontent.com/u/123", prefixes), true);
});

Deno.test("isAllowedURI - non-matching prefix", () => {
  const prefixes = ["https://objects.aya.is/"];
  assertEquals(isAllowedURI("https://evil.com/hack.png", prefixes), false);
  assertEquals(isAllowedURI("http://objects.aya.is/image.png", prefixes), false); // http vs https
});

Deno.test("isAllowedURI - empty prefix list allows all", () => {
  assertEquals(isAllowedURI("https://anything.com/file.png", []), true);
});
