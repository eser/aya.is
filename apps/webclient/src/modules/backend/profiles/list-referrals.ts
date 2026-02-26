import { fetcher } from "../fetcher";
import type { ProfileMembershipReferral } from "../types";

export async function listReferrals(
  locale: string,
  slug: string,
): Promise<ProfileMembershipReferral[] | null> {
  return await fetcher<ProfileMembershipReferral[]>(
    locale,
    `/profiles/${slug}/_referrals`,
  );
}
