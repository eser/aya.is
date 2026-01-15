import { fetcher } from "@/shared/modules/backend/fetcher.ts";

export type UpdateProfileTranslationRequest = {
  title: string;
  description: string;
  // deno-lint-ignore no-explicit-any
  properties?: Record<string, any> | null;
};

export type UpdateProfileTranslationResponse = {
  success: boolean;
  message: string;
};

export async function updateProfileTranslation(
  locale: string,
  slug: string,
  translationLocale: string,
  data: UpdateProfileTranslationRequest,
): Promise<UpdateProfileTranslationResponse | null> {
  const response = await fetcher<UpdateProfileTranslationResponse>(
    `/${locale}/profiles/${slug}/translations/${translationLocale}`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    },
  );

  return response;
}
