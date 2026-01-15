"use client";

import * as React from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/components/ui/card.tsx";
import { Button } from "@/shared/components/ui/button.tsx";
import { Alert, AlertDescription } from "@/shared/components/ui/alert.tsx";

type Status = "loading" | "success" | "error";

export default function AuthCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [status, setStatus] = React.useState<Status>("loading");
  const [errorMessage, setErrorMessage] = React.useState<string>("");
  const { t, localeCode } = useTranslations();

  React.useEffect(() => {
    const processCallback = async () => {
      try {
        const code = searchParams.get("code");
        const state = searchParams.get("state");

        console.log("Auth callback processing:", { code: code?.substring(0, 10) + "...", state, localeCode });

        if (code === null) {
          throw new Error("Authorization code not found");
        }

        // Use backend function to handle the callback
        console.log("Calling handleAuthCallback...");
        const authData = await backend.handleAuthCallback(localeCode, {
          code,
          state: state !== null ? state : undefined,
        });

        console.log("Auth callback successful, received data:", {
          hasToken: !!authData?.token,
          hasSession: !!authData?.session,
        });

        if (authData === null) {
          throw new Error("No authentication data received");
        }

        // Store JWT token in localStorage for now
        // TODO(@eser) Use secure httpOnly cookie in production
        if (globalThis.localStorage !== undefined) {
          localStorage.setItem("auth_token", authData.token);
          if (authData.session !== undefined) {
            localStorage.setItem("auth_session", JSON.stringify(authData.session));
          }
          // Store token expiration (24 hours from now)
          const expiresAt = Math.floor(Date.now() / 1000) + (24 * 60 * 60);
          localStorage.setItem("auth_token_expires_at", expiresAt.toString());
          console.log("Authentication data stored successfully");
        }

        setStatus("success");

        // Redirect to home or intended destination
        setTimeout(() => {
          router.push("/");
        }, 1000);
      } catch (error) {
        console.error("Auth callback error:", {
          error,
          message: error instanceof Error ? error.message : "Unknown error",
          stack: error instanceof Error ? error.stack : undefined,
          searchParams: Object.fromEntries(searchParams.entries()),
          currentUrl: globalThis.location?.href,
        });
        setStatus("error");
        setErrorMessage(error instanceof Error ? error.message : "Authentication failed");
      }
    };

    processCallback();
  }, [searchParams, router]);

  return (
    <div className="container mx-auto max-w-md px-4 py-8">
      <div className="flex min-h-[50vh] items-center justify-center">
        <Card className="w-full">
          <CardHeader className="text-center">
            <CardTitle>
              {status === "loading" && t("Auth", "Authenticating...")}
              {status === "success" && t("Auth", "✓ Successfully logged in!")}
              {status === "error" && t("Auth", "✗ Authentication failed")}
            </CardTitle>
            <CardDescription>
              {status === "loading" && t("Auth", "Please wait")}
              {status === "success" && t("Auth", "Redirecting to homepage...")}
              {status === "error" && "Please try again or return to homepage"}
            </CardDescription>
          </CardHeader>
          <CardContent className="text-center space-y-4">
            {status === "loading" && (
              <div className="flex justify-center">
                <div className="w-6 h-6 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin" />
              </div>
            )}

            {status === "success" && <div className="text-green-600 text-4xl">✓</div>}

            {status === "error" && (
              <>
                <Alert variant="destructive">
                  <AlertDescription>{errorMessage}</AlertDescription>
                </Alert>
                <Button
                  type="button"
                  onClick={() => router.push("/")}
                  variant="outline"
                  className="w-full"
                >
                  {t("Auth", "Return to Homepage")}
                </Button>
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
