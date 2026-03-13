import { fetcher } from "../fetcher";
import type { ApplicationForm } from "../types";

export async function getApplicationForm(
  locale: string,
  slug: string,
): Promise<ApplicationForm | null> {
  return await fetcher<ApplicationForm>(
    locale,
    `/profiles/${slug}/_application-form`,
  );
}
