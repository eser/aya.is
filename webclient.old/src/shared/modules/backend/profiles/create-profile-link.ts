import { getBackendUri } from "@/shared/config.ts";
import { getToken } from "@/shared/modules/auth/token-refresh.ts";
import type { ProfileLink } from "./list-profile-links.ts";

export type CreateProfileLinkRequest = {
  kind: string;
  uri?: string | null;
  title: string;
  is_hidden: boolean;
};

export async function createProfileLink(
  locale: string,
  slug: string,
  data: CreateProfileLinkRequest,
): Promise<ProfileLink | null> {
  const backendUri = getBackendUri();
  const token = getToken();

  const response = await fetch(`${backendUri}/${locale}/profiles/${slug}/_links`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(token !== null ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const errorText = await response.text().catch(() => "");
    throw new Error(errorText !== "" ? errorText : `Failed to create profile link: ${response.status}`);
  }

  const result = await response.json();
  return result.data as ProfileLink;
}
