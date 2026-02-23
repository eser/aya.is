import { fetcher } from "../fetcher.ts";
import type { DiscussionVoteResponse } from "../types.ts";

export interface VoteCommentInput {
  direction: 1 | -1;
}

export async function voteComment(
  locale: string,
  commentId: string,
  input: VoteCommentInput,
): Promise<DiscussionVoteResponse | null> {
  const response = await fetcher<DiscussionVoteResponse>(
    locale,
    `/discussions/comments/${commentId}/vote`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
