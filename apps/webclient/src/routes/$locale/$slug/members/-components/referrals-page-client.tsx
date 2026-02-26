"use client";

import * as React from "react";
import { useTranslation } from "react-i18next";
import { ChevronDown, ChevronUp, Plus, X } from "lucide-react";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland/site-avatar";
import { backend } from "@/modules/backend/backend";
import type { ProfileMembershipReferral, ProfileTeam, ReferralVote } from "@/modules/backend/types";
import styles from "./referrals-page-client.module.css";

type ReferralsPageClientProps = {
  referrals: ProfileMembershipReferral[];
  teams: ProfileTeam[];
  locale: string;
  slug: string;
};

const VOTE_LABELS = [
  "Referrals.Strongly Disagree",
  "Referrals.Disagree",
  "Referrals.Neutral",
  "Referrals.Agree",
  "Referrals.Strongly Agree",
] as const;

const SCORE_COLORS = [
  "bg-red-500",
  "bg-orange-400",
  "bg-yellow-400",
  "bg-lime-500",
  "bg-green-500",
] as const;

export function ReferralsPageClient(props: ReferralsPageClientProps) {
  const { t } = useTranslation();
  const [referrals, setReferrals] = React.useState(props.referrals);
  const [showCreateDialog, setShowCreateDialog] = React.useState(false);

  const handleReferralCreated = React.useCallback(
    async () => {
      setShowCreateDialog(false);
      // Re-fetch the full list to get complete data with joined profiles
      const updated = await backend.listReferrals(props.locale, props.slug);
      if (updated !== null) {
        setReferrals(updated);
      }
    },
    [props.locale, props.slug],
  );

  const handleVoteUpdated = React.useCallback(
    (referralId: string, score: number) => {
      setReferrals((prev) =>
        prev.map((r) => {
          if (r.id !== referralId) return r;
          return {
            ...r,
            viewer_vote_score: score as 1 | 2 | 3 | 4 | 5,
          };
        })
      );
    },
    [],
  );

  return (
    <>
      <div className={styles.header}>
        <div className={styles.headerText}>
          <h2>{t("Layout.Referrals")}</h2>
          <p>{t("Referrals.Referral proposals for new members.")}</p>
        </div>
        <button
          type="button"
          className={styles.referButton}
          onClick={() => setShowCreateDialog(true)}
        >
          <span className="flex items-center gap-1.5">
            <Plus className="size-4" />
            {t("Referrals.Refer Someone")}
          </span>
        </button>
      </div>

      {referrals.length === 0
        ? (
          <div className={styles.emptyState}>
            <p className={styles.emptyStateText}>
              {t("Referrals.No referrals yet")}
            </p>
          </div>
        )
        : (
          <div className="flex flex-col gap-4">
            {referrals.map((referral) => (
              <ReferralCard
                key={referral.id}
                referral={referral}
                locale={props.locale}
                slug={props.slug}
                onVoteUpdated={handleVoteUpdated}
              />
            ))}
          </div>
        )}

      {showCreateDialog && (
        <CreateReferralDialog
          locale={props.locale}
          slug={props.slug}
          teams={props.teams}
          onCreated={handleReferralCreated}
          onClose={() => setShowCreateDialog(false)}
        />
      )}
    </>
  );
}

// ─── Create Referral Dialog ──────────────────────────────────────────

type CreateReferralDialogProps = {
  locale: string;
  slug: string;
  teams: ProfileTeam[];
  onCreated: () => void;
  onClose: () => void;
};

