import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function deleteProfileResource(
  locale: string,
  slug: string,
  resourceId: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_resources/${resourceId}`,
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
