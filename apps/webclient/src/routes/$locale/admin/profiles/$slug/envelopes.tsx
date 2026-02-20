// Admin profile envelopes tab — send envelopes to profiles
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Field, FieldDescription, FieldError, FieldLabel } from "@/components/ui/field";
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
  const [fieldErrors, setFieldErrors] = useState<Record<string, string | null>>({});
  const [sendSuccess, setSendSuccess] = useState(false);

  if (profile === null || profile === undefined) {
    return null;
  }

  const handleSend = async () => {
    const errors: Record<string, string | null> = {};
    if (targetSlug.trim() === "") {
      errors.targetSlug = t("Common.This field is required");
    }
    if (envelopeKind === "telegram_group" && inviteCode.trim() === "") {
      errors.inviteCode = t("Common.This field is required");
    }
    if (title.trim() === "") {
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
            <FieldLabel htmlFor="title">{t("Common.Title")}</FieldLabel>
            <Input
              id="title"
              type="text"
              placeholder={t("ProfileSettings.Title for the message")}
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                setFieldErrors((prev) => ({ ...prev, title: null }));
              }}
            />
            {fieldErrors.title !== null && fieldErrors.title !== undefined && (
              <FieldError>{fieldErrors.title}</FieldError>
            )}
          </Field>
          <Field>
            <FieldLabel htmlFor="description">{t("Common.Description")}</FieldLabel>
            <Textarea
              id="description"
              placeholder={t("ProfileSettings.Optional message text")}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
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
        </CardContent>
      </Card>
    </div>
  );
}
