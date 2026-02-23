import { fetcher } from "../fetcher.ts";

export interface PinCommentInput {
  is_pinned: boolean;
  profile_slug: string;
}

export interface PinCommentResult {
  status: string;
}

export async function pinComment(
  locale: string,
  commentId: string,
  input: PinCommentInput,
): Promise<PinCommentResult | null> {
  const response = await fetcher<PinCommentResult>(
    locale,
    `/discussions/comments/${commentId}/pin`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );

  return response;
}
