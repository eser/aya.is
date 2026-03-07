// Individual story page (handles all story kinds including activities)
import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { dateProposalsQueryOptions, storyQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { StoryContent } from "@/components/widgets/story-content";
import { DiscussionThread } from "@/components/widgets/discussion-thread";
import { RSVPButtons } from "@/components/widgets/rsvp-buttons";
import { DatePoll } from "@/components/widgets/date-poll";
import { compileMdx, compileMdxLite } from "@/lib/mdx";
import { siteConfig } from "@/config";
import { useAuth } from "@/lib/auth/auth-context";
import {
  computeContentLanguage,
  computeStoryCanonicalUrl,
  generateCanonicalLink,
  generateMetaTags,
  truncateDescription,
} from "@/lib/seo";
import { setServerResponseHeader } from "@/lib/server-headers";
import { PageNotFound } from "@/components/page-not-found";
import type {
  ActivityProperties,
  DateMode,
  DateProposalListResponse,
  DiscussionComment,
  DiscussionListResponse,
  InteractionCount,
  RSVPMode,
  StoryEx,
  StoryInteraction,
} from "@/modules/backend/types";

export const Route = createFileRoute("/$locale/stories/$storyslug/")({
  loader: async ({ params, context }) => {
    const { locale, storyslug } = params;
    const story = await context.queryClient.ensureQueryData(storyQueryOptions(locale, storyslug));

    if (story === null || story === undefined) {
      return {
        story: null,
        compiledContent: null,
        currentUrl: null,
        locale,
        initialDiscussion: null,
        // Activity-specific fields
        myInteractions: null,
        interactionCounts: null,
        dateProposals: null,
        dateMode: "fixed" as DateMode,
        notFound: true as const,
      };
    }

    // Compile MDX content on the server
    let compiledContent: string | null = null;
    if (story.content !== null && story.content !== undefined) {
      try {
        compiledContent = await compileMdx(story.content);
      } catch (error) {
        console.error("Failed to compile MDX:", error);
        compiledContent = null;
      }
    }

    // Get current URL for sharing
    const currentUrl = `${siteConfig.host}/${locale}/stories/${storyslug}`;

    // Activity-specific data
    const isActivity = story.kind === "activity";
    const activityProperties = isActivity ? (story.properties ?? {}) as unknown as ActivityProperties : null;
    const dateMode: DateMode = activityProperties?.date_mode ?? "fixed";

    // Fetch activity-specific data in parallel with discussion data
    const [initialDiscussion, myInteractions, interactionCounts, dateProposals] = await Promise.all([
      // Pre-fetch discussion data for SSR (all story kinds)
      story.feat_discussions === true ? fetchDiscussionData(locale, storyslug) : Promise.resolve(null),
      // Activity: fetch interaction data
      isActivity ? backend.getMyInteractions(locale, storyslug).catch(() => null) : Promise.resolve(null),
      isActivity ? backend.getInteractionCounts(locale, storyslug).catch(() => null) : Promise.resolve(null),
      // Activity with undecided date: fetch date proposals
      isActivity && dateMode === "undecided"
        ? context.queryClient.ensureQueryData(dateProposalsQueryOptions(locale, storyslug)).catch(() => null)
        : Promise.resolve(null),
    ]);

    // Set Content-Language header with content locale awareness
    setServerResponseHeader("Content-Language", computeContentLanguage(locale, story.locale_code));

    return {
      story,
      compiledContent,
      currentUrl,
      locale,
      initialDiscussion,
      myInteractions,
      interactionCounts,
      dateProposals,
      dateMode,
      notFound: false as const,
    };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined || loaderData.notFound || loaderData.story === null) {
      return { meta: [] };
    }
    const { story, currentUrl, locale } = loaderData;
    const canonicalUrl = computeStoryCanonicalUrl(story, locale, "stories");
    return {
      meta: generateMetaTags({
        title: story.title ?? "Story",
        description: truncateDescription(story.summary),
        url: currentUrl,
        image: story.story_picture_uri,
        locale,
        type: "article",
        publishedTime: story.published_at ?? story.created_at,
        modifiedTime: story.updated_at,
        author: story.author_profile?.title ?? null,
      }),
      links: [generateCanonicalLink(canonicalUrl)],
    };
  },
  errorComponent: QueryError,
  component: StoryPage,
  notFoundComponent: PageNotFound,
});

