import { fetcher } from "@/shared/modules/backend/fetcher.ts";
import type { Session } from "./types.ts";

export async function getCurrentSession(locale: string): Promise<Session | null> {
  const response = await fetcher<Session>(`/${locale}/sessions/current`);

  return response;
}
