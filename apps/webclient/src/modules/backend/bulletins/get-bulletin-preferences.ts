import { fetcher } from "../fetcher";
import type { BulletinPreferences } from "./types";

export async function getBulletinPreferences(
  locale: string,
): Promise<BulletinPreferences | null> {
  return await fetcher<BulletinPreferences>(
    locale,
    "/bulletin/preferences",
  );
}
