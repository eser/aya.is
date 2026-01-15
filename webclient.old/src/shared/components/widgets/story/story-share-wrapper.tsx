import React, { Suspense } from "react";

import type { StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { StoryShare } from "./story-share.tsx";

export type StoryShareWrapperProps = {
  story: StoryEx;
  currentUrl: string;
};

export async function StoryShareWrapper(props: StoryShareWrapperProps) {
  const { t } = await getTranslations();

  return (
    <Suspense fallback={<div>{t("Loading", "Loading share options...")}</div>}>
      <StoryShare story={props.story} currentUrl={props.currentUrl} />
    </Suspense>
  );
}
