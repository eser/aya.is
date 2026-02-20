// Admin profile envelopes tab — send envelopes to profiles
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Send, Loader2, CheckCircle } from "lucide-react";

const ENVELOPE_KINDS = [
  { value: "telegram_group", labelKey: "ProfileSettings.Telegram Group Invite" },
] as const;

export const Route = createFileRoute("/$locale/admin/profiles/$slug/envelopes")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const profile = await backend.getAdminProfile(locale, slug);
    return { profile };
  },
  component: AdminProfileEnvelopes,
});

function AdminProfileEnvelopes() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const { profile } = Route.useLoaderData();

  const [envelopeKind, setEnvelopeKind] = useState("telegram_group");
  const [targetSlug, setTargetSlug] = useState("");
  const [inviteCode, setInviteCode] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [isSending, setIsSending] = useState(false);
  const [sendError, setSendError] = useState<string | null>(null);
  const [sendSuccess, setSendSuccess] = useState(false);

  if (profile === null || profile === undefined) {
    return null;
  }

  const handleSend = async () => {
    if (targetSlug.trim() === "") {
      setSendError(t("ProfileSettings.Target Profile Slug") + " is required");
      return;
    }
    if (envelopeKind === "telegram_group" && inviteCode.trim() === "") {
      setSendError(t("ProfileSettings.Invite Code") + " is required");
      return;
    }
    if (title.trim() === "") {
      setSendError(t("Common.Title") + " is required");
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
        title: title.trim(),
        description: description.trim() !== "" ? description.trim() : undefined,
        inviteCode: envelopeKind === "telegram_group" ? inviteCode.trim() : undefined,
      });

      if (result !== null) {
        setTargetSlug("");
        setInviteCode("");
        setTitle("");
        setDescription("");
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

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Send className="h-5 w-5" />
            {t("ProfileSettings.Send In-mail")}
          </CardTitle>
          <CardDescription>
            {t("ProfileSettings.Send an envelope to a profile")}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="envelopeKind">{t("ProfileSettings.Envelope Kind")}</Label>
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
            </div>
            <div className="space-y-2">
              <Label htmlFor="targetSlug">{t("ProfileSettings.Target Profile Slug")}</Label>
              <Input
                id="targetSlug"
                type="text"
                placeholder="someone"
                value={targetSlug}
                onChange={(e) => setTargetSlug(e.target.value)}
              />
            </div>
          </div>
          {envelopeKind === "telegram_group" && (
            <div className="space-y-2">
              <Label htmlFor="inviteCode">{t("ProfileSettings.Invite Code")}</Label>
              <Input
                id="inviteCode"
                type="text"
                placeholder="ABC123"
                value={inviteCode}
                onChange={(e) => setInviteCode(e.target.value.toUpperCase())}
              />
              <p className="text-xs text-muted-foreground">
                {t("ProfileSettings.Use /invite in a Telegram group to get a code")}
              </p>
            </div>
          )}
          <div className="space-y-2">
            <Label htmlFor="title">{t("Common.Title")}</Label>
            <Input
              id="title"
              type="text"
              placeholder={t("ProfileSettings.Title for the message")}
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">{t("Common.Description")}</Label>
            <Textarea
              id="description"
              placeholder={t("ProfileSettings.Optional message text")}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={2}
            />
          </div>
          {sendError !== null && (
            <p className="text-sm text-destructive">{sendError}</p>
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
        </CardContent>
      </Card>
    </div>
  );
}
