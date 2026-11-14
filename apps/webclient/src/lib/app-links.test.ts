// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { expect, test } from "vitest";
import { isAppLinkableUrl } from "./app-links.ts";

// YouTube domains
test("isAppLinkableUrl - youtube.com variants", () => {
  expect(isAppLinkableUrl("https://youtube.com/watch?v=abc123")).toEqual(true);
  expect(isAppLinkableUrl("https://www.youtube.com/watch?v=abc123")).toEqual(true);
  expect(isAppLinkableUrl("https://m.youtube.com/watch?v=abc123")).toEqual(true);
  expect(isAppLinkableUrl("https://youtu.be/abc123")).toEqual(true);
  expect(isAppLinkableUrl("https://music.youtube.com/watch?v=abc123")).toEqual(true);
});

test("isAppLinkableUrl - YouTube channels and playlists", () => {
  expect(isAppLinkableUrl("https://www.youtube.com/@channel")).toEqual(true);
  expect(isAppLinkableUrl("https://www.youtube.com/playlist?list=PLabc")).toEqual(true);
  expect(isAppLinkableUrl("https://www.youtube.com/shorts/abc123")).toEqual(true);
  expect(isAppLinkableUrl("https://www.youtube.com/live/abc123")).toEqual(true);
});

test("isAppLinkableUrl - YouTube with query params", () => {
  expect(isAppLinkableUrl("https://www.youtube.com/watch?v=abc123&si=xyz&t=42")).toEqual(true);
  expect(isAppLinkableUrl("https://youtu.be/abc123?feature=share")).toEqual(true);
});

// Non-YouTube domains
test("isAppLinkableUrl - non-app-linkable domains", () => {
  expect(isAppLinkableUrl("https://github.com/user/repo")).toEqual(false);
  expect(isAppLinkableUrl("https://twitter.com/user")).toEqual(false);
  expect(isAppLinkableUrl("https://example.com")).toEqual(false);
  expect(isAppLinkableUrl("https://notyoutube.com/watch")).toEqual(false);
});

// Edge cases
test("isAppLinkableUrl - empty and invalid URLs", () => {
  expect(isAppLinkableUrl("")).toEqual(false);
  expect(isAppLinkableUrl("not-a-url")).toEqual(false);
  expect(isAppLinkableUrl("javascript:alert(1)")).toEqual(false);
});

test("isAppLinkableUrl - protocol variations", () => {
  expect(isAppLinkableUrl("http://www.youtube.com/watch?v=abc123")).toEqual(true);
  expect(isAppLinkableUrl("https://www.youtube.com/watch?v=abc123")).toEqual(true);
});
