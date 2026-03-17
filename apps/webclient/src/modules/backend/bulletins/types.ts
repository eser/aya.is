// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
export interface BulletinPreferences {
  frequency: "daily" | "bidaily" | "weekly" | null;
  preferred_time: number | null;
  channels: string[];
  email: string | null;
  telegram_connected: boolean;
  last_bulletin_at: string | null;
}

export interface UpdateBulletinPreferencesRequest {
  frequency: string;
  preferred_time: number;
  channels: string[];
}
