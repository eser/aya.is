import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function deleteProfileTeam(
  locale: string,
  slug: string,
  teamId: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_teams/${teamId}`,
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
