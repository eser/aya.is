import { fetcher } from "../fetcher";
import type { StorySeries } from "../types";

export async function createSeries(
  locale: string,
  slug: string,
  title: string,
  description: string,
): Promise<StorySeries | null> {
  const response = await fetcher<StorySeries>(
    locale,
    "/series",
    {
      method: "POST",
      body: JSON.stringify({ slug, title, description }),
    },
  );

  return response;
}
