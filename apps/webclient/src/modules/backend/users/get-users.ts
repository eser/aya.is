import { fetcher } from "../fetcher";
import type { User } from "../types";

export async function getUsers(locale: string): Promise<User[] | null> {
  return await fetcher<User[]>(`/${locale}/users`);
}
