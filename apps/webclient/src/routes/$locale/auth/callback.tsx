// Auth callback page (locale-aware) - handles GitHub OAuth callback
import * as React from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { backend } from "@/modules/backend/backend";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";

type Status = "processing" | "error";

export const Route = createFileRoute("/$locale/auth/callback")({
  validateSearch: (search: Record<string, unknown>) => ({
    auth_token: (search.auth_token as string) ?? "",
    code: (search.code as string) ?? "",
    state: (search.state as string) ?? undefined,
    redirect: (search.redirect as string) ?? undefined,
  }),
  component: AuthCallbackPage,
});

function AuthCallbackPage() {
  const navigate = useNavigate();
  const { locale } = Route.useParams();
  const { auth_token, code, state, redirect } = Route.useSearch();
  const [status, setStatus] = React.useState<Status>("processing");
  const [errorMessage, setErrorMessage] = React.useState<string>("");
  const { t } = useTranslation();
  const processedRef = React.useRef(false);

  React.useEffect(() => {
    // Prevent double processing in React StrictMode
    if (processedRef.current) return;
    processedRef.current = true;

    const processCallback = async () => {
      try {
        // Case 1: Backend returned auth_token directly (new flow)
        if (auth_token !== "") {
          console.log("Auth callback processing with auth_token");

          // Parse JWT to get expiration
          let expiresAt = Date.now() + 24 * 60 * 60 * 1000; // Default 24 hours
          try {
            const payload = JSON.parse(atob(auth_token.split(".")[1]));
            if (payload.exp !== undefined) {
              expiresAt = payload.exp * 1000;
            }
          } catch {
            // Use default expiration if JWT parsing fails
          }

          // Store JWT token in localStorage
          if (typeof window !== "undefined" && globalThis.localStorage) {
            localStorage.setItem("auth_token", auth_token);
            localStorage.setItem("auth_token_expires_at", expiresAt.toString());
            console.log("Authentication token stored successfully");
          }

          // Redirect immediately to intended destination or home
          const destination = redirect !== undefined ? redirect : `/${locale}`;
          navigate({ to: destination });
          return;
        }

        // Case 2: Backend returned code to exchange (legacy flow)
        if (code !== "") {
          console.log("Auth callback processing with code:", {
            code: code.substring(0, 10) + "...",
            state,
            locale,
          });

          // Use backend function to handle the callback
          const authData = await backend.handleAuthCallback(locale, {
            code,
            state,
          });

          if (authData === null) {
            throw new Error("No authentication data received");
          }

          // Store JWT token in localStorage
          if (typeof window !== "undefined" && globalThis.localStorage) {
            localStorage.setItem("auth_token", authData.token);
            if (authData.session !== undefined) {
              localStorage.setItem(
                "auth_session",
                JSON.stringify(authData.session),
              );
            }
            const expiresAt = Date.now() + 24 * 60 * 60 * 1000;
            localStorage.setItem("auth_token_expires_at", expiresAt.toString());
            console.log("Authentication data stored successfully");
          }

          // Redirect immediately to intended destination or home
          const destination = redirect !== undefined ? redirect : `/${locale}`;
          navigate({ to: destination });
          return;
        }

        // Neither auth_token nor code found
        throw new Error("Authorization token or code not found");
      } catch (error) {
        console.error("Auth callback error:", {
          error,
          message: error instanceof Error ? error.message : "Unknown error",
          auth_token: auth_token !== "" ? "present" : "missing",
          code: code !== "" ? "present" : "missing",
          currentUrl: typeof window !== "undefined" ? globalThis.location.href : "",
        });
        setStatus("error");
        setErrorMessage(
          error instanceof Error ? error.message : "Authentication failed",
        );
      }
    };

    processCallback();
  }, [auth_token, code, state, redirect, locale, navigate]);

  const handleReturnHome = () => {
    navigate({ to: `/${locale}` });
  };

  // Only show UI if there's an error - otherwise we redirect immediately
  if (status === "processing") {
    return null;
  }

  return (
    <div className="container mx-auto max-w-md px-4 py-8">
      <div className="flex min-h-[50vh] items-center justify-center">
        <Card className="w-full">
          <CardHeader className="text-center">
            <CardTitle>{t("Auth.LoginFailed")}</CardTitle>
            <CardDescription>
              {t("Auth.TryAgainOrReturnToHomepage")}
            </CardDescription>
          </CardHeader>
          <CardContent className="text-center space-y-4">
            <Alert variant="destructive">
              <AlertDescription>{errorMessage}</AlertDescription>
            </Alert>
            <Button
              type="button"
              onClick={handleReturnHome}
              variant="outline"
              className="w-full"
            >
              {t("Auth.ReturnToHomepage")}
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
