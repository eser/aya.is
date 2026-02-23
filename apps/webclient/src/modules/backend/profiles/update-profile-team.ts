import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function updateProfileTeam(
  locale: string,
  slug: string,
  teamId: string,
  name: string,
  description?: string | null,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_teams/${teamId}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({ name, description: description ?? null }),
    },
  );

  return response.ok;
}
