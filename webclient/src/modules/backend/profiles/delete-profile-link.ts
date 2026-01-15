import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export type DeleteProfileLinkResponse = {
  success: boolean;
  message: string;
};

export async function deleteProfileLink(
  locale: string,
  slug: string,
  linkId: string
): Promise<DeleteProfileLinkResponse | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/${linkId}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    }
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
