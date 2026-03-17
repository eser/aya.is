// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";

export async function removeInteraction(
  locale: string,
  storySlug: string,
  kind: string,
): Promise<unknown> {
  const response = await fetcher<unknown>(
    locale,
    `/stories/${storySlug}/interactions/${kind}`,
    {
      method: "DELETE",
    },
  );

  return response;
}
