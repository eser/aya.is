// Centralized query option factories for React Query.
// Used in route loaders (ensureQueryData) and components (useSuspenseQuery).
//
// Key structure: [entity, locale, ...identifiers, { filters }]
// The { filters } object is reserved for future pagination params.
import { queryOptions } from "@tanstack/react-query";
import { getDailySeed } from "@/lib/seed-utils";
import { backend } from "./backend";

// === Stories ===

export const storiesQueryOptions = (locale: string) =>
  queryOptions({
    queryKey: ["stories", locale],
    queryFn: () => backend.getStories(locale),
  });

export const storiesByKindsQueryOptions = (locale: string, kinds: string[]) =>
  queryOptions({
    queryKey: ["stories", locale, { kinds }],
    queryFn: () => backend.getStoriesByKinds(locale, kinds),
  });

export const storyQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["story", locale, slug],
    queryFn: () => backend.getStory(locale, slug),
  });

export const storyDiscussionQueryOptions = (locale: string, storySlug: string) =>
  queryOptions({
    queryKey: ["story-discussion", locale, storySlug],
    queryFn: () => backend.getStoryDiscussion(locale, storySlug),
  });

// === Activities ===

export const activitiesQueryOptions = (locale: string) =>
  queryOptions({
    queryKey: ["activities", locale],
    queryFn: () => backend.getActivities(locale),
  });

export const activityQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["activity", locale, slug],
    queryFn: () => backend.getActivity(locale, slug),
  });

// === Profiles ===

export const profilesByKindsQueryOptions = (
  locale: string,
  kinds: string[],
  options?: { seed?: string; limit?: number; offset?: number; q?: string },
) => {
  const seed = options?.seed ?? getDailySeed();
  const limit = options?.limit ?? 24;
  const offset = options?.offset ?? 0;
  const q = options?.q ?? "";

  return queryOptions({
    queryKey: ["profiles", locale, { kinds, seed, limit, offset, q }],
    queryFn: () => backend.getProfilesByKinds(locale, kinds, seed, limit, offset, q),
    placeholderData: (prev: unknown) => prev, // keep stale data visible while refetching
  });
};

export const profileQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["profile", locale, slug],
    queryFn: () => backend.getProfile(locale, slug),
  });

export const profileStoriesQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["profile-stories", locale, slug],
    queryFn: () => backend.getProfileStories(locale, slug),
  });

export const profileMembersQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["profile-members", locale, slug],
    queryFn: () => backend.getProfileMembers(locale, slug),
  });

export const profileLinksQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["profile-links", locale, slug],
    queryFn: () => backend.getProfileLinks(locale, slug),
  });

export const profileContributionsQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["profile-contributions", locale, slug],
    queryFn: () => backend.getProfileContributions(locale, slug),
  });

export const profilePermissionsQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["profile-permissions", locale, slug],
    queryFn: () => backend.getProfilePermissions(locale, slug),
    retry: false, // 401 for unauthenticated — don't retry
    staleTime: 0, // Always refetch — SSR may lack auth on custom domains
  });

export const profileQuestionsQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["profile-questions", locale, slug],
    queryFn: () => backend.getProfileQuestions(locale, slug),
  });

// === Search ===

export const searchQueryOptions = (locale: string, query: string, profileSlug?: string) =>
  queryOptions({
    queryKey: ["search", locale, query, { profileSlug }],
    queryFn: () => backend.search(locale, query, profileSlug),
  });

// === Live ===

export const liveNowQueryOptions = (locale: string) =>
  queryOptions({
    queryKey: ["live-now", locale],
    queryFn: () => backend.getLiveNow(locale),
    staleTime: 30_000, // 30s — live data refreshes more frequently
  });

// === Profile Story ===

export const profileStoryQueryOptions = (locale: string, slug: string, storySlug: string) =>
  queryOptions({
    queryKey: ["profile-story", locale, slug, storySlug],
    queryFn: () => backend.getProfileStory(locale, slug, storySlug),
  });

// === Profile Page ===

export const profilePageQueryOptions = (locale: string, slug: string, pageSlug: string) =>
  queryOptions({
    queryKey: ["profile-page", locale, slug, pageSlug],
    queryFn: () => backend.getProfilePage(locale, slug, pageSlug),
  });

// === Referrals ===

export const referralsQueryOptions = (locale: string, slug: string) =>
  queryOptions({
    queryKey: ["referrals", locale, slug],
    queryFn: () => backend.listReferrals(locale, slug),
  });

// === Interactions ===

export const myInteractionsQueryOptions = (locale: string, storySlug: string) =>
  queryOptions({
    queryKey: ["my-interactions", locale, storySlug],
    queryFn: () => backend.getMyInteractions(locale, storySlug),
    retry: false,
  });

export const interactionCountsQueryOptions = (locale: string, storySlug: string) =>
  queryOptions({
    queryKey: ["interaction-counts", locale, storySlug],
    queryFn: () => backend.getInteractionCounts(locale, storySlug),
  });
