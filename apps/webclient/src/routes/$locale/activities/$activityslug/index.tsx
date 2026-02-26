// Activity detail page
import * as React from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Calendar, Clock, ExternalLink, PencilLine, User } from "lucide-react";
import { PageLayout } from "@/components/page-layouts/default";
import { backend } from "@/modules/backend/backend";
import { compileMdx } from "@/lib/mdx";
import { MdxContent } from "@/components/userland/mdx-content";
import { siteConfig } from "@/config";
import { useAuth } from "@/lib/auth/auth-context";
import { buildUrl, computeContentLanguage, generateCanonicalLink, generateMetaTags, truncateDescription } from "@/lib/seo";
import { setContentLanguageHeader } from "@/lib/set-content-language";
import { formatDateTimeLong } from "@/lib/date";
import { LocaleLink } from "@/components/locale-link";
import { PageNotFound } from "@/components/page-not-found";
import type { ActivityProperties, RSVPMode } from "@/modules/backend/types";
import { RSVPButtons } from "../_components/-rsvp-buttons";

export const Route = createFileRoute("/$locale/activities/$activityslug/")({
  loader: async ({ params }) => {
    const { locale, activityslug } = params;
    const activity = await backend.getActivity(locale, activityslug);

    if (activity === null || activity === undefined) {
      return {
        activity: null,
        compiledContent: null,
        currentUrl: null,
        locale,
        myInteractions: null,
        interactionCounts: null,
        notFound: true as const,
      };
    }

    // Compile MDX content on the server
    let compiledContent: string | null = null;
    if (activity.content !== null && activity.content !== undefined) {
      try {
        compiledContent = await compileMdx(activity.content);
      } catch (error) {
        console.error("Failed to compile MDX:", error);
        compiledContent = null;
      }
    }

    const currentUrl = `${siteConfig.host}/${locale}/activities/${activityslug}`;

    // Fetch interaction data
    const [myInteractions, interactionCounts] = await Promise.all([
      backend.getMyInteractions(locale, activityslug).catch(() => null),
      backend.getInteractionCounts(locale, activityslug).catch(() => null),
    ]);

    // Set Content-Language header with content locale awareness
    if (import.meta.env.SSR) {
      await setContentLanguageHeader({ data: computeContentLanguage(locale, activity.locale_code) });
    }

    return {
      activity,
      compiledContent,
      currentUrl,
      locale,
      myInteractions,
      interactionCounts,
      notFound: false as const,
    };
  },
  head: ({ loaderData }) => {
    if (loaderData === undefined || loaderData.notFound || loaderData.activity === null) {
      return { meta: [] };
    }
    const { activity, currentUrl, locale } = loaderData;
    const activityslug = activity.slug ?? "";
    return {
      meta: generateMetaTags({
        title: activity.title ?? "Activity",
        description: truncateDescription(activity.summary),
        url: currentUrl,
        image: activity.story_picture_uri,
        locale,
        type: "article",
        publishedTime: activity.published_at ?? activity.created_at,
        modifiedTime: activity.updated_at,
        author: activity.author_profile?.title ?? null,
      }),
      links: [generateCanonicalLink(buildUrl(locale, "activities", activityslug))],
    };
  },
  component: ActivityDetailPage,
  notFoundComponent: PageNotFound,
});

function ActivityDetailPage() {
  const params = Route.useParams();
  const loaderData = Route.useLoaderData();
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const auth = useAuth();
  const [canEdit, setCanEdit] = React.useState(false);

  if (loaderData.notFound || loaderData.activity === null) {
    return <PageNotFound />;
  }

  const { activity, compiledContent, myInteractions, interactionCounts } = loaderData;

  // Get author profile slug for permissions check
  const authorProfileSlug = activity.author_profile?.slug ?? null;

  // Check edit permissions
  React.useEffect(() => {
    if (auth.isAuthenticated && !auth.isLoading && authorProfileSlug !== null) {
      backend
        .getStoryPermissions(params.locale, authorProfileSlug, activity.id)
        .then((perms) => {
          if (perms !== null) {
            setCanEdit(perms.can_edit);
          }
        });
    } else {
      setCanEdit(false);
    }
  }, [auth.isAuthenticated, auth.isLoading, params.locale, authorProfileSlug, activity.id]);

  const activityProps = (activity.properties ?? {}) as unknown as ActivityProperties;
  const rsvpMode: RSVPMode = activityProps.rsvp_mode ?? "enabled";

  const timeStart = activityProps.activity_time_start !== undefined
    ? new Date(activityProps.activity_time_start)
    : null;
  const timeEnd = activityProps.activity_time_end !== undefined
    ? new Date(activityProps.activity_time_end)
    : null;

  const activityKindLabels: Record<string, string> = {
    meetup: "Activities.Meetup",
    workshop: "Activities.Workshop",
    conference: "Activities.Conference",
    broadcast: "Activities.Broadcast",
    meeting: "Activities.Meeting",
  };
  const kindLabel = activityKindLabels[activityProps.activity_kind ?? "meetup"] ?? "Activities.Meetup";

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto max-w-4xl">
        <article className="content">
          {/* Activity kind badge */}
          <span className="inline-block text-xs font-medium px-2 py-0.5 rounded-full bg-primary/10 text-primary mb-2">
            {t(kindLabel)}
          </span>

          <h1>{activity.title}</h1>

          {/* Activity metadata */}
          <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground mb-6 not-prose">
            {timeStart !== null && (
              <span className="flex items-center gap-1.5">
                <Clock className="size-4" />
                {formatDateTimeLong(timeStart, locale)}
                {timeEnd !== null && (
                  <> – {formatDateTimeLong(timeEnd, locale)}</>
                )}
              </span>
            )}

            {activity.author_profile !== null && activity.author_profile !== undefined && (
              <LocaleLink
                to={`/${activity.author_profile.slug}`}
                className="flex items-center gap-1.5 text-muted-foreground hover:text-foreground"
              >
                <User className="size-4" />
                {t("Activities.Organized by")} {activity.author_profile.title}
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
          </div>

          {/* Edit link */}
          {canEdit && (
            <div className="flex justify-end mb-4 not-prose">
              <Link
                to="/$locale/stories/$storyslug/edit"
                params={{ locale: params.locale, storyslug: params.activityslug }}
                className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors no-underline"
              >
                <PencilLine className="size-3.5" />
                {t("ContentEditor.Edit Story")}
              </Link>
            </div>
          )}

          {/* RSVP section */}
          <RSVPButtons
            locale={params.locale}
            storySlug={params.activityslug}
            rsvpMode={rsvpMode}
            externalAttendanceUri={activityProps.external_attendance_uri}
            isAuthenticated={auth.isAuthenticated}
            initialMyInteractions={myInteractions}
            initialCounts={interactionCounts}
          />

          {/* Activity content */}
          {compiledContent !== null
            ? <MdxContent compiledSource={compiledContent} />
            : activity.content !== null && activity.content !== undefined && (
              <div className="prose dark:prose-invert">{activity.content}</div>
            )}
        </article>
      </section>
    </PageLayout>
  );
}
