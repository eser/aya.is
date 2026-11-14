// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { DeleteAccountRequest } from "./types";

export async function deletePrivacyAccount(
  locale: string,
  request: DeleteAccountRequest,
): Promise<{ deleted: boolean } | null> {
  return await fetcher<{ deleted: boolean }>(locale, "/privacy/delete", {
    method: "POST",
    body: JSON.stringify(request),
  });
}
