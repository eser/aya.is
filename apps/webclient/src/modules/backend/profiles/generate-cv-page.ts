import type { ProfilePage } from "@/modules/backend/types";
import { fetcher } from "@/modules/backend/fetcher";

export type GenerateCVPageResult =
  | { ok: true; data: ProfilePage }
  | { ok: false; error: string };

export async function generateCVPage(
  locale: string,
  profileSlug: string,
): Promise<GenerateCVPageResult> {
  try {
    const response = await fetcher<ProfilePage>(
      locale,
      `/profiles/${profileSlug}/_pages/generate-cv`,
      {
        method: "POST",
      },
    );

    if (response === null) {
      return { ok: false, error: "Failed to generate CV page" };
    }

    return { ok: true, data: response };
  } catch (err: unknown) {
    const message = err instanceof Error
      ? err.message
      : "Failed to generate CV page";

    return { ok: false, error: message };
  }
}
