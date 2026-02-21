import type { ProfilePage } from "@/modules/backend/types";
import { fetcher } from "@/modules/backend/fetcher";

export async function generateCVPage(
  locale: string,
  profileSlug: string,
): Promise<ProfilePage | null> {
  const response = await fetcher<ProfilePage>(
    locale,
    `/profiles/${profileSlug}/_pages/generate-cv`,
    {
      method: "POST",
    },
  );
  return response;
}
