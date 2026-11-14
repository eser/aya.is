// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import * as React from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { EyeOff, Loader2 } from "lucide-react";
import { backend } from "@/modules/backend/backend";
import { clearAuthData } from "@/modules/backend/fetcher";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

type HideSectionProps = {
  locale: string;
};

export function HideSection(props: Readonly<HideSectionProps>) {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = React.useState(false);
  const [isBusy, setIsBusy] = React.useState(false);

  const handleHide = async () => {
    setIsBusy(true);

    try {
      await backend.hidePrivacyProfile(props.locale);
      toast.success(t("Privacy.Your profile is now hidden."));

      // Every session was ended on the server; drop the local token and start over from the home page.
      clearAuthData();
      globalThis.location.assign(`/${props.locale}`);
    } catch {
      toast.error(t("Privacy.Could not complete the request. Please try again."));
      setIsBusy(false);
    }
  };

  return (
    <Card className="border-amber-500/50 bg-amber-50/40 dark:bg-amber-950/10">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <EyeOff className="size-5 text-amber-600" aria-hidden="true" />
          {t("Privacy.Hide your profile")}
        </CardTitle>
        <CardDescription>
          {t(
            "Privacy.Hide your profile from everyone until you log in again. You are signed out on all your devices, and your profile, stories and pages disappear from the site until your next login.",
          )}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button variant="outline" onClick={() => setIsOpen(true)}>
          {t("Privacy.Hide my profile")}
        </Button>

        <AlertDialog open={isOpen} onOpenChange={setIsOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>{t("Privacy.Hide your profile?")}</AlertDialogTitle>
              <AlertDialogDescription>
                {t("Privacy.You will be signed out now. Your profile becomes visible again the next time you log in.")}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel disabled={isBusy}>{t("Common.Cancel")}</AlertDialogCancel>
              <Button onClick={() => void handleHide()} disabled={isBusy}>
                {isBusy && <Loader2 className="mr-2 size-4 animate-spin" aria-hidden="true" />}
                {t("Privacy.Hide and sign out")}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </CardContent>
    </Card>
  );
}
