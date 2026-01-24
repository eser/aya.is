import { fetcher } from "@/modules/backend/fetcher";
import type { InsertStoryInput, Story } from "@/modules/backend/types";

export type InsertStoryData = Story;

export async function insertStory(
  locale: string,
  profileSlug: string,
  input: InsertStoryInput,
): Promise<Story | null> {
  const response = await fetcher<InsertStoryData>(
    locale,
    `/profiles/${profileSlug}/_stories`,
    {
      method: "POST",
      body: JSON.stringify(input),
    },
  );
  return response;
}
