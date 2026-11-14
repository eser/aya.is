// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import * as React from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { Download, FileDown, Loader2 } from "lucide-react";
import { backend } from "@/modules/backend/backend";
import { readPrivacyError } from "@/modules/backend/privacy/errors";
import type { PrivacyExportState } from "@/modules/backend/privacy/types";
import { describeExport } from "@/lib/privacy";
import { formatDateTimeShort } from "@/lib/date";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

type ExportSectionProps = {
  locale: string;
  state: PrivacyExportState;
  onChanged: () => void;
};

export function ExportSection(props: Readonly<ExportSectionProps>) {
  const { t } = useTranslation();
  const [isBusy, setIsBusy] = React.useState(false);

  const model = describeExport(props.state);
  const formatDate = (iso: string) => formatDateTimeShort(new Date(iso), props.locale);

  const handleRequest = async () => {
    setIsBusy(true);

    try {
      await backend.requestPrivacyExport(props.locale);
      toast.success(t("Privacy.Request received. We will email you when your archive is ready."));
    } catch (error) {
      const details = readPrivacyError(error);

      if (details.code === "export_cooldown" && details.nextRequestAt !== null) {
        toast.info(
          t("Privacy.You can request a new copy after {{date}}.", { date: formatDate(details.nextRequestAt) }),
        );
      } else if (details.code === "export_unavailable") {
        toast.error(t("Privacy.Downloading your data is not available on this server."));
      } else if (details.code !== "export_in_progress") {
        toast.error(t("Privacy.Could not complete the request. Please try again."));
      }
    }

    setIsBusy(false);
    props.onChanged();
  };

  const handleDownload = async () => {
    if (props.state.id === null) {
      return;
    }

    setIsBusy(true);

    try {
      const result = await backend.getPrivacyExportDownloadUrl(props.locale, props.state.id);

      if (result === null) {
        toast.error(t("Privacy.Could not complete the request. Please try again."));
      } else {
        // The link is signed and forces a download (Content-Disposition: attachment), so the page stays put.
        globalThis.location.assign(result.url);
      }
    } catch {
      // The archive may have expired between the last refresh and the click.
      toast.error(t("Privacy.Your last archive has expired. Request a new one to get a fresh copy."));
      props.onChanged();
    }

    setIsBusy(false);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileDown className="size-5" aria-hidden="true" />
          {t("Privacy.Download your data")}
        </CardTitle>
        <CardDescription>
          {t(
            "Privacy.Request a copy of the personal data we store about you. We prepare it in the background and email you when it is ready to download.",
          )}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {model.view === "unavailable" && (
          <Alert>
            <AlertDescription>{t("Privacy.Downloading your data is not available on this server.")}</AlertDescription>
          </Alert>
        )}

        {model.view === "preparing" && (
          <Alert>
            <AlertDescription className="flex items-center gap-2">
              <Loader2 className="size-4 animate-spin" aria-hidden="true" />
              {t(
                "Privacy.Your archive is being prepared. This can take a few minutes. We will email you when it is ready.",
              )}
            </AlertDescription>
          </Alert>
        )}

        {model.view === "ready" && props.state.expires_at !== null && (
          <Alert>
            <AlertDescription>
              <span className="block font-medium">{t("Privacy.Your archive is ready")}</span>
              {t("Privacy.Available to download until {{date}}.", { date: formatDate(props.state.expires_at) })}
            </AlertDescription>
          </Alert>
        )}

        {model.view === "expired" && (
          <Alert>
            <AlertDescription>
              {t("Privacy.Your last archive has expired. Request a new one to get a fresh copy.")}
            </AlertDescription>
          </Alert>
        )}

        {model.view === "failed" && (
          <Alert variant="destructive">
            <AlertDescription>{t("Privacy.The last attempt failed. You can try again.")}</AlertDescription>
          </Alert>
        )}

        <div className="flex flex-wrap items-center gap-3">
          {model.canDownload && (
            <Button onClick={() => void handleDownload()} disabled={isBusy}>
              <Download className="mr-2 size-4" aria-hidden="true" />
              {t("Privacy.Download archive")}
            </Button>
          )}

          {model.canRequest && (
            <Button
              variant={model.canDownload ? "outline" : "default"}
              onClick={() => void handleRequest()}
              disabled={isBusy}
            >
              {isBusy && <Loader2 className="mr-2 size-4 animate-spin" aria-hidden="true" />}
              {t("Privacy.Request my data")}
            </Button>
          )}
        </div>

        <div className="space-y-1 text-sm text-muted-foreground">
          <p>
            {props.state.requested_at === null
              ? t("Privacy.You have not requested your data yet.")
              : t("Privacy.Last requested on {{date}}", { date: formatDate(props.state.requested_at) })}
          </p>

          {model.nextRequestAt !== null && (
            <p>{t("Privacy.You can request a new copy after {{date}}.", { date: formatDate(model.nextRequestAt) })}</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
