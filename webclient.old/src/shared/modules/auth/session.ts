import "server-only";

import { getCookie, setCookie, unsetCookie } from "@/shared/lib/cookies.ts";
import { decrypt, encrypt } from "./jtw-auth.ts";

const SESSION_COOKIE_NAME = "session";
const SESSION_COOKIE_EXPIRATION_TIME = 7 * 24 * 60 * 60 * 1000;

export async function createSession(userId: string) {
  const expiresAt = new Date(Date.now() + SESSION_COOKIE_EXPIRATION_TIME);
  const session = await encrypt({ userId, expiresAt });

  await setCookie(SESSION_COOKIE_NAME, session, expiresAt);
}

export function getSession() {
  return getCookie(SESSION_COOKIE_NAME);
}

export async function updateSession() {
  const session = await getCookie(SESSION_COOKIE_NAME);

  if (session === null) {
    return null;
  }

  const payload = await decrypt(session);

  if (payload === undefined) {
    return null;
  }

  const expiresAt = new Date(Date.now() + SESSION_COOKIE_EXPIRATION_TIME);

  await setCookie(SESSION_COOKIE_NAME, session, expiresAt);
}

export function deleteSession() {
  return unsetCookie(SESSION_COOKIE_NAME);
}
