// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
export type PrivacyExportStatus = "none" | "pending" | "processing" | "ready" | "failed" | "expired";

export interface PrivacyExportState {
  requested_at: string | null;
  ready_at: string | null;
  expires_at: string | null;
  // Earliest moment a new export can be requested; null when the user never requested one (or the last one failed).
  next_request_at: string | null;
  id: string | null;
  status: PrivacyExportStatus;
  size_bytes: number | null;
  // False when the server has no private storage configured for archives.
  available: boolean;
  can_request: boolean;
  can_download: boolean;
}

export interface PrivacyState {
  profile_slug: string;
  export: PrivacyExportState;
  is_hidden: boolean;
}

export interface PrivacyOwnedProfile {
  id: string;
  slug: string;
  kind: string;
  title: string;
}

export interface RequestedPrivacyExport {
  id: string;
  status: PrivacyExportStatus;
  requested_at: string;
}

export interface DeleteAccountRequest {
  confirm_slug: string;
  delete_content: boolean;
}

// Machine-readable reasons the backend attaches to a failed privacy request.
export type PrivacyErrorCode =
  | "export_cooldown"
  | "export_in_progress"
  | "export_unavailable"
  | "export_not_ready"
  | "export_expired"
  | "confirmation_mismatch"
  | "owns_profiles";
