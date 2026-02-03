// Backend API Facade
// This module provides a unified interface to all backend API calls

import { getProfilesByKinds } from "./profiles/get-profiles-by-kinds";
import { getProfile } from "./profiles/get-profile";
import { getProfilePage } from "./profiles/get-profile-page";
import { getProfilePermissions } from "./profiles/get-profile-permissions";
import { getProfileStories } from "./profiles/get-profile-stories";
import { getProfileStory } from "./profiles/get-profile-story";
import { getProfileMembers } from "./profiles/get-profile-members";
import { getProfileContributions } from "./profiles/get-profile-contributions";
import { checkProfileSlug } from "./profiles/check-profile-slug";
import { checkPageSlug } from "./profiles/check-page-slug";
import { createProfile } from "./profiles/create-profile";
import { updateProfile } from "./profiles/update-profile";
import { updateProfileTranslation } from "./profiles/update-profile-translation";
import { uploadProfilePicture } from "./profiles/upload-profile-picture";
import { updateProfilePicture } from "./profiles/update-profile-picture";
import { getProfileTranslations } from "./profiles/get-profile-translations";
import { listProfileLinks } from "./profiles/list-profile-links";
import { getProfileLinks } from "./profiles/get-profile-links";
import { createProfileLink } from "./profiles/create-profile-link";
import { updateProfileLink } from "./profiles/update-profile-link";
import { deleteProfileLink } from "./profiles/delete-profile-link";
import { initiateProfileLinkOAuth } from "./profiles/initiate-profile-link-oauth";
import { getGitHubAccounts } from "./profiles/get-github-accounts";
import { finalizeGitHubConnection } from "./profiles/finalize-github-connection";
import { listProfilePages } from "./profiles/list-profile-pages";
import { listProfileMemberships } from "./profiles/list-profile-memberships";
import { searchUsersForMembership } from "./profiles/search-users-for-membership";
import { addProfileMembership } from "./profiles/add-profile-membership";
import { updateProfileMembership } from "./profiles/update-profile-membership";
import { deleteProfileMembership } from "./profiles/delete-profile-membership";
import { createProfilePage } from "./profiles/create-profile-page";
import { updateProfilePage } from "./profiles/update-profile-page";
import { updateProfilePageTranslation } from "./profiles/update-profile-page-translation";
import { deleteProfilePage } from "./profiles/delete-profile-page";
import { listProfilePointTransactions } from "./profile-points/list-profile-point-transactions";
import { getStoriesByKinds } from "./stories/get-stories-by-kinds";
import { getStory } from "./stories/get-story";
import { checkStorySlug } from "./stories/check-story-slug";
import { insertStory } from "./stories/insert-story";
import { getStoryForEdit } from "./stories/get-story-for-edit";
import { updateStory } from "./stories/update-story";
import { updateStoryTranslation } from "./stories/update-story-translation";
import { removeStory } from "./stories/remove-story";
import { getStoryPermissions } from "./stories/get-story-permissions";
import { listStoryPublications } from "./stories/list-story-publications";
import { addStoryPublication } from "./stories/add-story-publication";
import { updateStoryPublication } from "./stories/update-story-publication";
import { removeStoryPublication } from "./stories/remove-story-publication";
import { getPresignedURL } from "./uploads/get-presigned-url";
import { uploadToPresignedURL } from "./uploads/upload-to-presigned-url";
import { removeUpload } from "./uploads/remove-upload";
import { getCurrentSession } from "./sessions/get-current-session";
import { createSession } from "./sessions/create-session";
import { getSessionCurrent } from "./sessions/get-session-current";
import { updateSessionPreferences } from "./sessions/update-preferences";
import { getPOWChallenge } from "./protection/get-pow-challenge";
import { isPOWSolverSupported, solvePOW } from "./protection/solve-pow";
import { getSpotlight } from "./site/get-spotlight";
import { searchBackgroundImages } from "./site/search-background-images";
import { handleAuthCallback } from "./auth/handle-callback";
import { getUser } from "./users/get-user";
import { getUsers } from "./users/get-users";
import { search } from "./search/search";
import { getPendingAwards } from "./admin/get-pending-awards";
import { getPendingAwardsStats } from "./admin/get-pending-awards-stats";
import { approvePendingAward } from "./admin/approve-pending-award";
import { rejectPendingAward } from "./admin/reject-pending-award";
import { bulkApprovePendingAwards } from "./admin/bulk-approve-pending-awards";
import { bulkRejectPendingAwards } from "./admin/bulk-reject-pending-awards";
import { getAdminProfiles } from "./admin/get-admin-profiles";
import { getAdminProfile } from "./admin/get-admin-profile";
import { addAdminPoints } from "./admin/add-admin-points";

