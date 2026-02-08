import { fetcher } from "../fetcher.ts";
import type { ProfileQuestion } from "../types.ts";

export interface CreateQuestionInput {
  content: string;
  is_anonymous: boolean;
}

export async function createQuestion(
  locale: string,
  slug: string,
  input: CreateQuestionInput,
): Promise<ProfileQuestion | null> {
  const response = await fetcher<ProfileQuestion>(
    locale,
    `/profiles/${slug}/_questions`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
