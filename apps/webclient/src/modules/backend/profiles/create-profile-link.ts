import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileLink, LinkVisibility } from "../types";

export type CreateProfileLinkRequest = {
  kind: string;
  uri?: string | null;
  title: string;
  group?: string | null;
  description?: string | null;
  is_featured?: boolean;
  visibility?: LinkVisibility;
};

export async function createProfileLink(
  locale: string,
  slug: string,
  data: CreateProfileLinkRequest,
): Promise<ProfileLink | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify(data),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
