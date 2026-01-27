// Admin pending awards management page
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import type { PendingAward, PendingAwardStatus } from "@/modules/backend/types";
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
import { CheckCircle, XCircle, Loader2 } from "lucide-react";

export const Route = createFileRoute("/$locale/admin/points/pending")({
  validateSearch: (search: Record<string, unknown>) => ({
    status: (search.status as PendingAwardStatus | undefined) ?? "pending",
  }),
  loaderDeps: ({ search: { status } }) => ({ status }),
  loader: async ({ deps: { status } }) => {
    const result = await backend.getPendingAwards({ status });
    return { awards: result?.data ?? [], nextCursor: result?.next_cursor };
  },
  component: AdminPendingAwards,
});

function AdminPendingAwards() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { status } = Route.useSearch();
  const { awards: initialAwards } = Route.useLoaderData();
  const params = Route.useParams();

  const [awards, setAwards] = useState<PendingAward[]>(initialAwards);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [isProcessing, setIsProcessing] = useState(false);

  const handleStatusChange = (newStatus: string) => {
    navigate({
      to: `/${params.locale}/admin/points/pending`,
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

  const getStatusBadge = (awardStatus: PendingAwardStatus) => {
    switch (awardStatus) {
      case "pending":
        return <Badge variant="secondary">{t("Admin.Pending")}</Badge>;
      case "approved":
        return (
          <Badge variant="default" className="bg-green-500">
            {t("Admin.Approved")}
          </Badge>
        );
      case "rejected":
        return <Badge variant="destructive">{t("Admin.Rejected")}</Badge>;
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
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="font-serif text-xl font-bold">
          {t("Admin.Pending Awards")}
        </h2>

        <Select value={status} onValueChange={handleStatusChange}>
          <SelectTrigger className="w-40">
            <SelectValue placeholder={t("Admin.Filter by status")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="pending">{t("Admin.Pending")}</SelectItem>
            <SelectItem value="approved">{t("Admin.Approved")}</SelectItem>
            <SelectItem value="rejected">{t("Admin.Rejected")}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {status === "pending" && selectedIds.size > 0 && (
        <div className="flex items-center gap-2 p-3 bg-muted rounded-md">
          <span className="text-sm text-muted-foreground">
            {t("Admin.Selected")}: {selectedIds.size}
          </span>
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
              <TableHead>{t("Admin.Amount")}</TableHead>
              <TableHead>{t("Admin.Date")}</TableHead>
              <TableHead>{t("Admin.Status")}</TableHead>
              {status === "pending" && (
                <TableHead className="text-right">{t("Admin.Actions")}</TableHead>
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
                <TableCell className="font-mono text-sm text-muted-foreground">
                  {award.target_profile_id.slice(0, 8)}...
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
  );
}
