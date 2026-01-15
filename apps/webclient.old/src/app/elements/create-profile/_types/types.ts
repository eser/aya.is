export type CheckSlugResponse = {
  exists: boolean;
};

export type CreateProfileRequest = {
  kind: string;
  slug: string;
  title: string;
  description: string;
};

export type CreateProfileResponse = {
  id: string;
  slug: string;
  kind: string;
  title: string;
  description: string;
  created_at: string;
};
