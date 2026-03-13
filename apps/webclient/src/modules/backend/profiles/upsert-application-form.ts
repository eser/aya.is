import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ApplicationForm } from "../types";

export interface UpsertApplicationFormInput {
  preset_key?: string | null;
  fields: {
    label: string;
    field_type: string;
    is_required: boolean;
    sort_order: number;
    placeholder?: string | null;
  }[];
  responses_visibility: "members" | "leads";
}

export async function upsertApplicationForm(
  locale: string,
  slug: string,
  input: UpsertApplicationFormInput,
): Promise<ApplicationForm | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_application-form`,
    {
      method: "PUT",
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
  return result.data ?? null;
}
