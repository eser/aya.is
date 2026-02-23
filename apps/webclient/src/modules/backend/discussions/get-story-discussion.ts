import { fetcher } from "../fetcher.ts";
import type { DiscussionListResponse, DiscussionSortMode } from "../types.ts";

export async function getStoryDiscussion(
  locale: string,
  storySlug: string,
  sort?: DiscussionSortMode,
  limit?: number,
  offset?: number,
): Promise<DiscussionListResponse | null> {
  const params = new URLSearchParams();

  if (sort !== undefined) {
    params.set("sort", sort);
  }

  if (limit !== undefined) {
    params.set("limit", String(limit));
  }

  if (offset !== undefined) {
    params.set("offset", String(offset));
  }

  const query = params.toString();
  const path = `/stories/${storySlug}/_discussions${query !== "" ? `?${query}` : ""}`;

  const response = await fetcher<DiscussionListResponse>(locale, path);

  return response;
}
