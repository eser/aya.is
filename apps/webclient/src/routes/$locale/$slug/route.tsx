// Profile layout
import { createFileRoute, notFound, Outlet, useParams } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
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
  beforeLoad: ({ params }) => {
    const { slug } = params;
    return { profileSlug: slug };
  },
  loader: async ({ params }) => {
    const { slug, locale } = params;

    const profile = await backend.getProfile(locale, slug);

    if (profile === null) {
      return { profile: null, notFound: true };
    }

    return { profile, notFound: false };
  },
  component: ProfileLayout,
  notFoundComponent: ProfileNotFound,
});

function ProfileNotFound() {
  const { t } = useTranslation();

  // Simple 404 page - PageLayout handles the header/footer
  return (
    <PageLayout>
      <div className="container mx-auto py-16 px-4 text-center">
        <h1 className="text-4xl font-bold mb-4">{t("Layout.Page not found")}</h1>
        <p className="text-muted-foreground">
          {t("Layout.The page you are looking for does not exist. Please check your spelling and try again.")}
        </p>
      </div>
    </PageLayout>
  );
}

function ProfileLayout() {
  const params = useParams({ strict: false });
  const loaderData = Route.useLoaderData();
  const slug = (params as { slug?: string }).slug ?? "";

  // If notFound flag is set, render 404 page
  if (loaderData.notFound || loaderData.profile === null) {
    return <ProfileNotFound />;
  }

  const { profile } = loaderData;

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

type ProfileSidebarProps = {
  profile: Profile;
  slug: string;
};

function ProfileSidebar(props: ProfileSidebarProps) {
  const { t } = useTranslation();

  return (
    <aside className="flex flex-col gap-4">
      {/* Profile Picture */}
      {props.profile.profile_picture_uri !== null &&
        props.profile.profile_picture_uri !== undefined && (
        <div className="flex justify-center md:justify-start">
          <img
            src={props.profile.profile_picture_uri}
            alt={`${props.profile.title}'s profile picture`}
            width={280}
            height={280}
            className="border rounded-full"
          />
        </div>
      )}

      {/* Hero Section */}
      <div>
        <h1 className="mt-0 mb-2 font-serif text-base font-semibold leading-none sm:text-lg md:text-xl lg:text-2xl">
          {props.profile.title}
        </h1>

        <div className="mt-0 mb-4 font-sans text-sm font-light leading-none sm:text-base md:text-lg lg:text-xl text-muted-foreground">
          {props.profile.slug}
          {props.profile.pronouns !== null &&
            props.profile.pronouns !== undefined &&
            ` · ${props.profile.pronouns}`}
        </div>

        {props.profile.links !== null &&
          props.profile.links !== undefined &&
          props.profile.links.length > 0 && (
          <div className="flex gap-4 mb-3 text-sm text-muted">
            {props.profile.links.map((link) => {
              const Icon = findIcon(link.kind);
              return (
                <a
                  key={link.id}
                  href={link.uri}
                  title={link.title !== null && link.title !== undefined ? link.title : link.kind}
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

        {props.profile.description !== null &&
          props.profile.description !== undefined && (
          <p className="mt-0 mb-4 font-sans text-sm font-normal leading-snug text-left">
            {props.profile.description}
          </p>
        )}
      </div>

      {/* Navigation */}
      <nav className="flex justify-center font-serif md:justify-start">
        <ul className="flex flex-row flex-wrap justify-center p-0 space-y-0 md:space-y-3 lg:space-y-4 list-none md:flex-col">
          <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
            <LocaleLink
              to={`/${props.slug}`}
              className="no-underline text-muted-foreground hover:text-foreground"
            >
              {t("Layout.Profile")}
            </LocaleLink>
          </li>

          {props.profile.kind === "individual" && (
            <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
              <LocaleLink
                to={`/${props.slug}/contributions`}
                className="no-underline text-muted-foreground hover:text-foreground"
              >
                {t("Layout.Contributions")}
              </LocaleLink>
            </li>
          )}

          {(props.profile.kind === "organization" ||
            props.profile.kind === "project" ||
            props.profile.kind === "product") && (
            <li className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none">
              <LocaleLink
                to={`/${props.slug}/members`}
                className="no-underline text-muted-foreground hover:text-foreground"
              >
                {t("Layout.Members")}
              </LocaleLink>
            </li>
          )}

          {props.profile.pages?.map((page) => (
            <li
              key={page.slug}
              className="relative text-base leading-none sm:text-lg md:text-xl lg:text-2xl after:px-2 after:content-['·'] md:after:content-none"
            >
              <LocaleLink
                to={`/${props.slug}/${page.slug}`}
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
