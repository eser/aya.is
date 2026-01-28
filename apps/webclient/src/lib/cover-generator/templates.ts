// Template Configurations for Cover Generator

import type { TemplateId, CoverOptions } from "./types.ts";

export interface TemplateConfig {
  id: TemplateId;
  name: string;
  description: string;
  // Default options specific to this template
  defaults: Partial<CoverOptions>;
}

export const templates: Record<TemplateId, TemplateConfig> = {
  classic: {
    id: "classic",
    name: "Classic",
    description: "Clean design with centered title and author info at bottom",
    defaults: {
      backgroundColor: "#1a1a2e",
      accentColor: "#e94560",
      textColor: "#ffffff",
      themePreset: "dark",
      titleSize: 64,
      subtitleSize: 22,
      showAuthor: true,
      showDate: true,
      showStoryKind: false,
      showLogo: true,
      logoPosition: "top-right",
      backgroundPattern: "none",
      padding: 60,
    },
  },
  bold: {
    id: "bold",
    name: "Bold",
    description: "Attention-grabbing design with accent highlights and badges",
    defaults: {
      backgroundColor: "#0f0f23",
      accentColor: "#ff6b6b",
      textColor: "#ffffff",
      themePreset: "brand",
      titleSize: 52,
      subtitleSize: 20,
      showAuthor: true,
      showDate: false,
      showStoryKind: true,
      showLogo: true,
      logoPosition: "top-right",
      backgroundPattern: "diagonal",
      padding: 50,
    },
  },
  minimal: {
    id: "minimal",
    name: "Minimal",
    description: "Simple, elegant design with focus on typography",
    defaults: {
      backgroundColor: "#ffffff",
      accentColor: "#1a1a2e",
      textColor: "#1a1a2e",
      themePreset: "light",
      titleSize: 64,
      subtitleSize: 22,
      showAuthor: false,
      showDate: false,
      showStoryKind: false,
      showLogo: true,
      logoPosition: "top-right",
      backgroundPattern: "none",
      padding: 80,
    },
  },
  featured: {
    id: "featured",
    name: "Featured",
    description: "Premium look with prominent author photo and featured badge",
    defaults: {
      backgroundColor: "#16213e",
      accentColor: "#f9a825",
      textColor: "#ffffff",
      themePreset: "dark",
      titleSize: 48,
      subtitleSize: 20,
      showAuthor: true,
      showDate: true,
      showStoryKind: true,
      showLogo: true,
      logoPosition: "top-right",
      backgroundPattern: "dots",
      padding: 60,
    },
  },
};

// Get template list for UI
export function getTemplateList(): TemplateConfig[] {
  return Object.values(templates);
}

// Get template by ID
export function getTemplate(id: TemplateId): TemplateConfig {
  return templates[id];
}

// Merge template defaults with user options
export function mergeWithTemplateDefaults(
  templateId: TemplateId,
  options: Partial<CoverOptions>,
): CoverOptions {
  const template = templates[templateId];
  const { defaultCoverOptions } = require("./types.ts");

  return {
    ...defaultCoverOptions,
    ...template.defaults,
    ...options,
    templateId,
  };
}
