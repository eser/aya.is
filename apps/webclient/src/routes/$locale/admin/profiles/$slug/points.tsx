// Admin profile points tab
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import type { ProfilePointTransaction } from "@/modules/backend/types";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Coins, Plus, Loader2, CheckCircle } from "lucide-react";

export const Route = createFileRoute("/$locale/admin/profiles/$slug/points")({
  loader: async ({ params }) => {
    const { locale, slug } = params;
    const transactions = await backend.listProfilePointTransactions(locale, slug);
    return { transactions: transactions ?? [] };
  },
  component: AdminProfilePoints,
});

function AdminProfilePoints() {
  const { t } = useTranslation();
  const params = Route.useParams();
  const { transactions: initialTransactions } = Route.useLoaderData();

  // Get profile from parent route
  const parentData = Route.useRouteContext();
  // @ts-expect-error - accessing parent route data
  const profile = parentData.profile;

  const [transactions, setTransactions] = useState<ProfilePointTransaction[]>(initialTransactions);
  const [currentBalance, setCurrentBalance] = useState(profile?.points ?? 0);
  const [isAwarding, setIsAwarding] = useState(false);
  const [awardAmount, setAwardAmount] = useState("");
  const [awardDescription, setAwardDescription] = useState("");
  const [awardError, setAwardError] = useState<string | null>(null);
  const [awardSuccess, setAwardSuccess] = useState(false);

  if (profile === undefined) {
    return null;
  }

  const handleAwardPoints = async () => {
    const amount = Number.parseInt(awardAmount, 10);
    if (Number.isNaN(amount) || amount <= 0) {
      setAwardError(t("Admin.Please enter a valid amount"));
      return;
    }
    if (awardDescription.trim() === "") {
      setAwardError(t("Admin.Please enter a description"));
      return;
    }

    setIsAwarding(true);
    setAwardError(null);
    setAwardSuccess(false);

    try {
      const newTransaction = await backend.awardAdminPoints({
        slug: params.slug,
        amount,
        description: awardDescription.trim(),
      });

      if (newTransaction !== null) {
        // Add new transaction to the list
        setTransactions((prev) => [newTransaction, ...prev]);
        // Update balance
        setCurrentBalance(newTransaction.balance_after);
        // Clear form
        setAwardAmount("");
        setAwardDescription("");
        setAwardSuccess(true);
        // Hide success message after 3 seconds
        setTimeout(() => setAwardSuccess(false), 3000);
      } else {
        setAwardError(t("Admin.Failed to award points"));
      }
    } catch (error) {
      setAwardError(error instanceof Error ? error.message : t("Admin.Failed to award points"));
    } finally {
      setIsAwarding(false);
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

  const getTransactionTypeBadge = (type: string) => {
    switch (type) {
      case "GAIN":
        return <Badge className="bg-green-500">{t("Profile.Gain")}</Badge>;
      case "SPEND":
        return <Badge variant="destructive">{t("Profile.Spend")}</Badge>;
      case "TRANSFER":
        return <Badge variant="secondary">{t("Profile.Transfer")}</Badge>;
      default:
        return <Badge variant="outline">{type}</Badge>;
    }
  };

  return (
    <div className="space-y-6">
      {/* Current Balance Card */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Coins className="h-5 w-5" />
            {t("Admin.Current Balance")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-4xl font-bold">{currentBalance.toLocaleString()}</div>
          <p className="text-muted-foreground">{t("Admin.Total points")}</p>
        </CardContent>
      </Card>

      {/* Award Points Card */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Plus className="h-5 w-5" />
            {t("Admin.Award Points")}
          </CardTitle>
          <CardDescription>
            {t("Admin.Manually award points to this profile")}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="amount">{t("Admin.Amount")}</Label>
              <Input
                id="amount"
                type="number"
                min="1"
                placeholder="100"
                value={awardAmount}
                onChange={(e) => setAwardAmount(e.target.value)}
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">{t("Admin.Description")}</Label>
            <Textarea
              id="description"
              placeholder={t("Admin.Reason for awarding points...")}
              value={awardDescription}
              onChange={(e) => setAwardDescription(e.target.value)}
              rows={2}
            />
          </div>
          {awardError !== null && (
            <p className="text-sm text-destructive">{awardError}</p>
          )}
          {awardSuccess && (
            <p className="text-sm text-green-600 flex items-center gap-2">
              <CheckCircle className="h-4 w-4" />
              {t("Admin.Points awarded successfully")}
            </p>
          )}
          <Button onClick={handleAwardPoints} disabled={isAwarding}>
            {isAwarding ? (
              <Loader2 className="h-4 w-4 animate-spin mr-2" />
            ) : (
              <Plus className="h-4 w-4 mr-2" />
            )}
            {t("Admin.Award Points")}
          </Button>
        </CardContent>
      </Card>

      {/* Transaction History Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t("Admin.Transaction History")}</CardTitle>
          <CardDescription>
            {t("Admin.All point transactions for this profile")}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {transactions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              {t("Admin.No transactions found")}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t("Admin.Type")}</TableHead>
                  <TableHead>{t("Admin.Description")}</TableHead>
                  <TableHead className="text-right">{t("Admin.Amount")}</TableHead>
                  <TableHead className="text-right">{t("Admin.Balance After")}</TableHead>
                  <TableHead>{t("Admin.Date")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {transactions.map((tx) => (
                  <TableRow key={tx.id}>
                    <TableCell>{getTransactionTypeBadge(tx.transaction_type)}</TableCell>
                    <TableCell className="max-w-xs truncate">
                      {tx.description}
                      {tx.triggering_event !== null && (
                        <span className="ml-2 text-xs text-muted-foreground">
                          ({tx.triggering_event})
                        </span>
                      )}
                    </TableCell>
                    <TableCell className="text-right font-mono">
                      <span
                        className={
                          tx.transaction_type === "GAIN"
                            ? "text-green-600"
                            : tx.transaction_type === "SPEND"
                              ? "text-red-600"
                              : ""
                        }
                      >
                        {tx.transaction_type === "GAIN" ? "+" : "-"}
                        {tx.amount.toLocaleString()}
                      </span>
                    </TableCell>
                    <TableCell className="text-right font-mono">
                      {tx.balance_after.toLocaleString()}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(tx.created_at)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
