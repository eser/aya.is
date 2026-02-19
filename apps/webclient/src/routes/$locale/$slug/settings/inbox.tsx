// Profile inbox (envelope) settings
import * as React from "react";
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  Mail,
  Check,
  X,
  Send,
} from "lucide-react";
import { backend, type ProfileEnvelope, type EnvelopeStatus } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
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

  const _settingsData = settingsRoute.useLoaderData();

  const [envelopes, setEnvelopes] = React.useState<ProfileEnvelope[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [actionInProgress, setActionInProgress] = React.useState<string | null>(null);

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
  );
}
