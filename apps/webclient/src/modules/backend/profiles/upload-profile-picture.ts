// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { uploadFetcher } from "../fetcher";

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

  return await uploadFetcher<UploadProfilePictureResponse>(
    locale,
    `/profiles/${slug}/_upload-picture`,
    formData,
  );
}
