import { fetcher } from "@/shared/modules/backend/fetcher.ts";

export async function checkProfileSlug(locale: string, slug: string): Promise<boolean> {
  const response = await fetcher<{ exists: boolean }>(`/${locale}/profiles/${encodeURIComponent(slug)}/_check`);

  if (response === null) {
    throw new Error("Failed to check slug availability");
  }

  return response.exists;
}
