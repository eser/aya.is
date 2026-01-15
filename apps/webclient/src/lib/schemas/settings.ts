import { z } from "zod";

// Profile link schema
export const profileLinkSchema = z.object({
  platform: z.string().min(1, "Platform is required"),
  url: z.string().url("Must be a valid URL"),
  label: z.string().max(100).optional(),
});

// Update profile link schema
export const updateProfileLinkSchema = z.object({
  platform: z.string().min(1).optional(),
  url: z.string().url("Must be a valid URL").optional(),
  label: z.string().max(100).optional(),
  sort_order: z.number().int().min(0).optional(),
});

// Profile page schema
export const profilePageSchema = z.object({
  slug: z
    .string()
    .min(1, "Slug is required")
    .max(50, "Slug is too long")
    .regex(
      /^[a-z0-9-]+$/,
      "Slug can only contain lowercase letters, numbers, and hyphens",
    ),
  title: z.string().min(1, "Title is required").max(100, "Title is too long"),
  content: z.string().min(1, "Content is required"),
});

// Update profile page schema
export const updateProfilePageSchema = z.object({
  title: z.string().min(1).max(100).optional(),
  content: z.string().optional(),
  sort_order: z.number().int().min(0).optional(),
});

// Profile translation schema
export const profileTranslationSchema = z.object({
  locale: z.string().min(2, "Locale is required").max(10),
  title: z.string().min(1, "Title is required").max(100),
  description: z.string().max(500).optional(),
});

// Page translation schema
export const pageTranslationSchema = z.object({
  locale: z.string().min(2, "Locale is required").max(10),
  title: z.string().min(1, "Title is required").max(100),
  content: z.string().min(1, "Content is required"),
});

// Type exports
export type ProfileLinkInput = z.infer<typeof profileLinkSchema>;
export type UpdateProfileLinkInput = z.infer<typeof updateProfileLinkSchema>;
export type ProfilePageInput = z.infer<typeof profilePageSchema>;
export type UpdateProfilePageInput = z.infer<typeof updateProfilePageSchema>;
export type ProfileTranslationInput = z.infer<typeof profileTranslationSchema>;
export type PageTranslationInput = z.infer<typeof pageTranslationSchema>;
