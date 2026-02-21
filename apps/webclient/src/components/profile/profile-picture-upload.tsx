import * as React from "react";
import { useTranslation } from "react-i18next";
import { Camera, Loader2, Trash2, User } from "lucide-react";
import { toast } from "sonner";
import { backend } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { cn } from "@/lib/utils";

type ProfilePictureUploadProps = {
  currentImageUri: string | null | undefined;
  profileSlug: string;
  profileTitle: string;
  locale: string;
  onUploadComplete: (newUri: string) => void;
  onRemoveComplete?: () => void;
};

const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
const ACCEPTED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/gif", "image/webp"];

export function ProfilePictureUpload(props: ProfilePictureUploadProps) {
  const { t } = useTranslation();
  const fileInputRef = React.useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = React.useState(false);
  const [isRemoving, setIsRemoving] = React.useState(false);
  const [previewUrl, setPreviewUrl] = React.useState<string | null>(null);

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file === undefined) {
      return;
    }

    // Validate file type
    if (!ACCEPTED_IMAGE_TYPES.includes(file.type)) {
      toast.error(t("Profile.Invalid file type"));
      return;
    }

    // Validate file size
    if (file.size > MAX_FILE_SIZE) {
      toast.error(t("Profile.File too large"));
      return;
    }

    // Create preview
    const objectUrl = URL.createObjectURL(file);
    setPreviewUrl(objectUrl);

    // Upload the file
    handleUpload(file);
  };

  const handleUpload = async (file: File) => {
    setIsUploading(true);
    try {
      // Step 1: Get presigned URL
      const presignResponse = await backend.getPresignedURL(props.locale, {
        filename: file.name,
        content_type: file.type,
        purpose: "profile-picture",
      });

      if (presignResponse === null) {
        toast.error(t("Profile.Failed to upload picture"));
        setPreviewUrl(null);
        return;
      }

      // Step 2: Upload file directly to S3/R2
      const uploadSuccess = await backend.uploadToPresignedURL(
        presignResponse.upload_url,
        file,
        file.type,
      );

      if (!uploadSuccess) {
        toast.error(t("Profile.Failed to upload picture"));
        setPreviewUrl(null);
        return;
      }

      // Step 3: Update profile with new picture URI
      const updateResult = await backend.updateProfilePicture(
        props.locale,
        props.profileSlug,
        presignResponse.public_url,
      );

      if (updateResult !== null) {
        toast.success(t("Profile.Picture uploaded successfully"));
        props.onUploadComplete(presignResponse.public_url);
        setPreviewUrl(null);
      } else {
        toast.error(t("Profile.Failed to upload picture"));
        setPreviewUrl(null);
      }
    } catch {
      toast.error(t("Profile.Failed to upload picture"));
      setPreviewUrl(null);
    } finally {
      setIsUploading(false);
      // Reset file input
      if (fileInputRef.current !== null) {
        fileInputRef.current.value = "";
      }
    }
  };

  const handleRemove = async () => {
    setIsRemoving(true);
    try {
      // Send empty string to signal "remove picture" to the backend
      const result = await backend.updateProfilePicture(
        props.locale,
        props.profileSlug,
        "",
      );

      if (result !== null) {
        toast.success(t("Profile.Picture removed successfully"));
        setPreviewUrl(null);
        if (props.onRemoveComplete !== undefined) {
          props.onRemoveComplete();
        }
      } else {
        toast.error(t("Profile.Failed to remove picture"));
      }
    } catch {
      toast.error(t("Profile.Failed to remove picture"));
    } finally {
      setIsRemoving(false);
    }
  };

  const handleClick = () => {
    if (!isUploading && !isRemoving && fileInputRef.current !== null) {
      fileInputRef.current.click();
    }
  };

  const isBusy = isUploading || isRemoving;

  // Determine which image to display
  const displayImageUri = previewUrl ?? props.currentImageUri;
  const hasImage = displayImageUri !== null && displayImageUri !== undefined;

  return (
    <div className="flex flex-col items-center gap-4">
      <div className="relative group">
        {/* Profile Picture or Placeholder */}
        <button
          type="button"
          onClick={handleClick}
          disabled={isBusy}
          className={cn(
            "relative size-32 rounded-full overflow-hidden border-2 border-dashed border-muted-foreground/25 hover:border-muted-foreground/50 transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
            isBusy && "cursor-not-allowed opacity-50",
          )}
        >
          {hasImage ? (
            <img
              src={displayImageUri}
              alt={props.profileTitle}
              className="size-full object-cover"
            />
          ) : (
            <div className="size-full flex items-center justify-center bg-muted">
              <User className="size-12 text-muted-foreground" />
            </div>
          )}

          {/* Overlay */}
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
            {isBusy ? (
              <Loader2 className="size-8 text-white animate-spin" />
            ) : (
              <Camera className="size-8 text-white" />
            )}
          </div>
        </button>

        {/* Hidden File Input */}
        <input
          ref={fileInputRef}
          type="file"
          accept={ACCEPTED_IMAGE_TYPES.join(",")}
          onChange={handleFileSelect}
          className="hidden"
          aria-label={t("Profile.Change Picture")}
        />
      </div>

      {/* Action Buttons */}
      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleClick}
          disabled={isBusy}
        >
          {isUploading ? (
            <>
              <Loader2 className="size-4 mr-2 animate-spin" />
              {t("Profile.Uploading...")}
            </>
          ) : (
            t("Profile.Change Picture")
          )}
        </Button>

        {hasImage && (
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button
                type="button"
                variant="outline"
                size="sm"
                disabled={isBusy}
              >
                {isRemoving ? (
                  <>
                    <Loader2 className="size-4 mr-2 animate-spin" />
                    {t("Profile.Removing...")}
                  </>
                ) : (
                  <>
                    <Trash2 className="size-4 mr-2" />
                    {t("Profile.Remove Picture")}
                  </>
                )}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>{t("Profile.Remove Picture")}</AlertDialogTitle>
                <AlertDialogDescription>
                  {t("Profile.Are you sure you want to remove your profile picture?")}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
                <AlertDialogAction onClick={handleRemove}>
                  {t("Profile.Remove Picture")}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        )}
      </div>
    </div>
  );
}
