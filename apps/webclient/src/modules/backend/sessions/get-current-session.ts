import { fetcher } from "../fetcher";
import type { Session } from "../types";

export async function getCurrentSession(
  locale: string,
): Promise<Session | null> {
  return await fetcher<Session>(`/${locale}/auth/session`);
}
