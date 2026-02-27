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
