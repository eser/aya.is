import { fetcher } from "@/shared/modules/backend/fetcher.ts";

export type ProfilePermissions = {
  can_edit: boolean;
};

export async function getProfilePermissions(
  locale: string,
  slug: string,
): Promise<ProfilePermissions | null> {
  const response = await fetcher<ProfilePermissions>(
    `/${locale}/profiles/${slug}/_permissions`,
  );

  return response;
}
