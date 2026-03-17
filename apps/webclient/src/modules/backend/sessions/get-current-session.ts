// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { Session } from "../types";

export async function getCurrentSession(
  locale: string,
): Promise<Session | null> {
  return await fetcher<Session>(locale, `/auth/session`);
}
