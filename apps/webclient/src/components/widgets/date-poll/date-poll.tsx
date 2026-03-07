"use client";

import * as React from "react";
import { useTranslation } from "react-i18next";
import { Check, ThumbsDown, ThumbsUp, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { cn } from "@/lib/utils";
import { backend } from "@/modules/backend/backend";
import type { DateProposal } from "@/modules/backend/types";
import { formatDateTimeLong, formatDateTimeRange } from "@/lib/date";
import styles from "./date-poll.module.css";

export type DatePollProps = {
  locale: string;
  storySlug: string;
  canPropose: boolean;
  canVote: boolean;
  canEdit: boolean;
  initialProposals: DateProposal[] | null;
};

export function DatePoll(props: DatePollProps) {
  const { t } = useTranslation();
  const [proposals, setProposals] = React.useState<DateProposal[]>(
    props.initialProposals ?? [],
  );
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [datetimeStart, setDatetimeStart] = React.useState("");
  const [datetimeEnd, setDatetimeEnd] = React.useState("");

  const handleVote = async (proposalId: string, direction: number) => {
    if (!props.canVote || isSubmitting) {
      return;
    }

    setIsSubmitting(true);
    try {
      const result = await backend.voteDateProposal(
        props.locale,
        props.storySlug,
        proposalId,
        direction,
      );

      if (result !== null) {
        setProposals((prev) =>
          prev.map((p) => {
            if (p.id !== proposalId) {
              return p;
            }

            const oldDir = p.viewer_vote_direction;
            const newDir = result.viewer_vote_direction;
            const upDelta = (newDir === 1 ? 1 : 0) - (oldDir === 1 ? 1 : 0);
            const downDelta = (newDir === -1 ? 1 : 0) - (oldDir === -1 ? 1 : 0);

            return {
              ...p,
              vote_score: result.vote_score,
              upvote_count: p.upvote_count + upDelta,
              downvote_count: p.downvote_count + downDelta,
              viewer_vote_direction: result.viewer_vote_direction,
            };
          })
        );
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const handlePropose = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!props.canPropose || isSubmitting || datetimeStart === "") {
      return;
    }

    setIsSubmitting(true);
    try {
      const startISO = new Date(datetimeStart).toISOString();
      const endISO = datetimeEnd !== "" ? new Date(datetimeEnd).toISOString() : undefined;

      const proposal = await backend.createDateProposal(
        props.locale,
        props.storySlug,
        startISO,
        endISO,
      );

      if (proposal !== null) {
        // Refresh the full list to get profile info
        const updated = await backend.getDateProposals(props.locale, props.storySlug);
        if (updated !== null) {
          setProposals(updated.proposals);
        }

        setDatetimeStart("");
        setDatetimeEnd("");
        toast.success(t("Activities.Proposal created successfully"));
      }
    } catch {
      toast.error(t("Activities.Failed to create proposal"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleRemove = async (proposalId: string) => {
    if (isSubmitting) {
      return;
    }

    if (!globalThis.confirm(t("Activities.Confirm remove proposal"))) {
      return;
    }

    setIsSubmitting(true);
    try {
      await backend.removeDateProposal(props.locale, props.storySlug, proposalId);
      setProposals((prev) => prev.filter((p) => p.id !== proposalId));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleFinalize = async (proposalId: string) => {
    if (isSubmitting) {
      return;
    }

    if (!globalThis.confirm(t("Activities.Confirm finalize"))) {
      return;
    }

    setIsSubmitting(true);
    try {
      await backend.finalizeDateProposal(props.locale, props.storySlug, proposalId);
      toast.success(t("Activities.Proposal finalized successfully"));

      // Refresh proposals to see the finalized state
      const updated = await backend.getDateProposals(props.locale, props.storySlug);
      if (updated !== null) {
        setProposals(updated.proposals);
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const hasFinalized = proposals.some((p) => p.is_finalized);

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h3 className={styles.title}>{t("Activities.Date Proposals")}</h3>
      </div>

      {proposals.length === 0
        ? (
          <p className={styles.emptyState}>
            {t("Activities.No proposals yet")}
          </p>
        )
        : (
          <div className={styles.proposalList}>
            {proposals.map((proposal) => (
              <div key={proposal.id} className={styles.proposalCard}>
                <div className={styles.proposalInfo}>
                  <span className={styles.proposalDate}>
                    {proposal.datetime_end !== null
                      ? formatDateTimeRange(
                        new Date(proposal.datetime_start),
                        new Date(proposal.datetime_end),
                        props.locale,
                      )
                      : formatDateTimeLong(new Date(proposal.datetime_start), props.locale)}
                  </span>
                  <span className={styles.proposerInfo}>
                    {proposal.proposer_profile_picture_uri !== null && (
                      <img
                        src={proposal.proposer_profile_picture_uri}
                        alt=""
                        className={styles.proposerAvatar}
                      />
                    )}
                    {proposal.proposer_profile_title}
                  </span>
                </div>

                <div className={styles.voteSection}>
                  {proposal.is_finalized
                    ? (
                      <span className={styles.finalizedBadge}>
                        <Check className="size-3" />
                        {t("Activities.Finalized")}
                      </span>
                    )
                    : (
                      <>
                        {props.canVote && (
                          <>
                            <button
                              type="button"
                              className={cn(
                                styles.voteButton,
                                proposal.viewer_vote_direction === 1 ? styles.voteActive : styles.voteInactive,
                              )}
                              onClick={() => handleVote(proposal.id, 1)}
                              disabled={isSubmitting}
                              aria-label={t("Activities.Agree")}
                            >
                              <ThumbsUp className="size-3" />
                              {proposal.upvote_count}
                            </button>

                            <button
                              type="button"
                              className={cn(
                                styles.voteButton,
                                proposal.viewer_vote_direction === -1 ? styles.voteActive : styles.voteInactive,
                              )}
                              onClick={() => handleVote(proposal.id, -1)}
                              disabled={isSubmitting}
                              aria-label={t("Activities.Disagree")}
                            >
                              <ThumbsDown className="size-3" />
                              {proposal.downvote_count}
                            </button>
                          </>
                        )}

                        <span className={styles.voteScore}>
                          {proposal.vote_score > 0 ? "+" : ""}
                          {proposal.vote_score}
                        </span>

                        {props.canEdit && !hasFinalized && (
                          <>
                            <button
                              type="button"
                              className={styles.finalizeButton}
                              onClick={() => handleFinalize(proposal.id)}
                              disabled={isSubmitting}
                            >
                              {t("Activities.Finalize This Date")}
                            </button>
                            <button
                              type="button"
                              className={styles.removeButton}
                              onClick={() => handleRemove(proposal.id)}
                              disabled={isSubmitting}
                              aria-label={t("Common.Delete")}
                            >
                              <Trash2 className="size-3.5" />
                            </button>
                          </>
                        )}
                      </>
                    )}
                </div>
              </div>
            ))}
          </div>
        )}

      {/* Proposal form — only show if viewer can propose and not yet finalized */}
      {props.canPropose && !hasFinalized && (
        <form className={styles.proposalForm} onSubmit={handlePropose}>
          <p className={styles.formTitle}>{t("Activities.Propose a Date")}</p>
          <div className={styles.formFields}>
            <div className={styles.formField}>
              <label htmlFor="dp-start" className={styles.formLabel}>
                {t("Activities.Proposed Start")}
              </label>
              <input
                id="dp-start"
                type="datetime-local"
                className={styles.formInput}
                value={datetimeStart}
                onChange={(e) => setDatetimeStart(e.target.value)}
                required
              />
            </div>
            <div className={styles.formField}>
              <label htmlFor="dp-end" className={styles.formLabel}>
                {t("Activities.Proposed End (optional)")}
              </label>
              <input
                id="dp-end"
                type="datetime-local"
                className={styles.formInput}
                value={datetimeEnd}
                onChange={(e) => setDatetimeEnd(e.target.value)}
              />
            </div>
            <button
              type="submit"
              className={styles.submitButton}
              disabled={isSubmitting || datetimeStart === ""}
            >
              {t("Activities.Submit Proposal")}
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
