import { fetcher } from "@/modules/backend/fetcher";

export interface UpdateStoryPublicationInput {
  is_featured: boolean;
}

export interface UpdateStoryPublicationResult {
  success: boolean;
  message: string;
}

export async function updateStoryPublication(
  locale: string,
  profileSlug: string,
  storyId: string,
  publicationId: string,
  input: UpdateStoryPublicationInput,
): Promise<UpdateStoryPublicationResult | null> {
  const response = await fetcher<UpdateStoryPublicationResult>(
    locale,
    `/profiles/${profileSlug}/_stories/${storyId}/publications/${publicationId}`,
    {
      method: "PATCH",
      body: JSON.stringify(input),
    },
  );

  return response;
}
