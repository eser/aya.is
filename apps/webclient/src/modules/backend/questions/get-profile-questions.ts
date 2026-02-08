import { fetcher } from "../fetcher.ts";
import type { ProfileQuestion } from "../types.ts";

export async function getProfileQuestions(
  locale: string,
  slug: string,
): Promise<ProfileQuestion[] | null> {
  const response = await fetcher<ProfileQuestion[]>(
    locale,
    `/profiles/${slug}/_questions`,
  );

  return response;
}
