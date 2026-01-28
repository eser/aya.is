// Background Image Picker Component
// Search and select background images from Unsplash

import * as React from "react";
import { useTranslation } from "react-i18next";
import { Search, X, Check, Loader2 } from "lucide-react";
import { backend, type UnsplashPhoto } from "@/modules/backend/backend.ts";
import { Button } from "@/components/ui/button.tsx";
import { Input } from "@/components/ui/input.tsx";
import styles from "./cover-generator.module.css";

interface BackgroundImagePickerProps {
  locale: string;
  selectedImageUrl: string | null;
  onSelect: (imageUrl: string | null, photographerName: string | null) => void;
}

export function BackgroundImagePicker(props: BackgroundImagePickerProps) {
  const { t } = useTranslation();
  const { locale, selectedImageUrl, onSelect } = props;

  const [query, setQuery] = React.useState("");
  const [searchResults, setSearchResults] = React.useState<UnsplashPhoto[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [hasSearched, setHasSearched] = React.useState(false);
  const [selectedPhotographer, setSelectedPhotographer] = React.useState<string | null>(null);

  const handleSearch = React.useCallback(async () => {
    if (query.trim().length === 0) {
      return;
    }

    setIsLoading(true);
    setHasSearched(true);

    try {
      const result = await backend.searchBackgroundImages(locale, {
        query: query.trim(),
        perPage: 12,
      });

      if (result !== null) {
        setSearchResults(result.results);
      } else {
        setSearchResults([]);
      }
    } catch (error) {
      console.error("Failed to search background images:", error);
      setSearchResults([]);
    } finally {
      setIsLoading(false);
    }
  }, [locale, query]);

  const handleKeyDown = React.useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === "Enter") {
        handleSearch();
      }
    },
    [handleSearch],
  );

  const handleSelectImage = React.useCallback(
    (photo: UnsplashPhoto) => {
      // Use the "regular" size for the cover (1080px wide)
      onSelect(photo.urls.regular, photo.user.name);
      setSelectedPhotographer(photo.user.name);
    },
    [onSelect],
  );

  const handleClearImage = React.useCallback(() => {
    onSelect(null, null);
    setSelectedPhotographer(null);
  }, [onSelect]);

  // Show selected image preview
  if (selectedImageUrl !== null) {
    return (
      <div className={styles.backgroundImagePicker}>
        <div className={styles.selectedImagePreview}>
          <img
            src={selectedImageUrl}
            alt={t("CoverDesigner.Selected background")}
            className={styles.selectedImagePreviewImg}
          />
          <Button
            variant="destructive"
            size="icon"
            className={styles.selectedImageClear}
            onClick={handleClearImage}
          >
            <X className="size-4" />
          </Button>
        </div>
        {selectedPhotographer !== null && (
          <p className={styles.searchHint}>
            {t("CoverDesigner.Photo by")} {selectedPhotographer}
          </p>
        )}
      </div>
    );
  }

  return (
    <div className={styles.backgroundImagePicker}>
      <div className={styles.searchInputWrapper}>
        <Input
          className={styles.searchInput}
          placeholder={t("CoverDesigner.Search images...")}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <Button
          variant="outline"
          size="icon"
          className={styles.searchButton}
          onClick={handleSearch}
          disabled={isLoading || query.trim().length === 0}
        >
          {isLoading ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <Search className="size-4" />
          )}
        </Button>
      </div>

      {isLoading && (
        <div className={styles.loadingSpinner}>
          <Loader2 className="size-6 animate-spin text-muted-foreground" />
        </div>
      )}

      {!isLoading && searchResults.length > 0 && (
        <div className={styles.imageGrid}>
          {searchResults.map((photo) => (
            <button
              key={photo.id}
              type="button"
              className={`${styles.imageGridItem} ${
                selectedImageUrl === photo.urls.regular ? styles.imageGridItemSelected : ""
              }`}
              onClick={() => handleSelectImage(photo)}
            >
              <img
                src={photo.urls.small}
                alt={photo.alt_description ?? photo.description ?? "Background image"}
                className={styles.imageGridItemImg}
              />
              <div className={styles.imageGridItemOverlay}>
                <Check className={styles.imageGridItemCheck} />
              </div>
              <span className={styles.photographerAttribution}>
                {photo.user.name}
              </span>
            </button>
          ))}
        </div>
      )}

      {!isLoading && hasSearched && searchResults.length === 0 && (
        <p className={styles.noResults}>
          {t("CoverDesigner.No images found")}
        </p>
      )}

      {!hasSearched && (
        <p className={styles.searchHint}>
          {t("CoverDesigner.Search for images to use as background")}
        </p>
      )}
    </div>
  );
}
