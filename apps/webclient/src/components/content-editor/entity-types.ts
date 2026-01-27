// Entity type definitions for ContentEditor multi-entity architecture

import { allowedURIPrefixes } from "@/config";

// Entity types supported by ContentEditor
export type ContentEntityType = "story" | "page";

// Image field configuration
export interface ImageFieldConfig {
  key: string;           // Internal state key (e.g., "storyPictureUri")
  apiField: string;      // Backend API field name (e.g., "story_picture_uri")
  labelKey: string;      // Translation key (e.g., "Editor.Story Picture URI")
  required: boolean;
  allowedPrefixes: string[];  // From config.ts
}

// Entity-specific configuration
export interface EntityConfig {
  type: ContentEntityType;

  // Image fields (cover image, OG image, etc.)
  imageFields: ImageFieldConfig[];

  // Entity-specific features
  features: {
    hasKind: boolean;           // Stories have kind (article, announcement, etc.)
    hasFeatured: boolean;       // Stories can be featured
    slugDatePrefix: boolean;    // Stories require YYYYMMDD prefix
  };
}

// Story configuration
export const STORY_CONFIG: EntityConfig = {
  type: "story",
  imageFields: [
    {
      key: "storyPictureUri",
      apiField: "story_picture_uri",
      labelKey: "Editor.Story Picture URI",
      required: false,
      allowedPrefixes: allowedURIPrefixes.stories,
    },
  ],
  features: {
    hasKind: true,
    hasFeatured: true,
    slugDatePrefix: true,
  },
};

// Page configuration
export const PAGE_CONFIG: EntityConfig = {
  type: "page",
  imageFields: [
    {
      key: "storyPictureUri",  // Same internal key for backwards compatibility
      apiField: "cover_picture_uri",
      labelKey: "ContentEditor.Cover Picture URI",
      required: false,
      allowedPrefixes: allowedURIPrefixes.profiles,
    },
  ],
  features: {
    hasKind: false,
    hasFeatured: false,
    slugDatePrefix: false,
  },
};

// Get entity config by type
export function getEntityConfig(type: ContentEntityType): EntityConfig {
  switch (type) {
    case "story":
      return STORY_CONFIG;
    case "page":
      return PAGE_CONFIG;
    default:
      throw new Error(`Unknown entity type: ${type}`);
  }
}
