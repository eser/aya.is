import * as React from "react";
import { useTranslation } from "react-i18next";
import { useRouter } from "@tanstack/react-router";
import { ChevronDown, ChevronUp, MoreHorizontal, Plus, X } from "lucide-react";
import { toast } from "sonner";
import { useAuth } from "@/lib/auth/auth-context";
import { LocaleLink } from "@/components/locale-link";
import { SiteAvatar } from "@/components/userland/site-avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { backend } from "@/modules/backend/backend";
import type {
  CandidateFormResponse,
  CandidateStatus,
  CandidateVote,
  ProfileMembershipCandidate,
  ProfileTeam,
} from "@/modules/backend/types";
import styles from "./candidates-page-client.module.css";

type CandidatesPageClientProps = {
  candidates: ProfileMembershipCandidate[];
  locale: string;
  slug: string;
  viewerMembershipKind: string | null;
};

const LEAD_PLUS_KINDS = new Set(["maintainer", "lead", "owner"]);

const VOTE_LABELS = [
  "Candidates.Strongly Disagree",
  "Candidates.Disagree",
  "Candidates.Neutral",
  "Candidates.Agree",
  "Candidates.Strongly Agree",
] as const;

export function CandidatesPageClient(props: CandidatesPageClientProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const { user } = useAuth();
  const [showCreateDialog, setShowCreateDialog] = React.useState(false);

  const teams = React.useMemo(() => {
    const match = user?.accessible_profiles?.find(
      (p) => p.slug === props.slug,
    );
    return match?.teams ?? [];
  }, [user?.accessible_profiles, props.slug]);

  const handleCandidateCreated = React.useCallback(
    () => {
      setShowCreateDialog(false);
      router.invalidate();
    },
    [router],
  );

  const handleStatusChange = React.useCallback(
    async (candidateId: string, status: CandidateStatus) => {
      const result = await backend.updateCandidateStatus(
        props.locale,
        props.slug,
        candidateId,
        status,
      );

      if (result) {
        toast.success(t("Candidates.Actions.StatusUpdated"));
        router.invalidate();
      } else {
        toast.error(t("Candidates.Actions.StatusUpdateFailed"));
      }
    },
    [props.locale, props.slug, router, t],
  );

  const isLeadPlus = props.viewerMembershipKind !== null &&
    LEAD_PLUS_KINDS.has(props.viewerMembershipKind);

  return (
    <>
      <div className={styles.header}>
        <div className={styles.headerText}>
          <h2>{t("Layout.Candidates")}</h2>
          <p>{t("Candidates.Candidate proposals for new members.")}</p>
        </div>
        <button
          type="button"
          className={styles.referButton}
          onClick={() => setShowCreateDialog(true)}
        >
          <span className="flex items-center gap-1.5">
            <Plus className="size-4" />
            {t("Candidates.Refer Someone")}
          </span>
        </button>
      </div>

      {props.candidates.length === 0
        ? (
          <div className={styles.emptyState}>
            <p className={styles.emptyStateText}>
              {t("Candidates.No candidates yet")}
            </p>
          </div>
        )
        : (
          <div className="flex flex-col gap-4">
            {props.candidates.map((candidate) => (
              <CandidateCard
                key={candidate.id}
                candidate={candidate}
                locale={props.locale}
                slug={props.slug}
                isLeadPlus={isLeadPlus}
                onStatusChange={handleStatusChange}
              />
            ))}
          </div>
        )}

      {showCreateDialog && (
        <CreateCandidateDialog
          locale={props.locale}
          slug={props.slug}
          teams={teams}
          onCreated={handleCandidateCreated}
          onClose={() => setShowCreateDialog(false)}
        />
      )}
    </>
  );
}

// ─── Create Candidate Dialog ──────────────────────────────────────────

type CreateCandidateDialogProps = {
  locale: string;
  slug: string;
  teams: ProfileTeam[];
  onCreated: () => void;
  onClose: () => void;
};

