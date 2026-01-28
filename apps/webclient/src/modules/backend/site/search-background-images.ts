import { fetcher } from "../fetcher";

export interface UnsplashPhotoUrls {
  raw: string;
  full: string;
  regular: string;
  small: string;
  thumb: string;
}

export interface UnsplashPhotoUser {
  id: string;
  username: string;
  name: string;
}

export interface UnsplashPhoto {
  id: string;
  description: string | null;
  alt_description: string | null;
  width: number;
  height: number;
  color: string;
  urls: UnsplashPhotoUrls;
  user: UnsplashPhotoUser;
}

export interface UnsplashSearchResult {
  total: number;
  total_pages: number;
  results: UnsplashPhoto[];
}

export interface SearchBackgroundImagesParams {
  query: string;
  page?: number;
  perPage?: number;
}

export async function searchBackgroundImages(
  locale: string,
  params: SearchBackgroundImagesParams,
): Promise<UnsplashSearchResult | null> {
  const searchParams = new URLSearchParams();
  searchParams.set("query", params.query);

  if (params.page !== undefined) {
    searchParams.set("page", String(params.page));
  }

  if (params.perPage !== undefined) {
    searchParams.set("per_page", String(params.perPage));
  }

  return await fetcher<UnsplashSearchResult>(
    locale,
    `/site/generated-backgrounds/?${searchParams.toString()}`,
  );
}
