// Unified Mailbox — aggregates envelopes from all profiles where user is maintainer+
import * as React from "react";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  Mail,
  Check,
  X,
  Send,
  Loader2,
  CheckCircle,
  ChevronDown,
  ChevronUp,
} from "lucide-react";
import { backend, type ProfileEnvelope, type EnvelopeStatus, type MailboxEnvelope } from "@/modules/backend/backend";
import { PageLayout } from "@/components/page-layouts/default";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
import { formatDateString } from "@/lib/date";
import { useAuth } from "@/lib/auth/auth-context";
import { getCurrentLanguage } from "@/modules/i18n/i18n";
import styles from "./mailbox.module.css";

const ENVELOPE_KINDS = [
  { value: "standard", labelKey: "ProfileSettings.Standard Message" },
  { value: "telegram_group", labelKey: "ProfileSettings.Telegram Group Invite" },
] as const;

const statusColors: Record<EnvelopeStatus, string> = {
  pending: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200",
  accepted: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
  rejected: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
  revoked: "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200",
  redeemed: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
};

export const Route = createFileRoute("/$locale/mailbox")({
  ssr: false,
  component: MailboxPage,
});

function SenderInfo(props: { envelope: ProfileEnvelope; locale: string }) {
  const { t } = useTranslation();

  if (props.envelope.sender_profile_slug !== null && props.envelope.sender_profile_slug !== "") {
    return (
      <Link
        to="/$locale/$slug"
        params={{ locale: props.locale, slug: props.envelope.sender_profile_slug }}
        className="text-sm text-primary hover:underline"
      >
        {props.envelope.sender_profile_title ?? props.envelope.sender_profile_slug}
      </Link>
    );
  }

  if (props.envelope.sender_profile_id !== null) {
    return <span className="text-sm text-muted-foreground">{t("Common.Unknown")}</span>;
  }

  return <span className="text-sm text-muted-foreground italic">{t("Common.System")}</span>;
}

