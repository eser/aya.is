import { cookies } from "next/headers";

export async function getCookie(name: string) {
  const cookieStore = await cookies();
  const cookie = cookieStore.get(name);

  if (cookie === undefined) {
    return null;
  }

  return cookie.value;
}

export async function setCookie(name: string, value: string, expiresAt: Date) {
  const cookieStore = await cookies();

  cookieStore.set(name, value, {
    httpOnly: true,
    secure: true,
    expires: expiresAt,
    sameSite: "lax",
    path: "/",
  });
}

export async function unsetCookie(name: string) {
  const cookieStore = await cookies();

  cookieStore.delete(name);
}