async function fetchDiscussionData(
  locale: string,
  storyslug: string,
): Promise<{ thread: DiscussionListResponse["thread"]; comments: DiscussionComment[] } | null> {
  try {
    const discussion = await backend.getStoryDiscussion(locale, storyslug, "hot");
    if (discussion === null || discussion === undefined) {
      return null;
    }

    const compiledComments = await Promise.all(
      discussion.comments.map(async (comment) => {
        if (comment.content === "") {
          return { ...comment, compiledContent: null };
        }
        try {
          const compiled = await compileMdxLite(comment.content);
          return { ...comment, compiledContent: compiled };
        } catch {
          return { ...comment, compiledContent: null };
        }
      }),
    );
    return { thread: discussion.thread, comments: compiledComments };
  } catch {
    return null;
  }
}

function StoryPage() {
  const params = Route.useParams();
  const auth = useAuth();
  const loaderData = Route.useLoaderData();
  const [canEdit, setCanEdit] = React.useState(false);

  if (loaderData.notFound || loaderData.story === null) {
    return <PageNotFound />;
  }

  const { story, compiledContent, currentUrl, initialDiscussion } = loaderData;

  // Get author profile slug for permissions check
  const authorProfileSlug = story.author_profile?.slug ?? null;

  // Check edit permissions
  React.useEffect(() => {
    if (auth.isAuthenticated && !auth.isLoading && authorProfileSlug !== null) {
      backend
        .getStoryPermissions(params.locale, authorProfileSlug, story.id)
        .then((perms) => {
          if (perms !== null) {
            setCanEdit(perms.can_edit);
          }
        });
    } else {
      setCanEdit(false);
    }
  }, [auth.isAuthenticated, auth.isLoading, params.locale, authorProfileSlug, story.id]);

  const editUrl = canEdit ? `/${params.locale}/stories/${params.storyslug}/edit` : undefined;
  const coverUrl = canEdit ? `/${params.locale}/stories/${params.storyslug}/cover` : undefined;
  const shareUrl = canEdit ? `/${params.locale}/stories/${params.storyslug}/share` : undefined;

  // Activity-specific widgets injected between metadata and content
  const beforeContent = story.kind === "activity"
    ? buildActivityBeforeContent(story, params.locale, params.storyslug, auth.isAuthenticated, canEdit, loaderData)
    : undefined;

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto max-w-4xl">
        <StoryContent
          story={story}
          compiledContent={compiledContent}
          currentUrl={currentUrl ?? ""}
          locale={params.locale}
          showAuthor
          editUrl={editUrl}
          coverUrl={coverUrl}
          shareUrl={shareUrl}
          beforeContent={beforeContent}
        />

        {story.feat_discussions === true && (
          <DiscussionThread
            storySlug={params.storyslug}
            locale={params.locale}
            profileId={story.author_profile?.id ?? ""}
            profileKind={story.author_profile?.kind ?? "individual"}
            initialData={initialDiscussion}
          />
        )}
      </section>
    </PageLayout>
  );
}

// --- Activity-specific before-content (RSVP + DatePoll) ---

function buildActivityBeforeContent(
  story: StoryEx,
  locale: string,
  storySlug: string,
  isAuthenticated: boolean,
  canEdit: boolean,
  loaderData: {
    myInteractions: StoryInteraction[] | null;
    interactionCounts: InteractionCount[] | null;
    dateProposals: DateProposalListResponse | null;
    dateMode: DateMode;
  },
): React.ReactNode {
  const activityProps = (story.properties ?? {}) as unknown as ActivityProperties;
  const rsvpMode: RSVPMode = activityProps.rsvp_mode ?? "enabled";

  return (
    <>
      <RSVPButtons
        locale={locale}
        storySlug={storySlug}
        rsvpMode={rsvpMode}
        externalAttendanceUri={activityProps.external_attendance_uri}
        isAuthenticated={isAuthenticated}
        initialMyInteractions={loaderData.myInteractions}
        initialCounts={loaderData.interactionCounts}
      />

      {loaderData.dateMode === "undecided" && (
        <DatePoll
          locale={locale}
          storySlug={storySlug}
          canPropose={loaderData.dateProposals?.viewer_can_propose ?? false}
          canVote={loaderData.dateProposals?.viewer_can_vote ?? false}
          canEdit={canEdit}
          initialProposals={loaderData.dateProposals?.proposals ?? null}
        />
      )}
    </>
  );
}
