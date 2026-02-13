// Admin workers dashboard - background worker management
import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { formatDateTimeShort } from "@/lib/date";
import type { AdminWorkerStatus } from "@/modules/backend/backend";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Loader2, Play, Power } from "lucide-react";

export const Route = createFileRoute("/$locale/admin/workers/")({
  loader: async () => {
    const workers = await backend.getAdminWorkers();
    return { workers: workers ?? [] };
  },
  component: AdminWorkersDashboard,
});

function AdminWorkersDashboard() {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const { workers: initialWorkers } = Route.useLoaderData();

  const [workers, setWorkers] = useState<AdminWorkerStatus[]>(initialWorkers);
  const [processingWorker, setProcessingWorker] = useState<string | null>(null);

  // Auto-refresh every 30 seconds
  useEffect(() => {
    const interval = setInterval(async () => {
      const updated = await backend.getAdminWorkers();
      if (updated !== null) {
        setWorkers(updated);
      }
    }, 30000);

    return () => clearInterval(interval);
  }, []);

  const handleToggle = async (name: string) => {
    setProcessingWorker(name);
    const result = await backend.toggleAdminWorker(name);
    if (result !== null) {
      setWorkers((prev) =>
        prev.map((w) =>
          w.name === name ? { ...w, is_enabled: result.is_enabled } : w,
        ),
      );
    }
    setProcessingWorker(null);
  };

  const handleTrigger = async (name: string) => {
    setProcessingWorker(name);
    await backend.triggerAdminWorker(name);
    // Refresh after a short delay to get updated status
    setTimeout(async () => {
      const updated = await backend.getAdminWorkers();
      if (updated !== null) {
        setWorkers(updated);
      }
      setProcessingWorker(null);
    }, 1000);
  };

  const getStatusBadge = (worker: AdminWorkerStatus) => {
    if (!worker.is_enabled) {
      return <Badge variant="destructive">{t("Admin.Disabled")}</Badge>;
    }

    if (worker.is_running) {
      return (
        <Badge variant="default" className="bg-green-500">
          {t("Admin.Running")}
        </Badge>
      );
    }

    if (worker.last_error !== null) {
      return <Badge variant="secondary">{t("Admin.Idle")}</Badge>;
    }

    return <Badge variant="secondary">{t("Admin.Idle")}</Badge>;
  };

  const formatTime = (timeStr: string | null) => {
    if (timeStr === null) {
      return t("Admin.Never");
    }

    return formatDateTimeShort(new Date(timeStr), locale);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="font-serif text-xl font-bold">
          {t("Admin.Workers")}
        </h2>
        <Button
          variant="outline"
          size="sm"
          onClick={async () => {
            const updated = await backend.getAdminWorkers();
            if (updated !== null) {
              setWorkers(updated);
            }
          }}
        >
          {t("Common.Refresh")}
        </Button>
      </div>

      {workers.length === 0 ? (
        <div className="text-center py-8 text-muted-foreground">
          {t("Admin.No workers found")}
        </div>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t("Admin.Worker Name")}</TableHead>
              <TableHead>{t("Common.Status")}</TableHead>
              <TableHead>{t("Admin.Last Run")}</TableHead>
              <TableHead>{t("Admin.Next Run")}</TableHead>
              <TableHead>{t("Admin.Last Error")}</TableHead>
              <TableHead className="text-center">
                {t("Admin.Success Count")}
              </TableHead>
              <TableHead className="text-center">
                {t("Admin.Skip Count")}
              </TableHead>
              <TableHead className="text-center">
                {t("Admin.Error Count")}
              </TableHead>
              <TableHead className="text-right">
                {t("Common.Actions")}
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {workers.map((worker) => (
              <TableRow key={worker.name}>
                <TableCell className="font-mono text-sm">
                  {worker.name}
                </TableCell>
                <TableCell>{getStatusBadge(worker)}</TableCell>
                <TableCell className="text-sm text-muted-foreground">
                  {formatTime(worker.last_run)}
                </TableCell>
                <TableCell className="text-sm text-muted-foreground">
                  {formatTime(worker.next_run)}
                </TableCell>
                <TableCell className="text-sm text-muted-foreground max-w-48 truncate">
                  {worker.last_error ?? "-"}
                </TableCell>
                <TableCell className="text-center font-medium">
                  {worker.success_count > 0 ? (
                    <span className="text-green-600">{worker.success_count}</span>
                  ) : (
                    worker.success_count
                  )}
                </TableCell>
                <TableCell className="text-center font-medium text-muted-foreground">
                  {worker.skip_count}
                </TableCell>
                <TableCell className="text-center font-medium">
                  {worker.error_count > 0 ? (
                    <span className="text-red-500">{worker.error_count}</span>
                  ) : (
                    worker.error_count
                  )}
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-1">
                    <Button
                      size="sm"
                      variant={worker.is_enabled ? "outline" : "default"}
                      onClick={() => handleToggle(worker.name)}
                      disabled={processingWorker === worker.name}
                    >
                      {processingWorker === worker.name ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Power className="h-4 w-4" />
                      )}
                      <span className="ml-1">
                        {worker.is_enabled
                          ? t("Admin.Disable")
                          : t("Admin.Enable")}
                      </span>
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleTrigger(worker.name)}
                      disabled={
                        processingWorker === worker.name || !worker.is_enabled
                      }
                    >
                      {processingWorker === worker.name ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Play className="h-4 w-4" />
                      )}
                      <span className="ml-1">{t("Admin.Run Now")}</span>
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
  );
}
