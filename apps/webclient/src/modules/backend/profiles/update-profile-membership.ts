import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { MembershipKind } from "../types";

export type UpdateMembershipInput = {
  kind: MembershipKind;
};

export async function updateProfileMembership(
  locale: string,
  slug: string,
  membershipId: string,
  input: UpdateMembershipInput,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_memberships/${membershipId}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify(input),
    },
  );

  return response.ok;
}
