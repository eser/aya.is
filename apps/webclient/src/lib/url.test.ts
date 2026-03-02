/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals } from "@std/assert";
import { getAlternateUrls, localizedUrl, parseLocaleFromPath, sanitizeImageSrc } from "./url.ts";

// sanitizeImageSrc — security-critical
Deno.test("sanitizeImageSrc - safe protocols", async (t) => {
  await assertSnapshot(t, sanitizeImageSrc("https://example.com/image.png"));
  await assertSnapshot(t, sanitizeImageSrc("http://example.com/image.png"));
  await assertSnapshot(t, sanitizeImageSrc("data:image/png;base64,abc123"));
  await assertSnapshot(t, sanitizeImageSrc("blob:https://example.com/abc"));
});

Deno.test("sanitizeImageSrc - blocks dangerous protocols", () => {
  assertEquals(sanitizeImageSrc("javascript:alert(1)"), "");
  assertEquals(sanitizeImageSrc("vbscript:msgbox"), "");
  assertEquals(sanitizeImageSrc("file:///etc/passwd"), "");
});

Deno.test("sanitizeImageSrc - null/empty/undefined", () => {
  assertEquals(sanitizeImageSrc(null), "");
  assertEquals(sanitizeImageSrc(undefined), "");
  assertEquals(sanitizeImageSrc(""), "");
});

Deno.test("sanitizeImageSrc - relative paths", () => {
  assertEquals(sanitizeImageSrc("/images/logo.png"), "/images/logo.png");
  // Protocol-relative URLs resolve to https: via the base URL, so they pass the check
  assertEquals(sanitizeImageSrc("//evil.com/image.png"), "//evil.com/image.png");
});

// localizedUrl
Deno.test("localizedUrl - basic paths", async (t) => {
  await assertSnapshot(t, localizedUrl("/about", { locale: "en" }));
  await assertSnapshot(t, localizedUrl("/about", { locale: "tr" }));
  await assertSnapshot(t, localizedUrl("/", { locale: "fr" }));
  await assertSnapshot(t, localizedUrl("about", { locale: "de" })); // without leading slash
});

Deno.test("localizedUrl - defaults to en", async (t) => {
  await assertSnapshot(t, localizedUrl("/about"));
  await assertSnapshot(t, localizedUrl("/about", {}));
});

Deno.test("localizedUrl - uses currentLocale as fallback", async (t) => {
  await assertSnapshot(t, localizedUrl("/about", { currentLocale: "ja" }));
});

// parseLocaleFromPath
Deno.test("parseLocaleFromPath - valid locales", async (t) => {
  await assertSnapshot(t, parseLocaleFromPath("/en/about"));
  await assertSnapshot(t, parseLocaleFromPath("/tr/stories/my-story"));
  await assertSnapshot(t, parseLocaleFromPath("/pt-PT/page"));
  await assertSnapshot(t, parseLocaleFromPath("/zh-CN/"));
});

Deno.test("parseLocaleFromPath - invalid locale in path", async (t) => {
  await assertSnapshot(t, parseLocaleFromPath("/xx/about"));
  await assertSnapshot(t, parseLocaleFromPath("/about"));
  await assertSnapshot(t, parseLocaleFromPath("/"));
});

// getAlternateUrls
Deno.test("getAlternateUrls - generates for all locales", async (t) => {
  await assertSnapshot(t, getAlternateUrls("/about"));
});
