import { fetcher } from "../fetcher.ts";

export interface EditAnswerInput {
  answer_content: string;
}

export interface EditAnswerResult {
  status: string;
}

export async function editAnswer(
  locale: string,
  slug: string,
  questionId: string,
  input: EditAnswerInput,
): Promise<EditAnswerResult | null> {
  const response = await fetcher<EditAnswerResult>(
    locale,
    `/profiles/${slug}/_questions/${questionId}/answer`,
    {
      method: "PUT",
      body: JSON.stringify(input),
    },
  );

  return response;
}
