import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import { SiteAvatar } from "@/components/userland";
import type { StoryEx } from "@/modules/backend/types";

export type StoryFooterProps = {
  story: StoryEx;
};

export function StoryFooter(props: StoryFooterProps) {
  const { t } = useTranslation();

  if (
    props.story.author_profile === null ||
    props.story.author_profile === undefined
  ) {
    return null;
  }

  // Filter out publications where the id matches the author profile id
  const filteredPublications = props.story.publications.filter(
    (publication) => publication.id !== props.story.author_profile.id,
  );

  return (
    <div className="w-full flex items-start gap-12 pt-20 pb-10">
      <div className="w-4/6 flex flex-row gap-4 items-center">
        <div className="flex-none p-0.5 bg-card rounded-full shadow-md">
          <SiteAvatar
            src={props.story.author_profile.profile_picture_uri}
            name={props.story.author_profile.title ?? "Author"}
            fallbackName={props.story.author_profile.slug}
            className="size-24 border-2 border-background"
          />
        </div>
        <div>
          <div className="text-sm text-muted-foreground">
            {t("Stories.Written by")}
          </div>
          <div className="text-lg font-semibold text-foreground">
            <Link to={`/${props.story.author_profile.slug}`}>
              {props.story.author_profile.title}
            </Link>
          </div>
          <div className="text-sm text-muted-foreground">
            {props.story.author_profile.description}
          </div>
        </div>
      </div>

      {filteredPublications.length > 0 && (
        <div className="w-2/6 flex flex-col">
          <div className="text-sm text-muted-foreground">
            {t("Stories.Publications")}
          </div>
          <div className="flex flex-row gap-2 flex-wrap">
            {filteredPublications.map((publication) => (
              <div key={publication.id} className="text-sm text-foreground">
                <Link to={`/${publication.slug}`}>{publication.title}</Link>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
