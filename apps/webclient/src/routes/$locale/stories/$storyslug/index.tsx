// Individual story page (handles all story kinds including activities)
import * as React from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Clock, ExternalLink, PencilLine, User } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { dateProposalsQueryOptions, storyQueryOptions } from "@/modules/backend/queries";
import { QueryError } from "@/components/query-error";
import { StoryContent } from "@/components/widgets/story-content";
import { DiscussionThread } from "@/components/widgets/discussion-thread";
import { RSVPButtons } from "@/components/widgets/rsvp-buttons";
import { DatePoll } from "@/components/widgets/date-poll";
import { compileMdx, compileMdxLite } from "@/lib/mdx";
import { MdxContent } from "@/components/userland/mdx-content";
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
import { formatDateTimeLong, formatDateTimeRange } from "@/lib/date";
import { LocaleLink } from "@/components/locale-link";
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
    const activityProperties = isActivity
      ? (story.properties ?? {}) as unknown as ActivityProperties
      : null;
    const dateMode: DateMode = activityProperties?.date_mode ?? "fixed";

    // Fetch activity-specific data in parallel with discussion data
    const [initialDiscussion, myInteractions, interactionCounts, dateProposals] = await Promise.all([
      // Pre-fetch discussion data for SSR (all story kinds)
      story.feat_discussions === true
        ? fetchDiscussionData(locale, storyslug)
        : Promise.resolve(null),
      // Activity: fetch interaction data
      isActivity
        ? backend.getMyInteractions(locale, storyslug).catch(() => null)
        : Promise.resolve(null),
      isActivity
        ? backend.getInteractionCounts(locale, storyslug).catch(() => null)
        : Promise.resolve(null),
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

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto max-w-4xl">
        {story.kind === "activity"
          ? (
            <ActivityDetail
              story={story}
              compiledContent={compiledContent}
              locale={params.locale}
              storySlug={params.storyslug}
              isAuthenticated={auth.isAuthenticated}
              canEdit={canEdit}
              myInteractions={loaderData.myInteractions}
              interactionCounts={loaderData.interactionCounts}
              dateProposals={loaderData.dateProposals}
              dateMode={loaderData.dateMode}
            />
          )
          : (
            <StoryContent
              story={story}
              compiledContent={compiledContent}
              currentUrl={currentUrl ?? ""}
              showAuthor
              editUrl={editUrl}
              coverUrl={coverUrl}
            />
          )}

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

// --- Activity-specific rendering ---

type ActivityDetailProps = {
  story: StoryEx;
  compiledContent: string | null;
  locale: string;
  storySlug: string;
  isAuthenticated: boolean;
  canEdit: boolean;
  myInteractions: StoryInteraction[] | null;
  interactionCounts: InteractionCount[] | null;
  dateProposals: DateProposalListResponse | null;
  dateMode: DateMode;
};

const activityKindLabels: Record<string, string> = {
  meetup: "Activities.Meetup",
  workshop: "Activities.Workshop",
  conference: "Activities.Conference",
  broadcast: "Activities.Broadcast",
  meeting: "Activities.Meeting",
};

function ActivityDetail(props: ActivityDetailProps) {
  const { t } = useTranslation();
  const locale = props.locale;

  const activityProps = (props.story.properties ?? {}) as unknown as ActivityProperties;
  const rsvpMode: RSVPMode = activityProps.rsvp_mode ?? "enabled";

  const timeStart = activityProps.activity_time_start !== undefined
    ? new Date(activityProps.activity_time_start)
    : null;
  const timeEnd = activityProps.activity_time_end !== undefined ? new Date(activityProps.activity_time_end) : null;

  const kindLabel = activityKindLabels[activityProps.activity_kind ?? "meetup"] ?? "Activities.Meetup";

  return (
    <article className="content">
      {/* Activity kind badge */}
      <span className="inline-block text-xs font-medium px-2 py-0.5 rounded-full bg-primary/10 text-primary mb-2">
        {t(kindLabel)}
      </span>

      <h1>{props.story.title}</h1>

      {/* Activity metadata */}
      <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground mb-6 not-prose">
        {props.dateMode === "undecided"
          ? (
            <span className="flex items-center gap-1.5">
              <Clock className="size-4" />
              <span className="inline-block px-2 py-0.5 rounded-full bg-amber-100 text-amber-800 text-xs font-medium dark:bg-amber-900 dark:text-amber-200">
                {t("Activities.Date Undecided")}
              </span>
            </span>
          )
          : timeStart !== null && (
            <span className="flex items-center gap-1.5">
              <Clock className="size-4" />
              {timeEnd !== null
                ? formatDateTimeRange(timeStart, timeEnd, locale)
                : formatDateTimeLong(timeStart, locale)}
            </span>
          )}

        {props.story.author_profile !== null && props.story.author_profile !== undefined && (
          <LocaleLink
            to={`/${props.story.author_profile.slug}`}
            className="flex items-center gap-1.5 text-muted-foreground hover:text-foreground"
          >
            <User className="size-4" />
            {t("Activities.Organized by")} {props.story.author_profile.title}
          </LocaleLink>
        )}

        {activityProps.external_activity_uri !== undefined &&
          activityProps.external_activity_uri !== "" && (
          <a
            href={activityProps.external_activity_uri}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 text-primary hover:underline"
          >
            <ExternalLink className="size-4" />
            {t("Common.View")}
          </a>
        )}

        {/* Edit link */}
        {props.canEdit && (
          <Link
            to="/$locale/stories/$storyslug/edit"
            params={{ locale, storyslug: props.storySlug }}
            className="ml-auto flex items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors no-underline"
          >
            <PencilLine className="size-3.5" />
            {t("ContentEditor.Edit Story")}
          </Link>
        )}
      </div>

      {/* RSVP section */}
      <RSVPButtons
        locale={locale}
        storySlug={props.storySlug}
        rsvpMode={rsvpMode}
        externalAttendanceUri={activityProps.external_attendance_uri}
        isAuthenticated={props.isAuthenticated}
        initialMyInteractions={props.myInteractions}
        initialCounts={props.interactionCounts}
      />

      {/* Date Proposals (when date is undecided) */}
      {props.dateMode === "undecided" && (
        <DatePoll
          locale={locale}
          storySlug={props.storySlug}
          canPropose={props.dateProposals?.viewer_can_propose ?? false}
          canVote={props.dateProposals?.viewer_can_vote ?? false}
          canEdit={props.canEdit}
          initialProposals={props.dateProposals?.proposals ?? null}
        />
      )}

      {/* Activity content */}
      {props.compiledContent !== null
        ? <MdxContent compiledSource={props.compiledContent} />
        : props.story.content !== null && props.story.content !== undefined && (
          <div className="prose dark:prose-invert">{props.story.content}</div>
        )}
    </article>
  );
}
