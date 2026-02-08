import { fetcher } from "../fetcher.ts";
import type { CursoredResponse, ProfileQuestion } from "../types.ts";

export type GetProfileQuestionsData = CursoredResponse<ProfileQuestion[]>;

export async function getProfileQuestions(
  locale: string,
  slug: string,
): Promise<GetProfileQuestionsData | null> {
  const response = await fetcher<GetProfileQuestionsData>(
    locale,
    `/profiles/${slug}/_questions`,
  );

  return response;
}
