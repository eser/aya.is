// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "@/modules/backend/fetcher";

export interface RemoveUploadResult {
  success: boolean;
  message: string;
}

export async function removeUpload(
  locale: string,
  key: string,
): Promise<RemoveUploadResult | null> {
  const response = await fetcher<RemoveUploadResult>(
    locale,
    `/site/uploads/${key}`,
    {
      method: "DELETE",
    },
  );
  return response;
}
