import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { Profile } from "../types";

export type CreateProfileRequest = {
  kind: "individual" | "organization" | "project" | "product";
  slug: string;
  title: string;
  description: string;
};

export async function createProfile(
  locale: string,
  data: CreateProfileRequest
): Promise<Profile | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(`${getBackendUri()}/${locale}/profiles/_create`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    credentials: "include",
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const errorData = await response.json();
    if (
      errorData.error?.includes("duplicate") ||
      errorData.error?.includes("already exists")
    ) {
      throw new Error("DUPLICATE_INDIVIDUAL_PROFILE");
    }
    throw new Error(errorData.error || "Failed to create profile");
  }

  const result = await response.json();
  return result.data;
}
