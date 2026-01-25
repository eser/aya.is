import { fetcher } from "../fetcher";

export type CheckPageSlugResponse = {
  available: boolean;
  message?: string;
  severity?: "error" | "warning" | "";
};

export async function checkPageSlug(
  locale: string,
  profileSlug: string,
  pageSlug: string,
  excludeId?: string,
): Promise<CheckPageSlugResponse | null> {
  const queryParams = excludeId !== undefined ? `?exclude_id=${excludeId}` : "";
  const result = await fetcher<CheckPageSlugResponse>(
    locale,
    `/profiles/${profileSlug}/pages/${pageSlug}/_check${queryParams}`,
    {
      method: "GET",
    },
  );

  return result;
}
