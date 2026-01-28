// Cover Generator Type Definitions

export type TemplateId = "classic" | "bold" | "minimal" | "featured";

export type ThemePreset = "light" | "dark" | "brand" | "custom";

export type FontFamily = "bree-serif" | "nunito-sans" | "system";

export type LogoPosition =
  | "top-left"
  | "top-right"
  | "bottom-left"
  | "bottom-right";

export type BackgroundPattern = "none" | "dots" | "grid" | "diagonal";

// Story data passed to the generator
export interface StoryData {
  title: string;
  summary: string | null;
  kind: string;
  kindLabel: string; // Localized label for the kind
  authorName: string | null;
  authorAvatarUrl: string | null;
  publishedAt: string | null;
}

// Customization options for the cover
export interface CoverOptions {
  // Locale for date formatting
  locale: string;

  // Template
  templateId: TemplateId;

  // Colors
  backgroundColor: string;
  accentColor: string;
  textColor: string;
  themePreset: ThemePreset;

  // Typography
  headingFont: FontFamily;
  bodyFont: FontFamily;
  titleSize: number; // font size in pixels (template-specific default)
  subtitleSize: number; // font size in pixels (template-specific default)
  lineHeight: number; // multiplier for line height (like CSS line-height), 100 = 1.0, 120 = 1.2
  lineSpacing: number; // margin between title and subtitle in pixels

  // Content
  titleOverride: string;
  subtitleOverride: string;
  showAuthor: boolean;
  showDate: boolean;
  showStoryKind: boolean;

  // Branding
  showLogo: boolean;
  logoPosition: LogoPosition;
  logoOpacity: number; // 0-100

  // Layout
  padding: number; // pixels
  contentOffsetY: number; // vertical offset for content position in pixels
  backgroundPattern: BackgroundPattern;

  // Background Image (from Unsplash)
  backgroundImageUrl: string | null; // URL of the background image
  backgroundImageOpacity: number; // 0-100
}

// Default options
export const defaultCoverOptions: CoverOptions = {
  locale: "en",
  templateId: "classic",
  backgroundColor: "#1a1a2e",
  accentColor: "#e94560",
  textColor: "#ffffff",
  themePreset: "dark",
  headingFont: "bree-serif",
  bodyFont: "nunito-sans",
  titleSize: 56,
  subtitleSize: 20,
  lineHeight: 120,
  lineSpacing: 24,
  titleOverride: "",
  subtitleOverride: "",
  showAuthor: true,
  showDate: true,
  showStoryKind: true,
  showLogo: true,
  logoPosition: "top-right",
  logoOpacity: 80,
  padding: 60,
  contentOffsetY: 0,
  backgroundPattern: "none",
  backgroundImageUrl: null,
  backgroundImageOpacity: 50,
};

// Canvas dimensions for the cover (16:9 aspect ratio)
export const COVER_WIDTH = 1200;
export const COVER_HEIGHT = 675;

// Theme presets
export const themePresets: Record<
  ThemePreset,
  { backgroundColor: string; accentColor: string; textColor: string }
> = {
  light: {
    backgroundColor: "#ffffff",
    accentColor: "#e94560",
    textColor: "#1a1a2e",
  },
  dark: {
    backgroundColor: "#1a1a2e",
    accentColor: "#e94560",
    textColor: "#ffffff",
  },
  brand: {
    backgroundColor: "#0f0f23",
    accentColor: "#00d4aa",
    textColor: "#ffffff",
  },
  custom: {
    backgroundColor: "#1a1a2e",
    accentColor: "#e94560",
    textColor: "#ffffff",
  },
};

// Font mappings for canvas
export const fontFamilyMap: Record<FontFamily, string> = {
  "bree-serif": '"Bree Serif", Georgia, serif',
  "nunito-sans": '"Nunito Sans Variable", system-ui, sans-serif',
  system: "system-ui, -apple-system, sans-serif",
};
