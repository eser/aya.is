// Admin points dashboard - statistics overview
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Clock, CheckCircle, XCircle, Coins } from "lucide-react";

export const Route = createFileRoute("/$locale/admin/points/")({
  loader: async () => {
    const stats = await backend.getPendingAwardsStats();
    return { stats };
  },
  component: AdminPointsDashboard,
});

function AdminPointsDashboard() {
  const { t } = useTranslation();
  const { stats } = Route.useLoaderData();

  if (stats === null) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        {t("Admin.Failed to load statistics")}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h2 className="font-serif text-xl font-bold">
        {t("Admin.Points Overview")}
      </h2>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("Admin.Pending")}
            </CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total_pending}</div>
            <p className="text-xs text-muted-foreground">
              {t("Admin.Awards awaiting review")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("Admin.Approved")}
            </CardTitle>
            <CheckCircle className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total_approved}</div>
            <p className="text-xs text-muted-foreground">
              {t("Admin.Total approved awards")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("Admin.Rejected")}
            </CardTitle>
            <XCircle className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total_rejected}</div>
            <p className="text-xs text-muted-foreground">
              {t("Admin.Total rejected awards")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("Admin.Points Awarded")}
            </CardTitle>
            <Coins className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {stats.points_awarded.toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground">
              {t("Admin.Total points distributed")}
            </p>
          </CardContent>
        </Card>
      </div>

      {Object.keys(stats.by_event_type).length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">
              {t("Admin.Pending by Event Type")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {Object.entries(stats.by_event_type).map(([event, count]) => (
                <div
                  key={event}
                  className="flex items-center justify-between py-1"
                >
                  <span className="text-sm text-muted-foreground">{event}</span>
                  <span className="font-medium">{count}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
