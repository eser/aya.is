// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fetcher } from "../fetcher";
import type { CandidateFormResponse } from "../types";

export async function getCandidateResponses(
  locale: string,
  slug: string,
  candidateId: string,
): Promise<CandidateFormResponse[] | null> {
  return await fetcher<CandidateFormResponse[]>(
    locale,
    `/profiles/${slug}/_candidates/${candidateId}/responses`,
  );
}
