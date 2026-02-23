import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function setMembershipTeams(
  locale: string,
  slug: string,
  membershipId: string,
  teamIds: string[],
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_memberships/${membershipId}/teams`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ team_ids: teamIds }),
    },
  );

  return response.ok;
}
