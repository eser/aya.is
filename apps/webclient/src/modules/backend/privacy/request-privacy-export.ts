// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { RequestedPrivacyExport } from "./types";

export async function requestPrivacyExport(locale: string): Promise<RequestedPrivacyExport | null> {
  return await fetcher<RequestedPrivacyExport>(locale, "/privacy/exports", { method: "POST" });
}
