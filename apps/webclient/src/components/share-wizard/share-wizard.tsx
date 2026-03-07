import { useCallback, useState } from "react";
import type { StoryEx } from "@/modules/backend/types";
import { ShareActions } from "./share-actions";
import { ComposePanel } from "./compose-panel";
import { PlatformPreview } from "./platform-preview";
import { AiAssistPanel } from "./ai-assist-panel";
import { ImageComposer } from "./image-composer";
import styles from "./share-wizard.module.css";

export type ShareWizardProps = {
  story: StoryEx;
  locale: string;
  currentUrl: string;
  onBack: () => void;
};

export function ShareWizard(props: ShareWizardProps) {
  const [text, setText] = useState("");
  const [platformOverrides, setPlatformOverrides] = useState<Map<string, string>>(new Map());
  const [activePlatform, setActivePlatform] = useState("x");

  const handlePlatformOverrideChange = useCallback((platformId: string, overrideText: string | null) => {
    setPlatformOverrides((prev) => {
      const next = new Map(prev);
      if (overrideText === null) {
        next.delete(platformId);
      } else {
        next.set(platformId, overrideText);
      }
      return next;
    });
  }, []);

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <ShareActions
          storyTitle={props.story.title ?? ""}
          currentUrl={props.currentUrl}
          onBack={props.onBack}
        />
      </div>

      {/* Main content */}
      <div className={styles.content}>
        {/* Left: Compose panel */}
        <div className={styles.mainPanel}>
          <ComposePanel
            text={text}
            story={props.story}
            currentUrl={props.currentUrl}
            onTextChange={setText}
          />

          <PlatformPreview
            text={text}
            story={props.story}
            currentUrl={props.currentUrl}
            platformOverrides={platformOverrides}
            activePlatform={activePlatform}
            onActivePlatformChange={setActivePlatform}
            onPlatformOverrideChange={handlePlatformOverrideChange}
          />
        </div>

        {/* Right: AI Assist + Image Composer */}
        <div className={styles.sidePanel}>
          <AiAssistPanel
            text={text}
            storyContent={props.story.content ?? null}
            storySlug={props.story.slug ?? ""}
            locale={props.locale}
            activePlatform={activePlatform}
          />

          <ImageComposer
            story={props.story}
            locale={props.locale}
          />
        </div>
      </div>
    </div>
  );
}
