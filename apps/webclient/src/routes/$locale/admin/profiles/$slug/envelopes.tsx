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
  const [telegramChatId, setTelegramChatId] = useState("");
  const [groupName, setGroupName] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [isSending, setIsSending] = useState(false);
  const [sendError, setSendError] = useState<string | null>(null);
  const [sendSuccess, setSendSuccess] = useState(false);

  if (profile === null || profile === undefined) {
    return null;
  }

  const handleGroupNameChange = (value: string) => {
    setGroupName(value);
    // Auto-generate title if not manually edited
    if (title === "" || title === `${groupName} Telegram Group`) {
      setTitle(value !== "" ? `${value} Telegram Group` : "");
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
    if (title.trim() === "") {
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
              <Label htmlFor="title">{t("Common.Title")}</Label>
              <Input
                id="title"
                type="text"
                placeholder="ajanstack dev Telegram Group"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>
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
