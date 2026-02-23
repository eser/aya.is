import { fetcher } from "../fetcher.ts";

export interface EditCommentInput {
  content: string;
}

export interface EditCommentResult {
  status: string;
}

export async function editComment(
  locale: string,
  commentId: string,
  input: EditCommentInput,
): Promise<EditCommentResult | null> {
  const response = await fetcher<EditCommentResult>(
    locale,
    `/discussions/comments/${commentId}`,
    {
      method: "PUT",
      body: JSON.stringify(input),
    },
  );

  return response;
}
