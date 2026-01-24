import { fetcher } from "@/modules/backend/fetcher";
import type {
  GetPresignedURLRequest,
  GetPresignedURLResponse,
} from "@/modules/backend/types";

export async function getPresignedURL(
  locale: string,
  request: GetPresignedURLRequest,
): Promise<GetPresignedURLResponse | null> {
  const response = await fetcher<GetPresignedURLResponse>(
    locale,
    `/site/uploads/presign`,
    {
      method: "POST",
      body: JSON.stringify(request),
    },
  );
  return response;
}
