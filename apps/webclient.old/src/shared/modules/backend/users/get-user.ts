import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { User } from "./types.ts";

export type GetUserData = User;

export async function getUser(locale: string, id: string) {
  const response = await fetcher<GetUserData>(`/${locale}/users/${id}`);

  return response;
}
