import { fetcher } from "../fetcher.ts";

export interface VoteQuestionResult {
  voted: boolean;
}

export async function voteQuestion(
  locale: string,
  slug: string,
  questionId: string,
): Promise<VoteQuestionResult | null> {
  const response = await fetcher<VoteQuestionResult>(
    locale,
    `/profiles/${slug}/_questions/${questionId}/vote`,
    {
      method: "POST",
    },
  );

  return response;
}