function CreateReferralDialog(props: CreateReferralDialogProps) {
  const { t } = useTranslation();
  const [username, setUsername] = React.useState("");
  const [selectedTeamIds, setSelectedTeamIds] = React.useState<string[]>([]);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const handleTeamToggle = React.useCallback((teamId: string) => {
    setSelectedTeamIds((prev) => prev.includes(teamId) ? prev.filter((id) => id !== teamId) : [...prev, teamId]);
  }, []);

  const handleSubmit = React.useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      const trimmedUsername = username.trim();
      if (trimmedUsername.length === 0) return;

      setIsSubmitting(true);
      setError(null);

      const result = await backend.createReferral(
        props.locale,
        props.slug,
        trimmedUsername,
        selectedTeamIds,
      );

      if (result === null) {
        setError(t("Referrals.Failed to create referral. Please check the username and try again."));
        setIsSubmitting(false);
        return;
      }

      props.onCreated();
    },
    [username, selectedTeamIds, props.locale, props.slug, props.onCreated, t],
  );

  return (
    <div className={styles.dialogOverlay} onClick={props.onClose}>
      <div
        className={styles.dialogContent}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between">
          <h3 className={styles.dialogTitle}>
            {t("Referrals.Refer Someone")}
          </h3>
          <button type="button" onClick={props.onClose}>
            <X className="size-4 text-muted-foreground hover:text-foreground" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className={styles.formGroup}>
            <label htmlFor="referral-username" className={styles.label}>
              {t("Referrals.Username")}
            </label>
            <input
              id="referral-username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder={t("Referrals.Enter profile slug")}
              className={styles.input}
              autoFocus
            />
          </div>

          {props.teams.length > 0 && (
            <div className={styles.formGroup}>
              <label className={styles.label}>
                {t("Referrals.Suggested Teams")}
              </label>
              <div className={styles.teamCheckboxes}>
                {props.teams.map((team) => (
                  <label key={team.id} className={styles.teamCheckbox}>
                    <input
                      type="checkbox"
                      checked={selectedTeamIds.includes(team.id)}
                      onChange={() => handleTeamToggle(team.id)}
                    />
                    <span>{team.name}</span>
                  </label>
                ))}
              </div>
            </div>
          )}

          {error !== null && <p className={styles.errorMessage}>{error}</p>}

          <div className={styles.dialogActions}>
            <button
              type="button"
              onClick={props.onClose}
              className={styles.cancelButton}
            >
              {t("Common.Cancel")}
            </button>
            <button
              type="submit"
              disabled={isSubmitting || username.trim().length === 0}
              className={styles.submitButton}
            >
              {t("Referrals.Submit Referral")}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// ─── Referral Card ───────────────────────────────────────────────────

type ReferralCardProps = {
  referral: ProfileMembershipReferral;
  locale: string;
  slug: string;
  onVoteUpdated: (referralId: string, score: number) => void;
};

function ReferralCard(props: ReferralCardProps) {
  const { t } = useTranslation();
  const initialScore = props.referral.viewer_vote_score !== 0
    ? props.referral.viewer_vote_score
    : null;
  const [viewerScore, setViewerScore] = React.useState<number | null>(
    initialScore ?? null,
  );
  const [comment, setComment] = React.useState(
    props.referral.viewer_vote_comment ?? "",
  );
  const [showCommentInput, setShowCommentInput] = React.useState(
    initialScore !== null,
  );
  const [showVotes, setShowVotes] = React.useState(false);
  const [votes, setVotes] = React.useState<ReferralVote[] | null>(null);
  const [isVoting, setIsVoting] = React.useState(false);
  const [isLoadingVotes, setIsLoadingVotes] = React.useState(false);

  const referred = props.referral.referred_profile;
  const referrer = props.referral.referrer_profile;

  const formattedDate = new Date(props.referral.created_at).toLocaleDateString(
    props.locale,
    { year: "numeric", month: "short", day: "numeric" },
  );

  const handleVote = React.useCallback(
    async (score: number) => {
      if (isVoting) return;

      setIsVoting(true);
      setViewerScore(score);
      setShowCommentInput(true);

      try {
        const result = await backend.voteReferral(
          props.locale,
          props.slug,
          props.referral.id,
          score,
          comment.trim().length > 0 ? comment.trim() : null,
        );

        if (result !== null) {
          props.onVoteUpdated(props.referral.id, score);
        }
      } finally {
        setIsVoting(false);
      }
    },
    [
      isVoting,
      comment,
      props.locale,
      props.slug,
      props.referral.id,
      props.onVoteUpdated,
    ],
  );

  const handleCommentSubmit = React.useCallback(async () => {
    if (viewerScore === null || viewerScore === 0) return;

    setIsVoting(true);
    try {
      const result = await backend.voteReferral(
        props.locale,
        props.slug,
        props.referral.id,
        viewerScore,
        comment.trim().length > 0 ? comment.trim() : null,
      );

      if (result !== null) {
        setShowCommentInput(false);
      }
    } finally {
      setIsVoting(false);
    }
  }, [
    viewerScore,
    comment,
    props.locale,
    props.slug,
    props.referral.id,
  ]);

  const handleToggleVotes = React.useCallback(async () => {
    if (showVotes) {
      setShowVotes(false);
      return;
    }

    if (votes === null) {
      setIsLoadingVotes(true);
      const result = await backend.getReferralVotes(
        props.locale,
        props.slug,
        props.referral.id,
      );
      setVotes(result ?? []);
      setIsLoadingVotes(false);
    }

    setShowVotes(true);
  }, [
    showVotes,
    votes,
    props.locale,
    props.slug,
    props.referral.id,
  ]);

  return (
    <div className={styles.referralCard}>
      {/* Header: referred profile + status */}
      <div className={styles.referralHeader}>
        <div className={styles.referralProfileInfo}>
          {referred !== undefined && referred !== null && (
            <>
              <LocaleLink to={`/${referred.slug}`}>
                <SiteAvatar
                  src={referred.profile_picture_uri}
                  name={referred.title}
                  fallbackName={referred.slug}
                  className={styles.referralAvatar}
                />
              </LocaleLink>
              <div>
                <LocaleLink
                  to={`/${referred.slug}`}
                  className={styles.referralName}
                >
                  {referred.title}
                </LocaleLink>
                <div className={styles.referralSlug}>@{referred.slug}</div>
              </div>
            </>
          )}
        </div>
        <span className={styles.statusBadge}>
          {t(`Referrals.Status.${props.referral.status}`)}
        </span>
      </div>

      {/* Referrer info */}
      {referrer !== undefined && referrer !== null && (
        <div className={styles.referrerInfo}>
          {t("Referrals.Referred by")}{" "}
          <LocaleLink
            to={`/${referrer.slug}`}
            className="hover:underline"
          >
            {referrer.title}
          </LocaleLink>{" "}
          &middot; {formattedDate}
        </div>
      )}

      {/* Team badges */}
      {props.referral.teams !== undefined &&
        props.referral.teams !== null &&
        props.referral.teams.length > 0 && (
        <div className={styles.teamBadges}>
          {props.referral.teams.map((team) => (
            <span key={team.id} className={styles.teamBadge}>
              {team.name}
            </span>
          ))}
        </div>
      )}

      {/* Vote section */}
      <div className={styles.voteSection}>
        {/* Vote summary */}
        <div className={styles.voteSummary}>
          <span>
            {props.referral.total_votes} {t("Referrals.votes")}
          </span>
          {props.referral.total_votes > 0 && (
            <>
              <span>&middot;</span>
              <span>
                {t("Referrals.average")}: {props.referral.average_score.toFixed(1)}
              </span>
              <div className={styles.scoreBreakdown}>
                {SCORE_COLORS.map((color, i) => {
                  const score = i + 1;
                  const count = props.referral.vote_breakdown?.[score] ?? 0;
                  const pct = props.referral.total_votes > 0 ? (count / props.referral.total_votes) * 100 : 0;
                  return (
                    <div
                      key={score}
                      className={`${styles.scoreSegment} ${color}`}
                      style={{ width: `${pct}%` }}
                      title={`${VOTE_LABELS[i]}: ${count}`}
                    />
                  );
                })}
              </div>
            </>
          )}
        </div>

        {/* Vote buttons */}
        <div className={styles.voteButtons}>
          {VOTE_LABELS.map((label, i) => {
            const score = i + 1;
            const isActive = viewerScore === score;
            return (
              <button
                key={score}
                type="button"
                className={isActive ? styles.voteButtonActive : styles.voteButton}
                onClick={() => handleVote(score)}
                disabled={isVoting}
                title={t(label)}
              >
                {t(label)}
              </button>
            );
          })}
        </div>

        {/* Comment input */}
        {showCommentInput && viewerScore !== null && (
          <div className={styles.commentSection}>
            <textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder={t("Referrals.Add a comment (optional)")}
              className={styles.commentTextarea}
              rows={2}
            />
            <button
              type="button"
              onClick={handleCommentSubmit}
              disabled={isVoting}
              className={styles.commentSubmit}
            >
              {t("Referrals.Save comment")}
            </button>
          </div>
        )}

        {/* View votes toggle */}
        {props.referral.total_votes > 0 && (
          <div className={styles.voteDetails}>
            <button
              type="button"
              onClick={handleToggleVotes}
              className={styles.voteDetailsToggle}
            >
              <span className="flex items-center gap-1">
                {showVotes ? <ChevronUp className="size-3" /> : <ChevronDown className="size-3" />}
                {t("Referrals.View All Votes")} ({props.referral.total_votes})
              </span>
            </button>

            {showVotes && votes !== null && (
              <div className={styles.voteList}>
                {votes.map((vote) => (
                  <div key={vote.id} className={styles.voteItem}>
                    {vote.voter_profile !== undefined &&
                      vote.voter_profile !== null && (
                      <LocaleLink to={`/${vote.voter_profile.slug}`}>
                        <SiteAvatar
                          src={vote.voter_profile.profile_picture_uri}
                          name={vote.voter_profile.title}
                          fallbackName={vote.voter_profile.slug}
                          className={styles.voteItemAvatar}
                        />
                      </LocaleLink>
                    )}
                    <div className={styles.voteItemContent}>
                      <div className={styles.voteItemHeader}>
                        {vote.voter_profile !== undefined &&
                          vote.voter_profile !== null && (
                          <LocaleLink
                            to={`/${vote.voter_profile.slug}`}
                            className={styles.voteItemName}
                          >
                            {vote.voter_profile.title}
                          </LocaleLink>
                        )}
                        <span className={styles.voteItemScore}>
                          {t(VOTE_LABELS[vote.score - 1])}
                        </span>
                      </div>
                      {vote.comment !== null &&
                        vote.comment !== undefined && (
                        <p className={styles.voteItemComment}>
                          {vote.comment}
                        </p>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}

            {showVotes && isLoadingVotes && (
              <p className="text-xs text-muted-foreground mt-2">
                {t("Common.Loading")}...
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