function MailboxPage() {
  const { t } = useTranslation();
  const locale = getCurrentLanguage();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading, user, refreshAuth } = useAuth();

  const [envelopes, setEnvelopes] = React.useState<MailboxEnvelope[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [actionInProgress, setActionInProgress] = React.useState<string | null>(null);
  const [activeFilter, setActiveFilter] = React.useState<string>("all");

  // Reject confirmation dialog state
  const [rejectDialogOpen, setRejectDialogOpen] = React.useState(false);
  const [rejectTargetId, setRejectTargetId] = React.useState<string | null>(null);
  const [rejectTargetSlug, setRejectTargetSlug] = React.useState<string | null>(null);

  // Send form state
  const [sendFormOpen, setSendFormOpen] = React.useState(false);
  const [sendFromSlug, setSendFromSlug] = React.useState("");
  const [envelopeKind, setEnvelopeKind] = React.useState("standard");
  const [targetSlug, setTargetSlug] = React.useState("");
  const [inviteCode, setInviteCode] = React.useState("");
  const [sendTitle, setSendTitle] = React.useState("");
  const [sendDescription, setSendDescription] = React.useState("");
  const [isSending, setIsSending] = React.useState(false);
  const [sendError, setSendError] = React.useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = React.useState<Record<string, string | null>>({});
  const [sendSuccess, setSendSuccess] = React.useState(false);

  // Compute maintainer+ profiles from the user's accessible profiles (for filter bar & send form)
  const maintainerProfiles = React.useMemo(() => {
    if (user === null) {
      return [];
    }

    const profiles: Array<{ slug: string; title: string; kind: string }> = [];

    // Add user's own individual profile
    if (user.individual_profile !== undefined && user.individual_profile !== null) {
      profiles.push({
        slug: user.individual_profile.slug,
        title: user.individual_profile.title,
        kind: user.individual_profile.kind,
      });
    }

    // Add accessible profiles with maintainer+ membership
    const maintainerPlusKinds = new Set(["maintainer", "lead", "owner"]);
    if (user.accessible_profiles !== undefined) {
      for (const profile of user.accessible_profiles) {
        if (maintainerPlusKinds.has(profile.membership_kind)) {
          profiles.push({
            slug: profile.slug,
            title: profile.title,
            kind: profile.kind,
          });
        }
      }
    }

    return profiles;
  }, [user]);

  // Non-individual profiles that can send in-mail
  const sendableProfiles = React.useMemo(() => {
    return maintainerProfiles.filter((p) => p.kind !== "individual");
  }, [maintainerProfiles]);

  // Set default sendFromSlug when sendable profiles load
  React.useEffect(() => {
    if (sendFromSlug === "" && sendableProfiles.length > 0) {
      setSendFromSlug(sendableProfiles[0].slug);
    }
  }, [sendableProfiles, sendFromSlug]);

  const loadEnvelopes = React.useCallback(async () => {
    setIsLoading(true);

    const result = await backend.listMailboxEnvelopes(locale);

    if (result !== null) {
      // Sort by created_at descending (newest first)
      result.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
      setEnvelopes(result);
    } else {
      setEnvelopes([]);
    }

    setIsLoading(false);
  }, [locale]);

  // Redirect if not authenticated once auth is loaded
  React.useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate({ to: `/${locale}` });
    }
  }, [authLoading, isAuthenticated, navigate, locale]);

  // Load envelopes when authenticated
  React.useEffect(() => {
    if (!authLoading && isAuthenticated) {
      loadEnvelopes();
    }
  }, [authLoading, isAuthenticated, loadEnvelopes]);

  const handleAccept = async (envelopeId: string, profileSlug: string) => {
    setActionInProgress(envelopeId);
    const success = await backend.acceptProfileEnvelope(locale, profileSlug, envelopeId);
    if (success) {
      await loadEnvelopes();
      await refreshAuth();
    }
    setActionInProgress(null);
  };

  const promptReject = (envelopeId: string, profileSlug: string) => {
    setRejectTargetId(envelopeId);
    setRejectTargetSlug(profileSlug);
    setRejectDialogOpen(true);
  };

  const handleRejectConfirm = async () => {
    if (rejectTargetId === null || rejectTargetSlug === null) {
      return;
    }

    setRejectDialogOpen(false);
    setActionInProgress(rejectTargetId);

    const success = await backend.rejectProfileEnvelope(locale, rejectTargetSlug, rejectTargetId);
    if (success) {
      await loadEnvelopes();
      await refreshAuth();
    }

    setActionInProgress(null);
    setRejectTargetId(null);
    setRejectTargetSlug(null);
  };

  const handleSend = async () => {
    const errors: Record<string, string | null> = {};
    if (sendFromSlug === "") {
      errors.sendFromSlug = t("Common.This field is required");
    }
    if (targetSlug.trim() === "") {
      errors.targetSlug = t("Common.This field is required");
    }
    if (envelopeKind === "telegram_group" && inviteCode.trim() === "") {
      errors.inviteCode = t("Common.This field is required");
    }
    if (sendTitle.trim() === "") {
      errors.title = t("Common.This field is required");
    }
    setFieldErrors(errors);

    if (Object.keys(errors).length > 0) {
      return;
    }

    setIsSending(true);
    setSendError(null);
    setSendSuccess(false);

    try {
      const targetProfile = await backend.getProfile(locale, targetSlug.trim());
      if (targetProfile === null) {
        setSendError(t("ProfileSettings.Profile not found"));
        setIsSending(false);
        return;
      }

      const result = await backend.sendProfileEnvelope({
        locale,
        senderSlug: sendFromSlug,
        targetProfileId: targetProfile.id,
        kind: "invitation",
        title: sendTitle.trim(),
        description: sendDescription.trim() !== "" ? sendDescription.trim() : undefined,
        inviteCode: envelopeKind === "telegram_group" ? inviteCode.trim() : undefined,
      });

      if (result !== null) {
        setTargetSlug("");
        setInviteCode("");
        setSendTitle("");
        setSendDescription("");
        setSendSuccess(true);
        setTimeout(() => setSendSuccess(false), 3000);
      } else {
        setSendError(t("ProfileSettings.Failed to send invitation"));
      }
    } catch (error) {
      setSendError(
        error instanceof Error ? error.message : t("ProfileSettings.Failed to send invitation"),
      );
    } finally {
      setIsSending(false);
    }
  };

  // Filtered envelopes
  const filteredEnvelopes = activeFilter === "all"
    ? envelopes
    : envelopes.filter((e) => e.owning_profile_slug === activeFilter);

  if (authLoading) {
    return (
      <PageLayout>
        <div className={styles.container}>
          <Skeleton className="h-8 w-48 mb-4" />
          <Skeleton className="h-4 w-72 mb-6" />
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center gap-3 p-4 border rounded-lg">
                <Skeleton className="size-10 rounded" />
                <div className="flex-1">
                  <Skeleton className="h-5 w-48 mb-2" />
                  <Skeleton className="h-4 w-24" />
                </div>
                <Skeleton className="h-5 w-20" />
              </div>
            ))}
          </div>
        </div>
      </PageLayout>
    );
  }

  if (!isAuthenticated || user === null) {
    return null;
  }

  return (
    <PageLayout>
      <div className={styles.container}>
        {/* Header */}
        <div className={styles.header}>
          <h2 className={styles.title}>{t("Layout.Mailbox")}</h2>
          <p className={styles.subtitle}>
            {t("Profile.Your inbox items and invitations.")}
          </p>
        </div>

        {/* Profile Filter Bar */}
        {maintainerProfiles.length > 1 && (
          <div className={styles.filterBar}>
            <button
              type="button"
              className={`${styles.filterBadge} ${activeFilter === "all" ? styles.filterBadgeActive : ""}`}
              onClick={() => setActiveFilter("all")}
            >
              {t("Common.All")}
            </button>
            {maintainerProfiles.map((profile) => (
              <button
                key={profile.slug}
                type="button"
                className={`${styles.filterBadge} ${activeFilter === profile.slug ? styles.filterBadgeActive : ""}`}
                onClick={() => setActiveFilter(profile.slug)}
              >
                {profile.title}
              </button>
            ))}
          </div>
        )}

        {/* Send In-mail Section — only when user has org/product profiles */}
        {sendableProfiles.length > 0 && (
          <Card className={styles.sendSection}>
            <button
              type="button"
              className={styles.sendToggle}
              onClick={() => setSendFormOpen(!sendFormOpen)}
            >
              <div className={styles.sendToggleLabel}>
                <Send className={styles.sendToggleIcon} />
                <h3 className={styles.sendToggleTitle}>
                  {t("ProfileSettings.Send In-mail")}
                </h3>
              </div>
              {sendFormOpen ? (
                <ChevronUp className={styles.sendToggleIcon} />
              ) : (
                <ChevronDown className={styles.sendToggleIcon} />
              )}
            </button>

            {sendFormOpen && (
              <div className={styles.sendForm}>
                <div className={styles.sendFormGrid}>
                  <Field>
                    <FieldLabel htmlFor="sendFromSlug">{t("Profile.From")}</FieldLabel>
                    <Select value={sendFromSlug} onValueChange={setSendFromSlug}>
                      <SelectTrigger id="sendFromSlug" className="w-full">
                        <SelectValue>
                          {sendableProfiles.find((p) => p.slug === sendFromSlug)?.title ?? sendFromSlug}
                        </SelectValue>
                      </SelectTrigger>
                      <SelectContent>
                        {sendableProfiles.map((profile) => (
                          <SelectItem key={profile.slug} value={profile.slug}>
                            {profile.title}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </Field>
                  <Field data-invalid={fieldErrors.targetSlug !== undefined && fieldErrors.targetSlug !== null}>
                    <FieldLabel htmlFor="targetSlug">{t("ProfileSettings.Target Profile Slug")}</FieldLabel>
                    <Input
                      id="targetSlug"
                      type="text"
                      placeholder="someone"
                      value={targetSlug}
                      onChange={(e) => {
                        setTargetSlug(e.target.value);
                        setFieldErrors((prev) => ({ ...prev, targetSlug: null }));
                      }}
                    />
                    {fieldErrors.targetSlug !== null && fieldErrors.targetSlug !== undefined && (
                      <FieldError>{fieldErrors.targetSlug}</FieldError>
                    )}
                  </Field>
                </div>
                <Field>
                  <FieldLabel htmlFor="envelopeKind">{t("ProfileSettings.Envelope Kind")}</FieldLabel>
                  <Select value={envelopeKind} onValueChange={setEnvelopeKind}>
                    <SelectTrigger id="envelopeKind" className="w-full">
                      <SelectValue>
                        {t(ENVELOPE_KINDS.find((k) => k.value === envelopeKind)?.labelKey ?? "")}
                      </SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                      {ENVELOPE_KINDS.map((kind) => (
                        <SelectItem key={kind.value} value={kind.value}>
                          {t(kind.labelKey)}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </Field>
                {envelopeKind === "telegram_group" && (
                  <Field data-invalid={fieldErrors.inviteCode !== undefined && fieldErrors.inviteCode !== null}>
                    <FieldLabel htmlFor="inviteCode">{t("ProfileSettings.Invite Code")}</FieldLabel>
                    <Input
                      id="inviteCode"
                      type="text"
                      placeholder="ABC123"
                      value={inviteCode}
                      onChange={(e) => {
                        setInviteCode(e.target.value.toUpperCase());
                        setFieldErrors((prev) => ({ ...prev, inviteCode: null }));
                      }}
                    />
                    <FieldDescription>
                      {t("ProfileSettings.Use /invite in a Telegram group to get a code")}
                    </FieldDescription>
                    {fieldErrors.inviteCode !== null && fieldErrors.inviteCode !== undefined && (
                      <FieldError>{fieldErrors.inviteCode}</FieldError>
                    )}
                  </Field>
                )}
                <Field data-invalid={fieldErrors.title !== undefined && fieldErrors.title !== null}>
                  <FieldLabel htmlFor="sendTitle">{t("Common.Title")}</FieldLabel>
                  <Input
                    id="sendTitle"
                    type="text"
                    placeholder={t("ProfileSettings.Title for the message")}
                    value={sendTitle}
                    onChange={(e) => {
                      setSendTitle(e.target.value);
                      setFieldErrors((prev) => ({ ...prev, title: null }));
                    }}
                  />
                  {fieldErrors.title !== null && fieldErrors.title !== undefined && (
                    <FieldError>{fieldErrors.title}</FieldError>
                  )}
                </Field>
                <Field>
                  <FieldLabel htmlFor="sendDescription">{t("Common.Description")}</FieldLabel>
                  <Textarea
                    id="sendDescription"
                    placeholder={t("ProfileSettings.Optional message text")}
                    value={sendDescription}
                    onChange={(e) => setSendDescription(e.target.value)}
                    rows={2}
                  />
                </Field>
                {sendError !== null && (
                  <FieldError>{sendError}</FieldError>
                )}
                {sendSuccess && (
                  <p className={styles.successMessage}>
                    <CheckCircle className="h-4 w-4" />
                    {t("ProfileSettings.Invitation sent successfully")}
                  </p>
                )}
                <div className="flex justify-end">
                  <Button onClick={handleSend} disabled={isSending}>
                    {isSending ? (
                      <Loader2 className="h-4 w-4 animate-spin mr-2" />
                    ) : (
                      <Send className="h-4 w-4 mr-2" />
                    )}
                    {t("ProfileSettings.Send In-mail")}
                  </Button>
                </div>
              </div>
            )}
          </Card>
        )}

        {/* Envelope List */}
        <Card className="p-6">
          {isLoading ? (
            <div className="space-y-2">
              {[1, 2, 3].map((i) => (
                <div key={i} className="flex items-center gap-3 p-4 border rounded-lg">
                  <Skeleton className="size-10 rounded" />
                  <div className="flex-1">
                    <Skeleton className="h-5 w-48 mb-2" />
                    <Skeleton className="h-4 w-24" />
                  </div>
                  <Skeleton className="h-5 w-20" />
                </div>
              ))}
            </div>
          ) : filteredEnvelopes.length === 0 ? (
            <div className={styles.emptyState}>
              <Mail className={styles.emptyIcon} />
              <p className={styles.emptyText}>{t("Profile.No inbox items yet.")}</p>
            </div>
          ) : (
            <div className="space-y-2">
              {filteredEnvelopes.map((envelope) => {
                const isPending = envelope.status === "pending";
                const isAccepted = envelope.status === "accepted";
                const isProcessing = actionInProgress === envelope.id;
                const groupName = envelope.properties !== null && envelope.properties !== undefined
                  ? (envelope.properties as Record<string, unknown>).group_name as string | undefined
                  : undefined;

                return (
                  <div key={envelope.id} className={styles.envelopeCard}>
                    <div className={`${styles.envelopeIcon} ${statusColors[envelope.status]}`}>
                      <Send className="size-5" />
                    </div>
                    <div className={styles.envelopeBody}>
                      <p className={styles.envelopeTitle}>{envelope.title}</p>
                      {groupName !== undefined && groupName !== "" && (
                        <p className={styles.envelopeGroup}>{groupName}</p>
                      )}
                      {envelope.description !== null && (
                        <p className={styles.envelopeDescription}>{envelope.description}</p>
                      )}
                      <div className={styles.envelopeMeta}>
                        <span>{formatDateString(envelope.created_at, locale)}</span>
                        <span className="mx-1">&middot;</span>
                        <span>{t("Profile.From")}:</span>
                        <SenderInfo envelope={envelope} locale={locale} />
                        {maintainerProfiles.length > 1 && (
                          <>
                            <span className="mx-1">&middot;</span>
                            <span>{t("Profile.To")}:</span>
                            <span className={styles.profileBadge}>
                              {envelope.owning_profile_title}
                            </span>
                          </>
                        )}
                      </div>
                    </div>
                    <div className={styles.envelopeActions}>
                      <Badge variant="outline" className={statusColors[envelope.status]}>
                        {t(`Profile.EnvelopeStatus.${envelope.status}`)}
                      </Badge>
                      {isPending && (
                        <>
                          <Button
                            size="sm"
                            variant="outline"
                            className="text-green-700 border-green-300 hover:bg-green-50"
                            disabled={isProcessing}
                            onClick={() => handleAccept(envelope.id, envelope.owning_profile_slug)}
                          >
                            <Check className="size-4 mr-1" />
                            {t("Profile.Accept")}
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            className="text-red-700 border-red-300 hover:bg-red-50"
                            disabled={isProcessing}
                            onClick={() => promptReject(envelope.id, envelope.owning_profile_slug)}
                          >
                            <X className="size-4 mr-1" />
                            {t("Profile.Reject")}
                          </Button>
                        </>
                      )}
                      {isAccepted && envelope.kind === "invitation" && (
                        <span className="text-xs text-muted-foreground">
                          {t("Profile.Use /invitations in bot to redeem")}
                        </span>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </Card>

        {/* Reject Confirmation Dialog */}
        <AlertDialog open={rejectDialogOpen} onOpenChange={setRejectDialogOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t("Profile.Reject Envelope")}</AlertDialogTitle>
              <AlertDialogDescription>
                {t("Profile.Are you sure you want to reject this? This action cannot be undone.")}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>{t("Common.Cancel")}</AlertDialogCancel>
              <AlertDialogAction
                variant="destructive"
                onClick={handleRejectConfirm}
              >
                {t("Profile.Reject")}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </PageLayout>
  );
}
