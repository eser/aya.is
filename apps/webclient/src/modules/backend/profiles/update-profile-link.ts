import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileLink, LinkVisibility } from "../types";

export type UpdateProfileLinkRequest = {
  kind: string;
  order: number;
  uri?: string | null;
  title: string;
  icon?: string | null;
  group?: string | null;
  description?: string | null;
  is_featured?: boolean;
  visibility?: LinkVisibility;
};

export async function updateProfileLink(
  locale: string,
  slug: string,
  linkId: string,
  data: UpdateProfileLinkRequest,
): Promise<ProfileLink | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_links/${linkId}`,
    {
      method: "PATCH",
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
