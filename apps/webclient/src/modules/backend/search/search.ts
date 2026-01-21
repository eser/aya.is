import { fetcher } from "../fetcher";

export interface SearchResult {
  type: "profile" | "story" | "page";
  id: string;
  slug: string;
  title: string;
  summary: string | null;
  image_uri: string | null;
  profile_slug: string | null;
  rank: number;
}

export async function search(
  locale: string,
  query: string,
  limit: number = 20,
): Promise<SearchResult[] | null> {
  const params = new URLSearchParams({
    q: query,
    limit: String(limit),
  });

  return await fetcher<SearchResult[]>(`/${locale}/search?${params}`);
}
