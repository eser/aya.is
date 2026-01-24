import { fetcher } from "../fetcher";

export type UpdateProfileTranslationRequest = {
  title: string;
  description: string;
  properties?: Record<string, unknown> | null;
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
  return await fetcher<UpdateProfileTranslationResponse>(
    locale,
    `/profiles/${slug}/translations/${translationLocale}`,
    {
      method: "PATCH",
      body: JSON.stringify(data),
    },
  );
}
