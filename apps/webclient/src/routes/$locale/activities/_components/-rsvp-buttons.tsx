"use client";

import * as React from "react";
import { useTranslation } from "react-i18next";
import { ExternalLink, Check, Star, X, Users } from "lucide-react";
import { cn } from "@/lib/utils";
import { backend } from "@/modules/backend/backend";
import type {
  RSVPMode,
  StoryInteraction,
  InteractionCount,
} from "@/modules/backend/types";
import styles from "./-rsvp-buttons.module.css";

const RSVP_KINDS = ["attending", "interested", "not_attending"] as const;

const rsvpIcons: Record<string, React.ElementType> = {
  attending: Check,
  interested: Star,
  not_attending: X,
};

const rsvpLabels: Record<string, string> = {
  attending: "Activities.Attending",
  interested: "Activities.Interested",
  not_attending: "Activities.Not Attending",
};

export type RSVPButtonsProps = {
  locale: string;
  storySlug: string;
  rsvpMode: RSVPMode;
  externalAttendanceUri?: string;
  isAuthenticated: boolean;
  initialMyInteractions: StoryInteraction[] | null;
  initialCounts: InteractionCount[] | null;
};

export function RSVPButtons(props: RSVPButtonsProps) {
  const { t } = useTranslation();
  const [myInteractions, setMyInteractions] = React.useState<StoryInteraction[]>(
    props.initialMyInteractions ?? [],
  );
  const [counts, setCounts] = React.useState<InteractionCount[]>(
    props.initialCounts ?? [],
  );
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  if (props.rsvpMode === "disabled") {
    return null;
  }

  if (props.rsvpMode === "managed_externally") {
    if (props.externalAttendanceUri === undefined || props.externalAttendanceUri === "") {
      return null;
    }

    return (
      <div className={styles.container}>
        <a
          href={props.externalAttendanceUri}
          target="_blank"
          rel="noopener noreferrer"
          className={styles.externalLink}
        >
          <ExternalLink className="size-4" />
          {t("Activities.Register Externally")}
        </a>
      </div>
    );
  }

  const activeKinds = new Set(myInteractions.map((i) => i.kind));

  const handleRSVP = async (kind: string) => {
    if (!props.isAuthenticated || isSubmitting) {
      return;
    }

    setIsSubmitting(true);
    try {
      if (activeKinds.has(kind)) {
        await backend.removeInteraction(props.locale, props.storySlug, kind);
        setMyInteractions((prev) => prev.filter((i) => i.kind !== kind));
      } else {
        const result = await backend.setInteraction(props.locale, props.storySlug, kind);
        if (result !== null) {
          setMyInteractions((prev) => {
            const withoutRsvp = prev.filter(
              (i) => !RSVP_KINDS.includes(i.kind as typeof RSVP_KINDS[number]),
            );
            return [...withoutRsvp, result];
          });
        }
      }

      const updatedCounts = await backend.getInteractionCounts(
        props.locale,
        props.storySlug,
      );
      if (updatedCounts !== null) {
        setCounts(updatedCounts);
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const getCount = (kind: string): number => {
    const found = counts.find((c) => c.kind === kind);
    return found !== undefined ? found.count : 0;
  };

  return (
    <div className={styles.container}>
      {props.isAuthenticated && (
        <div className={styles.buttonGroup}>
          {RSVP_KINDS.map((kind) => {
            const Icon = rsvpIcons[kind];
            const isActive = activeKinds.has(kind);
            return (
              <button
                key={kind}
                type="button"
                className={cn(
                  styles.rsvpButton,
                  isActive ? styles.active : styles.inactive,
                )}
                onClick={() => handleRSVP(kind)}
                disabled={isSubmitting}
              >
                {Icon !== undefined && <Icon className="size-4 inline mr-1" />}
                {t(rsvpLabels[kind])}
              </button>
            );
          })}
        </div>
      )}

      <div className={styles.counts}>
        {RSVP_KINDS.map((kind) => {
          const count = getCount(kind);
          if (count === 0 && !activeKinds.has(kind)) {
            return null;
          }
          return (
            <span key={kind} className={styles.countItem}>
              <Users className="size-3.5" />
              {count} {t(rsvpLabels[kind])}
            </span>
          );
        })}
      </div>
    </div>
  );
}
