import { fetcher } from "../fetcher.ts";

export interface DeleteCommentResult {
  status: string;
}

export async function deleteComment(
  locale: string,
  commentId: string,
  profileSlug?: string,
): Promise<DeleteCommentResult | null> {
  const query = profileSlug !== undefined ? `?profile_slug=${encodeURIComponent(profileSlug)}` : "";

  const response = await fetcher<DeleteCommentResult>(
    locale,
    `/discussions/comments/${commentId}${query}`,
    {
      method: "DELETE",
    },
  );

  return response;
}
