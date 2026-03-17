// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { User } from "../types";

export async function getUser(
  locale: string,
  id: string,
): Promise<User | null> {
  return await fetcher<User>(locale, `/users/${id}`);
}
