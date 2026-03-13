import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { MembershipKind } from "../types";

export type AddMembershipInput = {
  member_profile_id: string;
  kind: MembershipKind;
};

export async function addProfileMembership(
  locale: string,
  slug: string,
  input: AddMembershipInput,
): Promise<{ id: string } | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_memberships`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify(input),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data ?? null;
}
