import { fetcher } from "../fetcher";
import type { DateProposal } from "../types";

export async function createDateProposal(
  locale: string,
  storySlug: string,
  datetimeStart: string,
  datetimeEnd?: string,
): Promise<DateProposal | null> {
  const response = await fetcher<DateProposal>(
    locale,
    `/stories/${storySlug}/date-proposals`,
    {
      method: "POST",
      body: JSON.stringify({
        datetime_start: datetimeStart,
        datetime_end: datetimeEnd ?? null,
      }),
    },
  );

  return response;
}
