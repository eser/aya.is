/// <reference lib="deno.ns" />
import { assertEquals } from "@std/assert";
import { isAppLinkableUrl } from "./app-links.ts";

// YouTube domains
Deno.test("isAppLinkableUrl - youtube.com variants", () => {
  assertEquals(isAppLinkableUrl("https://youtube.com/watch?v=abc123"), true);
  assertEquals(isAppLinkableUrl("https://www.youtube.com/watch?v=abc123"), true);
  assertEquals(isAppLinkableUrl("https://m.youtube.com/watch?v=abc123"), true);
  assertEquals(isAppLinkableUrl("https://youtu.be/abc123"), true);
  assertEquals(isAppLinkableUrl("https://music.youtube.com/watch?v=abc123"), true);
});

Deno.test("isAppLinkableUrl - YouTube channels and playlists", () => {
  assertEquals(isAppLinkableUrl("https://www.youtube.com/@channel"), true);
  assertEquals(isAppLinkableUrl("https://www.youtube.com/playlist?list=PLabc"), true);
  assertEquals(isAppLinkableUrl("https://www.youtube.com/shorts/abc123"), true);
  assertEquals(isAppLinkableUrl("https://www.youtube.com/live/abc123"), true);
});

Deno.test("isAppLinkableUrl - YouTube with query params", () => {
  assertEquals(isAppLinkableUrl("https://www.youtube.com/watch?v=abc123&si=xyz&t=42"), true);
  assertEquals(isAppLinkableUrl("https://youtu.be/abc123?feature=share"), true);
});

// Non-YouTube domains
Deno.test("isAppLinkableUrl - non-app-linkable domains", () => {
  assertEquals(isAppLinkableUrl("https://github.com/user/repo"), false);
  assertEquals(isAppLinkableUrl("https://twitter.com/user"), false);
  assertEquals(isAppLinkableUrl("https://example.com"), false);
  assertEquals(isAppLinkableUrl("https://notyoutube.com/watch"), false);
});

// Edge cases
Deno.test("isAppLinkableUrl - empty and invalid URLs", () => {
  assertEquals(isAppLinkableUrl(""), false);
  assertEquals(isAppLinkableUrl("not-a-url"), false);
  assertEquals(isAppLinkableUrl("javascript:alert(1)"), false);
});

Deno.test("isAppLinkableUrl - protocol variations", () => {
  assertEquals(isAppLinkableUrl("http://www.youtube.com/watch?v=abc123"), true);
  assertEquals(isAppLinkableUrl("https://www.youtube.com/watch?v=abc123"), true);
});
