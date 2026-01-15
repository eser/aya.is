import { fetcher } from "../fetcher";
import type { Session } from "../types";

export async function getCurrentSession(
  locale: string
): Promise<Session | null> {
  return fetcher<Session>(`/${locale}/auth/session`);
}
