/// <reference lib="deno.ns" />

import { assertEquals } from "jsr:@std/assert@1.0.14";

import { type Locale, localeMatch } from "./locales.ts";

export const supportedLocales: readonly Locale[] = [
  {
    code: "tr",
    matches: ["@(tr)?(-*)"],
    name: "Türkçe",
    flag: "🇹🇷",
    dir: "ltr",
  },
  {
    code: "en",
    matches: ["@(en)?(-*)"],
    name: "English",
    flag: "🇺🇸",
    dir: "ltr",
  },
];

Deno.test({
  name: "example test",
  fn() {
    const locales = ["de", "de-DE", "tr-TR", "tr"];

    const locale = localeMatch(locales, supportedLocales);

    assertEquals(locale?.code, "tr");
  },
});

Deno.test({
  name: "example test",
  fn() {
    const locales = ["de", "de-DE", "tr", "tr-TR"];

    const locale = localeMatch(locales, supportedLocales);

    assertEquals(locale?.code, "tr");
  },
});
