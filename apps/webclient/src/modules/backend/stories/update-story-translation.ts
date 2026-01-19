import { fetcher } from "@/modules/backend/fetcher";
import type { UpdateStoryTranslationInput } from "@/modules/backend/types";

export interface UpdateStoryTranslationResult {
  success: boolean;
  message: string;
}

export async function updateStoryTranslation(
  locale: string,
  profileSlug: string,
  storyId: string,
  translationLocale: string,
  input: UpdateStoryTranslationInput,
): Promise<UpdateStoryTranslationResult | null> {
  const response = await fetcher<UpdateStoryTranslationResult>(
    `/${locale}/profiles/${profileSlug}/_stories/${storyId}/translations/${translationLocale}`,
    {
      method: "PATCH",
      body: JSON.stringify(input),
    },
  );
  return response;
}
