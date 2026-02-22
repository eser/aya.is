// Backend API Facade
// This module provides a unified interface to all backend API calls

import { getProfilesByKinds } from "./profiles/get-profiles-by-kinds";
import { getProfile } from "./profiles/get-profile";
import { getProfilePage } from "./profiles/get-profile-page";
import { getProfilePermissions } from "./profiles/get-profile-permissions";
import { getProfileStories } from "./profiles/get-profile-stories";
import { getProfileAuthoredStories } from "./profiles/get-profile-authored-stories";
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
import { getLinkedInAccounts } from "./profiles/get-linkedin-accounts";
import { finalizeLinkedInConnection } from "./profiles/finalize-linkedin-connection";
import { connectSpeakerDeck } from "./profiles/connect-speakerdeck";
import { verifyTelegramCode } from "./profiles/verify-telegram-code";
import { listProfilePages } from "./profiles/list-profile-pages";
import { listProfileMemberships } from "./profiles/list-profile-memberships";
import { searchUsersForMembership } from "./profiles/search-users-for-membership";
import { addProfileMembership } from "./profiles/add-profile-membership";
import { updateProfileMembership } from "./profiles/update-profile-membership";
import { deleteProfileMembership } from "./profiles/delete-profile-membership";
import { followProfile } from "./profiles/follow-profile";
import { unfollowProfile } from "./profiles/unfollow-profile";
import { createProfilePage } from "./profiles/create-profile-page";
import { updateProfilePage } from "./profiles/update-profile-page";
import { updateProfilePageTranslation } from "./profiles/update-profile-page-translation";
import { deleteProfilePage } from "./profiles/delete-profile-page";
import { listProfilePageTranslationLocales } from "./profiles/list-profile-page-translation-locales";
import { deleteProfilePageTranslation } from "./profiles/delete-profile-page-translation";
import { autoTranslateProfilePage } from "./profiles/auto-translate-profile-page";
import { generateCVPage } from "./profiles/generate-cv-page";
import { listProfileResources } from "./profiles/list-profile-resources";
import { createProfileResource } from "./profiles/create-profile-resource";
import { deleteProfileResource } from "./profiles/delete-profile-resource";
import { listGitHubRepos } from "./profiles/list-github-repos";
import { listProfilePointTransactions } from "./profile-points/list-profile-point-transactions";
import { listProfileEnvelopes } from "./profile-envelopes/list-profile-envelopes";
import { acceptProfileEnvelope } from "./profile-envelopes/accept-profile-envelope";
import { rejectProfileEnvelope } from "./profile-envelopes/reject-profile-envelope";
import { sendProfileEnvelope } from "./profile-envelopes/send-profile-envelope";
import { listMailboxEnvelopes } from "./mailbox/list-mailbox-envelopes";
import { listConversations, type ListConversationsResult } from "./mailbox/list-conversations";
import { getConversation } from "./mailbox/get-conversation";
import { markConversationRead } from "./mailbox/mark-conversation-read";
import { archiveConversation } from "./mailbox/archive-conversation";
import { unarchiveConversation } from "./mailbox/unarchive-conversation";
import { sendMailboxMessage } from "./mailbox/send-mailbox-message";
import { acceptMailboxMessage } from "./mailbox/accept-mailbox-message";
import { rejectMailboxMessage } from "./mailbox/reject-mailbox-message";
import { addReaction } from "./mailbox/add-reaction";
import { removeReaction } from "./mailbox/remove-reaction";
import { removeConversation } from "./mailbox/remove-conversation";
import { getUnreadCount } from "./mailbox/get-unread-count";
import { getProfileQuestions } from "./questions/get-profile-questions";
import { createQuestion } from "./questions/create-question";
import { voteQuestion } from "./questions/vote-question";
import { answerQuestion } from "./questions/answer-question";
import { editAnswer } from "./questions/edit-answer";
import { hideQuestion } from "./questions/hide-question";
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
import { listStoryTranslationLocales } from "./stories/list-story-translation-locales";
import { deleteStoryTranslation } from "./stories/delete-story-translation";
import { autoTranslateStory } from "./stories/auto-translate-story";
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
import { getActivities } from "./activities/get-activities";
import { getActivity } from "./activities/get-activity";
import { setInteraction } from "./interactions/set-interaction";
import { removeInteraction } from "./interactions/remove-interaction";
import { getInteractions } from "./interactions/get-interactions";
import { getMyInteractions } from "./interactions/get-my-interactions";
import { getInteractionCounts } from "./interactions/get-interaction-counts";
import { getPendingAwards } from "./admin/get-pending-awards";
import { getPendingAwardsStats } from "./admin/get-pending-awards-stats";
import { approvePendingAward } from "./admin/approve-pending-award";
import { rejectPendingAward } from "./admin/reject-pending-award";
import { bulkApprovePendingAwards } from "./admin/bulk-approve-pending-awards";
import { bulkRejectPendingAwards } from "./admin/bulk-reject-pending-awards";
import { getAdminProfiles } from "./admin/get-admin-profiles";
import { getAdminProfile } from "./admin/get-admin-profile";
import { addAdminPoints } from "./admin/add-admin-points";
import { getAdminWorkers } from "./admin/get-admin-workers";
import { toggleAdminWorker } from "./admin/toggle-admin-worker";
import { triggerAdminWorker } from "./admin/trigger-admin-worker";

