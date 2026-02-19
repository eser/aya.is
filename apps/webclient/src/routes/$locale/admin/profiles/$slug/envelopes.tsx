// Admin profile envelopes tab — send invitations to profiles
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Send, Loader2, CheckCircle } from "lucide-react";

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
      setSendError(t("Admin.Target Profile Slug") + " is required");
      return;
    }
    if (inviteCode.trim() === "") {
      setSendError(t("Admin.Invite Code") + " is required");
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
      // Resolve target slug to profile ID
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
        title: title.trim(),
        description: description.trim() !== "" ? description.trim() : undefined,
        inviteCode: inviteCode.trim(),
      });

      if (result !== null) {
        setTargetSlug("");
        setInviteCode("");
        setTitle("");
        setDescription("");
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

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Send className="h-5 w-5" />
            {t("Admin.Send Telegram Invitation")}
          </CardTitle>
          <CardDescription>
            {t("Admin.Send an invitation to a profile")}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="targetSlug">{t("Admin.Target Profile Slug")}</Label>
              <Input
                id="targetSlug"
                type="text"
                placeholder="someone"
                value={targetSlug}
                onChange={(e) => setTargetSlug(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="inviteCode">{t("Admin.Invite Code")}</Label>
              <Input
                id="inviteCode"
                type="text"
                placeholder="ABC123"
                value={inviteCode}
                onChange={(e) => setInviteCode(e.target.value.toUpperCase())}
              />
              <p className="text-xs text-muted-foreground">
                {t("Admin.Use /invite in a Telegram group to get a code")}
              </p>
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="title">{t("Common.Title")}</Label>
              <Input
                id="title"
                type="text"
                placeholder="Telegram Group Invitation"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">{t("Common.Description")}</Label>
              <Textarea
                id="description"
                placeholder={t("Admin.Optional description for the invitation")}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={2}
              />
            </div>
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
            {t("Admin.Send Telegram Invitation")}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
