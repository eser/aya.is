import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileResource } from "../types";

export interface CreateProfileResourceInput {
  kind: string;
  remote_id?: string;
  public_id?: string;
  url?: string;
  title: string;
  description?: string | null;
  properties?: Record<string, unknown>;
}

export async function createProfileResource(
  locale: string,
  slug: string,
  input: CreateProfileResourceInput,
): Promise<ProfileResource | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_resources`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify(input),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data;
}
