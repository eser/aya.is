import { fetcher } from "../fetcher";

export type CheckStorySlugResponse = {
  available: boolean;
  message?: string;
  severity?: "error" | "warning" | "";
};

export type CheckStorySlugOptions = {
  excludeId?: string;
  storyId?: string;
  publishedAt?: string | null;
};

export async function checkStorySlug(
  locale: string,
  slug: string,
  options?: CheckStorySlugOptions,
): Promise<CheckStorySlugResponse | null> {
  const params = new URLSearchParams();

  if (options?.excludeId !== undefined) {
    params.set("exclude_id", options.excludeId);
  }

  if (options?.storyId !== undefined) {
    params.set("story_id", options.storyId);
  }

  if (options?.publishedAt !== undefined && options.publishedAt !== null) {
    params.set("published_at", options.publishedAt);
  }

  const queryString = params.toString();
  const queryParams = queryString.length > 0 ? `?${queryString}` : "";

  const result = await fetcher<CheckStorySlugResponse>(
    locale,
    `/stories/${slug}/_check${queryParams}`,
    {
      method: "GET",
    },
  );

  return result;
}
