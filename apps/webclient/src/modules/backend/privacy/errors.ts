// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type { PrivacyErrorCode, PrivacyOwnedProfile } from "./types";

export interface PrivacyErrorDetails {
  code: PrivacyErrorCode | null;
  message: string;
  nextRequestAt: string | null;
  profiles: PrivacyOwnedProfile[];
}

const KNOWN_CODES: readonly string[] = [
  "export_cooldown",
  "export_in_progress",
  "export_unavailable",
  "export_not_ready",
  "export_expired",
  "confirmation_mismatch",
  "owns_profiles",
];

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function isOwnedProfile(value: unknown): value is PrivacyOwnedProfile {
  return isRecord(value) && typeof value["slug"] === "string" && typeof value["id"] === "string";
}

// The fetcher throws `Error(message, { cause: { response } })` for non-OK JSON answers; the backend puts a
// machine-readable `code` (and extra data such as `next_request_at` or `profiles`) next to the message.
export function readPrivacyError(error: unknown): PrivacyErrorDetails {
  const details: PrivacyErrorDetails = { code: null, message: "", nextRequestAt: null, profiles: [] };

  if (error instanceof Error) {
    details.message = error.message;
  }

  const cause = error instanceof Error ? error.cause : undefined;
  const response = isRecord(cause) ? cause["response"] : undefined;

  if (!isRecord(response)) {
    return details;
  }

  const code = response["code"];
  if (typeof code === "string" && KNOWN_CODES.includes(code)) {
    details.code = code as PrivacyErrorCode;
  }

  const nextRequestAt = response["next_request_at"];
  if (typeof nextRequestAt === "string") {
    details.nextRequestAt = nextRequestAt;
  }

  const profiles = response["profiles"];
  if (Array.isArray(profiles)) {
    details.profiles = profiles.filter(isOwnedProfile);
  }

  return details;
}
