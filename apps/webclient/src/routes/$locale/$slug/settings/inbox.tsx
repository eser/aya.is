// Profile inbox (envelope) settings
import * as React from "react";
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
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
import { backend, type ProfileEnvelope, type EnvelopeStatus } from "@/modules/backend/backend";
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
import { formatDateString } from "@/lib/date";

const settingsRoute = getRouteApi("/$locale/$slug/settings");

const ENVELOPE_KINDS = [
  { value: "telegram_group", labelKey: "ProfileSettings.Telegram Group Invite" },
] as const;

const statusColors: Record<EnvelopeStatus, string> = {
  pending: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200",
  accepted: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
  rejected: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
  revoked: "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200",
  redeemed: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
};

export const Route = createFileRoute("/$locale/$slug/settings/inbox")({
  component: InboxSettingsPage,
});

function InboxSettingsPage() {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const params = Route.useParams();

  const settingsData = settingsRoute.useLoaderData();
  const profile = settingsData.profile;
  const isNonIndividual = profile !== null && profile.kind !== "individual";

  const [envelopes, setEnvelopes] = React.useState<ProfileEnvelope[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [actionInProgress, setActionInProgress] = React.useState<string | null>(null);

  // Send form state
  const [sendFormOpen, setSendFormOpen] = React.useState(false);
  const [envelopeKind, setEnvelopeKind] = React.useState("telegram_group");
  const [targetSlug, setTargetSlug] = React.useState("");
  const [inviteCode, setInviteCode] = React.useState("");
  const [sendTitle, setSendTitle] = React.useState("");
  const [sendDescription, setSendDescription] = React.useState("");
  const [isSending, setIsSending] = React.useState(false);
  const [sendError, setSendError] = React.useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = React.useState<Record<string, string | null>>({});
  const [sendSuccess, setSendSuccess] = React.useState(false);

  React.useEffect(() => {
    loadEnvelopes();
  }, [params.locale, params.slug]);

  const loadEnvelopes = async () => {
    setIsLoading(true);
    const result = await backend.listProfileEnvelopes(params.locale, params.slug);
    if (result !== null) {
      setEnvelopes(result);
    }
    setIsLoading(false);
  };

  const handleAccept = async (envelopeId: string) => {
    setActionInProgress(envelopeId);
    const success = await backend.acceptProfileEnvelope(params.locale, params.slug, envelopeId);
    if (success) {
      await loadEnvelopes();
    }
    setActionInProgress(null);
  };

  const handleReject = async (envelopeId: string) => {
    setActionInProgress(envelopeId);
    const success = await backend.rejectProfileEnvelope(params.locale, params.slug, envelopeId);
    if (success) {
      await loadEnvelopes();
    }
    setActionInProgress(null);
  };

  const handleSend = async () => {
    const errors: Record<string, string | null> = {};
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
      const targetProfile = await backend.getProfile(params.locale, targetSlug.trim());
      if (targetProfile === null) {
        setSendError(t("ProfileSettings.Profile not found"));
        setIsSending(false);
        return;
      }

      const result = await backend.sendProfileEnvelope({
        locale: params.locale,
        senderSlug: params.slug,
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

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="mb-6">
          <Skeleton className="h-7 w-40 mb-2" />
          <Skeleton className="h-4 w-72" />
        </div>
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="flex items-center gap-3 p-4 border rounded-lg"
            >
              <Skeleton className="size-10 rounded" />
              <div className="flex-1">
                <Skeleton className="h-5 w-48 mb-2" />
                <Skeleton className="h-4 w-24" />
              </div>
              <Skeleton className="h-5 w-20" />
            </div>
          ))}
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Send In-mail Section — only for org/product profiles */}
      {isNonIndividual && (
        <Card className="p-6">
          <button
            type="button"
            className="flex items-center justify-between w-full text-left"
            onClick={() => setSendFormOpen(!sendFormOpen)}
          >
            <div className="flex items-center gap-2">
              <Send className="size-5 text-muted-foreground" />
              <h3 className="font-serif text-lg font-semibold text-foreground">
                {t("ProfileSettings.Send In-mail")}
              </h3>
            </div>
            {sendFormOpen ? (
              <ChevronUp className="size-5 text-muted-foreground" />
            ) : (
              <ChevronDown className="size-5 text-muted-foreground" />
            )}
          </button>

          {sendFormOpen && (
            <div className="mt-4 space-y-4">
              <div className="grid grid-cols-2 gap-4">
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
                <p className="text-sm text-green-600 flex items-center gap-2">
                  <CheckCircle className="h-4 w-4" />
                  {t("ProfileSettings.Invitation sent successfully")}
                </p>
              )}
              <Button onClick={handleSend} disabled={isSending}>
                {isSending ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                ) : (
                  <Send className="h-4 w-4 mr-2" />
                )}
                {t("ProfileSettings.Send In-mail")}
              </Button>
            </div>
          )}
        </Card>
      )}

      {/* Inbox List */}
      <Card className="p-6">
        <div>
          <h3 className="font-serif text-xl font-semibold text-foreground">{t("Common.Inbox")}</h3>
          <p className="text-muted-foreground text-sm mt-1">
            {t("Profile.Your inbox items and invitations.")}
          </p>
        </div>

        {envelopes.length === 0 ? (
          <div className="text-center py-12 border-2 border-dashed rounded-lg">
            <Mail className="size-12 mx-auto text-muted-foreground mb-4" />
            <p className="text-muted-foreground">{t("Profile.No inbox items yet.")}</p>
          </div>
        ) : (
          <div className="space-y-2">
            {envelopes.map((envelope) => {
              const isPending = envelope.status === "pending";
              const isAccepted = envelope.status === "accepted";
              const isProcessing = actionInProgress === envelope.id;
              const groupName = envelope.properties !== null && envelope.properties !== undefined
                ? (envelope.properties as Record<string, unknown>).group_name as string | undefined
                : undefined;

              return (
                <div
                  key={envelope.id}
                  className="flex items-center gap-3 p-4 border rounded-lg hover:bg-muted/50 transition-colors"
                >
                  <div className={`flex items-center justify-center size-10 rounded ${statusColors[envelope.status]}`}>
                    <Send className="size-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">{envelope.title}</p>
                    {groupName !== undefined && groupName !== "" && (
                      <p className="text-sm text-foreground/70 truncate">{groupName}</p>
                    )}
                    {envelope.description !== null && (
                      <p className="text-sm text-muted-foreground truncate">{envelope.description}</p>
                    )}
                    <p className="text-xs text-muted-foreground">
                      {formatDateString(envelope.created_at, locale)}
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
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
                          onClick={() => handleAccept(envelope.id)}
                        >
                          <Check className="size-4 mr-1" />
                          {t("Profile.Accept")}
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          className="text-red-700 border-red-300 hover:bg-red-50"
                          disabled={isProcessing}
                          onClick={() => handleReject(envelope.id)}
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
    </div>
  );
}