// Re-export types
export * from "./types";
export type { AdminWorkerStatus } from "./admin/get-admin-workers";
export type { AccessibleProfile } from "./sessions/types";
export type { ToggleWorkerResult } from "./admin/toggle-admin-worker";
export type { TriggerWorkerResult } from "./admin/trigger-admin-worker";
export type { GitHubAccount, GitHubAccountsResponse } from "./profiles/get-github-accounts";
export type { VerifyTelegramCodeResponse } from "./profiles/verify-telegram-code";
export type { SendProfileEnvelopeParams } from "./profile-envelopes/send-profile-envelope";
export type { MailboxEnvelope } from "./mailbox/list-mailbox-envelopes";
export type { SendMailboxMessageParams } from "./mailbox/send-mailbox-message";
export type { ListConversationsResult } from "./mailbox/list-conversations";
export type { CreateProfileResourceInput } from "./profiles/create-profile-resource";
export type { CreateQuestionInput } from "./questions/create-question";
export type { VoteQuestionResult } from "./questions/vote-question";
export type { AnswerQuestionInput } from "./questions/answer-question";
export type { EditAnswerInput } from "./questions/edit-answer";
export type { HideQuestionInput } from "./questions/hide-question";
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
  getProfileAuthoredStories,
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
  getLinkedInAccounts,
  finalizeLinkedInConnection,
  connectSpeakerDeck,
  verifyTelegramCode,

  // Profile Resources
  listProfileResources,
  createProfileResource,
  deleteProfileResource,
  listGitHubRepos,

  // Profile Memberships
  listProfileMemberships,
  searchUsersForMembership,
  addProfileMembership,
  updateProfileMembership,
  deleteProfileMembership,
  followProfile,
  unfollowProfile,

  // Profile Pages
  listProfilePages,
  createProfilePage,
  updateProfilePage,
  updateProfilePageTranslation,
  deleteProfilePage,
  listProfilePageTranslationLocales,
  deleteProfilePageTranslation,
  autoTranslateProfilePage,
  generateCVPage,

  // Profile Points
  listProfilePointTransactions,

  // Profile Envelopes (Inbox)
  listProfileEnvelopes,
  acceptProfileEnvelope,
  rejectProfileEnvelope,
  sendProfileEnvelope,

  // Mailbox (conversations)
  listMailboxEnvelopes,
  listConversations,
  getConversation,
  markConversationRead,
  archiveConversation,
  unarchiveConversation,
  removeConversation,
  sendMailboxMessage,
  acceptMailboxMessage,
  rejectMailboxMessage,
  addReaction,
  removeReaction,
  getUnreadCount,

  // Profile Questions (Q&A)
  getProfileQuestions,
  createQuestion,
  voteQuestion,
  answerQuestion,
  editAnswer,
  hideQuestion,

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
  listStoryTranslationLocales,
  deleteStoryTranslation,
  autoTranslateStory,

  // Activities
  getActivities,
  getActivity,

  // Story Interactions
  setInteraction,
  removeInteraction,
  getInteractions,
  getMyInteractions,
  getInteractionCounts,

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
  getAdminWorkers,
  toggleAdminWorker,
  triggerAdminWorker,
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
  connectSpeakerDeck,
  createProfile,
  createProfileLink,
  createProfilePage,
  createSession,
  deleteProfileLink,
  deleteProfilePage,
  finalizeGitHubConnection,
  finalizeLinkedInConnection,
  verifyTelegramCode,
  getCurrentSession,
  getGitHubAccounts,
  getLinkedInAccounts,
  getPOWChallenge,
  getPresignedURL,
  getProfile,
  getProfileContributions,
  getProfileMembers,
  getProfilePage,
  getProfilePermissions,
  getProfilesByKinds,
  getProfileAuthoredStories,
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
  listProfileResources,
  createProfileResource,
  deleteProfileResource,
  listGitHubRepos,
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

  // Translation management
  listStoryTranslationLocales,
  deleteStoryTranslation,
  autoTranslateStory,
  listProfilePageTranslationLocales,
  deleteProfilePageTranslation,
  autoTranslateProfilePage,
  generateCVPage,

  // Profile Questions (Q&A)
  getProfileQuestions,
  createQuestion,
  voteQuestion,
  answerQuestion,
  editAnswer,
  hideQuestion,

  // Activities
  getActivities,
  getActivity,

  // Story Interactions
  setInteraction,
  removeInteraction,
  getInteractions,
  getMyInteractions,
  getInteractionCounts,

  // Profile Envelopes (Inbox)
  listProfileEnvelopes,
  acceptProfileEnvelope,
  rejectProfileEnvelope,
  sendProfileEnvelope,

  // Mailbox (conversations)
  listConversations,
  getConversation,
  markConversationRead,
  archiveConversation,
  unarchiveConversation,
  removeConversation,
  sendMailboxMessage,
  acceptMailboxMessage,
  rejectMailboxMessage,
  addReaction,
  removeReaction,
  getUnreadCount,

  // Admin
  approvePendingAward,
  addAdminPoints,
  bulkApprovePendingAwards,
  bulkRejectPendingAwards,
  getAdminProfile,
  getAdminProfiles,
  getAdminWorkers,
  toggleAdminWorker,
  triggerAdminWorker,
  getPendingAwards,
  getPendingAwardsStats,
  rejectPendingAward,
};
