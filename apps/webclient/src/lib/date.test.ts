// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { assertSnapshot } from "@std/testing/snapshot";
import { assertEquals } from "@std/assert";
import {
  formatDateLong,
  formatDateShort,
  formatDateString,
  formatDateTimeLong,
  formatDateTimeRange,
  formatDateTimeShort,
  formatMonthShort,
  formatMonthYear,
  formatTime,
  formatTimeString,
  parseDateFromSlug,
} from "./date.ts";

// Use a fixed date for deterministic snapshots
const FIXED_DATE = new Date(2026, 0, 28, 14, 30, 0); // Jan 28, 2026, 14:30

Deno.test("formatDateLong - multiple locales", async (t) => {
  await assertSnapshot(t, formatDateLong(FIXED_DATE, "en"));
  await assertSnapshot(t, formatDateLong(FIXED_DATE, "tr"));
  await assertSnapshot(t, formatDateLong(FIXED_DATE, "ja"));
});

Deno.test("formatDateShort - multiple locales", async (t) => {
  await assertSnapshot(t, formatDateShort(FIXED_DATE, "en"));
  await assertSnapshot(t, formatDateShort(FIXED_DATE, "tr"));
  await assertSnapshot(t, formatDateShort(FIXED_DATE, "de"));
});

Deno.test("formatDateTimeLong", async (t) => {
  await assertSnapshot(t, formatDateTimeLong(FIXED_DATE, "en"));
});

Deno.test("formatDateTimeShort", async (t) => {
  await assertSnapshot(t, formatDateTimeShort(FIXED_DATE, "en"));
});

Deno.test("formatDateTimeRange - same day", async (t) => {
  const start = new Date(2026, 1, 23, 20, 30);
  const end = new Date(2026, 1, 23, 21, 30);
  await assertSnapshot(t, formatDateTimeRange(start, end, "en"));
});

Deno.test("formatDateTimeRange - different days", async (t) => {
  const start = new Date(2026, 1, 23, 20, 30);
  const end = new Date(2026, 1, 24, 1, 30);
  await assertSnapshot(t, formatDateTimeRange(start, end, "en"));
});

Deno.test("formatMonthYear", async (t) => {
  await assertSnapshot(t, formatMonthYear(FIXED_DATE, "en"));
  await assertSnapshot(t, formatMonthYear(FIXED_DATE, "tr"));
});

Deno.test("formatMonthShort", async (t) => {
  await assertSnapshot(t, formatMonthShort(FIXED_DATE, "en"));
});

Deno.test("formatTime", async (t) => {
  await assertSnapshot(t, formatTime(FIXED_DATE, "en"));
});

Deno.test("formatTimeString - valid and null", async (t) => {
  await assertSnapshot(t, formatTimeString("2026-01-28T14:30:00Z", "en"));
  await assertSnapshot(t, formatTimeString(null, "en"));
});

Deno.test("formatDateString - valid and null", async (t) => {
  await assertSnapshot(t, formatDateString("2026-01-28T14:30:00Z", "en"));
  await assertSnapshot(t, formatDateString(null, "en"));
});

// parseDateFromSlug — unit tests with exact assertions (not snapshot)
// because Date objects don't serialize deterministically across locales

Deno.test("parseDateFromSlug - valid slug", () => {
  const result = parseDateFromSlug("20160610-microsoftun-yari-vizyonu");
  assertEquals(result !== null, true);
  assertEquals(result!.getFullYear(), 2016);
  assertEquals(result!.getMonth(), 5); // June = 5
  assertEquals(result!.getDate(), 10);
});

Deno.test("parseDateFromSlug - valid slug with only date", () => {
  const result = parseDateFromSlug("20260128");
  assertEquals(result !== null, true);
  assertEquals(result!.getFullYear(), 2026);
  assertEquals(result!.getMonth(), 0); // January = 0
  assertEquals(result!.getDate(), 28);
});

Deno.test("parseDateFromSlug - null input", () => {
  assertEquals(parseDateFromSlug(null), null);
});

Deno.test("parseDateFromSlug - empty string", () => {
  assertEquals(parseDateFromSlug(""), null);
});

Deno.test("parseDateFromSlug - no date prefix", () => {
  assertEquals(parseDateFromSlug("hello-world"), null);
});

Deno.test("parseDateFromSlug - invalid date (Feb 31)", () => {
  assertEquals(parseDateFromSlug("20260231-invalid"), null);
});

Deno.test("parseDateFromSlug - year out of range", () => {
  assertEquals(parseDateFromSlug("18990101-too-old"), null);
  assertEquals(parseDateFromSlug("21010101-too-far"), null);
});
