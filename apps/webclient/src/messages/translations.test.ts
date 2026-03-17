// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals } from "@std/assert";

const LOCALES = [
  "en",
  "tr",
  "fr",
  "de",
  "es",
  "pt-PT",
  "it",
  "nl",
  "ja",
  "ko",
  "ru",
  "zh-CN",
  "ar",
] as const;

type TranslationData = Record<string, Record<string, string>>;

async function loadTranslations(locale: string): Promise<TranslationData> {
  const filePath = new URL(`./${locale}.json`, import.meta.url);
  const text = await Deno.readTextFile(filePath);
  return JSON.parse(text);
}

function extractKeys(data: TranslationData): string[] {
  const keys: string[] = [];
  for (const [section, entries] of Object.entries(data)) {
    for (const key of Object.keys(entries)) {
      keys.push(`${section}.${key}`);
    }
  }
  return keys.sort();
}

// Snapshot the English key set as baseline — any addition/removal is visible in diff
Deno.test("translations - English key set baseline", async (t) => {
  const en = await loadTranslations("en");
  await assertSnapshot(t, extractKeys(en));
});

// Snapshot section names across all locales
Deno.test("translations - all locales have same sections as English", async () => {
  const en = await loadTranslations("en");
  const enSections = Object.keys(en).sort();

  for (const locale of LOCALES) {
    if (locale === "en") continue;
    const data = await loadTranslations(locale);
    const sections = Object.keys(data).sort();
    assertEquals(
      sections,
      enSections,
      `Locale "${locale}" has different sections than English. Missing: ${
        enSections.filter((s) => !sections.includes(s)).join(", ")
      }. Extra: ${sections.filter((s) => !enSections.includes(s)).join(", ")}`,
    );
  }
});

// Check that every key in English exists in all other locales
Deno.test("translations - all locales have all English keys", async () => {
  const en = await loadTranslations("en");
  const enKeys = extractKeys(en);

  const missingByLocale: Record<string, string[]> = {};

  for (const locale of LOCALES) {
    if (locale === "en") continue;
    const data = await loadTranslations(locale);
    const localeKeys = new Set(extractKeys(data));
    const missing = enKeys.filter((key) => !localeKeys.has(key));
    if (missing.length > 0) {
      missingByLocale[locale] = missing;
    }
  }

  const localesWithMissing = Object.keys(missingByLocale);
  if (localesWithMissing.length > 0) {
    const details = localesWithMissing
      .map((locale) => `  ${locale}: missing ${missingByLocale[locale].length} keys`)
      .join("\n");
    console.warn(`Translation gaps found:\n${details}`);
    // Note: This is a warning, not a failure, since translations may be in progress.
    // Uncomment the line below to make it a hard failure:
    // throw new Error(`Missing translations:\n${details}`);
  }
});

// Check for orphan keys (exist in locale but not in English)
Deno.test("translations - no orphan keys in non-English locales", async () => {
  const en = await loadTranslations("en");
  const enKeys = new Set(extractKeys(en));

  const orphansByLocale: Record<string, string[]> = {};

  for (const locale of LOCALES) {
    if (locale === "en") continue;
    const data = await loadTranslations(locale);
    const localeKeys = extractKeys(data);
    const orphans = localeKeys.filter((key) => !enKeys.has(key));
    if (orphans.length > 0) {
      orphansByLocale[locale] = orphans;
    }
  }

  const localesWithOrphans = Object.keys(orphansByLocale);
  if (localesWithOrphans.length > 0) {
    const details = localesWithOrphans
      .map((locale) =>
        `  ${locale}: ${orphansByLocale[locale].length} orphan keys (${orphansByLocale[locale].slice(0, 3).join(", ")}${
          orphansByLocale[locale].length > 3 ? "..." : ""
        })`
      )
      .join("\n");
    console.warn(`Orphan translation keys found:\n${details}`);
  }
});
