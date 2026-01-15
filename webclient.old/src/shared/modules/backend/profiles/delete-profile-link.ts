import { getBackendUri } from "@/shared/config.ts";
import { getToken } from "@/shared/modules/auth/token-refresh.ts";

export type DeleteProfileLinkResponse = {
  success: boolean;
  message: string;
};

export async function deleteProfileLink(
  locale: string,
  slug: string,
  linkId: string,
): Promise<DeleteProfileLinkResponse | null> {
  const backendUri = getBackendUri();
  const token = getToken();

  const response = await fetch(`${backendUri}/${locale}/profiles/${slug}/_links/${linkId}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      ...(token !== null ? { Authorization: `Bearer ${token}` } : {}),
    },
  });

  if (!response.ok) {
    const errorText = await response.text().catch(() => "");
    throw new Error(errorText !== "" ? errorText : `Failed to delete profile link: ${response.status}`);
  }

  const result = await response.json();
  return result.data as DeleteProfileLinkResponse;
}
