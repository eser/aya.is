/**
 * Upload a file directly to S3/R2 using a pre-signed URL
 * This bypasses our backend - the file goes directly to storage
 */
export async function uploadToPresignedURL(
  presignedURL: string,
  file: File,
  contentType: string,
): Promise<boolean> {
  try {
    const response = await fetch(presignedURL, {
      method: "PUT",
      body: file,
      headers: {
        "Content-Type": contentType,
      },
    });

    return response.ok;
  } catch {
    return false;
  }
}
