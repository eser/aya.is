import type { Profile } from "@/shared/modules/backend/profiles/types.ts";

export type StoryKind = "status" | "announcement" | "article" | "content" | "presentation";

export type Publication = Profile;

export type Story = {
  id: string;
  kind: StoryKind;
  slug: string | null;

  story_picture_uri: string | null;
  title: string | null;
  summary: string | null;
  content: string;

  is_featured: boolean;
  author_profile_id: string | null;
  author_profile: Profile | null;

  created_at: string;
  updated_at: string | null;
  deleted_at: string | null;
};

export type StoryEx = Omit<Story, "author_profile_id"> & {
  author_profile: {
    id: string;
    slug: string;
    title: string;
    profile_picture_uri: string | null;
  };
  publications: Publication[];
};
