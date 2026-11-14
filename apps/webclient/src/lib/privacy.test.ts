// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { readFile } from "node:fs/promises";
import { describe, expect, test } from "vitest";
import { describeExport, isDeletionConfirmed } from "./privacy.ts";
import { readPrivacyError } from "@/modules/backend/privacy/errors.ts";
import type { PrivacyExportState } from "@/modules/backend/privacy/types.ts";

function exportState(overrides: Partial<PrivacyExportState>): PrivacyExportState {
  return {
    requested_at: null,
    ready_at: null,
    expires_at: null,
    next_request_at: null,
    id: null,
    status: "none",
    size_bytes: null,
    available: true,
    can_request: true,
    can_download: false,
    ...overrides,
  };
}

describe("isDeletionConfirmed", () => {
  test("accepts the username regardless of case and surrounding spaces", () => {
    expect(isDeletionConfirmed("jane", "jane")).toBe(true);
    expect(isDeletionConfirmed("  JANE ", "jane")).toBe(true);
  });

  test("rejects anything else, including an empty username", () => {
    expect(isDeletionConfirmed("", "jane")).toBe(false);
    expect(isDeletionConfirmed("jan", "jane")).toBe(false);
    expect(isDeletionConfirmed("jane2", "jane")).toBe(false);
    expect(isDeletionConfirmed("", "")).toBe(false);
    expect(isDeletionConfirmed("anything", "")).toBe(false);
  });
});

describe("describeExport", () => {
  test("never requested", () => {
    expect(describeExport(exportState({}))).toEqual({
      view: "idle",
      canRequest: true,
      canDownload: false,
      nextRequestAt: null,
    });
  });

  test("feature not configured on the server", () => {
    const model = describeExport(exportState({ available: false, can_request: false }));
    expect(model.view).toBe("unavailable");
    expect(model.canRequest).toBe(false);
  });

  test("pending and processing both read as preparing and block new requests", () => {
    for (const status of ["pending", "processing"] as const) {
      const model = describeExport(
        exportState({ status, can_request: false, next_request_at: "2026-12-14T00:00:00Z" }),
      );
      expect(model).toEqual({ view: "preparing", canRequest: false, canDownload: false, nextRequestAt: null });
    }
  });

  test("ready archive can be downloaded, and the cooldown is reported", () => {
    const model = describeExport(
      exportState({ status: "ready", can_download: true, can_request: false, next_request_at: "2026-12-14T00:00:00Z" }),
    );
    expect(model).toEqual({
      view: "ready",
      canRequest: false,
      canDownload: true,
      nextRequestAt: "2026-12-14T00:00:00Z",
    });
  });

  test("expired archive cannot be downloaded", () => {
    const model = describeExport(exportState({ status: "expired", can_download: false, can_request: true }));
    expect(model.view).toBe("expired");
    expect(model.canDownload).toBe(false);
    expect(model.canRequest).toBe(true);
  });

  test("a failed export can be retried right away", () => {
    const model = describeExport(exportState({ status: "failed", can_request: true }));
    expect(model).toEqual({ view: "failed", canRequest: true, canDownload: false, nextRequestAt: null });
  });
});

describe("readPrivacyError", () => {
  function fetcherError(response: Record<string, unknown>): Error {
    return new Error(String(response["error"]), { cause: { url: "x", response } });
  }

  test("reads the machine-readable code and the cooldown date", () => {
    const details = readPrivacyError(
      fetcherError({ error: "nope", code: "export_cooldown", next_request_at: "2026-12-14T20:11:17Z" }),
    );
    expect(details.code).toBe("export_cooldown");
    expect(details.nextRequestAt).toBe("2026-12-14T20:11:17Z");
    expect(details.message).toBe("nope");
  });

  test("reads the profiles that block a deletion and ignores malformed entries", () => {
    const details = readPrivacyError(
      fetcherError({
        error: "transfer",
        code: "owns_profiles",
        profiles: [{ id: "1", slug: "acme", kind: "organization", title: "Acme" }, { id: 2 }, "oops", null],
      }),
    );
    expect(details.code).toBe("owns_profiles");
    expect(details.profiles).toEqual([{ id: "1", slug: "acme", kind: "organization", title: "Acme" }]);
  });

  test("ignores a code it does not know", () => {
    expect(readPrivacyError(fetcherError({ error: "x", code: "something_new" })).code).toBeNull();
  });

  test("survives errors that are not from the fetcher", () => {
    expect(readPrivacyError(new Error("network down"))).toEqual({
      code: null,
      message: "network down",
      nextRequestAt: null,
      profiles: [],
    });
    expect(readPrivacyError("a string")).toEqual({ code: null, message: "", nextRequestAt: null, profiles: [] });
    expect(readPrivacyError(undefined).code).toBeNull();
  });
});

describe("privacy translations", () => {
  const LOCALES = ["en", "tr", "fr", "de", "es", "pt-PT", "it", "nl", "ja", "ko", "ru", "zh-CN", "ar"];

  async function load(locale: string): Promise<Record<string, string>> {
    const text = await readFile(new URL(`../messages/${locale}.json`, import.meta.url), "utf8");
    return JSON.parse(text).Privacy;
  }

  // i18next splits "namespace:key" at the first colon, so a lookup key with a colon renders only the part after it
  // (the delete-account label once showed just the username). Values may contain colons; keys must not.
  test("lookup keys never contain the i18next namespace separator", async () => {
    const keys = Object.keys(await load("en"));
    expect(keys.filter((key) => key.includes(":"))).toEqual([]);
  });

  test("every locale has exactly the English keys and keeps the interpolation placeholders", async () => {
    const en = await load("en");

    for (const locale of LOCALES) {
      const entries = await load(locale);
      expect(Object.keys(entries).sort(), locale).toEqual(Object.keys(en).sort());

      for (const [key, value] of Object.entries(entries)) {
        for (const placeholder of ["{{date}}", "{{slug}}"]) {
          expect(value.includes(placeholder), `${locale}: "${key}" placeholder ${placeholder}`).toBe(
            key.includes(placeholder),
          );
        }
      }
    }
  });

  test("no non-English locale ships the English text", async () => {
    const en = await load("en");
    // Words that legitimately stay in English everywhere, such as the placeholder profile name.
    const shared = new Set(["Privacy"]);

    for (const locale of LOCALES.filter((l) => l !== "en")) {
      const entries = await load(locale);
      const untranslated = Object.entries(entries).filter(([key, value]) => value === en[key] && !shared.has(key));
      expect(
        untranslated.map(([key]) => key),
        locale,
      ).toEqual([]);
    }
  });
});
