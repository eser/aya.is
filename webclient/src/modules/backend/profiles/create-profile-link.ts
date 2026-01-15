import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileLink } from "../types";

export type CreateProfileLinkRequest = {
  kind: string;
  uri?: string | null;
  title: string;
  is_hidden: boolean;
};

export async function createProfileLink(
  locale: string,
  slug: string,
  data: CreateProfileLinkRequest
): Promise<ProfileLink | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(`${getBackendUri()}/${locale}/profiles/${slug}/_links`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    credentials: "include",
    body: JSON.stringify(data),
  });

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
