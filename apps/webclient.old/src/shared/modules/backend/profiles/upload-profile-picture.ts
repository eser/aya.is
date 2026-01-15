import { fetcher } from "@/shared/modules/backend/fetcher.ts";

export type UploadProfilePictureResponse = {
  profile_picture_uri: string;
  message: string;
};

export async function uploadProfilePicture(
  locale: string,
  slug: string,
  file: File,
): Promise<UploadProfilePictureResponse | null> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetcher<UploadProfilePictureResponse>(
    `/${locale}/profiles/${slug}/_upload-picture`,
    {
      method: "POST",
      body: formData,
    },
  );

  return response;
}
