// Admin points dashboard - statistics overview and pending awards management
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { LocaleLink } from "@/components/locale-link";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import type { PendingAward, PendingAwardStatus } from "@/modules/backend/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Clock, CheckCircle, XCircle, Coins, Loader2 } from "lucide-react";

export const Route = createFileRoute("/$locale/admin/points/")({
  validateSearch: (search: Record<string, unknown>) => ({
    status: (search.status as PendingAwardStatus | undefined) ?? "pending",
  }),
  loaderDeps: ({ search: { status } }) => ({ status }),
  loader: async ({ deps: { status } }) => {
    const [stats, awardsResult] = await Promise.all([
      backend.getPendingAwardsStats(),
      backend.getPendingAwards({ status }),
    ]);
    return {
      stats,
      awards: awardsResult?.data ?? [],
      nextCursor: awardsResult?.next_cursor,
    };
  },
  component: AdminPointsDashboard,
});

function AdminPointsDashboard() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { status } = Route.useSearch();
  const { stats, awards: initialAwards } = Route.useLoaderData();
  const params = Route.useParams();

  const [awards, setAwards] = useState<PendingAward[]>(initialAwards);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [isProcessing, setIsProcessing] = useState(false);

  const handleStatusChange = (newStatus: string) => {
    navigate({
      to: `/${params.locale}/admin/points`,
      search: { status: newStatus as PendingAwardStatus },
    });
  };

  const toggleSelect = (id: string) => {
    setSelectedIds((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const toggleSelectAll = () => {
    if (selectedIds.size === awards.length) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(awards.map((a) => a.id)));
    }
  };

  const handleApprove = async (id: string) => {
    setIsProcessing(true);
    const result = await backend.approvePendingAward(id);
    if (result !== null) {
      setAwards((prev) => prev.filter((a) => a.id !== id));
      setSelectedIds((prev) => {
        const newSet = new Set(prev);
        newSet.delete(id);
        return newSet;
      });
    }
    setIsProcessing(false);
  };

  const handleReject = async (id: string) => {
    setIsProcessing(true);
    const result = await backend.rejectPendingAward(id);
    if (result !== null) {
      setAwards((prev) => prev.filter((a) => a.id !== id));
      setSelectedIds((prev) => {
        const newSet = new Set(prev);
        newSet.delete(id);
        return newSet;
      });
    }
    setIsProcessing(false);
  };

  const handleBulkApprove = async () => {
    if (selectedIds.size === 0) return;
    setIsProcessing(true);
    const result = await backend.bulkApprovePendingAwards(
      Array.from(selectedIds),
    );
    if (result !== null) {
      const approvedIds = new Set(
        result.results
          .filter((r) => r.status === "approved")
          .map((r) => r.id),
      );
      setAwards((prev) => prev.filter((a) => !approvedIds.has(a.id)));
      setSelectedIds(new Set());
    }
    setIsProcessing(false);
  };

  const handleBulkReject = async () => {
    if (selectedIds.size === 0) return;
    setIsProcessing(true);
    const result = await backend.bulkRejectPendingAwards(
      Array.from(selectedIds),
    );
    if (result !== null) {
      const rejectedIds = new Set(
        result.results
          .filter((r) => r.status === "rejected")
          .map((r) => r.id),
      );
      setAwards((prev) => prev.filter((a) => !rejectedIds.has(a.id)));
      setSelectedIds(new Set());
    }
    setIsProcessing(false);
  };

  const getStatusLabel = (statusValue: PendingAwardStatus) => {
    switch (statusValue) {
      case "pending":
        return t("Common.Pending");
      case "approved":
        return t("Common.Approved");
      case "rejected":
        return t("Common.Rejected");
      default:
        return statusValue;
    }
  };

  const getStatusBadge = (awardStatus: PendingAwardStatus) => {
    switch (awardStatus) {
      case "pending":
        return <Badge variant="secondary">{t("Common.Pending")}</Badge>;
      case "approved":
        return (
          <Badge variant="default" className="bg-green-500">
            {t("Common.Approved")}
          </Badge>
        );
      case "rejected":
        return <Badge variant="destructive">{t("Common.Rejected")}</Badge>;
      default:
        return null;
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <div className="space-y-6">
      {stats === null ? (
        <div className="text-center py-8 text-muted-foreground">
          {t("Admin.Failed to load statistics")}
        </div>
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {t("Admin.Pending Awards")}
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
                  {t("Admin.Approved Awards")}
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
                  {t("Admin.Rejected Awards")}
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
                      <span className="font-medium">{count as number}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </>
      )}

      {/* Pending Awards Section */}
      <div className="space-y-4 pt-4">
        <div className="flex items-center justify-between">
          <h2 className="font-serif text-xl font-bold">
            {t("Admin.Awards")}
          </h2>

          <Select value={status} onValueChange={handleStatusChange}>
            <SelectTrigger className="w-40">
              <SelectValue>{getStatusLabel(status)}</SelectValue>
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="pending">{t("Common.Pending")}</SelectItem>
              <SelectItem value="approved">{t("Common.Approved")}</SelectItem>
              <SelectItem value="rejected">{t("Common.Rejected")}</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {status === "pending" && selectedIds.size > 0 && (
          <div className="flex items-center justify-between p-3 bg-muted rounded-md">
            <span className="text-sm text-muted-foreground">
              {t("Common.Selected")}: {selectedIds.size}
            </span>
            <div className="flex gap-2">
              <Button
                size="sm"
                variant="outline"
                onClick={handleBulkApprove}
                disabled={isProcessing}
              >
                {isProcessing ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-1" />
                ) : (
                  <CheckCircle className="h-4 w-4 mr-1 text-green-500" />
                )}
                {t("Admin.Approve Selected")}
              </Button>
              <Button
                size="sm"
                variant="outline"
                onClick={handleBulkReject}
                disabled={isProcessing}
              >
                {isProcessing ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-1" />
                ) : (
                  <XCircle className="h-4 w-4 mr-1 text-red-500" />
                )}
                {t("Admin.Reject Selected")}
              </Button>
            </div>
          </div>
        )}

        {awards.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            {t("Admin.No awards found")}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                {status === "pending" && (
                  <TableHead className="w-12">
                    <Checkbox
                      checked={
                        selectedIds.size === awards.length && awards.length > 0
                      }
                      onCheckedChange={toggleSelectAll}
                    />
                  </TableHead>
                )}
                <TableHead>{t("Admin.Event")}</TableHead>
                <TableHead>{t("Admin.Profile")}</TableHead>
                <TableHead>{t("Common.Amount")}</TableHead>
                <TableHead>{t("Common.Date")}</TableHead>
                <TableHead>{t("Common.Status")}</TableHead>
                {status === "pending" && (
                  <TableHead className="text-right">{t("Common.Actions")}</TableHead>
                )}
              </TableRow>
            </TableHeader>
            <TableBody>
              {awards.map((award) => (
                <TableRow key={award.id}>
                  {status === "pending" && (
                    <TableCell>
                      <Checkbox
                        checked={selectedIds.has(award.id)}
                        onCheckedChange={() => toggleSelect(award.id)}
                      />
                    </TableCell>
                  )}
                  <TableCell className="font-mono text-sm">
                    {award.triggering_event}
                  </TableCell>
                  <TableCell className="text-sm">
                    {award.target_profile !== undefined ? (
                      <LocaleLink
                        to={`/${award.target_profile.slug}`}
                        className="text-foreground hover:underline"
                      >
                        {award.target_profile.title}
                      </LocaleLink>
                    ) : (
                      <span className="font-mono text-muted-foreground">
                        {award.target_profile_id.slice(0, 8)}...
                      </span>
                    )}
                  </TableCell>
                  <TableCell className="font-medium">{award.amount}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {formatDate(award.created_at)}
                  </TableCell>
                  <TableCell>{getStatusBadge(award.status)}</TableCell>
                  {status === "pending" && (
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-1">
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleApprove(award.id)}
                          disabled={isProcessing}
                        >
                          <CheckCircle className="h-4 w-4 text-green-500" />
                        </Button>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleReject(award.id)}
                          disabled={isProcessing}
                        >
                          <XCircle className="h-4 w-4 text-red-500" />
                        </Button>
                      </div>
                    </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  );
}
