// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// Pure helpers of the privacy settings page (kept free of React so they are unit-testable).
import type { PrivacyExportState } from "@/modules/backend/privacy/types";

/**
 * Deleting an account requires typing the username. The comparison ignores case and surrounding whitespace, the
 * same way the backend does, so a user is never blocked by the UI for something the server would accept.
 */
export function isDeletionConfirmed(typed: string, slug: string): boolean {
  if (slug === "") {
    return false;
  }

  return typed.trim().toLowerCase() === slug.toLowerCase();
}

export type ExportView = "unavailable" | "idle" | "preparing" | "ready" | "expired" | "failed";

export interface ExportViewModel {
  view: ExportView;
  canRequest: boolean;
  canDownload: boolean;
  // Set when a new request is blocked by the cooldown; the page tells the user when it opens up.
  nextRequestAt: string | null;
}

/** Reduces the backend's export state to the one thing the page should show. */
export function describeExport(state: PrivacyExportState): ExportViewModel {
  if (!state.available) {
    return { view: "unavailable", canRequest: false, canDownload: false, nextRequestAt: null };
  }

  const blockedByCooldown = !state.can_request && state.next_request_at !== null;
  const nextRequestAt = blockedByCooldown ? state.next_request_at : null;

  switch (state.status) {
    case "pending":
    case "processing":
      return { view: "preparing", canRequest: false, canDownload: false, nextRequestAt: null };
    case "ready":
      return { view: "ready", canRequest: state.can_request, canDownload: state.can_download, nextRequestAt };
    case "expired":
      return { view: "expired", canRequest: state.can_request, canDownload: false, nextRequestAt };
    case "failed":
      return { view: "failed", canRequest: state.can_request, canDownload: false, nextRequestAt: null };
    default:
      return { view: "idle", canRequest: state.can_request, canDownload: false, nextRequestAt };
  }
}
