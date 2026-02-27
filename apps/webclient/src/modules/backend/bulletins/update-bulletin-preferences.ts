import { fetcher } from "../fetcher";
import type { UpdateBulletinPreferencesRequest } from "./types";

export async function updateBulletinPreferences(
  locale: string,
  data: UpdateBulletinPreferencesRequest,
): Promise<{ updated: boolean } | null> {
  return await fetcher<{ updated: boolean }>(
    locale,
    "/bulletin/preferences",
    {
      method: "PUT",
      body: JSON.stringify(data),
    },
  );
}
