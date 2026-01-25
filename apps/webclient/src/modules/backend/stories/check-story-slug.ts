import { fetcher } from "../fetcher";

export type CheckStorySlugResponse = {
  available: boolean;
  message?: string;
};

export async function checkStorySlug(
  locale: string,
  slug: string,
  excludeId?: string,
): Promise<CheckStorySlugResponse | null> {
  const queryParams = excludeId !== undefined ? `?exclude_id=${excludeId}` : "";
  const result = await fetcher<CheckStorySlugResponse>(
    locale,
    `/stories/${slug}/_check${queryParams}`,
    {
      method: "GET",
    },
  );

  return result;
}
