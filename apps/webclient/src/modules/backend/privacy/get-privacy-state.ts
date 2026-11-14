// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { PrivacyState } from "./types";

export async function getPrivacyState(locale: string): Promise<PrivacyState | null> {
  return await fetcher<PrivacyState>(locale, "/privacy");
}
