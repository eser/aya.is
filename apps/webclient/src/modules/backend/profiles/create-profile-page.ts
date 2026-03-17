// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { ContentVisibility, ProfilePage } from "../types";

export type CreateProfilePageRequest = {
  slug: string;
  title: string;
  summary: string;
  content: string;
  cover_picture_uri: string | null;
  published_at: string | null;
  visibility: ContentVisibility;
};

export async function createProfilePage(
  locale: string,
  profileSlug: string,
  data: CreateProfilePageRequest,
): Promise<ProfilePage | null> {
  return await fetcher<ProfilePage>(locale, `/profiles/${profileSlug}/_pages`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}
