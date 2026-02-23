import { fetcher } from "../fetcher.ts";
import type { DiscussionComment } from "../types.ts";

export interface CommentRepliesResponse {
  comments: DiscussionComment[];
}

export async function getCommentReplies(
  locale: string,
  commentId: string,
  limit?: number,
  offset?: number,
): Promise<CommentRepliesResponse | null> {
  const params = new URLSearchParams();

  if (limit !== undefined) {
    params.set("limit", String(limit));
  }

  if (offset !== undefined) {
    params.set("offset", String(offset));
  }

  const query = params.toString();
  const path = `/discussions/comments/${commentId}/replies${query !== "" ? `?${query}` : ""}`;

  const response = await fetcher<CommentRepliesResponse>(locale, path);

  return response;
}
