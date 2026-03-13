import { fetcher } from "../fetcher";
import type { ApplicationPreset } from "../types";

export async function listApplicationPresets(
  locale: string,
  slug: string,
): Promise<ApplicationPreset[] | null> {
  return await fetcher<ApplicationPreset[]>(
    locale,
    `/profiles/${slug}/_application-presets`,
  );
}
