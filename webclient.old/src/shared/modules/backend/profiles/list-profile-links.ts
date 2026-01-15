import { getBackendUri } from "@/shared/config.ts";
import { getToken } from "@/shared/modules/auth/token-refresh.ts";

export type ProfileLink = {
  id: string;
  kind: string;
  profile_id: string;
  order: number;
  is_managed: boolean;
  is_verified: boolean;
  is_hidden: boolean;
  remote_id?: string | null;
  public_id?: string | null;
  uri?: string | null;
  title: string;
  created_at: string;
  updated_at?: string | null;
  deleted_at?: string | null;
};

export async function listProfileLinks(
  locale: string,
  slug: string,
): Promise<ProfileLink[] | null> {
  const backendUri = getBackendUri();
  const token = getToken();

  const response = await fetch(`${backendUri}/${locale}/profiles/${slug}/_links`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      ...(token !== null ? { Authorization: `Bearer ${token}` } : {}),
    },
  });

  if (!response.ok) {
    const errorText = await response.text().catch(() => "");
    throw new Error(errorText !== "" ? errorText : `Failed to list profile links: ${response.status}`);
  }

  const result = await response.json();
  return result.data as ProfileLink[];
}
