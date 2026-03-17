// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher.ts";

export interface LockThreadInput {
  is_locked: boolean;
  profile_slug: string;
}

export interface LockThreadResult {
  status: string;
}

export async function lockThread(
  locale: string,
  threadId: string,
  input: LockThreadInput,
): Promise<LockThreadResult | null> {
  const response = await fetcher<LockThreadResult>(
    locale,
    `/discussions/threads/${threadId}/lock`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
