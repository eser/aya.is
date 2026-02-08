import { fetcher } from "../fetcher.ts";

export interface AnswerQuestionInput {
  answer_content: string;
}

export interface AnswerQuestionResult {
  status: string;
}

export async function answerQuestion(
  locale: string,
  slug: string,
  questionId: string,
  input: AnswerQuestionInput,
): Promise<AnswerQuestionResult | null> {
  const response = await fetcher<AnswerQuestionResult>(
    locale,
    `/profiles/${slug}/_questions/${questionId}/answer`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
