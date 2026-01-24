import { fetcher } from "../fetcher";
import type { Profile } from "../types";

export async function getProfilesByKinds(
  locale: string,
  kinds: string[],
): Promise<Profile[] | null> {
  const kindsParam = kinds.join(",");
  return await fetcher<Profile[]>(locale, `/profiles?filter_kind=${kindsParam}`);
}
