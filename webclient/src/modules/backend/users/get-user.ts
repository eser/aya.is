import { fetcher } from "../fetcher";
import type { User } from "../types";

export async function getUser(locale: string, id: string): Promise<User | null> {
  return fetcher<User>(`/${locale}/users/${id}`);
}
