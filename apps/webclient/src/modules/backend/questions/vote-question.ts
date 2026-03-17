// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
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
