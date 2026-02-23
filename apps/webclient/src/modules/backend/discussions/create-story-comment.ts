import { fetcher } from "../fetcher.ts";
import type { DiscussionComment } from "../types.ts";

export interface CreateCommentInput {
  content: string;
  parent_id?: string | null;
}

export async function createStoryComment(
  locale: string,
  storySlug: string,
  input: CreateCommentInput,
): Promise<DiscussionComment | null> {
  const response = await fetcher<DiscussionComment>(
    locale,
    `/stories/${storySlug}/_discussions`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
