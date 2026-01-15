import { getBackendUri } from "@/shared/config.ts";

export interface AuthCallbackRequest {
  code: string;
  state?: string;
}

export interface AuthCallbackResponse {
  token: string;
  session?: unknown;
}

export async function handleAuthCallback(
  locale: string,
  { code, state }: AuthCallbackRequest,
): Promise<AuthCallbackResponse | null> {
  const backendUri = getBackendUri();

  const response = await fetch(
    `${backendUri}/${locale}/auth/github/callback?code=${code}&state=${state !== null ? state : ""}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    },
  );

  if (!response.ok) {
    let errorMessage = `Authentication failed (${response.status})`;

    try {
      const errorData = await response.text();
      if (errorData) {
        errorMessage += `: ${errorData}`;
      }
    } catch {
      // If we can't read the error response, use the generic message
    }

    throw new Error(errorMessage);
  }

  const data = await response.json();

  if (data.data?.token === undefined) {
    // Include the full response data for debugging
    throw new Error(`No token received from server. Response: ${JSON.stringify(data)}`);
  }

  return data.data as AuthCallbackResponse;
}
