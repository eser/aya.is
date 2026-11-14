// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";

export async function getPrivacyExportDownloadUrl(locale: string, exportId: string): Promise<{ url: string } | null> {
  return await fetcher<{ url: string }>(locale, `/privacy/exports/${encodeURIComponent(exportId)}/download`);
}
