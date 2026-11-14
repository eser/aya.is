// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
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
test("SUPPORTED_LOCALES - complete list", () => {
  expect(SUPPORTED_LOCALES).toMatchSnapshot();
});

// Snapshot the full locale record structure — catches property changes
test("supportedLocales - all entries", () => {
  expect(supportedLocales).toMatchSnapshot();
});

test("DEFAULT_LOCALE and FALLBACK_LOCALE", () => {
  expect(DEFAULT_LOCALE).toEqual("en");
  expect(FALLBACK_LOCALE).toEqual("en");
});

// isValidLocale
test("isValidLocale - all valid locales return true", () => {
  for (const locale of SUPPORTED_LOCALES) {
    expect(isValidLocale(locale), `Expected ${locale} to be valid`).toEqual(true);
  }
});

test("isValidLocale - invalid locales return false", () => {
  expect(isValidLocale("xx")).toEqual(false);
  expect(isValidLocale("")).toEqual(false);
  expect(isValidLocale("EN")).toEqual(false); // case-sensitive
  expect(isValidLocale("pt")).toEqual(false); // only pt-PT is valid
  expect(isValidLocale("zh")).toEqual(false); // only zh-CN is valid
});

// getLocaleData
test("getLocaleData - returns data for valid locale", () => {
  expect(getLocaleData("en")).toMatchSnapshot();
  expect(getLocaleData("ar")).toMatchSnapshot(); // RTL locale
  expect(getLocaleData("pt-PT")).toMatchSnapshot(); // hyphenated locale
});

test("getLocaleData - returns undefined for invalid locale", () => {
  expect(getLocaleData("invalid")).toEqual(undefined);
  expect(getLocaleData("")).toEqual(undefined);
});

// isAllowedURI
test("isAllowedURI - null/empty/undefined are allowed", () => {
  expect(isAllowedURI(null, ["https://example.com/"])).toEqual(true);
  expect(isAllowedURI(undefined, ["https://example.com/"])).toEqual(true);
  expect(isAllowedURI("", ["https://example.com/"])).toEqual(true);
});

test("isAllowedURI - matching prefix", () => {
  const prefixes = ["https://objects.aya.is/", "https://avatars.githubusercontent.com/"];
  expect(isAllowedURI("https://objects.aya.is/image.png", prefixes)).toEqual(true);
  expect(isAllowedURI("https://avatars.githubusercontent.com/u/123", prefixes)).toEqual(true);
});

test("isAllowedURI - non-matching prefix", () => {
  const prefixes = ["https://objects.aya.is/"];
  expect(isAllowedURI("https://evil.com/hack.png", prefixes)).toEqual(false);
  expect(isAllowedURI("http://objects.aya.is/image.png", prefixes)).toEqual(false); // http vs https
});

test("isAllowedURI - empty prefix list allows all", () => {
  expect(isAllowedURI("https://anything.com/file.png", [])).toEqual(true);
});
