import { getCustomDomain } from "./profiles/get-custom-domain.ts";
import { getProfilesByKinds } from "./profiles/get-profiles-by-kinds.ts";
import { getProfile } from "./profiles/get-profile.ts";
import { getProfilePage } from "./profiles/get-profile-page.ts";
import { getProfileContributions } from "./profiles/get-profile-contributions.ts";
import { getProfileMembers } from "./profiles/get-profile-members.ts";
import { getProfileStories } from "./profiles/get-profile-stories.ts";
import { getProfileStory } from "./profiles/get-profile-story.ts";
import { checkProfileSlug } from "./profiles/check-profile-slug.ts";
import { createProfile } from "./profiles/create-profile.ts";
import { updateProfile } from "./profiles/update-profile.ts";
import { updateProfileTranslation } from "./profiles/update-profile-translation.ts";
import { getProfilePermissions } from "./profiles/get-profile-permissions.ts";
import { uploadProfilePicture } from "./profiles/upload-profile-picture.ts";
import { listProfileLinks } from "./profiles/list-profile-links.ts";
import { createProfileLink } from "./profiles/create-profile-link.ts";
import { updateProfileLink } from "./profiles/update-profile-link.ts";
import { deleteProfileLink } from "./profiles/delete-profile-link.ts";
import { listProfilePages } from "./profiles/list-profile-pages.ts";
import { createProfilePage } from "./profiles/create-profile-page.ts";
import { updateProfilePage } from "./profiles/update-profile-page.ts";
import { updateProfilePageTranslation } from "./profiles/update-profile-page-translation.ts";
import { deleteProfilePage } from "./profiles/delete-profile-page.ts";
import { getProfileTranslations } from "./profiles/get-profile-translations.ts";
import { getStory } from "./stories/get-story.ts";
import { getStoriesByKinds } from "./stories/get-stories-by-kinds.ts";
import { getSpotlight } from "./site/get-spotlight.ts";
import { handleAuthCallback } from "./auth/handle-callback.ts";
import { getCurrentSession } from "./sessions/get-current-session.ts";

export const backend = {
  // Profiles
  getCustomDomain,
  getProfilesByKinds,
  getProfile,
  getProfilePage,
  getProfileContributions,
  getProfileMembers,
  getProfileStories,
  getProfileStory,
  checkProfileSlug,
  createProfile,
  updateProfile,
  updateProfileTranslation,
  getProfilePermissions,
  uploadProfilePicture,
  listProfileLinks,
  createProfileLink,
  updateProfileLink,
  deleteProfileLink,
  listProfilePages,
  createProfilePage,
  updateProfilePage,
  updateProfilePageTranslation,
  deleteProfilePage,
  getProfileTranslations,

  // Stories
  getStory,
  getStoriesByKinds,

  // Site
  getSpotlight,

  // Auth
  handleAuthCallback,

  // Sessions
  getCurrentSession,
};
