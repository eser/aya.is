// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "@/modules/backend/fetcher";

export interface DeleteStoryTranslationResult {
  success: boolean;
}

export async function deleteStoryTranslation(
  locale: string,
  profileSlug: string,
  storyId: string,
  translationLocale: string,
): Promise<DeleteStoryTranslationResult | null> {
  const response = await fetcher<DeleteStoryTranslationResult>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/translations/${translationLocale}`,
    {
      method: "DELETE",
    },
  );
  return response;
}