function CreateCandidateDialog(props: CreateCandidateDialogProps) {
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

      const result = await backend.createCandidate(
        props.locale,
        props.slug,
        trimmedUsername,
        selectedTeamIds,
      );

      if (result === null) {
        setError(t("Candidates.Failed to create candidate. Please check the username and try again."));
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
            {t("Candidates.Refer Someone")}
          </h3>
          <button type="button" onClick={props.onClose}>
            <X className="size-4 text-muted-foreground hover:text-foreground" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className={styles.formGroup}>
            <label htmlFor="candidate-username" className={styles.label}>
              {t("Candidates.Username")}
            </label>
            <input
              id="candidate-username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder={t("Candidates.Enter profile slug")}
              className={styles.input}
              autoFocus
            />
          </div>

          {props.teams.length > 0 && (
            <div className={styles.formGroup}>
              <label className={styles.label}>
                {t("Candidates.Suggested Teams")}
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
              {t("Candidates.Submit Candidate")}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// ─── Candidate Card ───────────────────────────────────────────────────

type CandidateCardProps = {
  candidate: ProfileMembershipCandidate;
  locale: string;
  slug: string;
  isLeadPlus: boolean;
  onStatusChange: (candidateId: string, status: CandidateStatus) => Promise<void>;
};

type ConfirmAction = {
  status: CandidateStatus;
  titleKey: string;
  descriptionKey: string;
  descriptionParams?: Record<string, string>;
};

function CandidateCard(props: CandidateCardProps) {
  const { t } = useTranslation();
  const router = useRouter();
  const initialScore = props.candidate.viewer_vote_score !== -1 ? props.candidate.viewer_vote_score : null;
  const [viewerScore, setViewerScore] = React.useState<number | null>(
    initialScore ?? null,
  );
  const [showVotes, setShowVotes] = React.useState(false);
  const [votes, setVotes] = React.useState<CandidateVote[] | null>(null);
  const [isLoadingVotes, setIsLoadingVotes] = React.useState(false);
  const [totalVotes, setTotalVotes] = React.useState(props.candidate.total_votes);
  const [averageScore, setAverageScore] = React.useState(props.candidate.average_score);
  const [comment, setComment] = React.useState(props.candidate.viewer_vote_comment ?? "");
  const [confirmAction, setConfirmAction] = React.useState<ConfirmAction | null>(null);
  const [isUpdatingStatus, setIsUpdatingStatus] = React.useState(false);

  const handleConfirmAction = React.useCallback(async () => {
    if (confirmAction === null) return;

    setIsUpdatingStatus(true);
    await props.onStatusChange(props.candidate.id, confirmAction.status);
    setIsUpdatingStatus(false);
    setConfirmAction(null);
  }, [confirmAction, props.onStatusChange, props.candidate.id]);

  const referred = props.candidate.referred_profile;
  const referrer = props.candidate.referrer_profile;

  const formattedDate = new Date(props.candidate.created_at).toLocaleDateString(
    props.locale,
    { year: "numeric", month: "short", day: "numeric" },
  );

  const refreshVotes = React.useCallback(async () => {
    const freshVotes = await backend.getCandidateVotes(
      props.locale,
      props.slug,
      props.candidate.id,
    );

    if (freshVotes !== null) {
      setVotes(freshVotes);
      setTotalVotes(freshVotes.length);

      if (freshVotes.length > 0) {
        const avg = freshVotes.reduce((sum, v) => sum + v.score, 0) / freshVotes.length;
        setAverageScore(Math.round(avg * 10) / 10);
      }
    }
  }, [props.locale, props.slug, props.candidate.id]);

  const handleReset = React.useCallback(() => {
    setViewerScore(initialScore ?? null);
    setComment(props.candidate.viewer_vote_comment ?? "");
  }, [initialScore, props.candidate.viewer_vote_comment]);

  const isDirty = viewerScore !== (initialScore ?? null) ||
    comment !== (props.candidate.viewer_vote_comment ?? "");

  const [, saveAction, isSaving] = React.useActionState(
    async (_prev: null): Promise<null> => {
      if (viewerScore === null) return null;

      const trimmedComment = comment.trim();
      const result = await backend.voteCandidate(
        props.locale,
        props.slug,
        props.candidate.id,
        viewerScore,
        trimmedComment.length > 0 ? trimmedComment : null,
      );

      if (result !== null) {
        toast.success(t("Candidates.Comment saved"));
        await refreshVotes();
        router.invalidate();
      } else {
        toast.error(t("Candidates.Failed to submit vote"));
      }

      return null;
    },
    null,
  );

  const handleToggleVotes = React.useCallback(async () => {
    if (showVotes) {
      setShowVotes(false);
      return;
    }

    if (votes === null) {
      setIsLoadingVotes(true);
      const result = await backend.getCandidateVotes(
        props.locale,
        props.slug,
        props.candidate.id,
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
    props.candidate.id,
  ]);

  return (
    <div className={styles.candidateCard}>
      {/* Header: referred profile + status */}
      <div className={styles.candidateHeader}>
        <div className={styles.candidateProfileInfo}>
          {referred !== undefined && referred !== null && (
            <>
              <LocaleLink to={`/${referred.slug}`}>
                <SiteAvatar
                  src={referred.profile_picture_uri}
                  name={referred.title}
                  fallbackName={referred.slug}
                  className={styles.candidateAvatar}
                />
              </LocaleLink>
              <div>
                <LocaleLink
                  to={`/${referred.slug}`}
                  className={styles.candidateName}
                >
                  {referred.title}
                </LocaleLink>
                <div className={styles.candidateSlug}>@{referred.slug}</div>
              </div>
            </>
          )}
        </div>
        <div className={styles.statusActions}>
          <span className={styles.statusBadge}>
            {t(`Candidates.Status.${props.candidate.status}`)}
          </span>
          <CandidateActionsMenu
            candidate={props.candidate}
            isLeadPlus={props.isLeadPlus}
            onAction={setConfirmAction}
          />
        </div>
      </div>

      {/* Referrer / applicant info */}
      {props.candidate.source === "application"
        ? (
          <div className={styles.referrerInfo}>
            {t("Candidates.Applied on")} {formattedDate}
          </div>
        )
        : referrer !== undefined && referrer !== null && (
          <div className={styles.referrerInfo}>
            {t("Candidates.Referred by")}{" "}
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
      {props.candidate.teams !== undefined &&
        props.candidate.teams !== null &&
        props.candidate.teams.length > 0 && (
        <div className={styles.teamBadges}>
          {props.candidate.teams.map((team) => (
            <span key={team.id} className={styles.teamBadge}>
              {team.name}
            </span>
          ))}
        </div>
      )}

      {/* Form responses viewer — for application-type candidates */}
      {props.candidate.source === "application" && props.isLeadPlus && (
        <FormResponsesViewer
          locale={props.locale}
          slug={props.slug}
          candidateId={props.candidate.id}
          applicantMessage={props.candidate.applicant_message ?? null}
        />
      )}

      {/* Vote section — only shown when candidate is in "voting" state */}
      {props.candidate.status === "voting" && (
        <form action={saveAction} className={styles.voteSection}>
          {/* Vote buttons (local selection only, no auto-submit) */}
          <div className={styles.voteButtons}>
            {VOTE_LABELS.map((label, i) => {
              const score = i;
              const isActive = viewerScore === score;
              return (
                <button
                  key={score}
                  type="button"
                  className={isActive ? styles.voteButtonActive : styles.voteButton}
                  onClick={() => setViewerScore(score)}
                  disabled={isSaving}
                  title={t(label)}
                >
                  {t(label)}
                </button>
              );
            })}
          </div>

          {/* Comment */}
          <textarea
            name="comment"
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            placeholder={t("Candidates.Add a comment (optional)")}
            className={styles.commentTextarea}
            rows={2}
          />

          {/* Form actions */}
          <div className={styles.formActions}>
            <button
              type="button"
              onClick={handleReset}
              disabled={isSaving || !isDirty}
              className={styles.resetButton}
            >
              {t("Common.Reset")}
            </button>
            <button
              type="submit"
              disabled={isSaving || viewerScore === null || !isDirty}
              className={styles.commentSubmit}
            >
              {t("Common.Save")}
            </button>
          </div>
        </form>
      )}

      {/* View all votes toggle */}
      {totalVotes > 0 && (
        <div className={styles.voteDetails}>
          <button
            type="button"
            onClick={handleToggleVotes}
            className={styles.voteDetailsToggle}
          >
            <span className="flex items-center gap-1">
              {showVotes ? <ChevronUp className="size-3" /> : <ChevronDown className="size-3" />}
              {t("Candidates.View All")} ({totalVotes} {t("Candidates.votes")}, {averageScore.toFixed(1)}{" "}
              {t("Candidates.average")})
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
                        {t(VOTE_LABELS[vote.score])}
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

      {/* Confirmation dialog for status changes */}
      <AlertDialog
        open={confirmAction !== null}
        onOpenChange={(open) => {
          if (!open) setConfirmAction(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {confirmAction !== null ? t(confirmAction.titleKey) : ""}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {confirmAction !== null ? t(confirmAction.descriptionKey, confirmAction.descriptionParams) : ""}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>
              {t("Candidates.Actions.Cancel")}
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleConfirmAction}
              disabled={isUpdatingStatus}
            >
              {t("Candidates.Actions.Confirm")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}

// ─── Candidate Actions Menu ──────────────────────────────────────────

type CandidateActionsMenuProps = {
  candidate: ProfileMembershipCandidate;
  isLeadPlus: boolean;
  onAction: (action: ConfirmAction) => void;
};

function CandidateActionsMenu(props: CandidateActionsMenuProps) {
  const { t } = useTranslation();
  const { status } = props.candidate;

  // Only show menu for lead+ roles and non-terminal statuses.
  if (!props.isLeadPlus) return null;
  if (
    status === "invitation_pending_response" ||
    status === "reference_rejected" ||
    status === "invitation_accepted" ||
    status === "invitation_rejected" ||
    status === "application_accepted"
  ) {
    return null;
  }

  const referredName = props.candidate.referred_profile?.title ?? props.candidate.referred_profile?.slug ?? "";

  return (
    <DropdownMenu>
      <DropdownMenuTrigger className={styles.actionsMenuTrigger}>
        <MoreHorizontal className="size-4" />
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {status === "voting" && (
          <>
            <DropdownMenuItem
              onClick={() =>
                props.onAction({
                  status: "frozen",
                  titleKey: "Candidates.Actions.Freeze",
                  descriptionKey: "Candidates.Actions.FreezeConfirm",
                })}
            >
              {t("Candidates.Actions.Freeze")}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            {props.candidate.source === "application"
              ? (
                <DropdownMenuItem
                  onClick={() =>
                    props.onAction({
                      status: "application_accepted",
                      titleKey: "Candidates.Actions.AcceptApplication",
                      descriptionKey: "Candidates.Actions.AcceptApplicationConfirm",
                      descriptionParams: { name: referredName },
                    })}
                >
                  {t("Candidates.Actions.AcceptApplication")}
                </DropdownMenuItem>
              )
              : (
                <DropdownMenuItem
                  onClick={() =>
                    props.onAction({
                      status: "invitation_pending_response",
                      titleKey: "Candidates.Actions.SendInvite",
                      descriptionKey: "Candidates.Actions.SendInviteConfirm",
                      descriptionParams: { name: referredName },
                    })}
                >
                  {t("Candidates.Actions.SendInvite")}
                </DropdownMenuItem>
              )}
            <DropdownMenuItem
              variant="destructive"
              onClick={() =>
                props.onAction({
                  status: "reference_rejected",
                  titleKey: "Candidates.Actions.Reject",
                  descriptionKey: "Candidates.Actions.RejectConfirm",
                })}
            >
              {t("Candidates.Actions.Reject")}
            </DropdownMenuItem>
          </>
        )}

        {status === "frozen" && (
          <>
            <DropdownMenuItem
              onClick={() =>
                props.onAction({
                  status: "voting",
                  titleKey: "Candidates.Actions.Unfreeze",
                  descriptionKey: "Candidates.Actions.UnfreezeConfirm",
                })}
            >
              {t("Candidates.Actions.Unfreeze")}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              variant="destructive"
              onClick={() =>
                props.onAction({
                  status: "reference_rejected",
                  titleKey: "Candidates.Actions.Reject",
                  descriptionKey: "Candidates.Actions.RejectConfirm",
                })}
            >
              {t("Candidates.Actions.Reject")}
            </DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

// ─── Form Responses Viewer ──────────────────────────────────────────

type FormResponsesViewerProps = {
  locale: string;
  slug: string;
  candidateId: string;
  applicantMessage: string | null;
};

function FormResponsesViewer(props: FormResponsesViewerProps) {
  const { t } = useTranslation();
  const [expanded, setExpanded] = React.useState(false);
  const [responses, setResponses] = React.useState<CandidateFormResponse[] | null>(null);
  const [isLoading, setIsLoading] = React.useState(false);

  const handleToggle = React.useCallback(async () => {
    if (expanded) {
      setExpanded(false);
      return;
    }

    if (responses === null) {
      setIsLoading(true);
      const result = await backend.getCandidateResponses(
        props.locale,
        props.slug,
        props.candidateId,
      );
      setResponses(result ?? []);
      setIsLoading(false);
    }

    setExpanded(true);
  }, [expanded, responses, props.locale, props.slug, props.candidateId]);

  return (
    <div className={styles.formResponses}>
      <button
        type="button"
        onClick={handleToggle}
        className={styles.formResponsesToggle}
      >
        <span className="flex items-center gap-1">
          {expanded ? <ChevronUp className="size-3" /> : <ChevronDown className="size-3" />}
          {t("Candidates.Form Responses")}
        </span>
      </button>

      {expanded && isLoading && (
        <p className="text-xs text-muted-foreground mt-2">
          {t("Common.Loading")}...
        </p>
      )}

      {expanded && responses !== null && (responses.length > 0 || props.applicantMessage !== null) && (
        <div className={styles.formResponsesList}>
          {responses.map((response) => (
            <div key={response.id} className={styles.formResponseItem}>
              <div className={styles.formResponseLabel}>
                {t(`ApplicationFields.${response.field_label}`, response.field_label)}
              </div>
              <div className={styles.formResponseValue}>
                {response.value.length > 0 ? response.value : <span className="italic">{t("Common.Empty")}</span>}
              </div>
            </div>
          ))}
          {props.applicantMessage !== null && (
            <div className={styles.formResponseItem}>
              <div className={styles.formResponseLabel}>
                {t("Applications.Additional message")}
              </div>
              <div className={styles.formResponseValue}>
                {props.applicantMessage}
              </div>
            </div>
          )}
        </div>
      )}

      {expanded && responses !== null && responses.length === 0 && props.applicantMessage === null && (
        <p className="text-xs text-muted-foreground mt-2">
          {t("Candidates.No form responses")}
        </p>
      )}
    </div>
  );
}
