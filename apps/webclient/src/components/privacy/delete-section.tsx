// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import * as React from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Loader2, Trash2 } from "lucide-react";
import { backend } from "@/modules/backend/backend";
import { clearAuthData } from "@/modules/backend/fetcher";
import { readPrivacyError } from "@/modules/backend/privacy/errors";
import type { PrivacyOwnedProfile } from "@/modules/backend/privacy/types";
import { isDeletionConfirmed } from "@/lib/privacy";
import { LocaleLink } from "@/components/locale-link";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";

type DeleteSectionProps = {
  locale: string;
  slug: string;
};

export function DeleteSection(props: Readonly<DeleteSectionProps>) {
  const { t } = useTranslation();
  const [typed, setTyped] = React.useState("");
  const [deleteContent, setDeleteContent] = React.useState(false);
  const [isOpen, setIsOpen] = React.useState(false);
  const [isBusy, setIsBusy] = React.useState(false);
  const [blockingProfiles, setBlockingProfiles] = React.useState<PrivacyOwnedProfile[]>([]);

  const confirmed = isDeletionConfirmed(typed, props.slug);

  const handleDelete = async () => {
    setIsBusy(true);

    try {
      await backend.deletePrivacyAccount(props.locale, { confirm_slug: typed, delete_content: deleteContent });
      toast.success(t("Privacy.Your account has been deleted."));

      clearAuthData();
      globalThis.location.assign(`/${props.locale}`);
    } catch (error) {
      const details = readPrivacyError(error);

      setIsOpen(false);
      setIsBusy(false);

      if (details.code === "owns_profiles") {
        setBlockingProfiles(details.profiles);
      } else if (details.code === "confirmation_mismatch") {
        toast.error(t("Privacy.The username you typed does not match."));
      } else {
        toast.error(t("Privacy.Could not complete the request. Please try again."));
      }
    }
  };

  return (
    <Card className="border-destructive/50">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-destructive">
          <Trash2 className="size-5" aria-hidden="true" />
          {t("Privacy.Delete your account")}
        </CardTitle>
        <CardDescription>
          {t("Privacy.Permanently delete your account and erase your personal data. This cannot be undone.")}{" "}
          {t(
            "Privacy.Your name, e-mail address and linked accounts are erased, and your username becomes available for anyone to register again.",
          )}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {blockingProfiles.length > 0 && (
          <Alert variant="destructive">
            <AlertDescription className="space-y-2">
              <p>{t("Privacy.Transfer ownership of these profiles before you delete your account.")}</p>
              <ul className="list-disc pl-5">
                {blockingProfiles.map((profile) => (
                  <li key={profile.id}>
                    <LocaleLink to={`/${profile.slug}/settings/access`} className="underline">
                      {profile.title !== "" ? profile.title : profile.slug}
                    </LocaleLink>
                  </li>
                ))}
              </ul>
            </AlertDescription>
          </Alert>
        )}

        <label className="flex items-start gap-3">
          <Checkbox
            checked={deleteContent}
            onCheckedChange={(checked) => setDeleteContent(checked === true)}
            className="mt-0.5"
          />
          <span className="space-y-1 text-sm">
            <span className="block font-medium">
              {t("Privacy.Also delete the stories, pages and comments I wrote")}
            </span>
            <span className="block text-muted-foreground">
              {t('Privacy.If unchecked, your content stays on the site under "Deleted user", without your name.')}
            </span>
          </span>
        </label>

        <div className="space-y-2">
          <label htmlFor="delete-account-confirm" className="block text-sm">
            {t("Privacy.Type your username ({{slug}}) to confirm", { slug: props.slug })}
          </label>
          <Input
            id="delete-account-confirm"
            value={typed}
            onChange={(event) => setTyped(event.target.value)}
            autoComplete="off"
            autoCapitalize="off"
            spellCheck={false}
          />
        </div>

        <Button variant="destructive" disabled={!confirmed} onClick={() => setIsOpen(true)}>
          {t("Privacy.Delete my account")}
        </Button>

        <AlertDialog open={isOpen} onOpenChange={setIsOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t("Privacy.Delete your account permanently?")}</AlertDialogTitle>
              <AlertDialogDescription>
                {t("Privacy.This is permanent. Your account and personal data cannot be restored afterwards.")}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel disabled={isBusy}>{t("Common.Cancel")}</AlertDialogCancel>
              <Button variant="destructive" onClick={() => void handleDelete()} disabled={isBusy || !confirmed}>
                {isBusy && <Loader2 className="mr-2 size-4 animate-spin" aria-hidden="true" />}
                {t("Privacy.Delete permanently")}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </CardContent>
    </Card>
  );
}
