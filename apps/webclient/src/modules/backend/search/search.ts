import { fetcher } from "../fetcher";

export interface SearchResult {
  type: "profile" | "story" | "page";
  id: string;
  slug: string;
  title: string;
  summary: string | null;
  image_uri: string | null;
  profile_slug: string | null;
  profile_title: string | null;
  kind: string | null;
  rank: number;
}

export async function search(
  locale: string,
  query: string,
  profileSlug?: string,
  limit: number = 20,
): Promise<SearchResult[] | null> {
  const params = new URLSearchParams({
    q: query,
    limit: String(limit),
  });

  if (profileSlug !== undefined && profileSlug !== "") {
    params.set("profile", profileSlug);
  }

  return await fetcher<SearchResult[]>(`/${locale}/search?${params}`);
}
