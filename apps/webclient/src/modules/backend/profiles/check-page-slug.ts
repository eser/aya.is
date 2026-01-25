import { fetcher } from "../fetcher";

export type CheckPageSlugResponse = {
  available: boolean;
  message?: string;
  severity?: "error" | "warning" | "";
};

export type CheckPageSlugOptions = {
  excludeId?: string;
};

export async function checkPageSlug(
  locale: string,
  profileSlug: string,
  pageSlug: string,
  options?: CheckPageSlugOptions,
): Promise<CheckPageSlugResponse | null> {
  const params = new URLSearchParams();

  if (options?.excludeId !== undefined) {
    params.set("exclude_id", options.excludeId);
  }

  const queryString = params.toString();
  const queryParams = queryString.length > 0 ? `?${queryString}` : "";

  const result = await fetcher<CheckPageSlugResponse>(
    locale,
    `/profiles/${profileSlug}/pages/${pageSlug}/_check${queryParams}`,
    {
      method: "GET",
    },
  );

  return result;
}
