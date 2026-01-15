// Profile layout
import {
  createFileRoute,
  Outlet,
  notFound,
  useParams,
} from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { forbiddenSlugs } from "@/config";
import { backend, type Profile } from "@/modules/backend/backend";
import { LocaleLink } from "@/components/locale-link";
import { Icons } from "@/components/icons";
import { PageLayout } from "@/components/page-layouts/default";

function findIcon(kind: string) {
  switch (kind) {
    case "github":
      return Icons.github;
    case "twitter":
    case "x":
      return Icons.twitter;
    case "linkedin":
      return Icons.linkedin;
    case "instagram":
      return Icons.instagram;
    case "youtube":
      return Icons.youtube;
    default:
      return Icons.link;
  }
}

export const Route = createFileRoute("/$locale/$slug")({
  beforeLoad: async ({ params }) => {
    const { slug } = params;

    // Check for forbidden slugs
    if (forbiddenSlugs.includes(slug.toLowerCase())) {
      throw notFound();
    }

    return { profileSlug: slug };
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;

    const profile = await backend.getProfile(locale, slug);

    if (profile === null) {
      throw notFound();
    }

    return { profile };
  },
  component: ProfileLayout,
  notFoundComponent: ProfileNotFound,
});

function ProfileNotFound() {
  return (
    <PageLayout>
      <div className="container mx-auto py-16 px-4 text-center">
        <h1 className="text-4xl font-bold mb-4">Profile Not Found</h1>
        <p className="text-muted-foreground">
          The profile you're looking for doesn't exist.
        </p>
      </div>
    </PageLayout>
  );
}

function ProfileLayout() {
  const params = useParams({ strict: false });
  const { profile } = Route.useLoaderData();
  const slug = (params as { slug?: string }).slug || "";

  return (
    <PageLayout>
      <section className="container px-4 py-8 mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-[280px_1fr] gap-8 items-start">
          <ProfileSidebar profile={profile} slug={slug} />
          <main>
            <Outlet />
          </main>
        </div>
      </section>
    </PageLayout>
  );
}

interface ProfileSidebarProps {
  profile: Profile;
  slug: string;
}

function ProfileSidebar({ profile, slug }: ProfileSidebarProps) {
  const { t } = useTranslation();

  return (
    <aside className="flex flex-col gap-4">
      {/* Profile Picture */}
      {profile.profile_picture_uri && (
        <div className="flex justify-center md:justify-start">
          <img
            src={profile.profile_picture_uri}
            alt={`${profile.title}'s profile picture`}
            width={280}
            height={280}
            className="border rounded-full"
          />
        </div>
      )}

      {/* Hero Section */}
      <div>
        <h1 className="mt-0 mb-2 font-serif text-base font-semibold leading-none sm:text-lg md:text-xl lg:text-2xl">
          {profile.title}
        </h1>

        <div className="mt-0 mb-4 font-sans text-sm font-light leading-none sm:text-base md:text-lg lg:text-xl text-muted-foreground">
          {profile.slug}
          {profile.pronouns && ` · ${profile.pronouns}`}
        </div>

        {profile.links && profile.links.length > 0 && (
          <div className="flex gap-4 mb-3 text-sm text-muted">
            {profile.links.map((link) => {
              const Icon = findIcon(link.kind);
              return (
                <a
                  key={link.id}
                  href={link.uri}
                  title={link.title || link.kind}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="no-underline"
                >
                  <Icon className="stroke-muted-foreground hover:stroke-foreground h-5 w-5" />
                </a>
              );
            })}
          </div>
        )}

        {profile.description && (
          <p className="mt-0 mb-4 font-sans text-sm font-normal leading-snug text-left">
            {profile.description}
          </p>
        )}
      </div>

      {/* Navigation */}
      <nav className="flex justify-center font-serif md:justify-start">
        <ul className="flex flex-row flex-wrap justify-center p-0 space-y-0 md:space-y-3 lg:space-y-4 list-none md:flex-col">
          <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
            <LocaleLink
              to={`/${slug}`}
              className="no-underline text-muted-foreground hover:text-foreground"
            >
              {t("Layout.Profile")}
            </LocaleLink>
          </li>

          {profile.kind === "individual" && (
            <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
              <LocaleLink
                to={`/${slug}/contributions`}
                className="no-underline text-muted-foreground hover:text-foreground"
              >
                {t("Layout.Contributions")}
              </LocaleLink>
            </li>
          )}

          {(profile.kind === "organization" ||
            profile.kind === "project" ||
            profile.kind === "product") && (
            <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
              <LocaleLink
                to={`/${slug}/members`}
                className="no-underline text-muted-foreground hover:text-foreground"
              >
                {t("Layout.Members")}
              </LocaleLink>
            </li>
          )}

          {profile.pages?.map((page) => (
            <li
              key={page.slug}
              className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none"
            >
              <LocaleLink
                to={`/${slug}/${page.slug}`}
                className="no-underline text-muted-foreground hover:text-foreground"
              >
                {page.title}
              </LocaleLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
