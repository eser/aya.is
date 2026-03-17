// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import type * as React from "react";
import type { StoryEx } from "@/modules/backend/types";
import { TextContent } from "@/components/widgets/text-content";
import { LocaleBadge } from "@/components/locale-badge";
import { StoryInformation } from "./story-information";
import { StoryFooter } from "./story-footer";

export type StoryContentProps = {
  story: StoryEx;
  compiledContent: string | null;
  currentUrl: string;
  locale: string;
  showAuthor?: boolean;
  showPublications?: boolean;
  showShare?: boolean;
  headingOffset?: number;
  editUrl?: string;
  coverUrl?: string;
  shareUrl?: string;
  beforeContent?: React.ReactNode;
};

const headingTags = ["h1", "h2", "h3", "h4", "h5", "h6"] as const;

export function StoryContent(props: StoryContentProps) {
  const {
    story,
    compiledContent,
    currentUrl,
    locale,
    showAuthor = true,
    showPublications,
    showShare = true,
    headingOffset = 1,
    editUrl,
    coverUrl,
    shareUrl,
    beforeContent,
  } = props;

  const TitleTag = headingTags[headingOffset - 1] ?? "h1";

  return (
    <article className="content">
      <TitleTag>
        {story.title}
        <LocaleBadge
          localeCode={story.locale_code}
          className="text-sm font-medium px-2 py-0.5 rounded-full bg-primary/10 text-primary ml-3 align-middle"
        />
      </TitleTag>

      <StoryInformation story={story} locale={locale} editUrl={editUrl} coverUrl={coverUrl} shareUrl={shareUrl} />

      {beforeContent}

      <TextContent
        compiledContent={compiledContent}
        rawContent={story.content}
        shareOptions={showShare
          ? {
            title: story.title ?? "",
            summary: story.summary,
            slug: story.slug,
            currentUrl,
          }
          : undefined}
        bare
        headingOffset={headingOffset}
      />

      {showAuthor && <StoryFooter story={story} showPublications={showPublications} />}
    </article>
  );
}
