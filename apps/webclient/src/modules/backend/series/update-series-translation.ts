import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function updateSeriesTranslation(
  locale: string,
  seriesId: string,
  txLocale: string,
  data: { title: string; description: string },
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/series/${seriesId}/translations/${txLocale}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
      body: JSON.stringify(data),
    },
  );

  return response.ok;
}
