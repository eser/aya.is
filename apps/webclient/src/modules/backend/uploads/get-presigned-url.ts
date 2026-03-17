// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "@/modules/backend/fetcher";
import type { GetPresignedURLRequest, GetPresignedURLResponse } from "@/modules/backend/types";

export async function getPresignedURL(
  locale: string,
  request: GetPresignedURLRequest,
): Promise<GetPresignedURLResponse | null> {
  const response = await fetcher<GetPresignedURLResponse>(
    locale,
    `/site/uploads/presign`,
    {
      method: "POST",
      body: JSON.stringify(request),
    },
  );
  return response;
}
