import { getBackendUri } from "@/shared/config.ts";
import { getToken } from "@/shared/modules/auth/token-refresh.ts";

export interface CreateProfileRequest {
  kind: string;
  slug: string;
  title: string;
  description: string;
}

export interface CreateProfileResponse {
  id: string;
  slug: string;
  kind: string;
  title: string;
  description: string;
  created_at: string;
}

export async function createProfile(
  locale: string,
  data: CreateProfileRequest,
): Promise<CreateProfileResponse | null> {
  const backendUri = getBackendUri();
  const token = getToken();

  const response = await fetch(`${backendUri}/${locale}/profiles/_create`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(token !== null ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const errorText = await response.text().catch(() => "");

    // Handle specific error for duplicate individual profile
    if (response.status === 400 && errorText.includes("User already has an individual profile")) {
      throw new Error("You already have an individual profile. Please select a different profile type.");
    }

    // Handle other specific errors
    if (response.status === 400 && errorText.includes("Slug is already taken")) {
      throw new Error("This profile URL is already taken. Please choose a different slug.");
    }

    // Generic error handling
    throw new Error(errorText !== "" ? errorText : `Failed to create profile: ${response.status}`);
  }

  const result = await response.json();
  return result.data as CreateProfileResponse;
}