// Re-export types
export * from "./types";
export type { GitHubAccount, GitHubAccountsResponse } from "./profiles/get-github-accounts";
export type {
  UnsplashPhoto,
  UnsplashPhotoUrls,
  UnsplashPhotoUser,
  UnsplashSearchResult,
  SearchBackgroundImagesParams,
} from "./site/search-background-images";

// Backend API facade
export const backend = {
  // Profiles
  getProfilesByKinds,
  getProfile,
  getProfilePage,
  getProfilePermissions,
  getProfileStories,
  getProfileStory,
  getProfileMembers,
  getProfileContributions,
  checkProfileSlug,
  checkPageSlug,
  createProfile,
  updateProfile,
  updateProfileTranslation,
  uploadProfilePicture,
  updateProfilePicture,
  getProfileTranslations,

  // Profile Links
  listProfileLinks,
  getProfileLinks,
  createProfileLink,
  updateProfileLink,
  deleteProfileLink,
  initiateProfileLinkOAuth,
  getGitHubAccounts,
  finalizeGitHubConnection,

  // Profile Memberships
  listProfileMemberships,
  searchUsersForMembership,
  addProfileMembership,
  updateProfileMembership,
  deleteProfileMembership,

  // Profile Pages
  listProfilePages,
  createProfilePage,
  updateProfilePage,
  updateProfilePageTranslation,
  deleteProfilePage,

  // Profile Points
  listProfilePointTransactions,

  // Stories
  getStoriesByKinds,
  getStory,
  checkStorySlug,
  insertStory,
  getStoryForEdit,
  updateStory,
  updateStoryTranslation,
  removeStory,
  getStoryPermissions,
  listStoryPublications,
  addStoryPublication,
  updateStoryPublication,
  removeStoryPublication,

  // Uploads
  getPresignedURL,
  uploadToPresignedURL,
  removeUpload,

  // Sessions
  getCurrentSession,
  createSession,
  getSessionCurrent,
  updateSessionPreferences,

  // Protection
  getPOWChallenge,
  solvePOW,
  isPOWSolverSupported,

  // Site
  getSpotlight,
  searchBackgroundImages,

  // Auth
  handleAuthCallback,

  // Users
  getUser,
  getUsers,

  // Search
  search,

  // Admin
  addAdminPoints,
  getAdminProfile,
  getAdminProfiles,
  getPendingAwards,
  getPendingAwardsStats,
  approvePendingAward,
  rejectPendingAward,
  bulkApprovePendingAwards,
  bulkRejectPendingAwards,
};

// Individual exports for tree-shaking
export {
  checkPageSlug,
  checkProfileSlug,
  checkStorySlug,
  createProfile,
  createProfileLink,
  createProfilePage,
  createSession,
  deleteProfileLink,
  deleteProfilePage,
  finalizeGitHubConnection,
  getCurrentSession,
  getGitHubAccounts,
  getPOWChallenge,
  getPresignedURL,
  getProfile,
  getProfileContributions,
  getProfileMembers,
  getProfilePage,
  getProfilePermissions,
  getProfilesByKinds,
  getProfileStories,
  getProfileStory,
  getProfileTranslations,
  getSessionCurrent,
  getSpotlight,
  getStoriesByKinds,
  getStory,
  getStoryForEdit,
  getStoryPermissions,
  getUser,
  getUsers,
  handleAuthCallback,
  insertStory,
  isPOWSolverSupported,
  getProfileLinks,
  listProfileLinks,
  listProfileMemberships,
  listProfilePages,
  listProfilePointTransactions,
  addProfileMembership,
  deleteProfileMembership,
  searchUsersForMembership,
  updateProfileMembership,
  removeStory,
  removeUpload,
  search,
  searchBackgroundImages,
  solvePOW,
  updateProfile,
  updateProfileLink,
  updateProfilePage,
  updateProfilePageTranslation,
  updateProfilePicture,
  updateProfileTranslation,
  updateSessionPreferences,
  updateStory,
  updateStoryTranslation,
  uploadProfilePicture,
  uploadToPresignedURL,

  // Admin
  approvePendingAward,
  addAdminPoints,
  bulkApprovePendingAwards,
  bulkRejectPendingAwards,
  getAdminProfile,
  getAdminProfiles,
  getPendingAwards,
  getPendingAwardsStats,
  rejectPendingAward,
};
