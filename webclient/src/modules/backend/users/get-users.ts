import { fetcher } from "../fetcher";
import type { User } from "../types";

export async function getUsers(locale: string): Promise<User[] | null> {
  return fetcher<User[]>(`/${locale}/users`);
}
