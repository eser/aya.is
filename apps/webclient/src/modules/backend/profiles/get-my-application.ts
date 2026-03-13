import { fetcher } from "../fetcher";
import type { ProfileMembershipCandidate } from "../types";

export async function getMyApplication(
  locale: string,
  slug: string,
): Promise<ProfileMembershipCandidate | null> {
  return await fetcher<ProfileMembershipCandidate>(
    locale,
    `/profiles/${slug}/_candidates/my-application`,
  );
}
