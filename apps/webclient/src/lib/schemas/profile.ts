import { z } from "zod";

// Profile creation schema
export const createProfileSchema = z.object({
  slug: z
    .string()
    .min(3, "Slug must be at least 3 characters")
    .max(50, "Slug must be at most 50 characters")
    .regex(
      /^[a-z0-9-]+$/,
      "Slug can only contain lowercase letters, numbers, and hyphens",
    ),
  title: z.string().min(1, "Title is required").max(100, "Title is too long"),
  description: z.string().max(500, "Description is too long").optional(),
  kind: z.enum(["individual", "organization", "product"], {
    required_error: "Please select a profile type",
  }),
});

// Profile update schema (all fields optional)
export const updateProfileSchema = z.object({
  title: z.string().min(1, "Title is required").max(100).optional(),
  description: z.string().max(500).optional(),
  pronouns: z.string().max(50).optional(),
});

// Profile slug check schema
export const checkSlugSchema = z.object({
  slug: z
    .string()
    .min(3, "Slug must be at least 3 characters")
    .max(50, "Slug must be at most 50 characters")
    .regex(
      /^[a-z0-9-]+$/,
      "Slug can only contain lowercase letters, numbers, and hyphens",
    ),
});

// Type exports
export type CreateProfileInput = z.infer<typeof createProfileSchema>;
export type UpdateProfileInput = z.infer<typeof updateProfileSchema>;
export type CheckSlugInput = z.infer<typeof checkSlugSchema>;
