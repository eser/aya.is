// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { getAlternateUrls, localizedUrl, parseLocaleFromPath, sanitizeImageSrc } from "./url.ts";

// sanitizeImageSrc — security-critical
test("sanitizeImageSrc - safe protocols", () => {
  expect(sanitizeImageSrc("https://example.com/image.png")).toMatchSnapshot();
  expect(sanitizeImageSrc("http://example.com/image.png")).toMatchSnapshot();
  expect(sanitizeImageSrc("data:image/png;base64,abc123")).toMatchSnapshot();
  expect(sanitizeImageSrc("blob:https://example.com/abc")).toMatchSnapshot();
});

test("sanitizeImageSrc - blocks dangerous protocols", () => {
  expect(sanitizeImageSrc("javascript:alert(1)")).toEqual("");
  expect(sanitizeImageSrc("vbscript:msgbox")).toEqual("");
  expect(sanitizeImageSrc("file:///etc/passwd")).toEqual("");
});

test("sanitizeImageSrc - null/empty/undefined", () => {
  expect(sanitizeImageSrc(null)).toEqual("");
  expect(sanitizeImageSrc(undefined)).toEqual("");
  expect(sanitizeImageSrc("")).toEqual("");
});

test("sanitizeImageSrc - relative paths", () => {
  expect(sanitizeImageSrc("/images/logo.png")).toEqual("/images/logo.png");
  // Protocol-relative URLs resolve to https: via the base URL, so they pass the check
  expect(sanitizeImageSrc("//evil.com/image.png")).toEqual("//evil.com/image.png");
});

// localizedUrl
test("localizedUrl - basic paths", () => {
  expect(localizedUrl("/about", { locale: "en" })).toMatchSnapshot();
  expect(localizedUrl("/about", { locale: "tr" })).toMatchSnapshot();
  expect(localizedUrl("/", { locale: "fr" })).toMatchSnapshot();
  expect(localizedUrl("about", { locale: "de" })).toMatchSnapshot(); // without leading slash
});

test("localizedUrl - defaults to en", () => {
  expect(localizedUrl("/about")).toMatchSnapshot();
  expect(localizedUrl("/about", {})).toMatchSnapshot();
});

test("localizedUrl - uses currentLocale as fallback", () => {
  expect(localizedUrl("/about", { currentLocale: "ja" })).toMatchSnapshot();
});

// parseLocaleFromPath
test("parseLocaleFromPath - valid locales", () => {
  expect(parseLocaleFromPath("/en/about")).toMatchSnapshot();
  expect(parseLocaleFromPath("/tr/stories/my-story")).toMatchSnapshot();
  expect(parseLocaleFromPath("/pt-PT/page")).toMatchSnapshot();
  expect(parseLocaleFromPath("/zh-CN/")).toMatchSnapshot();
});

test("parseLocaleFromPath - invalid locale in path", () => {
  expect(parseLocaleFromPath("/xx/about")).toMatchSnapshot();
  expect(parseLocaleFromPath("/about")).toMatchSnapshot();
  expect(parseLocaleFromPath("/")).toMatchSnapshot();
});

// getAlternateUrls
test("getAlternateUrls - generates for all locales", () => {
  expect(getAlternateUrls("/about")).toMatchSnapshot();
});
