/**
 * Registry of domains that support mobile app deep linking
 * (Universal Links on iOS, App Links on Android).
 *
 * Adding a domain here means: on touch devices, links to that domain
 * will navigate in the same tab instead of opening a new tab,
 * allowing the OS to intercept and open the native app.
 */
const APP_LINK_DOMAINS: ReadonlySet<string> = new Set([
  "youtube.com",
  "www.youtube.com",
  "m.youtube.com",
  "youtu.be",
  "music.youtube.com",
]);

/**
 * Returns true if the URL belongs to a domain that has a native mobile app
 * registered via Universal Links (iOS) or App Links (Android).
 */
export function isAppLinkableUrl(url: string): boolean {
  if (url === "" || !URL.canParse(url)) {
    return false;
  }

  const parsed = new URL(url);

  return APP_LINK_DOMAINS.has(parsed.hostname);
}
