import { getBackendUri } from "@/config";
import { getAuthToken } from "../fetcher";

export async function deleteSeries(
  locale: string,
  id: string,
): Promise<boolean> {
  const token = getAuthToken();
  if (token === null) return false;

  const response = await fetch(
    `${getBackendUri()}/${locale}/series/${id}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    },
  );

  return response.ok;
}
