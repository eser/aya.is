"use client";

import { useRef, useState } from "react";

import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { Alert, AlertDescription } from "@/shared/components/ui/alert.tsx";

import styles from "./profile-picture-upload.module.css";

const SAFE_PROTOCOLS = ["https:", "http:", "data:", "blob:"];
function sanitizeImageSrc(src: string | null | undefined): string {
  if (src === null || src === undefined || src === "") {
    return "";
  }
  try {
    const parsed = new URL(src, "https://placeholder.invalid");
    if (SAFE_PROTOCOLS.includes(parsed.protocol)) {
      return src;
    }
  } catch {
    if (src.startsWith("/") && !src.startsWith("//")) {
      return src;
    }
  }
  return "";
}

type ProfilePictureUploadProps = {
  currentProfilePictureURI?: string | null;
  onUploadComplete: (newProfilePictureURI: string) => void;
  slug: string;
  locale: string;
};

export function ProfilePictureUpload(props: ProfilePictureUploadProps) {
  const { t } = useTranslations();
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = async (file: File) => {
    // Validate file type
    if (!file.type.startsWith("image/")) {
      setError(
        t("Profile", "Please select an image file") ??
          "Please select an image file",
      );
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      setError(
        t("Profile", "File size must be less than 5MB") ??
          "File size must be less than 5MB",
      );
      return;
    }

    setError(null);
    setIsUploading(true);

    try {
      // Create preview
      const previewUrl = URL.createObjectURL(file);
      setPreview(previewUrl);

      // TODO: Upload file to server
      // For now, we'll simulate the upload process
      await new Promise((resolve) => setTimeout(resolve, 2000));

      // This would be the actual upload:
      // const response = await backend.uploadProfilePicture(props.locale, props.slug, file);
      // if (response?.profile_picture_uri) {
      //   props.onUploadComplete(response.profile_picture_uri);
      // }

      // Simulated success
      const simulatedURI = previewUrl; // In real implementation, this would be the server-returned URI
      props.onUploadComplete(simulatedURI);
    } catch (error) {
      console.error("Upload failed:", error);
      setError(
        t("Profile", "Failed to upload image") ?? "Failed to upload image",
      );
      setPreview(null);
    } finally {
      setIsUploading(false);
    }
  };

  const handleFileInputChange = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0];
    if (file) {
      handleFileSelect(file);
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const handleRemove = () => {
    setIsUploading(true);
    setError(null);

    try {
      // Remove profile picture by setting it to null - parent will handle the backend call
      props.onUploadComplete("");
      setPreview(null);
    } catch (error) {
      console.error("Failed to remove profile picture:", error);
      setError(
        t("Profile", "Failed to remove profile picture") ??
          "Failed to remove profile picture",
      );
    } finally {
      setIsUploading(false);
    }
  };

  const currentImage = preview || props.currentProfilePictureURI;

  return (
    <div className={styles.container}>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileInputChange}
        style={{ display: "none" }}
      />

      {error && (
        <Alert variant="destructive" className={styles.alert}>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className={styles.uploadArea}>
        {currentImage
          ? (
            <div className={styles.imagePreview}>
              <img
                src={sanitizeImageSrc(currentImage)}
                alt={t("Profile", "Profile picture") ?? "Profile picture"}
                className={styles.image}
              />
              {isUploading && (
                <div className={styles.uploadingOverlay}>
                  <div className={styles.spinner} />
                  <span className={styles.uploadingText}>
                    {t("Profile", "Uploading...") ?? "Uploading..."}
                  </span>
                </div>
              )}
            </div>
          )
          : (
            <div
              className={`${styles.placeholder} ${isUploading ? styles.placeholderDisabled : ""}`}
              onClick={handleUploadClick}
            >
              <div className={styles.placeholderContent}>
                <div className={styles.uploadIcon}>📷</div>
                <p className={styles.placeholderText}>
                  {t("Profile", "Click to upload an image") ??
                    "Click to upload an image"}
                </p>
                <p className={styles.placeholderHint}>
                  {t("Profile", "Supports: JPG, PNG, GIF, WebP (max 5MB)") ??
                    "Supports: JPG, PNG, GIF, WebP (max 5MB)"}
                </p>
              </div>
            </div>
          )}
      </div>

      <div className={styles.actions}>
        {currentImage
          ? (
            <>
              <Button
                type="button"
                variant="outline"
                onClick={handleUploadClick}
                disabled={isUploading}
              >
                {t("Profile", "Change Image") ?? "Change Image"}
              </Button>

              <Button
                type="button"
                variant="destructive"
                onClick={handleRemove}
                disabled={isUploading}
              >
                {t("Profile", "Remove") ?? "Remove"}
              </Button>
            </>
          )
          : (
            <Button
              type="button"
              onClick={handleUploadClick}
              disabled={isUploading}
            >
              {t("Profile", "Upload Image") ?? "Upload Image"}
            </Button>
          )}
      </div>
    </div>
  );
}
