import { fetcher } from "../fetcher";
import type { StorySeries } from "../types";

export async function getSeriesList(
  locale: string,
): Promise<StorySeries[] | null> {
  const response = await fetcher<StorySeries[]>(
    locale,
    "/series",
  );

  return response;
}
