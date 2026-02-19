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
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { formatDateString } from "@/lib/date";

const settingsRoute = getRouteApi("/$locale/$slug/settings");

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
  const [targetSlug, setTargetSlug] = React.useState("");
  const [telegramChatId, setTelegramChatId] = React.useState("");
  const [groupName, setGroupName] = React.useState("");
  const [sendTitle, setSendTitle] = React.useState("");
  const [sendDescription, setSendDescription] = React.useState("");
  const [isSending, setIsSending] = React.useState(false);
  const [sendError, setSendError] = React.useState<string | null>(null);
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

  const handleGroupNameChange = (value: string) => {
    setGroupName(value);
    if (sendTitle === "" || sendTitle === `${groupName} Telegram Group`) {
      setSendTitle(value !== "" ? `${value} Telegram Group` : "");
    }
  };

  const handleSend = async () => {
    if (targetSlug.trim() === "") {
      setSendError(t("Admin.Target Profile Slug") + " is required");
      return;
    }
    if (telegramChatId.trim() === "") {
      setSendError(t("Admin.Telegram Chat ID") + " is required");
      return;
    }
    if (sendTitle.trim() === "") {
      setSendError(t("Common.Title") + " is required");
      return;
    }

    const chatIdNum = Number.parseInt(telegramChatId, 10);
    if (Number.isNaN(chatIdNum)) {
      setSendError(t("Admin.Telegram Chat ID") + " must be a number");
      return;
    }

    setIsSending(true);
    setSendError(null);
    setSendSuccess(false);

    try {
      const targetProfile = await backend.getProfile(params.locale, targetSlug.trim());
      if (targetProfile === null) {
        setSendError(t("Admin.Profile not found"));
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
        properties: {
          invitation_kind: "telegram_group",
          telegram_chat_id: chatIdNum,
          group_profile_slug: params.slug,
          group_name: groupName.trim(),
        },
      });

      if (result !== null) {
        setTargetSlug("");
        setTelegramChatId("");
        setGroupName("");
        setSendTitle("");
        setSendDescription("");
        setSendSuccess(true);
        setTimeout(() => setSendSuccess(false), 3000);
      } else {
        setSendError(t("Admin.Failed to send invitation"));
      }
    } catch (error) {
      setSendError(
        error instanceof Error ? error.message : t("Admin.Failed to send invitation"),
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
      {/* Send Invitation Section — only for org/product profiles */}
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
                {t("Profile.Send Invitation")}
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
                <div className="space-y-2">
                  <Label htmlFor="targetSlug">{t("Admin.Target Profile Slug")}</Label>
                  <Input
                    id="targetSlug"
                    type="text"
                    placeholder="seyma"
                    value={targetSlug}
                    onChange={(e) => setTargetSlug(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="telegramChatId">{t("Admin.Telegram Chat ID")}</Label>
                  <Input
                    id="telegramChatId"
                    type="text"
                    placeholder="-100123456789"
                    value={telegramChatId}
                    onChange={(e) => setTelegramChatId(e.target.value)}
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="groupName">{t("Admin.Group Name")}</Label>
                  <Input
                    id="groupName"
                    type="text"
                    placeholder="ajanstack dev"
                    value={groupName}
                    onChange={(e) => handleGroupNameChange(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="sendTitle">{t("Common.Title")}</Label>
                  <Input
                    id="sendTitle"
                    type="text"
                    placeholder="ajanstack dev Telegram Group"
                    value={sendTitle}
                    onChange={(e) => setSendTitle(e.target.value)}
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="sendDescription">{t("Common.Description")}</Label>
                <Textarea
                  id="sendDescription"
                  placeholder={t("Admin.Optional description for the invitation")}
                  value={sendDescription}
                  onChange={(e) => setSendDescription(e.target.value)}
                  rows={2}
                />
              </div>
              {sendError !== null && (
                <p className="text-sm text-destructive">{sendError}</p>
              )}
              {sendSuccess && (
                <p className="text-sm text-green-600 flex items-center gap-2">
                  <CheckCircle className="h-4 w-4" />
                  {t("Admin.Invitation sent successfully")}
                </p>
              )}
              <Button onClick={handleSend} disabled={isSending}>
                {isSending ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                ) : (
                  <Send className="h-4 w-4 mr-2" />
                )}
                {t("Profile.Send Invitation")}
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
