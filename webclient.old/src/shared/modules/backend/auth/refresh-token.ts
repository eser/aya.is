import { getBackendUri } from "@/shared/config.ts";

export interface RefreshTokenData {
  token: string;
  expiresAt: number;
}

export async function refreshTokenRequest(locale: string, currentToken: string) {
  const backendUri = getBackendUri();
  const response = await fetch(`${backendUri}/${locale}/auth/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${currentToken}`,
    },
  });

  if (!response.ok) {
    throw new Error(`Token refresh failed: ${response.status}`);
  }

  const data = await response.json();
  return data.data as RefreshTokenData;
}
