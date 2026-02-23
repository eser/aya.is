import { fetcher } from "../fetcher.ts";
import type { DiscussionComment } from "../types.ts";
import type { CreateCommentInput } from "./create-story-comment.ts";

export async function createProfileComment(
  locale: string,
  profileSlug: string,
  input: CreateCommentInput,
): Promise<DiscussionComment | null> {
  const response = await fetcher<DiscussionComment>(
    locale,
    `/profiles/${profileSlug}/_discussions`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
