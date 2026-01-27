// Profile points settings
import * as React from "react";
import { createFileRoute, getRouteApi } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  Coins,
  TrendingUp,
  TrendingDown,
  ArrowRightLeft,
} from "lucide-react";
import { backend, type ProfilePointTransaction, type ProfilePointTransactionType } from "@/modules/backend/backend";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { formatDateString } from "@/lib/date";

const settingsRoute = getRouteApi("/$locale/$slug/settings");

const transactionTypeIcons: Record<ProfilePointTransactionType, React.ElementType> = {
  GAIN: TrendingUp,
  SPEND: TrendingDown,
  TRANSFER: ArrowRightLeft,
};

const transactionTypeColors: Record<ProfilePointTransactionType, string> = {
  GAIN: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
  SPEND: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
  TRANSFER: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
};

export const Route = createFileRoute("/$locale/$slug/settings/points")({
  component: PointsSettingsPage,
});

function PointsSettingsPage() {
  const { t, i18n } = useTranslation();
  const locale = i18n.language;
  const params = Route.useParams();

  // Get profile from parent settings route loader
  const { profile } = settingsRoute.useLoaderData();

  const [transactions, setTransactions] = React.useState<ProfilePointTransaction[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);

  const currentBalance = profile?.points ?? 0;

  // Load transactions on mount
  React.useEffect(() => {
    loadTransactions();
  }, [params.locale, params.slug]);

  const loadTransactions = async () => {
    setIsLoading(true);
    const result = await backend.listProfilePointTransactions(params.locale, params.slug);
    if (result !== null) {
      setTransactions(result);
    }
    setIsLoading(false);
  };

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="mb-6">
          <Skeleton className="h-7 w-40 mb-2" />
          <Skeleton className="h-4 w-72" />
        </div>
        <div className="mb-6 p-4 rounded-lg bg-muted">
          <Skeleton className="h-4 w-24 mb-2" />
          <Skeleton className="h-8 w-32" />
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
              <Skeleton className="h-5 w-16" />
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
        <h3 className="font-serif text-xl font-semibold text-foreground">{t("Profile.Points")}</h3>
        <p className="text-muted-foreground text-sm mt-1">
          {t("Profile.View your point balance and transaction history.")}
        </p>
      </div>

      {/* Current Balance */}
      <div className="p-4 rounded-lg bg-muted">
        <p className="text-sm text-muted-foreground mb-1">{t("Profile.Current Balance")}</p>
        <div className="flex items-center gap-2">
          <Coins className="size-6 text-primary" />
          <span className="text-2xl font-bold">{currentBalance.toLocaleString()}</span>
          <span className="text-muted-foreground">{t("Profile.points")}</span>
        </div>
      </div>

      {/* Transaction History */}
      <div>
        <h4 className="font-medium text-foreground mb-3">{t("Profile.Transaction History")}</h4>

        {transactions.length === 0 ? (
          <div className="text-center py-12 border-2 border-dashed rounded-lg">
            <Coins className="size-12 mx-auto text-muted-foreground mb-4" />
            <p className="text-muted-foreground">{t("Profile.No transactions yet.")}</p>
          </div>
        ) : (
          <div className="space-y-2">
            {transactions.map((transaction) => {
              const TypeIcon = transactionTypeIcons[transaction.transaction_type];
              const typeColor = transactionTypeColors[transaction.transaction_type];
              const isPositive = transaction.transaction_type === "GAIN" ||
                (transaction.transaction_type === "TRANSFER" && transaction.amount > 0);

              return (
                <div
                  key={transaction.id}
                  className="flex items-center gap-3 p-4 border rounded-lg hover:bg-muted/50 transition-colors"
                >
                  <div className={`flex items-center justify-center size-10 rounded ${typeColor}`}>
                    {TypeIcon !== undefined ? <TypeIcon className="size-5" /> : <Coins className="size-5" />}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">{transaction.description}</p>
                    <p className="text-sm text-muted-foreground">
                      {formatDateString(transaction.created_at, locale)}
                    </p>
                  </div>
                  <div className="text-right">
                    <Badge variant="outline" className={typeColor}>
                      {t(`Profile.TransactionType.${transaction.transaction_type}`)}
                    </Badge>
                  </div>
                  <div className="text-right min-w-[80px]">
                    <p className={`font-medium ${isPositive ? "text-green-600" : "text-red-600"}`}>
                      {isPositive ? "+" : "-"}{Math.abs(transaction.amount).toLocaleString()}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {t("Profile.Balance")}: {transaction.balance_after.toLocaleString()}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </Card>
  );
}
