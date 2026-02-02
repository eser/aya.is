import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function deleteProfileMembership(
  locale: string,
  slug: string,
  membershipId: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_memberships/${membershipId}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    },
  );

  return response.ok;
}
