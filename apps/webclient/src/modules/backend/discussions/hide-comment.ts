// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher.ts";

export interface HideCommentInput {
  is_hidden: boolean;
  profile_slug: string;
}

export interface HideCommentResult {
  status: string;
}

export async function hideComment(
  locale: string,
  commentId: string,
  input: HideCommentInput,
): Promise<HideCommentResult | null> {
  const response = await fetcher<HideCommentResult>(
    locale,
    `/discussions/comments/${commentId}/hide`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
