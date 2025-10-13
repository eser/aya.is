"use client";

import * as React from "react";
import { useAuth } from "./auth-context.tsx";
import { Alert, AlertDescription } from "@/shared/components/ui/alert.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

export function AuthErrorBanner() {
  const { authError, clearError } = useAuth();
  const { t } = useTranslations();

  if (authError === null) {
    return null;
  }

  return (
    <div className="fixed top-0 left-0 right-0 z-50 border-b bg-background">
      <div className="container mx-auto px-4 py-2">
        <Alert variant="destructive">
          <AlertDescription className="flex items-center justify-between">
            <span>{authError}</span>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={clearError}
              className="ml-4 h-auto p-1 text-destructive hover:text-primary-foreground"
            >
              {t("Layout", "Dismiss") || "Dismiss"}
            </Button>
          </AlertDescription>
        </Alert>
      </div>
    </div>
  );
}
