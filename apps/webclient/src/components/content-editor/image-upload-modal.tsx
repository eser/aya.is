import * as React from "react";
import { Upload, ImageIcon, Loader2, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { backend } from "@/modules/backend/backend";
import type { UploadPurpose } from "@/modules/backend/types";
import styles from "./content-editor.module.css";
import { cn } from "@/lib/utils";

type ImageUploadModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onImageInsert: (url: string, alt: string) => void;
  locale: string;
  purpose?: UploadPurpose;
};

type UploadState =
  | { status: "idle" }
  | { status: "uploading"; progress: number }
  | { status: "success"; url: string }
  | { status: "error"; message: string };

const ALLOWED_TYPES = ["image/jpeg", "image/png", "image/gif", "image/webp"];
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB

export function ImageUploadModal(props: ImageUploadModalProps) {
  const {
    open,
    onOpenChange,
    onImageInsert,
    locale,
    purpose = "content-image",
  } = props;

  const [uploadState, setUploadState] = React.useState<UploadState>({
    status: "idle",
  });
  const [selectedFile, setSelectedFile] = React.useState<File | null>(null);
  const [altText, setAltText] = React.useState("");
  const [isDragging, setIsDragging] = React.useState(false);
  const fileInputRef = React.useRef<HTMLInputElement>(null);

  const resetState = () => {
    setUploadState({ status: "idle" });
    setSelectedFile(null);
    setAltText("");
    setIsDragging(false);
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      resetState();
    }
    onOpenChange(newOpen);
  };

  const validateFile = (file: File): string | null => {
    if (!ALLOWED_TYPES.includes(file.type)) {
      return "Please select a valid image file (JPEG, PNG, GIF, or WebP)";
    }
    if (file.size > MAX_FILE_SIZE) {
      return "File size must be less than 5MB";
    }
    return null;
  };

  const handleFileSelect = (file: File) => {
    const error = validateFile(file);
    if (error !== null) {
      setUploadState({ status: "error", message: error });
      return;
    }
    setSelectedFile(file);
    setUploadState({ status: "idle" });
    // Auto-populate alt text from filename
    const nameWithoutExt = file.name.replace(/\.[^/.]+$/, "").replace(/[-_]/g, " ");
    setAltText(nameWithoutExt);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files[0];
    if (file !== undefined) {
      handleFileSelect(file);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file !== undefined) {
      handleFileSelect(file);
    }
  };

  const handleUpload = async () => {
    if (selectedFile === null) return;

    setUploadState({ status: "uploading", progress: 0 });

    try {
      // Step 1: Get pre-signed URL
      const presignResponse = await backend.getPresignedURL(locale, {
        filename: selectedFile.name,
        content_type: selectedFile.type,
        purpose,
      });

      if (presignResponse === null) {
        setUploadState({
          status: "error",
          message: "Failed to get upload URL. Please try again.",
        });
        return;
      }

      setUploadState({ status: "uploading", progress: 50 });

      // Step 2: Upload file directly to S3/R2
      const uploadSuccess = await backend.uploadToPresignedURL(
        presignResponse.upload_url,
        selectedFile,
        selectedFile.type,
      );

      if (!uploadSuccess) {
        setUploadState({
          status: "error",
          message: "Failed to upload file. Please try again.",
        });
        return;
      }

      setUploadState({ status: "success", url: presignResponse.public_url });
    } catch {
      setUploadState({
        status: "error",
        message: "An unexpected error occurred. Please try again.",
      });
    }
  };

  const handleInsert = () => {
    if (uploadState.status === "success") {
      onImageInsert(uploadState.url, altText);
      handleOpenChange(false);
    }
  };

  const handleRemoveFile = () => {
    setSelectedFile(null);
    setUploadState({ status: "idle" });
    setAltText("");
    if (fileInputRef.current !== null) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Insert Image</DialogTitle>
          <DialogDescription>
            Upload an image to insert into your content.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {selectedFile === null ? (
            <div
              className={cn(
                styles.uploadDropzone,
                isDragging && styles.uploadDropzoneActive,
              )}
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onClick={() => fileInputRef.current?.click()}
            >
              <Upload className="mb-2 size-8 text-muted-foreground" />
              <p className="text-sm font-medium">
                Drop an image here or click to browse
              </p>
              <p className="text-xs text-muted-foreground">
                JPEG, PNG, GIF, or WebP up to 5MB
              </p>
              <input
                ref={fileInputRef}
                type="file"
                className="hidden"
                accept={ALLOWED_TYPES.join(",")}
                onChange={handleInputChange}
              />
            </div>
          ) : (
            <div className="space-y-4">
              <div className={styles.uploadPreview}>
                <img
                  src={URL.createObjectURL(selectedFile)}
                  alt="Preview"
                  className={styles.uploadPreviewImage}
                />
                <Button
                  variant="ghost"
                  size="icon-xs"
                  className="absolute top-2 right-2 bg-background/80"
                  onClick={handleRemoveFile}
                >
                  <X className="size-3" />
                </Button>
              </div>

              <div className="space-y-2">
                <Label htmlFor="alt-text">Alt Text</Label>
                <Input
                  id="alt-text"
                  value={altText}
                  onChange={(e) => setAltText(e.target.value)}
                  placeholder="Describe the image..."
                />
                <p className="text-xs text-muted-foreground">
                  Alt text helps with accessibility and SEO
                </p>
              </div>
            </div>
          )}

          {uploadState.status === "error" && (
            <p className="text-sm text-destructive">{uploadState.message}</p>
          )}

          {uploadState.status === "uploading" && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Loader2 className="size-4 animate-spin" />
              Uploading...
            </div>
          )}

          {uploadState.status === "success" && (
            <div className="flex items-center gap-2 text-sm text-green-600 dark:text-green-400">
              <ImageIcon className="size-4" />
              Image uploaded successfully
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => handleOpenChange(false)}>
            Cancel
          </Button>
          {uploadState.status === "success" ? (
            <Button onClick={handleInsert}>Insert Image</Button>
          ) : (
            <Button
              onClick={handleUpload}
              disabled={
                selectedFile === null || uploadState.status === "uploading"
              }
            >
              {uploadState.status === "uploading" ? (
                <>
                  <Loader2 className="mr-1.5 size-4 animate-spin" />
                  Uploading...
                </>
              ) : (
                <>
                  <Upload className="mr-1.5 size-4" />
                  Upload
                </>
              )}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
