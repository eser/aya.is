import { fetcher } from "../fetcher";

export type CheckProfileSlugResponse = {
  available: boolean;
  message?: string;
};

export async function checkProfileSlug(
  locale: string,
  slug: string,
): Promise<CheckProfileSlugResponse | null> {
  return await fetcher<CheckProfileSlugResponse>(
    locale,
    `/profiles/${slug}/_check`,
    {
      method: "GET",
    },
  );
}
