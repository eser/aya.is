// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";
import type { ProfileMembershipCandidate } from "../types";

export async function createApplication(
  locale: string,
  slug: string,
  applicantMessage: string | null,
  formResponses: Record<string, string>,
): Promise<ProfileMembershipCandidate | null> {
  const token = getAuthToken();
  if (token === null) return null;

  const response = await fetch(
    `${getBackendUri()}/${locale}/profiles/${slug}/_candidates/apply`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify({
        applicant_message: applicantMessage,
        form_responses: formResponses,
      }),
    },
  );

  if (!response.ok) return null;
  const result = await response.json();
  return result.data ?? null;
}
