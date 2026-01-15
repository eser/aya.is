import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { User } from "./types.ts";

export type GetUsersData = User[];

export async function getUsers(locale: string) {
  const response = await fetcher<GetUsersData>(`/${locale}/users`);

  return response;
}
