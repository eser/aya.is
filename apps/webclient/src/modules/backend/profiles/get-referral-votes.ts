import { fetcher } from "../fetcher";
import type { ReferralVote } from "../types";

export async function getReferralVotes(
  locale: string,
  slug: string,
  referralId: string,
): Promise<ReferralVote[] | null> {
  return await fetcher<ReferralVote[]>(
    locale,
    `/profiles/${slug}/_referrals/${referralId}/votes`,
  );
}
