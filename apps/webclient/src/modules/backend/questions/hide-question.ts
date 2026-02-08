import { fetcher } from "../fetcher.ts";

export interface HideQuestionInput {
  is_hidden: boolean;
}

export interface HideQuestionResult {
  status: string;
}

export async function hideQuestion(
  locale: string,
  slug: string,
  questionId: string,
  input: HideQuestionInput,
): Promise<HideQuestionResult | null> {
  const response = await fetcher<HideQuestionResult>(
    locale,
    `/profiles/${slug}/_questions/${questionId}/hide`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
