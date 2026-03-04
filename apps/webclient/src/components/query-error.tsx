import { useQueryClient } from "@tanstack/react-query";
import { type ErrorComponentProps, useRouter } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { AlertCircle, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";

export function QueryError(props: Readonly<ErrorComponentProps>) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const router = useRouter();

  const handleRetry = () => {
    // Invalidate all queries to force refetch, then reload the route
    // and reset the error boundary so the component re-renders
    queryClient.invalidateQueries();
    router.invalidate();
    props.reset();
  };

  return (
    <div className="flex flex-col items-center justify-center gap-4 py-12">
      <AlertCircle className="size-10 text-destructive" />
      <h2 className="text-lg font-semibold">{t("Error.Something went wrong")}</h2>
      <p className="text-sm text-muted-foreground max-w-md text-center">
        {props.error.message}
      </p>
      <Button variant="outline" onClick={handleRetry}>
        <RefreshCw className="mr-2 size-4" />
        {t("Error.Try again")}
      </Button>
    </div>
  );
}
