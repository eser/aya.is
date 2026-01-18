// Backend API Facade
// This module provides a unified interface to all backend API calls

import { getCustomDomain } from "./profiles/get-custom-domain";
import { getProfilesByKinds } from "./profiles/get-profiles-by-kinds";
import { getProfile } from "./profiles/get-profile";
import { getProfilePage } from "./profiles/get-profile-page";
import { getProfilePermissions } from "./profiles/get-profile-permissions";
import { getProfileStories } from "./profiles/get-profile-stories";
import { getProfileStory } from "./profiles/get-profile-story";
import { getProfileMembers } from "./profiles/get-profile-members";
import { getProfileContributions } from "./profiles/get-profile-contributions";
import { checkProfileSlug } from "./profiles/check-profile-slug";
import { createProfile } from "./profiles/create-profile";
import { updateProfile } from "./profiles/update-profile";
import { updateProfileTranslation } from "./profiles/update-profile-translation";
import { uploadProfilePicture } from "./profiles/upload-profile-picture";
import { getProfileTranslations } from "./profiles/get-profile-translations";
import { listProfileLinks } from "./profiles/list-profile-links";
import { createProfileLink } from "./profiles/create-profile-link";
import { updateProfileLink } from "./profiles/update-profile-link";
import { deleteProfileLink } from "./profiles/delete-profile-link";
import { listProfilePages } from "./profiles/list-profile-pages";
import { createProfilePage } from "./profiles/create-profile-page";
import { updateProfilePage } from "./profiles/update-profile-page";
import { updateProfilePageTranslation } from "./profiles/update-profile-page-translation";
import { deleteProfilePage } from "./profiles/delete-profile-page";
import { getStoriesByKinds } from "./stories/get-stories-by-kinds";
import { getStory } from "./stories/get-story";
import { getCurrentSession } from "./sessions/get-current-session";
import { createSession } from "./sessions/create-session";
import { getSessionPreferences } from "./sessions/get-preferences";
import { updateSessionPreferences } from "./sessions/update-preferences";
import { getPOWChallenge } from "./protection/get-pow-challenge";
import { solvePOW, isPOWSolverSupported } from "./protection/solve-pow";
import { getSpotlight } from "./site/get-spotlight";
import { handleAuthCallback } from "./auth/handle-callback";
import { checkSessionViaCookie } from "./auth/session-check";
import { getUser } from "./users/get-user";
import { getUsers } from "./users/get-users";

// Re-export types
export * from "./types";

// Backend API facade
export const backend = {
  // Profiles
  getCustomDomain,
  getProfilesByKinds,
  getProfile,
  getProfilePage,
  getProfilePermissions,
  getProfileStories,
  getProfileStory,
  getProfileMembers,
  getProfileContributions,
  checkProfileSlug,
  createProfile,
  updateProfile,
  updateProfileTranslation,
  uploadProfilePicture,
  getProfileTranslations,

  // Profile Links
  listProfileLinks,
  createProfileLink,
  updateProfileLink,
  deleteProfileLink,

  // Profile Pages
  listProfilePages,
  createProfilePage,
  updateProfilePage,
  updateProfilePageTranslation,
  deleteProfilePage,

  // Stories
  getStoriesByKinds,
  getStory,

  // Sessions
  getCurrentSession,
  createSession,
  getSessionPreferences,
  updateSessionPreferences,

  // Protection
  getPOWChallenge,
  solvePOW,
  isPOWSolverSupported,

  // Site
  getSpotlight,

  // Auth
  handleAuthCallback,
  checkSessionViaCookie,

  // Users
  getUser,
  getUsers,
};

// Individual exports for tree-shaking
export {
  checkProfileSlug,
  checkSessionViaCookie,
  createProfile,
  createProfileLink,
  createProfilePage,
  createSession,
  deleteProfileLink,
  deleteProfilePage,
  getCurrentSession,
  getCustomDomain,
  getPOWChallenge,
  getProfile,
  getProfileContributions,
  getProfileMembers,
  getProfilePage,
  getProfilePermissions,
  getProfilesByKinds,
  getProfileStories,
  getProfileStory,
  getProfileTranslations,
  getSessionPreferences,
  getSpotlight,
  getStoriesByKinds,
  getStory,
  getUser,
  getUsers,
  handleAuthCallback,
  isPOWSolverSupported,
  listProfileLinks,
  listProfilePages,
  solvePOW,
  updateProfile,
  updateProfileLink,
  updateProfilePage,
  updateProfilePageTranslation,
  updateProfileTranslation,
  updateSessionPreferences,
  uploadProfilePicture,
};
