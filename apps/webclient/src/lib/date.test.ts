// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
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
  normalizeSpaces,
  parseDateFromSlug,
} from "./date.ts";

// Use a fixed date for deterministic snapshots
const FIXED_DATE = new Date(2026, 0, 28, 14, 30, 0); // Jan 28, 2026, 14:30

test("formatDateLong - multiple locales", () => {
  expect(formatDateLong(FIXED_DATE, "en")).toMatchSnapshot();
  expect(formatDateLong(FIXED_DATE, "tr")).toMatchSnapshot();
  expect(formatDateLong(FIXED_DATE, "ja")).toMatchSnapshot();
});

test("formatDateShort - multiple locales", () => {
  expect(formatDateShort(FIXED_DATE, "en")).toMatchSnapshot();
  expect(formatDateShort(FIXED_DATE, "tr")).toMatchSnapshot();
  expect(formatDateShort(FIXED_DATE, "de")).toMatchSnapshot();
});

test("formatDateTimeLong", () => {
  expect(formatDateTimeLong(FIXED_DATE, "en")).toMatchSnapshot();
});

test("formatDateTimeShort", () => {
  expect(formatDateTimeShort(FIXED_DATE, "en")).toMatchSnapshot();
});

test("formatDateTimeRange - same day", () => {
  const start = new Date(2026, 1, 23, 20, 30);
  const end = new Date(2026, 1, 23, 21, 30);
  const result = formatDateTimeRange(start, end, "en");
  // Intl.DateTimeFormat.formatRange output varies by platform ICU data
  expect(typeof result).toEqual("string");
  expect(result.length > 0).toEqual(true);
  expect(result.includes("23")).toEqual(true);
});

test("formatDateTimeRange - different days", () => {
  const start = new Date(2026, 1, 23, 20, 30);
  const end = new Date(2026, 1, 24, 1, 30);
  const result = formatDateTimeRange(start, end, "en");
  expect(typeof result).toEqual("string");
  expect(result.length > 0).toEqual(true);
  expect(result.includes("23")).toEqual(true);
  expect(result.includes("24")).toEqual(true);
});

test("formatMonthYear", () => {
  expect(formatMonthYear(FIXED_DATE, "en")).toMatchSnapshot();
  expect(formatMonthYear(FIXED_DATE, "tr")).toMatchSnapshot();
});

test("formatMonthShort", () => {
  expect(formatMonthShort(FIXED_DATE, "en")).toMatchSnapshot();
});

test("formatTime", () => {
  expect(formatTime(FIXED_DATE, "en")).toMatchSnapshot();
});

test("formatTimeString - valid and null", () => {
  const result = formatTimeString("2026-01-28T14:30:00Z", "en");
  // Time format varies by platform ICU data (e.g. "2:30 PM" vs "2:30 pm")
  expect(typeof result).toEqual("string");
  expect(result.length > 0).toEqual(true);
  expect(result.includes("30")).toEqual(true);
  expect(formatTimeString(null, "en")).toEqual("");
});

test("formatDateString - valid and null", () => {
  expect(formatDateString("2026-01-28T14:30:00Z", "en")).toMatchSnapshot();
  expect(formatDateString(null, "en")).toMatchSnapshot();
});

// parseDateFromSlug — unit tests with exact assertions (not snapshot)
// because Date objects don't serialize deterministically across locales

test("parseDateFromSlug - valid slug", () => {
  const result = parseDateFromSlug("20160610-microsoftun-yari-vizyonu");
  expect(result !== null).toEqual(true);
  expect(result!.getFullYear()).toEqual(2016);
  expect(result!.getMonth()).toEqual(5); // June = 5
  expect(result!.getDate()).toEqual(10);
});

test("parseDateFromSlug - valid slug with only date", () => {
  const result = parseDateFromSlug("20260128");
  expect(result !== null).toEqual(true);
  expect(result!.getFullYear()).toEqual(2026);
  expect(result!.getMonth()).toEqual(0); // January = 0
  expect(result!.getDate()).toEqual(28);
});

test("parseDateFromSlug - null input", () => {
  expect(parseDateFromSlug(null)).toEqual(null);
});

test("parseDateFromSlug - empty string", () => {
  expect(parseDateFromSlug("")).toEqual(null);
});

test("parseDateFromSlug - no date prefix", () => {
  expect(parseDateFromSlug("hello-world")).toEqual(null);
});

test("parseDateFromSlug - invalid date (Feb 31)", () => {
  expect(parseDateFromSlug("20260231-invalid")).toEqual(null);
});

test("parseDateFromSlug - year out of range", () => {
  expect(parseDateFromSlug("18990101-too-old")).toEqual(null);
  expect(parseDateFromSlug("21010101-too-far")).toEqual(null);
});

test("normalizeSpaces - replaces thin, narrow no-break and no-break spaces with plain spaces", () => {
  expect(normalizeSpaces("10:00\u202fAM\u2009\u2013\u2009noon\u00a0today")).toBe("10:00 AM \u2013 noon today");
});

test("formatDateTimeRange - output never contains runtime-specific spaces (hydration safety)", () => {
  const start = new Date(2026, 6, 10, 10, 0, 0);
  const end = new Date(2026, 6, 12, 17, 0, 0);

  for (const locale of ["en", "tr", "de", "fr", "ja", "ar"]) {
    expect(formatDateTimeRange(start, end, locale)).not.toMatch(/[\u2009\u202f\u00a0]/);
  }
});
