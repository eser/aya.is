import { fetcher } from "@/modules/backend/fetcher";

export interface RemoveStoryPublicationResult {
  success: boolean;
  message: string;
}

export async function removeStoryPublication(
  locale: string,
  profileSlug: string,
  storyId: string,
  publicationId: string,
): Promise<RemoveStoryPublicationResult | null> {
  const response = await fetcher<RemoveStoryPublicationResult>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/publications/${publicationId}`,
    {
      method: "DELETE",
    },
  );

  return response;
}
