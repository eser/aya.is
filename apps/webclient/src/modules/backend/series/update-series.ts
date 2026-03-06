import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function updateSeries(
  locale: string,
  id: string,
  data: { slug: string; series_picture_uri?: string | null },
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/series/${id}`,
    {
      method: "PATCH",
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
