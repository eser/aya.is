// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";

export async function hidePrivacyProfile(locale: string): Promise<{ hidden: boolean } | null> {
  return await fetcher<{ hidden: boolean }>(locale, "/privacy/hide", { method: "POST" });
}
